package create_report

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/os/gtime"
	"net/http"
	"strconv"
	"strings"
	"td_report/app/repo"
	"td_report/app/service/report_v2/product_server"
	"td_report/app/service/report_v2/varible"
	"td_report/boot"
	"td_report/pkg/limiter"
	"td_report/pkg/logger"
	region2 "td_report/pkg/region"
	"td_report/pkg/requests"
	"td_report/vars"
	"time"
)

func SBCreateReport(profileMsg *varible.ProfileMsg, reportType, segment, creativeType, metrics string, ctx context.Context) (string, error) {
	path := fmt.Sprintf("/v2/hsa/%s/report", reportType)
	return CreateReport(profileMsg, reportType, segment, creativeType, metrics, path, vars.SB, ctx)
}

func SPCreateReport(profileMsg *varible.ProfileMsg, reportType, segment, metrics string, ctx context.Context) (string, error) {
	var path string
	if reportType == "campaignsPlacement" {
		path = fmt.Sprintf("/v2/sp/%s/report", "campaigns")
	} else {
		path = fmt.Sprintf("/v2/sp/%s/report", reportType)
	}

	return CreateReport(profileMsg, reportType, segment, "", metrics, path, vars.SP, ctx)
}

func SDCreateReport(profileMsg *varible.ProfileMsg, reportType, tactic, metrics string, ctx context.Context) (string, error) {
	path := fmt.Sprintf("/sd/%s/report", reportType)
	return CreateReport(profileMsg, reportType, tactic, "", metrics, path, vars.SD, ctx)
}

func CreateReport(profileMsg *varible.ProfileMsg, reportType, segment, creativeType, metrics, path, report string, ctx context.Context) (string, error) {
	endpoint, ok := region2.ApiUrl[profileMsg.Region]
	if !ok {
		return "", fmt.Errorf("region not found: %s", profileMsg.Region)
	}

	// 添加一道过滤，避免短时间内重复请求token
	if exists, _ := vars.Cache.Contains(profileMsg.RefreshToken); exists {
		logger.Logger.InfoWithContext(ctx, "error token")
		return "", errors.New("error token")
	}

	url := endpoint + path
	//todo test
	crt := varible.ClientRefreshToken{
		ProfileId: profileMsg.ProfileId,
		RefreshToken: profileMsg.RefreshToken,
		ClientId: profileMsg.ClientId,
		ClientSecret: profileMsg.ClientSecret,
	}
	headers, err := FakeHeaders(&crt)
	if err != nil {
		//针对token失效的情况，缓存下来，直接跳过
		vars.Cache.SetIfNotExist(profileMsg.RefreshToken, 1, time.Second*600)
		logger.Logger.ErrorWithContext(ctx, err, reportType, report, "get header error")
		return "", err
	}

	params := map[string]interface{}{
		"reportDate": profileMsg.ReportDate,
		"metrics":    metrics,
	}

	if report != "sd" {
		//campaigns
		if report == vars.SP && reportType == "campaigns" {
		} else {
			params["segment"] = segment
		}

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
	var reportId string
	limit := limiter.PerSecond(vars.LimitMap[report])
	limitKey := fmt.Sprintf("%s_%s_%s", report, profileMsg.ClientTag, "create_report")
	for {

		//添加频率率的限制，避免请求频次太多
		for {
			res, err := boot.Rlimiter.Allow(limitKey, limit)
			if err != nil {
				return "", err
			}

			if res.Allowed == 1 {
				break
			}

			time.Sleep(10 * time.Millisecond)
		}

		//todo test
		resp, err = requests.Post(url, requests.WithHeaders(headers), requests.WithJson(params), requests.WithTimeout(time.Second*40))
		if err != nil || resp == nil {
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
					count++
					if count >= 3 {
						return "", err
					}
					time.Sleep(time.Second)
					continue
				} else {
					errmap := map[string]interface{}{"err": errstr}
					logger.Logger.ErrorWithContextMap(ctx, errmap, "create report get post data error")
				}
			}

			return "", err
		}

		logger.Logger.InfoWithContext(ctx, fmt.Sprintf("%s报表第一次请求结果，statusCode:%d, body:%s", report, resp.StatusCode, string(resp.Body)))

		if resp.StatusCode == http.StatusTooManyRequests {
			logger.Logger.ErrorWithContext(ctx, "StatusTooManyRequests")
			return "", errors.New(varible.Retry429)
		}

		// 除了429之外的错误信息直接返回错误信息
		if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				da, _ := json.Marshal(params)
				logger.Logger.ErrorWithContext(ctx, fmt.Sprintf("resp error, StatusCode=%d, body=%s, parama=%s", resp.StatusCode, string(resp.Body), string(da)))
				return "", fmt.Errorf("resp error, StatusCode=%d, body=%s", resp.StatusCode, string(resp.Body))
			}
		}

		break
	}

	val, err := resp.JsonAndValueIsAny()
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
			"desc": "json error",
			"error": err.Error(),
		})
		return "", err
	}

	reportId, ok = val["reportId"].(string)
	if !ok {
		return "", errors.New("reportId is empty ,get reportId error. ")
	} else {
		return reportId, nil
	}
}

func BatchStart(reportIdMsg *varible.ReportIdMsg) {
	rds := boot.RedisCommonClient.GetClient()
	rds.SAdd(reportIdMsg.BatchKey, getBatchVal(reportIdMsg))

	batchStartTime := varible.GetBatchTimeKey(reportIdMsg.BatchKey)
	if rds.SetNX(batchStartTime, gtime.Timestamp(), 24 * time.Hour).Val() {
		rds.Expire(reportIdMsg.BatchKey, 24*time.Hour)
		if checkBatchKey(reportIdMsg.BatchKey) {
			repo.NewReportSchduleRepository().StartSchdule(reportIdMsg.BatchKey)
		}
	}
}

func BatchEnd(reportIdMsg *varible.ReportIdMsg) {
	rds := boot.RedisCommonClient.GetClient()
	rds.SRem(reportIdMsg.BatchKey, getBatchVal(reportIdMsg))
	if rds.Exists(reportIdMsg.BatchKey).Val() == 0 {
		//通知大数据
		//s.ReportTaskService.AddMutile(reportIdMsg.ReportType, reportIdMsg.ReportDate)
		//防止并发修改
		if rds.SetNX(reportIdMsg.BatchKey + "_end", 1, time.Second * 3).Val() {
			if checkBatchKey(reportIdMsg.BatchKey) {
				repo.NewReportSchduleRepository().EndSchdule(reportIdMsg.BatchKey)
			} else {
				product_server.CompleteBatch(reportIdMsg.BatchKey)
			}
		}
		return
	}
	//OvertimeNotice(reportIdMsg)
}

func checkBatchKey(key string) bool {
	return strings.Contains(key, "report:")
}

func getBatchVal(reportIdMsg *varible.ReportIdMsg) string {
	if checkBatchKey(reportIdMsg.BatchKey) {
		return fmt.Sprintf("%s:%s:%s", reportIdMsg.ProfileId, reportIdMsg.ReportName, reportIdMsg.ReportTactic)
	}
	return fmt.Sprintf("%s:%s:%s:%s", reportIdMsg.ProfileId, reportIdMsg.ReportName, reportIdMsg.ReportTactic, reportIdMsg.ReportDate)
}

func OvertimeNotice(reportIdMsg *varible.ReportIdMsg) {
	rds := boot.RedisCommonClient.GetClient()
	noticeKey := "notice:" + reportIdMsg.BatchKey
	batchTimeKey := varible.GetBatchTimeKey(reportIdMsg.BatchKey)
	batchTimeVal := rds.Get(batchTimeKey).Val()
	batchTime, _ := strconv.Atoi(batchTimeVal)
	timestamp := gtime.Timestamp()
	timespan := timestamp - int64(batchTime)
	overtime := varible.OvertimeNoticeMap[reportIdMsg.ReportType]
	//超过时，通知大数据
	if timespan > overtime {
		if rds.SetNX(noticeKey, timestamp, time.Hour * 8).Val() {
			//ReportTaskService.AddMutile(reportIdMsg.ReportType, reportIdMsg.ReportDate)
		}
	}
}
