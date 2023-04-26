package sqs_client

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"td_report/pkg/logger"
)

// SQSGetLPMsgAPI defines the interface for the GetQueueUrl and ReceiveMessage functions.
// We use this interface to test the functions using a mocked service.
type SQSGetLPMsgAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	ReceiveMessage(ctx context.Context,
		params *sqs.ReceiveMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
}

// GetQueueURLLp gets the URL of an Amazon SQS queue.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a GetQueueUrlOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to GetQueueUrl.
func GetQueueURLLp(c context.Context, api SQSGetLPMsgAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}

// GetLPMessages gets the messages from an Amazon SQS long polling queue.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a ReceiveMessageOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to ReceiveMessage.
func GetLPMessages(c context.Context, api SQSGetLPMsgAPI, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return api.ReceiveMessage(c, input)
}



func ReceiveLp(ctx context.Context, client *sqs.Client, topic *string, receiveNum int32, handelFunc func(data *sqs.ReceiveMessageOutput) error) (error, bool) {
	timeout := 18
	gQInput := &sqs.GetQueueUrlInput{
		QueueName: topic,
	}

	result, err := GetQueueURLLp(context.TODO(), client, gQInput)
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, "sqs_client get Got an error getting the queue URL:", err.Error())
		return err, false
	}

	mInput := &sqs.ReceiveMessageInput{
		QueueUrl:  result.QueueUrl,
		AttributeNames: []types.QueueAttributeName{
			"SentTimestamp",
		},
		MaxNumberOfMessages: receiveNum,
		MessageAttributeNames: []string{
			"All",
		},
		WaitTimeSeconds: int32(timeout),
	}

	resp, err := GetLPMessages(context.TODO(), client, mInput)
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, "sqs_client get Got an error receiving messages:", err.Error())
		return err, false
	}

    err = handelFunc(resp)
	if err != nil {
		return err, false
	}
	return nil, true
}
