package sqs_client

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gogf/gf/frame/g"
	"td_report/pkg/logger"
)

type SqsClient struct {
	Client *sqs.Client
}

func NewSqsClient() (*SqsClient, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(g.Cfg().GetString("sqs.regin")),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(g.Cfg().GetString("sqs.key"), g.Cfg().GetString("sqs.secret"), ""),
		))

	if err != nil {
		logger.Logger.Error("sqs_client get config error:", err.Error())
		return nil, err
	}

	client := sqs.NewFromConfig(cfg)
	return &SqsClient{
		Client: client,
	}, nil
}

// Receive 接收消息
func (t *SqsClient) Receive(ctx context.Context, topic *string, receiveNum int32, handelFunc func(data *sqs.ReceiveMessageOutput) error) (error, bool) {
	return Receive(ctx, t.Client, topic, receiveNum, handelFunc)
}

func (t *SqsClient)Receivelp(ctx context.Context, topic *string, receiveNum int32, handelFunc func(data *sqs.ReceiveMessageOutput) error)(error, bool){
	return ReceiveLp(ctx, t.Client, topic, receiveNum, handelFunc)
}

// ListQueue 打印出来，暂时没有返回值
func (t *SqsClient) ListQueue() {
	ListQueue(t.Client)
}

// DelMessage 删除消息
func (t *SqsClient) DelMessage(ctx context.Context, topic *string, messageHandle *string) error {
	return DelMessage(ctx, t.Client, topic, messageHandle)
}
