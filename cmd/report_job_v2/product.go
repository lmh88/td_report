package report_job_v2

import (
	"fmt"
	"github.com/spf13/cobra"
	"td_report/app/service/report_v2/product_server"
	"td_report/app/service/report_v2/varible"
	"td_report/pkg/logger"
	"td_report/vars"
)

// productDay 生产的天数 2:今天和昨天 7:最近7天  14: 最近14天
var productDay int

var productCmd = &cobra.Command{
	Use:   "product",
	Short: "生成profile",
	Long:  `生成profile, 如：product --report_type=sp`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Init("product_v2", false)
		logger.Logger.Info("product_v2 called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		if ReportType == "" {
			fmt.Println("report_type不能为空")
			return
		}

		_, ok := varible.ReportQueueMap[ReportType]
		if !ok {
			fmt.Println("report_type不正确")
			return
		}

		if LimitProductTime() {
			fmt.Println("limit product time")
			return
		}

		if productDay == vars.ProductDay30 {
			start, end := product_server.NewProductServer().GetLastMonth()
			if start != "" && end != "" {
				product_server.NewProductServer().ProductWithDate(ReportType, start, end)
			}
		} else {
			product_server.NewProductServer().Product(ReportType, productDay)
		}

	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("product_v2 run over")
	},
}

func init() {
	productCmd.PersistentFlags().IntVar(&productDay, "day", 2, "生产的天数")
	RootCmd.AddCommand(productCmd)
}


