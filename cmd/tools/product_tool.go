package tools

import (
	"fmt"
	"github.com/spf13/cobra"
	"td_report/app"
	"td_report/app/repo"
	"td_report/app/server"
	"td_report/app/service/report"
	"td_report/pkg/logger"
)

// 工具针对指定生产日期周期的情况，一般的情况用不上，主要用户手动拉取指定日期的数据
var productToolCmd = &cobra.Command{
	Use:   "product_tool",
	Short: "product_tool",
	Long:  `product_tool`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("product_tool", false)
		logger.Logger.Info("product_tool called")
		profileService := app.InitializeProfileService()
		schduleRepo := repo.NewReportSchduleRepository()
		retryRepo := repo.NewReportCheckRetryDetailRepository()
		addQueueService := report.NewAdddataService(profileService, schduleRepo, retryRepo)
		productServer := server.NewProductServer(addQueueService)
		if ReportType == "" || StartDate == "" || EndDate == "" {
			fmt.Println("paramas is empty ")
			return
		} else {
			productServer.ProductWithDate(ReportType, StartDate, EndDate)
		}

	},
}

func init() {
	RootCmd.AddCommand(productToolCmd)
}
