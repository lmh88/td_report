package tools

import (
	"fmt"
	"github.com/gogf/gf/os/gfile"
	"github.com/spf13/cobra"
	"strings"
	"td_report/app/bean"
	"td_report/common/file"
	"td_report/common/tool"
	"td_report/pkg/logger"
	"td_report/vars"
)

// CheckProfileDataCmd 给定profileid和日期范围，检测报表生成的时间直接输出，如果没有生成，输出日期
var CheckProfileDataCmd = &cobra.Command{
	Use:   "check_profile",
	Short: "check_profile",
	Long:  `check_profile`,
	Run: func(cmd *cobra.Command, args []string) {
		if StartDate == "" || EndDate == "" {
			fmt.Println("开始日期或者结束日期错误")
			return
		}
		logger.Init("check_profile", false)
		logger.Logger.Info("check_profile called")
		checkFunc(ProfileId, ReportType, ReportName, StartDate, EndDate)
	},
}

func init() {
	RootCmd.AddCommand(CheckProfileDataCmd)
}


func checkFunc(ProfileId, ReportType, ReportName, StartDate, EndDate string) {
	daylist, err := tool.GetDays(StartDate, EndDate, vars.TimeLayout)
	if err != nil {
		logger.Logger.Error(err.Error())
		fmt.Println("获取时间格式错误")
		return
	}

	var (
		path    string
		output  = fmt.Sprintf("%s:%s:%s\n", ReportType, ReportName, ProfileId)
		data = make([]bean.Ckeckdata,0)
	)

	for _, item := range daylist {
		var temp bean.Ckeckdata
		if ReportType == vars.SB && (ReportName == vars.BrandMetricsWeekly || ReportName == vars.BrandMetricsMonthly) {
			path = file.GetPath(ReportType, ReportName, item)
		} else {
			path = fmt.Sprintf("%s/%s/%s", vars.MypathMap[ReportType], ReportName, item)
		}

		if gfile.Exists(path) == false {
			temp = bean.Ckeckdata{
				ProfileId: ProfileId,
				ErrorType: 1,
				ReportDate: item,
			}

			data = append(data, temp)

		} else {


			if ReportType == vars.SD {
				for _, tactic := range []string{"T00020", "T00030"} {
					filenameList := getfilename(ReportType, ReportName, ProfileId, item, tactic)
					for _, filename := range filenameList {
						fullpath := fmt.Sprintf("%s/%s", path, filename)
						if gfile.Exists(fullpath) == false {
							temp = bean.Ckeckdata{
								ProfileId: ProfileId,
								ErrorType: 1,
								Extrant: tactic,
								ReportDate: item,
							}
						} else {
							temp = bean.Ckeckdata{
								ProfileId: ProfileId,
								ErrorType: 0,
								Extrant: tactic,
								FileUpdate: gfile.MTime(fullpath).Format(vars.TIMEFORMAT),
								ReportDate: item,
							}
						}

						data = append(data, temp)
					}
				}

			} else {

				filenameList := getfilename(ReportType, ReportName, ProfileId, item, "")
				for _, filename := range filenameList {
					fullpath := fmt.Sprintf("%s/%s", path, filename)
					if gfile.Exists(fullpath) == false {
						temp = bean.Ckeckdata{
							ProfileId: ProfileId,
							ErrorType: 1,
							ReportDate: item,
						}
					} else {
						temp = bean.Ckeckdata{
							ProfileId: ProfileId,
							ErrorType: 0,
							FileUpdate: gfile.MTime(fullpath).Format(vars.TIMEFORMAT),
							ReportDate: item,
						}
					}

					data = append(data, temp)
				}
			}
		}
	}

	for _, myitem:= range data {
       output += fmt.Sprintf("reportDate:%s,isempty:%d, upatetime:%s, (扩展类型:%s)\n", myitem.ReportDate, myitem.ErrorType, myitem.FileUpdate,myitem.Extrant)
	}

	fmt.Println(output)
}

func getfilename(reportType, reportName, profileId, date,  tactic string) []string {
	var fileName = make([]string, 0)
	switch reportType {
	case vars.DSP:
		fileName = append(fileName, fmt.Sprintf("%s_%s.csv", date, profileId))
	case vars.SP:
		fileName = append(fileName, fmt.Sprintf("%s_%s.gz", date, profileId))
	case vars.SD:
		fileName = append(fileName, fmt.Sprintf("%s_%s_%s.gz", date, profileId, tactic))

	case vars.SB:
		if reportType == vars.SB && (reportName == vars.BrandMetricsWeekly || reportName == vars.BrandMetricsMonthly) {
			path := file.GetPath(reportType, reportName, date)
			index := strings.LastIndex(path, "/")
			prefix := path[index+1:]
			temp := fmt.Sprintf("%s_%s.json", prefix, profileId)
			fileName = append(fileName, temp)
		} else {
			fileName = append(fileName, fmt.Sprintf("%s_%s.gz", date, profileId))
		}
	}

	return fileName
}
