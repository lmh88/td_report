package consumer_server

import (
	"fmt"
	rediS "github.com/go-redis/redis"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"td_report/app/service/report_v2"
	"td_report/app/service/report_v2/varible"
	"td_report/boot"
	"td_report/common/sendmsg/wechart"
	"td_report/pkg/logger"
	"time"
)

type CheckBusyServer struct {
	Rds *rediS.Client
	Mq *report_v2.MqServer
	Notice string
}

func NewCheckBusyServer() *CheckBusyServer {
	s := CheckBusyServer{
		Rds: boot.RedisCommonClient.GetClient(),
		Mq: report_v2.NewMqServer(),
		Notice: env,
	}
	return &s
}

var env = g.Cfg().GetString("server.Env", "local") + "\n"

func (s *CheckBusyServer) CheckQueueBusy() {

	for reportType, limitLen := range varible.LimitRetryQueueMap {
		for clientTag, _ := range varible.ClientMap {
			s.dealBusy(reportType, clientTag, "retry", limitLen)
		}
	}

	for reportType, limitLen := range varible.LimitReportQueueMap {
		for clientTag, _ := range varible.ClientMap {
			s.dealBusy(reportType, clientTag, "report", limitLen)
		}
	}

	if s.Notice != env {
		key := "0ac12699-c90a-4fd8-9bf6-5fb75cba6459"
		send := wechart.NewSendMsg(key, true)
		send.Send("", s.Notice)
	}
}

var busyNotice = map[string]string{
	"report": "report_id too much",
	"retry": "429 too much",
}

var checkFunc = map[string]func(string, string, int) bool {
	"report": NewCheckBusyServer().CheckReportQueue,
	"retry": NewCheckBusyServer().CheckRetryQueue,
}

func (s *CheckBusyServer) dealBusy(reportType, clientTag, kind string, limitLen int) {
	key := varible.GetQueueBusyKey(reportType, clientTag, kind)
	fn := checkFunc[kind]
	if !fn(reportType, clientTag, limitLen) {
		logger.Logger.Info(key)
		s.Notice += fmt.Sprintf("%s:%s:%s:%s\n", kind, clientTag, reportType, busyNotice[kind])
		s.Rds.Set(key, gtime.Timestamp(), time.Minute * 5)
	}
}

func (s *CheckBusyServer) CheckRetryQueue(reportType, clientTag string, limitLen int) bool {
	queueName := varible.AddQueuePre(varible.RetryQueueMap[reportType], clientTag)
	if s.Mq.GetQueueLen(queueName) > limitLen {
		return false
	}
	return true
}

func (s *CheckBusyServer) CheckReportQueue(reportType, clientTag string, limitLen int) bool {
	queueName := varible.AddQueuePre(varible.ReportQueueMap[reportType], clientTag)
	if s.Mq.GetQueueLen(queueName) > limitLen {
		return false
	}
	return true
}

func AllowRun(reportType, clientTag string) bool {
	key1 := varible.GetQueueBusyKey(reportType, clientTag, "retry")
	key2 := varible.GetQueueBusyKey(reportType, clientTag, "report")
	rds := boot.RedisCommonClient.GetClient()
	if rds.Exists(key1).Val() == 1 || rds.Exists(key2).Val() == 1  {
		return false
	}
	return true
}
