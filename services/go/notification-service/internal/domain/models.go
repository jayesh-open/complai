package domain

import (
	"time"

	"github.com/google/uuid"
)

type NotificationTemplate struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	Name      string    `json:"name"`
	Channel   string    `json:"channel"`
	Subject   *string   `json:"subject,omitempty"`
	Body      string    `json:"body"`
	Variables string    `json:"variables"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NotificationPreference struct {
	ID               uuid.UUID  `json:"id"`
	TenantID         uuid.UUID  `json:"tenant_id"`
	UserID           uuid.UUID  `json:"user_id"`
	EmailEnabled     bool       `json:"email_enabled"`
	DigestEnabled    bool       `json:"digest_enabled"`
	QuietHoursStart  *string    `json:"quiet_hours_start,omitempty"`
	QuietHoursEnd    *string    `json:"quiet_hours_end,omitempty"`
	EmailAddress     *string    `json:"email_address,omitempty"`
	EmailValid       bool       `json:"email_valid"`
	BounceCount      int        `json:"bounce_count"`
	UnsubscribeToken uuid.UUID  `json:"unsubscribe_token"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type Notification struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	UserID       uuid.UUID  `json:"user_id"`
	TemplateID   *uuid.UUID `json:"template_id,omitempty"`
	Channel      string     `json:"channel"`
	Subject      *string    `json:"subject,omitempty"`
	Body         *string    `json:"body,omitempty"`
	Recipient    string     `json:"recipient"`
	Status       string     `json:"status"`
	SentAt       *time.Time `json:"sent_at,omitempty"`
	FailedReason *string    `json:"failed_reason,omitempty"`
	DigestGroup  *string    `json:"digest_group,omitempty"`
	DigestBatchID *uuid.UUID `json:"digest_batch_id,omitempty"`
	Metadata     string     `json:"metadata"`
	CreatedAt    time.Time  `json:"created_at"`
}

type NotificationBounce struct {
	ID             uuid.UUID  `json:"id"`
	TenantID       uuid.UUID  `json:"tenant_id"`
	NotificationID *uuid.UUID `json:"notification_id,omitempty"`
	BounceType     string     `json:"bounce_type"`
	BounceSubtype  *string    `json:"bounce_subtype,omitempty"`
	EmailAddress   string     `json:"email_address"`
	Diagnostic     *string    `json:"diagnostic,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

// Request types

type SendNotificationRequest struct {
	UserID      uuid.UUID  `json:"user_id"`
	TemplateID  *uuid.UUID `json:"template_id,omitempty"`
	Channel     string     `json:"channel"`
	Subject     string     `json:"subject"`
	Body        string     `json:"body"`
	Recipient   string     `json:"recipient"`
	DigestGroup string     `json:"digest_group,omitempty"`
	Metadata    string     `json:"metadata,omitempty"`
}

type UpdatePreferencesRequest struct {
	EmailEnabled    *bool   `json:"email_enabled,omitempty"`
	DigestEnabled   *bool   `json:"digest_enabled,omitempty"`
	QuietHoursStart *string `json:"quiet_hours_start,omitempty"`
	QuietHoursEnd   *string `json:"quiet_hours_end,omitempty"`
	EmailAddress    *string `json:"email_address,omitempty"`
}

type CreateTemplateRequest struct {
	Name      string  `json:"name"`
	Channel   string  `json:"channel"`
	Subject   *string `json:"subject,omitempty"`
	Body      string  `json:"body"`
	Variables string  `json:"variables,omitempty"`
}

type ProcessBounceRequest struct {
	NotificationID *uuid.UUID `json:"notification_id,omitempty"`
	BounceType     string     `json:"bounce_type"`
	BounceSubtype  string     `json:"bounce_subtype"`
	EmailAddress   string     `json:"email_address"`
	Diagnostic     string     `json:"diagnostic"`
}

type SendDigestRequest struct {
	DigestGroup string `json:"digest_group"`
	CutoffTime  string `json:"cutoff_time"` // RFC3339
}

type DigestResult struct {
	UserID            uuid.UUID  `json:"user_id"`
	NotificationCount int        `json:"notification_count"`
	DigestSent        bool       `json:"digest_sent"`
	DigestID          *uuid.UUID `json:"digest_id,omitempty"`
}
