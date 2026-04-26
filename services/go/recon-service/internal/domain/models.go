package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type MatchType string

const (
	MatchTypeDirect    MatchType = "DIRECT"
	MatchTypeProbable  MatchType = "PROBABLE"
	MatchTypePartial   MatchType = "PARTIAL"
	MatchTypeMissing2B MatchType = "MISSING_2B"
	MatchTypeMissingPR MatchType = "MISSING_PR"
	MatchTypeDuplicate MatchType = "DUPLICATE"
)

type MatchStatus string

const (
	MatchStatusUnreviewed MatchStatus = "UNREVIEWED"
	MatchStatusAccepted   MatchStatus = "ACCEPTED"
	MatchStatusRejected   MatchStatus = "REJECTED"
	MatchStatusFlagged    MatchStatus = "FLAGGED"
)

type ReconRun struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	GSTIN        string     `json:"gstin"`
	ReturnPeriod string     `json:"return_period"`
	Status       string     `json:"status"` // RUNNING, COMPLETED, FAILED
	PRCount      int        `json:"pr_count"`
	GSTR2BCount  int        `json:"gstr2b_count"`
	Matched      int        `json:"matched"`
	Mismatch     int        `json:"mismatch"`
	Partial      int        `json:"partial"`
	Missing2B    int        `json:"missing_2b"`
	MissingPR    int        `json:"missing_pr"`
	Duplicate    int        `json:"duplicate"`
	StartedAt    time.Time  `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	RequestID    uuid.UUID  `json:"request_id"`
	CreatedAt    time.Time  `json:"created_at"`
}

type ReconMatch struct {
	ID           uuid.UUID `json:"id"`
	TenantID     uuid.UUID `json:"tenant_id"`
	RunID        uuid.UUID `json:"run_id"`
	GSTIN        string    `json:"gstin"`
	ReturnPeriod string    `json:"return_period"`

	PRInvoiceNumber string          `json:"pr_invoice_number,omitempty"`
	PRInvoiceDate   string          `json:"pr_invoice_date,omitempty"`
	PRVendorGSTIN   string          `json:"pr_vendor_gstin,omitempty"`
	PRAmount        decimal.Decimal `json:"pr_amount"`
	PRHSN           string          `json:"pr_hsn,omitempty"`
	PRSourceID      string          `json:"pr_source_id,omitempty"`

	GSTR2BInvoiceNumber string          `json:"gstr2b_invoice_number,omitempty"`
	GSTR2BInvoiceDate   string          `json:"gstr2b_invoice_date,omitempty"`
	GSTR2BSupplierGSTIN string          `json:"gstr2b_supplier_gstin,omitempty"`
	GSTR2BAmount        decimal.Decimal `json:"gstr2b_amount"`
	GSTR2BHSN           string          `json:"gstr2b_hsn,omitempty"`

	MatchType       MatchType       `json:"match_type"`
	MatchConfidence decimal.Decimal `json:"match_confidence"`
	ReasonCodes     []string        `json:"reason_codes"`
	Status          MatchStatus     `json:"status"`
	AcceptedBy      *uuid.UUID      `json:"accepted_by,omitempty"`
	AcceptedAt      *time.Time      `json:"accepted_at,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type IMSAction struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	GSTIN        string     `json:"gstin"`
	ReturnPeriod string     `json:"return_period"`
	InvoiceID    string     `json:"invoice_id"`
	Action       string     `json:"action"` // ACCEPT, REJECT, PENDING
	Reason       string     `json:"reason,omitempty"`
	SyncedToGSTN bool       `json:"synced_to_gstn"`
	SyncedAt     *time.Time `json:"synced_at,omitempty"`
	CreatedBy    uuid.UUID  `json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
}

// Gateway DTOs

type PurchaseRegisterEntry struct {
	InvoiceNumber string          `json:"invoice_number"`
	InvoiceDate   string          `json:"invoice_date"`
	VendorGSTIN   string          `json:"vendor_gstin"`
	VendorName    string          `json:"vendor_name"`
	TaxableValue  decimal.Decimal `json:"taxable_value"`
	CGSTAmount    decimal.Decimal `json:"cgst_amount"`
	SGSTAmount    decimal.Decimal `json:"sgst_amount"`
	IGSTAmount    decimal.Decimal `json:"igst_amount"`
	TotalValue    decimal.Decimal `json:"total_value"`
	HSN           string          `json:"hsn"`
	PlaceOfSupply string          `json:"place_of_supply"`
	ReverseCharge bool            `json:"reverse_charge"`
	SourceID      string          `json:"source_id"`
}

type GSTR2BEntry struct {
	SupplierGSTIN string          `json:"supplier_gstin"`
	InvoiceNumber string          `json:"invoice_number"`
	InvoiceDate   string          `json:"invoice_date"`
	TaxableValue  decimal.Decimal `json:"taxable_value"`
	CGSTAmount    decimal.Decimal `json:"cgst_amount"`
	SGSTAmount    decimal.Decimal `json:"sgst_amount"`
	IGSTAmount    decimal.Decimal `json:"igst_amount"`
	TotalValue    decimal.Decimal `json:"total_value"`
	HSN           string          `json:"hsn"`
	PlaceOfSupply string          `json:"place_of_supply"`
	ReverseCharge bool            `json:"reverse_charge"`
	IMSAction     string          `json:"ims_action"`
}

type BucketSummary struct {
	Matched   int `json:"matched"`
	Mismatch  int `json:"mismatch"`
	Partial   int `json:"partial"`
	Missing2B int `json:"missing_2b"`
	MissingPR int `json:"missing_pr"`
	Duplicate int `json:"duplicate"`
}

// Request/Response types

type RunReconRequest struct {
	GSTIN        string `json:"gstin" validate:"required"`
	ReturnPeriod string `json:"return_period" validate:"required"`
}

type RunReconResponse struct {
	RunID  uuid.UUID `json:"run_id"`
	Status string    `json:"status"`
}

type AcceptMatchRequest struct {
	MatchID uuid.UUID `json:"match_id" validate:"required"`
}

type BulkAcceptRequest struct {
	MatchIDs []uuid.UUID `json:"match_ids" validate:"required"`
}

type IMSActionRequest struct {
	InvoiceID string `json:"invoice_id" validate:"required"`
	Action    string `json:"action" validate:"required,oneof=ACCEPT REJECT PENDING"`
	Reason    string `json:"reason,omitempty"`
}
