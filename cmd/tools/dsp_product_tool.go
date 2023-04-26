package tools

import (
	"fmt"
	"github.com/spf13/cobra"
	"td_report/app"
	"td_report/app/repo"
	"td_report/app/server"
	"td_report/app/service/report"
	datetool "td_report/common/tool"
	"td_report/pkg/logger"
	"td_report/vars"
)

var dspproductToolCmd = &cobra.Command{
	Use:   "dsp_product_tool",
	Short: "dsp_product_tool",
	Long:  `dsp_product_tool`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("dsp_product_tool", false)
		logger.Logger.Info("dsp_product_tool called")
		//ProfileId = "1898361487593447"
		//ReportName = "order"
		//StartDate = "20220808"
		//EndDate = "20220821"
		if ProfileId == "" || ReportName == "" || StartDate == "" || EndDate == "" {
			fmt.Println("参数缺失")
			return
		}

		dspProductTool(ProfileId, ReportName, StartDate, EndDate)
	},
}

func init() {
	RootCmd.AddCommand(dspproductToolCmd)
}

func dspProductTool(profileid, reportName, startdate, enddate string) {
	if daylist, err := datetool.GetDays(startdate, enddate, vars.TimeLayout); err != nil {
		logger.Logger.Info(vars.DSP, "product error", err.Error())
	} else {

		prefixQueue := vars.FastQueue
		profileService := app.InitializeProfileService()
		schduleRepo := repo.NewReportSchduleRepository()
		retryRepo := repo.NewReportCheckRetryDetailRepository()
		addQueueService := report.NewAdddataService(profileService, schduleRepo, retryRepo)
		productServer := server.NewProductServer(addQueueService)
		profileList := make([]string, 0)
		profileList = append(profileList, profileid)
		groupData := productServer.AdddataService.DspAdgroupData(profileList)
		if groupData == nil {
			logger.Logger.Info(err, "dsp get profileList groupdata error")
			return
		}
		for _, reportday := range daylist {
			productServer.AdddataService.NewAddDspData(reportday, vars.DSP, 0, groupData, prefixQueue)
		}
	}
}
