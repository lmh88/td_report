package product_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/os/gtime"
	"strings"
	"td_report/app"
	"td_report/app/bean"
	"td_report/app/model"
	"td_report/app/repo"
	"td_report/app/service/amazon"
	"td_report/app/service/check_report_file"
	"td_report/app/service/report_v2"
	"td_report/app/service/report_v2/varible"
	"td_report/boot"
	datetool "td_report/common/tool"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

type ProductServer struct {
	ProfileServer *amazon.ProfileService
	ScheduleRepo *repo.ReportSchduleRepository
	SellerProfileRepo *repo.SellerProfileRepository
	ScheduleType int
}

func NewProductServer() *ProductServer {
	res := ProductServer{
		ProfileServer: app.InitializeProfileService(),
		ScheduleRepo: repo.NewReportSchduleRepository(),
		SellerProfileRepo: repo.NewSellerProfileRepository(),
		ScheduleType: 0,
	}
	return &res
}

// GetFilterProfile 过滤黑名单profile  目前暂定2天的需要过滤，其他的不需要过滤
func (t *ProductServer) GetFilterProfile(reportType string, date string) ([]string, error) {
	key := check_report_file.GetRdsKeyOfdate(reportType, date)
	return boot.RedisCommonClient.GetClient().SMembers(key).Result()
}

func (t *ProductServer) Product(reportType string, day int) {
	var (
		startDay, endDay time.Time
		queueName = varible.GetQueueByDay(reportType, day)
	)

	current := time.Now()
	startDay = current.Add(time.Duration(-1*day+1) * time.Hour * 24)
	if day == vars.ProductDay14 {
		// 去掉后7天的数据
		endDay = current.Add(time.Duration(-1*7) * time.Hour * 24)
	} else {
		endDay = current
	}

	if daylist, err := datetool.GetDaysWithTimeRever(startDay, endDay, vars.TimeLayout); err != nil {
		logger.Logger.Error(reportType, " product error ", err.Error())
	} else {
		fmt.Println(daylist)
		//过滤掉刚刚生成的某个日期的profile数据
		daylist = t.CheckRepeatProduct(reportType, day, daylist)
		fmt.Println(daylist)
		filterProfiles := make([]string, 0)
		if reportType != vars.DSP {
			if day == varible.SpecialDay {
				for _, reportDay := range daylist {
					if filterProfiles, err = t.GetFilterProfile(reportType, reportDay); err != nil {
						filterProfiles = make([]string, 0)
					}
					groupData := t.GetPpcDataByFilter(filterProfiles)
					t.AddPpcData(reportDay, reportType, queueName, groupData)
				}
			} else {
				t.ScheduleType = 1
				groupData := t.GetPpcDataByFilter(filterProfiles)
				for _, reportDay := range daylist {
					t.AddPpcData(reportDay, reportType, queueName, groupData)
				}
			}
		}
	}
}

func (t *ProductServer) CheckRepeatProduct(reportType string, num int, dayList []string) []string {
	//key := GetRepeatProductKey(reportType)
	rds := boot.RedisCommonClient.GetClient()
	list := make([]string, 0)
	for _, day := range dayList {
		key := varible.GetRepeatProductKey(reportType, day)
		if rds.Exists(key).Val() == 0 {
			if num != varible.SpecialDay {
				rds.Set(key, gtime.Timestamp(), time.Minute * 30)
			}
			list = append(list, day)
		}
	}
	return list
}

func (t *ProductServer) GetPpcDataByFilter(ProfileIdList []string) []*bean.ProfileToken {
	res, err := t.SellerProfileRepo.GetProfileAndRefreshTokenByFilter(ProfileIdList)
	if err != nil {
		logger.Logger.Error(map[string]interface{}{
			"flag": "ListProfile error",
			"err":  err.Error(),
		})
		return nil
	}
	return res
}

func (t *ProductServer) AddPpcData(reportDate, reportType, queueName string, profileList []*bean.ProfileToken) error {

	if len(profileList) > 0 {
		timestamp := gtime.Timestamp()
		batchKeys := make(map[string]int)
		mq := report_v2.NewMqServer()
		//TODO test
		//m := rand.Intn(len(profileList) - 20)
		//profileList = profileList[m:m+10]

		for _, item := range profileList {
			//if item.ProfileType == varible.Vendor || item.Tag == 2 {
			//
			//} else {
			//	continue
			//}
			clientTag := varible.GetClientTag(item.Tag)
			batchKey := varible.GetBatchKey(reportType, reportDate, timestamp, clientTag)
			batchKeys[batchKey]++
			profileMsg := varible.ProfileMsg{
				ReportType:   reportType,
				ProfileId:    item.ProfileId,
				ProfileType:  item.ProfileType,
				Region:       item.Region,
				RefreshToken: item.RefreshToken,
				Timestamp:    timestamp,
				ReportDate:   reportDate,
				BatchKey:     batchKey,
				ClientTag:    clientTag,
				ClientId:     item.ClientId,
				ClientSecret: item.ClientSecret,
			}

			msg, _ := json.Marshal(profileMsg)
			queue := queueName
			queue = varible.AddQueuePre(queue, clientTag)
			err := mq.SendMsg(varible.ReportDefaultExchange, queue, msg)
			if err != nil {
				logger.Logger.Error("profile推送失败：", err.Error(), "数据结构：", profileMsg)
			}
		}

		for key, num := range batchKeys {
			itemReportSchdule := &model.ReportSchdule{
				Batch:          key,
				ReportType:     reportType,
				ReportnameList: "all",
				ReportDate:     reportDate,
				ProfileNum:     num,
				SchduleType:    t.ScheduleType,
				CreateDate:     gtime.Now(),
				//EndTime: gtime.Now(),
			}
			t.ScheduleRepo.AddSchdule([]*model.ReportSchdule{itemReportSchdule})
		}
	}
	return nil
}

func (t *ProductServer) ProductWithDate(reportType, startDate, endDate string) error {

	var (
		queueName = varible.QueueMap[reportType][varible.SlowLevel]
	)

	if dayList, err := datetool.GetDays(startDate, endDate, vars.TimeLayout); err != nil {
		logger.Logger.Error(reportType, "product error", err.Error())
		return err
	} else {
		fmt.Println(dayList)
		dayList = t.CheckRepeatProduct(reportType, 1, dayList)
		fmt.Println(dayList)
		t.ScheduleType = 2
		profileList := make([]string, 0)
		if reportType != vars.DSP {
			groupData := t.GetPpcDataByFilter(profileList)
			if groupData == nil {
				logger.Logger.Error("ppc get profileList data error")
				return errors.New("ppc get profileList empty")
			}
			for _, reportDay := range dayList {
				t.AddPpcData(reportDay, reportType, queueName, groupData)
			}
		}
	}
	return nil
}

func (t *ProductServer) ProductOneProfile(reportType, startDate, endDate, profileId string) error {
	var (
		queueName = varible.QueueMap[reportType][varible.FastLevel]
	)

	if dayList, err := datetool.GetDays(startDate, endDate, vars.TimeLayout); err != nil {
		logger.Logger.Error(reportType, "product error", err.Error())
		return err
	} else {
		fmt.Println(dayList, profileId)
		t.ScheduleType = 2
		profileList := make([]string, 0)
		if strings.Contains(profileId, ",") {
			ids := strings.Split(profileId, ",")
			profileList = append(profileList, ids...)
		} else {
			profileList = append(profileList, profileId)
		}

		if reportType != vars.DSP {
			groupData, err := repo.NewSellerProfileRepository().GetProfileAndRefreshTokenById(profileList)
			fmt.Println(groupData)
			if groupData == nil || err != nil {
				logger.Logger.Error("ppc get profileList data error")
				return errors.New("ppc get profileList empty")
			}
			for _, reportDay := range dayList {
				t.AddPpcData(reportDay, reportType, queueName, groupData)
			}
		}
	}
	return nil
}

func (t *ProductServer) GetLastMonth() (string, string) {

	//r1, _ := time.Parse("2006-01-02", "2022-05-04")
	//year, month, day := r1.Date()

	year, month, day := time.Now().Date()
	thisM := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	lastM := thisM.AddDate(0, -1, 0)

	var startTime, endTime time.Time
	if day < 4 {
		start := (day - 1) * 5
		end := day * 5 - 1
		startTime = lastM.AddDate(0, 0, start)
		endTime = lastM.AddDate(0, 0, end)

	} else if day == 4 {
		start := (day - 1) * 5
		end := day * 5
		startTime = lastM.AddDate(0, 0, start)
		endTime = lastM.AddDate(0, 0, end)
	} else {
		return "", ""
	}

	//fmt.Println(startTime, endTime)
	return startTime.Format(vars.TimeLayout), endTime.Format(vars.TimeLayout)
}

