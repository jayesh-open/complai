package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type FilingStatus string

const (
	FilingStatusDraft     FilingStatus = "draft"
	FilingStatusIngested  FilingStatus = "ingested"
	FilingStatusValidated FilingStatus = "validated"
	FilingStatusApproved  FilingStatus = "approved"
	FilingStatusSaved     FilingStatus = "saved"
	FilingStatusSubmitted FilingStatus = "submitted"
	FilingStatusFiled     FilingStatus = "filed"
	FilingStatusFailed    FilingStatus = "failed"
)

type SalesRegisterEntry struct {
	ID             uuid.UUID       `json:"id" db:"id"`
	TenantID       uuid.UUID       `json:"tenant_id" db:"tenant_id"`
	GSTIN          string          `json:"gstin" db:"gstin"`
	ReturnPeriod   string          `json:"return_period" db:"return_period"`
	DocumentNumber string          `json:"document_number" db:"document_number"`
	DocumentDate   string          `json:"document_date" db:"document_date"`
	DocumentType   string          `json:"document_type" db:"document_type"`
	SupplyType     string          `json:"supply_type" db:"supply_type"`
	ReverseCharge  bool            `json:"reverse_charge" db:"reverse_charge"`
	SupplierGSTIN  string          `json:"supplier_gstin" db:"supplier_gstin"`
	BuyerGSTIN     string          `json:"buyer_gstin" db:"buyer_gstin"`
	BuyerName      string          `json:"buyer_name" db:"buyer_name"`
	BuyerState     string          `json:"buyer_state" db:"buyer_state"`
	PlaceOfSupply  string          `json:"place_of_supply" db:"place_of_supply"`
	HSN            string          `json:"hsn" db:"hsn"`
	TaxableValue   decimal.Decimal `json:"taxable_value" db:"taxable_value"`
	CGSTRate       decimal.Decimal `json:"cgst_rate" db:"cgst_rate"`
	CGSTAmount     decimal.Decimal `json:"cgst_amount" db:"cgst_amount"`
	SGSTRate       decimal.Decimal `json:"sgst_rate" db:"sgst_rate"`
	SGSTAmount     decimal.Decimal `json:"sgst_amount" db:"sgst_amount"`
	IGSTRate       decimal.Decimal `json:"igst_rate" db:"igst_rate"`
	IGSTAmount     decimal.Decimal `json:"igst_amount" db:"igst_amount"`
	GrandTotal     decimal.Decimal `json:"grand_total" db:"grand_total"`
	SourceSystem   string          `json:"source_system" db:"source_system"`
	SourceID       string          `json:"source_id" db:"source_id"`
	Section        string          `json:"section" db:"section"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
}

type GSTR1Filing struct {
	ID           uuid.UUID    `json:"id" db:"id"`
	TenantID     uuid.UUID    `json:"tenant_id" db:"tenant_id"`
	GSTIN        string       `json:"gstin" db:"gstin"`
	ReturnPeriod string       `json:"return_period" db:"return_period"`
	Status       FilingStatus `json:"status" db:"status"`
	TotalCount   int          `json:"total_count" db:"total_count"`
	ErrorCount   int          `json:"error_count" db:"error_count"`
	ARN          string       `json:"arn,omitempty" db:"arn"`
	FiledAt      *time.Time   `json:"filed_at,omitempty" db:"filed_at"`
	FiledBy      *uuid.UUID   `json:"filed_by,omitempty" db:"filed_by"`
	ApprovedBy   *uuid.UUID   `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt   *time.Time   `json:"approved_at,omitempty" db:"approved_at"`
	CreatedBy    *uuid.UUID   `json:"created_by,omitempty" db:"created_by"`
	RequestID    uuid.UUID    `json:"request_id" db:"request_id"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
}

type GSTR1Section struct {
	ID           uuid.UUID       `json:"id" db:"id"`
	TenantID     uuid.UUID       `json:"tenant_id" db:"tenant_id"`
	FilingID     uuid.UUID       `json:"filing_id" db:"filing_id"`
	Section      string          `json:"section" db:"section"`
	InvoiceCount int             `json:"invoice_count" db:"invoice_count"`
	TaxableValue decimal.Decimal `json:"taxable_value" db:"taxable_value"`
	CGST         decimal.Decimal `json:"cgst" db:"cgst"`
	SGST         decimal.Decimal `json:"sgst" db:"sgst"`
	IGST         decimal.Decimal `json:"igst" db:"igst"`
	TotalTax     decimal.Decimal `json:"total_tax" db:"total_tax"`
	Status       string          `json:"status" db:"status"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
}

type ValidationError struct {
	ID         uuid.UUID `json:"id" db:"id"`
	TenantID   uuid.UUID `json:"tenant_id" db:"tenant_id"`
	FilingID   uuid.UUID `json:"filing_id" db:"filing_id"`
	EntryID    uuid.UUID `json:"entry_id" db:"entry_id"`
	Field      string    `json:"field" db:"field"`
	Code       string    `json:"code" db:"code"`
	Message    string    `json:"message" db:"message"`
	Severity   string    `json:"severity" db:"severity"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// GSTR-1 section identifiers
const (
	SectionB2B    = "b2b"
	SectionB2CL   = "b2cl"
	SectionB2CS   = "b2cs"
	SectionCDNR   = "cdnr"
	SectionCDNUR  = "cdnur"
	SectionEXP    = "exp"
	SectionAT     = "at"
	SectionATAdj  = "atadj"
	SectionNIL    = "nil"
	SectionHSN    = "hsn"
	SectionDOCS   = "docs"
)

var AllSections = []string{
	SectionB2B, SectionB2CL, SectionB2CS, SectionCDNR, SectionCDNUR,
	SectionEXP, SectionAT, SectionATAdj, SectionNIL, SectionHSN, SectionDOCS,
}

type GSTR1Summary struct {
	Filing   GSTR1Filing    `json:"filing"`
	Sections []GSTR1Section `json:"sections"`
	Errors   int            `json:"error_count"`
}

// GSTR-3B types

type GSTR3BStatus string

const (
	GSTR3BStatusDraft     GSTR3BStatus = "draft"
	GSTR3BStatusPopulated GSTR3BStatus = "populated"
	GSTR3BStatusReviewed  GSTR3BStatus = "reviewed"
	GSTR3BStatusApproved  GSTR3BStatus = "approved"
	GSTR3BStatusSaved     GSTR3BStatus = "saved"
	GSTR3BStatusSubmitted GSTR3BStatus = "submitted"
	GSTR3BStatusFiled     GSTR3BStatus = "filed"
	GSTR3BStatusFailed    GSTR3BStatus = "failed"
)

type GSTR3BFiling struct {
	ID           uuid.UUID    `json:"id" db:"id"`
	TenantID     uuid.UUID    `json:"tenant_id" db:"tenant_id"`
	GSTIN        string       `json:"gstin" db:"gstin"`
	ReturnPeriod string       `json:"return_period" db:"return_period"`
	Status       GSTR3BStatus `json:"status" db:"status"`
	DataJSON     string       `json:"data_json,omitempty" db:"data_json"`
	ARN          string       `json:"arn,omitempty" db:"arn"`
	FiledAt      *time.Time   `json:"filed_at,omitempty" db:"filed_at"`
	FiledBy      *uuid.UUID   `json:"filed_by,omitempty" db:"filed_by"`
	ApprovedBy   *uuid.UUID   `json:"approved_by,omitempty" db:"approved_by"`
	ApprovedAt   *time.Time   `json:"approved_at,omitempty" db:"approved_at"`
	CreatedBy    *uuid.UUID   `json:"created_by,omitempty" db:"created_by"`
	RequestID    uuid.UUID    `json:"request_id" db:"request_id"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
}

type GSTR3BAutoFillData struct {
	GSTR1Summary   GSTR3BOutwardSupply `json:"gstr1_summary"`
	InwardSupply   GSTR3BInwardSupply  `json:"inward_supply"`
	IMSActions     GSTR3BIMSSummary    `json:"ims_actions"`
	EligibleITC    GSTR3BITC           `json:"eligible_itc"`
	GrossLiability GSTR3BTaxAmount     `json:"gross_liability"`
	NetLiability   GSTR3BTaxAmount     `json:"net_liability"`
	Flags          []string            `json:"flags"`
}

type GSTR3BOutwardSupply struct {
	B2B        GSTR3BTaxRow `json:"b2b"`
	B2CL       GSTR3BTaxRow `json:"b2cl"`
	B2CS       GSTR3BTaxRow `json:"b2cs"`
	Exports    GSTR3BTaxRow `json:"exports"`
	NIL        GSTR3BTaxRow `json:"nil_exempt"`
	CreditNote GSTR3BTaxRow `json:"credit_notes"`
	Advances   GSTR3BTaxRow `json:"advances"`
}

type GSTR3BTaxRow struct {
	TaxableValue decimal.Decimal `json:"taxable_value"`
	CGST         decimal.Decimal `json:"cgst"`
	SGST         decimal.Decimal `json:"sgst"`
	IGST         decimal.Decimal `json:"igst"`
	Cess         decimal.Decimal `json:"cess"`
}

type GSTR3BInwardSupply struct {
	TotalValue   decimal.Decimal `json:"total_value"`
	TaxableAt5   decimal.Decimal `json:"taxable_at_5"`
	TaxableAt12  decimal.Decimal `json:"taxable_at_12"`
	TaxableAt18  decimal.Decimal `json:"taxable_at_18"`
	TaxableAt28  decimal.Decimal `json:"taxable_at_28"`
	ExemptAndRCM decimal.Decimal `json:"exempt_and_rcm"`
	ITCAvailable GSTR3BTaxAmount `json:"itc_available"`
}

type GSTR3BIMSSummary struct {
	Accepted      int             `json:"accepted"`
	Rejected      int             `json:"rejected"`
	Pending       int             `json:"pending"`
	AcceptedValue decimal.Decimal `json:"accepted_value"`
	RejectedValue decimal.Decimal `json:"rejected_value"`
	PendingValue  decimal.Decimal `json:"pending_value"`
}

type GSTR3BITC struct {
	ImportGoods     GSTR3BTaxAmount `json:"import_goods"`
	ImportServices  GSTR3BTaxAmount `json:"import_services"`
	InwardRCM       GSTR3BTaxAmount `json:"inward_rcm"`
	InwardISD       GSTR3BTaxAmount `json:"inward_isd"`
	AllOther        GSTR3BTaxAmount `json:"all_other"`
	Total           GSTR3BTaxAmount `json:"total"`
}

type GSTR3BTaxAmount struct {
	CGST decimal.Decimal `json:"cgst"`
	SGST decimal.Decimal `json:"sgst"`
	IGST decimal.Decimal `json:"igst"`
	Cess decimal.Decimal `json:"cess"`
}

type GSTR3BAutoFillRequest struct {
	GSTIN        string `json:"gstin" validate:"required"`
	ReturnPeriod string `json:"return_period" validate:"required"`
}

type GSTR3BAutoFillResponse struct {
	FilingID uuid.UUID          `json:"filing_id"`
	Data     GSTR3BAutoFillData `json:"data"`
}

type GSTR3BApproveRequest struct {
	FilingID   uuid.UUID `json:"filing_id" validate:"required"`
	ApprovedBy uuid.UUID `json:"approved_by" validate:"required"`
}

type GSTR3BFileRequest struct {
	FilingID uuid.UUID `json:"filing_id" validate:"required"`
	SignType string    `json:"sign_type" validate:"required,oneof=DSC EVC"`
	OTP      string    `json:"otp,omitempty"`
	FiledBy  uuid.UUID `json:"filed_by" validate:"required"`
}

type GSTR3BFileResponse struct {
	FilingID uuid.UUID    `json:"filing_id"`
	Status   GSTR3BStatus `json:"status"`
	ARN      string       `json:"arn,omitempty"`
}

type IngestRequest struct {
	GSTIN        string `json:"gstin" validate:"required"`
	ReturnPeriod string `json:"return_period" validate:"required"`
}

type IngestResponse struct {
	FilingID   uuid.UUID `json:"filing_id"`
	Ingested   int       `json:"ingested_count"`
	Duplicates int       `json:"duplicate_count"`
}

type ValidateRequest struct {
	FilingID uuid.UUID `json:"filing_id" validate:"required"`
}

type ValidateResponse struct {
	FilingID   uuid.UUID `json:"filing_id"`
	TotalCount int       `json:"total_count"`
	ErrorCount int       `json:"error_count"`
	Sections   []GSTR1Section `json:"sections"`
}

type ApproveRequest struct {
	FilingID   uuid.UUID `json:"filing_id" validate:"required"`
	ApprovedBy uuid.UUID `json:"approved_by" validate:"required"`
}

type FileRequest struct {
	FilingID uuid.UUID `json:"filing_id" validate:"required"`
	SignType string    `json:"sign_type" validate:"required,oneof=DSC EVC"`
	OTP      string    `json:"otp,omitempty"`
	FiledBy  uuid.UUID `json:"filed_by" validate:"required"`
}

type FileResponse struct {
	FilingID uuid.UUID    `json:"filing_id"`
	Status   FilingStatus `json:"status"`
	ARN      string       `json:"arn,omitempty"`
}
