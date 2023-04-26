package audience

import (
	"context"
	"td_report/pkg/dsp/dsp_report"
)

func CreateAudience(region, dateStr, profileId string, ctx context.Context) (string, error) {
	return dsp_report.CreateDspReport(region, "audience", dateStr, profileId, ctx)
}

func CreateAudiencePeriod(region, startDate, endDate, profileId string, ctx context.Context) (string, error) {
	return dsp_report.CreateDspReportPeriod(region, "audience", startDate, endDate, profileId, ctx)
}
