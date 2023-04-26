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

var spConsumerCmd = &cobra.Command{
	Use:   "sp_consumer_one",
	Short: "sp消费第一步",
	Long:  `sp消费第一步，如：sp_consumer_one --queue_type=fast`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Init("sp_consumer_one_v2", false)
		logger.Logger.Info("sp_consumer_one_v2 called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		if QueueType != "" {
			if !varible.CheckQueueLevel(vars.SP, QueueType) {
				logger.Logger.Error("queueType error")
				fmt.Println(QueueType, ": is queueType error")
				return
			}
		}
		checkClientTag(ClientTag)
		consumer_server.NewConsumerOneServer().ConsumerOne(QueueType, vars.SP, ClientTag)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("sp_consumer_one_v2 over")
	},
}

func init() {
	rand.Seed(time.Now().UnixNano())
	RootCmd.AddCommand(spConsumerCmd)
}
