package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReconcileTDS_AllMatched(t *testing.T) {
	ais := []AISEntry{
		{DeductorTAN: "MUMB12345A", Section: "392", TDSAmount: d(50000)},
		{DeductorTAN: "DELH67890B", Section: "393(1)", TDSAmount: d(10000)},
	}
	claims := []TDSCreditEntry{
		{DeductorTAN: "MUMB12345A", Section: "392", TDSAmount: d(50000)},
		{DeductorTAN: "DELH67890B", Section: "393(1)", TDSAmount: d(10000)},
	}

	result := ReconcileTDS(ais, claims)
	assert.Equal(t, d(60000), result.TotalAIS)
	assert.Equal(t, d(60000), result.TotalClaim)
	assert.True(t, result.Difference.IsZero())
	require.Len(t, result.Matched, 2)
	assert.Empty(t, result.Unmatched)
	for _, m := range result.Matched {
		assert.Equal(t, "MATCHED", m.Status)
	}
}

func TestReconcileTDS_Discrepancy(t *testing.T) {
	ais := []AISEntry{
		{DeductorTAN: "MUMB12345A", Section: "392", TDSAmount: d(50000)},
	}
	claims := []TDSCreditEntry{
		{DeductorTAN: "MUMB12345A", Section: "392", TDSAmount: d(45000)},
	}

	result := ReconcileTDS(ais, claims)
	require.Len(t, result.Matched, 1)
	assert.Equal(t, "DISCREPANCY", result.Matched[0].Status)
	assert.Equal(t, d(5000), result.Matched[0].Discrepancy)
}

func TestReconcileTDS_AISOnly(t *testing.T) {
	ais := []AISEntry{
		{DeductorTAN: "MUMB12345A", Section: "392", TDSAmount: d(50000)},
	}
	var claims []TDSCreditEntry

	result := ReconcileTDS(ais, claims)
	require.Len(t, result.Unmatched, 1)
	assert.Equal(t, "AIS", result.Unmatched[0].Source)
	assert.Contains(t, result.Unmatched[0].Issue, "not claimed")
}

func TestReconcileTDS_ClaimOnly(t *testing.T) {
	var ais []AISEntry
	claims := []TDSCreditEntry{
		{DeductorTAN: "MUMB12345A", Section: "392", TDSAmount: d(50000)},
	}

	result := ReconcileTDS(ais, claims)
	require.Len(t, result.Unmatched, 1)
	assert.Equal(t, "CLAIM", result.Unmatched[0].Source)
	assert.Contains(t, result.Unmatched[0].Issue, "not found in AIS")
}

func TestReconcileTDS_MultipleAIS_SameTAN(t *testing.T) {
	ais := []AISEntry{
		{DeductorTAN: "MUMB12345A", Section: "392", TDSAmount: d(25000)},
		{DeductorTAN: "MUMB12345A", Section: "392", TDSAmount: d(25000)},
	}
	claims := []TDSCreditEntry{
		{DeductorTAN: "MUMB12345A", Section: "392", TDSAmount: d(50000)},
	}

	result := ReconcileTDS(ais, claims)
	require.Len(t, result.Matched, 1)
	assert.Equal(t, "MATCHED", result.Matched[0].Status)
}
