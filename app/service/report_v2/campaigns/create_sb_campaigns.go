package campaigns

import (
	"context"
	"td_report/app/service/report_v2/create_report"
	"td_report/app/service/report_v2/varible"
)

func CreateSbCampaignsReport(profileMsg *varible.ProfileMsg, ctx context.Context) (string, error) {
	reportType := "campaigns"
	segment := "placement"
	creativeType := "all"
	metrics := "campaignName,campaignId,campaignStatus,campaignBudget,campaignBudgetType,campaignRuleBasedBudget,applicableBudgetRuleId,applicableBudgetRuleName,impressions,clicks,cost,attributedDetailPageViewsClicks14d,attributedSales14d,attributedSales14dSameSKU,attributedConversions14d,attributedConversions14dSameSKU,attributedOrdersNewToBrand14d,attributedOrdersNewToBrandPercentage14d,attributedOrderRateNewToBrand14d,attributedSalesNewToBrand14d,attributedSalesNewToBrandPercentage14d,attributedUnitsOrderedNewToBrand14d,attributedUnitsOrderedNewToBrandPercentage14d,unitsSold14d,dpv14d"

	return create_report.SBCreateReport(profileMsg, reportType, segment, creativeType, metrics, ctx)
}

func CreateSbCampaignsVideoReport(profileMsg *varible.ProfileMsg, ctx context.Context) (string, error) {
	reportType := "campaigns"
	segment := "placement"
	creativeType := "video"
	metrics := "campaignName,campaignId,campaignStatus,campaignBudget,campaignBudgetType,impressions,clicks,cost,attributedSales14d,attributedSales14dSameSKU,attributedConversions14d,attributedConversions14dSameSKU,viewableImpressions,videoFirstQuartileViews,videoMidpointViews,videoThirdQuartileViews,videoCompleteViews,video5SecondViews,video5SecondViewRate,videoUnmutes,vtr,vctr"
	return create_report.SBCreateReport(profileMsg, reportType, segment, creativeType, metrics, ctx)
}
