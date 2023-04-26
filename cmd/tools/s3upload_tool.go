package tools

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
	"td_report/app/bean"
	"td_report/app/repo"
	"td_report/boot"
	"td_report/common/file"
	"td_report/common/tool"
	"td_report/pkg/logger"
	"td_report/pkg/save_file"
	"td_report/vars"
	"time"
)

// S3uploadToolCmd  上传指定指定报表类型，指定日期文件到s3,
var S3uploadToolCmd = &cobra.Command{
	Use:   "s3tool",
	Short: "s3tool",
	Long: `s3tool,暂时只支持指定的报表类型，报表名称，开始日期，结束日期，
和指定的profile或者profile列表上传，最大限度的是单个报表单个日期的整个目录上传，避免多个目录一起上传，如果是多个profile，则用逗号隔开`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Init("s3tool", false)
		logger.Logger.Info("s3tool called", time.Now().Format(vars.TIMEFORMAT))
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

		s3toolFunc(StartDate, EndDate, ReportName, ReportType, ProfileId)
	},
}

func init() {
	RootCmd.AddCommand(S3uploadToolCmd)
}

func Send(uploaddataList []*bean.UploadS3Data) error {
	if client, err := boot.GetRabbitmqClient(); err != nil {
		return err
	} else {

		for _, uploaddata := range uploaddataList {
			err = client.Send(vars.UploadS3, uploaddata)
			if err != nil {
				logger.Logger.Error("push data to queue error ", err.Error())
				save_file.Saveerrordata(uploaddata)
			}

			time.Sleep(2 * time.Second)
		}

		return nil
	}
}

func s3toolFunc(startdate, enddate, reportName, reportType, profileId string) error {
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
				if uploadDataList, err := GetkeyAndPath(item, reportName, reportType, profileList, tactic); err == nil {
					if len(uploadDataList) > 0 {
						Send(uploadDataList)
						logger.Logger.Info("send sd ok")
					} else {
						logger.Logger.Info("no data ")
					}
				} else {
					logger.Logger.Error(err, "send sd error")
				}
			}
		}

	} else {

		for _, item := range dayList {
			if uploadDataList, err := GetkeyAndPath(item, reportName, reportType, profileList, ""); err == nil {
				if len(uploadDataList) > 0 {
					Send(uploadDataList)
					logger.Logger.Info("send ok")
				} else {
					logger.Logger.Info("no data")
				}

			} else {
				logger.Logger.Error(err, "send")
			}
		}
	}

	return nil
}

//func PathExists(path string) (bool, error) {
//	_, err := os.Stat(path)
//	if err == nil {
//		return true, nil
//	}
//	//isnotexist来判断，是不是不存在的错误
//	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
//		return false, nil
//	}
//	return false, err //如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
//}

// GetkeyAndPath 获取对应的key和path
func GetkeyAndPath(date string, reportName string, reportType string, profileIdList []string, tactic string) ([]*bean.UploadS3Data, error) {
	var (
		fileName  string
		dy        string // 前缀0
		year, w   int
		prefixKey string
		tempstr   string
		daw       string
		path      string
		retdata   []*bean.UploadS3Data = make([]*bean.UploadS3Data, 0)
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
				temp := &bean.UploadS3Data{
					Key:  prefixKey,
					Path: filepath,
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
				temp := &bean.UploadS3Data{
					Key:  prefixKey,
					Path: filepath,
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
				temp := &bean.UploadS3Data{
					Key:  prefixKey,
					Path: filepath,
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
					temp := &bean.UploadS3Data{
						Key:  prefixKey,
						Path: filepath,
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
					temp := &bean.UploadS3Data{
						Key:  prefixKey,
						Path: filepath,
					}
					retdata = append(retdata, temp)
				}
			}
		}
	}

	return retdata, nil
}
