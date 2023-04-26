package tool

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"td_report/app/repo"
	"td_report/common/sendmsg/wechart"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

type MonitorErrorServer struct {
	ErrorNotice string
}

func NewMonitorErrorServer() *MonitorErrorServer {
	instance := MonitorErrorServer{
		ErrorNotice: g.Cfg().GetString("server.Env", "local") + "\n" + "monitor_report_error:",
	}
	return &instance
}

func (s *MonitorErrorServer) Run() {
	year, month, day := time.Now().Date()
	date := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	date = date.AddDate(0, 0, -1)
	s.ErrorNotice += date.Format(vars.TimeLayout) + "\n"

	reportTypes := []string{vars.SP, vars.SB, vars.SD, vars.DSP}
	errorTypes := []int{repo.ReportErrorTypeOne, repo.ReportErrorTypeTwo, repo.ReportErrorTypeThree, repo.ReportErrorTypeTimeOut, repo.ReportErrorTypeRetry}

	for _, reportType := range reportTypes {
		for _, errorType := range errorTypes {
			total, pTotal, mProfile, _ := s.getErrorNum(reportType, date, errorType)
			s.ErrorNotice += fmt.Sprintf("%s,error_type:%d,total:%d,profile_total:%d,many_profile:%d\n", reportType, errorType, total, pTotal, mProfile)
		}
	}
	s.send()
}

func (s *MonitorErrorServer) send() {
	key := "0ac12699-c90a-4fd8-9bf6-5fb75cba6459"
	logger.Info(s.ErrorNotice)
	send := wechart.NewSendMsg(key, true)
	send.Send("", s.ErrorNotice)
}


func (s *MonitorErrorServer) getErrorNum(reportType string, date time.Time, errorType int) (int, int, int, error) {
	var (
		total, profileTotal, manyProfile int
		profileNums []repo.ProfileNum
		err error
	)

	profileNums, err = repo.NewReportErrorRepository().GetProfileNum(reportType, date, errorType)

	if err != nil {
		return total, profileTotal, manyProfile, err
	}

	if len(profileNums) > 0 {
		for _, item := range  profileNums {
			total += item.Num
			profileTotal++
			if item.Num > 4 {
				manyProfile++
			}
		}
	}

	return total, profileTotal, manyProfile, err
}

