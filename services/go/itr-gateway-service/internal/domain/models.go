package domain

type PANAadhaarLinkRequest struct {
	PAN      string `json:"pan"`
	TenantID string `json:"-"`
}

type PANAadhaarLinkResponse struct {
	PAN          string `json:"pan"`
	Linked       bool   `json:"linked"`
	LinkDate     string `json:"link_date,omitempty"`
	AadhaarLast4 string `json:"aadhaar_last4,omitempty"`
}

type AISRequest struct {
	PAN      string `json:"pan"`
	TaxYear  string `json:"tax_year"`
	TenantID string `json:"-"`
}

type AISTDSEntry struct {
	DeductorTAN  string  `json:"deductor_tan"`
	DeductorName string  `json:"deductor_name"`
	Section      string  `json:"section"`
	Amount       float64 `json:"amount"`
	TDSAmount    float64 `json:"tds_amount"`
	Quarter      string  `json:"quarter"`
}

type AISResponse struct {
	PAN               string        `json:"pan"`
	TaxYear           string        `json:"tax_year"`
	Form168Ref        string        `json:"form_168_ref"`
	TDSEntries        []AISTDSEntry `json:"tds_entries"`
	InterestIncome    float64       `json:"interest_income"`
	DividendIncome    float64       `json:"dividend_income"`
	SalaryIncome      float64       `json:"salary_income"`
	SecuritiesTrading float64       `json:"securities_trading"`
}

type ITRSubmitRequest struct {
	PAN      string `json:"pan"`
	TaxYear  string `json:"tax_year"`
	FormType string `json:"form_type"`
	Payload  string `json:"payload"`
	TenantID string `json:"-"`
}

type ITRSubmitResponse struct {
	ARN               string `json:"arn"`
	AcknowledgementNo string `json:"acknowledgement_no"`
	FilingDate        string `json:"filing_date"`
	Status            string `json:"status"`
}

type ITRVRequest struct {
	ARN      string `json:"arn"`
	TenantID string `json:"-"`
}

type ITRVResponse struct {
	ARN               string `json:"arn"`
	ITRVURL           string `json:"itrv_url"`
	AcknowledgementNo string `json:"acknowledgement_no"`
}

type EVerifyRequest struct {
	ARN      string `json:"arn"`
	Method   string `json:"method"`
	TenantID string `json:"-"`
}

type EVerifyResponse struct {
	ARN      string `json:"arn"`
	Verified bool   `json:"verified"`
	Method   string `json:"method,omitempty"`
	Date     string `json:"date,omitempty"`
}

type RefundStatusRequest struct {
	PAN      string `json:"pan"`
	TaxYear  string `json:"tax_year"`
	TenantID string `json:"-"`
}

type RefundStatusResponse struct {
	PAN        string  `json:"pan"`
	TaxYear    string  `json:"tax_year"`
	Status     string  `json:"status"`
	Amount     float64 `json:"amount"`
	BankRef    string  `json:"bank_ref,omitempty"`
	CreditDate string  `json:"credit_date,omitempty"`
}
