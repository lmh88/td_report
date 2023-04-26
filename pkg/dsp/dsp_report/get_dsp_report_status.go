package dsp_report

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"td_report/boot"
	"td_report/pkg/limiter"
	"td_report/pkg/logger"
	region2 "td_report/pkg/region"
	"td_report/pkg/report"
	"td_report/pkg/requests"
	"td_report/vars"
	"time"
)

func GetDspReportStatus(region, profileId, reportId string, ctx context.Context) (*report.ReportStatusBody, error) {
	path := fmt.Sprintf("/dsp/reports/%s", reportId)
	endpoint, ok := region2.ApiUrl[region]
	if !ok {
		return nil, fmt.Errorf("region not found: %s", region)
	}

	url := endpoint + path

	refreshToken, ok := DspRefreshToken[region]
	if !ok {
		return nil, fmt.Errorf("region token not found: %s", region)
	}

	headers, err := fakeHeaders(profileId, refreshToken)
	if err != nil {
		return nil, err
	}

	count := 0
	var resp *requests.Resp
	var rlimit = limiter.PerSecond(vars.LimitMap["dsp"])
	var sleepSec int
	for {

		for {

			res, err := boot.Rlimiter.Allow(url, rlimit)
			if err != nil {
				return nil, err
			}

			if res.Allowed == 1 {
				break
			}

			time.Sleep(10 * time.Millisecond)
		}

		count += 1
		if count > 40 {
			logger.Logger.InfoWithContext(ctx, map[string]interface{}{
				"num":       count,
				"desc":      "dsp sleep too long to break",
				"profileid": profileId,
				"reportId":  reportId,
			})
			break
		}

		resp, err = requests.Get(url, requests.WithHeaders(headers), requests.WithTimeout(time.Minute*1))
		if err != nil || resp == nil {
			sleepSec = rand.Intn(30) + 1 + count*20
			if resp != nil {
				errmap := map[string]interface{}{"code": resp.StatusCode, "body": string(resp.Body)}
				logger.Logger.ErrorWithContextMap(ctx, errmap, "dsp http post create report error")
			}
			if err != nil {
				errstr := err.Error()
				errmap := map[string]interface{}{"err": errstr, "reportType": vars.DSP}
				//超时休眠重试
				if strings.Contains(errstr, vars.Timeouttxt) || strings.Contains(errstr, vars.Timeouttxt1) {
					logger.Logger.ErrorWithContextMap(ctx, errmap, "dsp getReportStatus")
					time.Sleep(time.Duration(sleepSec))
					continue
				} else {
					logger.Logger.ErrorWithContextMap(ctx, errmap, "dsp getReportStatus")
				}
			}

			logger.Logger.ErrorWithContextMap(ctx, map[string]interface{}{
				"sleep": sleepSec,
				"count": count,
			}, "dsp http post create report error")
			// 当前如果网络不好或者抖动，则添加休眠，避免中间持续频繁请求导致恶化
			time.Sleep(time.Duration(sleepSec))
			return nil, err
		}

		logger.Logger.InfoWithContext(ctx, fmt.Sprintf("dsp报表获取下载地址结果，statusCode:%d, body:%s", resp.StatusCode, string(resp.Body)))

		if resp.StatusCode == http.StatusTooManyRequests {
			sleepSec = rand.Intn(20) + 1 + count*20
			logger.Logger.InfoWithContext(ctx, map[string]interface{}{
				"flag":     "StatusTooManyRequests",
				"count":    count,
				"sleepSec": sleepSec,
				"path":     path,
			})

			time.Sleep(time.Second * time.Duration(sleepSec))
			continue
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			logger.Logger.InfoWithContext(ctx, map[string]interface{}{
				"profileid":  profileId,
				"reportId":   reportId,
				"statuscode": resp.StatusCode,
				"action":     "getReportStatus",
			})

			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				logger.Logger.InfoWithContext(ctx, map[string]interface{}{
					"desc":     "get status get http error code",
					"count":    count,
					"sleepSec": sleepSec,
					"code":     resp.StatusCode,
				})

				return nil, fmt.Errorf("resp error, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body))
			} else {

				sleepSec = rand.Intn(20) + 1 + count*20
				logger.Logger.InfoWithContext(ctx, map[string]interface{}{
					"desc":     "get status get http error code, retry again",
					"count":    count,
					"sleepSec": sleepSec,
				})

				time.Sleep(time.Second * time.Duration(sleepSec))
				continue
			}

		} else {

			var reportStatusBody report.ReportStatusBody
			_, err = resp.Json(&reportStatusBody)
			if err != nil {
				return nil, err
			} else {
				return &reportStatusBody, nil
			}
		}
	}

	logger.Logger.InfoWithContext(ctx, map[string]interface{}{
		"desc": "get dsp report status data error",
	})
	return nil, errors.New("get dsp data error")
}
