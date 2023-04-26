package report

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"td_report/app/bean"
	"td_report/app/repo"
	"td_report/common/cryptodata"
	"td_report/common/reportsystem"
	"td_report/common/tool"
	"td_report/pkg/all_steps"
	"td_report/pkg/common"
	"td_report/pkg/dsp/dsp_all_steps"
	"td_report/pkg/logger"
	"td_report/pkg/save_file"
	"td_report/vars"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/google/wire"
	"github.com/guonaihong/gout"
	amqp "github.com/rabbitmq/amqp091-go"
)

var ToolReportServiceSet = wire.NewSet(wire.Struct(new(ToolReportService), "*"))

// 在dev或者pre环境限制拉数据，只是通知
const developModel = "prod"

type ToolReportService struct {
	SellerProfileRepository     *repo.SellerProfileRepository
	ProfileRepository           *repo.ProfileRepository
	ReportTaskRepository        *repo.ReportTaskRepository
	ReportBatchRepository       *repo.ReportBatchRepository
	ReportBatchDetailRepository *repo.ReportBatchDetailRepository
}

func NewToolReportService(SellerProfileRepository *repo.SellerProfileRepository,
	ProfileRepository *repo.ProfileRepository,
	ReportTaskRepository *repo.ReportTaskRepository,
	ReportBatchRepository *repo.ReportBatchRepository,
	ReportBatchDetailRepository *repo.ReportBatchDetailRepository,
) *ToolReportService {
	return &ToolReportService{
		SellerProfileRepository:     SellerProfileRepository,
		ProfileRepository:           ProfileRepository,
		ReportTaskRepository:        ReportTaskRepository,
		ReportBatchRepository:       ReportBatchRepository,
		ReportBatchDetailRepository: ReportBatchDetailRepository,
	}
}

func (t *ToolReportService) getPpcProfileToken(profileIdList []string) ([]*bean.ProfileToken, error) {
	return t.SellerProfileRepository.ListProfileAndRefreshtoken(profileIdList)
}

func (t *ToolReportService) Notice(id int, desc string, status int, time int64, url string) (error, string) {
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

func (t *ToolReportService) AddNotice(ReportType string, ReportName string, StartDate string) (bool, string) {
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
		logger.Logger.Info(map[string]interface{}{
			"flag":       "reportType not found",
			"reportType": ReportType,
		})

	}

	var localDir string
	//var reportDate string
	if ReportName == vars.BrandMetricsMonthly || ReportName == vars.BrandMetricsWeekly {
		localDir = save_file.GetSbPath(ReportName, StartDate)
		//splist := strings.Split(localDir, "/")
		//length := len(splist)
		//reportDate = splist[length-1]
	} else {
		localDir = fmt.Sprintf("%s/%s/%s", path, ReportName, StartDate)
		//reportDate = StartDate
	}

	num := common.GetDirFileNum(localDir)
	if num > 0 {
		//t.ReportTaskRepository.Add(localDir, ReportName, ReportType, reportDate)
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

func (t *ToolReportService) addSbNotice(reportTYpe string, day int) (bool, string) {
	today := time.Now()
	starttime := today.Add(-1 * time.Duration(day) * time.Hour * 24)
	//获取过去一年的数据，暂定365天
	daylist := tool.GetBetweenWeek(starttime.Format(vars.TimeFormatTpl), today.Format(vars.TimeFormatTpl), vars.TimeFormatTpl)
	reason := "report_type: sb;"
	errdayList := make([]string, 0)
	errMonthList := make([]string, 0)
	dataflag := true
	for _, dateStr := range daylist {
		re, _ := t.AddNotice(reportTYpe, vars.BrandMetricsWeekly, dateStr.StartDate)
		if re == false {
			errdayList = append(errdayList, dateStr.StartDate+"~"+dateStr.EndDate)
		}
	}

	startMonthTime := today.Add(-1 * time.Duration(day) * time.Hour * 24)
	monthlist := tool.GetBetweenmonth(startMonthTime.Format(vars.TimeFormatTpl), today.Format(vars.TimeFormatTpl), vars.TimeFormatTpl)
	for _, monthStr := range monthlist {
		re, _ := t.AddNotice(reportTYpe, vars.BrandMetricsMonthly, monthStr.StartDate)
		if re == false {
			errMonthList = append(errMonthList, monthStr.StartDate+"~"+monthStr.EndDate)
		}
	}

	if len(errdayList) != 0 {
		dataflag = false
		reason = reason + " report_name:" + vars.BrandMetricsWeekly + " date:" + strings.Join(errdayList, ",")
	}
	if len(errMonthList) != 0 {
		dataflag = false
		reason = reason + "report_name:" + vars.BrandMetricsMonthly + " date:" + strings.Join(errMonthList, ",")
	}

	if dataflag == true {
		return dataflag, ""
	}
	return dataflag, reason
}

// 运行那种一个时间周期的报表类型
func (t *ToolReportService) runSbrand(sps []*bean.ProfileToken, day int, wg *reportsystem.Pool) {
	today := time.Now()
	starttime := today.Add(-1 * time.Duration(day) * time.Hour * 24)
	daylist := tool.GetBetweenWeek(starttime.Format(vars.TimeFormatTpl), today.Format(vars.TimeFormatTpl), vars.TimeFormatTpl)
	for _, dateStr := range daylist {
		wg.Add(1)
		go func(date bean.Mydate, mywg *reportsystem.Pool, profileList []*bean.ProfileToken) {
			for _, profile := range profileList {
				ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s_%d", vars.SB, vars.BrandMetricsWeekly, profile.ProfileId, time.Now().Unix()))
				all_steps.SbAllSteps(profile, vars.BrandMetricsWeekly, date.StartDate, date.EndDate, wg, ctx)
			}
		}(dateStr, wg, sps)
	}

	startMonthTime := today.Add(-1 * time.Duration(day) * time.Hour * 24)
	monthlist := tool.GetBetweenmonth(startMonthTime.Format(vars.TimeFormatTpl), today.Format(vars.TimeFormatTpl), vars.TimeFormatTpl)
	for _, monthStr := range monthlist {
		wg.Add(1)
		go func(date bean.Mydate, mywg *reportsystem.Pool, profileList []*bean.ProfileToken) {
			for _, profile := range profileList {
				ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s_%d", vars.SB, vars.BrandMetricsMonthly, profile.ProfileId, time.Now().Unix()))
				all_steps.SbAllSteps(profile, vars.BrandMetricsMonthly, date.StartDate, date.EndDate, wg, ctx)
			}
		}(monthStr, wg, sps)
	}
}

func (t *ToolReportService) runPpc(msg *amqp.Delivery) error {
	report, err := t.receiveCommon(msg)
	if err != nil {
		logger.Logger.Error(err.Error(), "==================get queue data error ")
		return err
	}

	// 避免在测试和pre环境拉数据
	if g.Cfg().GetString("server.Env") != developModel {
		t.Notice(report.ProcessId, "推送成功", 2, time.Now().Unix(), report.CallBackUrl)
		return nil
	}
	tokenList, err := t.getPpcProfileToken(report.Profileids)
	if err != nil || (tokenList == nil || len(tokenList) == 0) {
		logger.Logger.Info("error1======", err.Error())
		return err
	}

	daylist, err := tool.GetDays(report.StartDate, report.EndDate, vars.TimeLayout)
	if err != nil {
		logger.Logger.Info("get daylist error", err.Error())
		return err
	}

	AllReport := make([]string, 0)
	//sbnewbrandNum := 60
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
		for _, item := range tokenList {
			for _, reportName := range report.ReportName {
				switch report.ReportType {
				case vars.SP:
					pool.Add(1)
					go func(pitem *bean.ProfileToken, preportName string, pdayitem string) {
						ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s_%d", vars.SP, preportName, pitem.ProfileId, time.Now().Unix()))
						all_steps.SPAllSteps(pitem, preportName, pdayitem, pool, ctx)
					}(item, reportName, dayitem)

				case vars.SD:
					for _, tactic := range []string{"T00020", "T00030"} {
						pool.Add(1)
						go func(pitem *bean.ProfileToken, preportName string, pdayitem string, ptactic string) {
							ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s_%d", vars.SD, preportName, pitem.ProfileId, time.Now().Unix()))
							all_steps.SDAllSteps(pitem, preportName, pdayitem, ptactic, pool, ctx)
						}(item, reportName, dayitem, tactic)
					}
				case vars.SB:
					pool.Add(1)
					go func(pitem *bean.ProfileToken, preportName string, pdayitem string) {
						ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s_%d", vars.SB, preportName, pitem.ProfileId, time.Now().Unix()))
						result, errstruct := all_steps.SbAllSteps(pitem, preportName, pdayitem, pdayitem, pool, ctx)
						if result {
							fmt.Println("success sb")
						} else {
							fmt.Println(errstruct.ErrorReason, "=====error type")
						}
					}(item, reportName, dayitem)

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

	//dataflag, str := t.addSbNotice("sb", sbnewbrandNum)
	//if dataflag == false {
	//	sbbrandreason = reason + str
	//	//根据时间段拉取数据，有时候是空的拉取不到，所以检测到空目录不报错，通知下游还是成功的
	//	//status = 3
	//}

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

func (t *ToolReportService) getDspProfileToken(profileIdList []string) ([]*bean.DspRegionProfile, error) {
	return t.ProfileRepository.ListDspRegionProfile(profileIdList)
}

func (t *ToolReportService) receiveCommon(msg *amqp.Delivery) (*bean.ReportData, error) {
	var paramas bean.ReceiveMsgBody
	var report *bean.ReportData
	if err := json.Unmarshal(msg.Body, &paramas); err != nil {
		logger.Logger.Error(err)
		return nil, err
	}

	json.Unmarshal([]byte(paramas.MessageBody), &report)
	return report, nil
}

func (t *ToolReportService) ReceiveSp(msg *amqp.Delivery) error {
	return t.runPpc(msg)
}

func (t *ToolReportService) ReceiveDsp(msg *amqp.Delivery) error {
	logger.Logger.Info("get queue dsp:", string(msg.Body))
	report, err := t.receiveCommon(msg)
	if err != nil {
		logger.Logger.Error(err.Error(), "==================get queue data error ")
		return err
	}

	tokenList, err := t.getDspProfileToken(report.Profileids)
	if err != nil {
		logger.Logger.Error("error1======", err.Error())
		return err
	}

	if report.ReportDataType == 0 {
		report.ReportName = vars.ReportList[report.ReportType]
	}

	reportListStr := strings.Join(report.ReportName, ",")
	id, err := t.ReportBatchDetailRepository.Addone(report.Batch, reportListStr, report.ReportType, report.StartDate, report.EndDate, vars.TimeLayout)
	daylist, err := tool.GetDays(report.StartDate, report.EndDate, vars.TimeLayout)
	wg := reportsystem.NewPool(15)
	for _, dayitem := range daylist {
		for _, item := range tokenList {
			for _, reportName := range report.ReportName {
				wg.Add(1)
				go func(pt *bean.DspRegionProfile, ReportName string, StartDate string, fg *reportsystem.Pool) {
					ctx := logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s_%s_%s_%d", vars.DSP, ReportName, pt.ProfileId, time.Now().Unix()))
					dsp_all_steps.DspAllSteps(pt.Region, pt.ProfileId, ReportName, StartDate, fg, ctx)
				}(item, reportName, dayitem, wg)
			}
		}
	}

	wg.Wait()
	var reason = fmt.Sprintf("report_type: %s ", report.ReportType)
	var temp = 0
	var total = 0
	var status = 2
	var str = ""
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

	t.ReportBatchDetailRepository.ChangeStatus(status, id, reason)
	t.ReportBatchRepository.ChangeStatus(status, report.Batch)
	t.Notice(report.ProcessId, reason, status, time.Now().Unix(), report.CallBackUrl)
	return nil
}

func (t *ToolReportService) ReceiveSb(msg *amqp.Delivery) error {
	return t.runPpc(msg)
}

func (t *ToolReportService) ReceiveSd(msg *amqp.Delivery) error {
	return t.runPpc(msg)
}

func (t *ToolReportService) ReceiveLog(msg *amqp.Delivery) error {
	logger.Logger.Info("get queue:", string(msg.Body))
	return nil
}
