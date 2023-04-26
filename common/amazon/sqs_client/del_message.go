package sqs_client

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"td_report/pkg/logger"
)

// SQSDeleteMessageAPI defines the interface for the GetQueueUrl and DeleteMessage functions.
// We use this interface to test the functions using a mocked service.
type SQSDeleteMessageAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	DeleteMessage(ctx context.Context,
		params *sqs.DeleteMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}

// GetQueueURL gets the URL of an Amazon SQS queue.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a GetQueueUrlOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to GetQueueUrl.
//func GetQueueURL(c context.Context, api SQSDeleteMessageAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
//	return api.GetQueueUrl(c, input)
//}

// RemoveMessage deletes a message from an Amazon SQS queue.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a DeleteMessageOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to DeleteMessage.
func RemoveMessage(c context.Context, api SQSDeleteMessageAPI, input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return api.DeleteMessage(c, input)
}

func DelMessage(ctx context.Context, client *sqs.Client, topic *string, messageHandle *string) error {
	qUInput := &sqs.GetQueueUrlInput{
		QueueName: topic,
	}

	// Get URL of queue
	result, err := GetQueueURL(context.TODO(), client, qUInput)
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, "Got an error getting the queue URL:", err)
		return err
	}

	queueURL := result.QueueUrl

	dMInput := &sqs.DeleteMessageInput{
		QueueUrl:      queueURL,
		ReceiptHandle: messageHandle,
	}

	_, err = RemoveMessage(context.TODO(), client, dMInput)
	if err != nil {
		logger.Logger.ErrorWithContext(ctx, "Got an error deleting the message:", err)
		return err
	}

	//fmt.Println("Deleted message from queue with URL " + *queueURL)
	return nil
}
