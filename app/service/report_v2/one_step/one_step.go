package one_step

import (
	"context"
	"fmt"
	"td_report/app/bean"
	"td_report/app/repo"
	"td_report/app/service/report_v2/adGroups"
	"td_report/app/service/report_v2/asins"
	"td_report/app/service/report_v2/campaigns"
	"td_report/app/service/report_v2/keywords"
	"td_report/app/service/report_v2/productAds"
	"td_report/app/service/report_v2/targets"
	"td_report/app/service/report_v2/varible"
	"td_report/pkg/logger"
	"td_report/vars"
)

// SpOneStep sp第一步请求
func SpOneStep(profileMsg *varible.ProfileMsg, reportType, dateStr string, ctx context.Context) (string, *bean.ReportErr) {

	var (
		reportId string
		err      error
	)
	errInfo := &bean.ReportErr{
		ReportType: vars.SP,
		ReportName: reportType,
		ReportDate: dateStr,
		ProfileId:  profileMsg.ProfileId,
		KeyParam:   logger.Logger.FromTraceIDContext(ctx),
	}

	//todo test
	//time.Sleep(time.Second * 1)
	//return "sp_" + stringx.Rand(), nil
	//if rand.Intn(2) > 0 {
	//	return "sp_" + stringx.Rand(), nil
	//}
	//errInfo.ErrorReason = report_v2.Retry429
	//return "", errInfo

	switch reportType {
	case "keywords":
		reportId, err = keywords.CreateSPKeywordsReport(profileMsg, ctx)
	case "keywordsQuery":
		reportId, err = keywords.CreateSPKeywordsQueryReport(profileMsg, ctx)
	case "adGroups":
		reportId, err = adGroups.CreateSPAdGroupsReport(profileMsg, ctx)
	case "asins":
		reportId, err = asins.CreateSPAsinsReport(profileMsg, ctx)
	case "campaigns", "campaignsPlacement":
		reportId, err = campaigns.CreateSPCampaignsReport(profileMsg, reportType, ctx)
	case "productAds":
		reportId, err = productAds.CreateSPProductAdsReport(profileMsg, ctx)
	case "targetQuerys":
		reportId, err = targets.CreateSPTargetsQueryAdsReport(profileMsg, ctx)
	case "targets":
		reportId, err = targets.CreateSPTargetsAdsReport(profileMsg, ctx)

	default:
		logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
			"flag": fmt.Sprintf("allstep reportType nut found: %s", reportType),
		})
		errInfo.ErrorType = repo.ReportErrorTypeOne
		errInfo.ErrorReason = "report name not exist"
		return "", errInfo
	}

	if err != nil {
		logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
			"flag":       "create report error",
			"err":        err.Error(),
			"reportType": reportType,
			"profileId":  profileMsg.ProfileId,
		})
		errInfo.ErrorType = repo.ReportErrorTypeOne
		errInfo.ErrorReason = err.Error()
		return "", errInfo
	}
	return reportId, nil
}

func SbOneStep(profileMsg *varible.ProfileMsg, reportType, dateStr string, ctx context.Context) (string, *bean.ReportErr) {

	var (
		reportId string
		err      error
	)

	errInfo := &bean.ReportErr{
		ReportType: vars.SB,
		ReportName: reportType,
		ReportDate: dateStr,
		ProfileId:  profileMsg.ProfileId,
		ErrorType:  repo.ReportErrorTypeOne,
		KeyParam:   logger.Logger.FromTraceIDContext(ctx),
	}

	//todo test
	//time.Sleep(time.Second * 1)
	//return "sb_" + stringx.Rand(), nil
	//if rand.Intn(2) > 0 {
	//	return "sb_" + stringx.Rand(), nil
	//}
	//errInfo.ErrorReason = report_v2.Retry429
	//return "", errInfo

	switch reportType {
	case "adGroup":
		reportId, err = adGroups.CreateSbAdGroupsReport(profileMsg, ctx)
	case "adGroupVideo":
		reportId, err = adGroups.CreateSbAdGroupsVideoReport(profileMsg, ctx)
	case "campaigns":
		reportId, err = campaigns.CreateSbCampaignsReport(profileMsg, ctx)
	case "campaignsVideo":
		reportId, err = campaigns.CreateSbCampaignsVideoReport(profileMsg, ctx)
	case "keywords":
		reportId, err = keywords.CreateSbKeywordsReport(profileMsg, ctx)
	case "keywordsVideo":
		reportId, err = keywords.CreateSbKeywordsVideoReport(profileMsg, ctx)
	case "keywordsQuery":
		reportId, err = keywords.CreateSbKeywordsQueryReport(profileMsg, ctx)
	case "keywordsQueryVideo":
		reportId, err = keywords.CreateSbKeywordsQueryVideoReport(profileMsg, ctx)
	case "targets":
		reportId, err = targets.CreateSbTargetsReport(profileMsg, ctx)
	case "targetsVideo":
		reportId, err = targets.CreateSbTargetsVideoReport(profileMsg, ctx)
	case "ads":
		reportId, err = productAds.CreateSbAdsReport(profileMsg, ctx)
	default:
		logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
			"flag": fmt.Sprintf("allstep reportType nout found: %s", reportType),
		})
		errInfo.ErrorReason = "report name not exist"
		return "", errInfo
	}

	if err != nil {
		logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
			"flag":      "create report error",
			"err":       err.Error(),
			"profileId": profileMsg.ProfileId,
			"reportId":  reportId,
		})
		errInfo.ErrorReason = err.Error()
		return "", errInfo
	}
	return reportId, nil
}

func SdOneStep(profileMsg *varible.ProfileMsg, reportType, tactic string, ctx context.Context) (string, *bean.ReportErr) {

	var (
		reportId string
		err      error
	)

	errInfo := &bean.ReportErr{
		ReportType: vars.SD,
		ReportName: reportType,
		ReportDate: profileMsg.ReportDate,
		ProfileId:  profileMsg.ProfileId,
		ErrorType:  repo.ReportErrorTypeOne,
		KeyParam:   logger.Logger.FromTraceIDContext(ctx),
		Extra:      tactic,
	}

	//todo test
	//time.Sleep(time.Second)
	//time.Sleep(time.Millisecond * 100)
	//return "sd_" + stringx.Rand(), nil
	//rand.Seed(time.Now().UnixNano())
	//if rand.Intn(2) > 0 {
	//	return "sd_" + stringx.Rand(), nil
	//}
	//errInfo.ErrorReason = varible.Retry429
	//return "", errInfo

	switch reportType {
	case "adGroups":
		reportId, err = adGroups.CreateSdAdGroupsReport(profileMsg, tactic, ctx)
	case "asins":
		reportId, err = asins.CreateSdAsinsReport(profileMsg, tactic, ctx)
	case "campaigns":
		reportId, err = campaigns.CreateSdCampaignsReport(profileMsg, tactic, ctx)
	case "productAds":
		reportId, err = productAds.CreateSdProductAdsReport(profileMsg, tactic, ctx)
	case "targets":
		reportId, err = targets.CreateSdTargetsReport(profileMsg, tactic, ctx)
	default:
		logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
			"flag": fmt.Sprintf("allstep reportType nout found: %s", reportType),
		})
		errInfo.ErrorReason = "report name not exist"
		return "", errInfo
	}

	if err != nil {
		logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
			"flag":      "create report error",
			"err":       err.Error(),
			"profileId": profileMsg.ProfileId,
			"reportId":  reportId,
		})
		errInfo.ErrorReason = err.Error()
		return "", errInfo
	}
	return reportId, nil
}
