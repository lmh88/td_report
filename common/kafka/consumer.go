package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"github.com/gogf/gf/frame/g"
	"reflect"
)

type EventConsumer struct {
	topic        string
	consumerName string
	client       *cluster.Consumer
	handler      interface{}
	quit         chan bool
	running      bool
}

func NewEventConsumer(client *cluster.Consumer, topic string, consumerName string, handler interface{}) *EventConsumer {
	return &EventConsumer{
		topic:        topic,
		consumerName: consumerName,
		client:       client,
		handler:      handler,
		quit:         make(chan bool),
		running:      false,
	}
}

func (c *EventConsumer) Start() {
	if c.running {
		return
	}

	c.running = true

	go func() {

		for {
			select {
			case <-c.quit:
				return
			default:

				waitChan := make(chan int)
				for msg:=range c.client.Messages() {
					go c.handleMessages(waitChan, msg)
					<-waitChan
				}
			}
		}
	}()
}

func (c *EventConsumer) handleMessages(waitChan chan int, msg *sarama.ConsumerMessage) {
	defer func() {
		//if r := recover(); r != nil {
		//	g.Log().Errorf("%s recover from panic, msg: %v", c.getKey(), r)
		//}
		waitChan <- 1
	}()
			// 消息转换
			rMethod := reflect.ValueOf(c.handler)
			rParamType := rMethod.Type().In(0)

			var paramVal reflect.Value
			if reflect.TypeOf(msg) == rParamType {
				paramVal = reflect.ValueOf(msg)
			} else if reflect.TypeOf(&msg) == rParamType {
				paramVal = reflect.ValueOf(&msg)
			} else { // 自定义消息体
				body := reflect.ValueOf(msg).Elem().Interface()
				paramVal = reflect.ValueOf(body)
			}

			// 消费消息
			ret := rMethod.Call([]reflect.Value{paramVal})
			if ret[0].Interface() != nil {
				err := ret[0].Interface().(error)
				g.Log().Errorf("%s, 执行消费消息失败: %s", c.getKey(), err.Error())
				return
			}
}

func (c *EventConsumer) Stop() {
	if !c.running {
		return
	}

	c.quit <- true
	c.running = false
	g.Log().Infof("结束轮询消息: " + c.getKey())
}

func (c *EventConsumer) Restart() {
	go func() {
		c.quit <- true
		c.running = false
		c.Start()
	}()
}

func (c *EventConsumer) getKey() string {
	return c.topic
}
