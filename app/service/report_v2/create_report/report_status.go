package create_report

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"td_report/app/service/report_v2/get_token"
	"td_report/app/service/report_v2/varible"
	"td_report/pkg/logger"
	region2 "td_report/pkg/region"
	"td_report/pkg/requests"
	"td_report/vars"
	"time"
)

type ReportStatusBody struct {
	ReportId string `json:"reportId"`
	Status   string `json:"status"`
	Location string `json:"location"`
}

func GetReportStatus(reportIdMsg *varible.ReportIdMsg, ctx context.Context) (*ReportStatusBody, error) {
	path := fmt.Sprintf("/v2/reports/%s", reportIdMsg.ReportId)
	endpoint, ok := region2.ApiUrl[reportIdMsg.Region]
	if !ok {
		return nil, fmt.Errorf("region not found: %s", reportIdMsg.Region)
	}

	url := endpoint + path
	//todo test
	crt := varible.ClientRefreshToken{
		ProfileId: reportIdMsg.ProfileId,
		RefreshToken: reportIdMsg.RefreshToken,
		ClientId: reportIdMsg.ClientId,
		ClientSecret: reportIdMsg.ClientSecret,
	}

	headers, err := FakeHeaders(&crt)
	if err != nil {
		return nil, err
	}

	sleepNum1, sleepNum2 := 0, 0
	var sleepSec int
	for {
		//todo test
		resp, err := requests.Get(url, requests.WithHeaders(headers), requests.WithTimeout(time.Second*50))

		if err != nil  {
			sleepNum1++
			sleepSec = rand.Intn(3) + 1
			errstr := err.Error()
			errmap := map[string]interface{}{"err": errstr, "reportType": reportIdMsg.ReportType}
			//超时休眠重试
			if checkTimeout(errstr) {
				logger.Logger.ErrorWithContextMap(ctx, errmap, "get report status  error,timeout retry")
				//尝试3次
				if sleepNum1 < 3 {
					time.Sleep(time.Duration(sleepSec) * time.Second)
					continue
				}
			} else {
				logger.Logger.ErrorWithContextMap(ctx, errmap, " get report status  error")
			}
			return nil, err
		}

		//logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s报表获取下载地址结果，statusCode:%d, body:%s", reportIdMsg.ReportType, resp.StatusCode, string(resp.Body)))

		if resp.StatusCode == http.StatusTooManyRequests {
			logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
				"flag":      "StatusTooManyRequests",
				"profileId": reportIdMsg.ProfileId,
				"reportId":  reportIdMsg.ReportId,
			})
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				sleepNum2++
				//尝试两次
				if sleepNum2 < 2 {
					sleepSec = rand.Intn(3) + 1
					time.Sleep(time.Duration(sleepSec) * time.Second)
					continue
				}

				return nil, fmt.Errorf("get status resp error, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body))
			}
		}

		var reportStatusBody ReportStatusBody
		_, err = resp.Json(&reportStatusBody)
		if err != nil {
			logger.Logger.InfoWithContext(ctx, reportIdMsg.ReportType + " json error ", resp)
			return nil, err
		}
		return &reportStatusBody, nil
	}
}

func checkTimeout(errStr string) bool {
	timeoutTxt := []string{vars.Timeouttxt, vars.Timeouttxt1, vars.Timeouttxt2}
	for _, txt := range timeoutTxt {
		if strings.Contains(errStr, txt) {
			return true
		}
	}
	return false
}

func FakeHeaders(crt *varible.ClientRefreshToken) (map[string]string, error) {

	accessToken, err := get_token.GetAccessToken(crt)

	if err != nil {
		logger.Logger.Info(err, "get token error")
		return nil, err
	}

	headers := map[string]string{
		"Amazon-Advertising-API-Scope":    crt.ProfileId,
		"Amazon-Advertising-API-ClientId": crt.ClientId,
		"Authorization":                   accessToken,
	}

	return headers, nil
}


func DownloadReport(reportIdMsg *varible.ReportIdMsg, location string, ctx context.Context) ([]byte, error) {

	crt := varible.ClientRefreshToken{
		ProfileId: reportIdMsg.ProfileId,
		RefreshToken: reportIdMsg.RefreshToken,
		ClientId: reportIdMsg.ClientId,
		ClientSecret: reportIdMsg.ClientSecret,
	}
	headers, err := FakeHeaders(&crt)
	if err != nil {
		return nil, err
	}

	num := 3
	for i := 0; i < num; i++ {
		resp, err := requests.Get(location, requests.WithHeaders(headers), requests.WithTimeout(time.Hour))
		if err != nil || resp == nil {
			errstr := err.Error()
			if strings.Contains(errstr, vars.Timeouttxt) || strings.Contains(errstr, vars.Timeouttxt1) {
				logger.Logger.ErrorWithContext(ctx, "network error")
				time.Sleep(time.Duration((i+1)*20) * time.Second)
				continue
			}

			logger.Logger.ErrorWithContext(ctx, "get download error:", err.Error())
			return nil, err
		}

		logger.Logger.InfoWithContext(ctx, fmt.Sprintf("报表下载数据结果，statusCode:%d", resp.StatusCode))

		if resp.StatusCode != http.StatusOK {
			logger.Logger.ErrorWithContext(ctx, fmt.Sprintf("download error status,code:%s, body:%s", resp.StatusCode, string(resp.Body)))
			if resp.StatusCode >= 500 {
				time.Sleep(time.Duration((i+1)*20) * time.Second)
				continue
			}
			return nil, fmt.Errorf("download report resp error, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body))
		}

		return resp.Body, nil
	}

	logger.Logger.ErrorWithContext(ctx, "get download error try 5 time:")
	return nil, errors.New("download file error")
}

//func getUrl(url string) (*requests.Resp, error) {
//	res := requests.Resp{
//		StatusCode: 500,
//		Body: []byte("{\"reportId\":\"haha\"}"),
//	}
//	//return &res, errors.New("chang_e")
//	return &res, nil
//}
