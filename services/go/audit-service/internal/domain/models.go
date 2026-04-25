package domain

import (
	"time"

	"github.com/google/uuid"
)

type AuditEvent struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	UserID       *uuid.UUID `json:"user_id,omitempty"`
	ResourceType string     `json:"resource_type"`
	ResourceID   *uuid.UUID `json:"resource_id,omitempty"`
	Action       string     `json:"action"`
	OldValue     *string    `json:"old_value,omitempty"`
	NewValue     *string    `json:"new_value,omitempty"`
	Status       string     `json:"status"`
	ErrorMessage *string    `json:"error_message,omitempty"`
	IPAddress    *string    `json:"ip_address,omitempty"`
	UserAgent    *string    `json:"user_agent,omitempty"`
	TraceID      *string    `json:"trace_id,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type MerkleChain struct {
	ID           uuid.UUID `json:"id"`
	TenantID     uuid.UUID `json:"tenant_id"`
	HourBucket   time.Time `json:"hour_bucket"`
	EventCount   int       `json:"event_count"`
	HashPayload  string    `json:"hash_payload"`
	PreviousHash string    `json:"previous_hash"`
	ComputedHash string    `json:"computed_hash"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateAuditEventRequest struct {
	UserID       *uuid.UUID `json:"user_id"`
	ResourceType string     `json:"resource_type"`
	ResourceID   *uuid.UUID `json:"resource_id"`
	Action       string     `json:"action"`
	OldValue     *string    `json:"old_value"`
	NewValue     *string    `json:"new_value"`
	Status       string     `json:"status"`
	ErrorMessage *string    `json:"error_message"`
	IPAddress    *string    `json:"ip_address"`
	UserAgent    *string    `json:"user_agent"`
	TraceID      *string    `json:"trace_id"`
}

type IntegrityCheckResult struct {
	Valid       bool       `json:"valid"`
	CheckedFrom time.Time  `json:"checked_from"`
	CheckedTo   time.Time  `json:"checked_to"`
	ChainLength int        `json:"chain_length"`
	BrokenAt    *time.Time `json:"broken_at,omitempty"`
	Message     string     `json:"message"`
}

type QueryParams struct {
	ResourceType string
	Action       string
	DateFrom     *time.Time
	DateTo       *time.Time
	Limit        int
	Offset       int
}
