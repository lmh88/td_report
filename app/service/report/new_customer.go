package report

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/google/wire"
	"github.com/guonaihong/gout"
	"strings"
	"td_report/app/bean"
	"td_report/app/repo"
	"td_report/common/cryptodata"
	"td_report/common/reportsystem"
	"td_report/common/tool"
	"td_report/pkg/all_steps"
	"td_report/pkg/common"
	"td_report/pkg/logger"
	"td_report/pkg/save_file"
	"td_report/vars"
	"time"
)

var CustomerReportServiceSet = wire.NewSet(wire.Struct(new(CustomerReportService), "*"))

// 在dev或者pre环境限制拉数据，只是通知
const newdevelopModel = "prod"

type CustomerReportService struct {
	SellerProfileRepository     *repo.SellerProfileRepository
	ProfileRepository           *repo.ProfileRepository
	ReportTaskRepository        *repo.ReportTaskRepository
	ReportBatchRepository       *repo.ReportBatchRepository
	ReportBatchDetailRepository *repo.ReportBatchDetailRepository
}

func NewCustomerReportService(SellerProfileRepository *repo.SellerProfileRepository,
	ProfileRepository *repo.ProfileRepository,
	ReportTaskRepository *repo.ReportTaskRepository,
	ReportBatchRepository *repo.ReportBatchRepository,
	ReportBatchDetailRepository *repo.ReportBatchDetailRepository,
) *CustomerReportService {
	return &CustomerReportService{
		SellerProfileRepository:     SellerProfileRepository,
		ProfileRepository:           ProfileRepository,
		ReportTaskRepository:        ReportTaskRepository,
		ReportBatchRepository:       ReportBatchRepository,
		ReportBatchDetailRepository: ReportBatchDetailRepository,
	}
}

func (t *CustomerReportService) getPpcProfileToken(profileIdList []string) ([]*bean.ProfileToken, error) {
	return t.SellerProfileRepository.ListProfileAndRefreshtoken(profileIdList)
}

func (t *CustomerReportService) getDspProfileToken(profileIdList []string) ([]*bean.DspRegionProfile, error) {
	return t.ProfileRepository.ListDspRegionProfile(profileIdList)
}

func (t *CustomerReportService) Notice(id int, desc string, status int, time int64, url string) (error, string) {
	code := 0
	key := g.Cfg().GetString("common.noticekey")
	singn := cryptodata.Sha1(fmt.Sprintf("%s%d%d%d%s", desc, id, status, time, key))
	result := new(map[string]interface{})
	logger.Logger.Info(gout.H{
		"id":     id,
		"desc":   desc,
		"status": status,
		"t":      time,
		"sign":   singn,
	})
	err := gout.
		// 设置POST方法和url
		POST(url).
		//打开debug模式
		// 设置非结构化数据到http body里面
		// 设置json需使用SetJSON
		SetWWWForm(
			gout.H{
				"id":     id,
				"desc":   desc,
				"status": status,
				"t":      time,
				"sign":   singn,
			},
		).
		BindJSON(&result).
		Code(&code).
		//结束函数
		Do()

	if err != nil {
		logger.Logger.Error(err)
		return err, ""
	}
	if code != 200 {
		re, _ := json.Marshal(result)
		logger.Logger.Error(re)
		return errors.New("error"), string(re)
	} else {
		logger.Logger.Info("send success")
		return nil, ""
	}
}

func (t *CustomerReportService) AddNotice(ReportType string, ReportName string, StartDate string) (bool, string) {
	var result = true
	path := ""
	commonPath := g.Cfg().GetString("common.datapath")
	switch ReportType {
	case vars.SP:
		path = commonPath + vars.SpPath

	case vars.SD:
		path = commonPath + vars.SdPath

	case vars.SB:
		path = commonPath + vars.SbPath
	case vars.DSP:
		path = commonPath + vars.DspPath
	default:
		logger.Info(map[string]interface{}{
			"flag":        "reportType not found",
			"*reportType": ReportType,
		})

	}

	var localDir string
	var reportDate string
	if ReportName == vars.BrandMetricsMonthly || ReportName == vars.BrandMetricsWeekly {
		localDir = save_file.GetSbPath(ReportName, StartDate)
		splist := strings.Split(localDir, "/")
		length := len(splist)
		reportDate = splist[length-1]
	} else {
		localDir = fmt.Sprintf("%s/%s/%s", path, ReportName, StartDate)
		reportDate = StartDate
	}

	num := common.GetDirFileNum(localDir)
	if num > 0 {
		t.ReportTaskRepository.Add(localDir, ReportName, ReportType, reportDate)
		logger.Logger.Info(map[string]interface{}{
			"flag":       "profile report done",
			"dateStr":    StartDate,
			"reportType": ReportType,
			"reportName": ReportName,
		})

	} else {

		logger.Logger.Info(map[string]interface{}{
			"flag":       "profile report  dir is empty",
			"dateStr":    StartDate,
			"reportType": ReportType,
			"reportName": ReportName,
		})

		result = false
	}

	return result, localDir
}

func (t *CustomerReportService) Report(report *bean.ReportData) error {
	// 避免在测试和pre环境拉数据
	if g.Cfg().GetString("server.Env") != developModel {
		t.Notice(report.ProcessId, "推送成功", 2, time.Now().Unix(), report.CallBackUrl)
		return nil
	}
	tokenList, err := t.getPpcProfileToken(report.Profileids)
	if err != nil {
		logger.Logger.Error("error1======", err.Error())
		return err
	}

	daylist, err := tool.GetDays(report.StartDate, report.EndDate, vars.TimeLayout)
	if err != nil {
		logger.Logger.Error("get daylist error", err.Error())
		return err
	}

	AllReport := make([]string, 0)
	// 如果是全部的，则获取全部对应的报表名称
	if report.ReportDataType == 0 {
		report.ReportName = vars.ReportList[report.ReportType]
		if report.ReportType == vars.SB {
			AllReport = vars.ReportList[vars.SB]
		} else {
			AllReport = report.ReportName
		}
	}

	reportListStr := strings.Join(AllReport, ",")
	id, err := t.ReportBatchDetailRepository.Addone(report.Batch, reportListStr, report.ReportType, report.StartDate, report.EndDate, vars.TimeLayout)
	if err != nil {
		logger.Logger.Error(err)
		return err
	}

	pool := reportsystem.NewPool(15)
	for _, dayitem := range daylist {
		for _, reportName := range report.ReportName {
			for _, item := range tokenList {
				switch report.ReportType {
				case "sp":
					pool.Add(1)
					ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s", "sp", reportName, item.ProfileId))
					go all_steps.SPAllSteps(item, reportName, dayitem, pool, ctx)
					//todo 权重写入
				case "sd":
					for _, tactic := range []string{"T00020", "T00030"} {
						pool.Add(1)
						ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s", "sd", reportName, item.ProfileId))
						go all_steps.SDAllSteps(item, reportName, dayitem, tactic, pool, ctx)
					}
				case "sb":
					pool.Add(1)
					ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s", "sd", reportName, item.ProfileId))
					go all_steps.SbAllSteps(item, reportName, dayitem, dayitem, pool, ctx)

				default:
					logger.Logger.Error(map[string]interface{}{
						"flag":       "reportType not found",
						"reportType": report.ReportType,
					})
				}
			}

		}
	}

	pool.Wait()
	var reason = fmt.Sprintf("report_type: %s ", report.ReportType)
	var temp = 0
	var total = 0
	var status = 2
	var str = ""
	var sbbrandreason = ""
	var reportList []string
	for _, dayitem := range daylist {
		temp = 0
		reportList = make([]string, 0)
		for _, reportName := range report.ReportName {
			re, _ := t.AddNotice(report.ReportType, reportName, dayitem)
			if re == false {
				temp = temp + 1
				total = total + 1
				reportList = append(reportList, reportName)

			}
		}
		if temp > 0 {
			str = str + fmt.Sprintf(" day:%s, reportName:%s", dayitem, strings.Join(reportList, ","))
		}
	}

	if total != 0 {
		status = 3
		reason = reason + str + " dir is empty"
	} else {
		reason = ""
	}

	if status == 3 {
		reason = "拉取对应的报表失败"
	} else {
		reason = "操作成功"
	}

	t.ReportBatchDetailRepository.ChangeStatus(status, id, sbbrandreason)
	t.ReportBatchRepository.ChangeStatus(status, report.Batch)
	t.Notice(report.ProcessId, reason, status, time.Now().Unix(), report.CallBackUrl)
	return nil
}
