package adGroups

import (
	"context"
	"td_report/app/bean"
	"td_report/pkg/report"
)

func CreateSbAdGroupsReport(profileToken *bean.ProfileToken, dateStr string, ctx context.Context) (string, error) {
	reportType := "adGroups"
	segment := ""
	creativeType := ""
	metrics := "campaignName,campaignId,campaignStatus,campaignBudget,campaignBudgetType,adGroupName,adGroupId,impressions,clicks,cost,attributedDetailPageViewsClicks14d,attributedSales14d,attributedSales14dSameSKU,attributedConversions14d,attributedConversions14dSameSKU,attributedOrdersNewToBrand14d,attributedOrdersNewToBrandPercentage14d,attributedOrderRateNewToBrand14d,attributedSalesNewToBrand14d,attributedSalesNewToBrandPercentage14d,attributedUnitsOrderedNewToBrand14d,attributedUnitsOrderedNewToBrandPercentage14d,unitsSold14d,dpv14d"
	return report.SBCreateReport(profileToken, reportType, dateStr, segment, creativeType, metrics, ctx)
}

func CreateSbAdGroupsVideoReport(profileToken *bean.ProfileToken, dateStr string, ctx context.Context) (string, error) {
	reportType := "adGroups"
	segment := ""
	creativeType := "video"
	metrics := "campaignName,campaignId,campaignStatus,campaignBudget,campaignBudgetType,adGroupName,adGroupId,impressions,clicks,cost,attributedSales14d,attributedSales14dSameSKU,attributedConversions14d,attributedConversions14dSameSKU,viewableImpressions,videoFirstQuartileViews,videoMidpointViews,videoThirdQuartileViews,videoCompleteViews,video5SecondViews,video5SecondViewRate,videoUnmutes,vtr,vctr"
	return report.SBCreateReport(profileToken, reportType, dateStr, segment, creativeType, metrics, ctx)
}
