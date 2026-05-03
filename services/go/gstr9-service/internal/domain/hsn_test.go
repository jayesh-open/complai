package domain

import (
	"sort"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sampleHSNEntries() []HSNEntry {
	return []HSNEntry{
		{HSNCode: "100110", Description: "Durum wheat", UQC: "KGS", Quantity: decimal.NewFromInt(1000), TaxableValue: decimal.NewFromInt(500000), CGST: decimal.NewFromInt(25000), SGST: decimal.NewFromInt(25000), IGST: decimal.Zero, Cess: decimal.Zero},
		{HSNCode: "100190", Description: "Other wheat", UQC: "KGS", Quantity: decimal.NewFromInt(2000), TaxableValue: decimal.NewFromInt(800000), CGST: decimal.NewFromInt(40000), SGST: decimal.NewFromInt(40000), IGST: decimal.Zero, Cess: decimal.Zero},
		{HSNCode: "100510", Description: "Maize seed", UQC: "KGS", Quantity: decimal.NewFromInt(500), TaxableValue: decimal.NewFromInt(200000), CGST: decimal.NewFromInt(10000), SGST: decimal.NewFromInt(10000), IGST: decimal.Zero, Cess: decimal.Zero},
		{HSNCode: "100590", Description: "Other maize", UQC: "KGS", Quantity: decimal.NewFromInt(300), TaxableValue: decimal.NewFromInt(120000), CGST: decimal.NewFromInt(6000), SGST: decimal.NewFromInt(6000), IGST: decimal.Zero, Cess: decimal.Zero},
		{HSNCode: "200100", Description: "Preserved vegetables", UQC: "KGS", Quantity: decimal.NewFromInt(150), TaxableValue: decimal.NewFromInt(90000), CGST: decimal.NewFromInt(4500), SGST: decimal.NewFromInt(4500), IGST: decimal.Zero, Cess: decimal.Zero},
	}
}

func TestBuildHSNSummary(t *testing.T) {
	entries := sampleHSNEntries()
	summary := BuildHSNSummary(entries)

	assert.Equal(t, 5, len(summary.Entries))
	assert.True(t, summary.TotalValue.Equal(decimal.NewFromInt(1710000)))

	expectedTax := decimal.NewFromInt(171000)
	assert.True(t, summary.TotalTax.Equal(expectedTax), "total tax: got %s, expected %s", summary.TotalTax, expectedTax)
}

func TestHSNDigitLevel_AboveFiveCrore(t *testing.T) {
	turnover := decimal.NewFromInt(60000000)
	assert.Equal(t, 6, HSNDigitLevel(turnover))
}

func TestHSNDigitLevel_ExactlyFiveCrore(t *testing.T) {
	turnover := decimal.NewFromInt(50000000)
	assert.Equal(t, 4, HSNDigitLevel(turnover), "exactly ₹5Cr is not > threshold, so 4-digit")
}

func TestHSNDigitLevel_BetweenOnePointFiveAndFiveCrore(t *testing.T) {
	turnover := decimal.NewFromInt(30000000)
	assert.Equal(t, 4, HSNDigitLevel(turnover))
}

func TestHSNDigitLevel_ExactlyOnePointFiveCrore(t *testing.T) {
	turnover := decimal.NewFromInt(15000000)
	assert.Equal(t, 2, HSNDigitLevel(turnover), "exactly ₹1.5Cr is not > threshold, so 2-digit")
}

func TestHSNDigitLevel_BelowOnePointFiveCrore(t *testing.T) {
	turnover := decimal.NewFromInt(10000000)
	assert.Equal(t, 2, HSNDigitLevel(turnover))
}

func TestAggregateHSN_6Digit_HighTurnover(t *testing.T) {
	entries := sampleHSNEntries()
	turnover := decimal.NewFromInt(60000000)
	summary := AggregateHSNByTurnover(entries, turnover)

	assert.Equal(t, 6, summary.DigitLevel)
	assert.Equal(t, 5, len(summary.Entries), "6-digit: no aggregation, all entries preserved")
	assert.True(t, summary.TotalValue.Equal(decimal.NewFromInt(1710000)))
}

func TestAggregateHSN_4Digit_MediumTurnover(t *testing.T) {
	entries := sampleHSNEntries()
	turnover := decimal.NewFromInt(30000000)
	summary := AggregateHSNByTurnover(entries, turnover)

	assert.Equal(t, 4, summary.DigitLevel)

	sort.Slice(summary.Entries, func(i, j int) bool {
		return summary.Entries[i].HSNCode < summary.Entries[j].HSNCode
	})

	assert.Equal(t, 3, len(summary.Entries), "4-digit: 1001+1001→1001, 1005+1005→1005, 2001→2001")

	entryMap := make(map[string]HSNEntry)
	for _, e := range summary.Entries {
		entryMap[e.HSNCode] = e
	}
	e1001 := entryMap["1001"]
	assert.True(t, e1001.TaxableValue.Equal(decimal.NewFromInt(1300000)),
		"1001 should aggregate 500K + 800K = 1.3M: got %s", e1001.TaxableValue)
	assert.True(t, e1001.Quantity.Equal(decimal.NewFromInt(3000)),
		"1001 quantity should be 1000 + 2000 = 3000")

	e1005 := entryMap["1005"]
	assert.True(t, e1005.TaxableValue.Equal(decimal.NewFromInt(320000)))
}

func TestAggregateHSN_2Digit_LowTurnover(t *testing.T) {
	entries := sampleHSNEntries()
	turnover := decimal.NewFromInt(10000000)
	summary := AggregateHSNByTurnover(entries, turnover)

	assert.Equal(t, 2, summary.DigitLevel)
	assert.Equal(t, 2, len(summary.Entries), "2-digit: 10xxxx→10, 20xxxx→20")

	entryMap := make(map[string]HSNEntry)
	for _, e := range summary.Entries {
		entryMap[e.HSNCode] = e
	}

	e10 := entryMap["10"]
	assert.True(t, e10.TaxableValue.Equal(decimal.NewFromInt(1620000)),
		"HSN 10 should aggregate all 10xxxx entries: got %s", e10.TaxableValue)

	e20 := entryMap["20"]
	assert.True(t, e20.TaxableValue.Equal(decimal.NewFromInt(90000)))
}

func TestAggregateHSN_MixedMonths(t *testing.T) {
	entries := []HSNEntry{
		{HSNCode: "841810", TaxableValue: decimal.NewFromInt(100000), CGST: decimal.NewFromInt(9000), SGST: decimal.NewFromInt(9000), IGST: decimal.Zero, Cess: decimal.Zero, Quantity: decimal.NewFromInt(10)},
		{HSNCode: "841810", TaxableValue: decimal.NewFromInt(200000), CGST: decimal.NewFromInt(18000), SGST: decimal.NewFromInt(18000), IGST: decimal.Zero, Cess: decimal.Zero, Quantity: decimal.NewFromInt(20)},
		{HSNCode: "841820", TaxableValue: decimal.NewFromInt(150000), CGST: decimal.NewFromInt(13500), SGST: decimal.NewFromInt(13500), IGST: decimal.Zero, Cess: decimal.Zero, Quantity: decimal.NewFromInt(15)},
		{HSNCode: "842110", TaxableValue: decimal.NewFromInt(300000), CGST: decimal.NewFromInt(27000), SGST: decimal.NewFromInt(27000), IGST: decimal.Zero, Cess: decimal.Zero, Quantity: decimal.NewFromInt(5)},
	}

	turnover := decimal.NewFromInt(30000000)
	summary := AggregateHSNByTurnover(entries, turnover)
	assert.Equal(t, 4, summary.DigitLevel)

	entryMap := make(map[string]HSNEntry)
	for _, e := range summary.Entries {
		entryMap[e.HSNCode] = e
	}

	require.Contains(t, entryMap, "8418")
	assert.True(t, entryMap["8418"].TaxableValue.Equal(decimal.NewFromInt(450000)),
		"8418 should be 100K + 200K + 150K = 450K: got %s", entryMap["8418"].TaxableValue)
	assert.True(t, entryMap["8418"].Quantity.Equal(decimal.NewFromInt(45)),
		"8418 quantity should be 10+20+15 = 45")

	require.Contains(t, entryMap, "8421")
	assert.True(t, entryMap["8421"].TaxableValue.Equal(decimal.NewFromInt(300000)))
}

func TestAggregateHSN_EmptyEntries(t *testing.T) {
	summary := AggregateHSNByTurnover(nil, decimal.NewFromInt(30000000))
	assert.Equal(t, 4, summary.DigitLevel)
	assert.Empty(t, summary.Entries)
	assert.True(t, summary.TotalValue.IsZero())
}

func TestAggregateHSN_ShortHSNCodes(t *testing.T) {
	entries := []HSNEntry{
		{HSNCode: "84", TaxableValue: decimal.NewFromInt(100000), CGST: decimal.NewFromInt(9000), SGST: decimal.NewFromInt(9000), Quantity: decimal.NewFromInt(10)},
		{HSNCode: "8418", TaxableValue: decimal.NewFromInt(200000), CGST: decimal.NewFromInt(18000), SGST: decimal.NewFromInt(18000), Quantity: decimal.NewFromInt(20)},
	}
	summary := AggregateHSNByTurnover(entries, decimal.NewFromInt(60000000))
	assert.Equal(t, 6, summary.DigitLevel)
	assert.Equal(t, 2, len(summary.Entries), "short codes should be preserved as-is")
}
