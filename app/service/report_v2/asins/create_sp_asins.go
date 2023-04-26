package asins

import (
	"context"
	"td_report/app/service/report_v2/create_report"
	"td_report/app/service/report_v2/varible"
)

func CreateSPAsinsReport(profileMsg *varible.ProfileMsg, ctx context.Context) (string, error) {
	reportType := "asins"
	segment := ""
	commonMetrics := "campaignName,campaignId,adGroupName,adGroupId,asin,otherAsin,currency,matchType,attributedUnitsOrdered1d,attributedUnitsOrdered7d,attributedUnitsOrdered14d,attributedUnitsOrdered30d,attributedUnitsOrdered1dOtherSKU,attributedUnitsOrdered7dOtherSKU,attributedUnitsOrdered14dOtherSKU,attributedUnitsOrdered30dOtherSKU,attributedSales1dOtherSKU,attributedSales7dOtherSKU,attributedSales14dOtherSKU,attributedSales30dOtherSKU,targetingText,targetingType"
	var metrics string
	if profileMsg.ProfileType == varible.Seller {
		metrics = commonMetrics + ",sku"
	} else {
		metrics = commonMetrics
	}
	return create_report.SPCreateReport(profileMsg, reportType, segment, metrics, ctx)
}
