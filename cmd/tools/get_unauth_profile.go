package tools

import (
	"fmt"
	"github.com/spf13/cobra"
	"td_report/app/service/get_unauth_profile"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

var needDate string

var getUnauthProfile = &cobra.Command{
	Use:   "get_unauth_profile",
	Short: "获取未授权的profile",
	Long:  `get_unauth_profile`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Init("get_unauth_profile", false)
		logger.Logger.Info("get_unauth_profile_called")
	},
	Run: func(cmd *cobra.Command, args []string) {

		var (
			date time.Time
			err  error
		)

		if needDate != "" {
			if date, err = time.Parse(vars.TimeLayout, needDate); err != nil {
				logger.Logger.Error(map[string]interface{}{
					"err":  err,
					"desc": "时间格式错误",
				})
				fmt.Println("时间格式错误")
				return
			}
		} else {
			date = time.Now()
		}

		if ReportType == "" {
			fmt.Println("report_type必填")
			return
		}
		get_unauth_profile.StoreProfile(ReportType, date)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("get_unauth_profile_over")
	},
}

func init() {
	RootCmd.AddCommand(getUnauthProfile)
	RootCmd.PersistentFlags().StringVar(&needDate, "need_date", "", "日期，如：20220322")
}


