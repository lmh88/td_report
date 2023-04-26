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

var NewspConsumerCmd = &cobra.Command{
	Use:   "new_sp_consumer",
	Short: "new_sp_consumer",
	Long:  `new_sp_consumer.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("new_sp_consumer", true)
		logger.Logger.Info("new_sp_consumer called")
		consumerServer := server.NewConsumerServer()
		if QueueType != "" {
			if reporttool.CheckQueueName(QueueType) == false {
				logger.Logger.Error("paramas error")
				return
			}
		}
		consumerServer.SpConsumer(QueueType)
	},
}

func init() {
	rand.Seed(time.Now().UnixNano())
	NewspConsumerCmd.PersistentFlags().StringVar(&QueueType, "queue_type", "", "指定队列类型queue_type(fast,middle,slow,back)")
}
