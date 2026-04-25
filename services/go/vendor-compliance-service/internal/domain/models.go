package domain

import (
	"time"

	"github.com/google/uuid"
)

type VendorSnapshot struct {
	ID                 uuid.UUID `json:"id"`
	TenantID           uuid.UUID `json:"tenant_id"`
	VendorID           string    `json:"vendor_id"`
	Name               string    `json:"name"`
	LegalName          string    `json:"legal_name"`
	TradeName          string    `json:"trade_name"`
	PAN                string    `json:"pan"`
	GSTIN              string    `json:"gstin"`
	TAN                string    `json:"tan,omitempty"`
	State              string    `json:"state"`
	StateCode          string    `json:"state_code"`
	Category           string    `json:"category"`
	RegistrationStatus string    `json:"registration_status"`
	MSMERegistered     bool      `json:"msme_registered"`
	Email              string    `json:"email"`
	Phone              string    `json:"phone"`
	Address            string    `json:"address"`
	SyncedAt           time.Time `json:"synced_at"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type ComplianceScore struct {
	ID                    uuid.UUID `json:"id"`
	TenantID              uuid.UUID `json:"tenant_id"`
	VendorID              string    `json:"vendor_id"`
	VendorSnapshotID      uuid.UUID `json:"vendor_snapshot_id"`
	TotalScore            int       `json:"total_score"`
	Category              string    `json:"category"`
	RiskLevel             string    `json:"risk_level"`
	FilingRegularityScore int       `json:"filing_regularity_score"`
	IRNComplianceScore    int       `json:"irn_compliance_score"`
	MismatchRateScore     int       `json:"mismatch_rate_score"`
	PaymentBehaviorScore  int       `json:"payment_behavior_score"`
	DocumentHygieneScore  int       `json:"document_hygiene_score"`
	FilingRegularityNote  string    `json:"filing_regularity_note,omitempty"`
	IRNComplianceNote     string    `json:"irn_compliance_note,omitempty"`
	MismatchRateNote      string    `json:"mismatch_rate_note,omitempty"`
	PaymentBehaviorNote   string    `json:"payment_behavior_note,omitempty"`
	DocumentHygieneNote   string    `json:"document_hygiene_note,omitempty"`
	ScoredAt              time.Time `json:"scored_at"`
	CreatedAt             time.Time `json:"created_at"`
}

type SyncStatus struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	SyncType     string     `json:"sync_type"`
	Status       string     `json:"status"`
	VendorCount  int        `json:"vendor_count"`
	ScoredCount  int        `json:"scored_count"`
	StartedAt    time.Time  `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	ErrorMessage string     `json:"error_message,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// Request/Response DTOs

type SyncRequest struct {
	RequestID string `json:"request_id"`
}

type SyncResponse struct {
	SyncID      uuid.UUID `json:"sync_id"`
	VendorCount int       `json:"vendor_count"`
	ScoredCount int       `json:"scored_count"`
	Status      string    `json:"status"`
}

type VendorScoreResponse struct {
	Vendor VendorSnapshot  `json:"vendor"`
	Score  ComplianceScore `json:"score"`
}

type VendorListResponse struct {
	Vendors []VendorScoreResponse `json:"vendors"`
	Total   int                   `json:"total"`
	Summary ScoreSummary          `json:"summary"`
}

type ScoreSummary struct {
	Total    int `json:"total"`
	CatA     int `json:"cat_a"`
	CatB     int `json:"cat_b"`
	CatC     int `json:"cat_c"`
	CatD     int `json:"cat_d"`
	AvgScore int `json:"avg_score"`
}
