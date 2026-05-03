package domain

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateMagicLinkToken(t *testing.T) {
	tok1, err := GenerateMagicLinkToken()
	require.NoError(t, err)
	assert.Len(t, tok1, 64) // 32 bytes → 64 hex chars

	tok2, err := GenerateMagicLinkToken()
	require.NoError(t, err)
	assert.NotEqual(t, tok1, tok2)
}

func TestMagicLinkExpiry(t *testing.T) {
	exp := MagicLinkExpiry()
	assert.True(t, exp.After(MagicLinkExpiry().Add(-7*24*60*60*1e9)))
}

func TestMaxBulkBatchSize(t *testing.T) {
	assert.Equal(t, 1000, MaxBulkBatchSize())
}

func TestDetermineFormType(t *testing.T) {
	tests := []struct {
		name        string
		assessee    AssesseeType
		residency   ResidencyStatus
		income      decimal.Decimal
		hasBusiness bool
		hasCapGains bool
		want        ITRFormType
	}{
		{"Business→ITR3", AssesseeIndividual, Resident, decimal.NewFromInt(300000), true, false, FormITR3},
		{"CapGains→ITR2", AssesseeIndividual, Resident, decimal.NewFromInt(300000), false, true, FormITR2},
		{"SimpleResident→ITR1", AssesseeIndividual, Resident, decimal.NewFromInt(300000), false, false, FormITR1},
		{"NRI→ITR2", AssesseeIndividual, NonResident, decimal.NewFromInt(300000), false, false, FormITR2},
		{"HUF→ITR2", AssesseeHUF, Resident, decimal.NewFromInt(300000), false, false, FormITR2},
		{"HighIncome→ITR2", AssesseeIndividual, Resident, decimal.NewFromInt(6000000), false, false, FormITR2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetermineFormType(tt.assessee, tt.residency, tt.income, tt.hasBusiness, tt.hasCapGains)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProcessEmployeeForBulkFiling_Clean(t *testing.T) {
	emp := BulkProcessInput{
		PAN:         "ABCDE1234F",
		Name:        "Test Employee",
		Email:       "test@example.com",
		GrossSalary: decimal.NewFromInt(800000),
		TDSDeducted: decimal.NewFromInt(50000),
	}
	ais := AISSourceData{
		PAN:          "ABCDE1234F",
		TaxYear:      "2026-27",
		SalaryIncome: decimal.NewFromInt(800000),
		TDSEntries: []AISEntry{
			{DeductorTAN: "EMPLOYER", Section: "392", TDSAmount: decimal.NewFromInt(50000)},
		},
	}

	result := ProcessEmployeeForBulkFiling(emp, ais)

	assert.Equal(t, "ABCDE1234F", result.PAN)
	assert.Equal(t, "Test Employee", result.Name)
	assert.Equal(t, FormITR1, result.FormType)
	assert.Equal(t, EmpPendingReview, result.Status)
	assert.NotNil(t, result.TaxComputation)
	assert.NotNil(t, result.Reconciliation)
}

func TestProcessEmployeeForBulkFiling_WithMismatches(t *testing.T) {
	emp := BulkProcessInput{
		PAN:         "XYZAB9876G",
		Name:        "Mismatch Emp",
		Email:       "mismatch@example.com",
		GrossSalary: decimal.NewFromInt(1200000),
		TDSDeducted: decimal.NewFromInt(100000),
	}
	ais := AISSourceData{
		PAN:            "XYZAB9876G",
		TaxYear:        "2026-27",
		SalaryIncome:   decimal.NewFromInt(1500000),
		InterestIncome: decimal.NewFromInt(50000),
		TDSEntries: []AISEntry{
			{DeductorTAN: "EMPLOYER", Section: "392", TDSAmount: decimal.NewFromInt(100000)},
		},
	}

	result := ProcessEmployeeForBulkFiling(emp, ais)

	assert.Equal(t, EmpMismatch, result.Status)
	assert.True(t, result.MismatchCount > 0)
	assert.True(t, result.Reconciliation.HasErrors)
}

func TestProcessEmployeeForBulkFiling_TaxComputed(t *testing.T) {
	emp := BulkProcessInput{
		PAN:         "PQRST5678H",
		Name:        "Tax Check",
		Email:       "tax@example.com",
		GrossSalary: decimal.NewFromInt(1000000),
		TDSDeducted: decimal.NewFromInt(80000),
	}
	ais := AISSourceData{
		PAN:          "PQRST5678H",
		TaxYear:      "2026-27",
		SalaryIncome: decimal.NewFromInt(1000000),
		TDSEntries: []AISEntry{
			{DeductorTAN: "EMPLOYER", Section: "392", TDSAmount: decimal.NewFromInt(80000)},
		},
	}

	result := ProcessEmployeeForBulkFiling(emp, ais)

	tc := result.TaxComputation
	assert.True(t, tc.GrossIncome.Equal(decimal.NewFromInt(1000000)))
	assert.True(t, tc.StandardDeduction.Equal(decimal.NewFromInt(75000)))
	assert.True(t, tc.TDSCredit.Equal(decimal.NewFromInt(80000)))
	assert.Equal(t, NewRegime, tc.Regime)
}

func TestProcessEmployeeForBulkFiling_ZeroSalary(t *testing.T) {
	emp := BulkProcessInput{
		PAN:         "MNOPQ1111A",
		Name:        "Zero Sal",
		Email:       "zero@example.com",
		GrossSalary: decimal.NewFromInt(0),
		TDSDeducted: decimal.NewFromInt(0),
	}
	ais := AISSourceData{PAN: "MNOPQ1111A", TaxYear: "2026-27"}

	result := ProcessEmployeeForBulkFiling(emp, ais)

	assert.Equal(t, EmpPendingReview, result.Status)
	assert.Equal(t, 0, result.MismatchCount)
}

func BenchmarkProcessEmployee1000(b *testing.B) {
	emp := BulkProcessInput{
		PAN:         "BENCH1234X",
		Name:        "Bench",
		Email:       "b@b.com",
		GrossSalary: decimal.NewFromInt(900000),
		TDSDeducted: decimal.NewFromInt(60000),
	}
	ais := AISSourceData{
		PAN:          "BENCH1234X",
		TaxYear:      "2026-27",
		SalaryIncome: decimal.NewFromInt(900000),
		TDSEntries: []AISEntry{
			{DeductorTAN: "EMP", Section: "392", TDSAmount: decimal.NewFromInt(60000)},
		},
	}
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			ProcessEmployeeForBulkFiling(emp, ais)
		}
	}
}
