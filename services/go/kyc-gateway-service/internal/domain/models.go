package domain

type PANVerifyRequest struct {
	PAN       string `json:"pan"`
	Name      string `json:"name"`
	RequestID string `json:"request_id"`
}

type PANVerifyResponse struct {
	PAN       string `json:"pan"`
	Name      string `json:"name"`
	Category  string `json:"category"` // Individual, Company, HUF, etc.
	Status    string `json:"status"`   // valid, invalid
	Valid     bool   `json:"valid"`
	RequestID string `json:"request_id"`
}

type GSTINVerifyRequest struct {
	GSTIN     string `json:"gstin"`
	RequestID string `json:"request_id"`
}

type GSTINVerifyResponse struct {
	GSTIN            string `json:"gstin"`
	LegalName        string `json:"legal_name"`
	TradeName        string `json:"trade_name"`
	Status           string `json:"status"` // Active, Cancelled, Suspended
	RegistrationType string `json:"registration_type"`
	StateCode        string `json:"state_code"`
	State            string `json:"state"`
	PAN              string `json:"pan"`
	Valid            bool   `json:"valid"`
	RequestID        string `json:"request_id"`
}

type TANVerifyRequest struct {
	TAN       string `json:"tan"`
	RequestID string `json:"request_id"`
}

type TANVerifyResponse struct {
	TAN       string `json:"tan"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	Valid     bool   `json:"valid"`
	RequestID string `json:"request_id"`
}

type BankVerifyRequest struct {
	AccountNumber string `json:"account_number"`
	IFSC          string `json:"ifsc"`
	RequestID     string `json:"request_id"`
}

type BankVerifyResponse struct {
	AccountNumber string `json:"account_number"`
	IFSC          string `json:"ifsc"`
	BankName      string `json:"bank_name"`
	BranchName    string `json:"branch_name"`
	NameAtBank    string `json:"name_at_bank"`
	Valid         bool   `json:"valid"`
	RequestID     string `json:"request_id"`
}

type GatewayResponse struct {
	Data interface{}  `json:"data"`
	Meta ResponseMeta `json:"meta"`
}

type ResponseMeta struct {
	RequestID      string `json:"request_id"`
	LatencyMs      int    `json:"latency_ms"`
	ProviderStatus string `json:"provider_status"`
}
