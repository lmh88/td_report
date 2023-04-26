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

var spConsumerTwoCmd = &cobra.Command{
	Use:   "sp_consumer_two",
	Short: "sp消费第二步",
	Long:  `sp消费第二步，如：sp_consumer_two`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Init("sp_consumer_two_v2", false)
		logger.Logger.Info("sp_consumer_two_v2 called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		checkClientTag(ClientTag)
		consumer_server.NewConsumerTwoServer().ConsumerTwo(vars.SP, ClientTag)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("sp_consumer_two_v2 over")
	},
}

func init() {
	rand.Seed(time.Now().UnixNano())
	RootCmd.AddCommand(spConsumerTwoCmd)
}
