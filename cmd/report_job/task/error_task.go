package task

import (
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/spf13/cobra"
	"math/rand"
	"td_report/app/bean"
	"td_report/app/repo"
	rabbitmq "td_report/common/rabbitmqV1"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

var ErrortaskCmd = &cobra.Command{
	Use:   "error_task",
	Short: "error_task",
	Long:  `error_task`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("rabbitmq_error_info", false)
		logger.Logger.Info("rabbitmq_error_info called", time.Now().Format(vars.TIMEFORMAT))
		errortaskCmdfunc()
	},
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func errortaskCmdfunc() {
	rabbitMq, err := rabbitmq.NewRabbitmq(g.Cfg().GetString("rabbitmq.address"))
	if err != nil {
		logger.Logger.Error(err)
		panic(err)
	}
	errorRepository := repo.NewReportErrorRepository()
	rabbitMq.Receive(vars.ErrorInfo, func(bytes []byte) error {
		var errorDetail *bean.ReportErr
		if err := json.Unmarshal(bytes, &errorDetail); err != nil {
			logger.Logger.Error(err)
			return err
		}

		errorRepository.AddOne(errorDetail)
		return nil
	})
}
