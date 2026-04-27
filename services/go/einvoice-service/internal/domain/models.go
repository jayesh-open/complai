package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type IRNStatus string

const (
	IRNStatusPending   IRNStatus = "PENDING"
	IRNStatusGenerated IRNStatus = "GENERATED"
	IRNStatusCancelled IRNStatus = "CANCELLED"
	IRNStatusFailed    IRNStatus = "FAILED"
)

type InvoiceType string

const (
	InvoiceTypeINV InvoiceType = "INV"
	InvoiceTypeCRN InvoiceType = "CRN"
	InvoiceTypeDBN InvoiceType = "DBN"
)

type SupplyType string

const (
	SupplyTypeB2B   SupplyType = "B2B"
	SupplyTypeB2CL  SupplyType = "B2CL"
	SupplyTypeB2CS  SupplyType = "B2CS"
	SupplyTypeCDNR  SupplyType = "CDNR"
	SupplyTypeCDNUR SupplyType = "CDNUR"
	SupplyTypeEXP   SupplyType = "EXP"
)

type EInvoice struct {
	ID                   uuid.UUID       `json:"id"`
	TenantID             uuid.UUID       `json:"tenant_id"`
	IRN                  string          `json:"irn"`
	AckNo                string          `json:"ack_no"`
	InvoiceNumber        string          `json:"invoice_number"`
	InvoiceDate          string          `json:"invoice_date"`
	InvoiceType          InvoiceType     `json:"invoice_type"`
	SupplierGSTIN        string          `json:"supplier_gstin"`
	SupplierName         string          `json:"supplier_name"`
	BuyerGSTIN           string          `json:"buyer_gstin"`
	BuyerName            string          `json:"buyer_name"`
	SupplyType           SupplyType      `json:"supply_type"`
	PlaceOfSupply        string          `json:"place_of_supply"`
	ReverseCharge        bool            `json:"reverse_charge"`
	TaxableValue         decimal.Decimal `json:"taxable_value"`
	CGSTAmount           decimal.Decimal `json:"cgst_amount"`
	SGSTAmount           decimal.Decimal `json:"sgst_amount"`
	IGSTAmount           decimal.Decimal `json:"igst_amount"`
	CessAmount           decimal.Decimal `json:"cess_amount"`
	TotalAmount          decimal.Decimal `json:"total_amount"`
	Status               IRNStatus       `json:"status"`
	IRNGeneratedAt       *time.Time      `json:"irn_generated_at,omitempty"`
	IRNCancelledAt       *time.Time      `json:"irn_cancelled_at,omitempty"`
	CancelReason         string          `json:"cancel_reason,omitempty"`
	SignedInvoice        string          `json:"signed_invoice,omitempty"`
	SignedQRCode         string          `json:"signed_qr_code,omitempty"`
	RequestID            uuid.UUID       `json:"request_id"`
	SourceSystem         string          `json:"source_system"`
	SourceID             string          `json:"source_id,omitempty"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

type EInvoiceLineItem struct {
	ID           uuid.UUID       `json:"id"`
	InvoiceID    uuid.UUID       `json:"invoice_id"`
	TenantID     uuid.UUID       `json:"tenant_id"`
	LineNumber   int             `json:"line_number"`
	Description  string          `json:"description"`
	HSNCode      string          `json:"hsn_code"`
	Quantity     decimal.Decimal `json:"quantity"`
	Unit         string          `json:"unit"`
	UnitPrice    decimal.Decimal `json:"unit_price"`
	Discount     decimal.Decimal `json:"discount"`
	TaxableValue decimal.Decimal `json:"taxable_value"`
	CGSTRate     decimal.Decimal `json:"cgst_rate"`
	CGSTAmount   decimal.Decimal `json:"cgst_amount"`
	SGSTRate     decimal.Decimal `json:"sgst_rate"`
	SGSTAmount   decimal.Decimal `json:"sgst_amount"`
	IGSTRate     decimal.Decimal `json:"igst_rate"`
	IGSTAmount   decimal.Decimal `json:"igst_amount"`
	CessRate     decimal.Decimal `json:"cess_rate"`
	CessAmount   decimal.Decimal `json:"cess_amount"`
	CreatedAt    time.Time       `json:"created_at"`
}

// --- API Request/Response types ---

type GenerateIRNRequest struct {
	InvoiceNumber string             `json:"invoice_number"`
	InvoiceDate   string             `json:"invoice_date"`
	InvoiceType   InvoiceType        `json:"invoice_type"`
	SupplierGSTIN string             `json:"supplier_gstin"`
	SupplierName  string             `json:"supplier_name"`
	BuyerGSTIN    string             `json:"buyer_gstin"`
	BuyerName     string             `json:"buyer_name"`
	SupplyType    SupplyType         `json:"supply_type"`
	PlaceOfSupply string             `json:"place_of_supply"`
	ReverseCharge bool               `json:"reverse_charge"`
	LineItems     []LineItemRequest  `json:"line_items"`
	TaxableValue  decimal.Decimal    `json:"taxable_value"`
	CGSTAmount    decimal.Decimal    `json:"cgst_amount"`
	SGSTAmount    decimal.Decimal    `json:"sgst_amount"`
	IGSTAmount    decimal.Decimal    `json:"igst_amount"`
	CessAmount    decimal.Decimal    `json:"cess_amount"`
	TotalAmount   decimal.Decimal    `json:"total_amount"`
	SourceSystem  string             `json:"source_system"`
	SourceID      string             `json:"source_id,omitempty"`
}

type LineItemRequest struct {
	Description  string          `json:"description"`
	HSNCode      string          `json:"hsn_code"`
	Quantity     decimal.Decimal `json:"quantity"`
	Unit         string          `json:"unit"`
	UnitPrice    decimal.Decimal `json:"unit_price"`
	Discount     decimal.Decimal `json:"discount"`
	TaxableValue decimal.Decimal `json:"taxable_value"`
	CGSTRate     decimal.Decimal `json:"cgst_rate"`
	CGSTAmount   decimal.Decimal `json:"cgst_amount"`
	SGSTRate     decimal.Decimal `json:"sgst_rate"`
	SGSTAmount   decimal.Decimal `json:"sgst_amount"`
	IGSTRate     decimal.Decimal `json:"igst_rate"`
	IGSTAmount   decimal.Decimal `json:"igst_amount"`
	CessRate     decimal.Decimal `json:"cess_rate"`
	CessAmount   decimal.Decimal `json:"cess_amount"`
}

type GenerateIRNResponse struct {
	ID            uuid.UUID `json:"id"`
	IRN           string    `json:"irn"`
	AckNo         string    `json:"ack_no"`
	Status        IRNStatus `json:"status"`
	SignedInvoice string    `json:"signed_invoice"`
	SignedQRCode  string    `json:"signed_qr_code"`
	GeneratedAt   time.Time `json:"generated_at"`
}

type CancelIRNRequest struct {
	Reason string `json:"reason"`
	Remark string `json:"remark"`
}

type CancelIRNResponse struct {
	ID          uuid.UUID `json:"id"`
	IRN         string    `json:"irn"`
	Status      IRNStatus `json:"status"`
	CancelledAt time.Time `json:"cancelled_at"`
}

type ListEInvoicesRequest struct {
	GSTIN      string `json:"gstin"`
	Status     string `json:"status,omitempty"`
	FromDate   string `json:"from_date,omitempty"`
	ToDate     string `json:"to_date,omitempty"`
	PageSize   int    `json:"page_size"`
	PageOffset int    `json:"page_offset"`
}

type ListEInvoicesResponse struct {
	Invoices   []EInvoice `json:"invoices"`
	TotalCount int        `json:"total_count"`
	PageSize   int        `json:"page_size"`
	PageOffset int        `json:"page_offset"`
}

type EInvoiceSummary struct {
	TotalCount     int             `json:"total_count"`
	GeneratedCount int             `json:"generated_count"`
	PendingCount   int             `json:"pending_count"`
	CancelledCount int             `json:"cancelled_count"`
	FailedCount    int             `json:"failed_count"`
	TotalValue     decimal.Decimal `json:"total_value"`
}
