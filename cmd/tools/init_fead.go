package tools

import (
	"github.com/spf13/cobra"
	"td_report/app/service/fead"
	"td_report/pkg/logger"
)

var feadProfile string

var initFead = &cobra.Command{
	Use:   "init_fead",
	Short: "初始化fead订阅",
	Long:  `初始化fead订阅`,
	PreRun: func(cmd *cobra.Command, args []string) {
		//logger.Init("init_fead")
		logger.Logger.Info("init_fead start")
	},
	Run: func(cmd *cobra.Command, args []string) {
		fead.InitSubscription(feadProfile)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("init_fead over")
	},
}

func init() {
	RootCmd.AddCommand(initFead)
	RootCmd.PersistentFlags().StringVar(&feadProfile, "fead_profile", "valid", "跑有效用户和无效用户，valid或invalid")
}


