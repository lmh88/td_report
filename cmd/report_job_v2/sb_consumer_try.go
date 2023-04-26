package report_job_v2

import (
	"github.com/spf13/cobra"
	"td_report/app/service/report_v2/one_step"
	"td_report/pkg/logger"
	"td_report/vars"
)

// 按照一个profile拉取所有的报表

var sbConsumerTryCmd = &cobra.Command{
	Use:   "sb_consumer_try",
	Short: "sb消费第一步重试",
	Long:  `sb消费第一步重试，如：sb_consumer_try`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Init("sb_consumer_try_v2", false)
		logger.Logger.Info("sb_consumer_try called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		checkClientTag(ClientTag)
		one_step.RetryConsumer(vars.SB, ClientTag)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("sb_consumer_try over")
	},
}

func init() {
	RootCmd.AddCommand(sbConsumerTryCmd)
}
