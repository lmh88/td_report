package report_job_v2

import (
	"github.com/spf13/cobra"
	"math/rand"
	"td_report/app/service/report_v2/consumer_server"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

// 按照一个profile拉取所有的报表

var sbConsumerTwoCmd = &cobra.Command{
	Use:   "sb_consumer_two",
	Short: "sb消费第二步",
	Long:  `sb消费第二步，如：sb_consumer_two`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Init("sb_consumer_two_v2", false)
		logger.Logger.Info("sb_consumer_two_v2 called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		checkClientTag(ClientTag)
		consumer_server.NewConsumerTwoServer().ConsumerTwo(vars.SB, ClientTag)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("sb_consumer_two_v2 over")
	},
}

func init() {
	rand.Seed(time.Now().UnixNano())
	RootCmd.AddCommand(sbConsumerTwoCmd)
}
