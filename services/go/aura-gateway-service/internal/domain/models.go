package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Invoice struct {
	ID             uuid.UUID       `json:"id"`
	TenantID       uuid.UUID       `json:"tenant_id"`
	DocumentNumber string          `json:"document_number"`
	DocumentDate   string          `json:"document_date"` // DD/MM/YYYY
	DocumentType   string          `json:"document_type"` // INV, CRN, DBN
	SupplyType     string          `json:"supply_type"`   // B2B, B2CL, B2CS, EXP
	ReverseCharge  bool            `json:"reverse_charge"`
	Supplier       Party           `json:"supplier"`
	Buyer          Party           `json:"buyer"`
	LineItems      []LineItem      `json:"line_items"`
	Totals         InvoiceTotals   `json:"totals"`
	PlaceOfSupply  string          `json:"place_of_supply"` // state code
	SourceSystem   string          `json:"source_system"`
	CreatedAt      time.Time       `json:"created_at"`
}

type Party struct {
	GSTIN     string `json:"gstin"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	StateCode string `json:"state_code"`
}

type LineItem struct {
	ItemID       string          `json:"item_id"`
	Description  string          `json:"description"`
	HSN          string          `json:"hsn"`
	Unit         string          `json:"unit"`
	Quantity     decimal.Decimal `json:"quantity"`
	UnitPrice    decimal.Decimal `json:"unit_price"`
	Discount     decimal.Decimal `json:"discount"`
	TaxableValue decimal.Decimal `json:"taxable_value"`
	CGSTRate     decimal.Decimal `json:"cgst_rate"`
	CGSTAmount   decimal.Decimal `json:"cgst_amount"`
	SGSTRate     decimal.Decimal `json:"sgst_rate"`
	SGSTAmount   decimal.Decimal `json:"sgst_amount"`
	IGSTRate     decimal.Decimal `json:"igst_rate"`
	IGSTAmount   decimal.Decimal `json:"igst_amount"`
}

type InvoiceTotals struct {
	TaxableValue decimal.Decimal `json:"taxable_value"`
	CGST         decimal.Decimal `json:"cgst"`
	SGST         decimal.Decimal `json:"sgst"`
	IGST         decimal.Decimal `json:"igst"`
	GrandTotal   decimal.Decimal `json:"grand_total"`
}

type InvoiceListResponse struct {
	Invoices   []Invoice `json:"invoices"`
	TotalCount int       `json:"total_count"`
	Summary    InvoiceSummary `json:"summary"`
}

type InvoiceSummary struct {
	B2BIntraCount int `json:"b2b_intra_count"`
	B2BInterCount int `json:"b2b_inter_count"`
	B2CLCount     int `json:"b2cl_count"`
	B2CSCount     int `json:"b2cs_count"`
	ExportCount   int `json:"export_count"`
	RCMCount      int `json:"rcm_count"`
	CreditNote    int `json:"credit_note_count"`
	DebitNote     int `json:"debit_note_count"`
}
