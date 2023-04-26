package detail

import (
	"context"
	"td_report/pkg/dsp/dsp_report"
)

func CreateDetail(region, dateStr, profileId string, ctx context.Context) (string, error) {
	return dsp_report.CreateDspReport(region, "detail", dateStr, profileId, ctx)
}

func CreateDetailPeriod(region, startDate, endDate, profileId string, ctx context.Context) (string, error) {
	return dsp_report.CreateDspReportPeriod(region, "detail", startDate, endDate, profileId, ctx)
}
