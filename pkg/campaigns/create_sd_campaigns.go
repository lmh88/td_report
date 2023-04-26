package campaigns

import (
	"context"
	"td_report/app/bean"
	"td_report/pkg/report"
)

func CreateSdCampaignsReport(token *bean.ProfileToken, dateStr, tactic string, ctx context.Context) (string, error) {
	reportType := "campaigns"
	metrics := "campaignName,campaignId,impressions,clicks,cost,currency,attributedConversions1d,attributedConversions7d,attributedConversions14d,attributedConversions30d,attributedConversions1dSameSKU,attributedConversions7dSameSKU,attributedConversions14dSameSKU,attributedConversions30dSameSKU,attributedUnitsOrdered1d,attributedUnitsOrdered7d,attributedUnitsOrdered14d,attributedUnitsOrdered30d,attributedSales1d,attributedSales7d,attributedSales14d,attributedSales30d,attributedSales1dSameSKU,attributedSales7dSameSKU,attributedSales14dSameSKU,attributedSales30dSameSKU,attributedOrdersNewToBrand14d,attributedSalesNewToBrand14d,attributedUnitsOrderedNewToBrand14d,costType,viewImpressions,viewAttributedConversions14d,viewAttributedSales14d,viewAttributedUnitsOrdered14d,attributedDetailPageView14d,viewAttributedDetailPageView14d"
	return report.SDCreateReport(token, reportType, dateStr, tactic,  metrics, ctx)
}
