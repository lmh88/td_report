package dsp_report

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"td_report/boot"
	"td_report/pkg/dbp_token"
	"td_report/pkg/limiter"
	"td_report/pkg/logger"
	region2 "td_report/pkg/region"
	"td_report/pkg/requests"
	"td_report/vars"
	"time"
)

const (
	REPORT_TYPE_ORDER     = "order"
	REPORT_TYPE_DETAIL    = "detail"
	REPORT_TYPE_AUDIENCE  = "audience"
	REPORT_TYPE_INVENTORY = "inventory"
	REPORT_TYPE_PRODUCT   = "product"
)

var dimensions = map[string][]string{
	REPORT_TYPE_ORDER:     {"ORDER"},
	REPORT_TYPE_DETAIL:    {"ORDER", "LINE_ITEM", "CREATIVE"},
	REPORT_TYPE_AUDIENCE:  {"ORDER", "LINE_ITEM"},
	REPORT_TYPE_INVENTORY: {"ORDER", "LINE_ITEM", "SITE", "SUPPLY"},
	REPORT_TYPE_PRODUCT:   {"ORDER", "LINE_ITEM"},
}

var metrics = map[string]string{
	REPORT_TYPE_ORDER:     "totalSales14d,clickThroughs,dpv14d,impressions,totalAddToCart14d,totalNewToBrandProductSales14d,sales14d,totalDetailPageViews14d,totalCost,totalPurchases14d,totalUnitsSold14d,supplyCost,amazonPlatformFee,amazonAudienceFee,agencyFee,3PFees,totalFee,viewableImpressions,measurableImpressions,totalPixel14d,totalPixelViews14d,totalPixelClicks14d,dpvViews14d,dpvClicks14d,atc14d,atcViews14d,atcClicks14d,purchases14d,purchasesViews14d,purchasesClicks14d,newToBrandPurchases14d,newToBrandPurchasesViews14d,newToBrandPurchasesClicks14d,unitsSold14d,newToBrandUnitsSold14d,newToBrandProductSales14d,totalDetailPageViewViews14d,totalDetailPageClicks14d,totalAddToCartViews14d,totalAddToCartClicks14d,totalPurchasesViews14d,totalPurchasesClicks14d,totalNewToBrandPurchases14d,totalNewToBrandPurchasesViews14d,totalNewToBrandPurchasesClicks14d,totalNewToBrandUnitsSold14d,totalPRPV14d,totalPRPVViews14d,totalPRPVClicks14d,totalAddToList14d,totalAddToListViews14d,totalAddToListClicks14d,totalSubscribeAndSaveSubscriptions14d,totalSubscribeAndSaveSubscriptionViews14d,totalSubscribeAndSaveSubscriptionClicks14d,newSubscribeAndSave14d,newSubscribeAndSaveViews14d,newSubscribeAndSaveClicks14d,pRPV14d,pRPVViews14d,pRPVClicks14d,atl14d,atlViews14d,atlClicks14d",
	REPORT_TYPE_DETAIL:    "impressions,clickThroughs,dpv14d,totalDetailPageViews14d,totalAddToCart14d,totalPurchases14d,totalSales14d,totalUnitsSold14d,totalNewToBrandProductSales14d,totalCost,sales14d",
	REPORT_TYPE_INVENTORY: "impressions,clickThroughs,dpv14d,totalDetailPageViews14d,totalAddToCart14d,totalPurchases14d,totalSales14d,totalUnitsSold14d,totalNewToBrandProductSales14d,totalCost,sales14d,supplyCost,placementName,placementSize",
	REPORT_TYPE_AUDIENCE:  "impressions,clickThroughs,dpv14d,atc14d,purchases14d,totalCost,segmentSource,segmentType,segmentMarketplaceID,lineitemtype",
	REPORT_TYPE_PRODUCT:   "productName,productGroup,productCategory,productSubcategory,dpv14d,dpvViews14d,dpvClicks14d,pRPV14d,pRPVViews14d,pRPVClicks14d,atl14d,atlViews14d,atlClicks14d,atc14d,atcViews14d,atcClicks14d,purchases14d,purchasesViews14d,purchasesClicks14d,newToBrandPurchases14d,newToBrandPurchasesViews14d,newToBrandPurchasesClicks14d,percentOfPurchasesNewToBrand14d,totalDetailPageViews14d,totalDetailPageViewViews14d,totalDetailPageClicks14d,totalPRPV14d,totalPRPVViews14d,totalPRPVClicks14d,totalAddToList14d,totalAddToListViews14d,totalAddToListClicks14d,totalAddToCart14d,totalAddToCartViews14d,totalAddToCartClicks14d,totalPurchases14d,totalPurchasesViews14d,totalPurchasesClicks14d,totalNewToBrandPurchases14d,totalNewToBrandPurchasesViews14d,totalNewToBrandPurchasesClicks14d,totalPercentOfPurchasesNewToBrand14d,newSubscribeAndSave14d,newSubscribeAndSaveViews14d,newSubscribeAndSaveClicks14d,totalSubscribeAndSaveSubscriptions14d,totalSubscribeAndSaveSubscriptionViews14d,totalSubscribeAndSaveSubscriptionClicks14d,unitsSold14d,sales14d,totalUnitsSold14d,totalSales14d,newToBrandUnitsSold14d,newToBrandProductSales14d,brandHaloDetailPage14d,brandHaloDetailPageViews14d,brandHaloDetailPageClicks14d,brandHaloProductReviewPage14d,brandHaloProductReviewPageViews14d,brandHaloProductReviewPageClicks14d,brandHaloAddToList14d,brandHaloAddToListViews14d,brandHaloAddToListClicks14d,brandHaloAddToCart14d,brandHaloAddToCartViews14d,brandHaloAddToCartClicks14d,brandHaloPurchases14d,brandHaloPurchasesViews14d,brandHaloPurchasesClicks14d,brandHaloNewToBrandPurchases14d,brandHaloNewToBrandPurchasesViews14d,brandHaloNewToBrandPurchasesClicks14d,brandHaloPercentOfPurchasesNewToBrand14d,brandHaloNewSubscribeAndSave14d,brandHaloNewSubscribeAndSaveViews14d,brandHaloNewSubscribeAndSaveClicks14d,brandHaloTotalUnitsSold14d,brandHaloTotalSales14d,brandHaloTotalNewToBrandSales14d,brandHaloTotalNewToBrandUnitsSold14d",
}

var DspRefreshToken = map[string]string{
	region2.RegionFe: "Atzr|IwEBIJqR9kHQts1y3pFLrTeLMA8afkEd_dNUAwsjSjqIhFEZHCmVGyhiOkSqZ1gNHUY2X1TSLJjhUeRm6PVe1FnzBAo7o6bOmE495OjbGKaOQAJ2fHh1rBS51dJLT4oX1l7-bWJYb2npLXmOjZQgY8t1ZCDR-45KY0SqzjrbtQMDjNh8uAKRiwhBwBu3H13FcJxxOBi8yjrS2zKuNb2wmWBTe9vUhPbigbvvAMzMeTFmMSK-_N5hd8ZglvJjkvobK5SD7BQeDXK49t5BN-05WSzR6kLfFa9gWSSYphoPdF9hUYaUHIDBUj4ZOeWURoFugZzDO408Xc0lee6GGz1CP21aumiWlSCTuEUrOAENGcNaipEEvWRvZDiHKhIX546ycrVt95ujrQgIlVm5H2-copN3Pkb1KAUBSmW7Ylgwha6OB5fJNg",
	region2.RegionEu: "Atzr|IwEBIJ9_AwBc6Pzv65zeD4ph1nrxNqceokx8IrqxRktMCrpeuC0ZdlqGqUjD8tO3TUCoD8804QB4ND-1Ng-9bY5gMbay9O69vgYN4c2A36y-ep_jzRF4kjyuU22ALAnZjjj6FQ0SoRTkstAZUd7zB3z7U2_i6SIRyaXpDhhIPSAleYEBlQAyM48Lph_VW5dU8pNL7OQrh9r1SW05yyuB7lYaA1juQDYdX4U8WRez0iLssUXhTbCSqTUSSFuqBedAzJMcge84iI_EYp6C9wtFFDFJasXmEfOeLA4S2x_OS07PW06nvXjL3EeLDxgCZqGRTHSVj3V0BbsypQdjYgD7chXeZpAs-YqJKE1dZQ4lPt6yomeJzFFmtO_3RNO1H7_II-nKW3OSHC3WjCG5V9on7_vR1sEI4PZmD0d1RUEEfOL2liTp5ThjRq809SnV9YZL38g9As8",
	region2.RegionNa: "Atzr|IwEBIFohGsLBoxln39hOsfSZfYN18O1YJINunLjrMwBv28y6rpH7EtK8PC19orgSR6TzQE5IkWVOB-e-QMIDCYvDcTCxg9Zqndo9pHLovbPLmXYPdt04DRyEZKgptREAhq8GmZvkZww8ldYm4nGnpFOpV7Czoy8CAjmZHQF9B-jzQhOzjxTKhCUeM2M-6kojkHLmJQqDs3n8QiZr4RVscL3Y41ZLvYwh8D5H1Z93gQNuIdsFzBq0bookPCWJuJeQjC484654ox5AMLfaKlQEejUZQNDdB-Dn6zuvEGudeDPbFh69cANaA7EVHUd4XSl6SP3l0QrW87qLCWz5iXnEnDATVa4-MyGCiMVE-zt4zxxkVJXyDo8UZ8zaFS3qAUdlHOJZ_nGHFHJN7XreJ-cw3foIRLFR1Q6TBYbt6nnmdhFwgj4K2OcePo-Z2T_tvDSIAuLdYn4",
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func fakeHeaders(profileId, refreshToken string) (map[string]string, error) {
	clientId:="amzn1.application-oa2-client.084663234c2143c3a3bf91fe34bbdf1e"
	clientSecret:="29f982cfd2571585c52db5ba462f302716c87fd52507dcc355582b8821a80abc"
	accessToken, err := dbp_token.GetAccessTokenWithClient(refreshToken, clientId, clientSecret)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Amazon-Advertising-API-Scope":    profileId,
		"Amazon-Advertising-API-ClientId": clientId,
		"Authorization":                   accessToken,
	}

	return headers, nil
}

func CreateDspReport(region, reportType, dateStr, profileId string, ctx context.Context) (string, error) {

	endpoint, ok := region2.ApiUrl[region]
	if !ok {
		return "", fmt.Errorf("region not found: %s", region)
	}

	url := endpoint + "/dsp/reports"

	refreshToken, ok := DspRefreshToken[region]
	if !ok {
		return "", fmt.Errorf("region token not found: %s", region)
	}

	headers, err := fakeHeaders(profileId, refreshToken)
	if err != nil {
		return "", err
	}

	params := map[string]interface{}{
		"format":     "CSV",
		"startDate":  dateStr,
		"endDate":    dateStr,
		"dimensions": dimensions[reportType],
		"metrics":    metrics[reportType],
	}

	// 根据报表类型设置参数
	switch reportType {
	case REPORT_TYPE_AUDIENCE:
		params["type"] = "AUDIENCE"
		params["timeUnit"] = "SUMMARY"
	case REPORT_TYPE_INVENTORY:
		params["type"] = "INVENTORY"
		params["timeUnit"] = "DAILY"
	case REPORT_TYPE_PRODUCT:
		params["type"] = "PRODUCTS"
		params["timeUnit"] = "DAILY"
	case REPORT_TYPE_DETAIL, REPORT_TYPE_ORDER:
		params["type"] = "CAMPAIGN"
		params["timeUnit"] = "DAILY"
	default:
		return "", fmt.Errorf("reportType not found:%s", reportType)
	}

	count := 0
	var resp *requests.Resp
	var rlimit = limiter.PerSecond(vars.LimitMap["dsp"])
	var sleepSec int
	for {

		for {

			res, err := boot.Rlimiter.Allow(url, rlimit)
			if err != nil {
				return "", err
			}

			if res.Allowed == 1 {
				break
			}

			time.Sleep(10 * time.Millisecond)
		}

		count += 1
		if count > 40 {
			logger.Logger.InfoWithContext(ctx, map[string]interface{}{
				"desc":       " get reportId error ppc get sleep too long to break",
				"count":      count,
				"reportName": reportType,
				"reportType": "dsp",
			})
			break
		}

		resp, err = requests.Post(url, requests.WithHeaders(headers), requests.WithJson(params), requests.WithTimeout(time.Minute))
		if err != nil || resp == nil {
			sleepSec = rand.Intn(30) + 1 + count*20
			if resp != nil {
				errmap := map[string]interface{}{"code": resp.StatusCode, "body": string(resp.Body)}
				logger.Logger.ErrorWithContextMap(ctx, errmap, "dsp http post create report error")
			}
			if err != nil {
				errstr := err.Error()
				errmap := map[string]interface{}{"err": errstr, "reportType": reportType}
				//超时休眠重试
				if strings.Contains(errstr, vars.Timeouttxt) || strings.Contains(errstr, vars.Timeouttxt1) {
					logger.Logger.ErrorWithContextMap(ctx, errmap, "dsp http post create report error")
					time.Sleep(time.Duration(sleepSec))
					continue
				} else {
					logger.Logger.ErrorWithContextMap(ctx, errmap, "dsp http post create report error")
				}
			}

			// 当前如果网络不好或者抖动，则添加休眠，避免中间持续频繁请求导致恶化
			time.Sleep(time.Duration(sleepSec))
			return "", err
		}

		logger.Logger.InfoWithContext(ctx, fmt.Sprintf("dsp报表第一次请求结果，statusCode:%d, body:%s", resp.StatusCode, string(resp.Body)))
		if resp.StatusCode == http.StatusTooManyRequests {
			sleepSec = rand.Intn(30) + 1 + count*20
			logger.Logger.InfoWithContext(ctx, map[string]interface{}{
				"flag":     "StatusTooManyRequests",
				"count":    count,
				"sleepSec": sleepSec,
			})

			time.Sleep(time.Second * time.Duration(sleepSec))
			continue
		}

		if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				logger.Logger.ErrorWithContext(ctx, fmt.Sprintf("get status resp error, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body)))
				return "", fmt.Errorf("resp error, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body))
			} else {
				sleepSec = rand.Intn(30) + 1 + count*20
				logger.Logger.InfoWithContext(ctx, map[string]interface{}{
					"status":  resp.StatusCode,
					"sleep":   sleepSec,
					"num":     count,
					"moudule": "get dsp report reportid",
					"desc":    "dsp http post create report error",
					"params":  params,
					"body":    string(resp.Body),
				})

				time.Sleep(time.Second * time.Duration(sleepSec))
				continue
			}

		} else {

			val, err := resp.Json()
			if err != nil {
				logger.Logger.ErrorWithContext(ctx, "dsp json error")
				return "", err
			}

			reportId, ok := val["reportId"].(string)
			if !ok {
				logger.Logger.InfoWithContext(ctx, map[string]interface{}{
					"desc": "dsp create reportId not found",
				})
				return "", fmt.Errorf("reportId not found, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body))
			} else {
				return reportId, nil
			}
		}
	}

	logger.Logger.InfoWithContext(ctx, map[string]interface{}{
		"desc": "get dsp reportId errors",
	})
	return "", errors.New("get dsp reportId errors ")
}

func CreateDspReportPeriod(region, reportType, startDate, endDate, profileId string, ctx context.Context) (string, error) {

	endpoint, ok := region2.ApiUrl[region]
	if !ok {
		return "", fmt.Errorf("region not found: %s", region)
	}

	url := endpoint + "/dsp/reports"

	refreshToken, ok := DspRefreshToken[region]
	if !ok {
		return "", fmt.Errorf("region token not found: %s", region)
	}

	headers, err := fakeHeaders(profileId, refreshToken )
	if err != nil {
		return "", err
	}

	params := map[string]interface{}{
		"format":     "CSV",
		"startDate":  startDate,
		"endDate":    endDate,
		"dimensions": dimensions[reportType],
		"metrics":    metrics[reportType],
	}

	// 根据报表类型设置参数
	switch reportType {
	case REPORT_TYPE_AUDIENCE:
		params["type"] = "AUDIENCE"
		params["timeUnit"] = "SUMMARY"
	case REPORT_TYPE_INVENTORY:
		params["type"] = "INVENTORY"
		params["timeUnit"] = "DAILY"
	case REPORT_TYPE_PRODUCT:
		params["type"] = "PRODUCTS"
		params["timeUnit"] = "DAILY"
	case REPORT_TYPE_DETAIL, REPORT_TYPE_ORDER:
		params["type"] = "CAMPAIGN"
		params["timeUnit"] = "DAILY"
	default:
		return "", fmt.Errorf("reportType not found:%s", reportType)
	}

	count := 0
	var resp *requests.Resp
	var rlimit = limiter.PerSecond(vars.LimitMap["dsp"])
	var sleepSec int
	for {

		for {

			res, err := boot.Rlimiter.Allow(url, rlimit)
			if err != nil {
				return "", err
			}

			if res.Allowed == 1 {
				break
			}

			time.Sleep(10 * time.Millisecond)
		}

		count += 1
		if count > 40 {
			logger.Logger.Info(map[string]interface{}{
				"desc":       " get reportId error ppc get sleep too long to break",
				"count":      count,
				"reportName": reportType,
				"reportType": "dsp",
			})
			break
		}

		resp, err = requests.Post(url, requests.WithHeaders(headers), requests.WithJson(params), requests.WithTimeout(time.Minute))
		if err != nil || resp == nil {
			sleepSec = rand.Intn(30) + 1 + count*20
			if resp != nil {
				errmap := map[string]interface{}{"code": resp.StatusCode, "body": string(resp.Body)}
				logger.Logger.ErrorWithContextMap(ctx, errmap, "dsp http post create report error")
			}
			if err != nil {
				errstr := err.Error()
				errmap := map[string]interface{}{"err": errstr, "reportType": reportType}
				//超时休眠重试
				if strings.Contains(errstr, vars.Timeouttxt) || strings.Contains(errstr, vars.Timeouttxt1) {
					logger.Logger.ErrorWithContextMap(ctx, errmap, "dsp http post create report error")
					time.Sleep(time.Duration(sleepSec))
					continue
				} else {
					logger.Logger.ErrorWithContextMap(ctx, errmap, "dsp http post create report error")
				}
			}

			// 当前如果网络不好或者抖动，则添加休眠，避免中间持续频繁请求导致恶化
			time.Sleep(time.Duration(sleepSec))
			return "", err
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			sleepSec = rand.Intn(30) + 1 + count*30
			logger.Info(map[string]interface{}{
				"flag":     "StatusTooManyRequests",
				"count":    count,
				"sleepSec": sleepSec,
				"path":     url,
			})

			time.Sleep(time.Second * time.Duration(sleepSec))
			continue
		}

		if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
			logger.Info(map[string]interface{}{
				"desc":      "dsp http post create report error",
				"params":    params,
				"http_code": resp.StatusCode,
				"body":      string(resp.Body),
			})

			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				return "", fmt.Errorf("resp error, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body))
			} else {
				sleepSec = rand.Intn(30) + 1 + count*20
				logger.Info(map[string]interface{}{
					"status":  resp.StatusCode,
					"sleep":   sleepSec,
					"num":     count,
					"moudule": "get report status",
				})
				time.Sleep(time.Second * time.Duration(sleepSec))
				continue
			}

		} else {

			val, err := resp.Json()
			if err != nil {
				return "", err
			}

			reportId, ok := val["reportId"].(string)
			if !ok {
				logger.Info(map[string]interface{}{
					"desc": "dsp create reportId not found",
				})

				return "", fmt.Errorf("reportId not found, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body))
			} else {
				return reportId, nil
			}
		}
	}

	return "", errors.New("get dsp reportId error")
}
