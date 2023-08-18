package services

import (
	"context"
	"encoding/json"
	"fmt"
	"fold/internal/models"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func SQS(doc *models.DenormalizedProject) error {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error loading AWS config:", err)
		return err
	}

	// Create an SQS client
	client := sqs.NewFromConfig(cfg)

	queueURL := "https://sqs.ap-south-1.amazonaws.com/182059941907/syncElasticsearch.fifo"

	jsonBytes, err := json.Marshal(doc)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	// Create a new message metadata
	messageGroupId := "sync-elastic"
	messageDeduplicationId := GenerateUniqueID()

	// Send message to the SQS FIFO queue
	sendMessageInput := &sqs.SendMessageInput{
		QueueUrl:               aws.String(queueURL),
		MessageBody:            aws.String(string(jsonBytes)),
		MessageGroupId:         aws.String(messageGroupId),
		MessageDeduplicationId: aws.String(messageDeduplicationId),
	}
	_, err = client.SendMessage(context.TODO(), sendMessageInput)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return err
	}

	fmt.Println("Message sent to SQS successfully!")
	return nil
}

func GenerateUniqueID() string {
	id := uuid.New()
	return id.String()
}
