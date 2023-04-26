package save_file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"td_report/app/bean"
	"td_report/boot"
	"td_report/common/redis"
	"td_report/common/tool"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

//const (
//	TimeDay   = "day"
//	TimeWeek  = "week"
//	TimeMonth = "month"
//)

func GetSbPath(reportName string, dateStr string) string {
	var (
		path     string
		daw      string
		dy       string // 前缀0
		year, w  int
		timetype string
	)

	if reportName == vars.BrandMetricsWeekly || reportName == vars.BrandMetricsMonthly {
		loc, _ := time.LoadLocation(vars.Timezon)
		t, _ := time.ParseInLocation(vars.TimeFormatTpl, dateStr, loc)
		currentTime := t
		currentTimeStr := currentTime.Format(vars.TimeFormatTpl)
		if reportName == vars.BrandMetricsMonthly {
			timetype = vars.TimeMonth
			daw = "M"
			year, w = tool.GetMonth(currentTimeStr)
		} else {
			timetype = vars.TimeWeek
			daw = "W"
			year, w = tool.GetWeek(currentTimeStr)
		}

		if w < 10 {
			dy = fmt.Sprintf("0%d", w)
		} else {
			dy = fmt.Sprintf("%d", w)
		}

		path = fmt.Sprintf("%s/%s/%s/%d%s%s", vars.MypathMap[vars.SB], "bm", timetype, year, daw, dy)
	} else {
		path = fmt.Sprintf("%s/%s/%s", vars.MypathMap[vars.SB], reportName, dateStr)
	}

	return path
}

func GetSbPathWithS3(reportName string, dateStr string, profileId string) (string, string, string) {
	var (
		path       string
		daw        string
		dy         string // 前缀0
		year, w    int
		prefixKey  string
		reportType string = vars.SB
		timePeriod string
		fileName   string
		tempstr    string
	)

	if reportName == vars.BrandMetricsWeekly || reportName == vars.BrandMetricsMonthly {
		loc, _ := time.LoadLocation(vars.Timezon)
		t, _ := time.ParseInLocation(vars.TimeFormatTpl, dateStr, loc)
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
		fileName = fmt.Sprintf("%s_%s.json", strdate, profileId)
		prefixKey = fmt.Sprintf("%s/%s/%s/%s/%s", reportType, timePeriod, tempstr, strdate, fileName)
	} else {

		timePeriod = vars.TimeDay
		path = fmt.Sprintf("%s/%s/%s", vars.MypathMap[reportType], reportName, dateStr)
		fileName = fmt.Sprintf("%s_%s.gz", dateStr, profileId)
		tempkey := fmt.Sprintf("%s_%s", reportType, reportName)
		s3ReportName, ok := vars.S3ReportMap[tempkey]
		if ok == false {
			panic("the reportname not exists")
		}
		prefixKey = fmt.Sprintf("%s/%s/%s/%s/%s", reportType, timePeriod, s3ReportName, dateStr, fileName)
	}

	return path, prefixKey, fileName
}

// SaveSBFile
func SaveSBFile(reportName, dateStr, profileId string, data []byte) error {
	path, prefixKey, fileName := GetSbPathWithS3(reportName, dateStr, profileId)
	return sacveFileCommon(path, fileName, data, prefixKey)
}

func SaveSPFile(reportName, dateStr, profileId string, data []byte) error {
	reportType := vars.SP
	timePeriod := vars.TimeDay
	path := fmt.Sprintf("%s/%s/%s", vars.MypathMap[reportType], reportName, dateStr)
	fileName := fmt.Sprintf("%s_%s.gz", dateStr, profileId)
	tempkey := fmt.Sprintf("%s_%s", reportType, reportName)
	s3ReportName, ok := vars.S3ReportMap[tempkey]
	if ok == false {
		panic("the reportname not exists")
	}
	prefixKey := fmt.Sprintf("%s/%s/%s/%s/%s", reportType, timePeriod, s3ReportName, dateStr, fileName)
	return sacveFileCommon(path, fileName, data, prefixKey)
}

func SaveDspFile(reportName, dateStr, profileId string, data []byte) error {
	reportType := vars.DSP
	timePeriod := vars.TimeDay
	path := fmt.Sprintf("%s/%s/%s", vars.MypathMap[reportType], reportName, dateStr)
	fileName := fmt.Sprintf("%s_%s.csv", dateStr, profileId)
	tempkey := fmt.Sprintf("%s_%s", reportType, reportName)
	s3ReportName, ok := vars.S3ReportMap[tempkey]
	if ok == false {
		panic("the reportname not exists")
	}
	prefixKey := fmt.Sprintf("%s/%s/%s/%s/%s", reportType, timePeriod, s3ReportName, dateStr, fileName)
	return sacveFileCommon(path, fileName, data, prefixKey)
}

func SaveDspFilePeriod(reportName, startdate, enddate, profileId string, data []byte) error {
	reportType := vars.DSP
	timePeriod := vars.TimeDay
	path := fmt.Sprintf("%s/%s/%s", vars.MypathMap[reportType], reportName, startdate)
	fileName := fmt.Sprintf("%s_%s_%s.csv", startdate, enddate, profileId)
	tempkey := fmt.Sprintf("%s_%s", reportType, reportName)
	s3ReportName, ok := vars.S3ReportMap[tempkey]
	if ok == false {
		panic("the reportname not exists")
	}
	prefixKey := fmt.Sprintf("%s/%s/%s/%s/%s", reportType, timePeriod, s3ReportName, startdate, fileName)
	return sacveFileCommon(path, fileName, data, prefixKey)
}

func SaveSDFile(reportName, dateStr, profileId, tactic string, data []byte) error {
	reportType := vars.SD
	timePeriod := vars.TimeDay
	path := fmt.Sprintf("%s/%s/%s", vars.MypathMap[reportType], reportName, dateStr)
	fileName := fmt.Sprintf("%s_%s_%s.gz", dateStr, profileId, tactic)
	tempkey := fmt.Sprintf("%s_%s", reportType, reportName)
	s3ReportName, ok := vars.S3ReportMap[tempkey]
	if ok == false {
		panic("the reportname not exists")
	}
	prefixKey := fmt.Sprintf("%s/%s/%s/%s/%s", reportType, timePeriod, s3ReportName, dateStr, fileName)
	return sacveFileCommon(path, fileName, data, prefixKey)
}

var uploaderror = redis.WithUploadFile("uploaderr")

func sacveFileCommon(path string, fileName string, data []byte, prekey string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			return err
		}
	}

	filepath := fmt.Sprintf("%s/%s", path, fileName)
	err = ioutil.WriteFile(filepath, data, 0644)
	if err != nil {
		return err
	} else {

		uploaddata := &bean.UploadS3Data{
			Key:  prekey,
			Path: filepath,
		}

		err = boot.RabitmqClient.Send(vars.UploadS3, uploaddata)
		if err != nil {
			logger.Logger.Error("push data to queue error ", err.Error())
			Saveerrordata(uploaddata)
		}

		return err
	}
}

// Saveerrordata 存redis， 是入列rabbitmq出错了或者中间环节出错，可能rabbitmq存储有问题，避免当前时刻数据都存入rabbitmq
func Saveerrordata(uploaddata *bean.UploadS3Data) {
	redisClient := boot.RedisCommonClient.GetClient()
	pipe := redisClient.Pipeline()
	uploadJson, _ := json.Marshal(uploaddata)
	pipe.SAdd(uploaderror, uploadJson)
	pipe.Expire(uploaderror, 24*time.Hour)
	pipe.Exec()
	pipe.Close()
}
