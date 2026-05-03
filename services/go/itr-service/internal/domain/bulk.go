package domain

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BulkBatchStatus string

const (
	BatchPending    BulkBatchStatus = "PENDING"
	BatchProcessing BulkBatchStatus = "PROCESSING"
	BatchCompleted  BulkBatchStatus = "COMPLETED"
	BatchFailed     BulkBatchStatus = "FAILED"
)

type EmployeeFilingStatus string

const (
	EmpPendingReview EmployeeFilingStatus = "PENDING_REVIEW"
	EmpApproved      EmployeeFilingStatus = "APPROVED"
	EmpRejected      EmployeeFilingStatus = "REJECTED"
	EmpSubmitted     EmployeeFilingStatus = "SUBMITTED"
	EmpEVerified     EmployeeFilingStatus = "E_VERIFIED"
	EmpMismatch      EmployeeFilingStatus = "MISMATCH"
)

type BulkFilingBatch struct {
	ID             uuid.UUID       `json:"id"`
	TenantID       uuid.UUID       `json:"tenant_id"`
	TaxYear        string          `json:"tax_year"`
	EmployerTAN    string          `json:"employer_tan"`
	EmployerName   string          `json:"employer_name"`
	TotalEmployees int             `json:"total_employees"`
	Processed      int             `json:"processed"`
	Ready          int             `json:"ready"`
	WithMismatches int             `json:"with_mismatches"`
	Status         BulkBatchStatus `json:"status"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type BulkFilingEmployee struct {
	ID             uuid.UUID            `json:"id"`
	TenantID       uuid.UUID            `json:"tenant_id"`
	BatchID        uuid.UUID            `json:"batch_id"`
	PAN            string               `json:"pan"`
	Name           string               `json:"name"`
	Email          string               `json:"email"`
	GrossSalary    decimal.Decimal      `json:"gross_salary"`
	TDSDeducted    decimal.Decimal      `json:"tds_deducted"`
	FormType       ITRFormType          `json:"form_type"`
	FilingID       *uuid.UUID           `json:"filing_id,omitempty"`
	Status         EmployeeFilingStatus `json:"status"`
	MismatchCount  int                  `json:"mismatch_count"`
	MagicLinkToken string               `json:"magic_link_token,omitempty"`
	TokenExpiresAt *time.Time           `json:"token_expires_at,omitempty"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
}

type MagicLinkToken struct {
	ID        uuid.UUID  `json:"id"`
	TenantID  uuid.UUID  `json:"tenant_id"`
	Token     string     `json:"token"`
	PAN       string     `json:"pan"`
	BatchID   uuid.UUID  `json:"batch_id"`
	FilingID  uuid.UUID  `json:"filing_id"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

const magicLinkTokenBytes = 32
const magicLinkTTL = 7 * 24 * time.Hour
const maxBulkBatchSize = 1000

func GenerateMagicLinkToken() (string, error) {
	b := make([]byte, magicLinkTokenBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func MagicLinkExpiry() time.Time {
	return time.Now().Add(magicLinkTTL)
}

func MaxBulkBatchSize() int {
	return maxBulkBatchSize
}

type BulkProcessInput struct {
	PAN         string          `json:"pan"`
	Name        string          `json:"name"`
	Email       string          `json:"email"`
	GrossSalary decimal.Decimal `json:"gross_salary"`
	TDSDeducted decimal.Decimal `json:"tds_deducted"`
}

type BulkProcessResult struct {
	PAN            string               `json:"pan"`
	Name           string               `json:"name"`
	FormType       ITRFormType          `json:"form_type"`
	Status         EmployeeFilingStatus `json:"status"`
	MismatchCount  int                  `json:"mismatch_count"`
	TaxComputation *TaxComputeResult    `json:"tax_computation,omitempty"`
	Reconciliation *AISReconcileResult  `json:"reconciliation,omitempty"`
}

func DetermineFormType(assesseeType AssesseeType, residency ResidencyStatus, totalIncome decimal.Decimal, hasBusiness bool, hasCapGains bool) ITRFormType {
	if hasBusiness {
		return FormITR3
	}
	if hasCapGains {
		return FormITR2
	}
	elig := CheckITR1Eligibility(assesseeType, residency, totalIncome, 0, false, zero, false, false, false, false, zero)
	if elig.Eligible {
		return FormITR1
	}
	return FormITR2
}

func ProcessEmployeeForBulkFiling(emp BulkProcessInput, aisData AISSourceData) BulkProcessResult {
	salaryResult := ComputeSalaryIncome(SalaryInput{GrossSalary: emp.GrossSalary})

	income := IncomeBreakdown{Salary: salaryResult.NetSalary}
	taxResult := ComputeTax(TaxComputeInput{
		Income:     income,
		Regime:     NewRegime,
		IsResident: true,
		TDSCredits: emp.TDSDeducted,
	})

	books := BookData{
		SalaryIncome: emp.GrossSalary,
		TDSClaims: []TDSCreditEntry{
			{DeductorTAN: "EMPLOYER", Section: "392", TDSAmount: emp.TDSDeducted},
		},
	}
	recon := ReconcileAIS(aisData, books, true)

	formType := DetermineFormType(AssesseeIndividual, Resident, income.Total(), false, false)

	status := EmpPendingReview
	if recon.HasErrors {
		status = EmpMismatch
	}

	return BulkProcessResult{
		PAN:            emp.PAN,
		Name:           emp.Name,
		FormType:       formType,
		Status:         status,
		MismatchCount:  recon.ErrorCount + recon.WarnCount,
		TaxComputation: &taxResult,
		Reconciliation: &recon,
	}
}

type BulkBatchSummary struct {
	TotalEmployees int `json:"total_employees"`
	Ready          int `json:"ready"`
	WithMismatches int `json:"with_mismatches"`
	Submitted      int `json:"submitted"`
	EVerified      int `json:"e_verified"`
	Rejected       int `json:"rejected"`
}
