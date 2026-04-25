package scorer

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/services/go/vendor-compliance-service/internal/domain"
)

func makeVendor(gstin, pan, tan, email string, msme bool) domain.VendorSnapshot {
	return domain.VendorSnapshot{
		ID:                 uuid.New(),
		TenantID:           uuid.New(),
		VendorID:           "V001",
		Name:               "Test Vendor",
		LegalName:          "Test Vendor Pvt Ltd",
		TradeName:          "Test Vendor",
		PAN:                pan,
		GSTIN:              gstin,
		TAN:                tan,
		State:              "Karnataka",
		StateCode:          "29",
		Category:           "Regular",
		RegistrationStatus: "Active",
		MSMERegistered:     msme,
		Email:              email,
		Phone:              "9876543210",
		Address:            "Bangalore",
		SyncedAt:           time.Now().UTC(),
		CreatedAt:          time.Now().UTC(),
		UpdatedAt:          time.Now().UTC(),
	}
}

func makeInvoices(count int, filingStatus, mismatchStatus, paymentStatus string, irn bool) []APInvoice {
	invoices := make([]APInvoice, count)
	for i := 0; i < count; i++ {
		invoices[i] = APInvoice{
			IRNGenerated:    irn,
			GSTFilingStatus: filingStatus,
			MismatchStatus:  mismatchStatus,
			PaymentStatus:   paymentStatus,
			PaymentDate:     "2026-03-15",
			DueDate:         "2026-03-20",
		}
	}
	return invoices
}

// Test 1: Vendor with 100% on-time GSTR-1 + no mismatches -> score >= 90 (Cat A)
func TestScore_PerfectFilingNoMismatches_CatA(t *testing.T) {
	vendor := makeVendor("29AABCA1234A1Z5", "AABCA1234A", "BLRA12345B", "test@example.com", true)
	invoices := makeInvoices(10, "filed", "matched", "paid", true)

	score := Score(vendor, invoices)

	assert.GreaterOrEqual(t, score.TotalScore, 90)
	assert.Equal(t, "A", score.Category)
	assert.Equal(t, "Low", score.RiskLevel)
}

// Test 2: Vendor with 3 late filings (out of 10) + mismatches + some overdue -> score 60-89 (Cat B)
func TestScore_SomeLateFilingsSmallMismatch_CatB(t *testing.T) {
	// Partial doc hygiene: no TAN, no MSME -> loses 6 points (score = 9/15)
	vendor := makeVendor("29AABCA1234A1Z5", "AABCA1234A", "", "test@example.com", false)

	// 7 filed on time, 3 late; mix of IRN, mismatch, payment issues
	invoices := make([]APInvoice, 10)
	for i := 0; i < 7; i++ {
		invoices[i] = APInvoice{
			IRNGenerated:    true,
			GSTFilingStatus: "filed",
			MismatchStatus:  "matched",
			PaymentStatus:   "paid",
			PaymentDate:     "2026-03-15",
			DueDate:         "2026-03-20",
		}
	}
	for i := 7; i < 10; i++ {
		invoices[i] = APInvoice{
			IRNGenerated:    false,
			GSTFilingStatus: "late",
			MismatchStatus:  "mismatched",
			PaymentStatus:   "overdue",
		}
	}

	score := Score(vendor, invoices)

	assert.GreaterOrEqual(t, score.TotalScore, 60)
	assert.LessOrEqual(t, score.TotalScore, 89)
	assert.Equal(t, "B", score.Category)
}

// Test 3: Vendor with no IRN compliance + missing GSTIN + 30% mismatch -> score < 40 (Cat D)
func TestScore_NoIRNMissingGSTINHighMismatch_CatD(t *testing.T) {
	// Missing GSTIN, PAN, TAN, email, not MSME -> document hygiene = 0
	vendor := makeVendor("", "", "", "", false)

	invoices := make([]APInvoice, 10)
	for i := 0; i < 10; i++ {
		invoices[i] = APInvoice{
			IRNGenerated:    false,
			GSTFilingStatus: "not_filed",
			MismatchStatus:  "mismatched",
			PaymentStatus:   "overdue",
		}
	}

	score := Score(vendor, invoices)

	assert.Less(t, score.TotalScore, 40)
	assert.Equal(t, "D", score.Category)
	assert.Equal(t, "Critical", score.RiskLevel)
}

// Test 4: Perfect vendor (all dimensions maxed) = 100 points
func TestScore_PerfectVendor_100Points(t *testing.T) {
	vendor := makeVendor("29AABCA1234A1Z5", "AABCA1234A", "BLRA12345B", "test@example.com", true)
	invoices := makeInvoices(20, "filed", "matched", "paid", true)

	score := Score(vendor, invoices)

	assert.Equal(t, 100, score.TotalScore)
	assert.Equal(t, "A", score.Category)
	assert.Equal(t, "Low", score.RiskLevel)
	assert.Equal(t, 30, score.FilingRegularityScore)
	assert.Equal(t, 20, score.IRNComplianceScore)
	assert.Equal(t, 20, score.MismatchRateScore)
	assert.Equal(t, 15, score.PaymentBehaviorScore)
	assert.Equal(t, 15, score.DocumentHygieneScore)
}

// Test 5: Empty invoices -> scoring defaults
func TestScore_EmptyInvoices_Defaults(t *testing.T) {
	vendor := makeVendor("29AABCA1234A1Z5", "AABCA1234A", "BLRA12345B", "test@example.com", true)
	invoices := []APInvoice{}

	score := Score(vendor, invoices)

	// Default scores: filing=15, irn=10, mismatch=10, payment=8, doc=15
	assert.Equal(t, 15, score.FilingRegularityScore)
	assert.Equal(t, 10, score.IRNComplianceScore)
	assert.Equal(t, 10, score.MismatchRateScore)
	assert.Equal(t, 8, score.PaymentBehaviorScore)
	assert.Equal(t, 15, score.DocumentHygieneScore)
	assert.Equal(t, 58, score.TotalScore)
}

// Test 6: All overdue payments -> low payment behavior score
func TestScore_AllOverduePayments_LowPaymentScore(t *testing.T) {
	vendor := makeVendor("29AABCA1234A1Z5", "AABCA1234A", "BLRA12345B", "test@example.com", true)
	invoices := makeInvoices(10, "filed", "matched", "overdue", true)

	score := Score(vendor, invoices)

	assert.Equal(t, 0, score.PaymentBehaviorScore)
	assert.Contains(t, score.PaymentBehaviorNote, "0/10 paid on time")
}

// Test 7: Score boundary - exactly 90 = A
func TestScore_BoundaryExactly90_CatA(t *testing.T) {
	assert.Equal(t, "A", categorize(90))
}

// Test 8: Score boundary - exactly 60 = B
func TestScore_BoundaryExactly60_CatB(t *testing.T) {
	assert.Equal(t, "B", categorize(60))
}

// Test 9: Score boundary - exactly 40 = C
func TestScore_BoundaryExactly40_CatC(t *testing.T) {
	assert.Equal(t, "C", categorize(40))
}

// Test 10: Score boundary - 39 = D
func TestScore_Boundary39_CatD(t *testing.T) {
	assert.Equal(t, "D", categorize(39))
}

// Test 11: Risk level boundaries
func TestScore_RiskLevelBoundaries(t *testing.T) {
	assert.Equal(t, "Low", riskLevel(80))
	assert.Equal(t, "Low", riskLevel(100))
	assert.Equal(t, "Medium", riskLevel(60))
	assert.Equal(t, "Medium", riskLevel(79))
	assert.Equal(t, "High", riskLevel(40))
	assert.Equal(t, "High", riskLevel(59))
	assert.Equal(t, "Critical", riskLevel(39))
	assert.Equal(t, "Critical", riskLevel(0))
}

// Test 12: Score populates all fields correctly
func TestScore_FieldsPopulated(t *testing.T) {
	vendor := makeVendor("29AABCA1234A1Z5", "AABCA1234A", "BLRA12345B", "test@example.com", true)
	invoices := makeInvoices(5, "filed", "matched", "paid", true)

	score := Score(vendor, invoices)

	require.NotEqual(t, uuid.Nil, score.ID)
	assert.Equal(t, vendor.TenantID, score.TenantID)
	assert.Equal(t, vendor.VendorID, score.VendorID)
	assert.Equal(t, vendor.ID, score.VendorSnapshotID)
	assert.NotEmpty(t, score.FilingRegularityNote)
	assert.NotEmpty(t, score.IRNComplianceNote)
	assert.NotEmpty(t, score.MismatchRateNote)
	assert.NotEmpty(t, score.PaymentBehaviorNote)
	assert.NotEmpty(t, score.DocumentHygieneNote)
	assert.False(t, score.ScoredAt.IsZero())
	assert.False(t, score.CreatedAt.IsZero())
}

// Test 13: Document hygiene with partial fields
func TestScore_DocumentHygiene_PartialFields(t *testing.T) {
	// Has GSTIN and PAN but no TAN, email, not MSME -> score = 6
	vendor := makeVendor("29AABCA1234A1Z5", "AABCA1234A", "", "", false)
	invoices := makeInvoices(5, "filed", "matched", "paid", true)

	score := Score(vendor, invoices)

	assert.Equal(t, 6, score.DocumentHygieneScore)
	assert.Contains(t, score.DocumentHygieneNote, "TAN")
	assert.Contains(t, score.DocumentHygieneNote, "Email")
	assert.Contains(t, score.DocumentHygieneNote, "MSME")
}

// Test 14: Mixed invoices with different statuses
func TestScore_MixedInvoiceStatuses(t *testing.T) {
	vendor := makeVendor("29AABCA1234A1Z5", "AABCA1234A", "BLRA12345B", "test@example.com", true)

	invoices := []APInvoice{
		{IRNGenerated: true, GSTFilingStatus: "filed", MismatchStatus: "matched", PaymentStatus: "paid"},
		{IRNGenerated: true, GSTFilingStatus: "filed", MismatchStatus: "matched", PaymentStatus: "paid"},
		{IRNGenerated: false, GSTFilingStatus: "late", MismatchStatus: "mismatched", PaymentStatus: "overdue"},
		{IRNGenerated: true, GSTFilingStatus: "filed", MismatchStatus: "matched", PaymentStatus: "paid"},
		{IRNGenerated: false, GSTFilingStatus: "not_filed", MismatchStatus: "pending", PaymentStatus: "unpaid"},
	}

	score := Score(vendor, invoices)

	// Filing: 3 filed + 1 late*0.5 + 1 not_filed*0 = 3.5/5 -> int(3.5/5 * 30) = 21
	assert.Equal(t, 21, score.FilingRegularityScore)
	// IRN: 3/5 -> int(3/5 * 20) = 12
	assert.Equal(t, 12, score.IRNComplianceScore)
	// Mismatch: 3/5 matched -> int(3/5 * 20) = 12
	assert.Equal(t, 12, score.MismatchRateScore)
	// Payment: 3/5 paid -> int(3/5 * 15) = 9
	assert.Equal(t, 9, score.PaymentBehaviorScore)
	// Doc hygiene: all fields present -> 15
	assert.Equal(t, 15, score.DocumentHygieneScore)
	// Total: 21 + 12 + 12 + 9 + 15 = 69
	assert.Equal(t, 69, score.TotalScore)
	assert.Equal(t, "B", score.Category)
	assert.Equal(t, "Medium", score.RiskLevel)
}
