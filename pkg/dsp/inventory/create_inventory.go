package inventory

import (
	"context"
	"td_report/pkg/dsp/dsp_report"
)

func CreateInventory(region, dateStr, profileId string, ctx context.Context) (string, error) {
	return dsp_report.CreateDspReport(region, "inventory", dateStr, profileId, ctx)
}

func CreateInventoryPeriod(region, startDate, endDate, profileId string, ctx context.Context) (string, error) {
	return dsp_report.CreateDspReportPeriod(region, "inventory", startDate, endDate, profileId, ctx)
}
