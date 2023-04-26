package rabbitmq

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"td_report/common/sendmsg"
	"td_report/pkg/logger"
	"time"
)

type RabbitMq struct {
	conn *amqp.Connection
}

func NewRabbitmq(url string) (*RabbitMq, error) {
	myconn, err := amqp.Dial(url)
	defer func() {
		if r := recover(); r != nil {
			sendObj:=sendmsg.New(sendmsg.Wechat)
			sendObj.SendMsg("rabbitmq not connect !!!")
		}
	}()
	if err != nil {
		logger.Infof("rabbitmq err:%s", err.Error())
		panic(err)
	}

	return &RabbitMq{
		conn: myconn,
	}, nil
}

func (r *RabbitMq) Close() bool {
	if !r.conn.IsClosed() {
		r.conn.Close()
	}

	return true
}

func (r *RabbitMq) SendWithNoClose(topic string, playload interface{}) error {
	ch, err := r.conn.Channel()
	queue, err := ch.QueueDeclare(topic, true, false, false, false, nil)
	if err != nil {
		logger.Infof("Failed to declare a queue err:%s", err.Error())
		return err
	}
	
	defer func() {
		if !ch.IsClosed() {
			ch.Close()
		}
	}()

	data, err := json.Marshal(playload)
	if err != nil {
		return err
	}

	err = ch.Publish("", queue.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	})

	if err != nil {
		logger.Infof("publish err:%s", err.Error())
		return err
	}

	return err
}

func (r *RabbitMq) Send(topic string, playload interface{}) error {
	ch, err := r.conn.Channel()
	if err != nil {
		logger.Infof("Failed to open a channel err:%s", err.Error())
		return err
	}

	defer func() {
		if !ch.IsClosed() {
			ch.Close()
		}
	}()

	queue, err := ch.QueueDeclare(topic, true, false, false, false, nil)
	if err != nil {
		logger.Infof("Failed to declare a queue err:%s", err.Error())
		return err
	}
	data, err := json.Marshal(playload)
	if err != nil {
		return err
	}

	err = ch.Publish("", queue.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	})
	if err != nil {
		logger.Infof("publish err:%s", err.Error())
		return err
	}

	return err
}

func (r *RabbitMq) GetLength(queueName string) int {
	ch, err := r.conn.Channel()
	if err != nil {
		logger.Infof("Failed to open a channel err:%s", err.Error())
		return -1
	}

	queue, err := ch.QueueDeclarePassive(queueName, true, false, false, false, nil)
	if err != nil {
		logger.Infof("Failed to  count :%s", err.Error())
		return -1
	}

	return queue.Messages
}

// ReceivWithAck 目前都是单线程接收
func (r *RabbitMq) ReceivWithAck(topic string, f func([]byte) error) {
	ch, err := r.conn.Channel()
	if err != nil {
		logger.Infof("Failed to open a channel err:%s", err.Error())
		return
	}
	// 限制一下每次rabbitmq拿取数据的数量
	ch.Qos(5, 0, false)
	queue, err := ch.QueueDeclare(topic, true, false, false, false, nil)
	if err != nil {
		logger.Logger.Error("ueue declare err:%s", err.Error())
		return
	}
	msgCh, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		logger.Logger.Errorf("Failed to receive message err:%s", err.Error())
		return
	}

	forever := make(chan bool)
	go func() {
		for d := range msgCh {
			err = f(d.Body)
			if err != nil {
				logger.Logger.Infof("receive message data:error", err.Error())
				continue
			} else {

				//logger.Logger.Infof("receive message data:%s", string(d.Body))
				d.Ack(false)
			}
		}
	}()

	// 检测rabbitmq是否正常，如果连接失败就通知退出
	go func() {
		var (
			length int
			n      = 0
		)

		for {

			// -1 是channel连接失败特定标记的长度，直接关闭整个进程
			length = r.GetLength(topic)
			if length == -1 || (length == 0 && n > 20) {
				r.Close()
				forever <- true
			}
			n++
			time.Sleep(10 * time.Second)
		}
	}()

	logger.Logger.Infof(" Waiting for messages. topic:%s", topic)
	<-forever
}

func (r *RabbitMq) Receive(topic string, f func([]byte) error) {
	ch, err := r.conn.Channel()
	if err != nil {
		logger.Infof("Failed to open a channel err:%s", err.Error())
		return
	}
	// 限制一下每次rabbitmq拿取数据的数量
	ch.Qos(5, 0, false)
	queue, err := ch.QueueDeclare(topic, true, false, false, false, nil)
	if err != nil {
		logger.Logger.Infof("queue declare err:%s", err.Error())
		return
	}
	msgCh, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		logger.Logger.Infof("Failed to receive message err:%s", err.Error())
		return
	}

	forever := make(chan bool)
	go func() {
		for d := range msgCh {
			err = f(d.Body)
			if err != nil {
				continue
			}
		}
	}()

	// 检测rabbitmq是否正常，如果连接失败就通知退出
	go func() {
		var (
			length int
			n      = 0
		)

		for {

			length = r.GetLength(topic)
			// -1 是channel连接失败特定标记的长度，直接关闭整个进程
			if length == -1 || (length == 0 && n > 20) {
				r.Close()
				forever <- true
			}
			n++

			time.Sleep(10 * time.Second)
		}
	}()

	logger.Logger.Infof(" Waiting for messages. topic:%s", topic)
	<-forever
}
