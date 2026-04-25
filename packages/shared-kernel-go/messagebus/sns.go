package messagebus

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	snstypes "github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type SNSPublisher struct {
	client *sns.Client
	logger zerolog.Logger
}

func NewSNSPublisher(cfg aws.Config, logger zerolog.Logger) *SNSPublisher {
	opts := func(o *sns.Options) {
		if endpoint := os.Getenv("AWS_ENDPOINT_URL"); endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	}

	return &SNSPublisher{
		client: sns.NewFromConfig(cfg, opts),
		logger: logger.With().Str("component", "sns-publisher").Logger(),
	}
}

func (p *SNSPublisher) Publish(ctx context.Context, topicARN string, event Event) error {
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("sns: marshal event: %w", err)
	}

	input := &sns.PublishInput{
		TopicArn: aws.String(topicARN),
		Message:  aws.String(string(body)),
		MessageAttributes: map[string]snstypes.MessageAttributeValue{
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

	output, err := p.client.Publish(ctx, input)
	if err != nil {
		return fmt.Errorf("sns: publish to %s: %w", topicARN, err)
	}

	p.logger.Debug().
		Str("topic", topicARN).
		Str("event_type", event.Type).
		Str("message_id", aws.ToString(output.MessageId)).
		Msg("published message to SNS")

	return nil
}
