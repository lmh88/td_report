package tools

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/spf13/cobra"
	"strings"
	"td_report/app/bean"
	"td_report/app/repo"
	"td_report/common/amazon/s3"
	"td_report/common/file"
	"td_report/common/sendmsg/wechart"
	"td_report/common/tool"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

var S3uploadCheckCmd = &cobra.Command{
	Use:   "s3check",
	Short: "s3check",
	Long:  `s3check检测`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("s3check", false)
		logger.Logger.Info("s3check called", time.Now().Format(vars.TIMEFORMAT))
		if StartDate == "" {
			fmt.Println("开始日期不能为空")
			return
		}

		if ReportType == "" {
			fmt.Println("报表类型不能为空")
			return
		}

		if ReportName == "" {
			fmt.Println("报表名称不能为空")
			return
		}

		if ProfileId == "" && EndDate != "" {
			fmt.Println("不能是全部profile的多个日期对应的目录一起上传，风险比较大")
			return
		}

		S3uploadCheckCmdFunc(StartDate, EndDate, ReportName, ReportType, ProfileId)
	},
}

func init() {
	RootCmd.AddCommand(S3uploadCheckCmd)
}

func S3uploadCheckCmdFunc(startdate, enddate, reportName, reportType, profileId string) error {
	var (
		profileList []string
		err         error
		layout      string
		dayList     []string
	)

	if profileId != "" {
		if strings.Contains(profileId, ",") {
			temp := strings.Split(profileId, ",")
			profileList = temp
		} else {
			profileList = make([]string, 0)
			profileList = append(profileList, profileId)
		}

	} else {

		// 从数据库中获取全部的profile
		if reportType != vars.DSP {
			profileList, err = repo.NewSellerProfileRepository().GetSellerProfileId()
		} else {

			profileList, err = repo.NewProfileRepository().GetDspProfileId()
		}

		if err != nil {
			logger.Logger.Error(err, "=====1，获取profile失败")
			return err
		}
	}

	if reportType == vars.SB && (reportName == vars.BrandMetricsMonthly || reportName == vars.BrandMetricsWeekly) {
		layout = vars.TimeFormatTpl
	} else {
		layout = vars.TimeLayout
	}

	if enddate == "" {
		dayList = make([]string, 0)
		dayList = append(dayList, startdate)
	} else {

		dayList, err = tool.GetDays(startdate, enddate, layout)
		if err != nil {
			logger.Logger.Error(err, "获取天数错误")
			return err
		}
	}

	if reportType == vars.SD {
		for _, item := range dayList {
			for _, tactic := range []string{"T00020", "T00030"} {
				if uploadDataList, err := GetkeyAndPathCheck(item, reportName, reportType, profileList, tactic); err == nil {
					if len(uploadDataList) > 0 {
						check(uploadDataList)
						logger.Logger.Info("send sd ok")
					} else {
						logger.Logger.Info("no data ")
						fmt.Println("no data1")
					}
				} else {
					fmt.Println("no data2")
				}
			}
		}

	} else {

		for _, item := range dayList {
			if uploadDataList, err := GetkeyAndPathCheck(item, reportName, reportType, profileList, ""); err == nil {
				if len(uploadDataList) > 0 {
					check(uploadDataList)
					logger.Logger.Info("send ok")
				} else {
					fmt.Println("no data1")
				}

			} else {
				fmt.Println("no data2")
			}
		}
	}

	return nil
}

func check(filelist []*bean.UploadS3DataCheck) {
	var da = make([]*bean.UploadS3DataCheck, 0)
	client, err := s3.NewClient(g.Cfg().GetString("s3.bucket"))
	if err != nil {
		fmt.Println(err, "s3 create client error")
		return
	}

	for _, item := range filelist {
		if num, err := client.GetFileNum(item.Key); err != nil || num == 0 {
			fmt.Println(err, num, item.Key, item.ProfileId, item.ReportType, item.ReportName, item.StartDate)
			da = append(da, item)
		}
	}

	key := g.Cfg().GetString("wechat.key")
	env := g.Cfg().GetString("server.Env")
	send := wechart.NewSendMsg(key, g.Cfg().GetBool("wechat.open"))
	if len(da) > 0 {
		var str = ""
		var co = 0
		for _, item := range da {
			str = str + fmt.Sprintf("%s:%s:%s:%s:%s", item.ReportType, item.ReportName, item.ProfileId, item.StartDate, item.Key) + "\n"
			co++
			if co > 50 {
				send.Send("", fmt.Sprintf("env:%s\n%s", env, str))
				str = ""
				co = 0
			}
		}

		if str != "" {
			send.Send("", fmt.Sprintf("env:%s\n%s", env, str))
		}

	} else {
		send.Send("", fmt.Sprintf("env:%s\n%s", env, "检测所有的本地文件和s3文件，保持一致"))
	}
}

func GetkeyAndPathCheck(date string, reportName string, reportType string, profileIdList []string, tactic string) ([]*bean.UploadS3DataCheck, error) {
	var (
		fileName  string
		dy        string // 前缀0
		year, w   int
		prefixKey string
		tempstr   string
		daw       string
		path      string
		retdata   = make([]*bean.UploadS3DataCheck, 0)
		filepath  string
	)
	timePeriod := vars.TimeDay
	tempkey := fmt.Sprintf("%s_%s", reportType, reportName)
	s3ReportName, ok := vars.S3ReportMap[tempkey]
	if ok == false {
		return nil, errors.New("对应的报表文件不存在")
	}

	switch reportType {
	case vars.DSP:
		path = fmt.Sprintf("%s/%s/%s", vars.MypathMap[reportType], reportName, date)
		for _, profileId := range profileIdList {
			fileName = fmt.Sprintf("%s_%s.csv", date, profileId)
			prefixKey = fmt.Sprintf("%s/%s/%s/%s/%s", reportType, timePeriod, s3ReportName, date, fileName)
			filepath = fmt.Sprintf("%s/%s", path, fileName)
			if result, _ := file.PathExists(filepath); result {
				temp := &bean.UploadS3DataCheck{
					Key:        prefixKey,
					Path:       filepath,
					ProfileId:  profileId,
					ReportType: reportType,
					ReportName: reportName,
					StartDate:  date,
				}
				retdata = append(retdata, temp)
			}
		}

	case vars.SP:
		path = fmt.Sprintf("%s/%s/%s", vars.MypathMap[reportType], reportName, date)
		for _, profileId := range profileIdList {
			fileName = fmt.Sprintf("%s_%s.gz", date, profileId)
			prefixKey = fmt.Sprintf("%s/%s/%s/%s/%s", reportType, timePeriod, s3ReportName, date, fileName)
			filepath = fmt.Sprintf("%s/%s", path, fileName)
			if result, _ := file.PathExists(filepath); result {
				temp := &bean.UploadS3DataCheck{
					Key:        prefixKey,
					Path:       filepath,
					ProfileId:  profileId,
					ReportType: reportType,
					ReportName: reportName,
					StartDate:  date,
				}
				retdata = append(retdata, temp)
			}
		}

	case vars.SD:
		path = fmt.Sprintf("%s/%s/%s", vars.MypathMap[reportType], reportName, date)
		for _, profileId := range profileIdList {
			fileName = fmt.Sprintf("%s_%s_%s.gz", date, profileId, tactic)
			prefixKey = fmt.Sprintf("%s/%s/%s/%s/%s", reportType, timePeriod, s3ReportName, date, fileName)
			filepath = fmt.Sprintf("%s/%s", path, fileName)
			if result, _ := file.PathExists(filepath); result {
				temp := &bean.UploadS3DataCheck{
					Key:        prefixKey,
					Path:       filepath,
					ProfileId:  profileId,
					ReportType: reportType,
					ReportName: reportName,
					StartDate:  date,
				}
				retdata = append(retdata, temp)
			}
		}

	case vars.SB:
		if reportName == vars.BrandMetricsWeekly || reportName == vars.BrandMetricsMonthly {
			loc, _ := time.LoadLocation(vars.Timezon)
			t, _ := time.ParseInLocation(vars.TimeFormatTpl, date, loc)
			currentTime := t
			currentTimeStr := currentTime.Format(vars.TimeFormatTpl)
			if reportName == vars.BrandMetricsMonthly {
				daw = "M"
				year, w = tool.GetMonth(currentTimeStr)
				timePeriod = vars.TimeMonth
				tempstr = "brand_metrics_month"
			} else {
				daw = "W"
				year, w = tool.GetWeek(currentTimeStr)
				timePeriod = vars.TimeWeek
				tempstr = "brand_metrics_week"
			}

			if w < 10 {
				dy = fmt.Sprintf("0%d", w)
			} else {
				dy = fmt.Sprintf("%d", w)
			}

			strdate := fmt.Sprintf("%d%s%s", year, daw, dy)
			path = fmt.Sprintf("%s/%s/%s/%s", vars.MypathMap[reportType], "bm", timePeriod, strdate)
			for _, profileId := range profileIdList {
				fileName = fmt.Sprintf("%s_%s.json", strdate, profileId)
				prefixKey = fmt.Sprintf("%s/%s/%s/%s/%s", reportType, timePeriod, tempstr, strdate, fileName)
				filepath = fmt.Sprintf("%s/%s", path, fileName)
				if result, _ := file.PathExists(filepath); result {
					temp := &bean.UploadS3DataCheck{
						Key:        prefixKey,
						Path:       filepath,
						ProfileId:  profileId,
						ReportType: reportType,
						ReportName: reportName,
						StartDate:  strdate,
					}
					retdata = append(retdata, temp)
				}
			}

		} else {

			path = fmt.Sprintf("%s/%s/%s", vars.MypathMap[reportType], reportName, date)
			for _, profileId := range profileIdList {
				fileName = fmt.Sprintf("%s_%s.gz", date, profileId)
				prefixKey = fmt.Sprintf("%s/%s/%s/%s/%s", reportType, timePeriod, s3ReportName, date, fileName)
				filepath = fmt.Sprintf("%s/%s", path, fileName)
				if result, _ := file.PathExists(filepath); result {
					temp := &bean.UploadS3DataCheck{
						Key:        prefixKey,
						Path:       filepath,
						ProfileId:  profileId,
						ReportType: reportType,
						ReportName: reportName,
						StartDate:  date,
					}
					retdata = append(retdata, temp)
				}
			}
		}
	}

	return retdata, nil
}
