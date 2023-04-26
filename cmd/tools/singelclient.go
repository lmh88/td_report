package tools

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"td_report/pkg/logger"
	"td_report/vars"
)

// 获取单个profile的所有的报表的某一天的数据
var singelClientCmd = &cobra.Command{
	Use:   "singel_client",
	Short: "singel_client",
	Long:  `singel_client`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("singel_client", false)
		logger.Logger.Info("singel_client called")
		if ProfileId == "" || ReportType == "" || StartDate == "" {
			logger.Logger.Info("report_type or profileId or start_date is empty")
			return
		}

		singelClienthFunc(ProfileId, ReportType)
	},
}

func init() {
	RootCmd.AddCommand(singelClientCmd)
}

// 后期添加队列里面，避免协成调度的控制
func singelClienthFunc(profileId, reportType string) {
	reportNameList := vars.ReportList[reportType]
	var cmd *exec.Cmd
	cmd.Dir = "/project/tool_job"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	for _, reportName := range reportNameList {
		cmd = exec.Command("toolmain", "singel_profile",
			fmt.Sprintf("--report_type=%s", reportType),
			fmt.Sprintf("--profile_id=%s", profileId),
			fmt.Sprintf("--report_name=%s", reportName),
			fmt.Sprintf("--startdate=%s", StartDate),
		)
		err := cmd.Run()
		if err != nil {
			logger.Logger.Info(err)
		}
	}
}
