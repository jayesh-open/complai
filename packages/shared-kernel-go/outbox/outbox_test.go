package outbox_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/packages/shared-kernel-go/outbox"
)

func TestNewOutboxRow(t *testing.T) {
	requestID := uuid.New()
	payload := json.RawMessage(`{"invoice_id":"123"}`)

	row := outbox.NewOutboxRow("invoice", "inv-001", "invoice.created", "invoice-events", payload, requestID)

	assert.NotEqual(t, uuid.Nil, row.ID)
	assert.Equal(t, "invoice", row.AggregateType)
	assert.Equal(t, "inv-001", row.AggregateID)
	assert.Equal(t, "invoice.created", row.EventType)
	assert.Equal(t, "invoice-events", row.TargetQueue)
	assert.Equal(t, requestID, row.RequestID)
	assert.Equal(t, outbox.StatusPending, row.Status)
	assert.Nil(t, row.PublishedAt)
	assert.WithinDuration(t, time.Now().UTC(), row.CreatedAt, 2*time.Second)

	require.NotNil(t, row.Payload)
	assert.JSONEq(t, `{"invoice_id":"123"}`, string(row.Payload))
}

func TestOutboxRow_StatusConstants(t *testing.T) {
	assert.Equal(t, outbox.Status("pending"), outbox.StatusPending)
	assert.Equal(t, outbox.Status("published"), outbox.StatusPublished)
	assert.Equal(t, outbox.Status("failed"), outbox.StatusFailed)
}
