package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	ExternalID    *string    `json:"external_id,omitempty"`
	Email         string     `json:"email"`
	EmailVerified bool       `json:"email_verified"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type UserCredential struct {
	ID             uuid.UUID `json:"id"`
	TenantID       uuid.UUID `json:"tenant_id"`
	UserID         uuid.UUID `json:"user_id"`
	Provider       string    `json:"provider"`
	ProviderUserID string    `json:"provider_user_id"`
	CreatedAt      time.Time `json:"created_at"`
}

type UserSession struct {
	ID         uuid.UUID  `json:"id"`
	TenantID   uuid.UUID  `json:"tenant_id"`
	UserID     uuid.UUID  `json:"user_id"`
	DeviceInfo *string    `json:"device_info,omitempty"`
	IPAddress  *string    `json:"ip_address,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  time.Time  `json:"expires_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
}

type MFAFactor struct {
	ID              uuid.UUID `json:"id"`
	TenantID        uuid.UUID `json:"tenant_id"`
	UserID          uuid.UUID `json:"user_id"`
	FactorType      string    `json:"factor_type"`
	SecretEncrypted *string   `json:"-"`
	PhoneNumber     *string   `json:"phone_number,omitempty"`
	Verified        bool      `json:"verified"`
	CreatedAt       time.Time `json:"created_at"`
}

type StepUpEvent struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	UserID      uuid.UUID `json:"user_id"`
	SessionID   uuid.UUID `json:"session_id"`
	ActionClass string    `json:"action_class"`
	VerifiedAt  time.Time `json:"verified_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	MFAMethod   string    `json:"mfa_method"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type StepUpRequest struct {
	ActionClass string `json:"action_class" validate:"required"`
	MFACode     string `json:"mfa_code" validate:"required"`
}

type StepUpCheckRequest struct {
	ActionClass string `json:"action_class" validate:"required"`
}
