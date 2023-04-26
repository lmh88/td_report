package job

import (
	"github.com/gogf/gf/util/gconv"
	"github.com/spf13/cobra"
	"strings"
	"td_report/app"
	"td_report/pkg/logger"
)

var FeadSubCmd = &cobra.Command{
	Use:   "fead",
	Short: "fead",
	Long:  `fead`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("fead_substr", false)
		logger.Logger.Info("call ")
		feadService := app.InitializeFeadService()
		profileIdList := make([]int64, 0)
		if len(profileId) != 0 {
			if strings.Contains(profileId, ",") {
				profilestr := strings.Split(profileId, ",")
				for _, item := range profilestr {
					profileIdList = append(profileIdList, gconv.Int64(item))
				}
			} else {
				profileIdList = append(profileIdList, gconv.Int64(profileId))
			}
		}

		feadService.Fead(1, profileIdList...)
		logger.Logger.Info("end  ")
	},
}

func init() {
	FeadSubCmd.PersistentFlags().StringVar(&profileId, "profileId", "", "profileid多个用逗号隔开")
}
