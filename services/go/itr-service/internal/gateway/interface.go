package gateway

import (
	"context"

	"github.com/shopspring/decimal"
)

type ITRGatewayProvider interface {
	CheckPANAadhaarLink(ctx context.Context, pan string) (*PANAadhaarStatus, error)
	FetchAIS(ctx context.Context, pan, taxYear string) (*AISData, error)
	SubmitITR(ctx context.Context, req ITRSubmitRequest) (*ITRSubmitResponse, error)
	GenerateITRV(ctx context.Context, arn string) (*ITRVResponse, error)
	CheckEVerification(ctx context.Context, arn string) (*EVerifyStatus, error)
	CheckRefundStatus(ctx context.Context, pan, taxYear string) (*RefundStatus, error)
}

type PANAadhaarStatus struct {
	PAN       string `json:"pan"`
	Linked    bool   `json:"linked"`
	LinkDate  string `json:"link_date,omitempty"`
	AadhaarLast4 string `json:"aadhaar_last4,omitempty"`
}

type AISData struct {
	PAN      string          `json:"pan"`
	TaxYear  string          `json:"tax_year"`
	Form168Ref string        `json:"form_168_ref"`
	TDSEntries []AISTDSEntry `json:"tds_entries"`
	InterestIncome decimal.Decimal `json:"interest_income"`
	DividendIncome decimal.Decimal `json:"dividend_income"`
	SalaryIncome   decimal.Decimal `json:"salary_income"`
	SecuritiesTrading decimal.Decimal `json:"securities_trading"`
}

type AISTDSEntry struct {
	DeductorTAN  string          `json:"deductor_tan"`
	DeductorName string          `json:"deductor_name"`
	Section      string          `json:"section"`
	Amount       decimal.Decimal `json:"amount"`
	TDSAmount    decimal.Decimal `json:"tds_amount"`
	Quarter      string          `json:"quarter"`
}

type ITRSubmitRequest struct {
	PAN      string `json:"pan"`
	TaxYear  string `json:"tax_year"`
	FormType string `json:"form_type"`
	Payload  string `json:"payload"`
}

type ITRSubmitResponse struct {
	ARN               string `json:"arn"`
	AcknowledgementNo string `json:"acknowledgement_no"`
	FilingDate        string `json:"filing_date"`
	Status            string `json:"status"`
}

type ITRVResponse struct {
	ARN              string `json:"arn"`
	ITRVURL          string `json:"itrv_url"`
	AcknowledgementNo string `json:"acknowledgement_no"`
}

type EVerifyStatus struct {
	ARN      string `json:"arn"`
	Verified bool   `json:"verified"`
	Method   string `json:"method,omitempty"`
	Date     string `json:"date,omitempty"`
}

type RefundStatus struct {
	PAN       string          `json:"pan"`
	TaxYear   string          `json:"tax_year"`
	Status    string          `json:"status"`
	Amount    decimal.Decimal `json:"amount"`
	BankRef   string          `json:"bank_ref,omitempty"`
	CreditDate string         `json:"credit_date,omitempty"`
}
