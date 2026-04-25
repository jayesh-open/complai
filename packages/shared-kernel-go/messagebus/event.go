package messagebus

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        uuid.UUID       `json:"id"`
	Type      string          `json:"type"`
	Source    string          `json:"source"`
	TenantID  uuid.UUID       `json:"tenant_id"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp time.Time       `json:"timestamp"`
	RequestID uuid.UUID       `json:"request_id"`
}
