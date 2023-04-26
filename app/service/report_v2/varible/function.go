package varible

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"td_report/app/bean"
	"td_report/common/stringx"
	"td_report/vars"
)

func getQueueName(reportType, level string) string {
	name, ok := QueueMap[reportType][level]
	if ok {
		return name
	}
	return ""
}

func GetQueueByDay(reportType string, day int) string {
	if day == vars.ProductDay2 {
		return getQueueName(reportType, "fast")
	} else if day == vars.ProductDay7 {
		return getQueueName(reportType, "middle")
	} else if day == vars.ProductDay14 {
		return getQueueName(reportType, "slow")
	}
	//todo 测试
	return getQueueName(reportType, "slow")
}

func AddQueuePre(queueName, clientTag string) string {
	if clientTag != "" {
		return clientTag + ":" + queueName
	}
	return queueName
}

func GetBatchKey(reportType, reportDate string, timestamp int64, clientTag string) string {
	return fmt.Sprintf("report:%s:%s:%s:%d", clientTag, reportType, reportDate, timestamp)
}

func GetClientTag(tag int) string {
	return fmt.Sprintf("c%d", tag)
}

func GetTraceId(profileMsg *ProfileMsg, reportName, tactic string) string {
	if tactic == "" {
		return fmt.Sprintf("%s_%s_%s_%s_%d_%s", profileMsg.ReportType, profileMsg.ProfileId, profileMsg.ReportDate, reportName, profileMsg.Timestamp, stringx.Randn(4))
	}
	return fmt.Sprintf("%s_%s_%s_%s_%s_%d_%s", profileMsg.ReportType, profileMsg.ProfileId, profileMsg.ReportDate, reportName, tactic, profileMsg.Timestamp, stringx.Randn(4))
}

func CheckQueueLevel(reportType, level string) bool {
	_, ok := QueueMap[reportType][level]
	return ok
}

func GetBatchTimeKey(rKey string) string {
	return fmt.Sprintf("batch_start_time:%s", rKey)
}

func GetRepeatProductKey(reportType, day string) string {
	return fmt.Sprintf("product_data:%s:%s:%s", reportType, gtime.Date(), day)
}

func GetQueueBusyKey(reportType, clientTag, kind string) string {
	return fmt.Sprintf("queue_busy:%s:%s:%s:%s", kind, clientTag, reportType, gtime.Date())
}

func GetRetryCount(reportType string)  map[string]int {
	key := "consumer_wait." + reportType
	res := RetryCount
	res[WaitTime120s] = g.Cfg().GetInt(key, 20)
	return res
}

func GetNewProfileKey(reportData *bean.ReportData, reportName, profile string) string {
	date := gtime.Date()
	if reportName != "" {
		return fmt.Sprintf("new_profile:%s:%s_%s_%s_%s_%s", date, reportData.ReportType, profile, reportName, reportData.StartDate, reportData.EndDate)
	} else {
		return fmt.Sprintf("new_profile:%s:%s_%s_%s_%s", date, reportData.ReportType, profile, reportData.StartDate, reportData.EndDate)
	}
}

