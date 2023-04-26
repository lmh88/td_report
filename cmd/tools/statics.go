package tools

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/spf13/cobra"
	"io/ioutil"
	"td_report/common/sendmsg/wechart"
	"td_report/vars"
	"time"
)

var staticsCmd = &cobra.Command{
	Use:   "statis",
	Short: "statis",
	Long:  `statis`,
	Run: func(cmd *cobra.Command, args []string) {
		if StartDate == "" {
			StartDate = time.Now().Format(vars.TimeLayout)
		}
		statisFunc(StartDate)
	},
}

func init() {
	RootCmd.AddCommand(staticsCmd)
}

func statisFunc(startdate string) {
	reportTypeList := []string{vars.SP, vars.SB, vars.SD, vars.DSP}
	var str string
	for _, reportType := range reportTypeList {
		reportNameList := vars.ReportList[reportType]
		path := vars.MypathMap[reportType]
		for _, reportName := range reportNameList {
			tempPath := path + "/" + reportName + "/" + startdate
			files, _ := ioutil.ReadDir(tempPath)
			num := len(files)
			str = str + fmt.Sprintf("report_type:%s,report_name:%s report_date:%s, num:%d \n", reportType, reportName, startdate, num)
		}
	}

	if str != "" {
		fmt.Println(str)
		key := g.Cfg().GetString("wechat.key")
		env := g.Cfg().GetString("server.Env")
		send := wechart.NewSendMsg(key, g.Cfg().GetBool("wechat.open"))
		send.Send("", fmt.Sprintf("env:%s\n%s", env, str))
	}
}
