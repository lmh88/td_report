package productAds

import (
	"context"
	"td_report/app/service/report_v2/create_report"
	"td_report/app/service/report_v2/varible"
)

func CreateSdProductAdsReport(profileMsg *varible.ProfileMsg, tactic string, ctx context.Context) (string, error) {
	reportType := "productAds"
	commonMetrics := "campaignName,campaignId,adGroupName,adGroupId,asin,adId,impressions,clicks,cost,currency,attributedConversions1d,attributedConversions7d,attributedConversions14d,attributedConversions30d,attributedConversions1dSameSKU,attributedConversions7dSameSKU,attributedConversions14dSameSKU,attributedConversions30dSameSKU,attributedUnitsOrdered1d,attributedUnitsOrdered7d,attributedUnitsOrdered14d,attributedUnitsOrdered30d,attributedSales1d,attributedSales7d,attributedSales14d,attributedSales30d,attributedSales1dSameSKU,attributedSales7dSameSKU,attributedSales14dSameSKU,attributedSales30dSameSKU,attributedOrdersNewToBrand14d,attributedSalesNewToBrand14d,attributedUnitsOrderedNewToBrand14d,viewImpressions,viewAttributedConversions14d,viewAttributedSales14d,viewAttributedUnitsOrdered14d,attributedDetailPageView14d,viewAttributedDetailPageView14d"
	var metrics string
	if profileMsg.ProfileType == varible.Seller {
		metrics = commonMetrics + ",sku"
	} else {
		metrics = commonMetrics
	}
	return create_report.SDCreateReport(profileMsg, reportType, tactic, metrics, ctx)
}
