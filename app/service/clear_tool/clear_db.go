package clear_tool

import (
	"fmt"
	"td_report/app/repo"
	"td_report/vars"
	"time"
)

const (
	ReportError = "report_error"
	ReportDetail = "report_detail"
)

type clearFunc func(date time.Time) error

var DbMap = map[string]clearFunc{
	ReportError: repo.NewReportErrorRepository().ClearByDate,
	ReportDetail: clearDetail,
}

func ClearDb(clearType string, days int) error {
	dateStr := time.Now().Format(vars.TimeFormatTpl)
	date, _ := time.Parse(vars.TimeFormatTpl, dateStr)
	sub, _ := time.ParseDuration(fmt.Sprintf("-%dh", days * 24))
	date = date.Add(sub)

	err := DbMap[clearType](date)
	return err
}

func clearDetail(date time.Time) error {
	err := repo.NewReportDspConsumerDetailRepository().ClearByDate(date)
	if err != nil {
		return err
	}
	err = repo.NewReportSdConsumerDetailRepository().ClearByDate(date)
	if err != nil {
		return err
	}
	err = repo.NewReportSbConsumerDetailRepository().ClearByDate(date)
	if err != nil {
		return err
	}
	err = repo.NewReportSpConsumerDetailRepository().ClearByDate(date)
	if err != nil {
		return err
	}
	return nil
}

