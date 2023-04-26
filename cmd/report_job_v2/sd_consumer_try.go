package report_job_v2

import (
	"github.com/spf13/cobra"
	"td_report/app/service/report_v2/one_step"
	"td_report/pkg/logger"
	"td_report/vars"
)

// 按照一个profile拉取所有的报表

var sdConsumerTryCmd = &cobra.Command{
	Use:   "sd_consumer_try",
	Short: "sd消费第一步重试",
	Long:  `sd消费第一步重试，如：sd_consumer_try`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Init("sd_consumer_try_v2", false)
		logger.Logger.Info("sd_consumer_try called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		checkClientTag(ClientTag)
		one_step.RetryConsumer(vars.SD, ClientTag)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("sd_consumer_try over")
	},
}

func init() {
	RootCmd.AddCommand(sdConsumerTryCmd)
}
