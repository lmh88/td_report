package old

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"td_report/app"
	"td_report/app/bean"
	comonservice "td_report/app/service/common"
	"td_report/common/file"
	"td_report/common/reportsystem"
	"td_report/common/tool"
	"td_report/pkg/all_steps"
	"td_report/pkg/common"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"

	"github.com/spf13/cobra"
)

// SbmetricsCmd 独立的sb metricsCmd 拉取时间是特定的时间区间
var SbmetricsCmd = &cobra.Command{
	Use:   "sbmetrics",
	Short: "sbmetrics",
	Long:  `sbmetrics`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("sbmetrics", false)
		logger.Logger.Info("sbmetrics called", time.Now().Format(vars.TIMEFORMAT))
		SbmetricsCmdFunc(ProfileId)
	},
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Jobs struct {
	dateStr     bean.Mydate
	reportName  string
	profileList []*bean.ProfileToken
}

type Result struct {
	reportName string
	StartDate  string
	reportType string
}

var jobs = make(chan Jobs)

//var jobResult = make(chan Result)

func worker(wg *reportsystem.Pool) {
	for job := range jobs {
		for _, profile := range job.profileList {
			wg.Add(1)
			go func(date bean.Mydate, mywg *reportsystem.Pool, myprofile *bean.ProfileToken, reportName string) {
				ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s_%d", "sb", reportName, myprofile.ProfileId, time.Now().UnixNano()))
				result, _ := all_steps.SbAllSteps(myprofile, reportName, date.StartDate, date.EndDate, mywg, ctx)
				if result {
					logger.Logger.InfoWithContext(ctx, "sb的metrics报告消费结束,执行成功")
				} else {
					logger.Logger.InfoWithContext(ctx, "sb报告消费结束，执行失败")
				}
			}(job.dateStr, wg, profile, job.reportName)
		}

		//jobResult <- Result{reportName: job.reportName, StartDate: job.dateStr.StartDate, reportType: vars.SB}
	}

}

func SbmetricsCmdFunc(dataProfileId string) {
	today := time.Now()
	day := 183
	starttime := today.Add(-1 * time.Duration(day) * time.Hour * 24)
	profileService := app.InitializeProfileService()
	profileIdList := make([]string, 0)
	if dataProfileId != "" {
		profileIdList = append(profileIdList, dataProfileId)
	}

	sps, _ := profileService.GetPpcProfile(profileIdList)
	wg := reportsystem.NewPool(10)
	daylist := tool.GetBetweenWeek(starttime.Format(vars.TimeFormatTpl), today.Format(vars.TimeFormatTpl), vars.TimeFormatTpl)
	dayjob := Jobs{
		profileList: sps,
		reportName:  vars.BrandMetricsWeekly,
	}
	startMonthTime := today.Add(-1 * time.Duration(day) * time.Hour * 24)
	monthlist := tool.GetBetweenmonth(startMonthTime.Format(vars.TimeFormatTpl), today.Format(vars.TimeFormatTpl), vars.TimeFormatTpl)
	monthjob := Jobs{
		profileList: sps,
		reportName:  vars.BrandMetricsMonthly,
	}

	//taskService := app.InitializeReportTaskService()

	go worker(wg)

	//go func() {
	//	for result := range jobResult {
	//		changestatus(taskService, result.reportName, result.StartDate, result.reportType)
	//	}
	//}()

	for _, item := range monthlist {
		monthjob.dateStr = item
		jobs <- monthjob

	}

	for _, item := range daylist {
		dayjob.dateStr = item
		jobs <- dayjob
	}

	close(jobs)
	wg.Wait()
}

func changestatus(taskService *comonservice.ReportTaskService, reportName string, StartDate string, reportType string) {
	localDir := file.GetPath(reportType, reportName, StartDate)
	num := common.GetDirFileNum(localDir)
	if num > 0 {
		splist := strings.Split(localDir, "/")
		fileLength := len(splist)
		mydate := splist[fileLength-1]
		taskService.Add(localDir, reportName, reportType, mydate)
	} else {

		logger.Error(map[string]interface{}{
			"flag": localDir + "current file dir is empty",
		})
	}
}
