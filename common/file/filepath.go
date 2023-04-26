package file

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"os"
	"td_report/vars"
	"time"
)

// GetWeek 获取当前时间是一年中第几周
func GetWeek(datetime string) (y, w int) {
	loc, _ := time.LoadLocation(vars.Timezon)
	tmp, _ := time.ParseInLocation(vars.TimeFormatTpl, datetime, loc)
	return tmp.ISOWeek()
}

// GetMonth 获取当前时间是一年中第几周
func GetMonth(datetime string) (y, w int) {
	loc, _ := time.LoadLocation(vars.Timezon)
	tmp, _ := time.ParseInLocation(vars.TimeFormatTpl, datetime, loc)
	return tmp.Year(), int(tmp.Month())
}

func GetPath(reportType string, reportName string, dateStr string) string {
	var (
		path      string
		timetype  string
		daw       string
		dy        string // 前缀0
		year, w   int
		reportDir string
	)
	if reportType == vars.DSP || reportType == vars.SD || reportType == vars.SP {
		reportDir = vars.PathMap[reportType]
		path = fmt.Sprintf("%s/%s/%s", g.Cfg().GetString("common.datapath")+reportDir, reportName, dateStr)
	} else {

		if reportName == vars.BrandMetricsWeekly || reportName == vars.BrandMetricsMonthly {
			loc, _ := time.LoadLocation(vars.Timezon)
			h, _ := time.ParseInLocation(vars.TimeFormatTpl, dateStr, loc)
			currentTime := h
			currentTimeStr := currentTime.Format(vars.TimeFormatTpl)
			if reportName == vars.BrandMetricsMonthly {
				timetype = "month"
				daw = "M"
				year, w = GetMonth(currentTimeStr)
			} else {
				timetype = "week"
				daw = "W"
				year, w = GetWeek(currentTimeStr)
			}

			if w < 10 {
				dy = fmt.Sprintf("0%d", w)
			} else {
				dy = fmt.Sprintf("%d", w)
			}

			path = fmt.Sprintf("%s/%s/%s/%d%s%s", g.Cfg().GetString("common.datapath")+vars.SbPath, "bm", timetype, year, daw, dy)
		} else {
			path = fmt.Sprintf("%s/%s/%s", g.Cfg().GetString("common.datapath")+vars.SbPath, reportName, dateStr)
		}
	}

	return path
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	//isnotexist来判断，是不是不存在的错误
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false, nil
	}
	return false, err //如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
}
