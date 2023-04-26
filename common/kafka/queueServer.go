package kafka

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"github.com/gogf/guuid"
	"math/rand"
	"sync"
	"time"
)

// EventServerIface 现阶段针对kafka 只有topic消息
type EventServerIface interface {
	// Publish 发布消息
	Publish(topic string, messageBody string) error

	// GetConsumer Subscribe 订阅消息
	//Subscribe(topic, subscriptionName, consumer string, handler interface{}) error

	// GetConsumer 获取消费者
	GetConsumer(topic, consumerName string) (*EventConsumer, error)
	// AddQueueConsumer 增加消费者
	AddQueueConsumer(topic, consumerName ,groupId string  , handler interface{}) error
	// CreateQueue 创建队列
	CreateQueue(topic string) error
}

type EventServer struct {
	config      Config
	producerMap sync.Map
	consumerMap sync.Map

	once    sync.Once
	running bool
}

func NewEventServer(config Config) (*EventServer, error){
	server := &EventServer{config: config}
	return server, nil
}

func (s *EventServer) CreateQueue(topic string) error {
	return nil
}

func (s *EventServer) Publish(topic string, messageBody string) error {
	if topic == "" {
		return errors.New("topic 不能为空")
	}
	if messageBody == "" {
		return errors.New("body 不能为空")
	}

	var (
		producer *Producer
		topicmsg TopicMessageRequest
		err error
	)

	queueKey := fmt.Sprintf("%s-%s", s.config.BusinessName, topic)
	if p, ok := s.producerMap.Load(queueKey); ok {
		producer = p.(*Producer)
	} else {

		producer, err = NewProducer(s.config)
		if err != nil {
			return  err
		}

		s.producerMap.Store(topic, producer)
	}

	topicmsg.Topic = topic
	topicmsg.Body.MessageBody = messageBody
	topicmsg.Body.MessageId = guuid.New().String()
	return producer.SendMsg(topicmsg)
}


func (s *EventServer) GetConsumer(topic, consumerName string) (*EventConsumer, error) {
	consumerKey := fmt.Sprintf("%s-%s", topic, consumerName)
	if consumer, ok := s.consumerMap.Load(consumerKey); ok {
		return consumer.(*EventConsumer), nil
	}

	return nil, nil
}

// AddQueueConsumer 队列增加消费者
func (s *EventServer) AddQueueConsumer(topic, consumerName string, groupId string, handler interface{}) error {
	consumer, err := s.addConsumer(topic, consumerName,groupId, handler)
	if err != nil {
		return err
	}

	if s.running {
		consumer.Start()
	}

	return nil
}

// AddConsumer 添加消费者
func (s *EventServer) addConsumer(topic, consumerName string, groupId string, handler interface{}) (*EventConsumer, error) {
	if handler == nil {
		return nil, errors.New("处理方法不能为空")
	}

	consumerKey := fmt.Sprintf("%s-%s", topic, consumerName)
	if _, ok := s.consumerMap.Load(consumerKey); ok {
		return nil, fmt.Errorf("消费者: %s 队列名称: %s, 消费者的消费行为必须保持一致！", consumerName, topic)
	}

    // 优化kafka参数，避免队列rebance
	config := cluster.NewConfig()
	config.Consumer.Offsets.Initial             = sarama.OffsetNewest
	config.Consumer.Fetch.Default               = 262144 //1024*512 // 这个参数需要观察一下，服务器用524288 测试环境用
	config.Consumer.MaxProcessingTime           = 1000 * time.Millisecond
	config.Consumer.MaxWaitTime                 = 1500 * time.Millisecond
	config.Consumer.Group.Session.Timeout       = 20  * time.Second
	config.Consumer.Group.Heartbeat.Interval    = 6  * time.Second
	topics:= []string{topic}

	client, err := cluster.NewConsumer(s.config.addres, groupId, topics, config)

	if err != nil {
		fmt.Printf("Failed to start consumer: %s\n", err)
		return nil, err
	}

	consumer := NewEventConsumer(client, topic, consumerName, handler)
	s.consumerMap.Store(consumerKey, consumer)
	return consumer, nil
}

// Start 开始事件处理
func (s *EventServer) Start() error {
	s.running = true
	s.once.Do(func() {
		s.consumerMap.Range(func(key, value interface{}) bool {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			consumer := value.(*EventConsumer)
			consumer.Start()
			return true
		})
	})

	return nil
}

