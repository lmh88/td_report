package report

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/os/gtime"
	"github.com/google/wire"
	"td_report/app/bean"
	"td_report/app/model"
	"td_report/app/repo"
	"td_report/app/service/amazon"
	"td_report/boot"
	"td_report/common/redis"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

// AdddataService 补数据服务
type AdddataService struct {
	profileService                   *amazon.ProfileService
	reportSchduleRepository          *repo.ReportSchduleRepository
	reportCheckRetryDetailRepository *repo.ReportCheckRetryDetailRepository
}

var AdddataServiceSet = wire.NewSet(wire.Struct(new(AdddataService), "*"))

func NewAdddataService(profileService *amazon.ProfileService,
	reportSchduleRepository *repo.ReportSchduleRepository,
	reportCheckRetryDetailRepository *repo.ReportCheckRetryDetailRepository) *AdddataService {
	return &AdddataService{profileService: profileService,
		reportSchduleRepository:          reportSchduleRepository,
		reportCheckRetryDetailRepository: reportCheckRetryDetailRepository,
	}
}

func (t *AdddataService) NewAddDspData(startdate string, ReportType string, schduletype int, groupData [][]*bean.DspRegionProfile, prefixQueue string) {
	var err error
	pipe := boot.RedisCommonClient.GetClient().Pipeline()
	schdulelist := make([]*model.ReportSchdule, 0)
	keysList := make([]string, 0)
	for _, v := range groupData {
		length := len(v)
		batchUuid := time.Now().Unix()
		batchList := make([]string, 0)
		key := redis.WithBasePrefix(fmt.Sprintf("new_batch_detail:%s:%s:%d_%d", ReportType, startdate, batchUuid, length))
		keysList = append(keysList, key)
		myitem := &model.ReportSchdule{
			Batch:          key,
			ReportType:     ReportType,
			ReportnameList: "all",
			ReportDate:     startdate,
			ProfileNum:     length,
			SchduleType:    schduletype,
			CreateDate:     gtime.Now(),
		}

		schdulelist = append(schdulelist, myitem)
		for _, item := range v {
			profiledata, _ := json.Marshal(item)
			batchList = append(batchList, string(profiledata))
		}

		pipe.LPush(key, batchList)
		pipe.Expire(key, 120*time.Hour)
		_, err = pipe.Exec()
		if err != nil {
			logger.Logger.Error(err.Error(), "==================redis create batch error")
		}
		//休眠，避免uuid
		time.Sleep(time.Second * 1)
	}

	pipe.Set(redis.WithNewSchdule(ReportType, prefixQueue), time.Now().Unix(), time.Hour*24)
	pipe.LPush(redis.WithNewBatch(ReportType, prefixQueue), keysList)
	pipe.Exec()

	if len(schdulelist) > 0 {
		t.reportSchduleRepository.AddSchdule(schdulelist)
	}
}

// NewAddPpcData   新模式生产者
func (t *AdddataService) NewAddPpcData(startdate string, ReportType string, schduletype int, groupData [][]*bean.ProfileToken, prefixQueue string) {
	var err error
	pipe := boot.RedisCommonClient.GetClient().Pipeline()
	schdulelist := make([]*model.ReportSchdule, 0)
	keysList := make([]string, 0)
	for _, v := range groupData {
		length := len(v)
		batchUuid := time.Now().Unix()
		batchList := make([]string, 0)
		key := redis.WithBasePrefix(fmt.Sprintf("new_batch_detail:%s:%s:%d_%d", ReportType, startdate, batchUuid, length))
		keysList = append(keysList, key)
		myitem := &model.ReportSchdule{
			Batch:          key,
			ReportType:     ReportType,
			ReportnameList: "all",
			ReportDate:     startdate,
			ProfileNum:     length,
			SchduleType:    schduletype,
			CreateDate:     gtime.Now(),
		}

		schdulelist = append(schdulelist, myitem)
		for _, item := range v {
			profiledata, _ := json.Marshal(item)
			batchList = append(batchList, string(profiledata))
		}

		pipe.LPush(key, batchList)
		pipe.Expire(key, 120*time.Hour) // 5天时间
		_, err = pipe.Exec()
		if err != nil {
			logger.Logger.Error(err.Error(), "==================redis create batch error")
		}
		//休眠，避免uuid
		time.Sleep(time.Second * 1)
	}

	pipe.Set(redis.WithNewSchdule(ReportType, prefixQueue), time.Now().Unix(), time.Hour*72)
	pipe.LPush(redis.WithNewBatch(ReportType, prefixQueue), keysList)
	pipe.Exec()

	if len(schdulelist) > 0 {
		t.reportSchduleRepository.AddSchdule(schdulelist)
	}
}

func (t *AdddataService) PpcAdgroupData(ProfileIdList []string) [][]*bean.ProfileToken {
	rps, err := t.profileService.GetPpcProfile(ProfileIdList)
	if err != nil {
		logger.Logger.Error(map[string]interface{}{
			"flag": "ListProfile error",
			"err":  err.Error(),
		})

		return nil
	}

	n := 100 // 数据切割成100份
	return t.profileService.SellerProfileRepository.ArrayInGroupsOf(rps, int64(n))
}

// PpcAdgroupDataFilter 过滤数据
func (t *AdddataService) PpcAdgroupDataFilter(ProfileIdList []string) [][]*bean.ProfileToken {
	rps, err := t.profileService.GetPpcProfileFilter(ProfileIdList)
	if err != nil {
		logger.Logger.Error(map[string]interface{}{
			"flag": "ListProfile error",
			"err":  err.Error(),
		})

		return nil
	}

	n := 100 // 数据切割成100份
	return t.profileService.SellerProfileRepository.ArrayInGroupsOf(rps, int64(n))
}

func (t *AdddataService) DspAdgroupData(ProfileIdList []string) [][]*bean.DspRegionProfile {
	rps, err := t.profileService.GetDspProfile(ProfileIdList)
	if err != nil {
		logger.Logger.Error(map[string]interface{}{
			"flag":       "ListdspProfile error",
			"err":        err.Error(),
			"reporttype": vars.DSP,
		})

		return nil
	}

	n := 100 // 数据切割成100份
	return t.profileService.ProfileRepository.ArrayInGroupsOf(rps, int64(n))
}

func (t *AdddataService) DspAdgroupDataFilter(ProfileIdList []string) [][]*bean.DspRegionProfile {
	rps, err := t.profileService.GetDspProfileFilter(ProfileIdList)
	if err != nil {
		logger.Logger.Error(map[string]interface{}{
			"flag":       "ListdspProfile error",
			"err":        err.Error(),
			"reporttype": vars.DSP,
		})

		return nil
	}

	n := 100 // 数据切割成100份
	return t.profileService.ProfileRepository.ArrayInGroupsOf(rps, int64(n))
}

//// AddPpcData 添加带补充的数据到队列中 dataDirection 1 lpush 2 rpush(插队)
//func (t *AdddataService) AddPpcData(startdate string, ReportType string, ReportNameList []string, groupData [][]*bean.ProfileToken, schduletype int, prefixQueue string) {
//	var err error
//	pipe := boot.RedisCommonClient.GetClient().Pipeline()
//	datalist := make([]*model.ReportSchdule, 0)
//	for _, reportname := range ReportNameList {
//		keysList := make([]string, 0)
//		for _, v := range groupData {
//			length := len(v)
//			batchUuid := time.Now().Unix()
//			batchList := make([]string, 0)
//			key := redis.WithBasePrefix(fmt.Sprintf("batch_detail:%s:%d:%d:%s:%s", ReportType, batchUuid, length, reportname, startdate))
//			keysList = append(keysList, key)
//			for _, item := range v {
//				batch := reportname + "," + startdate + "," + item.ProfileId + "," + item.Region
//				batchList = append(batchList, batch)
//			}
//
//			myitem := &model.ReportSchdule{
//				Batch:          key,
//				ReportType:     ReportType,
//				ReportnameList: reportname,
//				ReportDate:     startdate,
//				ProfileNum:     length,
//				SchduleType:    schduletype,
//				CreateDate:     gtime.Now(),
//			}
//
//			datalist = append(datalist, myitem)
//			pipe.LPush(key, batchList)
//			pipe.Expire(key, 8*time.Hour)
//			_, err = pipe.Exec()
//			if err != nil {
//				logger.Logger.Error(err.Error(), "==================redis create batch error")
//			}
//			//休眠，避免uuid
//			time.Sleep(time.Second * 1)
//		}
//
//		pipe.Set(redis.WithSchdule(ReportType, reportname, startdate), time.Now().Unix(), time.Hour*24)
//		pipe.LPush(redis.WithBatch(ReportType), keysList)
//		pipe.Exec()
//	}
//
//	if len(datalist) > 0 {
//		t.reportSchduleRepository.AddSchdule(datalist)
//	}
//}
//
//func (t *AdddataService) AddDspData(startdate string, ReportNameList []string, groupData [][]*bean.DspRegionProfile, schduletype int, dataDirection int) {
//	var err error
//	pipe := boot.RedisCommonClient.GetClient().Pipeline()
//	datalist := make([]*model.ReportSchdule, 0)
//	for _, reportname := range ReportNameList {
//		keysList := make([]string, 0)
//		for _, v := range groupData {
//			length := len(v)
//			batchUuid := time.Now().Unix()
//			batchList := make([]string, 0)
//			key := redis.WithBasePrefix(fmt.Sprintf("batch_detail:%s:%d:%d:%s:%s", vars.DSP, batchUuid, length, reportname, startdate))
//			keysList = append(keysList, key)
//			for _, item := range v {
//				batch := reportname + "," + startdate + "," + item.ProfileId + "," + item.Region
//				batchList = append(batchList, batch)
//			}
//
//			myitem := &model.ReportSchdule{
//				Batch:          key,
//				ReportType:     vars.DSP,
//				ReportnameList: reportname,
//				ReportDate:     startdate,
//				ProfileNum:     length,
//				SchduleType:    schduletype,
//				CreateDate:     gtime.Now(),
//			}
//
//			datalist = append(datalist, myitem)
//			pipe.LPush(key, batchList)
//			pipe.Expire(key, 8*time.Hour)
//			_, err = pipe.Exec()
//			if err != nil {
//				logger.Logger.Error(err.Error(), "==================redis create batch error")
//				continue
//			}
//			//休眠，避免uuid
//			time.Sleep(time.Second * 1)
//		}
//
//		pipe.Set(redis.WithSchdule(vars.DSP, reportname, startdate), time.Now().Unix(), time.Hour*24)
//		// 去掉插队机制，防止数据太多当天的数据一直无法消费
//		pipe.LPush(redis.WithBatch(vars.DSP), keysList)
//		pipe.Exec()
//	}
//
//	if len(datalist) > 0 {
//		t.reportSchduleRepository.AddSchdule(datalist)
//	}
//
//}

// PullData 需要重新拉取的记录到数据库
func (t *AdddataService) PullData(startdate string, ReportType string, ReportName string, tempdataList []*bean.Tempdata) {
	t.reportCheckRetryDetailRepository.AddoneRedis(ReportType, ReportName, startdate, tempdataList)
}

func (t *AdddataService) GetCount(key, reportType, reportName, reportDate, profileId, Extrant string) int {
	return t.reportCheckRetryDetailRepository.GetCountRedis(key, reportType, reportName, reportDate, profileId, Extrant)
}

// CheckNum 检测最近一次创建的批次间隔时间和当前时间的差值，如果小于则说明短时间内创建重复了
func (t *AdddataService) CheckNum(date, reportType, reportName string, gap int64) bool {
	//schduleDate, err := boot.RedisCommonClient.GetClient().Get(redis.WithNewSchdule(reportType, reportName, date)).Int64()
	//if err != nil {
	//	return false
	//}
	//current := time.Now().Unix()
	//if current-schduleDate < gap {
	//	return true
	//}
	//return false
	return false
}
