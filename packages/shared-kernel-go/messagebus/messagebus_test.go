package messagebus_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/complai/complai/packages/shared-kernel-go/messagebus"
)

func TestEvent_JSONRoundTrip(t *testing.T) {
	original := messagebus.Event{
		ID:        uuid.New(),
		Type:      "invoice.created",
		Source:    "invoice-service",
		TenantID:  uuid.New(),
		Payload:   json.RawMessage(`{"amount":1000}`),
		Timestamp: time.Now().UTC().Truncate(time.Second),
		RequestID: uuid.New(),
	}

	data, err := json.Marshal(original)
	assert.NoError(t, err)

	var decoded messagebus.Event
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, original.ID, decoded.ID)
	assert.Equal(t, original.Type, decoded.Type)
	assert.Equal(t, original.Source, decoded.Source)
	assert.Equal(t, original.TenantID, decoded.TenantID)
	assert.JSONEq(t, string(original.Payload), string(decoded.Payload))
	assert.Equal(t, original.Timestamp.Unix(), decoded.Timestamp.Unix())
	assert.Equal(t, original.RequestID, decoded.RequestID)
}

func TestMessageBusInterface_Compliance(t *testing.T) {
	// Verify SQSBus implements MessageBus at compile time
	var _ messagebus.MessageBus = (*messagebus.SQSBus)(nil)
}
