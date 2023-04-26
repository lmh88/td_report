package tool

import (
	"errors"
	"fmt"
	"math"
	"os"
	"td_report/app/bean"
	"td_report/app/repo"
	"td_report/common/file"
	"td_report/common/tool"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"

	"github.com/google/wire"
)

type StatisticsService struct {
	SellerProfileRepository  *repo.SellerProfileRepository
	ProfileRepository        *repo.ProfileRepository
	ReportSdCheckRepository  *repo.ReportSdCheckRepository
	ReportSpCheckRepository  *repo.ReportSpCheckRepository
	ReportSbCheckRepository  *repo.ReportSbCheckRepository
	ReportDspCheckRepository *repo.ReportDspCheckRepository
}

var StatisticsServiceSet = wire.NewSet(wire.Struct(new(StatisticsService), "*"))

func NewStatisticsService(
	sellerProfileRepository *repo.SellerProfileRepository,
	profileRepository *repo.ProfileRepository,
	ReportSdCheckRepository *repo.ReportSdCheckRepository,
	ReportSpCheckRepository *repo.ReportSpCheckRepository,
	ReportSbCheckRepository *repo.ReportSbCheckRepository,
	ReportDspCheckRepository *repo.ReportDspCheckRepository) *StatisticsService {
	return &StatisticsService{
		SellerProfileRepository:  sellerProfileRepository,
		ProfileRepository:        profileRepository,
		ReportSdCheckRepository:  ReportSdCheckRepository,
		ReportSpCheckRepository:  ReportSpCheckRepository,
		ReportSbCheckRepository:  ReportSbCheckRepository,
		ReportDspCheckRepository: ReportDspCheckRepository,
	}
}

// 针对普通类型报表的情况
func (t *StatisticsService) statis(path, fileExt, date string, profileIdList []string, check bool) []*bean.Result {
	var (
		filename string
		fileList = make([]*bean.Result, 0)
	)

	if check == true {
		for _, ig := range []string{"T00020", "T00030"} {
			for _, item := range profileIdList {
				filename = fmt.Sprintf("%s_%s_%s.%s", date, item, ig, fileExt)
				detailfile := &bean.Result{
					Exists:    false,
					StartDate: date,
					ProfileId: item,
					Extrant:   ig,
					FileName:  filename,
				}

				fileDetailpath := fmt.Sprintf("%s/%s", path, filename)
				info, err := os.Stat(fileDetailpath)
				if err != nil { // 文件不存在
					fileList = append(fileList, detailfile)
					continue
				} else {
					detailfile.FileName = info.Name()
					detailfile.ModifileDate = info.ModTime().Format(vars.TIMEFORMAT)
					detailfile.Exists = true
					fileList = append(fileList, detailfile)
				}
			}
		}

	} else {

		for _, item := range profileIdList {
			filename = fmt.Sprintf("%s_%s.%s", date, item, fileExt)
			detailfile := &bean.Result{
				Exists:    false,
				StartDate: date,
				ProfileId: item,
				FileName:  filename,
			}

			fileDetailpath := fmt.Sprintf("%s/%s", path, filename)
			info, err := os.Stat(fileDetailpath)
			if err != nil { // 文件不存在
				fileList = append(fileList, detailfile)
				continue
			} else {
				detailfile.FileName = info.Name()
				detailfile.ModifileDate = info.ModTime().Format(vars.TIMEFORMAT)
				detailfile.Exists = true
				fileList = append(fileList, detailfile)
			}
		}
	}

	return fileList
}

func (t *StatisticsService) GetPath(reportType string, reportName string, date string) string {
	return file.GetPath(reportType, reportName, date)
}

// 1 sb的特殊情况 2 普通的情况
func (t *StatisticsService) checkdata(reportType, reportName string) int {
	if reportType == vars.SB && (reportName == vars.BrandMetricsWeekly || reportName == vars.BrandMetricsMonthly) {
		return 1
	} else {
		return 2
	}
}

func (t *StatisticsService) Getperfied(start, end string, format string) ([]string, error) {
	return tool.GetDays(start, end, format)
}

// GetDate 根据时间格式校验时间并且转化时间对象
func (t *StatisticsService) GetDate(dateFormat string, startDate, endDate string) (time.Time, time.Time, error) {
	var (
		startT, endT time.Time
		err          error
	)

	startT, err = time.Parse(dateFormat, startDate)
	if err != nil {
		return startT, endT, errors.New("开始日期格式错误")
	}
	endT, err = time.Parse(dateFormat, endDate)
	if err != nil {
		return startT, endT, errors.New("结束日期格式错误")
	}

	if endT.Before(startT) {
		return startT, endT, errors.New("开始日期落后于结束日期，格式错误")
	}

	return startT, endT, nil
}

// GetDateFormat 获取时间格式
func (t *StatisticsService) GetDateFormat(reportName, reportType string) string {
	var dateFormat string
	if t.checkdata(reportType, reportName) == 1 {
		dateFormat = vars.TimeFormatTpl
	} else {
		dateFormat = vars.TimeLayout
	}
	return dateFormat
}

func (t *StatisticsService) getdata(startT time.Time, endT time.Time, reportType string, reportName string, startDate string, fileExt string, profileIdList []string, dateFormat string) (map[string][]*bean.Result, error) {
	var (
		checkResultData = make(map[string][]*bean.Result, 0)
	)

	if len(profileIdList) == 0 {
		return nil, errors.New("no profileid")
	}

	timeResult := endT.Sub(startT)
	d := timeResult.Hours() / 24
	otherCheck := false
	if reportType == vars.SD {
		otherCheck = true
	}
	if d < 1 {
		path := t.GetPath(reportType, reportName, startDate)
		result := t.statis(path, fileExt, startDate, profileIdList, otherCheck)
		checkResultData[startDate] = append(checkResultData[startDate], result...)

	} else {

		day := int(math.Ceil(d))
		for i := 0; i <= day; i++ {
			gg := i * 24
			dd, _ := time.ParseDuration(fmt.Sprintf("%dh", gg))
			start := startT.Add(dd).Format(dateFormat)
			path := t.GetPath(reportType, reportName, start)
			result := t.statis(path, fileExt, start, profileIdList, otherCheck)
			checkResultData[start] = append(checkResultData[start], result...)
		}

	}

	return checkResultData, nil
}

// 获取报表文件名后缀
func (t *StatisticsService) getfileext(reportType, reportName string) string {
	if reportName != vars.BrandMetricsMonthly && reportName != vars.BrandMetricsWeekly {
		return vars.ReportListFileExt[reportType]
	} else {
		return vars.ReportListFileExt[vars.SB_BRAND]
	}
}

// GetFileByProfile 统计服务器上面单个profile在某一个时间段内的情况
func (t *StatisticsService) GetFileByProfile(reportName, reportType, startDate string, endDate string, profileList []string) (map[string][]*bean.Result, error) {
	var (
		startT, endT time.Time
		err          error
		dateFormat   string
	)

	dateFormat = t.GetDateFormat(reportName, reportType)
	startT, endT, err = t.GetDate(dateFormat, startDate, endDate)
	if err != nil {
		return nil, err
	}
	fileExt := t.getfileext(reportType, reportName)
	return t.getdata(startT, endT, reportType, reportName, startDate, fileExt, profileList, dateFormat)
}

// GetFileWithProfileList 统计服务器上系列profile在某一个时间段内的情况
func (t *StatisticsService) GetFileWithProfileList(reportName, reportType, startDate string, endDate string) (map[string][]*bean.Result, error) {
	var (
		startT, endT  time.Time
		err           error
		dateFormat    string
		profileIdList = make([]string, 0)
	)

	fileExt := t.getfileext(reportType, reportName)
	dateFormat = t.GetDateFormat(reportName, reportType)
	startT, endT, err = t.GetDate(dateFormat, startDate, endDate)

	if err != nil {
		logger.Logger.Info(err)
		return nil, err
	}

	if reportType == vars.DSP {
		dspProfile, err := t.ProfileRepository.ListDspRegionProfile(profileIdList)
		if err != nil {
			logger.Logger.Error(err)
			return nil, err
		}

		for _, item := range dspProfile {
			profileIdList = append(profileIdList, item.ProfileId)
		}

		return t.getdata(startT, endT, reportType, reportName, startDate, fileExt, profileIdList, dateFormat)

	} else {

		var (
			profileList []*bean.ProfileToken
		)

		if t.checkdata(reportType, reportName) == 1 {
			profileList, err = t.SellerProfileRepository.GetProfile(2)
		} else {
			profileList, err = t.SellerProfileRepository.GetProfile(1)
		}

		if err != nil {
			logger.Logger.Error(err)
			return nil, err
		} else {

			for _, item := range profileList {
				profileIdList = append(profileIdList, item.ProfileId)
			}

			return t.getdata(startT, endT, reportType, reportName, startDate, fileExt, profileIdList, dateFormat)
		}
	}
}
