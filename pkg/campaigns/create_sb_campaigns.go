package campaigns

import (
	"context"
	"td_report/app/bean"
	"td_report/pkg/report"
)

func CreateSbCampainsReport(token *bean.ProfileToken, dateStr string, ctx context.Context) (string, error) {
	reportType := "campaigns"
	segment := "placement"
	creativeType := ""
	metrics := "campaignName,campaignId,campaignStatus,campaignBudget,campaignBudgetType,campaignRuleBasedBudget,applicableBudgetRuleId,applicableBudgetRuleName,impressions,clicks,cost,attributedDetailPageViewsClicks14d,attributedSales14d,attributedSales14dSameSKU,attributedConversions14d,attributedConversions14dSameSKU,attributedOrdersNewToBrand14d,attributedOrdersNewToBrandPercentage14d,attributedOrderRateNewToBrand14d,attributedSalesNewToBrand14d,attributedSalesNewToBrandPercentage14d,attributedUnitsOrderedNewToBrand14d,attributedUnitsOrderedNewToBrandPercentage14d,unitsSold14d,dpv14d"

	return report.SBCreateReport(token, reportType, dateStr, segment, creativeType, metrics, ctx)
}

func CreateSbCampainsVideoReport(token *bean.ProfileToken, dateStr string, ctx context.Context) (string, error) {
	reportType := "campaigns"
	segment := "placement"
	creativeType := "video"
	metrics := "campaignName,campaignId,campaignStatus,campaignBudget,campaignBudgetType,impressions,clicks,cost,attributedSales14d,attributedSales14dSameSKU,attributedConversions14d,attributedConversions14dSameSKU,viewableImpressions,videoFirstQuartileViews,videoMidpointViews,videoThirdQuartileViews,videoCompleteViews,video5SecondViews,video5SecondViewRate,videoUnmutes,vtr,vctr"
	return report.SBCreateReport(token, reportType, dateStr, segment, creativeType, metrics, ctx)
}
