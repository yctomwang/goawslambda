package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"time"
)

// Request is the structure of the lambda input
type Request struct {
	CodeSnippet string `json:"codeSnippet"`
}

func HandleRequest(ctx context.Context, req Request) (string, error) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-2")},
	)

	// Create a SQS service client.
	svc := sqs.New(sess)

	// URL of the SQS queue
	queueURL := "https://sqs.ap-southeast-2.amazonaws.com/989900959400/golangleetcode.fifo"

	// Serialize the request object to JSON
	reqJSON, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("error marshalling request: %v", err)
	}
	messageGroupId := "Group-" + time.Now().Format("20060102150405")
	// Send message to SQS queue
	_, err = svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds:   aws.Int64(10),
		MessageBody:    aws.String(string(reqJSON)),
		QueueUrl:       &queueURL,
		MessageGroupId: aws.String(messageGroupId), // Use the unique MessageGroupId
	})

	if err != nil {
		return "", fmt.Errorf("error sending SQS message: %v", err)
	}

	return "Message sent to SQS queue successfully", nil
}

func main() {
	lambda.Start(HandleRequest)
}
