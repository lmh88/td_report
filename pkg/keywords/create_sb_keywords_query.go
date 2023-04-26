package keywords

import (
	"context"
	"td_report/app/bean"
	"td_report/pkg/report"
)

func CreateSbKeywordsQueryReport(profileToken *bean.ProfileToken, dateStr string, ctx context.Context) (string, error) {
	reportType := "keywords"
	segment := "query"
	creativeType := ""
	metrics := "campaignName,campaignId,campaignStatus,campaignBudget,campaignBudgetType,adGroupName,adGroupId,keywordText,keywordBid,keywordStatus,searchTermImpressionRank,matchType,impressions,clicks,cost,attributedSales14d,attributedConversions14d"
	return report.SBCreateReport(profileToken, reportType, dateStr, segment, creativeType, metrics, ctx)
}

func CreateSbKeywordsQueryVideoReport(profileToken *bean.ProfileToken, dateStr string, ctx context.Context) (string, error) {
	reportType := "keywords"
	segment := "query"
	creativeType := "video"
	metrics := "campaignName,campaignId,campaignStatus,campaignBudget,campaignBudgetType,adGroupName,adGroupId,keywordText,keywordBid,keywordStatus,matchType,impressions,clicks,cost,attributedSales14d,attributedConversions14d,viewableImpressions,videoFirstQuartileViews,videoMidpointViews,videoThirdQuartileViews,videoCompleteViews,video5SecondViews,video5SecondViewRate,videoUnmutes,vtr,vctr"
	return report.SBCreateReport(profileToken, reportType, dateStr, segment, creativeType, metrics, ctx)
}
