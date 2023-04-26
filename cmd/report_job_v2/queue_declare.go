package report_job_v2

import (
	"fmt"
	"github.com/spf13/cobra"
	"td_report/app/service/report_v2"
	"td_report/app/service/report_v2/varible"
	"td_report/pkg/logger"
)

var declareCmd = &cobra.Command{
	Use:   "queue_declare",
	Short: "声明队列",
	Long:  `声明队列, 如：queue_declare --report_type=sp`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Init("queue_declare", false)
		logger.Logger.Info("queue_declare called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		if ReportType != "" {
			_, ok := varible.QueueMap[ReportType]
			if !ok {
				fmt.Println("report_type不正确")
				return
			}
		}

		checkClientTag(ClientTag)
		//fmt.Println(ClientTag)
		report_v2.ScanClient(ReportType, ClientTag)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("queue_declare run over")
	},
}

func init() {
	//declareCmd.PersistentFlags().StringVar(&clientTag, "client_tag", "", "client tag")
	RootCmd.AddCommand(declareCmd)
}
