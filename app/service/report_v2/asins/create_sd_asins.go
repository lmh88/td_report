package asins

import (
	"context"
	"td_report/app/service/report_v2/create_report"
	"td_report/app/service/report_v2/varible"
)

func CreateSdAsinsReport(profileMsg *varible.ProfileMsg, tactic string, ctx context.Context) (string, error) {
	reportType := "asins"
	commonMetrics := "campaignName,campaignId,adGroupName,adGroupId,asin,otherAsin,currency,attributedUnitsOrdered1dOtherSKU,attributedUnitsOrdered7dOtherSKU,attributedUnitsOrdered14dOtherSKU,attributedUnitsOrdered30dOtherSKU,attributedSales1dOtherSKU,attributedSales7dOtherSKU,attributedSales14dOtherSKU,attributedSales30dOtherSKU"
	var metrics string
	if profileMsg.ProfileType == varible.Seller {
		metrics = commonMetrics + ",sku"
	} else {
		metrics = commonMetrics
	}
	return create_report.SDCreateReport(profileMsg, reportType, tactic, metrics, ctx)
}
