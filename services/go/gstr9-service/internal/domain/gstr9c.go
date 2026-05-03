package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	SeverityInfoThreshold  = decimal.NewFromInt(100)
	SeverityWarnThreshold  = decimal.NewFromInt(10000)
)

type AuditedFinancials struct {
	Turnover           decimal.Decimal `json:"turnover"`
	TaxPayable         TaxBreakdown    `json:"tax_payable"`
	ITCClaimed         ITCBreakdown    `json:"itc_claimed"`
	UnbilledRevenue    decimal.Decimal `json:"unbilled_revenue"`
	UnadjustedAdvances decimal.Decimal `json:"unadjusted_advances"`
	DeemedSupply       decimal.Decimal `json:"deemed_supply"`
	CreditNotesAfterFY decimal.Decimal `json:"credit_notes_after_fy"`
	Sec15_3Adjustments decimal.Decimal `json:"sec_15_3_adjustments"`
}

type ReconciliationResult struct {
	GSTR9CFiling GSTR9CFiling   `json:"filing"`
	Mismatches   []GSTR9CMismatch `json:"mismatches"`
	CanSubmit    bool             `json:"can_submit"`
}

func ClassifySeverity(diff decimal.Decimal) MismatchSeverity {
	absDiff := diff.Abs()
	if absDiff.LessThan(SeverityInfoThreshold) {
		return SeverityInfo
	}
	if absDiff.LessThan(SeverityWarnThreshold) {
		return SeverityWarn
	}
	return SeverityError
}

func CanSubmit(mismatches []GSTR9CMismatch) bool {
	for _, m := range mismatches {
		if m.Severity == SeverityError && !m.Resolved {
			return false
		}
	}
	return true
}

func Reconcile(
	gstr9Filing *GSTR9Filing,
	tables []GSTR9TableData,
	audited AuditedFinancials,
	tenantID uuid.UUID,
	gstr9cFilingID uuid.UUID,
) []GSTR9CMismatch {
	var mismatches []GSTR9CMismatch
	now := time.Now()

	mismatches = append(mismatches, reconcileTurnover(gstr9Filing, tables, audited, tenantID, gstr9cFilingID, now)...)
	mismatches = append(mismatches, reconcileTax(tables, audited, tenantID, gstr9cFilingID, now)...)
	mismatches = append(mismatches, reconcileITC(tables, audited, tenantID, gstr9cFilingID, now)...)

	return mismatches
}

func reconcileTurnover(
	gstr9Filing *GSTR9Filing,
	tables []GSTR9TableData,
	audited AuditedFinancials,
	tenantID uuid.UUID,
	gstr9cFilingID uuid.UUID,
	now time.Time,
) []GSTR9CMismatch {
	var mismatches []GSTR9CMismatch

	gstr9Turnover := ComputeAggregateTurnover(tables)
	diff := audited.Turnover.Sub(gstr9Turnover)
	if !diff.IsZero() {
		mismatches = append(mismatches, GSTR9CMismatch{
			ID: uuid.New(), TenantID: tenantID, GSTR9CFilingID: gstr9cFilingID,
			Section: "II", Category: "turnover", Description: "Aggregate turnover mismatch",
			BooksAmount: audited.Turnover, GSTR9Amount: gstr9Turnover,
			Difference: diff, Severity: ClassifySeverity(diff),
			Reason: "Audited annual turnover differs from GSTR-9 aggregate turnover",
			SuggestedAction: "Verify turnover figures in books of accounts vs GSTR-9 Tables 4+5",
			CreatedAt: now,
		})
	}

	if !audited.UnbilledRevenue.IsZero() {
		mismatches = append(mismatches, GSTR9CMismatch{
			ID: uuid.New(), TenantID: tenantID, GSTR9CFilingID: gstr9cFilingID,
			Section: "II", Category: "turnover",
			Description: "Unbilled revenue not reflected in GSTR-9",
			BooksAmount: audited.UnbilledRevenue, GSTR9Amount: decimal.Zero,
			Difference: audited.UnbilledRevenue,
			Severity: ClassifySeverity(audited.UnbilledRevenue),
			Reason: "Revenue recognized in books but invoice not issued before FY end",
			SuggestedAction: "Include in next FY GSTR-1 when invoiced; no amendment needed if disclosed",
			CreatedAt: now,
		})
	}

	if !audited.UnadjustedAdvances.IsZero() {
		mismatches = append(mismatches, GSTR9CMismatch{
			ID: uuid.New(), TenantID: tenantID, GSTR9CFilingID: gstr9cFilingID,
			Section: "II", Category: "turnover",
			Description: "Unadjusted advances",
			BooksAmount: audited.UnadjustedAdvances, GSTR9Amount: decimal.Zero,
			Difference: audited.UnadjustedAdvances,
			Severity: ClassifySeverity(audited.UnadjustedAdvances),
			Reason:          "Advances received on which GST paid but supply not yet made",
			SuggestedAction: "Adjust against future invoices; report in Table 10/11 if applicable",
			CreatedAt: now,
		})
	}

	if !audited.DeemedSupply.IsZero() {
		mismatches = append(mismatches, GSTR9CMismatch{
			ID: uuid.New(), TenantID: tenantID, GSTR9CFilingID: gstr9cFilingID,
			Section: "II", Category: "turnover",
			Description: "Deemed supply under Schedule I",
			BooksAmount: audited.DeemedSupply, GSTR9Amount: decimal.Zero,
			Difference: audited.DeemedSupply,
			Severity: ClassifySeverity(audited.DeemedSupply),
			Reason:          "Supply between related persons / distinct persons deemed as supply",
			SuggestedAction: "Ensure deemed supplies are reported in GSTR-1 and reflected in Table 4",
			CreatedAt: now,
		})
	}

	if !audited.CreditNotesAfterFY.IsZero() {
		mismatches = append(mismatches, GSTR9CMismatch{
			ID: uuid.New(), TenantID: tenantID, GSTR9CFilingID: gstr9cFilingID,
			Section: "II", Category: "turnover",
			Description: "Credit notes issued after FY end",
			BooksAmount: decimal.Zero, GSTR9Amount: audited.CreditNotesAfterFY,
			Difference: audited.CreditNotesAfterFY.Neg(),
			Severity: ClassifySeverity(audited.CreditNotesAfterFY),
			Reason:          "Credit notes issued in the next FY pertaining to current FY supplies",
			SuggestedAction: "Report in Table 11 (amendments reducing turnover) of next FY GSTR-9",
			CreatedAt: now,
		})
	}

	if !audited.Sec15_3Adjustments.IsZero() {
		mismatches = append(mismatches, GSTR9CMismatch{
			ID: uuid.New(), TenantID: tenantID, GSTR9CFilingID: gstr9cFilingID,
			Section: "II", Category: "turnover",
			Description: "Adjustments per Section 15(3)",
			BooksAmount: audited.Sec15_3Adjustments, GSTR9Amount: decimal.Zero,
			Difference: audited.Sec15_3Adjustments,
			Severity: ClassifySeverity(audited.Sec15_3Adjustments),
			Reason:          "Value of supply adjusted as per Section 15(3) conditions",
			SuggestedAction: "Verify discount/adjustment eligibility under Section 15(3)",
			CreatedAt: now,
		})
	}

	return mismatches
}

func reconcileTax(
	tables []GSTR9TableData,
	audited AuditedFinancials,
	tenantID uuid.UUID,
	gstr9cFilingID uuid.UUID,
	now time.Time,
) []GSTR9CMismatch {
	var mismatches []GSTR9CMismatch

	var gstr9Tax TaxBreakdown
	for _, t := range tables {
		if t.TableNumber == "9" {
			gstr9Tax = TaxBreakdown{
				CGST: t.CGST, SGST: t.SGST, IGST: t.IGST, Cess: t.Cess,
			}
			break
		}
	}

	type taxCheck struct {
		name      string
		books     decimal.Decimal
		gstr9     decimal.Decimal
	}
	checks := []taxCheck{
		{"CGST", audited.TaxPayable.CGST, gstr9Tax.CGST},
		{"SGST", audited.TaxPayable.SGST, gstr9Tax.SGST},
		{"IGST", audited.TaxPayable.IGST, gstr9Tax.IGST},
		{"Cess", audited.TaxPayable.Cess, gstr9Tax.Cess},
	}

	for _, c := range checks {
		diff := c.books.Sub(c.gstr9)
		if !diff.IsZero() {
			mismatches = append(mismatches, GSTR9CMismatch{
				ID: uuid.New(), TenantID: tenantID, GSTR9CFilingID: gstr9cFilingID,
				Section: "III", Category: "tax",
				Description: c.name + " tax payable mismatch",
				BooksAmount: c.books, GSTR9Amount: c.gstr9,
				Difference: diff, Severity: ClassifySeverity(diff),
				Reason:          "Tax payable per audited books differs from GSTR-9 Table 9 " + c.name,
				SuggestedAction: "Reconcile " + c.name + " payable with GSTR-3B challan payments",
				CreatedAt: now,
			})
		}
	}

	return mismatches
}

func reconcileITC(
	tables []GSTR9TableData,
	audited AuditedFinancials,
	tenantID uuid.UUID,
	gstr9cFilingID uuid.UUID,
	now time.Time,
) []GSTR9CMismatch {
	var mismatches []GSTR9CMismatch

	var gstr9ITC ITCBreakdown
	for _, t := range tables {
		if t.TableNumber == "6F" {
			gstr9ITC = ITCBreakdown{
				CGST: t.CGST, SGST: t.SGST, IGST: t.IGST, Cess: t.Cess,
			}
			break
		}
	}

	type itcCheck struct {
		name  string
		books decimal.Decimal
		gstr9 decimal.Decimal
	}
	checks := []itcCheck{
		{"CGST", audited.ITCClaimed.CGST, gstr9ITC.CGST},
		{"SGST", audited.ITCClaimed.SGST, gstr9ITC.SGST},
		{"IGST", audited.ITCClaimed.IGST, gstr9ITC.IGST},
		{"Cess", audited.ITCClaimed.Cess, gstr9ITC.Cess},
	}

	for _, c := range checks {
		diff := c.books.Sub(c.gstr9)
		if !diff.IsZero() {
			sev := ClassifySeverity(diff)
			action := "Reconcile " + c.name + " ITC with GSTR-2B and books"
			reason := c.name + " ITC claimed in books differs from GSTR-9 Table 6F"

			if diff.LessThan(decimal.Zero) {
				action = "Excess ITC claimed in GSTR-9 — reverse via DRC-03 or adjust in next period"
				reason = c.name + " ITC in GSTR-9 exceeds audited books — potential excess claim"
			} else if diff.GreaterThan(decimal.Zero) {
				action = "Short ITC claimed — opportunity to reclaim if within time limit"
				reason = c.name + " ITC in books exceeds GSTR-9 — potential missed claim"
			}

			mismatches = append(mismatches, GSTR9CMismatch{
				ID: uuid.New(), TenantID: tenantID, GSTR9CFilingID: gstr9cFilingID,
				Section: "IV", Category: "itc",
				Description: c.name + " ITC mismatch",
				BooksAmount: c.books, GSTR9Amount: c.gstr9,
				Difference: diff, Severity: sev,
				Reason: reason, SuggestedAction: action,
				CreatedAt: now,
			})
		}
	}

	return mismatches
}
