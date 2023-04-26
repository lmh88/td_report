package productAds

import (
	"context"
	"td_report/app/service/report_v2/create_report"
	"td_report/app/service/report_v2/varible"
)

func CreateSPProductAdsReport(profileMsg *varible.ProfileMsg, ctx context.Context) (string, error) {
	reportType := "productAds"
	segment := ""
	commonMetrics := "campaignName,campaignId,adGroupName,adGroupId,impressions,clicks,cost,currency,asin,attributedConversions1d,attributedConversions7d,attributedConversions14d,attributedConversions30d,attributedConversions1dSameSKU,attributedConversions7dSameSKU,attributedConversions14dSameSKU,attributedConversions30dSameSKU,attributedUnitsOrdered1d,attributedUnitsOrdered7d,attributedUnitsOrdered14d,attributedUnitsOrdered30d,attributedSales1d,attributedSales7d,attributedSales14d,attributedSales30d,attributedSales1dSameSKU,attributedSales7dSameSKU,attributedSales14dSameSKU,attributedSales30dSameSKU,attributedUnitsOrdered1dSameSKU,attributedUnitsOrdered7dSameSKU,attributedUnitsOrdered14dSameSKU,attributedUnitsOrdered30dSameSKU"
	var metrics string
	if profileMsg.ProfileType == varible.Seller {
		metrics = commonMetrics + ",sku"
	} else {
		metrics = commonMetrics
	}
	return create_report.SPCreateReport(profileMsg, reportType, segment, metrics, ctx)
}
