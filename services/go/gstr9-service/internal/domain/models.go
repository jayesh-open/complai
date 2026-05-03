package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type FilingStatus string

const (
	FilingStatusDraft      FilingStatus = "draft"
	FilingStatusAggregated FilingStatus = "aggregated"
	FilingStatusReviewed   FilingStatus = "reviewed"
	FilingStatusApproved   FilingStatus = "approved"
	FilingStatusSaved      FilingStatus = "saved"
	FilingStatusSubmitted  FilingStatus = "submitted"
	FilingStatusFiled      FilingStatus = "filed"
	FilingStatusFailed     FilingStatus = "failed"
)

type GSTR9Filing struct {
	ID                uuid.UUID       `json:"id" db:"id"`
	TenantID          uuid.UUID       `json:"tenant_id" db:"tenant_id"`
	GSTIN             string          `json:"gstin" db:"gstin"`
	FinancialYear     string          `json:"financial_year" db:"financial_year"`
	Status            FilingStatus    `json:"status" db:"status"`
	AggregateTurnover decimal.Decimal `json:"aggregate_turnover" db:"aggregate_turnover"`
	IsMandatory       bool            `json:"is_mandatory" db:"is_mandatory"`
	ARN               string          `json:"arn,omitempty" db:"arn"`
	FiledAt           *time.Time      `json:"filed_at,omitempty" db:"filed_at"`
	FiledBy           *uuid.UUID      `json:"filed_by,omitempty" db:"filed_by"`
	ApprovedBy        *uuid.UUID      `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt        *time.Time      `json:"approved_at,omitempty" db:"approved_at"`
	RequestID         uuid.UUID       `json:"request_id" db:"request_id"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at" db:"updated_at"`
}

type GSTR9TableData struct {
	ID           uuid.UUID       `json:"id" db:"id"`
	TenantID     uuid.UUID       `json:"tenant_id" db:"tenant_id"`
	FilingID     uuid.UUID       `json:"filing_id" db:"filing_id"`
	PartNumber   int             `json:"part_number" db:"part_number"`
	TableNumber  string          `json:"table_number" db:"table_number"`
	Description  string          `json:"description" db:"description"`
	TaxableValue decimal.Decimal `json:"taxable_value" db:"taxable_value"`
	CGST         decimal.Decimal `json:"cgst" db:"cgst"`
	SGST         decimal.Decimal `json:"sgst" db:"sgst"`
	IGST         decimal.Decimal `json:"igst" db:"igst"`
	Cess         decimal.Decimal `json:"cess" db:"cess"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
}

type GSTR9AuditLog struct {
	ID        uuid.UUID `json:"id" db:"id"`
	TenantID  uuid.UUID `json:"tenant_id" db:"tenant_id"`
	FilingID  uuid.UUID `json:"filing_id" db:"filing_id"`
	Action    string    `json:"action" db:"action"`
	Details   string    `json:"details" db:"details"`
	ActorID   uuid.UUID `json:"actor_id" db:"actor_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type GSTR9CStatus string

const (
	GSTR9CStatusDraft      GSTR9CStatus = "draft"
	GSTR9CStatusReconciled GSTR9CStatus = "reconciled"
	GSTR9CStatusCertified  GSTR9CStatus = "certified"
	GSTR9CStatusSubmitted  GSTR9CStatus = "submitted"
)

type GSTR9CFiling struct {
	ID                 uuid.UUID       `json:"id" db:"id"`
	TenantID           uuid.UUID       `json:"tenant_id" db:"tenant_id"`
	GSTR9FilingID      uuid.UUID       `json:"gstr9_filing_id" db:"gstr9_filing_id"`
	Status             GSTR9CStatus    `json:"status" db:"status"`
	AuditedTurnover    decimal.Decimal `json:"audited_turnover" db:"audited_turnover"`
	UnreconciledAmount decimal.Decimal `json:"unreconciled_amount" db:"unreconciled_amount"`
	IsSelfCertified    bool            `json:"is_self_certified" db:"is_self_certified"`
	CertifiedAt        *time.Time      `json:"certified_at,omitempty" db:"certified_at"`
	CertifiedBy        *uuid.UUID      `json:"certified_by,omitempty" db:"certified_by"`
	CreatedAt          time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at" db:"updated_at"`
}

type MismatchSeverity string

const (
	SeverityInfo  MismatchSeverity = "INFO"
	SeverityWarn  MismatchSeverity = "WARN"
	SeverityError MismatchSeverity = "ERROR"
)

type GSTR9CMismatch struct {
	ID              uuid.UUID        `json:"id" db:"id"`
	TenantID        uuid.UUID        `json:"tenant_id" db:"tenant_id"`
	GSTR9CFilingID  uuid.UUID        `json:"gstr9c_filing_id" db:"gstr9c_filing_id"`
	Section         string           `json:"section" db:"section"`
	Category        string           `json:"category" db:"category"`
	Description     string           `json:"description" db:"description"`
	BooksAmount     decimal.Decimal  `json:"books_amount" db:"books_amount"`
	GSTR9Amount     decimal.Decimal  `json:"gstr9_amount" db:"gstr9_amount"`
	Difference      decimal.Decimal  `json:"difference" db:"difference"`
	Severity        MismatchSeverity `json:"severity" db:"severity"`
	Reason          string           `json:"reason" db:"reason"`
	SuggestedAction string           `json:"suggested_action" db:"suggested_action"`
	Resolved        bool             `json:"resolved" db:"resolved"`
	ResolvedReason  string           `json:"resolved_reason,omitempty" db:"resolved_reason"`
	ResolvedAt      *time.Time       `json:"resolved_at,omitempty" db:"resolved_at"`
	ResolvedBy      *uuid.UUID       `json:"resolved_by,omitempty" db:"resolved_by"`
	CreatedAt       time.Time        `json:"created_at" db:"created_at"`
}
