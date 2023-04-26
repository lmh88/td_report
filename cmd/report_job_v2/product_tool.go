package report_job_v2

import (
	"fmt"
	"github.com/spf13/cobra"
	"td_report/app/service/report_v2/product_server"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

var productTool = &cobra.Command{
	Use:   "product_tool",
	Short: "定制生成profile",
	Long:  `生成profile, 如：product --report_type=sp`,
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Init("product_tool_v2", false)
		logger.Logger.Info("product_tool_v2 called")
		fmt.Println("product_tool_v2 start")
	},
	Run: func(cmd *cobra.Command, args []string) {
		if ReportType == "" {
			fmt.Println("report_type不能为空")
			return
		}

		if StartDate == "" || EndDate == "" {
			fmt.Println("start_date,end_date不能为空")
			return
		}

		t1, err1 := time.Parse(vars.TimeLayout, StartDate)
		t2, err2 := time.Parse(vars.TimeLayout, EndDate)
		if err1 != nil || err2 != nil {
			fmt.Println("start_date,end_date格式不对，如：20220401")
			return
		}

		if t1.Unix() > t2.Unix() {
			fmt.Println("start_date要小于等于end_date")
			return
		}

		var err error
		if ProfileId == "" {
			err = product_server.NewProductServer().ProductWithDate(ReportType, StartDate, EndDate)
		} else {
			err = product_server.NewProductServer().ProductOneProfile(ReportType, StartDate, EndDate, ProfileId)
		}

		if err != nil {
			fmt.Println(err.Error())
		}
		return
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		logger.Logger.Info("product_tool_v2 run over")
		fmt.Println("product_tool_v2 end")
	},
}

func init() {
	//productCmd.PersistentFlags().IntVar(&productday, "day", 2, "生产的天数")
	RootCmd.AddCommand(productTool)
}
