package report_v3

import (
	"fmt"
	"log"
	"td_report/app/bean"
	"td_report/app/service/report_v2/create_report"
	"td_report/app/service/report_v2/varible"
	"td_report/pkg/logger"
	region2 "td_report/pkg/region"
	"td_report/pkg/requests"
	"time"
)


//func getReport()

func CreateReport(profileToken *bean.ProfileToken) {
	endpoint := region2.ApiUrl[profileToken.Region]
	url := endpoint + "/reporting/reports"
	fmt.Println(url)

	header := getReqHeader(profileToken)
	params := map[string]interface{}{
		"name": "gary-sp-campaign",
		"startDate": "2022-08-01",
		"endDate": "2022-08-05",
		"configuration": map[string]interface{}{
			"adProduct":"SPONSORED_PRODUCTS",
			"groupBy":[]string{"campaign"},
			"columns":[]string{"impressions","clicks","cost","campaignId","date"},
			"reportTypeId":"spCampaigns",
			"timeUnit":"DAILY",
			"format":"GZIP_JSON",
		},
	}

	resp, err := requests.Post(url, requests.WithHeaders(header),requests.WithJson(params), requests.WithTimeout(30 * time.Second))

	if err != nil {
		log.Println(err)
		logger.Logger.Error(err)
	}

	log.Println(resp.StatusCode, string(resp.Body))
	logger.Logger.Info(resp.StatusCode, string(resp.Body), resp.Header())

}

func getReqHeader(profileToken *bean.ProfileToken) map[string]string {
	crt := &varible.ClientRefreshToken{
		ClientId: profileToken.ClientId,
		ClientSecret: profileToken.ClientSecret,
		ProfileId: profileToken.ProfileId,
		RefreshToken: profileToken.RefreshToken,
	}
	headers, err := create_report.FakeHeaders(crt)

	if err != nil {
		logger.Logger.Error("getReqHeader_error:" + err.Error())
		return nil
	}
	headers["Content-Type"] = "application/vnd.createasyncreportrequest.v3+json"

	return headers
}

func GetReportStatus(profileToken *bean.ProfileToken, reportId string) {

	endpoint := region2.ApiUrl[profileToken.Region]
	url := endpoint + "/reporting/reports/" + reportId
	header := getReqHeader(profileToken)

	resp, err := requests.Get(url, requests.WithHeaders(header), requests.WithTimeout(30 * time.Second))

	if err != nil {
		log.Println(err)
		logger.Logger.Error(err)
	}

	log.Println(resp.StatusCode, string(resp.Body))
	logger.Logger.Info(resp.StatusCode, string(resp.Body), resp.Header())
}

