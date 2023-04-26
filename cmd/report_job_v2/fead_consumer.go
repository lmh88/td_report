package report_job_v2

import (
	"fmt"
	"github.com/spf13/cobra"
	"td_report/app/service/fead"
	"td_report/pkg/logger"
)

var feadQueue string

// fead队列数据消费
var feadConsumerCmd = &cobra.Command{
	Use:   "fead_consumer",
	Short: "fead数据消费",
	Long:  `fead数据消费，如：fead_consumer`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Init("fead_consumer", false)
		logger.Logger.Info("fead_consumer called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		if _, ok := fead.SqsMap[feadQueue]; !ok {
			fmt.Println(feadQueue)
			fmt.Println("队列名错误")
			return
		}
		err := fead.Receive(feadQueue)
		if err != nil {
			fmt.Println("error:", err.Error())
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("fead_consumer over")
	},
}

func init() {
	feadConsumerCmd.PersistentFlags().StringVar(&feadQueue, "fead_queue", "", "fead队列名")
	RootCmd.AddCommand(feadConsumerCmd)
}
