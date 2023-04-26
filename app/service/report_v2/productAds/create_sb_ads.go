package productAds

import (
	"context"
	"td_report/app/service/report_v2/create_report"
	"td_report/app/service/report_v2/varible"
)

func CreateSbAdsReport(profileMsg *varible.ProfileMsg, ctx context.Context) (string, error) {
	reportType := "ads"
	segment := ""
	creativeType := "all"
	metrics := "adGroupId,adGroupName,adId,applicableBudgetRuleId,applicableBudgetRuleName,attributedConversions14d,attributedConversions14dSameSKU,attributedDetailPageViewsClicks14d,attributedOrderRateNewToBrand14d,attributedOrdersNewToBrand14d,attributedOrdersNewToBrandPercentage14d,attributedSales14d,attributedSales14dSameSKU,attributedSalesNewToBrand14d,attributedSalesNewToBrandPercentage14d,attributedUnitsOrderedNewToBrand14d,attributedUnitsOrderedNewToBrandPercentage14d,campaignBudget,campaignBudgetType,campaignId,campaignName,campaignRuleBasedBudget,campaignStatus,clicks,cost,dpv14d,impressions,unitsSold14d,vctr,video5SecondViewRate,video5SecondViews,videoCompleteViews,videoFirstQuartileViews,videoMidpointViews,videoThirdQuartileViews,videoUnmutes,viewableImpressions,vtr,attributedBrandedSearches14d"
	return create_report.SBCreateReport(profileMsg, reportType, segment, creativeType, metrics, ctx)
}
