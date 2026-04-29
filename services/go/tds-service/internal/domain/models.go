package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Section string

const (
	Section392   Section = "392"
	Section393_1 Section = "393(1)"
	Section393_2 Section = "393(2)"
	Section393_3 Section = "393(3)"
)

func ValidSection(s Section) bool {
	switch s {
	case Section392, Section393_1, Section393_2, Section393_3:
		return true
	}
	return false
}

type PaymentCode string

const (
	CodeSalaryState     PaymentCode = "1001"
	CodeSalaryPrivate   PaymentCode = "1002"
	CodeSalaryCentral   PaymentCode = "1003"
	CodeEPFWithdrawal   PaymentCode = "1004"
	CodeRentPlant       PaymentCode = "1008"
	CodeRentLand        PaymentCode = "1009"
	CodeContractorIndiv PaymentCode = "1023"
	CodeContractorOther PaymentCode = "1024"
	CodeTechnical       PaymentCode = "1026"
	CodeProfessional    PaymentCode = "1027"
	CodeDirectorRem     PaymentCode = "1028"
	CodePurchaseGoods   PaymentCode = "1031"
	CodeNonResident     PaymentCode = "1057"
)

var paymentCodeSection = map[PaymentCode]Section{
	CodeSalaryState:     Section392,
	CodeSalaryPrivate:   Section392,
	CodeSalaryCentral:   Section392,
	CodeEPFWithdrawal:   Section392,
	CodeRentPlant:       Section393_1,
	CodeRentLand:        Section393_1,
	CodeContractorIndiv: Section393_1,
	CodeContractorOther: Section393_1,
	CodeTechnical:       Section393_1,
	CodeProfessional:    Section393_1,
	CodeDirectorRem:     Section393_1,
	CodePurchaseGoods:   Section393_1,
	CodeNonResident:     Section393_2,
}

var paymentCodeSubClause = map[PaymentCode]string{
	CodeRentPlant:       "Sl.2(ii).D(a)",
	CodeRentLand:        "Sl.2(ii).D(b)",
	CodeContractorIndiv: "Sl.6(i).D(a)",
	CodeContractorOther: "Sl.6(i).D(b)",
	CodeTechnical:       "Sl.6(iii).D(a)",
	CodeProfessional:    "Sl.6(iii).D(b)",
	CodeDirectorRem:     "Sl.6(iii).D(b)",
	CodePurchaseGoods:   "Sl.8(ii)",
	CodeNonResident:     "Sl.17",
}

func SectionForCode(code PaymentCode) Section {
	if s, ok := paymentCodeSection[code]; ok {
		return s
	}
	return ""
}

func SubClauseForCode(code PaymentCode) string {
	return paymentCodeSubClause[code]
}

func ValidPaymentCode(code PaymentCode) bool {
	_, ok := paymentCodeSection[code]
	return ok
}

type DeducteeType string

const (
	DeducteeIndividual DeducteeType = "INDIVIDUAL"
	DeducteeHUF        DeducteeType = "HUF"
	DeducteeCompany    DeducteeType = "COMPANY"
	DeducteeFirm       DeducteeType = "FIRM"
	DeducteeTrust      DeducteeType = "TRUST"
	DeducteeAOP        DeducteeType = "AOP"
	DeducteeLocalAuth  DeducteeType = "LOCAL_AUTHORITY"
	DeducteeGovernment DeducteeType = "GOVERNMENT"
)

type ResidentStatus string

const (
	Resident    ResidentStatus = "RESIDENT"
	NonResident ResidentStatus = "NON_RESIDENT"
)

type RentType string

const (
	RentLandBuilding   RentType = "LAND_BUILDING"
	RentPlantMachinery RentType = "PLANT_MACHINERY"
)

type EntryStatus string

const (
	StatusPending   EntryStatus = "PENDING"
	StatusDeposited EntryStatus = "DEPOSITED"
	StatusFiled     EntryStatus = "FILED"
	StatusRevised   EntryStatus = "REVISED"
)

type Deductee struct {
	ID                       uuid.UUID        `json:"id"`
	TenantID                 uuid.UUID        `json:"tenant_id"`
	VendorID                 uuid.UUID        `json:"vendor_id"`
	Name                     string           `json:"name"`
	PAN                      string           `json:"pan"`
	PANVerified              bool             `json:"pan_verified"`
	PANStatus                string           `json:"pan_status"`
	DeducteeType             DeducteeType     `json:"deductee_type"`
	ResidentStatus           ResidentStatus   `json:"resident_status"`
	SectionOverride          string           `json:"section_override,omitempty"`
	LowerDeductionCert       string           `json:"lower_deduction_cert,omitempty"`
	LowerDeductionRate       *decimal.Decimal `json:"lower_deduction_rate,omitempty"`
	LowerDeductionValidUntil *time.Time       `json:"lower_deduction_valid_until,omitempty"`
	CreatedAt                time.Time        `json:"created_at"`
	UpdatedAt                time.Time        `json:"updated_at"`
}

type TDSEntry struct {
	ID               uuid.UUID       `json:"id"`
	TenantID         uuid.UUID       `json:"tenant_id"`
	DeducteeID       uuid.UUID       `json:"deductee_id"`
	DeducteeName     string          `json:"deductee_name,omitempty"`
	Section          Section         `json:"section"`
	PaymentCode      PaymentCode     `json:"payment_code"`
	SubClause        string          `json:"sub_clause,omitempty"`
	FinancialYear    string          `json:"financial_year"`
	TaxYear          string          `json:"tax_year"`
	Quarter          string          `json:"quarter"`
	TransactionDate  time.Time       `json:"transaction_date"`
	PaymentDate      *time.Time      `json:"payment_date,omitempty"`
	GrossAmount      decimal.Decimal `json:"gross_amount"`
	TDSRate          decimal.Decimal `json:"tds_rate"`
	TDSAmount        decimal.Decimal `json:"tds_amount"`
	Surcharge        decimal.Decimal `json:"surcharge"`
	Cess             decimal.Decimal `json:"cess"`
	TotalTax         decimal.Decimal `json:"total_tax"`
	InvoiceNumber    string          `json:"invoice_number,omitempty"`
	NatureOfPayment  string          `json:"nature_of_payment"`
	PANAtDeduction   string          `json:"pan_at_deduction"`
	NoPANDeduction   bool            `json:"no_pan_deduction"`
	LowerCertApplied bool            `json:"lower_cert_applied"`
	ChallanNumber    string          `json:"challan_number,omitempty"`
	ChallanDate      *time.Time      `json:"challan_date,omitempty"`
	BSRCode          string          `json:"bsr_code,omitempty"`
	Status           EntryStatus     `json:"status"`
	Remarks          string          `json:"remarks,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

type TDSAggregate struct {
	ID               uuid.UUID       `json:"id"`
	TenantID         uuid.UUID       `json:"tenant_id"`
	DeducteeID       uuid.UUID       `json:"deductee_id"`
	PaymentCode      PaymentCode     `json:"payment_code"`
	FinancialYear    string          `json:"financial_year"`
	TotalPaid        decimal.Decimal `json:"total_paid"`
	TotalTDS         decimal.Decimal `json:"total_tds"`
	TransactionCount int             `json:"transaction_count"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

type TDSSummary struct {
	TotalDeductees       int                    `json:"total_deductees"`
	TotalEntries         int                    `json:"total_entries"`
	TotalTDSDeducted     decimal.Decimal        `json:"total_tds_deducted"`
	TotalTDSDeposited    decimal.Decimal        `json:"total_tds_deposited"`
	PendingDeposit       decimal.Decimal        `json:"pending_deposit"`
	EntriesByPaymentCode map[PaymentCode]int    `json:"entries_by_payment_code"`
	EntriesByStatus      map[EntryStatus]int    `json:"entries_by_status"`
}

func TaxYearFromFY(fy string) string {
	return fy
}
