package tools

import (
	"github.com/spf13/cobra"
	"td_report/app/server"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

// 新用户拉取失败重试
var newclientretryCmd = &cobra.Command{
	Use:   "retry",
	Short: "retry",
	Long:  `retry`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("retry", false)
		logger.Logger.Info("retry called")
		var (
			mydate time.Time
			err    error
		)

		if StartDate != "" {
			if mydate, err = time.Parse(vars.TimeLayout, StartDate); err != nil {
				logger.Logger.Error(map[string]interface{}{
					"err":  err,
					"desc": "时间格式错误",
				})
			}
		} else {
			mydate = time.Now()
		}

		retryFunc(mydate)
	},
}

func init() {
	RootCmd.AddCommand(newclientretryCmd)
}

func retryFunc(date time.Time) {
	customerServer := server.NewCustomerServer()
	customerServer.Retry(date.Format(vars.TimeFormatTpl))
}
