package consumer_server

import (
	"encoding/json"
	"fmt"
	"td_report/app/bean"
	"td_report/app/repo"
	"td_report/app/service/report_v2"
	"td_report/app/service/report_v2/varible"
	"td_report/boot"
	"td_report/pkg/logger"
	"td_report/vars"
	"time"
)

type ErrServer struct {
	Mq *report_v2.MqServer
}

func NewErrServer() *ErrServer {
	instance := ErrServer{
		Mq: report_v2.NewMqServer(),
	}
	return &instance
}

func (s *ErrServer) ConsumeFailQueue() {
	for _, item := range s.GetAllFailQueue() {
		if s.Mq.GetQueueLen(item) > 0 {
			logger.Logger.Infof("消费队列名称：%s", item)
			s.Mq.ReceiveMsg(item, s.ErrHandle)
		} else {
			fmt.Println(item + " wait time 10s")
			time.Sleep(time.Second * 10)
		}
	}
}

func (s *ErrServer) GetAllFailQueue() []string {
	allClient := repo.NewSellerClientRepository().GetAll()
	res := make([]string, 0)
	for _, item := range varible.FailQueueMap {
		for _, c := range allClient {
			res = append(res, fmt.Sprintf("c%d:%s", c, item))
		}
	}
	return res
}

func (s *ErrServer) ErrHandle(msg []byte) error {
	var (
		reportIdMsg varible.ReportIdMsg
		err error
	)

	err = json.Unmarshal(msg, &reportIdMsg)
	if err != nil {
		return err
	}

	errInfo := &bean.ReportErr{
		ReportType: reportIdMsg.ReportType,
		ReportName: reportIdMsg.ReportName,
		ReportDate: reportIdMsg.ReportDate,
		ProfileId:  reportIdMsg.ProfileId,
		ErrorType:  repo.ReportErrorTypeOne,
		KeyParam:   reportIdMsg.TraceId,
		Extra:      reportIdMsg.ReportTactic,
	}

	if reportIdMsg.FirstTryCount == 0 {
		errInfo.ErrorType = repo.ReportErrorTypeRetry
		errInfo.ErrorReason = "429尝试次数过多"
	}

	tab := true
	for _, item := range reportIdMsg.RetryCount {
		if item != 0 {
			tab = false
		}
	}

	if tab {
		errInfo.ErrorType = repo.ReportErrorTypeTimeOut
		errInfo.ErrorReason = "等待超时未完成"
	}

	err = boot.RabitmqClient.SendWithNoClose(vars.ErrorInfo, errInfo)
	return err
}


