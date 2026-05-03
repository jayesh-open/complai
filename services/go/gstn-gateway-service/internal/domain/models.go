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

// GSTR-2B types
type GSTR2BGetRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
	RequestID string `json:"request_id"`
}

type GSTR2BInvoice struct {
	SupplierGSTIN string  `json:"supplier_gstin"`
	InvoiceNumber string  `json:"invoice_number"`
	InvoiceDate   string  `json:"invoice_date"`
	InvoiceType   string  `json:"invoice_type"`
	TaxableValue  float64 `json:"taxable_value"`
	CGSTAmount    float64 `json:"cgst_amount"`
	SGSTAmount    float64 `json:"sgst_amount"`
	IGSTAmount    float64 `json:"igst_amount"`
	TotalValue    float64 `json:"total_value"`
	PlaceOfSupply string  `json:"place_of_supply"`
	ReverseCharge bool    `json:"reverse_charge"`
	HSN           string  `json:"hsn"`
	ITC           string  `json:"itc"`        // eligible, ineligible
	IMSAction     string  `json:"ims_action"` // ACCEPT, REJECT, PENDING, ""
}

type GSTR2BGetResponse struct {
	GSTIN       string          `json:"gstin"`
	RetPeriod   string          `json:"ret_period"`
	Invoices    []GSTR2BInvoice `json:"invoices"`
	TotalCount  int             `json:"total_count"`
	GeneratedOn string          `json:"generated_on"`
	Status      string          `json:"status"`
	RequestID   string          `json:"request_id"`
}

// GSTR-2A types
type GSTR2AGetRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
	Section   string `json:"section"` // B2B, B2BA, CDN, CDNA, ISD, ISDA, IMPG, IMPGSEZ, TDS, TCS
	RequestID string `json:"request_id"`
}

type GSTR2AGetResponse struct {
	GSTIN     string          `json:"gstin"`
	RetPeriod string          `json:"ret_period"`
	Section   string          `json:"section"`
	Invoices  []GSTR2BInvoice `json:"invoices"`
	Status    string          `json:"status"`
	RequestID string          `json:"request_id"`
}

// IMS types
type IMSGetRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
	RequestID string `json:"request_id"`
}

type IMSInvoice struct {
	InvoiceID     string  `json:"invoice_id"`
	SupplierGSTIN string  `json:"supplier_gstin"`
	InvoiceNumber string  `json:"invoice_number"`
	InvoiceDate   string  `json:"invoice_date"`
	TaxableValue  float64 `json:"taxable_value"`
	TotalValue    float64 `json:"total_value"`
	CGSTAmount    float64 `json:"cgst_amount"`
	SGSTAmount    float64 `json:"sgst_amount"`
	IGSTAmount    float64 `json:"igst_amount"`
	Action        string  `json:"action"` // ACCEPT, REJECT, PENDING
	ActionAt      string  `json:"action_at,omitempty"`
	ActionBy      string  `json:"action_by,omitempty"`
}

type IMSGetResponse struct {
	GSTIN      string       `json:"gstin"`
	RetPeriod  string       `json:"ret_period"`
	Invoices   []IMSInvoice `json:"invoices"`
	TotalCount int          `json:"total_count"`
	Summary    IMSSummary   `json:"summary"`
	Status     string       `json:"status"`
	RequestID  string       `json:"request_id"`
}

type IMSSummary struct {
	Accepted      int     `json:"accepted"`
	Rejected      int     `json:"rejected"`
	Pending       int     `json:"pending"`
	AcceptedValue float64 `json:"accepted_value"`
	RejectedValue float64 `json:"rejected_value"`
	PendingValue  float64 `json:"pending_value"`
}

type IMSActionRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
	InvoiceID string `json:"invoice_id"`
	Action    string `json:"action"` // ACCEPT, REJECT, PENDING
	Reason    string `json:"reason,omitempty"`
	RequestID string `json:"request_id"`
}

type IMSActionResponse struct {
	InvoiceID string `json:"invoice_id"`
	Action    string `json:"action"`
	Status    string `json:"status"`
	RequestID string `json:"request_id"`
	UpdatedAt string `json:"updated_at"`
}

type IMSBulkActionRequest struct {
	GSTIN      string   `json:"gstin"`
	RetPeriod  string   `json:"ret_period"`
	InvoiceIDs []string `json:"invoice_ids"`
	Action     string   `json:"action"`
	RequestID  string   `json:"request_id"`
}

type IMSBulkActionResponse struct {
	OperationID   string `json:"operation_id"`
	TotalInvoices int    `json:"total_invoices"`
	Status        string `json:"status"`
	RequestID     string `json:"request_id"`
}

// GSTR-3B types
type GSTR3BSaveRequest struct {
	GSTIN     string      `json:"gstin"`
	RetPeriod string      `json:"ret_period"`
	Data      interface{} `json:"data"`
	RequestID string      `json:"request_id"`
}

type GSTR3BSaveResponse struct {
	Status    string `json:"status"`
	RequestID string `json:"request_id"`
	Message   string `json:"message"`
}

type GSTR3BSubmitRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
	RequestID string `json:"request_id"`
}

type GSTR3BSubmitResponse struct {
	Status    string `json:"status"`
	RequestID string `json:"request_id"`
	Message   string `json:"message"`
}

type GSTR3BFileRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
	SignType  string `json:"sign_type"`
	EVOTP     string `json:"ev_otp,omitempty"`
	PAN       string `json:"pan"`
	RequestID string `json:"request_id"`
}

type GSTR3BFileResponse struct {
	Status    string    `json:"status"`
	ARN       string    `json:"arn"`
	RequestID string    `json:"request_id"`
	Message   string    `json:"message"`
	FiledAt   time.Time `json:"filed_at"`
}

type GSTR1SummaryRequest struct {
	GSTIN     string `json:"gstin"`
	RetPeriod string `json:"ret_period"`
	RequestID string `json:"request_id"`
}

type GSTR1SummaryResponse struct {
	GSTIN     string                 `json:"gstin"`
	RetPeriod string                 `json:"ret_period"`
	Summary   map[string]interface{} `json:"summary"`
	Status    string                 `json:"status"`
	RequestID string                 `json:"request_id"`
}

// Mock filing for GSTR-3B
type MockGSTR3BFiling struct {
	GSTIN     string
	RetPeriod string
	Status    FilingStatus
	Data      interface{}
	ARN       string
	FiledAt   *time.Time
}

// GSTR-9 types

type GSTR9SaveRequest struct {
	GSTIN         string      `json:"gstin"`
	FinancialYear string      `json:"financial_year"` // YYYY-YY e.g. 2025-26
	Data          interface{} `json:"data"`
	RequestID     string      `json:"request_id"`
}

type GSTR9SaveResponse struct {
	Status    string    `json:"status"`
	Reference string    `json:"reference"`
	RequestID string    `json:"request_id"`
	Message   string    `json:"message"`
	SavedAt   time.Time `json:"saved_at"`
}

type GSTR9SubmitRequest struct {
	GSTIN         string `json:"gstin"`
	FinancialYear string `json:"financial_year"`
	Reference     string `json:"reference"`
	RequestID     string `json:"request_id"`
}

type GSTR9SubmitResponse struct {
	Status    string `json:"status"`
	Reference string `json:"reference"`
	RequestID string `json:"request_id"`
	Message   string `json:"message"`
}

type GSTR9FileRequest struct {
	GSTIN         string `json:"gstin"`
	FinancialYear string `json:"financial_year"`
	Reference     string `json:"reference"`
	SignType      string `json:"sign_type"` // DSC or EVC
	EVOTP         string `json:"ev_otp,omitempty"`
	PAN           string `json:"pan"`
	RequestID     string `json:"request_id"`
}

type GSTR9FileResponse struct {
	Status    string    `json:"status"`
	ARN       string    `json:"arn"`
	RequestID string    `json:"request_id"`
	Message   string    `json:"message"`
	FiledAt   time.Time `json:"filed_at"`
}

type GSTR9StatusRequest struct {
	Reference string `json:"reference"`
	RequestID string `json:"request_id"`
}

type GSTR9StatusResponse struct {
	Reference string     `json:"reference"`
	Status    string     `json:"status"`
	ARN       string     `json:"arn,omitempty"`
	FiledAt   *time.Time `json:"filed_at,omitempty"`
	RequestID string     `json:"request_id"`
}

// GSTR-9C types

type GSTR9CSaveRequest struct {
	GSTIN         string      `json:"gstin"`
	FinancialYear string      `json:"financial_year"`
	Data          interface{} `json:"data"`
	RequestID     string      `json:"request_id"`
}

type GSTR9CSaveResponse struct {
	Status    string    `json:"status"`
	Reference string    `json:"reference"`
	RequestID string    `json:"request_id"`
	Message   string    `json:"message"`
	SavedAt   time.Time `json:"saved_at"`
}

type GSTR9CFileRequest struct {
	GSTIN         string `json:"gstin"`
	FinancialYear string `json:"financial_year"`
	Reference     string `json:"reference"`
	PAN           string `json:"pan"`
	RequestID     string `json:"request_id"`
}

type GSTR9CFileResponse struct {
	Status    string    `json:"status"`
	ARN       string    `json:"arn"`
	RequestID string    `json:"request_id"`
	Message   string    `json:"message"`
	FiledAt   time.Time `json:"filed_at"`
}

type GSTR9CStatusRequest struct {
	Reference string `json:"reference"`
	RequestID string `json:"request_id"`
}

type GSTR9CStatusResponse struct {
	Reference string     `json:"reference"`
	Status    string     `json:"status"`
	ARN       string     `json:"arn,omitempty"`
	FiledAt   *time.Time `json:"filed_at,omitempty"`
	RequestID string     `json:"request_id"`
}

// Mock filing for GSTR-9 annual return
type MockGSTR9Filing struct {
	GSTIN         string
	FinancialYear string
	Status        FilingStatus
	Reference     string
	Data          interface{}
	ARN           string
	FiledAt       *time.Time
	SavedAt       time.Time
}

// Mock filing for GSTR-9C reconciliation
type MockGSTR9CFiling struct {
	GSTIN         string
	FinancialYear string
	Status        FilingStatus
	Reference     string
	Data          interface{}
	ARN           string
	FiledAt       *time.Time
	SavedAt       time.Time
}
