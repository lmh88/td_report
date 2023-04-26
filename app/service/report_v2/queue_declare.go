package report_v2

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"td_report/app/repo"
	"td_report/app/service/report_v2/varible"
)

func DeclareQueue(reportType, clientTag string) {

	mq := NewMqServer()

	if mq.ExistQueue(varible.AddQueuePre(varible.FailQueueMap[reportType], clientTag)) {
		return
	}

	//交换机声明
	err := mq.ch.ExchangeDeclare(varible.ReportDefaultExchange, "direct", true, false, false, false, nil)
	if err != nil {
		fmt.Println("exchange err:", err.Error())
	} else {
		fmt.Println(varible.ReportDefaultExchange, "success")
	}

	//report id 队列
	reportQueue := varible.AddQueuePre(varible.ReportQueueMap[reportType], clientTag)
	_, err = mq.ch.QueueDeclare(reportQueue, true, false, false, false, nil)
	if err != nil {
		fmt.Println("queue err:", err.Error())
	} else {
		fmt.Println(reportQueue, "declare success")
	}
	//绑定
	err = mq.ch.QueueBind(reportQueue, reportQueue, varible.ReportDefaultExchange, false, nil)
	if err != nil {
		fmt.Println("queue binding err:", err.Error())
	} else {
		fmt.Println(reportQueue, "binding success")
	}

	//profile 队列声明和绑定
	for _, name := range varible.QueueMap[reportType] {
		name = varible.AddQueuePre(name, clientTag)
		_, err = mq.ch.QueueDeclare(name, true, false, false, false, nil)
		if err != nil {
			fmt.Println("queue err:", err.Error())
		} else {
			fmt.Println(name, "declare success")
		}

		err = mq.ch.QueueBind(name, name, varible.ReportDefaultExchange, false, nil)
		if err != nil {
			fmt.Println("queue binding err:", err.Error())
		} else {
			fmt.Println(name, "binding success")
		}
	}

	//延时队列声明
	arg := amqp.Table{"x-dead-letter-exchange": varible.ReportDefaultExchange, "x-dead-letter-routing-key": reportQueue}
	for _, name := range varible.DelayQueueMap[reportType] {
		name = varible.AddQueuePre(name, clientTag)
		_, err = mq.ch.QueueDeclare(name, true, false, false, false, arg)
		if err != nil {
			fmt.Println("queue err:", err.Error())
		} else {
			fmt.Println(name, "declare success")
		}

		err = mq.ch.QueueBind(name, name, varible.ReportDefaultExchange, false, nil)
		if err != nil {
			fmt.Println("queue binding err:", err.Error())
		} else {
			fmt.Println(name, "binding success")
		}
	}

	//失败队列声明
	failQueue := varible.AddQueuePre(varible.FailQueueMap[reportType], clientTag)
	_, err = mq.ch.QueueDeclare(failQueue, true, false, false, false, nil)
	if err != nil {
		fmt.Println("queue err:", err.Error())
	} else {
		fmt.Println(failQueue, "declare success")
	}
	//绑定
	err = mq.ch.QueueBind(failQueue, failQueue, varible.ReportDefaultExchange, false, nil)
	if err != nil {
		fmt.Println("queue binding err:", err.Error())
	} else {
		fmt.Println(failQueue, "binding success")
	}

	//重试队列声明
	retryQueue := varible.AddQueuePre(varible.RetryQueueMap[reportType], clientTag)
	_, err = mq.ch.QueueDeclare(retryQueue, true, false, false, false, nil)
	if err != nil {
		fmt.Println("queue err:", err.Error())
	} else {
		fmt.Println(retryQueue, "declare success")
	}
	//绑定
	err = mq.ch.QueueBind(retryQueue, retryQueue, varible.ReportDefaultExchange, false, nil)
	if err != nil {
		fmt.Println("queue binding err:", err.Error())
	} else {
		fmt.Println(retryQueue, "binding success")
	}

	//重试延时队列
	arg = amqp.Table{"x-dead-letter-exchange": varible.ReportDefaultExchange, "x-dead-letter-routing-key": retryQueue}
	retryDelayQueue := varible.AddQueuePre(varible.RetryDelayMap[reportType], clientTag)
	_, err = mq.ch.QueueDeclare(retryDelayQueue, true, false, false, false, arg)
	if err != nil {
		fmt.Println("queue err:", err.Error())
	} else {
		fmt.Println(retryDelayQueue, "declare success")
	}
	//绑定
	err = mq.ch.QueueBind(retryDelayQueue, retryDelayQueue, varible.ReportDefaultExchange, false, nil)
	if err != nil {
		fmt.Println("queue binding err:", err.Error())
	} else {
		fmt.Println(retryDelayQueue, "binding success")
	}
}


func ScanClient(reportType, clientTag string) {
	if reportType != "" && clientTag != "" {
		DeclareQueue(reportType, clientTag)
		return
	}

	ids := repo.NewSellerClientRepository().GetAll()
	for _, id := range ids {
		clientTag = varible.GetClientTag(id)
		for reportType, _ = range varible.QueueMap {
			DeclareQueue(reportType, clientTag)
			//fmt.Println(clientTag+reportType)
		}
	}
	return
}

