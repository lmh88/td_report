package tools

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"github.com/spf13/cobra"
	"td_report/common/file"
	"td_report/common/tool"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

// GetProfileDataCmd 给定profileid和日期范围，下载一定时间段的报表数据到本地，避免每一个都复制下来
var GetProfileDataCmd = &cobra.Command{
	Use:   "get_profile",
	Short: "get_profile",
	Long:  `get_profile`,
	Run: func(cmd *cobra.Command, args []string) {
		if StartDate == "" || EndDate == "" {
			fmt.Println("开始日期或者结束日期错误")
			return
		}
		logger.Init("get_profile", false)
		logger.Logger.Info("get_profile called")
		getFunc(ProfileId, ReportType, ReportName, StartDate, EndDate)
	},
}

func init() {
	RootCmd.AddCommand(GetProfileDataCmd)
}

func getFunc(ProfileId, ReportType, ReportName, StartDate, EndDate string) {
	daylist, err := tool.GetDays(StartDate, EndDate, vars.TimeLayout)
	if err != nil {
		logger.Logger.Error(err.Error())
		fmt.Println("获取时间格式错误")
		return
	}

	var (
		path         string
		uploadPrefix = g.Cfg().GetString("common.uploadpath")
		dstpath      string
	)

	if !gfile.Exists(uploadPrefix) {
		gfile.Mkdir(uploadPrefix)
	}

	date := time.Now().Format(vars.TimeLayout)
	uploaPath := uploadPrefix + "/" + date
	if gfile.Exists(uploaPath) {
		gfile.Remove(uploaPath)
	}

	gfile.Mkdir(uploaPath)

	for _, item := range daylist {
		if ReportType == vars.SB && (ReportName == vars.BrandMetricsWeekly || ReportName == vars.BrandMetricsMonthly) {
			path = file.GetPath(ReportType, ReportName, item)
		} else {
			path = fmt.Sprintf("%s/%s/%s", vars.MypathMap[ReportType], ReportName, item)
		}

		if gfile.Exists(path) == false {

		} else {

			if ReportType == vars.SD {
				for _, tactic := range []string{"T00020", "T00030"} {
					filenameList := getfilename(ReportType, ReportName, ProfileId, item, tactic)
					for _, filename := range filenameList {
						fullpath := fmt.Sprintf("%s/%s", path, filename)
						if gfile.Exists(fullpath) == true {
							dstpath = uploaPath + "/" + filename
							gfile.Copy(fullpath, dstpath)
						}
					}
				}

			} else {

				filenameList := getfilename(ReportType, ReportName, ProfileId, item, "")
				for _, filename := range filenameList {
					fullpath := fmt.Sprintf("%s/%s", path, filename)
					if gfile.Exists(fullpath) == true {
						dstpath = uploaPath + "/" + filename
						gfile.Copy(fullpath, dstpath)
					}
				}
			}
		}
	}
}
