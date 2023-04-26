package check_report_file

import (
	"td_report/boot"
	"time"
)

//获取某一天的黑名单redis key
func GetBackListKey(reportType string, date time.Time) string {
	prefix := GetTypePrefix(reportType)
	return GetRdsKey(prefix, date)
}

//黑名单中是否存在profile
func ExistProfile(reportType string, date time.Time, profile string) bool {
	rds := boot.RedisCommonClient.GetClient()
	key := GetBackListKey(reportType, date)
	return rds.SIsMember(key, profile).Val()
}

