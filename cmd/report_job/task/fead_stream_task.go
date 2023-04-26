package task

import (
	"encoding/json"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/spf13/cobra"
	"td_report/app"
	rabbitmq "td_report/common/rabbitmqV1"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

// FeadtaskCmd 处理新客户feadstream开通
var FeadtaskCmd = &cobra.Command{
	Use:   "fead_task",
	Short: "fead_task",
	Long:  `fead_task`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("fead_task", false)
		logger.Logger.Info("fead_task called", time.Now().Format(vars.TIMEFORMAT))
		feadtaskCmdfunc()
	},
}

type FeadProfileids struct {
	Profileids []string `json:"profileids"`
}

func feadtaskCmdfunc() {
	rabbitMq, err := rabbitmq.NewRabbitmq(g.Cfg().GetString("rabbitmq.address"))
	if err != nil {
		logger.Logger.Error(err)
		panic(err)
	}

	feadService := app.InitializeFeadService()
	rabbitMq.Receive(vars.FeadNewClient, func(bytes []byte) error {
		var feadProfile *FeadProfileids
		var profileIdList = make([]int64, 0)
		if err := json.Unmarshal(bytes, &feadProfile); err != nil {
			logger.Logger.Error(err)
			return err
		}

		if len(feadProfile.Profileids) != 0 {
			for _, item := range feadProfile.Profileids {
				profileIdList = append(profileIdList, gconv.Int64(item))
			}

			feadService.Fead(1, profileIdList...)
		}

		return nil
	})
}
