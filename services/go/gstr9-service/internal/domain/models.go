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
