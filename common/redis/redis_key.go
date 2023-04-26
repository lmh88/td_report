package redis

import (
	"fmt"
	"time"
)

func WithBasePrefix(key string) string {
	return fmt.Sprintf("dp:%s", key)
}

func WithEventDetailPrefix(key string) string {
	return WithBasePrefix(fmt.Sprintf("ed:%s", key))
}

func WithEventListKey() string {
	return WithBasePrefix("el:key")
}

func WithUidListKey() string {
	return WithBasePrefix("uidl:key")
}

func WithUidEventListPrefix(key string) string {
	return WithBasePrefix(fmt.Sprintf("uel:%s", key))
}

func WithTimerReportKeyPrefix(key, reportDate string) string {
	todayStr := time.Now().Format("2006-01-02")
	return WithBasePrefix(fmt.Sprintf("trk:%s:%s:%s", key, todayStr, reportDate))
}

func WithTimerReportRunningPrefix(key string) string {
	return WithBasePrefix(fmt.Sprintf("trr:lc:%s", key))
}

func WithTimerReportCountPrefix(key string) string {
	return WithBasePrefix(fmt.Sprintf("trc:%s", key))
}

func WithSyncLockPrefix(key string) string {
	return WithBasePrefix(fmt.Sprintf("sl:%s", key))
}

func WithAccessTokenPrefix(key string) string {
	return WithBasePrefix(fmt.Sprintf("at:%s", key))
}

func WithProfileTokenKey() string {
	return WithBasePrefix("ppc:key")
}

func WithDspRegionProfile() string {
	return WithBasePrefix("dsp:key")
}

// WithBatch 批次号通知模式，动态通知
func WithBatch(reportTYpe string) string {
	return fmt.Sprintf("dp:batch:%s", reportTYpe)
}

func WithNewBatch(reportType string, prefixQueue string) string {
	return fmt.Sprintf("dp:%s:%s", prefixQueue, reportType)
}

// WithDivide 切割文件批次
func WithDivide(reportType string) string {
	return fmt.Sprintf("dp:divide:%s", reportType)
}

func WithBatchCheck(reportType, reportName, date string) string {
	return WithBasePrefix(fmt.Sprintf("batch_check:%s:%s:%s", reportType, reportName, date))
}

func WithSchdule(reportType, reportName, date string) string {
	return fmt.Sprintf("dp:schdule:%s:%s:%s", reportType, reportName, date)
}

func WithNewSchdule(reportType string, prefixQueue string) string {
	return fmt.Sprintf("dp:new_schdule:%s:%s", prefixQueue, reportType)
}

func WithSchduleDetail(reportType, reportName, date string) string {
	return fmt.Sprintf("dp:schdule_detail:%s:%s:%s", reportType, reportName, date)
}

func WithUploadFile(key string) string {
	return WithBasePrefix(fmt.Sprintf("%s", key))
}
