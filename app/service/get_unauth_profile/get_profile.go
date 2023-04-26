package get_unauth_profile

import (
	"fmt"
	"td_report/app/repo"
	"td_report/boot"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

const (
	Prefix = "unauth_profile"
)

func StoreProfile(reportType string, date time.Time) {
	key := GetRdsKey(reportType, date.Format(vars.TimeLayout))
	rds := boot.RedisCommonClient.GetClient()
	profiles := getProfiles(reportType, date)
	if len(profiles) > 0 {
		rds.SAdd(key, profiles)
		rds.Expire(key, time.Hour * 24)
	}
}

func GetRdsKey(reportType, date string) string {
	return fmt.Sprintf("%s:%s:%s", Prefix, reportType, date)
}

func getProfiles(reportType string, date time.Time) (profiles []string) {
	profiles = make([]string, 0)
	profileNum, err := repo.NewReportErrorRepository().GetProfileNum(reportType, date, repo.ReportErrorTypeOne)
	if err != nil {
		logger.Error(err.Error())
	}

	if len(profileNum) > 0 {
		for _, item := range profileNum {
			if item.Num >= 4 {
				profiles = append(profiles, item.ProfileId)
			}
		}
	}

	return profiles
}