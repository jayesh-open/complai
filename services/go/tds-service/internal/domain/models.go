package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Section string

const (
	Section192  Section = "192"
	Section194C Section = "194C"
	Section194I Section = "194I"
	Section194J Section = "194J"
	Section194Q Section = "194Q"
	Section195  Section = "195"
)

func ValidSection(s Section) bool {
	switch s {
	case Section192, Section194C, Section194I, Section194J, Section194Q, Section195:
		return true
	}
	return false
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
	ID              uuid.UUID       `json:"id"`
	TenantID        uuid.UUID       `json:"tenant_id"`
	DeducteeID      uuid.UUID       `json:"deductee_id"`
	DeducteeName    string          `json:"deductee_name,omitempty"`
	Section         Section         `json:"section"`
	FinancialYear   string          `json:"financial_year"`
	Quarter         string          `json:"quarter"`
	TransactionDate time.Time       `json:"transaction_date"`
	PaymentDate     *time.Time      `json:"payment_date,omitempty"`
	GrossAmount     decimal.Decimal `json:"gross_amount"`
	TDSRate         decimal.Decimal `json:"tds_rate"`
	TDSAmount       decimal.Decimal `json:"tds_amount"`
	Surcharge       decimal.Decimal `json:"surcharge"`
	Cess            decimal.Decimal `json:"cess"`
	TotalTax        decimal.Decimal `json:"total_tax"`
	InvoiceNumber   string          `json:"invoice_number,omitempty"`
	NatureOfPayment string          `json:"nature_of_payment"`
	PANAtDeduction  string          `json:"pan_at_deduction"`
	NoPANDeduction  bool            `json:"no_pan_deduction"`
	LowerCertApplied bool          `json:"lower_cert_applied"`
	ChallanNumber   string          `json:"challan_number,omitempty"`
	ChallanDate     *time.Time      `json:"challan_date,omitempty"`
	BSRCode         string          `json:"bsr_code,omitempty"`
	Status          EntryStatus     `json:"status"`
	Remarks         string          `json:"remarks,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type TDSAggregate struct {
	ID               uuid.UUID       `json:"id"`
	TenantID         uuid.UUID       `json:"tenant_id"`
	DeducteeID       uuid.UUID       `json:"deductee_id"`
	Section          Section         `json:"section"`
	FinancialYear    string          `json:"financial_year"`
	TotalPaid        decimal.Decimal `json:"total_paid"`
	TotalTDS         decimal.Decimal `json:"total_tds"`
	TransactionCount int             `json:"transaction_count"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

type TDSSummary struct {
	TotalDeductees    int             `json:"total_deductees"`
	TotalEntries      int             `json:"total_entries"`
	TotalTDSDeducted  decimal.Decimal `json:"total_tds_deducted"`
	TotalTDSDeposited decimal.Decimal `json:"total_tds_deposited"`
	PendingDeposit    decimal.Decimal `json:"pending_deposit"`
	EntriesBySection  map[Section]int `json:"entries_by_section"`
	EntriesByStatus   map[EntryStatus]int `json:"entries_by_status"`
}
