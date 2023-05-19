package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// snippet-end:[sqs.go.send_message.imports]

// GetQueueURL gets the URL of an Amazon SQS queue
// Inputs:
//
//	sess is the current session, which provides configuration for the SDK's service clients
//	queueName is the name of the queue
//
// Output:
//
//	If success, the URL of the queue and nil
//	Otherwise, an empty string and an error from the call to
func GetQueueURL(sess *session.Session, queue *string) (*sqs.GetQueueUrlOutput, error) {
	// Create an SQS service client
	svc := sqs.New(sess)

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queue,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SendMsg sends a message to an Amazon SQS queue
// Inputs:
//
//	sess is the current session, which provides configuration for the SDK's service clients
//	queueURL is the URL of the queue
//
// Output:
//
//	If success, nil
//	Otherwise, an error from the call to SendMessage
func Publisher(sess *session.Session, queueURL *string) (string, error) {
	// Create an SQS service client
	// snippet-start:[sqs.go.send_message.call]
	svc := sqs.New(sess)

	sqsResponse, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Title": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("testing publisher"),
			},
		},
		MessageBody: aws.String("Hello from the auth service"),
		QueueUrl:    queueURL,
	})
	// snippet-end:[sqs.go.send_message.call]
	if err != nil {
		return "", err
	}

	return *sqsResponse.MessageId, nil
}

func SendMsg() {
	// snippet-start:[sqs.go.send_message.args]
	queue := os.Getenv("QUEUE_NAME")
	aws_region := os.Getenv("AWS_REGION")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(aws_region),
	})
	if err != nil {
		panic(err)
	}

	// snippet-end:[sqs.go.send_message.sess]

	// Get URL of queue
	result, err := GetQueueURL(sess, &queue)
	if err != nil {
		fmt.Println("Got an error getting the queue URL:")
		fmt.Println(err)
		return
	}

	queueURL := result.QueueUrl

	msgID, err := Publisher(sess, queueURL)
	if err != nil {
		fmt.Println("Got an error sending the message:")
		fmt.Println(err)
		return
	}

	log.Println("Message sent successfully to queue with id: ", msgID)
}
