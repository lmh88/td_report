package other

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/spf13/cobra"
	"net/http"
	"strings"
	"td_report/pkg/dbp_token"
	"td_report/pkg/dsp/dsp_report"
	"td_report/pkg/logger"
	region2 "td_report/pkg/region"
	"td_report/pkg/requests"
	"td_report/pkg/save_file"
	"td_report/vars"
	"time"
)

// 临时调研
var dspTempCmd = &cobra.Command{
	Use:   "dsp_temp",
	Short: "dsp_temp",
	Long:  `dsp_temp`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dspTemp called")
		dspTempCmdFunc()
	},
}

func init() {
	RootCmd.AddCommand(dspTempCmd)
}

const (
	ReportTypeOrder     = "order"
	ReportTypeDetail    = "detail"
	ReportTypeAudience  = "audience"
	ReportTypeInventory = "inventory"
	ReportTypeProduct   = "products"
	ReportTypeGeography = "geography"
)

var dimensions = map[string][]string{
	ReportTypeOrder:     {"ORDER"},
	ReportTypeDetail:    {"ORDER", "LINE_ITEM", "CREATIVE"},
	ReportTypeAudience:  {"ORDER", "LINE_ITEM"},
	ReportTypeInventory: {"ORDER", "LINE_ITEM", "SITE", "SUPPLY"},
	ReportTypeProduct:   {"ORDER", "LINE_ITEM"},
	ReportTypeGeography: {"ORDER", "LINE_ITEM", "COUNTRY", "STATE_COUNTY_REGION", "CITY"},
}

var metrics = map[string]string{
	ReportTypeOrder:     "totalSales14d,clickThroughs,dpv14d,impressions,totalAddToCart14d,totalNewToBrandProductSales14d,sales14d,totalDetailPageViews14d,totalCost,totalPurchases14d,totalUnitsSold14d,supplyCost,amazonPlatformFee,amazonAudienceFee,agencyFee,3PFees,totalFee,viewableImpressions,measurableImpressions,totalPixel14d,totalPixelViews14d,totalPixelClicks14d,dpvViews14d,dpvClicks14d,atc14d,atcViews14d,atcClicks14d,purchases14d,purchasesViews14d,purchasesClicks14d,newToBrandPurchases14d,newToBrandPurchasesViews14d,newToBrandPurchasesClicks14d,unitsSold14d,newToBrandUnitsSold14d,newToBrandProductSales14d,totalDetailPageViewViews14d,totalDetailPageClicks14d,totalAddToCartViews14d,totalAddToCartClicks14d,totalPurchasesViews14d,totalPurchasesClicks14d,totalNewToBrandPurchases14d,totalNewToBrandPurchasesViews14d,totalNewToBrandPurchasesClicks14d,totalNewToBrandUnitsSold14d,totalPRPV14d,totalPRPVViews14d,totalPRPVClicks14d,totalAddToList14d,totalAddToListViews14d,totalAddToListClicks14d,totalSubscribeAndSaveSubscriptions14d,totalSubscribeAndSaveSubscriptionViews14d,totalSubscribeAndSaveSubscriptionClicks14d,newSubscribeAndSave14d,newSubscribeAndSaveViews14d,newSubscribeAndSaveClicks14d,pRPV14d,pRPVViews14d,pRPVClicks14d,atl14d,atlViews14d,atlClicks14d",
	ReportTypeDetail:    "impressions,clickThroughs,dpv14d,totalDetailPageViews14d,totalAddToCart14d,totalPurchases14d,totalSales14d,totalUnitsSold14d,totalNewToBrandProductSales14d,totalCost,sales14d",
	ReportTypeInventory: "impressions,clickThroughs,dpv14d,totalDetailPageViews14d,totalAddToCart14d,totalPurchases14d,totalSales14d,totalUnitsSold14d,totalNewToBrandProductSales14d,totalCost,sales14d,supplyCost,placementName,placementSize",
	ReportTypeAudience:  "impressions,clickThroughs,dpv14d,atc14d,purchases14d,totalCost,segmentSource,segmentType,segmentMarketplaceID,lineitemtype,segmentClassCode,targetingMethod,CTR,eCPC,eCPM",
	ReportTypeProduct:   "productName,productGroup,productCategory,productSubcategory,dpv14d,dpvViews14d,dpvClicks14d,pRPV14d,pRPVViews14d,pRPVClicks14d,atl14d,atlViews14d,atlClicks14d,atc14d,atcViews14d,atcClicks14d,purchases14d,purchasesViews14d,purchasesClicks14d,newToBrandPurchases14d,newToBrandPurchasesViews14d,newToBrandPurchasesClicks14d,percentOfPurchasesNewToBrand14d,totalDetailPageViews14d,totalDetailPageViewViews14d,totalDetailPageClicks14d,totalPRPV14d,totalPRPVViews14d,totalPRPVClicks14d,totalAddToList14d,totalAddToListViews14d,totalAddToListClicks14d,totalAddToCart14d,totalAddToCartViews14d,totalAddToCartClicks14d,totalPurchases14d,totalPurchasesViews14d,totalPurchasesClicks14d,totalNewToBrandPurchases14d,totalNewToBrandPurchasesViews14d,totalNewToBrandPurchasesClicks14d,totalPercentOfPurchasesNewToBrand14d,newSubscribeAndSave14d,newSubscribeAndSaveViews14d,newSubscribeAndSaveClicks14d,totalSubscribeAndSaveSubscriptions14d,totalSubscribeAndSaveSubscriptionViews14d,totalSubscribeAndSaveSubscriptionClicks14d,unitsSold14d,sales14d,totalUnitsSold14d,totalSales14d,newToBrandUnitsSold14d,newToBrandProductSales14d,brandHaloDetailPage14d,brandHaloDetailPageViews14d,brandHaloDetailPageClicks14d,brandHaloProductReviewPage14d,brandHaloProductReviewPageViews14d,brandHaloProductReviewPageClicks14d,brandHaloAddToList14d,brandHaloAddToListViews14d,brandHaloAddToListClicks14d,brandHaloAddToCart14d,brandHaloAddToCartViews14d,brandHaloAddToCartClicks14d,brandHaloPurchases14d,brandHaloPurchasesViews14d,brandHaloPurchasesClicks14d,brandHaloNewToBrandPurchases14d,brandHaloNewToBrandPurchasesViews14d,brandHaloNewToBrandPurchasesClicks14d,brandHaloPercentOfPurchasesNewToBrand14d,brandHaloNewSubscribeAndSave14d,brandHaloNewSubscribeAndSaveViews14d,brandHaloNewSubscribeAndSaveClicks14d,brandHaloTotalUnitsSold14d,brandHaloTotalSales14d,brandHaloTotalNewToBrandSales14d,brandHaloTotalNewToBrandUnitsSold14d",
	ReportTypeGeography: "totalCost,impressions,viewableImpressions,clickThroughs",
}

var DspRefreshToken = map[string]string{
	region2.RegionFe: "Atzr|IwEBIJqR9kHQts1y3pFLrTeLMA8afkEd_dNUAwsjSjqIhFEZHCmVGyhiOkSqZ1gNHUY2X1TSLJjhUeRm6PVe1FnzBAo7o6bOmE495OjbGKaOQAJ2fHh1rBS51dJLT4oX1l7-bWJYb2npLXmOjZQgY8t1ZCDR-45KY0SqzjrbtQMDjNh8uAKRiwhBwBu3H13FcJxxOBi8yjrS2zKuNb2wmWBTe9vUhPbigbvvAMzMeTFmMSK-_N5hd8ZglvJjkvobK5SD7BQeDXK49t5BN-05WSzR6kLfFa9gWSSYphoPdF9hUYaUHIDBUj4ZOeWURoFugZzDO408Xc0lee6GGz1CP21aumiWlSCTuEUrOAENGcNaipEEvWRvZDiHKhIX546ycrVt95ujrQgIlVm5H2-copN3Pkb1KAUBSmW7Ylgwha6OB5fJNg",
	region2.RegionEu: "Atzr|IwEBIJ9_AwBc6Pzv65zeD4ph1nrxNqceokx8IrqxRktMCrpeuC0ZdlqGqUjD8tO3TUCoD8804QB4ND-1Ng-9bY5gMbay9O69vgYN4c2A36y-ep_jzRF4kjyuU22ALAnZjjj6FQ0SoRTkstAZUd7zB3z7U2_i6SIRyaXpDhhIPSAleYEBlQAyM48Lph_VW5dU8pNL7OQrh9r1SW05yyuB7lYaA1juQDYdX4U8WRez0iLssUXhTbCSqTUSSFuqBedAzJMcge84iI_EYp6C9wtFFDFJasXmEfOeLA4S2x_OS07PW06nvXjL3EeLDxgCZqGRTHSVj3V0BbsypQdjYgD7chXeZpAs-YqJKE1dZQ4lPt6yomeJzFFmtO_3RNO1H7_II-nKW3OSHC3WjCG5V9on7_vR1sEI4PZmD0d1RUEEfOL2liTp5ThjRq809SnV9YZL38g9As8",
	region2.RegionNa: "Atzr|IwEBIFohGsLBoxln39hOsfSZfYN18O1YJINunLjrMwBv28y6rpH7EtK8PC19orgSR6TzQE5IkWVOB-e-QMIDCYvDcTCxg9Zqndo9pHLovbPLmXYPdt04DRyEZKgptREAhq8GmZvkZww8ldYm4nGnpFOpV7Czoy8CAjmZHQF9B-jzQhOzjxTKhCUeM2M-6kojkHLmJQqDs3n8QiZr4RVscL3Y41ZLvYwh8D5H1Z93gQNuIdsFzBq0bookPCWJuJeQjC484654ox5AMLfaKlQEejUZQNDdB-Dn6zuvEGudeDPbFh69cANaA7EVHUd4XSl6SP3l0QrW87qLCWz5iXnEnDATVa4-MyGCiMVE-zt4zxxkVJXyDo8UZ8zaFS3qAUdlHOJZ_nGHFHJN7XreJ-cw3foIRLFR1Q6TBYbt6nnmdhFwgj4K2OcePo-Z2T_tvDSIAuLdYn4",
}

func getreportId(profileId, myregion, dateStr, accountId, reportType string) (string, error) {
	url := region2.ApiUrl[myregion] + fmt.Sprintf("/accounts/%s/dsp/reports", accountId)

	params := map[string]interface{}{
		"format":     "CSV",
		"startDate":  dateStr,
		"endDate":    dateStr,
		"dimensions": dimensions[reportType],
		"metrics":    strings.Split(metrics[reportType], ","),
	}

	// 根据报表类型设置参数
	switch reportType {
	case ReportTypeAudience:
		params["type"] = "AUDIENCE"
		params["timeUnit"] = "SUMMARY"
	case ReportTypeInventory:
		params["type"] = "INVENTORY"
		params["timeUnit"] = "DAILY"
	case ReportTypeProduct:
		params["type"] = "PRODUCTS"
		params["timeUnit"] = "DAILY"
	case ReportTypeDetail, ReportTypeOrder:
		params["type"] = "CAMPAIGN"
		params["timeUnit"] = "DAILY"
	case ReportTypeGeography:
		params["type"] = "GEOGRAPHY"
		params["timeUnit"] = "DAILY"
	default:
		g.Log().Println("error")
		return "", errors.New("report type error ")
	}

	count := 0
	var resp *requests.Resp
	refreshToken, ok := DspRefreshToken[myregion]
	if !ok {
		g.Log().Println("regin not found")
		return "", fmt.Errorf("region token not found: %s", myregion)
	}

	headers, err := fakeHeaders(profileId, refreshToken)
	if err != nil {
		g.Log().Println("header  not found")
		return "", err
	}

	headers["Accept"] = "application/vnd.dspcreatereports.v3+json"
	for {

		count += 1
		if count > 30 {
			break
		}

		resp, err = requests.Post(url, requests.WithHeaders(headers), requests.WithJson(params), requests.WithTimeout(time.Minute))
		if err != nil {
			g.Log().Println("get dsp http post error ")
			return "", err
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			sleepSec := count * 20

			g.Log().Info(map[string]interface{}{
				"flag":     "StatusTooManyRequests",
				"count":    count,
				"sleepSec": sleepSec,
				"path":     url,
			})

			time.Sleep(time.Second * time.Duration(sleepSec))
			continue
		}

		break
	}

	if resp.StatusCode != http.StatusAccepted {
		g.Log().Info("resp error, StatusCod")
		return "", fmt.Errorf("resp error, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body))
	}

	val, err := resp.Json()
	if err != nil {
		return "", err
	}

	reportId, ok := val["reportId"].(string)
	if !ok {
		g.Log().Info("reportId is empty")
		return "", fmt.Errorf("reportId not found, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body))
	}

	return reportId, nil
}

func dspTempCmdFunc() {
	var profileId = "2997968305172217"
	var myregion = "NA"
	var dateStr = "2022-02-25"
	var accountId = "586273131829674524"
	count := 1
	taskType := ReportTypeAudience
	reportId, err := getreportId(profileId, myregion, dateStr, accountId, taskType)
	if err != nil {
		g.Log().Info("get reportId empty")
		return
	}
	g.Log().Println("========reportid:", reportId)
	for {

		ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s_%s", vars.DSP, taskType, profileId, time.Now().Unix()))
		resp, err := dsp_report.GetDspReportStatus(myregion, profileId, reportId, ctx)
		if err != nil {
			g.Log().Error(map[string]interface{}{
				"flag":      "GetDspReportStatus error",
				"err":       err.Error(),
				"region":    myregion,
				"profileId": profileId,
				"taskType":  taskType,
				"reportId":  reportId,
			})

			return
		}

		if resp.Status == "FAILURE" {
			g.Log().Error(map[string]interface{}{
				"flag":      "dsp report status is FAILURE",
				"region":    myregion,
				"profileId": profileId,
				"taskType":  taskType,
				"reportId":  reportId,
			})
			return
		}

		if resp.Status == "SUCCESS" {
			val, err := dsp_report.DownloadDspReport(resp.Location, ctx)
			if err != nil {
				g.Log().Error(map[string]interface{}{
					"flag":      "DownloadDspReport error",
					"resp":      resp,
					"err":       err.Error(),
					"region":    myregion,
					"profileId": profileId,
					"taskType":  taskType,
					"reportId":  reportId,
				})

				return
			}

			//写入文件
			err = save_file.SaveDspFile(taskType, dateStr, profileId, val)
			if err != nil {
				g.Log().Error(map[string]interface{}{
					"flag":      "SaveDspFile error",
					"err":       err.Error(),
					"region":    myregion,
					"profileId": profileId,
					"taskType":  taskType,
					"reportId":  reportId,
				})
				return
			}

			break
		}

		sleepInt := count * 2
		if sleepInt > 30 {
			sleepInt = 30
		}
		time.Sleep(time.Second * time.Duration(sleepInt))
		count += 1

		//兜底
		if count > 400 {
			g.Log().Error(map[string]interface{}{
				"flag":      "count > 400",
				"region":    myregion,
				"profileId": profileId,
				"taskType":  taskType,
				"reportId":  reportId,
			})
			break
		}
	}
}

func fakeHeaders(profileId, refreshToken string) (map[string]string, error) {
	accessToken, err := dbp_token.GetAccessToken(refreshToken)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Amazon-Advertising-API-Scope":    profileId,
		"Amazon-Advertising-API-ClientId": "amzn1.application-oa2-client.084663234c2143c3a3bf91fe34bbdf1e",
		"Authorization":                   accessToken,
	}

	return headers, nil
}
