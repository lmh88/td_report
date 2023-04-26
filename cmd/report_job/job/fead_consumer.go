package job

import (
	"fmt"
	"github.com/spf13/cobra"
	"math/rand"
	"td_report/app"
	"td_report/pkg/logger"
	"time"
)

var QueueName string

// FeadConsumerCmd 用于消费数据：包含数据的确认和业务数据的订阅
var FeadConsumerCmd = &cobra.Command{
	Use:   "fead_data",
	Short: "fead_data",
	Long:  `fead_data`,
	Run: func(cmd *cobra.Command, args []string) {
		feadService := app.InitializeFeadService()
		if QueueName != "" {
			queueKey:=fmt.Sprintf("fead_data_%s",QueueName)
			logger.Init(queueKey, true)
			logger.Logger.Info(queueKey + " called")
			//feadService.Receive(QueueName)
			feadService.Receivelp(QueueName)
		} else {
			fmt.Println("队列参数缺失")
			return
		}
	},
}

func init() {
	rand.Seed(time.Now().UnixNano())
	FeadConsumerCmd.PersistentFlags().StringVar(&QueueName, "queue_name", "", "队列名称")
}
