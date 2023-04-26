package repo

import (
	"encoding/json"
	"fmt"
	"td_report/app/bean"
	"td_report/app/dao"
	"td_report/app/model"
	"td_report/boot"
	"td_report/common/redis"
	"time"

	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/guuid"
	"github.com/google/wire"
)

var ReportCheckRetryDetailRepositorySet = wire.NewSet(wire.Struct(new(ReportCheckRetryDetailRepository), "*"))

type ReportCheckRetryDetailRepository struct{}

func NewReportCheckRetryDetailRepository() *ReportCheckRetryDetailRepository {
	return &ReportCheckRetryDetailRepository{}
}

func (t *ReportCheckRetryDetailRepository) GetCountRedis(key string, reportType, reportName, reportDate, profileId, Extrant string) int {
	var filed string
	if Extrant != "" {
		filed = fmt.Sprintf("%s:%s:%s:%s:%s", reportType, reportName, reportDate, Extrant, profileId)
	} else {
		filed = fmt.Sprintf("%s:%s:%s:%s", reportType, reportName, reportDate, profileId)
	}

	str, err := boot.Redisclient.GetClient().HGet(key, filed).Result()
	if err != nil || str == "" {
		return 0
	}

	var temp *bean.CheckDetail
	json.Unmarshal([]byte(str), &temp)
	return temp.CheckTime
}

func (t *ReportCheckRetryDetailRepository) AddoneRedis(reportType string, reportName string, reportDate string, tempdataList []*bean.Tempdata) {
	key := redis.WithBatchCheck(reportType, reportName, reportDate)
	var filed string
	var exists = make(map[string]*bean.CheckDetail)
	redisclient := boot.Redisclient.GetClient()
	for _, item := range tempdataList {
		if item.Extrant != "" {
			filed = fmt.Sprintf("%s:%s:%s:%s:%s", reportType, reportName, reportDate, item.Extrant, item.ProfileId)
		} else {
			filed = fmt.Sprintf("%s:%s:%s:%s", reportType, reportName, reportDate, item.ProfileId)
		}

		str, err := redisclient.HGet(key, filed).Result()
		if err != nil {
			continue
		}

		if str != "" {
			var temp *bean.CheckDetail
			json.Unmarshal([]byte(str), &temp)
			temp.CheckTime = temp.CheckTime + 1
			temp.UpdateDate = gtime.Now()
			exists[filed] = temp
		}
	}

	pipe := boot.Redisclient.GetClient().Pipeline()
	var da *bean.CheckDetail
	var ok bool
	for _, item := range tempdataList {
		if item.Extrant != "" {
			filed = fmt.Sprintf("%s:%s:%s:%s:%s", reportType, reportName, reportDate, item.Extrant, item.ProfileId)
		} else {
			filed = fmt.Sprintf("%s:%s:%s:%s", reportType, reportName, reportDate, item.ProfileId)
		}

		da, ok = exists[filed]
		if !ok {
			da = &bean.CheckDetail{
				ReportType: reportType,
				ReportName: reportName,
				ReportDate: reportDate,
				ProfileId:  item.ProfileId,
				RetryType:  item.ErrorType,
				CreateDate: gtime.Now(),
				UpdateDate: gtime.Now(),
				CheckTime:  1,
				Extrant:    item.Extrant,
				FileUpdate: item.FileUpdate,
			}
		}

		bytedata, err := json.Marshal(da)
		if err != nil {
			return
		}

		pipe.HSet(key, filed, bytedata)

	}

	pipe.Expire(key, 24*time.Hour)
	pipe.Exec()
	//清空
	exists = make(map[string]*bean.CheckDetail)
}

func (t *ReportCheckRetryDetailRepository) Addone(reportType string, reportName string, reportDate string, tempdataList []*bean.Tempdata) {
	tempdata := make([]model.ReportCheckRetryDetail, 0)
	for _, v := range tempdataList {
		num, err := dao.ReportCheckRetryDetail.DB.Model(dao.ReportCheckRetryDetail.Table).
			Count("report_type=? and report_name=? and report_date=? and profile_id=? and extrant=?",
				reportType, reportName, reportDate, v.ProfileId, v.Extrant)
		if err != nil {
			continue
		} else {

			if num == 0 {
				uuidObj, _ := guuid.NewUUID()
				da := model.ReportCheckRetryDetail{
					Id:         uuidObj.String(),
					ReportType: reportType,
					ReportName: reportName,
					ReportDate: reportDate,
					ProfileId:  v.ProfileId,
					CheckTime:  1,
					CreateDate: gtime.Now(),
					UpdateDate: gtime.Now(),
					Extrant:    v.Extrant,
					RetryType:  v.ErrorType,
				}

				tempdata = append(tempdata, da)
			} else {

				dao.ReportCheckRetryDetail.DB.Model(dao.ReportCheckRetryDetail.Table).
					Where("report_type=? and report_name=? and report_date=? and profile_id=? and extrant=?", reportType, reportName, reportDate, v.ProfileId, v.Extrant).Increment("check_time", 1)
			}
		}

		if len(tempdata) >= 100 {
			dao.ReportCheckRetryDetail.DB.Model(dao.ReportCheckRetryDetail.Table).Data(tempdata).Insert()
			tempdata = make([]model.ReportCheckRetryDetail, 0)
			time.Sleep(3 * time.Second)
		}
	}

	if len(tempdata) > 0 {
		dao.ReportCheckRetryDetail.DB.Model(dao.ReportCheckRetryDetail.Table).Data(tempdata).Insert()
	}
}

func (t *ReportCheckRetryDetailRepository) GetCount(reportType, reportName, reportDate, profileId, Extrant string) int {
	num, _ := dao.ReportCheckRetryDetail.DB.Model(dao.ReportCheckRetryDetail.Table).
		Count("report_type=? and report_name=? and report_date=? and profile_id=? and extrant=?",
			reportType, reportName, reportDate, profileId, Extrant)
	return num
}
