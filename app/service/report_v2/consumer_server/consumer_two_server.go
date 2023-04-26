package consumer_server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"strconv"
	"sync"
	"td_report/app"
	"td_report/app/bean"
	"td_report/app/repo"
	"td_report/app/service/common"
	"td_report/app/service/report_v2"
	"td_report/app/service/report_v2/create_report"
	"td_report/app/service/report_v2/varible"
	"td_report/boot"
	"td_report/pkg/logger"
	"td_report/pkg/save_file"
	"td_report/vars"
)

type ConsumerTwoServer struct {
	ReportTaskService *common.ReportTaskService
	ReportSchedule    *repo.ReportSchduleRepository
	rw                sync.RWMutex
}

func NewConsumerTwoServer() *ConsumerTwoServer {
	return &ConsumerTwoServer{
		ReportTaskService: app.InitializeReportTaskService(),
		ReportSchedule:    repo.NewReportSchduleRepository(),
	}
}

func (s *ConsumerTwoServer) ConsumerTwo(reportType, clientTag string) {
	num := g.Cfg().GetString("consumer_quantity." + reportType, "10")
	i, _ := strconv.Atoi(num)
	mq := report_v2.NewMqServer()
	queue := varible.AddQueuePre(varible.ReportQueueMap[reportType], clientTag)
	mq.MulReceiveMsg(queue, s.ConsumerTwoHandle, i)
}

func (s *ConsumerTwoServer) ConsumerTwoHandle(msg []byte) error {
	var reportIdMsg *varible.ReportIdMsg
	err := json.Unmarshal(msg, &reportIdMsg)
	if err != nil {
		return err
	}
	ctx := logger.Logger.NewTraceIDContext(context.Background(), reportIdMsg.TraceId)
	logger.Logger.InfoWithContext(ctx, "获取报告开始处理")
	doing, errInfo := s.TwoStep(reportIdMsg, ctx)

	failQueue := s.RetryHandle(reportIdMsg, doing)
	// 判断是否完成
	if !doing || failQueue != "" {
		create_report.BatchEnd(reportIdMsg)
	}

	if err != nil {
		logger.Logger.ErrorWithContext(ctx, "处理失败：", err.Error())
		return err
	}
	if errInfo != nil {
		logger.Logger.ErrorWithContext(ctx, "获取报告处理失败")
		boot.RabitmqClient.SendWithNoClose(vars.ErrorInfo, errInfo)
		return fmt.Errorf(errInfo.ErrorReason)
	} else {
		logger.Logger.InfoWithContext(ctx, "获取报告处理成功")
		return nil
	}
}

func (s *ConsumerTwoServer) RetryHandle(reportIdMsg *varible.ReportIdMsg, doing bool) string {
	var(
		delayKey              string
		delayQueue, failQueue string
	)

	for _, item := range varible.WaitTimeArray {
		if reportIdMsg.RetryCount[item] > 0 {
			reportIdMsg.RetryCount[item]--
			break
		}
	}

	for _, item := range varible.WaitTimeArray {
		if reportIdMsg.RetryCount[item] > 0 {
			delayKey = item
			break
		}
	}

	if doing {
		if delayKey == "" {
			failQueue = varible.AddQueuePre(varible.FailQueueMap[reportIdMsg.ReportType], reportIdMsg.ClientTag)
		} else {
			delayQueue = varible.AddQueuePre(varible.DelayQueueMap[reportIdMsg.ReportType][delayKey], reportIdMsg.ClientTag)
		}
	}

	msg, _ := json.Marshal(reportIdMsg)
	mq := report_v2.NewMqServer()
	if delayQueue != "" {
		mq.SendMsgExp(varible.ReportDefaultExchange, delayQueue, msg, varible.DelayTimeMap[delayKey])
	}

	if failQueue != "" {
		mq.SendMsg(varible.ReportDefaultExchange, failQueue, msg)
	}
	return failQueue
}

func (s *ConsumerTwoServer) TwoStep(reportIdMsg *varible.ReportIdMsg, ctx context.Context) (doing bool, errInfo *bean.ReportErr) {

	errInfo = &bean.ReportErr{
		ReportType: reportIdMsg.ReportType,
		ReportName: reportIdMsg.ReportName,
		ReportDate: reportIdMsg.ReportDate,
		ProfileId:  reportIdMsg.ProfileId,
		KeyParam:   logger.Logger.FromTraceIDContext(ctx),
		Extra:      reportIdMsg.ReportTactic,
	}

	//time.Sleep(time.Millisecond * 100)
	//rand.Seed(time.Now().UnixNano())
	//if rand.Intn(2) > 0 {
	//	return true, nil
	//}
	//time.Sleep(time.Second)
	//time.Sleep(time.Millisecond * 100)
	//return false, nil

	resp, err := create_report.GetReportStatus(reportIdMsg, ctx)
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
			"flag": "获取报表地址错误",
			"err":  err.Error(),
		})
		errInfo.ErrorType = repo.ReportErrorTypeTwo
		errInfo.ErrorReason = err.Error()
		return false, errInfo
	}

	if resp.Status == "IN_PROGRESS" {
		return true, nil
	}

	if resp.Status == "FAILURE" {
		logger.Logger.ErrorWithContext(ctx, "获取报表地址失败", resp)
		errInfo.ErrorType = repo.ReportErrorTypeTwo
		errInfo.ErrorReason = "获取报表地址失败"
		return false, errInfo
	}

	if resp.Status == "SUCCESS" {
		val, err := create_report.DownloadReport(reportIdMsg, resp.Location, ctx)
		if err != nil {
			logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
				"flag":        "下载报告错误",
				"err":         err.Error(),
				"reportIdMsg": reportIdMsg,
			})
			errInfo.ErrorType = repo.ReportErrorTypeThree
			errInfo.ErrorReason = err.Error()
			return false, errInfo
		}

		//写入文件
		if reportIdMsg.ReportType == vars.SD {
			err = save_file.SaveSDFile(reportIdMsg.ReportName, reportIdMsg.ReportDate, reportIdMsg.ProfileId, reportIdMsg.ReportTactic, val)
		} else {
			dealFunc := SaveFileMap[reportIdMsg.ReportType]
			err = dealFunc(reportIdMsg.ReportName, reportIdMsg.ReportDate, reportIdMsg.ProfileId, val)
		}

		if err != nil {
			logger.Logger.ErrorWithContext(ctx, map[string]interface{}{
				"flag":        "文件保存失败",
				"err":         err.Error(),
				"reportIdMsg": reportIdMsg,
			})
			errInfo.ErrorType = repo.ReportErrorTypeOne
			errInfo.ErrorReason = "SaveFile error:" + err.Error()
			return false, errInfo
		} else {
			logger.Logger.InfoWithContext(ctx, "SaveFile success")

			return false, nil
		}
	}

	return false, nil
}

var SaveFileMap = map[string]func(string, string, string, []byte) error{
	vars.SP: save_file.SaveSPFile,
	vars.SB: save_file.SaveSBFile,
}
