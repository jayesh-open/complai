package domain

import (
	"time"

	"github.com/google/uuid"
)

type FormType string

const (
	FormType138 FormType = "138"
	FormType140 FormType = "140"
	FormType144 FormType = "144"
)

func ValidFormType(f FormType) bool {
	switch f {
	case FormType138, FormType140, FormType144:
		return true
	}
	return false
}

type FilingStatus string

const (
	FilingDraft     FilingStatus = "DRAFT"
	FilingSubmitted FilingStatus = "SUBMITTED"
	FilingFiled     FilingStatus = "FILED"
	FilingRejected  FilingStatus = "REJECTED"
)

type Filing struct {
	ID                    uuid.UUID    `json:"id"`
	TenantID              uuid.UUID    `json:"tenant_id"`
	FormType              FormType     `json:"form_type"`
	FinancialYear         string       `json:"financial_year"`
	Quarter               string       `json:"quarter"`
	TAN                   string       `json:"tan"`
	Status                FilingStatus `json:"status"`
	DeducteeCount         int          `json:"deductee_count"`
	TotalTDSAmount        string       `json:"total_tds_amount"`
	FVUContent            string       `json:"fvu_content,omitempty"`
	TokenNumber           string       `json:"token_number,omitempty"`
	AcknowledgementNumber string       `json:"acknowledgement_number,omitempty"`
	FilingDate            *time.Time   `json:"filing_date,omitempty"`
	ErrorMessage          string       `json:"error_message,omitempty"`
	IdempotencyKey        string       `json:"idempotency_key"`
	CreatedAt             time.Time    `json:"created_at"`
	UpdatedAt             time.Time    `json:"updated_at"`
}

func FilingIdempotencyKey(tenantID uuid.UUID, formType FormType, fy, quarter string) string {
	return tenantID.String() + ":" + string(formType) + ":" + fy + ":" + quarter
}
