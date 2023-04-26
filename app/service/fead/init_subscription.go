package fead

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/guuid"
	"strconv"
	"td_report/app/bean"
	"td_report/app/model"
	"td_report/app/repo"
	"td_report/app/service/report_v2/create_report"
	"td_report/app/service/report_v2/varible"
	"td_report/pkg/logger"
	region2 "td_report/pkg/region"
	"td_report/pkg/requests"
	"time"
)

func InitSubscription(feadProfile string) {
	//添加订阅信息给有效的用户
	if feadProfile == "valid" {
		AddSubscription("")
	}
}

// AddSubscription 订阅用户
func AddSubscription(profile string) {
	profiles, err := repo.NewSellerProfileRepository().GetProfileAndRefreshTokenByFilter([]string{})
	fmt.Println(len(profiles), err)
	//profiles = profiles[350:]
	//fmt.Println(len(profiles))
	//return
	if len(profiles) > 0 {
		for _, item := range profiles {
			if item.Region == "NA" {
				if profile == "" {
					dealAddProfileV2(item)
				} else {
					if profile == item.ProfileId {
						dealAddProfileV2(item)
					}
				}
			}
		}
	}
}

// CancelSubscription 取消已经订阅的用户
func CancelSubscription(profile string) {
	profiles, _ := repo.NewSellerProfileRepository().GetDisableProfile()

	if len(profiles) > 0 {
		for _, item := range profiles {
			if item.Region == "NA" {
				if profile == "" {
					dealCancelProfile(item)
				} else {
					if profile == item.ProfileId {
						dealCancelProfile(item)
					}
				}
			}
		}
	}
}

func dealCancelProfile(profileToken *bean.ProfileToken) {

	feadSubs := repo.NewFeadSubscriptionRepository().GetByProfile(profileToken.ProfileId)
	if len(feadSubs) == 0 {
		return
	}

	for _, item := range feadSubs {
		if item.SqsStatus == StatusActive {
			err := DelSubscription(item, profileToken)
			if err != nil {
				logger.Logger.Error(err)
				continue
			}
			DelSubInDb(item)
		}
	}
}

func DelSubInDb(sub *model.FeadSubscription) {
	repo.NewFeadSubscriptionRepository().DelOne(sub.Id)
	tmp := model.FeadSubscription{
		ProfileId: sub.ProfileId,
		DatasetId: sub.DatasetId,
		MessageSubscriptionId: sub.MessageSubscriptionId,
		Status: 0,
		SqsName: sub.SqsName,
		SqsArn: sub.SqsArn,
		SqsStatus: StatusArchived,
		CreateDate: gtime.Now(),
		//IsDel: 1,
	}
	add := []*model.FeadSubscription{&tmp}
	repo.NewFeadSubscriptionRepository().AddBatch(add)
}

func DelSubscription(sub *model.FeadSubscription, profileToken *bean.ProfileToken) error {
	domain := region2.ApiUrl[profileToken.Region]
	url := domain + "/streams/subscriptions/" + sub.MessageSubscriptionId
	headers := getReqHeader(profileToken)

	params := map[string]interface{}{
		"status" : StatusArchived,
		"notes": profileToken.ProfileId,
	}

	resp , err := requests.Put(url, requests.WithHeaders(headers), requests.WithJson(params), requests.WithTimeout(time.Second * 30))

	if err != nil {
		logger.Logger.Error("CancelSubscription:" + err.Error())
		return err
	}

	logger.Logger.Info(string(resp.Body), resp.StatusCode)

	if resp.StatusCode != 200 {
		return errors.New(strconv.Itoa(resp.StatusCode) + "," + string(resp.Body))
	}

	return nil
}

func getMyData(subs []Subscription) []Subscription {
	res := make([]Subscription, 0)
	for _, item := range subs {
		if isMySqs(item.DestinationArn) {
			res = append(res, item)
		}
	}
	return res
}

func dealAddProfileV2(profileToken *bean.ProfileToken) {
	subs, err := GetSubscriptionByProfile(profileToken)
	if err != nil {
		logger.Logger.Error(err.Error(), "profile:" + profileToken.ProfileId)
		return
	}

	subs = getMyData(subs)
	fmt.Println(subs)

	var (
		createDb = make([]*model.FeadSubscription, 0)
		createSub = make(map[string]*bean.ProfileToken)
		delDb = make([]uint, 0)
		updateDb = make(map[uint]*model.FeadSubscription)
	)

	if len(subs) > 0 {
		for k, _ := range SqsMap {
			createSub[k] = profileToken
		}

		for _, item := range subs {
			if isValidStatus(item.Status) {
				delete(createSub, item.DataSetId)
			}
			one := repo.NewFeadSubscriptionRepository().GetBySubId(item.SubscriptionId)
			if one.Id == 0 {
				if isValidStatus(item.Status) {
					profileId, _ := strconv.Atoi(profileToken.ProfileId)
					tmp := model.FeadSubscription{
						ProfileId:  int64(profileId),
						DatasetId: item.DataSetId,
						MessageSubscriptionId: item.SubscriptionId,
						Status: 1,
						SqsName: item.DataSetId,
						SqsArn: SqsMap[item.DataSetId],
						SqsStatus: item.Status,
						CreateDate: gtime.Now(),
					}
					createDb = append(createDb, &tmp)
				}
			} else {
				if isInValidStatus(item.Status) {
					delDb = append(delDb, one.Id)
				}
				if middleStatus(one.SqsStatus) && item.Status != one.SqsStatus {
					profileId, _ := strconv.Atoi(profileToken.ProfileId)
					updateDb[one.Id] = &model.FeadSubscription{
						ProfileId:  int64(profileId),
						DatasetId: item.DataSetId,
						MessageSubscriptionId: item.SubscriptionId,
						Status: 1,
						SqsName: item.DataSetId,
						SqsArn: SqsMap[item.DataSetId],
						SqsStatus: item.Status,
						CreateDate: gtime.Now(),
					}
				}
			}

		}
	} else {
		for k, _ := range SqsMap {
			createSub[k] = profileToken
		}
	}

	if len(createDb) > 0 {
		repo.NewFeadSubscriptionRepository().AddBatch(createDb)
	}

	if len(createSub) > 0 {
		for k, v := range createSub {
			createSubAndDb(v, SqsMap[k], k)
		}
	}

	if len(delDb) > 0 {
		for _, item := range delDb {
			repo.NewFeadSubscriptionRepository().DelOne(item)
		}
	}

	if len(updateDb) > 0 {
		tmpDb := make([]*model.FeadSubscription, 0)
		for id, item := range updateDb {
			repo.NewFeadSubscriptionRepository().DelOne(id)
			tmpDb = append(tmpDb, item)
		}
		repo.NewFeadSubscriptionRepository().AddBatch(tmpDb)
	}

}

func createSubAndDb(profileToken *bean.ProfileToken, arn, dataSetId string) {
	subId, err := CreateSubscription(profileToken, arn, dataSetId)
	if err != nil {
		logger.Logger.Error(map[string]string{
			"active": "CreateSubscription",
			"err": err.Error(),
			"profile": profileToken.ProfileId,
			"region": profileToken.Region,
		})
		return
	}
	profileId, _ := strconv.Atoi(profileToken.ProfileId)
	tmp := model.FeadSubscription{
		ProfileId: int64(profileId),
		DatasetId: dataSetId,
		MessageSubscriptionId: subId,
		Status: 1,
		SqsName: dataSetId,
		SqsArn: arn,
		SqsStatus: StatusProvisioning,
		CreateDate: gtime.Now(),
	}
	createDb := make([]*model.FeadSubscription, 0)
	createDb = append(createDb, &tmp)
	repo.NewFeadSubscriptionRepository().AddBatch(createDb)
}

func dealAddProfile(profileToken *bean.ProfileToken) {
	fmt.Println(profileToken.ProfileId, profileToken.Region)
	needCreate := make(map[string]string)
	checkMap := make(map[uint]string)
	needDel := make([]uint, 0)
	hadSub := make(map[string]string)
	feadSubs := repo.NewFeadSubscriptionRepository().GetByProfile(profileToken.ProfileId)
	//fmt.Println(feadSubs)
	if len(feadSubs) > 0 {
		for _, item := range feadSubs {
			if _, ok := SqsMap[item.DatasetId]; ok {
				if isInValidStatus(item.SqsStatus) {
					needCreate[item.DatasetId] = SqsMap[item.DatasetId]
					needDel = append(needDel, item.Id)
				}
				if middleStatus(item.SqsStatus) {
					checkMap[item.Id] = item.MessageSubscriptionId
				}
				hadSub[item.DatasetId] = SqsMap[item.DatasetId]
			}
		}
		for key, item := range SqsMap {
			if _, ok := hadSub[key]; !ok {
				needCreate[key] = item
			}
		}
	} else {
		needCreate = SqsMap
	}

	if len(needCreate) > 0 {
		for id, arn := range needCreate {
			GetOrCreateSubscription(profileToken, arn, id)
		}
	}

	if len(checkMap) > 0 {
		for id, subId := range checkMap {
			UpdateSubscriptionStatus(profileToken, id, subId)
		}
	}

	if len(needDel) > 0 {
		for _, id := range needDel {
			repo.NewFeadSubscriptionRepository().DelOne(id)
		}
	}
	fmt.Println(needCreate, checkMap, needDel)
}

func GetOrCreateSubscription(profileToken *bean.ProfileToken, arn, dataSetId string) {

	var (
		subs = make([]Subscription, 0)
		err error
		createSubDB = make([]*model.FeadSubscription, 0)
	)

	subs, err =  GetSubscriptionByProfile(profileToken)

	if err != nil {
		logger.Logger.Error(map[string]string{
			"active": "GetSubscriptionByProfile",
			"err": err.Error(),
			"profile": profileToken.ProfileId,
			"region": profileToken.Region,
		})
		return
	}

	profileId, _ := strconv.Atoi(profileToken.ProfileId)

	if len(subs) > 0 {
		for _, item := range subs {
			if isValidStatus(item.Status) && item.DataSetId == dataSetId {
				tmp := model.FeadSubscription{
					ProfileId: int64(profileId),
					DatasetId: dataSetId,
					MessageSubscriptionId: item.SubscriptionId,
					Status: 1,
					SqsName: dataSetId,
					SqsArn: arn,
					SqsStatus: item.Status,
					CreateDate: gtime.Now(),
				}
				createSubDB = append(createSubDB, &tmp)
			}

			if isInValidStatus(item.Status) && item.DataSetId == dataSetId {

			}
		}
	}

	if len(createSubDB) == 0 {
		subId, err := CreateSubscription(profileToken, arn, dataSetId)
		if err != nil {
			logger.Logger.Error(map[string]string{
				"active": "CreateSubscription",
				"err": err.Error(),
				"profile": profileToken.ProfileId,
				"region": profileToken.Region,
			})
			return
		}
		tmp := model.FeadSubscription{
			ProfileId: int64(profileId),
			DatasetId: dataSetId,
			MessageSubscriptionId: subId,
			Status: 1,
			SqsName: dataSetId,
			SqsArn: arn,
			SqsStatus: StatusProvisioning,
			CreateDate: gtime.Now(),
		}
		createSubDB = append(createSubDB, &tmp)
	}

	err = repo.NewFeadSubscriptionRepository().AddBatch(createSubDB)
	if err != nil {
		logger.Logger.Error(err)
	}
}

func UpdateSubscriptionStatus(profileToken *bean.ProfileToken, id uint, subId string) {
	var (
		subs = make([]Subscription, 0)
		err error
	)

	subs, err =  GetSubscriptionByProfile(profileToken)

	if err != nil {
		logger.Logger.Error(map[string]string{
			"active": "UpdateSubscriptionStatus",
			"err": err.Error(),
			"profile": profileToken.ProfileId,
			"region": profileToken.Region,
		})
		return
	}

	if len(subs) > 0 {
		for _, item := range subs {
			if item.SubscriptionId == subId {
				var change g.Map
				if isInValidStatus(item.Status) {
					change = g.Map{
						"sqs_status": item.Status,
						"status": 0,
					}
				} else {
					change = g.Map{
						"sqs_status": item.Status,
					}
				}
				err = repo.NewFeadSubscriptionRepository().UpdateOne(id, change)
				if err != nil {
					logger.Logger.Error(err)
				}
			}
		}
	}
}

func getReqHeader(profileToken *bean.ProfileToken) map[string]string {
	crt := &varible.ClientRefreshToken{
		ClientId: profileToken.ClientId,
		ClientSecret: profileToken.ClientSecret,
		ProfileId: profileToken.ProfileId,
		RefreshToken: profileToken.RefreshToken,
	}
	headers, err := create_report.FakeHeaders(crt)

	if err != nil {
		logger.Logger.Error("getReqHeader_error:" + err.Error())
		return nil
	}
	headers["Content-Type"] = "application/vnd.MarketingStreamSubscriptions.StreamSubscriptionResource.v1.0+json"

	return headers
}

func GetSubscriptionByProfile(profileToken *bean.ProfileToken) ([]Subscription, error) {

	var subList SubscriptionList

	domain := region2.ApiUrl[profileToken.Region]
	url := domain + "/streams/subscriptions"
	headers := getReqHeader(profileToken)

	resp, err := requests.Get(url, requests.WithHeaders(headers), requests.WithTimeout(time.Second * 30))

	if err != nil {
		logger.Logger.Error(err.Error())
		return nil, err
	}

	//fmt.Println(string(resp.Body))
	logger.Logger.Info(string(resp.Body), resp.StatusCode)

	if resp.StatusCode != 200 {
		return nil, errors.New(strconv.Itoa(resp.StatusCode) + "," + string(resp.Body))
	}

	_, err = resp.Json(&subList)
	//fmt.Println(subList)
	return subList.Subscriptions, err
}

func CreateSubscription(profileToken *bean.ProfileToken, arn, dataSetId string) (string, error) {

	domain := region2.ApiUrl[profileToken.Region]
	url := domain + "/streams/subscriptions"
	headers := getReqHeader(profileToken)

	params := map[string]interface{}{
		"destinationArn" : arn,
		"dataSetId" : dataSetId,
		"clientRequestToken": guuid.New().String(),
		"notes": profileToken.ProfileId,
	}

	resp , err := requests.Post(url, requests.WithHeaders(headers), requests.WithJson(params), requests.WithTimeout(time.Second * 30))

	if err != nil {
		logger.Logger.Error("CreateSubscription:" + err.Error())
		return "", err
	}

	logger.Logger.Info(string(resp.Body), resp.StatusCode)

	if resp.StatusCode != 200 {
		return "", errors.New(strconv.Itoa(resp.StatusCode) + "," + string(resp.Body))
	}

	var subId SubscriptionId
	_, err = resp.Json(&subId)
	if err != nil {
		logger.Logger.Error("CreateSubscription:" + err.Error())
		return "", err
	}

	return subId.SubscriptionId, nil
}




