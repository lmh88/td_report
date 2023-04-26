package service

import (
	"github.com/gogf/gf/frame/g"
	amqp "github.com/rabbitmq/amqp091-go"
)

type TestService struct {

}

func NewTestService()*TestService{
	return &TestService{}
}

func(s *TestService)ReceiveKeyword(msg  *amqp.Delivery)error {
    //fmt.Println("keyword:", string(msg.Body))
	g.Log().Println("keyword:", string(msg.Body))
	return nil
}

func(s *TestService)ReceiveTarget(msg  *amqp.Delivery) error {
	//fmt.Println("target:",string(msg.Body))
	g.Log().Println("target:", string(msg.Body))
	return nil
}

func(s *TestService)ReceiveLog(msg  *amqp.Delivery) error{
	g.Log().Println("log:",string(msg.Body))
	return nil
}
