package asins

import (
	"context"
	"td_report/app/bean"
	"td_report/pkg/report"
)

func CreateSPAsinsReport(profileToken *bean.ProfileToken, dateStr string, ctx context.Context) (string, error) {
	reportType := "asins"
	segment := ""
	metrics := "campaignName,campaignId,adGroupName,adGroupId,asin,otherAsin,sku,currency,matchType,attributedUnitsOrdered1d,attributedUnitsOrdered7d,attributedUnitsOrdered14d,attributedUnitsOrdered30d,attributedUnitsOrdered1dOtherSKU,attributedUnitsOrdered7dOtherSKU,attributedUnitsOrdered14dOtherSKU,attributedUnitsOrdered30dOtherSKU,attributedSales1dOtherSKU,attributedSales7dOtherSKU,attributedSales14dOtherSKU,attributedSales30dOtherSKU,targetingText,targetingType"
	return report.SPCreateReport(profileToken, reportType, dateStr, segment,  metrics, ctx)
}
