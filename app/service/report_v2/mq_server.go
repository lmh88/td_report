package report_v2

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
	"td_report/common/stringx"
	"td_report/pkg/logger"
	"time"
)

type MqServer struct {
	conn *amqp.Connection
	ch *amqp.Channel
}

var once sync.Once
var mqs *MqServer

func NewMqServer() *MqServer {
	once.Do(func () {
		uri := g.Cfg().GetString("rabbitmq.address")
		conn, err := amqp.Dial(uri)
		if err != nil {
			panic(err)
		}

		ch, err := conn.Channel()
		if err != nil {
			panic(err)
		}
		ch.Qos(1, 0, false)

		mqs = &MqServer{
			conn: conn,
			ch: ch,
		}
	})

	return mqs
}

// Reconnect 重连
func (mq *MqServer) Reconnect() *MqServer {
	uri := g.Cfg().GetString("rabbitmq.address")
	conn, err := amqp.Dial(uri)
	if err != nil {
		logger.Logger.Error("重连失败=>", err.Error())
		return mq
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Logger.Error("重连失败=>", err.Error())
		return mq
	}
	ch.Qos(1, 0, false)

	mq.conn = conn
	mq.ch = ch
	return mq
}

// SendMsg 发送消息
func (mq *MqServer) SendMsg(exchange, route string, msg []byte) error {
	return mq.ch.Publish(exchange, route, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body: msg,
	})
}

func (mq *MqServer) SendMsgExp(exchange, route string, msg []byte, expire int) error {
	return mq.ch.Publish(exchange, route, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body: msg,
		Expiration: fmt.Sprintf("%d", expire * 1000),
	})
}

// ReceiveMsg 接收消息
func (mq *MqServer) ReceiveMsg(queue string, handler func([]byte) error) {

	msg, err := mq.ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		logger.Errorf("Failed to receive message err:%s", err.Error())
		return
	}
	quit := make(chan bool)
	doing := make(chan bool, 1)
	go func(doing chan bool) {
		for d := range msg {
			doing <- true
			err := handler(d.Body)
			if err != nil {
				logger.Error("消息处理失败:", err.Error(), string(d.Body))
			}
			d.Ack(false)
			<- doing
		}
	}(doing)

	go func(quit chan bool) {
		n := 0
		for {
			qLen := mq.GetQueueLen(queue)
			if qLen == 0 && len(doing) == 0 {
				n++
				fmt.Println("num:", n)
				if n == 20 {
					quit <- true
				}
			} else {
				n = 0
			}
			fmt.Println(queue, "队列检查:", qLen)
			time.Sleep(time.Second * 5)
		}
	}(quit)
	<-quit
}

func (mq *MqServer) MulReceiveMsg(queue string, handler func([]byte) error, quantity int) {
	doing := make(chan bool, quantity)
	for i := 0; i < quantity; i++ {
		go func(i int) {
			consumer := fmt.Sprintf("%s_consumer_%d_%s", queue, i, stringx.Randn(4))
			msg, err := mq.ch.Consume(queue, consumer, false, false, false, false, nil)
			if err != nil {
				log.Println(consumer, err.Error())
				logger.Errorf("Failed to receive message err:%s", err.Error())
				err2 := mq.ch.Cancel(consumer, true)
				if err2 != nil {
					log.Println("err2", err2.Error())
					return
				}
				msg, err = mq.ch.Consume(queue, consumer, false, false, false, false, nil)
				if err != nil {
					log.Println(consumer, err.Error())
					return
				}
				//return
			}
			fmt.Println(consumer)
			for d := range msg {
				//log.Println(consumer)
				doing <- true
				err := handler(d.Body)
				if err != nil {
					logger.Error("消息处理失败:", err.Error(), string(d.Body))
				}
				d.Ack(false)
				<- doing
			}
		}(i)
	}

	//挂起
	forever := make(chan bool)
	//等待5分钟无数据退出
	go func () {
		checkNum := 0
		for {
			qLen := mq.GetQueueLen(queue)
			fmt.Println(queue, "队列检查长度：", qLen, "正在处理数：", len(doing))

			if qLen == 0 && len(doing) == 0  {
				if mq.CheckClosed() {
					forever <- false
				}
				checkNum++
				fmt.Println("num:", checkNum)
				if checkNum == 60 {
					mq.Shutdown()
					forever <- false
				}
			} else {
				checkNum = 0
			}

			time.Sleep(time.Second * 5)
		}
	}()
	<-forever
}

func (mq *MqServer) GetQueueLen(queueName string) int {
	queue, _ := mq.ch.QueueDeclarePassive(queueName, true, false, false, false, nil)
	return queue.Messages
}

func (mq *MqServer) ExistQueue(queueName string) bool {
	_, err := mq.ch.QueueDeclarePassive(queueName, true, false, false, false, nil)

	if err != nil {
		mq.ch, _ = mq.conn.Channel()
		return false
	}
	return true
}

func (mq *MqServer) Shutdown() {

	if !mq.ch.IsClosed() {
		mq.ch.Close()
	}

	if !mq.conn.IsClosed() {
		mq.conn.Close()
	}
}

func (mq *MqServer) CheckClosed() bool {
	if mq.ch.IsClosed() {
		return true
	}

	if mq.conn.IsClosed() {
		return true
	}

	return false
}
