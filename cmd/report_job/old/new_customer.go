package old

import (
	"github.com/gogf/gf/frame/g"
	"github.com/spf13/cobra"
	"math/rand"
	"td_report/app/server"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

var productTotalNum = g.Cfg().GetInt64("queue.queue_all_full")

// SchduleCmd 处理新客户检测和重试
var SchduleCmd = &cobra.Command{
	Use:   "newclient",
	Short: "newclient",
	Long:  `newclient`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("newclient", false)
		logger.Logger.Info("newclient called", time.Now().Format(vars.TIMEFORMAT))
		var startdate string
		if StartDate != "" {
			if mytime, err := time.Parse(vars.TimeLayout, StartDate); err != nil {
				logger.Logger.Error(map[string]interface{}{
					"err":  err,
					"desc": "newclient 转化时间错误",
				})

				return
			} else {
				startdate = mytime.Format(vars.TimeFormatTpl)
			}

		} else {
			startdate = time.Now().Format(vars.TimeFormatTpl)
		}

		customerServer := server.NewCustomerServer()
		customerServer.DealNewCustomer(startdate)
	},
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
