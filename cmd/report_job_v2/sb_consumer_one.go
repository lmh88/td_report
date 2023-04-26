package report_job_v2

import (
	"fmt"
	"github.com/spf13/cobra"
	"math/rand"
	"td_report/app/service/report_v2/consumer_server"
	"td_report/app/service/report_v2/varible"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

var sbConsumerCmd = &cobra.Command{
	Use:   "sb_consumer_one",
	Short: "sb消费第一步",
	Long:  `sb消费第一步，如：sb_consumer_one --queue_type=fast`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Init("sb_consumer_one_v2", false)
		logger.Logger.Info("sb_consumer_one_v2 called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		if QueueType != "" {
			if !varible.CheckQueueLevel(vars.SB, QueueType) {
				logger.Logger.Error("queueType error")
				fmt.Println(QueueType, ": is queueType error")
				return
			}
		}
		checkClientTag(ClientTag)
		consumer_server.NewConsumerOneServer().ConsumerOne(QueueType, vars.SB, ClientTag)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("sb_consumer_one_v2 over")
	},
}

func init() {
	rand.Seed(time.Now().UnixNano())
	RootCmd.AddCommand(sbConsumerCmd)
}
