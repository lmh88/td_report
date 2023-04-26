package startup

import (
	"log"
	"sync"
	"td_report/app"
	"td_report/common/rabbitmq"
	"td_report/pkg/logger"
	"td_report/vars"
)

// 存储注册到map的队列处理参数
var queuetoolList sync.Map

// InitWebMq  消息队列instance
func InitWebMq() {
	logger.Init("queue_consumer", false)
	//注册处理消息队列
	RigsterQueueToolList()

	go func() {
		err := rabbitmq.NewMqServer()

		if err != nil {
			log.Fatalf("queue init err: %v", err)
			return
		}

		reportToolService := app.InitializeToolReportService()
		queuetoolList.Range(func(k, v interface{}) bool {
			msgType := k.(string)
			queueName := v.(string)
			consume := queueName + msgType

			if queueName == vars.QueueMessageSp {
				err = rabbitmq.MqServer.AddQueueConsumer(queueName, consume, reportToolService.ReceiveSp)
			} else if queueName == vars.QueueMessageSd {
				err = rabbitmq.MqServer.AddQueueConsumer(queueName, consume, reportToolService.ReceiveSd)
			} else if queueName == vars.QueueMessageSb {
				err = rabbitmq.MqServer.AddQueueConsumer(queueName, consume, reportToolService.ReceiveSb)
			} else if queueName == vars.QueueMessageDsp {
				err = rabbitmq.MqServer.AddQueueConsumer(queueName, consume, reportToolService.ReceiveDsp)
			} else if queueName == vars.QueueMessageLog {
				err = rabbitmq.MqServer.AddQueueConsumer(queueName, consume, reportToolService.ReceiveLog)
			}

			return true
		})
	}()
}

// RigsterQueueToolList 由于队列是单机部署，只是一个分组添加多个消费者
func RigsterQueueToolList() {
	// 权限或者工具拉取数据的队列
	queuetoolList.Store(vars.QueueKeySp, vars.QueueMessageSp)
	queuetoolList.Store(vars.QueueKeySd, vars.QueueMessageSd)
	queuetoolList.Store(vars.QueueKeySb, vars.QueueMessageSb)
	queuetoolList.Store(vars.QueueKeyDsp, vars.QueueMessageDsp)
	// 公共日志队列
	queuetoolList.Store(vars.QueueKeyLog, vars.QueueMessageLog)
}
