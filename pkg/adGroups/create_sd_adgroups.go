package adGroups

import (
	"context"
	"td_report/app/bean"
	"td_report/pkg/report"
)

func CreateSdAdGroupsReport(profileToken *bean.ProfileToken, dateStr, tactic string, ctx context.Context) (string, error) {
	reportType := "adGroups"
	metrics := "campaignName,campaignId,adGroupName,adGroupId,impressions,clicks,cost,currency,attributedConversions1d,attributedConversions7d,attributedConversions14d,attributedConversions30d,attributedConversions1dSameSKU,attributedConversions7dSameSKU,attributedConversions14dSameSKU,attributedConversions30dSameSKU,attributedUnitsOrdered1d,attributedUnitsOrdered7d,attributedUnitsOrdered14d,attributedUnitsOrdered30d,attributedSales1d,attributedSales7d,attributedSales14d,attributedSales30d,attributedSales1dSameSKU,attributedSales7dSameSKU,attributedSales14dSameSKU,attributedSales30dSameSKU,attributedOrdersNewToBrand14d,attributedSalesNewToBrand14d,attributedUnitsOrderedNewToBrand14d,bidOptimization,viewImpressions,viewAttributedConversions14d,viewAttributedSales14d"
	return report.SDCreateReport(profileToken, reportType, dateStr, tactic, metrics, ctx)
}
