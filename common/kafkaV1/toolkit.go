package kafkaV1

import (
	"fmt"
	"github.com/Shopify/sarama"
	"strings"
	"sync"
)

var (
	wg sync.WaitGroup
)

type KafkaTool interface {
	MyAsyncProducer() (sarama.AsyncProducer, error)
	MyConsumer(processMessage func(msg []byte) error) error
}

type kafkaTool struct {
	Addrs  []string
	Config *sarama.Config
	Topic  string
}

func NewKafkaTool(broker string, config *sarama.Config, topic string) KafkaTool {
	addrs := strings.Split(broker, ",")
	return &kafkaTool{Addrs: addrs, Config: config, Topic: topic}
}

func DefaultConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V2_4_0_0
	return config
}

func (k *kafkaTool) MyAsyncProducer() (sarama.AsyncProducer, error) {
	producer, err := sarama.NewAsyncProducer(k.Addrs, k.Config)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}

	return producer, nil
}

func (k *kafkaTool) SetTopic(topic string) {
	k.Topic = topic
}

func (k *kafkaTool) MyConsumer(processMessage func(msg []byte) error) error {
	consumer, err := sarama.NewConsumer(k.Addrs, k.Config)
	if err != nil {
		return err
	}
	defer consumer.Close()
	//设置分区
	partitionList, err := consumer.Partitions(k.Topic)
	if err != nil {
		fmt.Println("Failed to get the list of partitions: ", err)
		return nil
	}

	fmt.Println(partitionList)
	//循环分区
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(k.Topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("Failed to start consumer for partition %d: %s\n", partition, err)
			return err
		}
		defer pc.AsyncClose()
		wg.Add(1)
		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				processMessage(msg.Value)
			}

		}(pc)
	}
	//time.Sleep(time.Hour)
	wg.Wait()
	consumer.Close()
	return err
}
