package domain

// Vendor represents a vendor record from Apex P2P system.
type Vendor struct {
	ID                 string `json:"id"`
	TenantID           string `json:"tenant_id"`
	Name               string `json:"name"`
	LegalName          string `json:"legal_name"`
	TradeName          string `json:"trade_name"`
	PAN                string `json:"pan"`
	GSTIN              string `json:"gstin"`
	TAN                string `json:"tan,omitempty"`
	State              string `json:"state"`
	StateCode          string `json:"state_code"`
	Category           string `json:"category"`            // Manufacturer, Service Provider, Trader, Import Supplier, Logistics
	RegistrationStatus string `json:"registration_status"` // Active, Cancelled, Suspended
	MSMERegistered     bool   `json:"msme_registered"`
	Email              string `json:"email"`
	Phone              string `json:"phone"`
	Address            string `json:"address"`
	BankAccount        string `json:"bank_account,omitempty"`
	IFSC               string `json:"ifsc,omitempty"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

// APInvoice represents an accounts-payable invoice from Apex P2P system.
type APInvoice struct {
	ID              string  `json:"id"`
	TenantID        string  `json:"tenant_id"`
	VendorID        string  `json:"vendor_id"`
	VendorGSTIN     string  `json:"vendor_gstin"`
	InvoiceNumber   string  `json:"invoice_number"`
	InvoiceDate     string  `json:"invoice_date"`
	DueDate         string  `json:"due_date"`
	TaxableValue    float64 `json:"taxable_value"`
	CGSTAmount      float64 `json:"cgst_amount"`
	SGSTAmount      float64 `json:"sgst_amount"`
	IGSTAmount      float64 `json:"igst_amount"`
	TotalAmount     float64 `json:"total_amount"`
	PaymentStatus   string  `json:"payment_status"`   // paid, unpaid, overdue, partial
	PaymentDate     string  `json:"payment_date,omitempty"`
	IRNGenerated    bool    `json:"irn_generated"`
	IRN             string  `json:"irn,omitempty"`
	GSTFilingStatus string  `json:"gst_filing_status"` // filed, pending, late, not_filed
	MismatchStatus  string  `json:"mismatch_status"`   // matched, mismatched, pending
	CreatedAt       string  `json:"created_at"`
}

// FetchVendorsRequest is the request to fetch vendors from Apex.
type FetchVendorsRequest struct {
	TenantID  string `json:"tenant_id"`
	Limit     int    `json:"limit,omitempty"`
	Offset    int    `json:"offset,omitempty"`
	RequestID string `json:"request_id"`
}

// FetchVendorsResponse is the response containing vendor data.
type FetchVendorsResponse struct {
	Vendors   []Vendor `json:"vendors"`
	Total     int      `json:"total"`
	RequestID string   `json:"request_id"`
}

// FetchAPInvoicesRequest is the request to fetch AP invoices from Apex.
type FetchAPInvoicesRequest struct {
	TenantID  string `json:"tenant_id"`
	VendorID  string `json:"vendor_id,omitempty"`
	FromDate  string `json:"from_date,omitempty"`
	ToDate    string `json:"to_date,omitempty"`
	RequestID string `json:"request_id"`
}

// FetchAPInvoicesResponse is the response containing AP invoice data.
type FetchAPInvoicesResponse struct {
	Invoices  []APInvoice `json:"invoices"`
	Total     int         `json:"total"`
	RequestID string      `json:"request_id"`
}

// GatewayResponse wraps all gateway responses with metadata.
type GatewayResponse struct {
	Data interface{}  `json:"data"`
	Meta ResponseMeta `json:"meta"`
}

// ResponseMeta contains request metadata for gateway responses.
type ResponseMeta struct {
	RequestID      string `json:"request_id"`
	LatencyMs      int    `json:"latency_ms"`
	ProviderStatus string `json:"provider_status"`
}
