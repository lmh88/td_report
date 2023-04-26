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
	"td_report/pkg/brand_metrics"
	"td_report/pkg/campaigns"
	"td_report/pkg/keywords"
	"td_report/pkg/logger"
	"td_report/pkg/productAds"
	"td_report/pkg/report"
	"td_report/pkg/save_file"
	"td_report/pkg/targets"
	"td_report/vars"
	"time"
)

func SbAllSteps(profileToken *bean.ProfileToken, reportType, dateStr, enddate string, wg *reportsystem.Pool, ctx context.Context) (bool, *bean.ReportErr) {
	defer wg.Done()

	var (
		reportId string
		err      error
	)

	errInfo := &bean.ReportErr{
		ReportType: "sb",
		ReportName: reportType,
		ReportDate: dateStr,
		ProfileId:  profileToken.ProfileId,
		KeyParam:   logger.Logger.FromTraceIDContext(ctx),
	}
	switch reportType {
	case "adGroup":
		reportId, err = adGroups.CreateSbAdGroupsReport(profileToken, dateStr, ctx)
	case "adGroupVideo":
		reportId, err = adGroups.CreateSbAdGroupsVideoReport(profileToken, dateStr, ctx)
	case "campaigns":
		reportId, err = campaigns.CreateSbCampainsReport(profileToken, dateStr, ctx)
	case "campaignsVideo":
		reportId, err = campaigns.CreateSbCampainsVideoReport(profileToken, dateStr, ctx)
	case "keywords":
		reportId, err = keywords.CreateSbKeywordsReport(profileToken, dateStr, ctx)
	case "keywordsVideo":
		reportId, err = keywords.CreateSbKeywordsVideoReport(profileToken, dateStr, ctx)
	case "keywordsQuery":
		reportId, err = keywords.CreateSbKeywordsQueryReport(profileToken, dateStr, ctx)
	case "keywordsQueryVideo":
		reportId, err = keywords.CreateSbKeywordsQueryVideoReport(profileToken, dateStr, ctx)
	case "targets":
		reportId, err = targets.CreateSbTargetsReport(profileToken, dateStr, ctx)
	case "targetsVideo":
		reportId, err = targets.CreateSbTargetsVideoReport(profileToken, dateStr, ctx)
	case vars.BrandMetricsWeekly:
		reportId, err = brand_metrics.CreateSbBrandMetricsReportPost(profileToken, dateStr, enddate, "1w", ctx)
	case vars.BrandMetricsMonthly:
		reportId, err = brand_metrics.CreateSbBrandMetricsReportPost(profileToken, dateStr, enddate, "1cm", ctx)
	default:
		logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
			"flag": fmt.Sprintf("allstep reportType nout found: %s", reportType),
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
			"reportId":  reportId,
		})
		errInfo.ErrorType = repo.ReportErrorTypeOne
		errInfo.ErrorReason = err.Error()
		return false, errInfo
	}

	//获取报表id到获取报表地址之间休眠
	intn := rand.Intn(4)
	time.Sleep(time.Second * time.Duration(intn))
	var num = 0
	for {

		if num > 40 {
			logger.Logger.ErrorWithContext(ctx, "retry too many times to break")
			errInfo.ErrorType = repo.ReportErrorTypeTwo
			errInfo.ErrorReason = "尝试次数超过40次"
			return false, errInfo
		}

		var (
			resp *report.ReportStatusBody
			val  []byte
		)
		if reportType == vars.BrandMetricsWeekly || reportType == vars.BrandMetricsMonthly {
			resp, err = report.GetReportSbbrandMetricsGet(profileToken, reportId, ctx)
		} else {
			resp, err = report.GetReportStatus(profileToken, reportId, "sb", ctx)
		}

		if err != nil {
			logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
				"flag":         "GetReportStatus error",
				"err":          err.Error(),
				"profileToken": profileToken.ProfileId,
				"reportId":     reportId,
			})
			errInfo.ErrorType = repo.ReportErrorTypeTwo
			errInfo.ErrorReason = err.Error()
			return false, errInfo
		}

		if resp.Status == "FAILURE" {
			logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
				"flag":         "report status is FAILURE",
				"profileToken": profileToken.ProfileId,
				"reportType":   reportType,
				"dateStr":      dateStr,
			})
			errInfo.ErrorType = repo.ReportErrorTypeTwo
			errInfo.ErrorReason = "report status is FAILURE"
			return false, errInfo
		}

		if resp.Status == "SUCCESS" || resp.Status == "SUCCESSFUL" {
			if resp.Location == "" {
				logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
					"flag":         "DownloadReport error, address lost",
					"resp":         resp.ReportId,
					"profileToken": profileToken.ProfileId,
				})
				errInfo.ErrorType = repo.ReportErrorTypeTwo
				errInfo.ErrorReason = "report status is FAILURE"
				return false, errInfo
			}

			if reportType == vars.BrandMetricsWeekly || reportType == vars.BrandMetricsMonthly {
				val, err = report.DownloadReportWithoutAuth(resp.Location, ctx)
			} else {
				val, err = report.DownloadReport(profileToken, resp.Location, ctx)
			}

			if err != nil {
				logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
					"flag":         "DownloadReport error",
					"resp":         resp.ReportId,
					"profileToken": profileToken.ProfileId,
					"err":          err.Error(),
				})
				errInfo.ErrorType = repo.ReportErrorTypeThree
				errInfo.ErrorReason = err.Error()
				return false, errInfo
			} else {
				logger.Logger.InfoWithContext(ctx, map[string]interface{}{
					"flag":         "DownloadReport success",
					"resp":         resp.ReportId,
					"profileToken": profileToken.ProfileId,
				})
			}

			//写入文件
			err = save_file.SaveSBFile(reportType, dateStr, profileToken.ProfileId, val)
			if err != nil {
				logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
					"flag":      "SaveFile error",
					"err":       err.Error(),
					"profileId": profileToken.ProfileId,
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

		logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
			"flag":       "wait resp.Status",
			"reportId":   resp.ReportId,
			"reportType": reportType,
			"dateStr":    dateStr,
		})

		num = num + 1
		time.Sleep(time.Second * 3)
	}

}

func SPAllSteps(profileToken *bean.ProfileToken, reportType, dateStr string, wg *reportsystem.Pool, ctx context.Context) (bool, *bean.ReportErr) {
	defer wg.Done()

	var (
		reportId string
		err      error
	)
	errInfo := &bean.ReportErr{
		ReportType: "sp",
		ReportName: reportType,
		ReportDate: dateStr,
		ProfileId:  profileToken.ProfileId,
		KeyParam:   logger.Logger.FromTraceIDContext(ctx),
	}

	switch reportType {
	case "keywords":
		reportId, err = keywords.CreateSPKeywordsReport(profileToken, dateStr, ctx)
	case "keywordsQuery":
		reportId, err = keywords.CreateSPKeywordsQueryReport(profileToken, dateStr, ctx)
	case "adGroups":
		reportId, err = adGroups.CreateSPAdGroupsReport(profileToken, dateStr, ctx)
	case "asins":
		reportId, err = asins.CreateSPAsinsReport(profileToken, dateStr, ctx)
	case "campaigns", "campaignsPlacement":
		reportId, err = campaigns.CreateSPCampaignsReport(profileToken, dateStr, reportType, ctx)
	case "productAds":
		reportId, err = productAds.CreateSPProductAdsReport(profileToken, dateStr, ctx)
	case "targetQuerys":
		reportId, err = targets.CreateSPTargetsQueryAdsReport(profileToken, dateStr, ctx)
	case "targets":
		reportId, err = targets.CreateSPTargetsAdsReport(profileToken, dateStr, ctx)

	default:
		logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
			"flag": fmt.Sprintf("allstep reportType nut found: %s", reportType),
		})
		errInfo.ErrorType = repo.ReportErrorTypeOne
		errInfo.ErrorReason = err.Error()
		return false, errInfo
	}

	if err != nil {
		logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
			"flag":       "create report error",
			"err":        err.Error(),
			"reportType": reportType,
			"profileid":  profileToken.ProfileId,
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

		resp, err := report.GetReportStatus(profileToken, reportId, "sp", ctx)
		if err != nil {
			logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
				"flag":      "GetReportStatus error",
				"err":       err.Error(),
				"profileId": profileToken.ProfileId,
				"reportId":  reportId,
			})
			errInfo.ErrorType = repo.ReportErrorTypeTwo
			errInfo.ErrorReason = err.Error()
			return false, errInfo
		}

		if resp.Status == "FAILURE" {
			logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
				"flag":       "report status is FAILURE",
				"profileId":  profileToken.ProfileId,
				"reportType": reportType,
				"dateStr":    dateStr,
			})
			errInfo.ErrorType = repo.ReportErrorTypeTwo
			errInfo.ErrorReason = "report status is FAILURE"
			return false, errInfo
		}

		if resp.Status == "SUCCESS" {
			val, err := report.DownloadReport(profileToken, resp.Location, ctx)
			if err != nil {
				logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
					"flag":     "DownloadReport error",
					"reportId": resp.ReportId,
					"profilid": profileToken.ProfileId,
					"err":      err.Error(),
				})
				errInfo.ErrorType = repo.ReportErrorTypeThree
				errInfo.ErrorReason = err.Error()
				return false, errInfo
			}

			//写入文件
			err = save_file.SaveSPFile(reportType, dateStr, profileToken.ProfileId, val)
			if err != nil {
				logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
					"flag":      "SaveFile error",
					"err":       err.Error(),
					"profileId": profileToken.ProfileId,
				})
				errInfo.ErrorType = repo.ReportErrorTypeOne
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
			"ReportId":   resp.ReportId,
			"reportType": reportType,
			"dateStr":    dateStr,
		})

		num = num + 1
		time.Sleep(time.Second * 3)
	}
}
