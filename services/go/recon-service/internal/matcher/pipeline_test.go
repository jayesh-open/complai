package matcher

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/services/go/recon-service/internal/domain"
)

func TestExactMatch(t *testing.T) {
	runID := uuid.New()
	pr := []domain.PurchaseRegisterEntry{
		{
			InvoiceNumber: "INV-001",
			InvoiceDate:   "15-01-2024",
			VendorGSTIN:   "29AABCU9603R1ZP",
			TotalValue:    decimal.NewFromFloat(10000.50),
			HSN:           "9954",
			SourceID:      "src-1",
		},
	}
	gstr2b := []domain.GSTR2BEntry{
		{
			SupplierGSTIN: "29AABCU9603R1ZP",
			InvoiceNumber: "INV-001",
			InvoiceDate:   "15-01-2024",
			TotalValue:    decimal.NewFromFloat(10000.50),
			HSN:           "9954",
		},
	}

	results := Run(pr, gstr2b, "27AABCU9603R1ZP", "012024", runID)

	require.Len(t, results, 1)
	assert.Equal(t, domain.MatchTypeDirect, results[0].MatchType)
	assert.Equal(t, decimal.NewFromFloat(1.0).String(), results[0].MatchConfidence.String())
	assert.Contains(t, results[0].ReasonCodes, "exact_invoice_number")
	assert.Contains(t, results[0].ReasonCodes, "exact_gstin")
	assert.Contains(t, results[0].ReasonCodes, "exact_amount")
	assert.Contains(t, results[0].ReasonCodes, "exact_date")
}

func TestExactMatch_NormalizedInvoiceNumber(t *testing.T) {
	// "INV-001" in PR matches "INV001" in 2B after normalization (exact match stage)
	runID := uuid.New()
	pr := []domain.PurchaseRegisterEntry{
		{
			InvoiceNumber: "INV-001",
			InvoiceDate:   "15-01-2024",
			VendorGSTIN:   "29AABCU9603R1ZP",
			TotalValue:    decimal.NewFromFloat(10000),
		},
	}
	gstr2b := []domain.GSTR2BEntry{
		{
			SupplierGSTIN: "29AABCU9603R1ZP",
			InvoiceNumber: "INV001",
			InvoiceDate:   "15-01-2024",
			TotalValue:    decimal.NewFromFloat(10000),
		},
	}

	results := Run(pr, gstr2b, "27AABCU9603R1ZP", "012024", runID)
	require.Len(t, results, 1)
	// Should be exact match since normalization strips hyphens
	assert.Equal(t, domain.MatchTypeDirect, results[0].MatchType)
}

func TestFuzzyMatch_InvoiceNumberTypo(t *testing.T) {
	// "INV001" vs "INV002" — DL distance 1, same vendor, same amount, same date => PROBABLE
	runID := uuid.New()
	pr := []domain.PurchaseRegisterEntry{
		{
			InvoiceNumber: "INV001",
			InvoiceDate:   "15-01-2024",
			VendorGSTIN:   "29AABCU9603R1ZP",
			TotalValue:    decimal.NewFromFloat(10000),
		},
	}
	gstr2b := []domain.GSTR2BEntry{
		{
			SupplierGSTIN: "29AABCU9603R1ZP",
			InvoiceNumber: "INV002",
			InvoiceDate:   "15-01-2024",
			TotalValue:    decimal.NewFromFloat(10000),
		},
	}

	results := Run(pr, gstr2b, "27AABCU9603R1ZP", "012024", runID)
	require.Len(t, results, 1)
	assert.Equal(t, domain.MatchTypeProbable, results[0].MatchType)
	assert.Contains(t, results[0].ReasonCodes, "fuzzy_invoice_number")
}

func TestFuzzyMatch_AmountTolerance(t *testing.T) {
	// Same invoice, same vendor, amount within 5%
	runID := uuid.New()
	pr := []domain.PurchaseRegisterEntry{
		{
			InvoiceNumber: "INV100",
			InvoiceDate:   "15-01-2024",
			VendorGSTIN:   "29AABCU9603R1ZP",
			TotalValue:    decimal.NewFromFloat(10000),
		},
	}
	gstr2b := []domain.GSTR2BEntry{
		{
			SupplierGSTIN: "29AABCU9603R1ZP",
			InvoiceNumber: "INV1O0", // "O" instead of "0" — DL distance 1
			InvoiceDate:   "15-01-2024",
			TotalValue:    decimal.NewFromFloat(10400), // 4% diff
		},
	}

	results := Run(pr, gstr2b, "27AABCU9603R1ZP", "012024", runID)
	require.Len(t, results, 1)
	assert.Equal(t, domain.MatchTypeProbable, results[0].MatchType)
	assert.Contains(t, results[0].ReasonCodes, "approx_amount")
}

func TestFuzzyMatch_DateTolerance(t *testing.T) {
	// Same invoice, same vendor, date within 5 days
	runID := uuid.New()
	pr := []domain.PurchaseRegisterEntry{
		{
			InvoiceNumber: "INV200",
			InvoiceDate:   "15-01-2024",
			VendorGSTIN:   "29AABCU9603R1ZP",
			TotalValue:    decimal.NewFromFloat(10000),
		},
	}
	gstr2b := []domain.GSTR2BEntry{
		{
			SupplierGSTIN: "29AABCU9603R1ZP",
			InvoiceNumber: "INV201", // DL dist 1
			InvoiceDate:   "18-01-2024",
			TotalValue:    decimal.NewFromFloat(10000),
		},
	}

	results := Run(pr, gstr2b, "27AABCU9603R1ZP", "012024", runID)
	require.Len(t, results, 1)
	assert.Equal(t, domain.MatchTypeProbable, results[0].MatchType)
	assert.Contains(t, results[0].ReasonCodes, "approx_date")
}

func TestMissing2B(t *testing.T) {
	// Invoice in PR but not in 2B
	runID := uuid.New()
	pr := []domain.PurchaseRegisterEntry{
		{
			InvoiceNumber: "INV-UNIQUE",
			InvoiceDate:   "15-01-2024",
			VendorGSTIN:   "29AABCU9603R1ZP",
			TotalValue:    decimal.NewFromFloat(5000),
		},
	}
	var gstr2b []domain.GSTR2BEntry

	results := Run(pr, gstr2b, "27AABCU9603R1ZP", "012024", runID)
	require.Len(t, results, 1)
	assert.Equal(t, domain.MatchTypeMissing2B, results[0].MatchType)
	assert.Contains(t, results[0].ReasonCodes, "no_matching_2b_entry")
}

func TestMissingPR(t *testing.T) {
	// Invoice in 2B but not in PR
	runID := uuid.New()
	var pr []domain.PurchaseRegisterEntry
	gstr2b := []domain.GSTR2BEntry{
		{
			SupplierGSTIN: "29AABCU9603R1ZP",
			InvoiceNumber: "INV-ONLY-2B",
			InvoiceDate:   "15-01-2024",
			TotalValue:    decimal.NewFromFloat(7000),
		},
	}

	results := Run(pr, gstr2b, "27AABCU9603R1ZP", "012024", runID)
	require.Len(t, results, 1)
	assert.Equal(t, domain.MatchTypeMissingPR, results[0].MatchType)
	assert.Contains(t, results[0].ReasonCodes, "no_matching_pr_entry")
}

func TestPartialMatch_SplitInvoice(t *testing.T) {
	// PR has 10,000 total; 2B has two invoices from same vendor: 6,000 + 4,000
	runID := uuid.New()
	pr := []domain.PurchaseRegisterEntry{
		{
			InvoiceNumber: "INV-SPLIT",
			InvoiceDate:   "15-01-2024",
			VendorGSTIN:   "29AABCU9603R1ZP",
			TotalValue:    decimal.NewFromFloat(10000),
		},
	}
	gstr2b := []domain.GSTR2BEntry{
		{
			SupplierGSTIN: "29AABCU9603R1ZP",
			InvoiceNumber: "INV-SPLIT-A",
			InvoiceDate:   "15-01-2024",
			TotalValue:    decimal.NewFromFloat(6000),
		},
		{
			SupplierGSTIN: "29AABCU9603R1ZP",
			InvoiceNumber: "INV-SPLIT-B",
			InvoiceDate:   "15-01-2024",
			TotalValue:    decimal.NewFromFloat(4000),
		},
	}

	results := Run(pr, gstr2b, "27AABCU9603R1ZP", "012024", runID)
	// Should produce 2 partial matches (one for each 2B entry in the split)
	partialCount := 0
	for _, m := range results {
		if m.MatchType == domain.MatchTypePartial {
			partialCount++
			assert.Contains(t, m.ReasonCodes, "split_invoice")
		}
	}
	assert.Equal(t, 2, partialCount)
}

func TestLargeDataset_Precision(t *testing.T) {
	// 1000 PR + 1000 2B with known overlap:
	// 850 exact matches, 100 fuzzy matches, 50 unmatched on each side
	runID := uuid.New()

	pr := make([]domain.PurchaseRegisterEntry, 1000)
	gstr2b := make([]domain.GSTR2BEntry, 1000)

	// First 850: exact matches
	for i := 0; i < 850; i++ {
		invNum := fmt.Sprintf("EXACT-%04d", i)
		amount := decimal.NewFromFloat(float64(1000 + i))
		pr[i] = domain.PurchaseRegisterEntry{
			InvoiceNumber: invNum,
			InvoiceDate:   "15-01-2024",
			VendorGSTIN:   "29AABCU9603R1ZP",
			TotalValue:    amount,
		}
		gstr2b[i] = domain.GSTR2BEntry{
			SupplierGSTIN: "29AABCU9603R1ZP",
			InvoiceNumber: invNum,
			InvoiceDate:   "15-01-2024",
			TotalValue:    amount,
		}
	}

	// Next 100: fuzzy matches (slight invoice number difference)
	for i := 0; i < 100; i++ {
		idx := 850 + i
		amount := decimal.NewFromFloat(float64(2000 + i))
		pr[idx] = domain.PurchaseRegisterEntry{
			InvoiceNumber: fmt.Sprintf("FUZZ-%04d", i),
			InvoiceDate:   "15-01-2024",
			VendorGSTIN:   "29AABCU9603R1ZP",
			TotalValue:    amount,
		}
		gstr2b[idx] = domain.GSTR2BEntry{
			SupplierGSTIN: "29AABCU9603R1ZP",
			InvoiceNumber: fmt.Sprintf("FUZZ-%04dX", i), // DL dist 1 (appended X)
			InvoiceDate:   "15-01-2024",
			TotalValue:    amount,
		}
	}

	// Last 50: unmatched on each side (different vendors to prevent cross-matching)
	for i := 0; i < 50; i++ {
		idx := 950 + i
		pr[idx] = domain.PurchaseRegisterEntry{
			InvoiceNumber: fmt.Sprintf("PR-ONLY-%04d", i),
			InvoiceDate:   "15-01-2024",
			VendorGSTIN:   "29BBBBX0000A1ZA",
			TotalValue:    decimal.NewFromFloat(float64(3000 + i)),
		}
		gstr2b[idx] = domain.GSTR2BEntry{
			SupplierGSTIN: "29CCCCX0000B1ZB",
			InvoiceNumber: fmt.Sprintf("2B-ONLY-%04d", i),
			InvoiceDate:   "15-01-2024",
			TotalValue:    decimal.NewFromFloat(float64(4000 + i)),
		}
	}

	results := Run(pr, gstr2b, "27AABCU9603R1ZP", "012024", runID)

	var directCount, probableCount, missing2bCount, missingPRCount int
	for _, m := range results {
		switch m.MatchType {
		case domain.MatchTypeDirect:
			directCount++
		case domain.MatchTypeProbable:
			probableCount++
		case domain.MatchTypeMissing2B:
			missing2bCount++
		case domain.MatchTypeMissingPR:
			missingPRCount++
		}
	}

	// Precision: exact matches should be 850
	assert.Equal(t, 850, directCount, "expected 850 exact matches")

	// Recall: fuzzy matches should be close to 100
	assert.GreaterOrEqual(t, probableCount, 90, "expected at least 90 fuzzy matches (>90%% recall)")

	// Unmatched
	assert.Equal(t, 50, missing2bCount, "expected 50 missing-2B")
	assert.Equal(t, 50, missingPRCount, "expected 50 missing-PR")

	// Overall precision: (direct + probable) / total matched >= 95%
	totalMatched := directCount + probableCount
	assert.GreaterOrEqual(t, totalMatched, 940, "expected at least 940 total matches (>95%% precision)")
}

func TestPerformance_10kx10k(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping performance test in short mode")
	}

	runID := uuid.New()
	pr := make([]domain.PurchaseRegisterEntry, 10000)
	gstr2b := make([]domain.GSTR2BEntry, 10000)

	for i := 0; i < 10000; i++ {
		invNum := fmt.Sprintf("PERF-%06d", i)
		amount := decimal.NewFromFloat(float64(1000 + i))
		pr[i] = domain.PurchaseRegisterEntry{
			InvoiceNumber: invNum,
			InvoiceDate:   "15-01-2024",
			VendorGSTIN:   fmt.Sprintf("29GSTIN%05d1ZP", i%100),
			TotalValue:    amount,
		}
		gstr2b[i] = domain.GSTR2BEntry{
			SupplierGSTIN: fmt.Sprintf("29GSTIN%05d1ZP", i%100),
			InvoiceNumber: invNum,
			InvoiceDate:   "15-01-2024",
			TotalValue:    amount,
		}
	}

	results := Run(pr, gstr2b, "27AABCU9603R1ZP", "012024", runID)

	var directCount int
	for _, m := range results {
		if m.MatchType == domain.MatchTypeDirect {
			directCount++
		}
	}

	assert.Equal(t, 10000, directCount, "all 10k should be exact matches")
	assert.Len(t, results, 10000)
}

func TestDamerauLevenshtein(t *testing.T) {
	tests := []struct {
		a, b     string
		expected int
	}{
		{"", "", 0},
		{"abc", "", 3},
		{"", "abc", 3},
		{"abc", "abc", 0},
		{"abc", "abd", 1},       // substitution
		{"abc", "abcd", 1},      // insertion
		{"abcd", "abc", 1},      // deletion
		{"abc", "bac", 1},       // transposition
		{"INV001", "INV002", 1}, // substitution
		{"INV001", "INV0001", 1},
		{"ABCDEF", "XYZXYZ", 6},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_%s", tt.a, tt.b), func(t *testing.T) {
			assert.Equal(t, tt.expected, damerauLevenshtein(tt.a, tt.b))
		})
	}
}

func TestNormalizeInvoiceNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"INV-001", "INV001"},
		{"INV/001", "INV001"},
		{"inv 001", "INV001"},
		{"INV-001/A", "INV001A"},
		{"abc", "ABC"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, normalizeInvoiceNumber(tt.input))
		})
	}
}

func TestDateDiffDays(t *testing.T) {
	tests := []struct {
		a, b     string
		expected int
	}{
		{"15-01-2024", "15-01-2024", 0},
		{"15-01-2024", "18-01-2024", 3},
		{"18-01-2024", "15-01-2024", 3},
		{"15-01-2024", "20-01-2024", 5},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_%s", tt.a, tt.b), func(t *testing.T) {
			assert.Equal(t, tt.expected, dateDiffDays(tt.a, tt.b))
		})
	}
}

func TestAmountDiffPercent(t *testing.T) {
	a := decimal.NewFromFloat(10000)
	b := decimal.NewFromFloat(10500)
	pct := amountDiffPercent(a, b)
	assert.InDelta(t, 5.0, pct, 0.01)
}

func TestMixedScenario(t *testing.T) {
	runID := uuid.New()
	pr := []domain.PurchaseRegisterEntry{
		// Will be exact match
		{InvoiceNumber: "INV-001", InvoiceDate: "15-01-2024", VendorGSTIN: "29AAA", TotalValue: decimal.NewFromFloat(1000)},
		// Will be fuzzy match
		{InvoiceNumber: "INV002", InvoiceDate: "15-01-2024", VendorGSTIN: "29BBB", TotalValue: decimal.NewFromFloat(2000)},
		// Will be missing 2B
		{InvoiceNumber: "INV003", InvoiceDate: "15-01-2024", VendorGSTIN: "29CCC", TotalValue: decimal.NewFromFloat(3000)},
	}
	gstr2b := []domain.GSTR2BEntry{
		// Exact match for INV-001
		{SupplierGSTIN: "29AAA", InvoiceNumber: "INV-001", InvoiceDate: "15-01-2024", TotalValue: decimal.NewFromFloat(1000)},
		// Fuzzy match for INV002 (typo)
		{SupplierGSTIN: "29BBB", InvoiceNumber: "INV0O2", InvoiceDate: "15-01-2024", TotalValue: decimal.NewFromFloat(2000)},
		// Missing PR
		{SupplierGSTIN: "29DDD", InvoiceNumber: "INV004", InvoiceDate: "15-01-2024", TotalValue: decimal.NewFromFloat(4000)},
	}

	results := Run(pr, gstr2b, "27TEST", "012024", runID)

	var types []domain.MatchType
	for _, m := range results {
		types = append(types, m.MatchType)
	}

	assert.Contains(t, types, domain.MatchTypeDirect)
	assert.Contains(t, types, domain.MatchTypeProbable)
	assert.Contains(t, types, domain.MatchTypeMissing2B)
	assert.Contains(t, types, domain.MatchTypeMissingPR)
}
