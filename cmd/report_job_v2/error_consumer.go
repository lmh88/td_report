package report_job_v2

import (
	"github.com/spf13/cobra"
	"td_report/app/service/report_v2/consumer_server"
	"td_report/pkg/logger"
)

// fail队列数据消费
var errorConsumerTryCmd = &cobra.Command{
	Use:   "error_consumer",
	Short: "失败数据消费",
	Long:  `失败数据消费，如：error_consumer`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Init("error_consumer", false)
		logger.Logger.Info("error_consumer called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		consumer_server.NewErrServer().ConsumeFailQueue()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("error_consumer over")
	},
}

func init() {
	RootCmd.AddCommand(errorConsumerTryCmd)
}
