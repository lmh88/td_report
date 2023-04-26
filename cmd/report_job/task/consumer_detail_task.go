package task

import (
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/spf13/cobra"
	"math/rand"
	"td_report/app/bean"
	"td_report/app/service/report"
	rabbitmq "td_report/common/rabbitmqV1"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

var ConsumerDetailaskCmd = &cobra.Command{
	Use:   "consumer_detail_task",
	Short: "consumer_detail_task",
	Long:  `consumer_detail_task`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("rabbitmq_consumer_detail", false)
		logger.Logger.Info("rabbitmq_consumer_detail called", time.Now().Format(vars.TIMEFORMAT))
		consumerDetailtaskCmdFunc()
	},
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func consumerDetailtaskCmdFunc() {
	rabbitMq, err := rabbitmq.NewRabbitmq(g.Cfg().GetString("rabbitmq.address"))
	if err != nil {
		logger.Logger.Error(err)
		panic(err)
	}

	ConsumerDetailService := report.NewConsumerDetailService()
	rabbitMq.Receive(vars.ProfileidDetail, func(bytes []byte) error {
		var consumerDetail *bean.ConsumerDetail
		if err := json.Unmarshal(bytes, &consumerDetail); err != nil {
			logger.Logger.Error(err)
			return err
		}

		ConsumerDetailService.AddConsumerDetail(consumerDetail)
		return nil
	})
}
