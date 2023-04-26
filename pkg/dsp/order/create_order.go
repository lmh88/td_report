package order

import (
	"context"
	"td_report/pkg/dsp/dsp_report"
)

func CreateOrder(region, dateStr, profileId string, ctx context.Context) (string, error) {
	return dsp_report.CreateDspReport(region, "order", dateStr, profileId, ctx)
}

func CreateOrderPeriod(region, startDate, endDate, profileId string, ctx context.Context) (string, error) {
	return dsp_report.CreateDspReportPeriod(region, "order", startDate, endDate, profileId, ctx)
}
