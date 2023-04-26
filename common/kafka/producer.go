package kafka

import (
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"math/rand"
)

type Producer  struct {
	Config Config
	client sarama.SyncProducer
}

func NewProducer(queurConfigonfig Config)(*Producer, error) {
	if len(queurConfigonfig.addres) == 0 {
		return nil, errors.New("缺少参数: 地址")
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner  = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Version = sarama.V0_10_0_1

	client, err := sarama.NewSyncProducer(queurConfigonfig.addres, config)
	if err != nil {
		g.Log().Error("producer close, err:", err.Error())
		return nil,err
	}

	return &Producer{
		Config: queurConfigonfig,
		client: client,
	}, nil
}

func(c *Producer)SendMsg(messageBody TopicMessageRequest) error {
	var (
		body []byte
		err error
	)
	body,_ = json.Marshal(messageBody.Body)

	msg := &sarama.ProducerMessage{}
	num:= rand.Intn(10000000)
	skey:= gconv.String(num)
	msg.Key = sarama.StringEncoder(skey)
	msg.Topic = messageBody.Topic
	msg.Value = sarama.StringEncoder(body)

	if _, _, err = c.client.SendMessage(msg);err!= nil {
		g.Log().Error("send message failed,", err.Error())
		return err
	}
	return nil
}



