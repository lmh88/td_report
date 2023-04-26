package keywords

import (
	"context"
	"td_report/app/bean"
	"td_report/pkg/report"
)

func CreateSPKeywordsReport(profileToken *bean.ProfileToken, dateStr string, ctx context.Context) (string, error) {
	reportType := "keywords"
	segment := ""
	metrics := "campaignName,campaignId,adGroupName,adGroupId,keywordId,keywordText,matchType,impressions,clicks,cost,attributedConversions1d,attributedConversions7d,attributedConversions14d,attributedConversions30d,attributedConversions1dSameSKU,attributedConversions7dSameSKU,attributedConversions14dSameSKU,attributedConversions30dSameSKU,attributedUnitsOrdered1d,attributedUnitsOrdered7d,attributedUnitsOrdered14d,attributedUnitsOrdered30d,attributedSales1d,attributedSales7d,attributedSales14d,attributedSales30d,attributedSales1dSameSKU,attributedSales7dSameSKU,attributedSales14dSameSKU,attributedSales30dSameSKU,attributedUnitsOrdered1dSameSKU,attributedUnitsOrdered7dSameSKU,attributedUnitsOrdered14dSameSKU,attributedUnitsOrdered30dSameSKU"
	return report.SPCreateReport(profileToken, reportType, dateStr, segment, metrics, ctx)
}

func CreateSPKeywords(profileToken *bean.ProfileToken, dateStr string, ctx context.Context) (string, error) {
	reportType := "keywords"
	segment := ""
	metrics := "campaignName,campaignId,adGroupName,adGroupId,keywordId,keywordText,matchType,impressions,clicks,cost,attributedConversions1d,attributedConversions7d,attributedConversions14d,attributedConversions30d,attributedConversions1dSameSKU,attributedConversions7dSameSKU,attributedConversions14dSameSKU,attributedConversions30dSameSKU,attributedUnitsOrdered1d,attributedUnitsOrdered7d,attributedUnitsOrdered14d,attributedUnitsOrdered30d,attributedSales1d,attributedSales7d,attributedSales14d,attributedSales30d,attributedSales1dSameSKU,attributedSales7dSameSKU,attributedSales14dSameSKU,attributedSales30dSameSKU,attributedUnitsOrdered1dSameSKU,attributedUnitsOrdered7dSameSKU,attributedUnitsOrdered14dSameSKU,attributedUnitsOrdered30dSameSKU"
	return report.SPCreateReportV1(profileToken, reportType, dateStr, segment, metrics, ctx)
}
