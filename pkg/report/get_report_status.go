package report

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"td_report/app/bean"
	"td_report/pkg/dbp_token"
	"td_report/pkg/logger"
	region2 "td_report/pkg/region"
	"td_report/pkg/requests"
	"td_report/vars"
	"time"
)

func fakeHeaders(profileId, refreshToken,clientId, clientSecret string) (map[string]string, error) {
	accessToken, err := dbp_token.GetAccessTokenWithClient(refreshToken,clientId, clientSecret)
	if err != nil {
		logger.Logger.Info(err, "get token error")
		return nil, err
	}

	headers := map[string]string{
		"Amazon-Advertising-API-Scope":    profileId,
		"Amazon-Advertising-API-ClientId": clientId,
		"Authorization":                   accessToken,
	}

	return headers, nil
}

func GetHeaderWithClient(profileId, refreshToken, clientId, clientSecret string) (map[string]string, error) {
	return fakeHeaders(profileId, refreshToken, clientId, clientSecret)
}

type ReportStatusBody struct {
	ReportId string `json:"reportId"`
	Status   string `json:"status"`
	Location string `json:"location"`
}

// ReportBrandMetrics  用于sb brand metrics
type ReportBrandMetrics struct {
	Expiration    int    `json:"expiration"`
	Format        string `json:"format"`
	Location      string `json:"location"`
	ReportId      string `json:"reportId"`
	Status        string `json:"status"`
	StatusDetails string `json:"statusDetails"`
}

func GetReportStatus(profileToken *bean.ProfileToken, reportId, reportType string, ctx context.Context) (*ReportStatusBody, error) {
	path := fmt.Sprintf("/v2/reports/%s", reportId)
	endpoint, ok := region2.ApiUrl[profileToken.Region]
	if !ok {
		return nil, fmt.Errorf("region not found: %s", profileToken.Region)
	}

	url := endpoint + path

	headers, err := fakeHeaders(profileToken.ProfileId, profileToken.RefreshToken, profileToken.ClientId, profileToken.ClientSecret)
	if err != nil {
		return nil, err
	}

	// 之前多次发现返回的Location 是空的不知道什么情况
	num := 0
	//var rlimit = limiter.PerSecond(vars.LimitMap[reportType])
	var sleepSec int
	for {

		//添加评率的限制，避免请求频次太多
		//for {
		//
		//	res, err := boot.Rlimiter.Allow(endpoint+"/v2/reports", rlimit)
		//	if err != nil {
		//		return nil, err
		//	}
		//
		//	if res.Allowed == 1 {
		//		break
		//	}
		//
		//	time.Sleep(10 * time.Millisecond)
		//}

		num++
		if num > 40 {
			logger.Logger.InfoWithContextMap(ctx, map[string]interface{}{
				"count":      num,
				"reportType": reportType,
			}, "get reportId error ppc get sleep too long to break")
			break
		}
		logger.Logger.InfoWithContext(ctx, headers, " get report status header====")
		resp, err := requests.Get(url, requests.WithHeaders(headers), requests.WithTimeout(time.Second*50))
		if err != nil || resp == nil {
			sleepSec = rand.Intn(20) + 1 + num*20
			if resp != nil {
				errmap := map[string]interface{}{"code": resp.StatusCode, "body": string(resp.Body)}
				logger.Logger.ErrorWithContextMap(ctx, errmap, "get report status  error")
			}
			if err != nil {
				errstr := err.Error()
				errmap := map[string]interface{}{"err": errstr, "reportType": reportType}
				//超时休眠重试
				if strings.Contains(errstr, vars.Timeouttxt) || strings.Contains(errstr, vars.Timeouttxt1) {
					logger.Logger.ErrorWithContextMap(ctx, errmap, "get report status  error,timeout retry")
					time.Sleep(time.Duration(sleepSec))
					continue
				} else {
					logger.Logger.ErrorWithContextMap(ctx, errmap, " get report status  error")
				}
			}

			// 当前如果网络不好或者抖动，则添加休眠，避免中间持续频繁请求导致恶化
			time.Sleep(time.Duration(sleepSec))
			return nil, err
		}

		logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s报表获取下载地址结果，statusCode:%d, body:%s", reportType, resp.StatusCode, string(resp.Body)))

		if resp.StatusCode == http.StatusTooManyRequests {
			sleepSec = rand.Intn(30) + 1 + num*20
			logger.Logger.InfoWithContext(ctx, map[string]interface{}{
				"flag":      "StatusTooManyRequests",
				"count":     num,
				"sleepSec":  sleepSec,
				"profileId": profileToken.ProfileId,
				"reportId":  reportId,
			})

			time.Sleep(time.Second * time.Duration(sleepSec))
			continue
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				logger.Logger.ErrorWithContext(ctx, fmt.Sprintf("get status resp error, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body)))
				return nil, fmt.Errorf("get status resp error, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body))
			} else {
				sleepSec = rand.Intn(30) + 1 + num*20
				logger.Logger.InfoWithContext(ctx, map[string]interface{}{
					"status":  resp.StatusCode,
					"num":     num,
					"sleep":   sleepSec,
					"moudule": "get report status",
					"desc":    string(resp.Body),
				})
				time.Sleep(time.Second * time.Duration(sleepSec))
				continue
			}
		}

		var reportStatusBody ReportStatusBody
		_, err = resp.Json(&reportStatusBody)
		if err != nil {
			logger.Logger.InfoWithContext(ctx, reportType+"json error")
			return nil, err
		}

		if reportStatusBody.Location != "" {
			return &reportStatusBody, nil
		} else {
			sleepSec = num*20 + rand.Intn(20)
			logger.Logger.InfoWithContext(ctx, map[string]interface{}{
				"reportType": reportType,
				"desc":       " get location is empty and sleep time",
				"profileId":  profileToken.ProfileId,
				"count":      num,
				"sleeptime":  sleepSec,
			})
			time.Sleep(time.Duration(sleepSec) * time.Second)
		}
	}
	logger.Logger.InfoWithContext(ctx, "cannot get data")
	return nil, errors.New("cannot get data")
}

func GetReportSbbrandMetricsGet(profileToken *bean.ProfileToken, reportId string, ctx context.Context) (*ReportStatusBody, error) {
	path := fmt.Sprintf("/insights/brandMetrics/report/%s", reportId)
	endpoint, ok := region2.ApiUrl[profileToken.Region]
	if !ok {
		return nil, fmt.Errorf("region not found: %s", profileToken.Region)
	}

	url := endpoint + path
	headers, err := fakeHeaders(profileToken.ProfileId, profileToken.RefreshToken, profileToken.ClientId, profileToken.ClientSecret)
	if err != nil {
		return nil, err
	}

	// 之前多次发现返回的Location 是空的不知道什么情况
	num := 0
	//var rlimit = limiter.PerSecond(vars.LimitMap["sb"] + 1)
	var sleepSec int
	for {

		//添加评率的限制，避免请求频次太多
		//for {
		//
		//	res, err := boot.Rlimiter.Allow("/insights/brandMetrics/report", rlimit)
		//	if err != nil {
		//		return nil, err
		//	}
		//
		//	if res.Allowed == 1 {
		//		break
		//	}
		//
		//	time.Sleep(10 * time.Millisecond)
		//}

		num++
		if num > 40 {
			logger.Logger.InfoWithContext(ctx, map[string]interface{}{
				"desc":       "get reportId error ppc get sleep too long to break",
				"count":      num,
				"reportType": "sb",
			})
			break
		}

		resp, err := requests.Get(url, requests.WithHeaders(headers), requests.WithTimeout(time.Second*50))
		if err != nil || resp == nil {
			sleepSec = rand.Intn(30) + 1 + num*20
			if resp != nil {
				errmap := map[string]interface{}{"code": resp.StatusCode, "body": string(resp.Body)}
				logger.Logger.ErrorWithContextMap(ctx, errmap, "get report status  error")
			}
			if err != nil {
				errstr := err.Error()
				errmap := map[string]interface{}{"err": errstr, "reportType": vars.SB}
				//超时休眠重试
				if strings.Contains(errstr, vars.Timeouttxt) || strings.Contains(errstr, vars.Timeouttxt1) {
					logger.Logger.ErrorWithContextMap(ctx, errmap, "get report status  error,timeout retry")
					time.Sleep(time.Duration(sleepSec))
					continue
				} else {
					logger.Logger.ErrorWithContextMap(ctx, errmap, " get report status  error")
				}
			}

			// 当前如果网络不好或者抖动，则添加休眠，避免中间持续频繁请求导致恶化
			time.Sleep(time.Duration(sleepSec))
			return nil, err
		}

		logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s报表获取下载地址结果，statusCode:%d, body:%s", vars.SB, resp.StatusCode, string(resp.Body)))
		if resp.StatusCode == http.StatusTooManyRequests {
			sleepSec = rand.Intn(20) + num*30
			logger.Logger.InfoWithContext(ctx, map[string]interface{}{
				"flag":      "StatusTooManyRequests",
				"count":     num,
				"sleepSec":  sleepSec,
				"profileid": profileToken.ProfileId,
			})

			time.Sleep(time.Second * time.Duration(sleepSec))
			continue
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				logger.Logger.ErrorWithContext(ctx, fmt.Sprintf("resp error, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body)))
				return nil, fmt.Errorf("resp error, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body))
			} else {
				sleepSec = rand.Intn(30) + 1 + num*20
				logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
					"status":  resp.StatusCode,
					"sleep":   sleepSec,
					"moudule": "get report status",
					"desc":    string(resp.Body),
				})

				time.Sleep(time.Second * time.Duration(sleepSec))
				continue
			}
		}

		var reportStatusBody ReportBrandMetrics
		_, err = resp.Json(&reportStatusBody)
		if err != nil {
			logger.Logger.InfoWithContext(ctx, "json error ")
			return nil, err
		}

		if reportStatusBody.Location != "" {
			var resResult ReportStatusBody
			resResult.ReportId = reportStatusBody.ReportId
			resResult.Status = reportStatusBody.Status
			resResult.Location = reportStatusBody.Location
			return &resResult, nil
		} else {

			sleepSec = num*20 + rand.Intn(40)
			logger.Logger.InfoWithContext(ctx, "body====", reportStatusBody)
			if reportStatusBody.Status == "SUCCESSFUL" || reportStatusBody.Status == "SUCCESS" {
				// 拉取不到数据，数据是空的，没有数据则中断没有必要重试
				break
			}
			logger.Logger.InfoWithContext(ctx, map[string]interface{}{
				"desc":     "sb brand get location is empty and sleep time ",
				"num":      num,
				"sleepsec": sleepSec,
			})
			time.Sleep(time.Second * time.Duration(sleepSec))
		}
	}

	logger.Logger.InfoWithContext(ctx, "cannot get data")
	return nil, errors.New("cannot get data")
}
