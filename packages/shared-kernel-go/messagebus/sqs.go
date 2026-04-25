package messagebus

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type SQSBus struct {
	client   *sqs.Client
	logger   zerolog.Logger
	handlers map[string]Handler
	mu       sync.RWMutex
}

func NewSQSBus(cfg aws.Config, logger zerolog.Logger) *SQSBus {
	opts := func(o *sqs.Options) {
		if endpoint := os.Getenv("AWS_ENDPOINT_URL"); endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	}

	return &SQSBus{
		client:   sqs.NewFromConfig(cfg, opts),
		logger:   logger.With().Str("component", "sqs-bus").Logger(),
		handlers: make(map[string]Handler),
	}
}

func (b *SQSBus) Publish(ctx context.Context, queueURL string, event Event) error {
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("sqs: marshal event: %w", err)
	}

	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(body)),
		MessageAttributes: map[string]sqstypes.MessageAttributeValue{
			"EventType": {
				DataType:    aws.String("String"),
				StringValue: aws.String(event.Type),
			},
			"TenantID": {
				DataType:    aws.String("String"),
				StringValue: aws.String(event.TenantID.String()),
			},
			"RequestID": {
				DataType:    aws.String("String"),
				StringValue: aws.String(event.RequestID.String()),
			},
		},
	}

	_, err = b.client.SendMessage(ctx, input)
	if err != nil {
		return fmt.Errorf("sqs: send message to %s: %w", queueURL, err)
	}

	b.logger.Debug().
		Str("queue", queueURL).
		Str("event_type", event.Type).
		Str("event_id", event.ID.String()).
		Msg("published message to SQS")

	return nil
}

func (b *SQSBus) Subscribe(queueURL string, handler Handler) error {
	b.mu.Lock()
	b.handlers[queueURL] = handler
	b.mu.Unlock()

	go b.poll(queueURL, handler)
	return nil
}

func (b *SQSBus) poll(queueURL string, handler Handler) {
	for {
		ctx := context.Background()
		output, err := b.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:              aws.String(queueURL),
			MaxNumberOfMessages:   10,
			WaitTimeSeconds:       20,
			MessageAttributeNames: []string{"All"},
		})
		if err != nil {
			b.logger.Error().Err(err).Str("queue", queueURL).Msg("sqs receive error")
			time.Sleep(5 * time.Second)
			continue
		}

		for _, msg := range output.Messages {
			var event Event
			if err := json.Unmarshal([]byte(aws.ToString(msg.Body)), &event); err != nil {
				b.logger.Error().Err(err).Msg("sqs: unmarshal event")
				continue
			}

			if err := handler(ctx, event); err != nil {
				b.logger.Error().Err(err).
					Str("event_id", event.ID.String()).
					Str("event_type", event.Type).
					Msg("sqs: handler error")
				continue
			}

			_, _ = b.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueURL),
				ReceiptHandle: msg.ReceiptHandle,
			})
		}
	}
}
