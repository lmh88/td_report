package tools

import (
	"fmt"
	"github.com/spf13/cobra"
	"td_report/app/service/clear_tool"
	"td_report/pkg/logger"
)

var clearType string

var clearDb = &cobra.Command{
	Use:   "clear_db",
	Short: "清理数据库14天以前的数据",
	Long:  `clear_db`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Init("clear_db", false)
		logger.Logger.Info("clear_db_called")
	},
	Run: func(cmd *cobra.Command, args []string) {

		if clearType == "" {
			fmt.Println("clearType不能为空")
			return
		}

		if _, ok := clear_tool.DbMap[clearType]; !ok {
			fmt.Println("clearType不正确")
			return
		}

		err := clear_tool.ClearDb(clearType, 14)
		if err != nil {
			logger.Errorf("清理数据库出错:clear_type:%s,error:%s", clearType, err.Error())
			fmt.Println(err.Error())
			return
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("clear_db_over")
	},
}

func init() {
	RootCmd.AddCommand(clearDb)
	RootCmd.PersistentFlags().StringVar(&clearType, "clear_type", "", "清理类型")
}


