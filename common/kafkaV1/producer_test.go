package kafkaV1

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/gogf/gf/frame/g"
	"testing"
)

type AmzTopic struct {
	QueueId string
}

func TestKafkaToolProducer(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V0_11_0_2

	producer, err := sarama.NewAsyncProducer([]string{g.Cfg().GetString("kafka.url")}, config)
	if err != nil {
		fmt.Printf("producer_test create producer error :%s\n", err.Error())
		return
	}

	defer producer.AsyncClose()
	tp := &AmzTopic{
		QueueId: "13123123",
	}

	data, _ := json.Marshal(tp)
	msg := &sarama.ProducerMessage{
		Topic: Topic,
		Key:   nil,
	}

	for {
		msg.Value = sarama.ByteEncoder(data)

		// send to chain
		producer.Input() <- msg

		select {
		case suc := <-producer.Successes():
			fmt.Printf("offset: %d,  timestamp: %s\n", suc.Offset, suc.Timestamp.String())
		case fail := <-producer.Errors():
			fmt.Printf("err: %s\n", fail.Err.Error())
		}
	}
}

func TestProductByTool(t *testing.T) {
	kafkaTool := NewKafkaTool(g.Cfg().GetString("kafka.url"), DefaultConfig(), Topic)
	producer, err := kafkaTool.MyAsyncProducer()
	if err != nil {
		t.Log(err)
		return
	}
	msg := &sarama.ProducerMessage{
		Topic: Topic,
		Key:   nil,
	}

	for i := 0; i < 5; i++ {
		msg.Value = sarama.ByteEncoder([]byte("2222"))
		producer.Input() <- msg
		select {
		case suc := <-producer.Successes():
			fmt.Printf("offset: %d,  timestamp: %s\n", suc.Offset, suc.Timestamp.String())
		case fail := <-producer.Errors():
			fmt.Printf("err: %s\n", fail.Err.Error())
		}
	}

}
