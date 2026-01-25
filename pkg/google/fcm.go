package google

import (
	"context"
	"fmt"
	"log"
	"os"

	"firebase.google.com/go/v4/messaging"
)

type FCM struct {
	client *messaging.Client
	logger *log.Logger
}

type IFCM interface {
	GetLogger() *log.Logger
	GetClient() *messaging.Client
	SendNotification(token, title, body, imageUrl string, data map[string]string) error
	SendBatchNotification(tokens []string, title, body, image string, data map[string]string)
	SendBulkNotification(tokens []string, title, body, image string, data map[string]string) error
}

// NewFCM creates a new FCM client.
func NewFCM(serviceAccountPath string) (*FCM, error) {
	app, err := newFirebaseApp(serviceAccountPath)
	if err != nil {
		return nil, err
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, err
	}
	return &FCM{
		client: client,
		logger: log.New(os.Stdout, "[FCM-CLIENT]", log.LstdFlags|log.Lshortfile),
	}, nil
}

func (c *FCM) GetLogger() *log.Logger {
	return c.logger
}

func (c *FCM) GetClient() *messaging.Client {
	return c.client
}

// SendNotification sends a push notification to a specific token.
func (s *FCM) SendNotification(token, title, body, imageUrl string, data map[string]string) error {
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: imageUrl,
		},
	}

	// if data payload is provided
	if len(data) > 0 {
		message.Data = data
	}

	response, err := s.client.Send(context.Background(), message)
	if err != nil {
		return fmt.Errorf("failed to send notification: %v", err)
	}

	fmt.Printf("Successfully sent notification: %s\n", response)
	return nil
}

// SendBatchNotification sends a push notification to multiple tokens in batches.
func (s *FCM) SendBatchNotification(tokens []string, title, body, image string, data map[string]string) {
	const maxTokensPerBatch = 500 // FCM limit: max 500 tokens per batch

	// Split tokens into batches
	for i := 0; i < len(tokens); i += maxTokensPerBatch {
		end := i + maxTokensPerBatch
		if end > len(tokens) {
			end = len(tokens)
		}

		batchTokens := tokens[i:end]

		// Send notification for the current batch
		err := s.SendBulkNotification(batchTokens, title, body, image, data)
		if err != nil {
			s.logger.Printf("[BATCH] error sending batch notification (tokens %d-%d): %v", i, end-1, err)
		}
	}
}

// SendBulkNotification sends a push notification to multiple tokens.
func (s *FCM) SendBulkNotification(tokens []string, title, body, image string, data map[string]string) error {
	// Create the multicast message
	message := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: image, 
		},
	}

	// Add data payload if provided
	if len(data) > 0 {
		message.Data = data
	}

	// Send the multicast message using FCM
	response, err := s.client.SendEachForMulticast(context.Background(), message)
	if err != nil {
		return fmt.Errorf("[BULK] failed to send notification: %v", err)
	}

	// Log the results
	fmt.Printf("[BULK] successfully sent notifications: %d successful, %d failed\n", response.SuccessCount, response.FailureCount)

	// Handle failed tokens
	if response.FailureCount > 0 {
		for i, resp := range response.Responses {
			if !resp.Success {
				fmt.Printf("[BULK] failed to send notification to token[%d]: %v\n", i, resp.Error)
			}
		}
	}

	return nil
}
