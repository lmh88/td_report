package report_job_v2

import (
	"github.com/spf13/cobra"
	"td_report/app/service/report_v2/consumer_server"
	"td_report/pkg/logger"
)

var checkBusyCmd = &cobra.Command{
	Use:   "check_busy",
	Short: "检查队列繁忙",
	Long:  `检查队列繁忙, 如：check_busy `,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Init("check_busy", false)
		logger.Logger.Info("check_busy called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		consumer_server.NewCheckBusyServer().CheckQueueBusy()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("check_busy run over")
	},
}

func init() {
	RootCmd.AddCommand(checkBusyCmd)
}
