package tools

import (
	"math/rand"
	"td_report/app"
	"td_report/app/repo"
	"td_report/app/server"
	"td_report/app/service/report"
	"td_report/common/reportsystem"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/spf13/cobra"
)

var checktype string // 0 全部报表  1 单个reportType的报表下面的全部报表名称 2 单个报表类型下单个报表名称 3 单个报表类型，报表名称下的单个profile
var checkday int     //  检测的天数， 报表日期是固定的，日期需要动态变化
var productNum = g.Cfg().GetInt64("queue.queue_all_full")

// statisticsToolCmd represents the statisticsTool command
// 统计工具，根据服务器上拉取文件情况做统计分析
var statisticsToolCmd = &cobra.Command{
	Use:   "statistics_tool",
	Short: "statistics_tool",
	Long:  `statistics_tool .`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("statistics_tool", false)
		logger.Logger.Info("statisticsTool called", time.Now().Format(vars.TIMEFORMAT))
		if checkday == 0 {
			logger.Logger.Info("check_day is  empty")
			return
		}
		if checkparamas() == false {
			logger.Logger.Info("paramas is  error")
			return
		}

		if StartDate == "" && EndDate == "" {
			//EndDate = time.Now().Add(24 * time.Hour * time.Duration(-2)).Format(vars.TimeLayout)
			EndDate = time.Now().Format(vars.TimeLayout)
			StartDate = time.Now().Add(24 * time.Hour * time.Duration(checkday*-1)).Format(vars.TimeLayout)
		}

		profileList := make([]string, 0)
		if ProfileId != "" {
			profileList = append(profileList, ProfileId)
		}

		service := app.InitializeStatisticsService()
		profileService := app.InitializeProfileService()
		schduleRepo := repo.NewReportSchduleRepository()
		retryRepo := repo.NewReportCheckRetryDetailRepository()
		addQueueService := report.NewAdddataService(profileService, schduleRepo, retryRepo)
		checkServer := server.NewCheckReportServer(addQueueService, service)

		wg := reportsystem.NewPool(30)
		switch checktype {
		case "0":
			reportTypeList := []string{vars.SP, vars.SB, vars.SD, vars.DSP}
			for _, reportType := range reportTypeList {
				reportNameList := vars.ReportList[reportType]
				for _, reportName := range reportNameList {
					wg.Add(1)
					go func(mReportType, mReportName, mStartDate, mEndDate string) {
						defer wg.Done()
						checkServer.CheckProfile(mReportType, mReportName, mStartDate, mEndDate, profileList, true)
					}(reportType, reportName, StartDate, EndDate)
				}
			}

		case "1":
			reportNameList := vars.ReportList[ReportType]
			for _, reportName := range reportNameList {
				wg.Add(1)
				go func(mReportType, mreportName, mStartDate, mEndDate string) {
					defer wg.Done()
					checkServer.CheckProfile(mReportType, mreportName, mStartDate, mEndDate, profileList, true)
				}(ReportType, reportName, StartDate, EndDate)
			}

			wg.Wait()

		case "2":
			profileList := make([]string, 0)
			checkServer.CheckProfile(ReportType, ReportName, StartDate, EndDate, profileList, true)
		case "3":
			if ProfileId == "" {
				logger.Logger.Info("profileid is empty")
				return
			}

			checkServer.CheckProfile(ReportType, ReportName, StartDate, EndDate, profileList, true)
		}

		logger.Logger.Info("statisticsTool called end ", time.Now().Format(vars.TIMEFORMAT))
	},
}

func checkparamas() bool {
	if checktype != "0" && checktype != "1" && checktype != "2" && checktype != "3" {
		logger.Logger.Println("checktype is in (0,1,2,3)")
		return false
	}

	if checktype == "1" && (ReportType != vars.SP && ReportType != vars.SB && ReportType != vars.SD && ReportType != vars.DSP) {
		logger.Logger.Println("ReportType is in (sp,sb,sd,dsp)")
		return false
	}

	if checktype == "2" {
		if ReportType != vars.SP && ReportType != vars.SB && ReportType != vars.SD && ReportType != vars.DSP {
			logger.Logger.Println("ReportType is in (sp,sb,sd,dsp)")
			return false
		}

	}

	return true
}

func init() {
	rand.Seed(time.Now().UnixNano())
	statisticsToolCmd.PersistentFlags().StringVar(&checktype, "check_type", "", "check_type(0,1,2,3)")
	statisticsToolCmd.PersistentFlags().IntVar(&checkday, "check_day", 14, "check_day:check how many days ")
	RootCmd.AddCommand(statisticsToolCmd)
}
