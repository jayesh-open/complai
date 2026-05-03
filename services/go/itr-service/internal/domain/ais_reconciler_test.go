package domain

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ais(salary, interest, dividend, securities, property float64) AISSourceData {
	return AISSourceData{
		PAN:               "ABCDE1234F",
		TaxYear:           "2026-27",
		SalaryIncome:      d(int64(salary)),
		InterestIncome:    d(int64(interest)),
		DividendIncome:    d(int64(dividend)),
		SecuritiesTrading: d(int64(securities)),
		PropertyTxnValue:  d(int64(property)),
	}
}

func books(salary, interest, dividend, securities, property float64) BookData {
	return BookData{
		SalaryIncome:   d(int64(salary)),
		InterestIncome: d(int64(interest)),
		DividendIncome: d(int64(dividend)),
		SecuritiesInCG: d(int64(securities)),
		PropertyInHP:   d(int64(property)),
	}
}

func TestReconcileAIS_AllMatched(t *testing.T) {
	result := ReconcileAIS(
		ais(1200000, 50000, 25000, 0, 0),
		books(1200000, 50000, 25000, 0, 0),
		true,
	)
	assert.Empty(t, result.Mismatches)
	assert.False(t, result.HasErrors)
	assert.False(t, result.SubmissionBlocked)
}

func TestReconcileAIS_SalaryMismatch(t *testing.T) {
	result := ReconcileAIS(
		ais(1200000, 0, 0, 0, 0),
		books(1100000, 0, 0, 0, 0),
		true,
	)
	require.Len(t, result.Mismatches, 1)
	assert.Equal(t, "SALARY", result.Mismatches[0].Category)
	assert.Equal(t, SeverityError, result.Mismatches[0].Severity)
	assert.Equal(t, d(100000), result.Mismatches[0].Delta)
	assert.True(t, result.HasErrors)
	assert.True(t, result.SubmissionBlocked)
}

func TestReconcileAIS_SalaryMismatch_NotBlocked(t *testing.T) {
	result := ReconcileAIS(
		ais(1200000, 0, 0, 0, 0),
		books(1100000, 0, 0, 0, 0),
		false,
	)
	assert.True(t, result.HasErrors)
	assert.False(t, result.SubmissionBlocked)
}

func TestReconcileAIS_InterestNotInBooks(t *testing.T) {
	result := ReconcileAIS(
		ais(0, 50000, 0, 0, 0),
		books(0, 0, 0, 0, 0),
		true,
	)
	require.Len(t, result.Mismatches, 1)
	assert.Equal(t, "INTEREST", result.Mismatches[0].Category)
	assert.Equal(t, SeverityError, result.Mismatches[0].Severity)
	assert.Contains(t, result.Mismatches[0].SuggestedAction, "Other Sources")
}

func TestReconcileAIS_InterestMismatch(t *testing.T) {
	result := ReconcileAIS(
		ais(0, 50000, 0, 0, 0),
		books(0, 45000, 0, 0, 0),
		true,
	)
	require.Len(t, result.Mismatches, 1)
	assert.Equal(t, "INTEREST", result.Mismatches[0].Category)
	assert.Equal(t, SeverityWarn, result.Mismatches[0].Severity)
}

func TestReconcileAIS_DividendNotInBooks(t *testing.T) {
	result := ReconcileAIS(
		ais(0, 0, 25000, 0, 0),
		books(0, 0, 0, 0, 0),
		true,
	)
	require.Len(t, result.Mismatches, 1)
	assert.Equal(t, "DIVIDEND", result.Mismatches[0].Category)
	assert.Equal(t, SeverityError, result.Mismatches[0].Severity)
}

func TestReconcileAIS_DividendMismatch(t *testing.T) {
	result := ReconcileAIS(
		ais(0, 0, 25000, 0, 0),
		books(0, 0, 20000, 0, 0),
		false,
	)
	require.Len(t, result.Mismatches, 1)
	assert.Equal(t, SeverityWarn, result.Mismatches[0].Severity)
}

func TestReconcileAIS_SecuritiesNotInCG(t *testing.T) {
	result := ReconcileAIS(
		ais(0, 0, 0, 500000, 0),
		books(0, 0, 0, 0, 0),
		true,
	)
	require.Len(t, result.Mismatches, 1)
	assert.Equal(t, "SECURITIES", result.Mismatches[0].Category)
	assert.Equal(t, SeverityError, result.Mismatches[0].Severity)
	assert.Contains(t, result.Mismatches[0].SuggestedAction, "Schedule CG")
}

func TestReconcileAIS_SecuritiesMismatch_Info(t *testing.T) {
	result := ReconcileAIS(
		ais(0, 0, 0, 500000, 0),
		books(0, 0, 0, 450000, 0),
		true,
	)
	require.Len(t, result.Mismatches, 1)
	assert.Equal(t, SeverityInfo, result.Mismatches[0].Severity)
}

func TestReconcileAIS_PropertyNotInHP(t *testing.T) {
	result := ReconcileAIS(
		ais(0, 0, 0, 0, 8000000),
		books(0, 0, 0, 0, 0),
		true,
	)
	require.Len(t, result.Mismatches, 1)
	assert.Equal(t, "PROPERTY", result.Mismatches[0].Category)
	assert.Equal(t, SeverityError, result.Mismatches[0].Severity)
}

func TestReconcileAIS_PropertyPresent_NoMismatch(t *testing.T) {
	result := ReconcileAIS(
		ais(0, 0, 0, 0, 8000000),
		books(0, 0, 0, 0, 8000000),
		true,
	)
	assert.Empty(t, result.Mismatches)
}

func TestReconcileAIS_MultipleMismatches(t *testing.T) {
	result := ReconcileAIS(
		ais(1200000, 50000, 25000, 500000, 8000000),
		books(1100000, 0, 0, 0, 0),
		true,
	)
	assert.Equal(t, 5, len(result.Mismatches))
	assert.Equal(t, 5, result.ErrorCount)
	assert.True(t, result.SubmissionBlocked)
}

func TestReconcileAIS_CountsSeverities(t *testing.T) {
	a := ais(1200000, 50000, 25000, 500000, 0)
	b := books(1100000, 45000, 20000, 450000, 0)
	result := ReconcileAIS(a, b, false)
	assert.Equal(t, 1, result.ErrorCount)
	assert.Equal(t, 2, result.WarnCount)
	assert.Equal(t, 1, result.InfoCount)
	assert.True(t, result.HasErrors)
	assert.False(t, result.SubmissionBlocked)
}

func TestReconcileAIS_TDSMismatch(t *testing.T) {
	a := AISSourceData{
		PAN:     "ABCDE1234F",
		TaxYear: "2026-27",
		TDSEntries: []AISEntry{
			{DeductorTAN: "MUMB12345A", Section: "392", TDSAmount: d(60000)},
		},
	}
	b := BookData{
		TDSClaims: []TDSCreditEntry{
			{DeductorTAN: "MUMB12345A", Section: "392", TDSAmount: d(50000)},
		},
	}
	result := ReconcileAIS(a, b, true)
	hasTDS := false
	for _, m := range result.Mismatches {
		if m.Category == "TDS" {
			hasTDS = true
		}
	}
	assert.True(t, hasTDS)
}

func TestReconcileAIS_TDSUnmatched_AISOnly(t *testing.T) {
	a := AISSourceData{
		PAN:     "ABCDE1234F",
		TaxYear: "2026-27",
		TDSEntries: []AISEntry{
			{DeductorTAN: "DELH67890B", Section: "393(1)", TDSAmount: d(10000)},
		},
	}
	b := BookData{}
	result := ReconcileAIS(a, b, true)
	require.NotEmpty(t, result.Mismatches)
	found := false
	for _, m := range result.Mismatches {
		if m.Category == "TDS" && m.Source == "AIS" {
			found = true
			assert.Equal(t, SeverityError, m.Severity)
		}
	}
	assert.True(t, found)
}

func TestReconcileAIS_ZeroAIS_ZeroBooks(t *testing.T) {
	result := ReconcileAIS(
		ais(0, 0, 0, 0, 0),
		books(0, 0, 0, 0, 0),
		true,
	)
	assert.Empty(t, result.Mismatches)
	assert.False(t, result.HasErrors)
}

func TestReconcileAIS_SalaryWithinThreshold(t *testing.T) {
	a := AISSourceData{PAN: "X", TaxYear: "2026-27", SalaryIncome: decimal.NewFromInt(1200000)}
	b := BookData{SalaryIncome: decimal.NewFromInt(1200001)}
	result := ReconcileAIS(a, b, true)
	assert.Empty(t, result.Mismatches)
}
