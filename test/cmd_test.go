package test

import (
	"fmt"
	"td_report/app"
	"td_report/app/repo"
	"td_report/app/service/report"
	"td_report/common/tool"
	"td_report/vars"
	"testing"
)

func TestCmd(t *testing.T) {
	tool.RunCommand("", "ping", "127.0.0.1")
}

func TestDsp(t *testing.T) {
	profileService := app.InitializeProfileService()
	schduleRepo := repo.NewReportSchduleRepository()
	retryRepo := repo.NewReportCheckRetryDetailRepository()
	addQueueService := report.NewAdddataService(profileService, schduleRepo, retryRepo)
	profileList := make([]string, 0)
	profileList = append(profileList, "2771045500938617")
	groupData := addQueueService.DspAdgroupData(profileList)
	str := "202208"
	var temp string
	for i := 1; i <= 31; i++ {
		if i < 10 {
			temp = fmt.Sprintf("%s0%d", str, i)
		} else {
			temp = fmt.Sprintf("%s%d", str, i)
		}

		addQueueService.NewAddDspData(temp, vars.DSP, 0, groupData, "fast")
	}
}
