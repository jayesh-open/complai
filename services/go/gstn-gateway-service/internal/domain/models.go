package domain

import (
	"time"

	"github.com/google/uuid"
)

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	JTI         string `json:"jti"`
}

type GSTR1SaveRequest struct {
	GSTIN     string          `json:"gstin"`
	RetPeriod string          `json:"ret_period"` // MMYYYY
	Section   string          `json:"section"`    // b2b, b2cl, b2cs, cdnr, cdnur, exp, at, atadj, nil, hsn, docs
	Data      interface{}     `json:"data"`
	RequestID string          `json:"request_id"`
}

type GSTR1SaveResponse struct {
	Status    string    `json:"status"`
	RequestID string    `json:"request_id"`
	Token     string    `json:"token"`
	Message   string    `json:"message"`
	SavedAt   time.Time `json:"saved_at"`
}

type GSTR1GetRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
	Section   string `json:"section,omitempty"`
	RequestID string `json:"request_id"`
}

type GSTR1GetResponse struct {
	GSTIN     string                 `json:"gstin"`
	RetPeriod string                 `json:"ret_period"`
	Data      map[string]interface{} `json:"data"`
	Status    string                 `json:"status"`
	RequestID string                 `json:"request_id"`
}

type GSTR1ResetRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
	RequestID string `json:"request_id"`
}

type GSTR1ResetResponse struct {
	Status    string `json:"status"`
	RequestID string `json:"request_id"`
	Message   string `json:"message"`
}

type GSTR1SubmitRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
	RequestID string `json:"request_id"`
}

type GSTR1SubmitResponse struct {
	Status    string `json:"status"`
	RequestID string `json:"request_id"`
	Token     string `json:"token"`
	Message   string `json:"message"`
}

type GSTR1FileRequest struct {
	GSTIN      string `json:"gstin"`
	RetPeriod  string `json:"ret_period"`
	SignType   string `json:"sign_type"` // DSC or EVC
	EVOTP      string `json:"ev_otp,omitempty"`
	PAN        string `json:"pan"`
	RequestID  string `json:"request_id"`
}

type GSTR1FileResponse struct {
	Status    string `json:"status"`
	ARN       string `json:"arn"`
	RequestID string `json:"request_id"`
	Message   string `json:"message"`
	FiledAt   time.Time `json:"filed_at"`
}

type GSTR1StatusRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
	RequestID string `json:"request_id"`
}

type GSTR1StatusResponse struct {
	GSTIN     string     `json:"gstin"`
	RetPeriod string     `json:"ret_period"`
	Status    string     `json:"status"`
	ARN       string     `json:"arn,omitempty"`
	FiledAt   *time.Time `json:"filed_at,omitempty"`
	RequestID string     `json:"request_id"`
}

type GatewayRequest struct {
	TenantID      uuid.UUID `json:"-"`
	IdempotencyKey string   `json:"-"`
}

type GatewayResponse struct {
	Data interface{} `json:"data"`
	Meta ResponseMeta `json:"meta"`
}

type ResponseMeta struct {
	RequestID      string `json:"request_id"`
	LatencyMs      int    `json:"latency_ms"`
	ProviderStatus string `json:"provider_status"`
}

type FilingStatus string

const (
	StatusDraft     FilingStatus = "draft"
	StatusSaved     FilingStatus = "saved"
	StatusSubmitted FilingStatus = "submitted"
	StatusFiled     FilingStatus = "filed"
)

type MockFiling struct {
	GSTIN     string
	RetPeriod string
	Status    FilingStatus
	Sections  map[string]interface{}
	ARN       string
	FiledAt   *time.Time
	Token     string
}
