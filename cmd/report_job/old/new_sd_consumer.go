package old

import (
	"github.com/spf13/cobra"
	"math/rand"
	"td_report/app/server"
	"td_report/common/reporttool"
	"td_report/pkg/logger"
	"time"
)

// 按照一个profile拉取所有的报表

var NewsdConsumerCmd = &cobra.Command{
	Use:   "new_sd_consumer",
	Short: "new_sd_consumer",
	Long:  `new_sd_consumer.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("new_sd_consumer", true)
		logger.Logger.Info("new_sd_consumer called")
		consumerServer := server.NewConsumerServer()
		if QueueType != "" {
			if reporttool.CheckQueueName(QueueType) == false {
				logger.Logger.Error("paramas error")
				return
			}
		}
		consumerServer.SdConsumer(QueueType)
	},
}

func init() {
	rand.Seed(time.Now().UnixNano())
	NewsdConsumerCmd.PersistentFlags().StringVar(&QueueType, "queue_type", "", "指定队列类型queue_type(fast,middle,slow,back)")
}
