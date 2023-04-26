package sqs_client

import (
	"context"
	"errors"
	"td_report/pkg/logger"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// SQSReceiveMessageAPI defines the interface for the GetQueueUrl function.
// We use this interface to test the function using a mocked service.
type SQSReceiveMessageAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	ReceiveMessage(ctx context.Context,
		params *sqs.ReceiveMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
}

// GetQueueURL gets the URL of an Amazon SQS queue.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a GetQueueUrlOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to GetQueueUrl.
func GetQueueURL(c context.Context, api SQSReceiveMessageAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}

// GetMessages gets the most recent message from an Amazon SQS queue.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a ReceiveMessageOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to ReceiveMessage.
func GetMessages(c context.Context, api SQSReceiveMessageAPI, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return api.ReceiveMessage(c, input)
}

func Receive(ctx context.Context, client *sqs.Client, topic *string, receiveNum int32, handelFunc func(data *sqs.ReceiveMessageOutput) error) (error, bool) {
	timeout := 18
	gQInput := &sqs.GetQueueUrlInput{
		QueueName: topic,
	}

	urlResult, err := GetQueueURL(context.TODO(), client, gQInput)
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, "sqs_client get Got an error getting the queue URL:", err.Error())
		return err, false
	}

	gMInput := &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl:            urlResult.QueueUrl,
		MaxNumberOfMessages: receiveNum,
		VisibilityTimeout:   int32(timeout),
	}

	msgResult, err := GetMessages(context.TODO(), client, gMInput)
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, "sqs_client get Got an error receiving messages:", err.Error())
		return err, false
	}

	if msgResult.Messages != nil {
		handelFunc(msgResult)
	} else {
		logger.Logger.ErrorWithContext(ctx, "sqs_client No messages found:")
		return errors.New("no messages found"), false
	}

	return nil, true
}
