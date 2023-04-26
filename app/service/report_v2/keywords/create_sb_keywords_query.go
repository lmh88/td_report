package keywords

import (
	"context"
	"td_report/app/service/report_v2/create_report"
	"td_report/app/service/report_v2/varible"
)

func CreateSbKeywordsQueryReport(profileMsg *varible.ProfileMsg, ctx context.Context) (string, error) {
	reportType := "keywords"
	segment := "query"
	creativeType := "all"
	metrics := "campaignName,campaignId,campaignStatus,campaignBudget,campaignBudgetType,adGroupName,adGroupId,keywordText,keywordBid,keywordStatus,searchTermImpressionRank,matchType,impressions,clicks,cost,attributedSales14d,attributedConversions14d"
	return create_report.SBCreateReport(profileMsg, reportType, segment, creativeType, metrics, ctx)
}

func CreateSbKeywordsQueryVideoReport(profileMsg *varible.ProfileMsg, ctx context.Context) (string, error) {
	reportType := "keywords"
	segment := "query"
	creativeType := "video"
	metrics := "campaignName,campaignId,campaignStatus,campaignBudget,campaignBudgetType,adGroupName,adGroupId,keywordText,keywordBid,keywordStatus,matchType,impressions,clicks,cost,attributedSales14d,attributedConversions14d,viewableImpressions,videoFirstQuartileViews,videoMidpointViews,videoThirdQuartileViews,videoCompleteViews,video5SecondViews,video5SecondViewRate,videoUnmutes,vtr,vctr"
	return create_report.SBCreateReport(profileMsg, reportType, segment, creativeType, metrics, ctx)
}
