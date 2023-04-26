package tools

import (
	"fmt"
	"github.com/spf13/cobra"
	"td_report/app"
	"td_report/app/repo"
	"td_report/app/server"
	"td_report/app/service/report"
	"td_report/app/service/report_v2/product_server"
	"td_report/pkg/logger"
	"td_report/vars"
)

var productMonthCmd = &cobra.Command{
	Use:   "product_month",
	Short: "生产上一个月",
	Long:  `生产上一个月， product_month --report_type=sd`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("product_tool", false)
		logger.Logger.Info("product_tool called")
		profileService := app.InitializeProfileService()
		schduleRepo := repo.NewReportSchduleRepository()
		retryRepo := repo.NewReportCheckRetryDetailRepository()
		addQueueService := report.NewAdddataService(profileService, schduleRepo, retryRepo)
		productServer := server.NewProductServer(addQueueService)
		if ReportType == "" {
			fmt.Println("report_type is empty ")
			return
		}
		allowType := map[string]bool{
			vars.DSP: true,
			vars.SD:  true,
		}
		if _, ok := allowType[ReportType]; !ok {
			fmt.Println("report_type is dsp or sd ")
			return
		}
		start, end := product_server.NewProductServer().GetLastMonth()
		if start != "" && end != "" {
			productServer.ProductWithDate(ReportType, start, end)
		}
	},
}

func init() {
	RootCmd.AddCommand(productMonthCmd)
}
