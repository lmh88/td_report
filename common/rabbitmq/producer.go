package rabbitmq

import (
	"encoding/json"
	"errors"
	"github.com/gogf/gf/frame/g"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer  struct {
	Config Config
	chanel *amqp.Channel
}

func NewProducer(queurConfigonfig Config)(*Producer, error) {
	if len(queurConfigonfig.addres) == 0 {
		return nil, errors.New("缺少参数: 地址")
	}
	conn, err := amqp.Dial(queurConfigonfig.addres)
	ch, err := conn.Channel()
	if err != nil {
		g.Log().Error("producer close, err:", err.Error())
		return nil,err
	}

	return &Producer{
		Config: queurConfigonfig,
		chanel: ch,
	}, nil
}

func(c *Producer)SendMsg(messageBody TopicMessageRequest) error {
	var (
		body []byte
		err error
	)
	body,_ = json.Marshal(messageBody.Body)

	queue, err:= c.chanel.QueueDeclare(messageBody.Queuename, true, false, false, false, nil)
	if err != nil {

	}
	return c.chanel.Publish("", queue.Name, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:body,
			AppId: messageBody.Body.MessageId,
	})
}

