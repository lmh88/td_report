package asins

import (
	"context"
	"td_report/app/bean"
	"td_report/pkg/report"
)

func CreateSdAsinsReport(token *bean.ProfileToken, dateStr, tactic string, ctx context.Context) (string, error) {
	reportType := "asins"
	metrics := "campaignName,campaignId,adGroupName,adGroupId,asin,otherAsin,sku,currency,attributedUnitsOrdered1dOtherSKU,attributedUnitsOrdered7dOtherSKU,attributedUnitsOrdered14dOtherSKU,attributedUnitsOrdered30dOtherSKU,attributedSales1dOtherSKU,attributedSales7dOtherSKU,attributedSales14dOtherSKU,attributedSales30dOtherSKU"
	return report.SDCreateReport(token, reportType, dateStr, tactic, metrics, ctx)
}
