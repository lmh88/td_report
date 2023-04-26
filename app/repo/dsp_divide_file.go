package repo

import "strings"

func GetDspFileFileds(reportName string) []string {
	ordersFileds := strings.Split("date,entityId,advertiserName,advertiserId,orderName,orderId,orderStartDate,orderEndDate,orderBudget,orderExternalId,"+
		"orderCurrency,newToBrandPurchasesClicks14d,totalCost,amazonAudienceFee,atcClicks14d,"+
		"amazonPlatformFee,purchases14d,newToBrandPurchasesViews14d,atlViews14d,supplyCost,newSubscribeAndSave14d,newSubscribeAndSaveViews14d,pRPVViews14d,pRPV14d,newToBrandPurchases14d,totalPixel14d,"+
		"dpvViews14d,newSubscribeAndSaveClicks14d,atlClicks14d,purchasesViews14d,impressions,atc14d,purchasesClicks14d,atcViews14d,"+
		"pRPVClicks14d,dpv14d,totalPixelClicks14d,clickThroughs,atl14d,dpvClicks14d,totalPixelViews14d,unitsSold14d,totalUnitsSold14d,"+
		"totalDetailPageViews14d,totalPurchases14d,totalAddToListViews14d,measurableImpressions,newToBrandUnitsSold14d,totalSubscribeAndSaveSubscriptionClicks14d,totalAddToCartViews14d,"+
		"totalSubscribeAndSaveSubscriptionViews14d,totalNewToBrandPurchasesClicks14d,totalAddToCartClicks14d,totalFee,3PFees,"+
		"totalPRPVViews14d,viewableImpressions,totalAddToListClicks14d,totalAddToCart14d,totalNewToBrandPurchasesViews14d,"+
		"newToBrandProductSales14d,totalAddToList14d,totalPurchasesViews14d,totalSales14d,totalPRPVClicks14d,agencyFee,"+
		"totalPRPV14d,totalNewToBrandUnitsSold14d,totalPurchasesClicks14d,totalSubscribeAndSaveSubscriptions14d,totalNewToBrandProductSales14d,"+
		"totalDetailPageViewViews14d,totalDetailPageClicks14d,totalNewToBrandPurchases14d,sales14d", ",")
	detailFileds := strings.Split("date,entityId,advertiserName,advertiserId,orderName,orderId,orderStartDate,orderEndDate,orderBudget,orderExternalId,orderCurrency,"+
		"lineItemName,lineItemId,lineItemStartDate,lineItemEndDate,lineItemBudget,lineItemExternalId,creativeName,creativeID,creativeAdId,creativeType,"+
		"creativeSize,totalCost,dpv14d,clickThroughs,impressions,totalAddToCart14d,totalNewToBrandProductSales14d,totalUnitsSold14d,totalDetailPageViews14d,totalPurchases14d,totalSales14d,sales14d", ",")
	inventoryFileds := strings.Split("date,entityId,advertiserName,advertiserId,orderName,orderId,orderStartDate,orderEndDate,orderBudget,orderExternalId,orderCurrency,lineItemName,lineItemId,"+
		"lineItemStartDate,lineItemEndDate,lineItemBudget,lineItemExternalId,siteName,placementSize,placementName,supplySourceName,totalCost,dpv14d,supplyCost,clickThroughs,"+
		"impressions,totalAddToCart14d,totalNewToBrandProductSales14d,totalUnitsSold14d,totalDetailPageViews14d,totalPurchases14d,totalSales14d,sales14d", ",")
	audienceFileds := strings.Split("intervalStart,intervalEnd,entityId,advertiserName,advertiserId,orderName,orderId,orderStartDate,orderEndDate,orderBudget,orderExternalId,orderCurrency,lineItemName,lineItemId,"+
		"lineItemStartDate,lineItemEndDate,lineItemBudget,lineItemExternalId,segment,segmentMarketplaceID,segmentSource,segmentType,purchases14d,totalCost,dpv14d,clickThroughs,impressions,atc14d,lineitemtype", ",")

	switch reportName {
	case "detail":
		return detailFileds
	case "order":
		return ordersFileds
	case "inventory":
		return inventoryFileds
	case "audience":
		return audienceFileds
	default:
		return nil
	}
}
