package all_steps

import (
	"context"
	"fmt"
	"math/rand"
	"td_report/app/bean"
	"td_report/app/repo"
	"td_report/common/reportsystem"
	"td_report/pkg/adGroups"
	"td_report/pkg/asins"
	"td_report/pkg/campaigns"
	"td_report/pkg/logger"
	"td_report/pkg/productAds"
	"td_report/pkg/report"
	"td_report/pkg/save_file"
	"td_report/pkg/targets"
	"time"
)

func SDAllSteps(profileToken *bean.ProfileToken, reportType, dateStr, tactic string, wg *reportsystem.Pool, ctx context.Context) (bool, *bean.ReportErr) {
	defer wg.Done()

	var (
		reportId string
		err      error
	)

	errInfo := &bean.ReportErr{
		ReportType: "sd",
		ReportName: reportType,
		ReportDate: dateStr,
		ProfileId:  profileToken.ProfileId,
		KeyParam:   logger.Logger.FromTraceIDContext(ctx),
	}

	switch reportType {
	case "adGroups":
		reportId, err = adGroups.CreateSdAdGroupsReport(profileToken, dateStr, tactic, ctx)
	case "asins":
		reportId, err = asins.CreateSdAsinsReport(profileToken, dateStr, tactic, ctx)
	case "campaigns":
		reportId, err = campaigns.CreateSdCampaignsReport(profileToken, dateStr, tactic, ctx)
	case "productAds":
		reportId, err = productAds.CreateSdProductAdsReport(profileToken, dateStr, tactic, ctx)
	case "targets":
		reportId, err = targets.CreateSdTargetsReport(profileToken, dateStr, tactic, ctx)

	default:
		logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
			"flag": fmt.Sprintf("allstep reportType not found: %s", reportType),
		})
		errInfo.ErrorType = repo.ReportErrorTypeOne
		errInfo.ErrorReason = err.Error()
		return false, errInfo
	}

	if err != nil {
		logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
			"flag":      "create report error",
			"err":       err.Error(),
			"profileid": profileToken.ProfileId,
		})
		errInfo.ErrorType = repo.ReportErrorTypeOne
		errInfo.ErrorReason = err.Error()
		return false, errInfo
	}

	intn := rand.Intn(4)
	time.Sleep(time.Second * time.Duration(intn))
	var num = 0
	for {

		if num > 40 {
			logger.Logger.ErrorWithContext(ctx, "retry too many times to break")
			errInfo.ErrorType = repo.ReportErrorTypeTwo
			errInfo.ErrorReason = "尝试数超过40次"
			return false, errInfo
		}
		resp, err := report.GetReportStatus(profileToken, reportId, "sd", ctx)
		if err != nil {
			logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
				"flag":         "GetReportStatus error",
				"err":          err.Error(),
				"profileToken": profileToken,
				"reportId":     reportId,
			})
			errInfo.ErrorType = repo.ReportErrorTypeTwo
			errInfo.ErrorReason = err.Error()
			return false, errInfo
		}

		if resp.Status == "FAILURE" {
			logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
				"flag":         "report status is FAILURE",
				"profileToken": profileToken,
				"reportType":   reportType,
				"dateStr":      dateStr,
			})
			errInfo.ErrorType = repo.ReportErrorTypeTwo
			errInfo.ErrorReason = "report status is FAILURE"
			return false, errInfo
		}

		if resp.Status == "SUCCESS" {
			val, err := report.DownloadReport(profileToken, resp.Location, ctx)
			if err != nil {
				logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
					"flag":         "DownloadReport error",
					"resp":         resp,
					"profileToken": profileToken,
					"err":          err.Error(),
				})
				errInfo.ErrorType = repo.ReportErrorTypeThree
				errInfo.ErrorReason = err.Error()
				return false, errInfo
			}

			//写入文件
			err = save_file.SaveSDFile(reportType, dateStr, profileToken.ProfileId, tactic, val)
			if err != nil {
				logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
					"flag":         "SaveFile error",
					"err":          err.Error(),
					"profileToken": profileToken,
				})
				errInfo.ErrorType = repo.ReportErrorTypeThree
				errInfo.ErrorReason = "SaveFile error:" + err.Error()
				return false, errInfo
			} else {
				logger.Logger.InfoWithContext(ctx, map[string]interface{}{
					"flag":      "SaveFile success",
					"profileId": profileToken.ProfileId,
				})

				return true, nil
			}
		}

		logger.Logger.InfoWithContext(ctx, map[string]interface{}{
			"flag":       "wait resp.Status",
			"resp":       resp,
			"reportType": reportType,
			"dateStr":    dateStr,
		})

		num = num + 1
		time.Sleep(time.Second * 3)
	}
}
