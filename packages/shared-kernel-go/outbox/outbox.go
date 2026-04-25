package outbox

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusPublished Status = "published"
	StatusFailed    Status = "failed"
)

type OutboxRow struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	AggregateType string          `json:"aggregate_type" db:"aggregate_type"`
	AggregateID   string          `json:"aggregate_id" db:"aggregate_id"`
	EventType     string          `json:"event_type" db:"event_type"`
	Payload       json.RawMessage `json:"payload" db:"payload"`
	TargetQueue   string          `json:"target_queue" db:"target_queue"`
	RequestID     uuid.UUID       `json:"request_id" db:"request_id"`
	Status        Status          `json:"status" db:"status"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	PublishedAt   *time.Time      `json:"published_at,omitempty" db:"published_at"`
}

func NewOutboxRow(aggregateType, aggregateID, eventType, targetQueue string, payload json.RawMessage, requestID uuid.UUID) OutboxRow {
	return OutboxRow{
		ID:            uuid.New(),
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		EventType:     eventType,
		Payload:       payload,
		TargetQueue:   targetQueue,
		RequestID:     requestID,
		Status:        StatusPending,
		CreatedAt:     time.Now().UTC(),
	}
}
