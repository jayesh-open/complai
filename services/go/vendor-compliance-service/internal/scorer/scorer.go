package scorer

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/vendor-compliance-service/internal/domain"
)

// APInvoice is the scoring-relevant subset of an AP invoice from Apex.
type APInvoice struct {
	IRNGenerated    bool   `json:"irn_generated"`
	GSTFilingStatus string `json:"gst_filing_status"` // filed, pending, late, not_filed
	MismatchStatus  string `json:"mismatch_status"`   // matched, mismatched, pending
	PaymentStatus   string `json:"payment_status"`    // paid, unpaid, overdue, partial
	PaymentDate     string `json:"payment_date"`
	DueDate         string `json:"due_date"`
}

// Score computes a ComplianceScore for a vendor given their snapshot and AP invoices.
// This is a pure function with no side effects.
func Score(vendor domain.VendorSnapshot, invoices []APInvoice) domain.ComplianceScore {
	now := time.Now().UTC()

	filingScore, filingNote := scoreFilingRegularity(invoices)
	irnScore, irnNote := scoreIRNCompliance(invoices)
	mismatchScore, mismatchNote := scoreMismatchRate(invoices)
	paymentScore, paymentNote := scorePaymentBehavior(invoices)
	docScore, docNote := scoreDocumentHygiene(vendor)

	total := filingScore + irnScore + mismatchScore + paymentScore + docScore

	return domain.ComplianceScore{
		ID:                    uuid.New(),
		TenantID:              vendor.TenantID,
		VendorID:              vendor.VendorID,
		VendorSnapshotID:      vendor.ID,
		TotalScore:            total,
		Category:              categorize(total),
		RiskLevel:             riskLevel(total),
		FilingRegularityScore: filingScore,
		IRNComplianceScore:    irnScore,
		MismatchRateScore:     mismatchScore,
		PaymentBehaviorScore:  paymentScore,
		DocumentHygieneScore:  docScore,
		FilingRegularityNote:  filingNote,
		IRNComplianceNote:     irnNote,
		MismatchRateNote:      mismatchNote,
		PaymentBehaviorNote:   paymentNote,
		DocumentHygieneNote:   docNote,
		ScoredAt:              now,
		CreatedAt:             now,
	}
}

// scoreFilingRegularity: 30 points max.
// Based on GST filing status of invoices.
// filed = full points, late = half points, pending/not_filed = 0 points.
func scoreFilingRegularity(invoices []APInvoice) (int, string) {
	if len(invoices) == 0 {
		return 15, "No invoices to evaluate; default score applied"
	}

	var filed, late, notFiled int
	for _, inv := range invoices {
		switch inv.GSTFilingStatus {
		case "filed":
			filed++
		case "late":
			late++
		default: // pending, not_filed
			notFiled++
		}
	}

	total := len(invoices)
	// filed gets full weight, late gets half weight
	points := float64(filed) + float64(late)*0.5
	score := int((points / float64(total)) * 30)

	if score > 30 {
		score = 30
	}

	note := fmt.Sprintf("%d/%d filed on time, %d late, %d not filed", filed, total, late, notFiled)
	return score, note
}

// scoreIRNCompliance: 20 points max.
// Based on IRN generation rate.
func scoreIRNCompliance(invoices []APInvoice) (int, string) {
	if len(invoices) == 0 {
		return 10, "No invoices to evaluate; default score applied"
	}

	irnCount := 0
	for _, inv := range invoices {
		if inv.IRNGenerated {
			irnCount++
		}
	}

	total := len(invoices)
	score := int((float64(irnCount) / float64(total)) * 20)

	if score > 20 {
		score = 20
	}

	note := fmt.Sprintf("%d/%d invoices have IRN generated", irnCount, total)
	return score, note
}

// scoreMismatchRate: 20 points max.
// Based on GSTR-2A/2B mismatch rate. 0% mismatch = 20, proportional reduction.
func scoreMismatchRate(invoices []APInvoice) (int, string) {
	if len(invoices) == 0 {
		return 10, "No invoices to evaluate; default score applied"
	}

	matchedCount := 0
	for _, inv := range invoices {
		if inv.MismatchStatus == "matched" {
			matchedCount++
		}
	}

	total := len(invoices)
	score := int((float64(matchedCount) / float64(total)) * 20)

	if score > 20 {
		score = 20
	}

	mismatchCount := total - matchedCount
	note := fmt.Sprintf("%d/%d invoices matched, %d mismatched/pending", matchedCount, total, mismatchCount)
	return score, note
}

// scorePaymentBehavior: 15 points max.
// Based on on-time payment rate.
func scorePaymentBehavior(invoices []APInvoice) (int, string) {
	if len(invoices) == 0 {
		return 8, "No invoices to evaluate; default score applied"
	}

	paid := 0
	overdue := 0
	for _, inv := range invoices {
		switch inv.PaymentStatus {
		case "paid":
			paid++
		case "overdue":
			overdue++
		}
	}

	total := len(invoices)
	score := int((float64(paid) / float64(total)) * 15)

	if score > 15 {
		score = 15
	}

	note := fmt.Sprintf("%d/%d paid on time, %d overdue", paid, total, overdue)
	return score, note
}

// scoreDocumentHygiene: 15 points max.
// Based on completeness of GSTIN, PAN, TAN, bank details (email), MSME status.
// Each field is worth 3 points.
func scoreDocumentHygiene(vendor domain.VendorSnapshot) (int, string) {
	score := 0
	missing := []string{}

	if vendor.GSTIN != "" {
		score += 3
	} else {
		missing = append(missing, "GSTIN")
	}

	if vendor.PAN != "" {
		score += 3
	} else {
		missing = append(missing, "PAN")
	}

	if vendor.TAN != "" {
		score += 3
	} else {
		missing = append(missing, "TAN")
	}

	if vendor.Email != "" {
		score += 3
	} else {
		missing = append(missing, "Email/Bank details")
	}

	if vendor.MSMERegistered {
		score += 3
	} else {
		missing = append(missing, "MSME registration")
	}

	if len(missing) == 0 {
		return score, "All document fields complete"
	}
	note := fmt.Sprintf("Missing: %v", missing)
	return score, note
}

// categorize returns A/B/C/D based on total score.
func categorize(total int) string {
	switch {
	case total >= 90:
		return "A"
	case total >= 60:
		return "B"
	case total >= 40:
		return "C"
	default:
		return "D"
	}
}

// riskLevel returns Low/Medium/High/Critical based on total score.
func riskLevel(total int) string {
	switch {
	case total >= 80:
		return "Low"
	case total >= 60:
		return "Medium"
	case total >= 40:
		return "High"
	default:
		return "Critical"
	}
}
