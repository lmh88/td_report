package consumer_server

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"td_report/app"
	"td_report/app/bean"
	"td_report/app/service/common"
	"td_report/app/service/report_v2"
	"td_report/app/service/report_v2/create_report"
	"td_report/app/service/report_v2/one_step"
	"td_report/app/service/report_v2/varible"
	"td_report/boot"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

type ConsumerOneServer struct {
	ReportTaskService *common.ReportTaskService
}

func NewConsumerOneServer() *ConsumerOneServer {
	reportTaskService := app.InitializeReportTaskService()
	return &ConsumerOneServer{
		ReportTaskService: reportTaskService,
	}
}

var currentBatchName string

//var oneStepMap = map[string]func(*report_v2.ProfileMsg, string, string, context.Context) (string, *bean.ReportErr) {
//	vars.SP: one_step.SpOneStep,
//	vars.SB: one_step.SbOneStep,
//}

// ConsumerOne 第一步消费profile生成report_id
func (t *ConsumerOneServer) ConsumerOne(queueType, reportType, clientTag string) {
	if AllowRun(reportType, clientTag) {
		//检查队列是否繁忙
		go func() {
			for {
				if !AllowRun(reportType, clientTag) {
					panic("队列繁忙，暂停消费")
				}
				fmt.Println("allow run")
				time.Sleep(10 * time.Second)
			}
		}()

		mq := report_v2.NewMqServer()
		queueName := t.SelectQueue(mq, queueType, reportType, clientTag)
		if queueName != "" {
			fmt.Println(queueName)
			mq.ReceiveMsg(queueName, t.ConsumerOneHandle)
			return
		} else {
			fmt.Println("无消费数据")
		}

	}
	time.Sleep(30 * time.Second)
}

// SelectQueue 启动是选择一个消费队列，如果有指定队列且不为空，就消费指定队列，否则按快慢队列顺序选择
func (t *ConsumerOneServer) SelectQueue(mq *report_v2.MqServer, queueType, reportType, clientTag string) string {

	if queueType != "" {
		queueName, ok := varible.QueueMap[reportType][queueType]
		if ok {
			queue := varible.AddQueuePre(queueName, clientTag)
			if mq.GetQueueLen(queue) > 0 {
				return queue
			}
		}
	}

	for _, item := range varible.QueueLevelArray {
		queueName := varible.QueueMap[reportType][item]
		queue := varible.AddQueuePre(queueName, clientTag)
		if mq.GetQueueLen(queue) > 0 {
			return queue
		}
	}

	return ""
}

func (t *ConsumerOneServer) ConsumerOneHandle(msg []byte) error {
	var (
		profileMsg *varible.ProfileMsg
		wg         sync.WaitGroup
	)
	err := json.Unmarshal(msg, &profileMsg)
	if err != nil {
		return err
	}

	mq := report_v2.NewMqServer()

	//消费到新的批次时，退出，并重新选择队列
	if currentBatchName != "" {
		if currentBatchName != profileMsg.BatchKey {
			panic(currentBatchName + "批次消费结束")
		}
	}

	for _, reportName := range vars.ReportList[profileMsg.ReportType] {
		if profileMsg.ReportType == vars.SD {
			for _, tactic := range []string{"T00020", "T00030"} {
				wg.Add(1)
				go dealProfile(profileMsg, reportName, tactic, mq, &wg)
			}
		} else {
			wg.Add(1)
			go dealProfile(profileMsg, reportName, "", mq, &wg)
		}
	}
	wg.Wait()
	currentBatchName = profileMsg.BatchKey
	return nil
}

func dealProfile(profileMsg *varible.ProfileMsg, reportName, tactic string, mq *report_v2.MqServer, wg *sync.WaitGroup) {
	defer wg.Done()
	var (
		reportId  string
		reportErr *bean.ReportErr
	)
	traceId := varible.GetTraceId(profileMsg, reportName, tactic)
	ctx := logger.Logger.NewTraceIDContext(context.Background(), traceId)
	logger.Logger.InfoWithContext(ctx, "report one step start")
	if profileMsg.ReportType == vars.SD {
		reportId, reportErr = one_step.SdOneStep(profileMsg, reportName, tactic, ctx)
	} else {
		dealFunc := one_step.OneStepFuncMap[profileMsg.ReportType]
		reportId, reportErr = dealFunc(profileMsg, reportName, profileMsg.ReportDate, ctx)
	}

	reportIdMsg := varible.ReportIdMsg{
		ReportType:    profileMsg.ReportType,
		ProfileId:     profileMsg.ProfileId,
		ProfileType:     profileMsg.ProfileType,
		Region:        profileMsg.Region,
		RefreshToken:  profileMsg.RefreshToken,
		Timestamp:     profileMsg.Timestamp,
		ReportDate:    profileMsg.ReportDate,
		BatchKey:      profileMsg.BatchKey,
		ClientTag:     profileMsg.ClientTag,
		ClientId:      profileMsg.ClientId,
		ClientSecret:  profileMsg.ClientSecret,
		TraceId:       traceId,
		ReportId:      reportId,
		ReportName:    reportName,
		ReportTactic:  tactic,
		RetryCount:    varible.GetRetryCount(profileMsg.ReportType),
		FirstTryCount: varible.FirstRetryCountMap[profileMsg.ReportType],
	}
	sendMsg, _ := json.Marshal(reportIdMsg)

	if reportErr != nil {
		if reportErr.ErrorReason == varible.Retry429 {
			mq.SendMsgExp(varible.ReportDefaultExchange,
				varible.AddQueuePre(varible.RetryDelayMap[profileMsg.ReportType], profileMsg.ClientTag),
				sendMsg, varible.RetryDelayTimeMap[profileMsg.ReportType])
		} else {
			//错误记录
			boot.RabitmqClient.SendWithNoClose(vars.ErrorInfo, reportErr)
		}
		logger.Logger.ErrorWithContextMap(ctx, map[string]interface{}{
			"flag":    "报表第一步失败",
			"errInfo": reportErr,
		})
		return
	} else {
		logger.Logger.InfoWithContext(ctx, "报表第一步成功", "reportId:", reportId)
		err := mq.SendMsgExp(varible.ReportDefaultExchange,
			varible.AddQueuePre(varible.DelayQueueMap[profileMsg.ReportType][varible.WaitTime10s], profileMsg.ClientTag),
			sendMsg, varible.DelayTimeMap[varible.WaitTime10s])
		if err != nil {
			logger.Logger.ErrorWithContext(ctx, "reportId发送失败：", err.Error())
			return
		}
		create_report.BatchStart(&reportIdMsg)
	}
	return
}
