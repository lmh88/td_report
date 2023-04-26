package dsp_all_steps

import (
	"context"
	"fmt"
	"math/rand"
	"td_report/app/bean"
	"td_report/app/repo"
	"td_report/boot"
	"td_report/common/redis"
	"td_report/common/reportsystem"
	"td_report/pkg/dsp/audience"
	"td_report/pkg/dsp/detail"
	"td_report/pkg/dsp/dsp_report"
	"td_report/pkg/dsp/inventory"
	"td_report/pkg/dsp/order"
	"td_report/pkg/logger"
	"td_report/pkg/save_file"
	"td_report/vars"
	"time"
)

func DspAllSteps(region, profileId, taskType, dateStr string, wg *reportsystem.Pool, ctx context.Context) (bool, *bean.ReportErr) {
	defer wg.Done()
	var (
		reportId string
		err      error
	)

	errInfo := &bean.ReportErr{
		ReportType: "dsp",
		ReportName: taskType,
		ReportDate: dateStr,
		ProfileId:  profileId,
		KeyParam:   logger.Logger.FromTraceIDContext(ctx),
	}

	switch taskType {
	case "audience":
		reportId, err = audience.CreateAudience(region, dateStr, profileId, ctx)
	case "detail":
		reportId, err = detail.CreateDetail(region, dateStr, profileId, ctx)
	case "inventory":
		reportId, err = inventory.CreateInventory(region, dateStr, profileId, ctx)
	case "order":
		reportId, err = order.CreateOrder(region, dateStr, profileId, ctx)
	default:
		logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
			"desc": "dsp allstep taskType nout found",
		})
		return false, nil
	}

	if err != nil {
		logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
			"flag":      "create dsp report error",
			"err":       err.Error(),
			"region":    region,
			"profileId": profileId,
			"taskType":  taskType,
		})
		// 错误数据记录
		errInfo.ErrorType = repo.ReportErrorTypeOne
		errInfo.ErrorReason = err.Error()
		return false, errInfo
	}

	intn := rand.Intn(4)
	time.Sleep(time.Second * time.Duration(intn))
	count := 1
	for {

		resp, err := dsp_report.GetDspReportStatus(region, profileId, reportId, ctx)
		if err != nil {
			logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
				"flag":      "GetDspReportStatus error",
				"err":       err.Error(),
				"region":    region,
				"profileId": profileId,
				"taskType":  taskType,
				"reportId":  reportId,
			})
			errInfo.ErrorType = repo.ReportErrorTypeTwo
			errInfo.ErrorReason = err.Error()
			return false, errInfo
		}

		if resp.Status == "FAILURE" {
			logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
				"flag":      "dsp report status is FAILURE",
				"region":    region,
				"profileId": profileId,
				"taskType":  taskType,
				"reportId":  reportId,
			})
			errInfo.ErrorType = repo.ReportErrorTypeTwo
			errInfo.ErrorReason = "dsp report status is FAILURE"
			return false, errInfo
		}

		if resp.Status == "SUCCESS" {
			val, err := dsp_report.DownloadDspReport(resp.Location, ctx)
			if err != nil {
				logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
					"flag":      "DownloadDspReport error",
					"resp":      resp,
					"err":       err.Error(),
					"region":    region,
					"profileId": profileId,
					"taskType":  taskType,
					"reportId":  reportId,
				})
				errInfo.ErrorType = repo.ReportErrorTypeThree
				errInfo.ErrorReason = err.Error()
				return false, errInfo
			}

			//写入文件
			err = save_file.SaveDspFile(taskType, dateStr, profileId, val)
			if err != nil {
				logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
					"flag":      "SaveDspFile error",
					"err":       err.Error(),
					"region":    region,
					"profileId": profileId,
					"taskType":  taskType,
					"reportId":  reportId,
				})
				errInfo.ErrorType = repo.ReportErrorTypeThree
				errInfo.ErrorReason = "SaveDspFile error:" + err.Error()
				return false, errInfo
			} else {
				logger.Logger.InfoWithContext(ctx, map[string]interface{}{
					"flag":      "SaveFile success",
					"profileId": reportId,
				})
				return true, nil
			}
		}

		sleepInt := rand.Intn(20) + 20*count
		logger.Logger.ErrorWithContext(ctx, fmt.Sprintf("dsp get error sleep:%d time and count：%d", sleepInt, count))
		time.Sleep(time.Second * time.Duration(sleepInt))
		count += 1

		//兜底
		if count > 40 {
			logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
				"flag":      "count > 40",
				"region":    region,
				"profileId": profileId,
				"taskType":  taskType,
				"reportId":  reportId,
			})
			errInfo.ErrorType = repo.ReportErrorTypeTwo
			errInfo.ErrorReason = "尝试数量超过40次"
			return false, errInfo
		}
	}

}

func DspAllStepsPeriod(region, profileId, taskType, startdate, enddate string, ctx context.Context) {
	var (
		reportId string
		err      error
		key      string
	)

	pipe := boot.RedisCommonClient.GetClient().Pipeline()

	switch taskType {
	case "audience":
		reportId, err = audience.CreateAudiencePeriod(region, startdate, enddate, profileId, ctx)
	case "detail":
		reportId, err = detail.CreateDetailPeriod(region, startdate, enddate, profileId, ctx)
	case "inventory":
		reportId, err = inventory.CreateInventoryPeriod(region, startdate, enddate, profileId, ctx)
	case "order":
		reportId, err = order.CreateOrderPeriod(region, startdate, enddate, profileId, ctx)
	default:
		logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
			"desc": "dsp allstep taskType nout found",
		})
		return
	}

	if err != nil {
		logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
			"flag":      "create dsp reporttool error",
			"err":       err.Error(),
			"region":    region,
			"profileId": profileId,
			"taskType":  taskType,
		})
		return
	}

	count := 1
	for {

		resp, err := dsp_report.GetDspReportStatus(region, profileId, reportId, ctx)
		if err != nil {
			logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
				"flag":      "GetDspReportStatus error",
				"err":       err.Error(),
				"region":    region,
				"profileId": profileId,
				"taskType":  taskType,
				"reportId":  reportId,
			})

			return
		}

		if resp.Status == "FAILURE" {
			logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
				"flag":      "dsp reporttool status is FAILURE",
				"region":    region,
				"profileId": profileId,
				"taskType":  taskType,
				"reportId":  reportId,
			})
			return
		}

		if resp.Status == "SUCCESS" {
			val, err := dsp_report.DownloadDspReport(resp.Location, ctx)
			if err != nil {
				logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
					"flag":      "DownloadDspReport error",
					"resp":      resp,
					"err":       err.Error(),
					"region":    region,
					"profileId": profileId,
					"taskType":  taskType,
					"reportId":  reportId,
				})

				return
			}

			//写入文件
			err = save_file.SaveDspFilePeriod(taskType, startdate, enddate, profileId, val)
			if err != nil {
				logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
					"flag":      "SaveDspFile error",
					"err":       err.Error(),
					"region":    region,
					"profileId": profileId,
					"taskType":  taskType,
					"reportId":  reportId,
				})

				return

			} else {

				fileName := fmt.Sprintf("%s_%s_%s.csv", startdate, enddate, profileId)
				key = fmt.Sprintf("%s_%s", taskType, fileName)
				rediskey := redis.WithDivide(vars.DSP)
				pipe.LPush(rediskey, key)
				pipe.Expire(rediskey, 8*time.Hour)
				_, err = pipe.Exec()
				if err != nil {
					logger.Logger.ErrorWithContext(ctx, err.Error(), "==================redis create divide error")
				}
			}

			break
		}

		sleepInt := rand.Intn(20) + 20*count
		logger.Logger.ErrorWithContext(ctx, fmt.Sprintf("dsp get error sleep:%d time and count：%d", sleepInt, count))
		time.Sleep(time.Second * time.Duration(sleepInt))
		count += 1

		//兜底
		if count > 40 {
			logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
				"flag":      "count > 40",
				"region":    region,
				"profileId": profileId,
				"taskType":  taskType,
				"reportId":  reportId,
			})

			break
		}
	}
}
