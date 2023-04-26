package tools

import (
	"github.com/spf13/cobra"
	"td_report/app/service/tool"
	"td_report/pkg/logger"
)

var monitorError = &cobra.Command{
	Use:   "monitor_error",
	Short: "监控每日的错误数据",
	Long:  `monitor_error`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Init("monitor_error", false)
		logger.Logger.Info("monitor_error_called")
	},
	Run: func(cmd *cobra.Command, args []string) {
		tool.NewMonitorErrorServer().Run()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("monitor_error_over")
	},
}

func init() {
	RootCmd.AddCommand(monitorError)
}


