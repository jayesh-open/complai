package domain

type PANVerifyRequest struct {
	PAN      string `json:"pan"`
	Name     string `json:"name"`
	TenantID string `json:"-"`
}

type PANVerifyResponse struct {
	PAN      string `json:"pan"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Category string `json:"category"`
}

type TANVerifyRequest struct {
	TAN      string `json:"tan"`
	TenantID string `json:"-"`
}

type TANVerifyResponse struct {
	TAN    string `json:"tan"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type ChallanRequest struct {
	TenantID       string  `json:"tenant_id"`
	TAN            string  `json:"tan"`
	Section        string  `json:"section"`
	Amount         float64 `json:"amount"`
	Surcharge      float64 `json:"surcharge"`
	Cess           float64 `json:"cess"`
	Interest       float64 `json:"interest"`
	Penalty        float64 `json:"penalty"`
	AssessmentYear string  `json:"assessment_year"`
}

type ChallanResponse struct {
	ChallanNumber string  `json:"challan_number"`
	BSRCode       string  `json:"bsr_code"`
	DepositDate   string  `json:"deposit_date"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
}

type Form26QRequest struct {
	TenantID      string         `json:"tenant_id"`
	TAN           string         `json:"tan"`
	FinancialYear string         `json:"financial_year"`
	Quarter       string         `json:"quarter"`
	Deductions    []Deduction26Q `json:"deductions"`
}

type Deduction26Q struct {
	DeducteePAN   string  `json:"deductee_pan"`
	DeducteeName  string  `json:"deductee_name"`
	Section       string  `json:"section"`
	PaymentDate   string  `json:"payment_date"`
	Amount        float64 `json:"amount"`
	TDSAmount     float64 `json:"tds_amount"`
	ChallanNumber string  `json:"challan_number"`
}

type Form24QRequest struct {
	TenantID      string       `json:"tenant_id"`
	TAN           string       `json:"tan"`
	FinancialYear string       `json:"financial_year"`
	Quarter       string       `json:"quarter"`
	Employees     []Employee24Q `json:"employees"`
}

type Employee24Q struct {
	PAN          string  `json:"pan"`
	Name         string  `json:"name"`
	Designation  string  `json:"designation"`
	GrossSalary  float64 `json:"gross_salary"`
	TDSDeducted  float64 `json:"tds_deducted"`
	TDSDeposited float64 `json:"tds_deposited"`
}

type FormFilingResponse struct {
	TokenNumber           string   `json:"token_number"`
	AcknowledgementNumber string   `json:"acknowledgement_number"`
	FilingDate            string   `json:"filing_date"`
	Status                string   `json:"status"`
	Errors                []string `json:"errors,omitempty"`
}
