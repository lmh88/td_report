package one_step

import (
	"context"
	"encoding/json"
	"td_report/app/bean"
	"td_report/app/service/report_v2"
	"td_report/app/service/report_v2/create_report"
	"td_report/app/service/report_v2/varible"
	"td_report/boot"
	"td_report/pkg/logger"
	"td_report/vars"
)

func RetryConsumer(reportType, clientTag string) {
	mq := report_v2.NewMqServer()
	queueName := varible.AddQueuePre(varible.RetryQueueMap[reportType], clientTag)
	mq.ReceiveMsg(queueName, retryHandle)
}

var OneStepFuncMap = map[string]func (*varible.ProfileMsg, string, string, context.Context) (string, *bean.ReportErr) {
	vars.SP: SpOneStep,
	vars.SB: SbOneStep,
}

func retryHandle(msg []byte) error {
	var (
		profileMsg *varible.ProfileMsg
		reportIdMsg *varible.ReportIdMsg
		err error
		reportId string
		errInfo *bean.ReportErr
	)

	err = json.Unmarshal(msg, &reportIdMsg)
	if err != nil {
		return err
	}

	profileMsg = &varible.ProfileMsg{
		ReportType:   reportIdMsg.ReportType,
		ProfileId:    reportIdMsg.ProfileId,
		Region:       reportIdMsg.Region,
		RefreshToken: reportIdMsg.RefreshToken,
		Timestamp:    reportIdMsg.Timestamp,
		ReportDate:   reportIdMsg.ReportDate,
		BatchKey:     reportIdMsg.BatchKey,
		ClientTag:    reportIdMsg.ClientTag,
		ClientId:     reportIdMsg.ClientId,
		ClientSecret: reportIdMsg.ClientSecret,
	}

	ctx := logger.Logger.NewTraceIDContext(context.Background(), reportIdMsg.TraceId)
	logger.Logger.InfoWithContext(ctx, "重试开始")
	if reportIdMsg.ReportType == vars.SD {
		reportId, errInfo = SdOneStep(profileMsg, reportIdMsg.ReportName, reportIdMsg.ReportTactic, ctx)
	} else {
		dealFunc := OneStepFuncMap[reportIdMsg.ReportType]
		reportId, errInfo = dealFunc(profileMsg, reportIdMsg.ReportName, reportIdMsg.ReportDate, ctx)
	}

	mq := report_v2.NewMqServer()
	if errInfo != nil {
		dealErr(reportIdMsg, mq, errInfo, ctx)
		logger.Logger.ErrorWithContext(ctx, "重试失败", reportIdMsg)
	} else {
		reportIdMsg.ReportId = reportId
		sendMsg, _ := json.Marshal(reportIdMsg)
		err = mq.SendMsgExp(varible.ReportDefaultExchange,
			varible.AddQueuePre(varible.DelayQueueMap[profileMsg.ReportType][varible.WaitTime10s], reportIdMsg.ClientTag),
			sendMsg, varible.DelayTimeMap[varible.WaitTime10s])
		if err != nil {
			logger.Logger.ErrorWithContext(ctx, "reportId发送失败：", err.Error())
			return err
		}
		create_report.BatchStart(reportIdMsg)
		logger.Logger.InfoWithContext(ctx, "重试成功", reportIdMsg)
	}
	return nil
}

func dealErr(reportIdMsg *varible.ReportIdMsg, mq *report_v2.MqServer, errInfo *bean.ReportErr, ctx context.Context) {
	if errInfo != nil {
		if errInfo.ErrorReason == varible.Retry429 {
			reportIdMsg.FirstTryCount--
			msg, _ := json.Marshal(reportIdMsg)
			if reportIdMsg.FirstTryCount == 0 {
				mq.SendMsg(varible.ReportDefaultExchange,
					varible.AddQueuePre(varible.FailQueueMap[reportIdMsg.ReportType], reportIdMsg.ClientTag), msg)
				boot.RabitmqClient.SendWithNoClose(vars.ErrorInfo, errInfo)
			} else {
				mq.SendMsgExp(varible.ReportDefaultExchange,
					varible.AddQueuePre(varible.RetryDelayMap[reportIdMsg.ReportType], reportIdMsg.ClientTag),
					msg, varible.RetryDelayTimeMap[reportIdMsg.ReportType])
			}
		} else {
			boot.RabitmqClient.SendWithNoClose(vars.ErrorInfo, errInfo)
			logger.Logger.ErrorWithContext(ctx, errInfo)
		}
	}
}
