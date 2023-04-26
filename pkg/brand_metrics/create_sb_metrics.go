package brand_metrics

import (
	"context"
	"td_report/app/bean"
	"td_report/pkg/report"
)

func CreateSbBrandMetricsReportPost(token *bean.ProfileToken, dateStr, enddate, lookBackPeriod string, ctx context.Context) (string, error) {
	return report.SBCreateBrandMetricsReportPost(token, dateStr, enddate, lookBackPeriod, ctx)
}