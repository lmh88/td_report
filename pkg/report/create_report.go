package report

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"td_report/app/bean"
	"td_report/boot"
	"td_report/pkg/limiter"
	"td_report/pkg/logger"
	region2 "td_report/pkg/region"
	"td_report/pkg/requests"
	"td_report/vars"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// SBCreateBrandMetricsReportPost  独立的报表
func SBCreateBrandMetricsReportPost(token *bean.ProfileToken, dateStr, enddate, lookBackPeriod string, ctx context.Context) (string, error) {
	path := "/insights/brandMetrics/report"
	return createReportSbPostMetrics(token, dateStr, enddate, path, lookBackPeriod, ctx)
}

func SBCreateReport(profileToken *bean.ProfileToken, reportType, dateStr, segment, creativeType, metrics string, ctx context.Context) (string, error) {
	path := fmt.Sprintf("/v2/hsa/%s/report", reportType)
	return createReport(profileToken, reportType, dateStr, segment, creativeType, metrics, path, vars.SB, ctx)
}

func SPCreateReport(profileToken *bean.ProfileToken, reportType, dateStr, segment, metrics string, ctx context.Context) (string, error) {
	var path string
	if reportType == "campaignsPlacement" {
		path = fmt.Sprintf("/v2/sp/%s/report", "campaigns")
	} else {
		path = fmt.Sprintf("/v2/sp/%s/report", reportType)
	}

	return createReport(profileToken, reportType, dateStr, segment, "", metrics, path, vars.SP, ctx)
}

func SPCreateReportV1(profileToken *bean.ProfileToken, reportType, dateStr, segment, metrics string, ctx context.Context) (string, error) {
	var path string
	if reportType == "campaignsPlacement" {
		path = fmt.Sprintf("/v2/sp/%s/report", "campaigns")
	} else {
		path = fmt.Sprintf("/v2/sp/%s/report", reportType)
	}

	return createReport(profileToken, reportType, dateStr, segment, "", metrics, path, vars.SP, ctx)
}

func SDCreateReport(profileToken *bean.ProfileToken, reportType, dateStr, tactic, metrics string, ctx context.Context) (string, error) {
	path := fmt.Sprintf("/sd/%s/report", reportType)
	return createReport(profileToken, reportType, dateStr, tactic, "", metrics, path, vars.SD, ctx)
}

func createReportSbPostMetrics(token *bean.ProfileToken, dateStr, enddate, path string, lookBackPeriod string, ctx context.Context) (string, error) {
	endpoint, ok := region2.ApiUrl[token.Region]
	if !ok {
		return "", fmt.Errorf("region not found: %s", token.Region)
	}

	url := endpoint + path
	headers, err := fakeHeaders(token.ProfileId, token.RefreshToken, token.ClientId, token.ClientSecret)
	if err != nil {
		return "", err
	}

	params := map[string]interface{}{
		"format": "JSON",
	}
	if dateStr != "" {
		params["reportStartDate"] = dateStr
	}
	if enddate != "" {
		params["reportEndDate"] = enddate
	}
	if lookBackPeriod != "" {
		params["lookBackPeriod"] = lookBackPeriod
	}

	count := 0
	var (
		resp      *requests.Resp
		rlimit    = limiter.PerSecond(vars.LimitMap["sb"])
		reportId  string
		sleepSec  int
		profileId = token.ProfileId
	)

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
			logger.Logger.InfoWithContext(ctx, "sleep too loog to break")
			break
		}

		resp, err = requests.Post(url, requests.WithHeaders(headers), requests.WithJson(params), requests.WithTimeout(time.Second*60))
		if err != nil || resp == nil {
			sleepSec = rand.Intn(30) + 1 + count*20
			if resp != nil {
				errmap := map[string]interface{}{"code": resp.StatusCode, "body": string(resp.Body)}
				logger.Logger.ErrorWithContextMap(ctx, errmap, "create report get post data error")
			}
			if err != nil {
				errstr := err.Error()
				//超时休眠重试
				if strings.Contains(errstr, vars.Timeouttxt) || strings.Contains(errstr, vars.Timeouttxt1) {
					errmap := map[string]interface{}{"err": errstr}
					logger.Logger.ErrorWithContextMap(ctx, errmap, "create report get post data error,timeout retry")
					time.Sleep(time.Duration(sleepSec))
					continue

				} else {
					errmap := map[string]interface{}{"err": errstr}
					logger.Logger.ErrorWithContextMap(ctx, errmap, "create report get post data error")
				}
			}

			// 当前如果网络不好或者抖动，则添加休眠，避免中间持续频繁请求导致恶化
			time.Sleep(time.Duration(sleepSec))
			return "", err
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			sleepSec = rand.Intn(30) + 1 + count*20
			logger.Info(map[string]interface{}{
				"flag":      "StatusTooManyRequests",
				"count":     count,
				"sleepSec":  sleepSec,
				"profileid": profileId,
			})

			time.Sleep(time.Second * time.Duration(sleepSec))
			continue
		}

		// 除了429之外的错误信息直接返回错误信息
		if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				logger.Logger.InfoWithContextMap(ctx, map[string]interface{}{"code": resp.StatusCode, "body": string(resp.Body)}, "resp error")
				return "", fmt.Errorf("resp error, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body))
			} else {
				sleepSec = rand.Intn(30) + 1 + count*20
				logger.Logger.InfoWithContext(ctx, map[string]interface{}{
					"status":  resp.StatusCode,
					"sleep":   sleepSec,
					"num":     count,
					"moudule": "get report status",
				})
				time.Sleep(time.Second * time.Duration(sleepSec))
				continue
			}

		}

		val, err := resp.JsonAndValueIsAny()
		if err != nil {
			logger.Logger.InfoWithContext(ctx, map[string]interface{}{
				"flag":      "json marson error",
				"count":     count,
				"err":       err,
				"profileid": profileId,
				"body":      string(resp.Body),
			})

			return "", err
		}

		reportId, ok = val["reportId"].(string)
		if !ok {
			time.Sleep(time.Second * time.Duration(rand.Intn(30)+count*20))
			continue
		} else {
			return reportId, nil
		}
	}

	return "", errors.New("reportid is empty ,get reportId error")
}

func createReport(profileToken *bean.ProfileToken, reportType, dateStr, segment, creativeType, metrics, path, report string, ctx context.Context) (string, error) {
	//func createReport(region, reportType, dateStr, segment, profileId, creativeType, metrics, refreshToken, path, report string, ctx context.Context) (string, error) {
	region := profileToken.Region
	refreshToken := profileToken.RefreshToken
	profileId := profileToken.ProfileId
	endpoint, ok := region2.ApiUrl[region]
	if !ok {
		return "", fmt.Errorf("region not found: %s", region)
	}

	// 添加一道过滤，避免短时间内重复请求token
	if exists, _ := vars.Cache.Contains(refreshToken); exists {
		logger.Logger.InfoWithContext(ctx, "error token")
		return "", errors.New("error token")
	}

	url := endpoint + path
	headers, err := fakeHeaders(profileId, refreshToken, profileToken.ClientId, profileToken.ClientSecret)
	if err != nil {
		//针对token失效的情况，缓存下来，直接跳过
		vars.Cache.SetIfNotExist(refreshToken, refreshToken, time.Second*600)
		logger.Logger.InfoWithContext(ctx, err, reportType, report, "get header error")
		return "", err
	}

	params := map[string]interface{}{
		"reportDate": dateStr,
		"metrics":    metrics,
	}

	if report != "sd" {
		//campaigns
		//if report == vars.SP && reportType == "campaigns" {
		//} else {
		//	params["segment"] = segment
		//}

	} else {
		params["tactic"] = segment
	}

	if report == "sb" {
		params["creativeType"] = creativeType
	}

	if report == "sp" && reportType == "asins" {
		params["campaignType"] = "sponsoredProducts"
	}

	count := 0
	var resp *requests.Resp
	//var rlimit = limiter.PerSecond(vars.LimitMap[report])
	var reportId string
	var sleepSec int
	for {

		//添加频率率的限制，避免请求频次太多
		//for {
		//
		//	res, err := boot.Rlimiter.Allow(url, rlimit)
		//	if err != nil {
		//		return "", err
		//	}
		//
		//	if res.Allowed == 1 {
		//		break
		//	}
		//
		//	time.Sleep(10 * time.Millisecond)
		//}

		count += 1
		if count > 40 {
			logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
				"desc":       " get reportId error ppc get sleep too long to break",
				"count":      count,
				"reportName": reportType,
				"reportType": report,
			})
			break
		}

		resp, err = requests.Post(url, requests.WithHeaders(headers), requests.WithJson(params), requests.WithTimeout(time.Second*40))
		if err != nil || resp == nil {
			sleepSec = rand.Intn(20) + 1 + count*20
			if resp != nil {
				errmap := map[string]interface{}{"code": resp.StatusCode, "body": string(resp.Body)}
				logger.Logger.ErrorWithContextMap(ctx, errmap, "create report get post data error")
			}
			if err != nil {
				errstr := err.Error()
				//超时休眠重试
				if strings.Contains(errstr, vars.Timeouttxt) || strings.Contains(errstr, vars.Timeouttxt1) {
					errmap := map[string]interface{}{"err": errstr}
					logger.Logger.ErrorWithContextMap(ctx, errmap, "create report get post data error,timeout retry")
					time.Sleep(time.Duration(sleepSec))
					continue

				} else {
					errmap := map[string]interface{}{"err": errstr}
					logger.Logger.ErrorWithContextMap(ctx, errmap, "create report get post data error")
				}
			}

			// 当前如果网络不好或者抖动，则添加休眠，避免中间持续频繁请求导致恶化
			time.Sleep(time.Duration(sleepSec))
			return "", err
		}

		logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s报表第一次请求结果，statusCode:%d, body:%s", report, resp.StatusCode, string(resp.Body)))

		if resp.StatusCode == http.StatusTooManyRequests {
			sleepSec = rand.Intn(30) + count*20
			logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
				"flag":      "StatusTooManyRequests",
				"count":     count,
				"sleepSec":  sleepSec,
				"profileid": profileId,
			})

			time.Sleep(time.Second * time.Duration(sleepSec))
			continue
		}

		// 除了429之外的错误信息直接返回错误信息
		if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				da, _ := json.Marshal(params)
				logger.Logger.ErrorWithContext(ctx, fmt.Sprintf("resp error, StatusCode=%d, body=%s, parama=%s", resp.StatusCode, string(resp.Body), string(da)))
				return "", fmt.Errorf("resp error, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body))
			} else {
				sleepSec = rand.Intn(30) + 1 + count*20
				logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
					"status": resp.StatusCode,
					"sleep":  sleepSec,
					"count":  count,
					"module": "get report status",
					"desc":   string(resp.Body),
				})

				time.Sleep(time.Second * time.Duration(sleepSec))
				continue
			}
		}

		break
	}

	val, err := resp.JsonAndValueIsAny()
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
			"desc": "json error",
		})
		return "", err
	}

	reportId, ok = val["reportId"].(string)
	if !ok {
		time.Sleep(time.Second * time.Duration(rand.Intn(30)+count*20))
		logger.Logger.InfoWithContext(ctx, "reportid is empty retry again ")
		return "", errors.New("reportid is empty ,get reportId error ")
	} else {
		return reportId, nil
	}

}
