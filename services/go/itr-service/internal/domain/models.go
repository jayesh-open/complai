package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type RegimeType string

const (
	NewRegime RegimeType = "NEW_REGIME"
	OldRegime RegimeType = "OLD_REGIME"
)

type ITRFormType string

const (
	FormITR1 ITRFormType = "ITR-1"
	FormITR2 ITRFormType = "ITR-2"
	FormITR3 ITRFormType = "ITR-3"
)

type FilingStatus string

const (
	StatusDraft     FilingStatus = "DRAFT"
	StatusValidated FilingStatus = "VALIDATED"
	StatusSubmitted FilingStatus = "SUBMITTED"
	StatusVerified  FilingStatus = "VERIFIED"
	StatusFiled     FilingStatus = "FILED"
	StatusProcessed FilingStatus = "PROCESSED"
	StatusDefective FilingStatus = "DEFECTIVE"
	StatusRejected  FilingStatus = "REJECTED"
)

type IncomeHead string

const (
	HeadSalary           IncomeHead = "SALARY"
	HeadHouseProperty    IncomeHead = "HOUSE_PROPERTY"
	HeadCapitalGains     IncomeHead = "CAPITAL_GAINS"
	HeadBusinessProf     IncomeHead = "BUSINESS_PROFESSION"
	HeadOtherSources     IncomeHead = "OTHER_SOURCES"
)

type ResidencyStatus string

const (
	Resident          ResidencyStatus = "RESIDENT"
	NonResident       ResidencyStatus = "NON_RESIDENT"
	ResidentNotOrdRes ResidencyStatus = "RNOR"
)

type AssesseeType string

const (
	AssesseeIndividual AssesseeType = "INDIVIDUAL"
	AssesseeHUF        AssesseeType = "HUF"
)

type CapitalGainType string

const (
	LTCG CapitalGainType = "LTCG"
	STCG CapitalGainType = "STCG"
)

type VerificationMethod string

const (
	VerifyAadhaarOTP VerificationMethod = "AADHAAR_OTP"
	VerifyEVC        VerificationMethod = "EVC"
	VerifyNetBanking VerificationMethod = "NET_BANKING"
	VerifyDSC        VerificationMethod = "DSC"
)

type Taxpayer struct {
	ID              uuid.UUID       `json:"id"`
	TenantID        uuid.UUID       `json:"tenant_id"`
	PAN             string          `json:"pan"`
	Name            string          `json:"name"`
	DateOfBirth     time.Time       `json:"date_of_birth"`
	AssesseeType    AssesseeType    `json:"assessee_type"`
	ResidencyStatus ResidencyStatus `json:"residency_status"`
	AadhaarLinked   bool            `json:"aadhaar_linked"`
	Email           string          `json:"email"`
	Mobile          string          `json:"mobile"`
	Address         string          `json:"address,omitempty"`
	EmployerTAN     string          `json:"employer_tan,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type ITRFiling struct {
	ID                    uuid.UUID          `json:"id"`
	TenantID              uuid.UUID          `json:"tenant_id"`
	TaxpayerID            uuid.UUID          `json:"taxpayer_id"`
	PAN                   string             `json:"pan"`
	TaxYear               string             `json:"tax_year"`
	FormType              ITRFormType        `json:"form_type"`
	RegimeSelected        RegimeType         `json:"regime_selected"`
	Form10IEARef          string             `json:"form_10iea_ref,omitempty"`
	Status                FilingStatus       `json:"status"`
	GrossIncome           decimal.Decimal    `json:"gross_income"`
	TotalDeductions       decimal.Decimal    `json:"total_deductions"`
	TaxableIncome         decimal.Decimal    `json:"taxable_income"`
	TaxPayable            decimal.Decimal    `json:"tax_payable"`
	TDSCredited           decimal.Decimal    `json:"tds_credited"`
	AdvanceTaxPaid        decimal.Decimal    `json:"advance_tax_paid"`
	SelfAssessmentTax     decimal.Decimal    `json:"self_assessment_tax"`
	RefundDue             decimal.Decimal    `json:"refund_due"`
	BalancePayable        decimal.Decimal    `json:"balance_payable"`
	VerificationMethod    VerificationMethod `json:"verification_method,omitempty"`
	ARN                   string             `json:"arn,omitempty"`
	AcknowledgementNumber string             `json:"acknowledgement_number,omitempty"`
	FiledAt               *time.Time         `json:"filed_at,omitempty"`
	IdempotencyKey        string             `json:"idempotency_key"`
	ErrorMessage          string             `json:"error_message,omitempty"`
	CreatedAt             time.Time          `json:"created_at"`
	UpdatedAt             time.Time          `json:"updated_at"`
}

type IncomeEntry struct {
	ID          uuid.UUID       `json:"id"`
	TenantID    uuid.UUID       `json:"tenant_id"`
	FilingID    uuid.UUID       `json:"filing_id"`
	Head        IncomeHead      `json:"head"`
	SubHead     string          `json:"sub_head,omitempty"`
	Section     string          `json:"section,omitempty"`
	Description string          `json:"description"`
	Amount      decimal.Decimal `json:"amount"`
	Exempt      bool            `json:"exempt"`
	CreatedAt   time.Time       `json:"created_at"`
}

type Deduction struct {
	ID        uuid.UUID       `json:"id"`
	TenantID  uuid.UUID       `json:"tenant_id"`
	FilingID  uuid.UUID       `json:"filing_id"`
	Section   string          `json:"section"`
	Label     string          `json:"label"`
	Claimed   decimal.Decimal `json:"claimed"`
	Allowed   decimal.Decimal `json:"allowed"`
	MaxLimit  decimal.Decimal `json:"max_limit"`
	CreatedAt time.Time       `json:"created_at"`
}

type TaxComputation struct {
	ID                uuid.UUID       `json:"id"`
	TenantID          uuid.UUID       `json:"tenant_id"`
	FilingID          uuid.UUID       `json:"filing_id"`
	RegimeType        RegimeType      `json:"regime_type"`
	GrossIncome       decimal.Decimal `json:"gross_income"`
	StandardDeduction decimal.Decimal `json:"standard_deduction"`
	TotalDeductions   decimal.Decimal `json:"total_deductions"`
	TaxableIncome     decimal.Decimal `json:"taxable_income"`
	BaseTax           decimal.Decimal `json:"base_tax"`
	Surcharge         decimal.Decimal `json:"surcharge"`
	SurchargeRate     decimal.Decimal `json:"surcharge_rate"`
	HealthEdCess      decimal.Decimal `json:"health_ed_cess"`
	Rebate87A         decimal.Decimal `json:"rebate_87a"`
	GrossTaxPayable   decimal.Decimal `json:"gross_tax_payable"`
	TDSCredit         decimal.Decimal `json:"tds_credit"`
	AdvanceTax        decimal.Decimal `json:"advance_tax"`
	SelfAssessmentTax decimal.Decimal `json:"self_assessment_tax"`
	NetTaxPayable     decimal.Decimal `json:"net_tax_payable"`
	RefundDue         decimal.Decimal `json:"refund_due"`
	CreatedAt         time.Time       `json:"created_at"`
}

type TDSCredit struct {
	ID              uuid.UUID       `json:"id"`
	TenantID        uuid.UUID       `json:"tenant_id"`
	FilingID        uuid.UUID       `json:"filing_id"`
	DeductorTAN     string          `json:"deductor_tan"`
	DeductorName    string          `json:"deductor_name"`
	Section         string          `json:"section"`
	TDSAmount       decimal.Decimal `json:"tds_amount"`
	GrossPayment    decimal.Decimal `json:"gross_payment"`
	TaxYear         string          `json:"tax_year"`
	MatchedWithAIS  bool            `json:"matched_with_ais"`
	AISAmount       decimal.Decimal `json:"ais_amount"`
	Discrepancy     decimal.Decimal `json:"discrepancy"`
	DiscrepancyNote string          `json:"discrepancy_note,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

type AISReconciliation struct {
	ID             uuid.UUID       `json:"id"`
	TenantID       uuid.UUID       `json:"tenant_id"`
	FilingID       uuid.UUID       `json:"filing_id"`
	PAN            string          `json:"pan"`
	TaxYear        string          `json:"tax_year"`
	SourceType     string          `json:"source_type"`
	ReportedAmount decimal.Decimal `json:"reported_amount"`
	AISAmount      decimal.Decimal `json:"ais_amount"`
	Discrepancy    decimal.Decimal `json:"discrepancy"`
	Status         string          `json:"status"`
	Notes          string          `json:"notes,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
}

var oldSectionAliases = map[string]bool{
	"115BAC":  true,
	"80C":     true,
	"80D":     true,
	"80E":     true,
	"80G":     true,
	"80GG":    true,
	"80TTA":   true,
	"80TTB":   true,
	"80CCD":   true,
	"10(13A)": true,
	"24(b)":   true,
}

func IsOldSectionRef(section string) bool {
	return oldSectionAliases[section]
}

var ita2025SectionMapping = map[string]string{
	"115BAC": "202 (New Tax Regime)",
	"80C":    "Schedule VI-A (old regime only via Form 10-IEA)",
	"80D":    "Schedule VI-A (old regime only via Form 10-IEA)",
}

func ITA2025Equivalent(oldSection string) string {
	if eq, ok := ita2025SectionMapping[oldSection]; ok {
		return eq
	}
	return ""
}
