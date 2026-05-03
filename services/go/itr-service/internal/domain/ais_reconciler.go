package domain

import "github.com/shopspring/decimal"

type MismatchSeverity string

const (
	SeverityInfo  MismatchSeverity = "INFO"
	SeverityWarn  MismatchSeverity = "WARN"
	SeverityError MismatchSeverity = "ERROR"
)

type AISSourceData struct {
	PAN               string          `json:"pan"`
	TaxYear           string          `json:"tax_year"`
	SalaryIncome      decimal.Decimal `json:"salary_income"`
	InterestIncome    decimal.Decimal `json:"interest_income"`
	DividendIncome    decimal.Decimal `json:"dividend_income"`
	SecuritiesTrading decimal.Decimal `json:"securities_trading"`
	PropertyTxnValue  decimal.Decimal `json:"property_txn_value"`
	TDSEntries        []AISEntry      `json:"tds_entries"`
}

type BookData struct {
	SalaryIncome      decimal.Decimal  `json:"salary_income"`
	InterestIncome    decimal.Decimal  `json:"interest_income"`
	DividendIncome    decimal.Decimal  `json:"dividend_income"`
	SecuritiesInCG    decimal.Decimal  `json:"securities_in_cg"`
	PropertyInHP      decimal.Decimal  `json:"property_in_hp"`
	TDSClaims         []TDSCreditEntry `json:"tds_claims"`
}

type AISMismatch struct {
	Category           string           `json:"category"`
	Source             string           `json:"source"`
	ExpectedValue      decimal.Decimal  `json:"expected_value"`
	ActualValue        decimal.Decimal  `json:"actual_value"`
	Delta              decimal.Decimal  `json:"delta"`
	Severity           MismatchSeverity `json:"severity"`
	SuggestedAction    string           `json:"suggested_action"`
}

type AISReconcileResult struct {
	PAN                string          `json:"pan"`
	TaxYear            string          `json:"tax_year"`
	Mismatches         []AISMismatch   `json:"mismatches"`
	TDSReconciliation  ReconciliationResult `json:"tds_reconciliation"`
	HasErrors          bool            `json:"has_errors"`
	ErrorCount         int             `json:"error_count"`
	WarnCount          int             `json:"warn_count"`
	InfoCount          int             `json:"info_count"`
	SubmissionBlocked  bool            `json:"submission_blocked"`
}

var (
	salaryThreshold    = decimal.NewFromInt(1)
	tdsThreshold       = decimal.NewFromInt(100)
)

func ReconcileAIS(ais AISSourceData, books BookData, blockOnErrors bool) AISReconcileResult {
	result := AISReconcileResult{
		PAN:     ais.PAN,
		TaxYear: ais.TaxYear,
	}

	result.TDSReconciliation = ReconcileTDS(ais.TDSEntries, books.TDSClaims)

	checkSalary(ais, books, &result)
	checkTDS(result.TDSReconciliation, &result)
	checkInterest(ais, books, &result)
	checkDividend(ais, books, &result)
	checkSecurities(ais, books, &result)
	checkProperty(ais, books, &result)

	for _, m := range result.Mismatches {
		switch m.Severity {
		case SeverityError:
			result.ErrorCount++
		case SeverityWarn:
			result.WarnCount++
		case SeverityInfo:
			result.InfoCount++
		}
	}
	result.HasErrors = result.ErrorCount > 0
	result.SubmissionBlocked = blockOnErrors && result.HasErrors

	return result
}

func checkSalary(ais AISSourceData, books BookData, r *AISReconcileResult) {
	if ais.SalaryIncome.IsZero() && books.SalaryIncome.IsZero() {
		return
	}
	delta := ais.SalaryIncome.Sub(books.SalaryIncome).Abs()
	if delta.GreaterThan(salaryThreshold) {
		r.Mismatches = append(r.Mismatches, AISMismatch{
			Category:        "SALARY",
			Source:           "Form 130 vs AIS",
			ExpectedValue:   ais.SalaryIncome,
			ActualValue:     books.SalaryIncome,
			Delta:           delta,
			Severity:        SeverityError,
			SuggestedAction: "Verify salary amount with Form 130 from employer and update accordingly",
		})
	}
}

func checkTDS(tdsRecon ReconciliationResult, r *AISReconcileResult) {
	if !tdsRecon.Difference.Abs().GreaterThan(tdsThreshold) {
		return
	}

	for _, gap := range tdsRecon.Unmatched {
		sev := SeverityWarn
		action := "Review unmatched TDS entry"
		if gap.Source == "AIS" {
			action = "TDS reported in AIS (Form 168) but not claimed — consider adding TDS credit"
			sev = SeverityError
		} else {
			action = "TDS claimed but not found in AIS (Form 168) — verify with deductor"
			sev = SeverityError
		}
		r.Mismatches = append(r.Mismatches, AISMismatch{
			Category:        "TDS",
			Source:           gap.Source,
			ExpectedValue:   gap.Amount,
			ActualValue:     zero,
			Delta:           gap.Amount,
			Severity:        sev,
			SuggestedAction: action,
		})
	}

	for _, m := range tdsRecon.Matched {
		if m.Status == "DISCREPANCY" && m.Discrepancy.Abs().GreaterThan(tdsThreshold) {
			r.Mismatches = append(r.Mismatches, AISMismatch{
				Category:        "TDS",
				Source:           "AIS vs Claim",
				ExpectedValue:   m.AISAmount,
				ActualValue:     m.ClaimAmount,
				Delta:           m.Discrepancy.Abs(),
				Severity:        SeverityWarn,
				SuggestedAction: "TDS amount differs between AIS and claim for TAN " + m.DeductorTAN,
			})
		}
	}
}

func checkInterest(ais AISSourceData, books BookData, r *AISReconcileResult) {
	if ais.InterestIncome.IsZero() {
		return
	}
	if books.InterestIncome.IsZero() {
		r.Mismatches = append(r.Mismatches, AISMismatch{
			Category:        "INTEREST",
			Source:           "AIS",
			ExpectedValue:   ais.InterestIncome,
			ActualValue:     zero,
			Delta:           ais.InterestIncome,
			Severity:        SeverityError,
			SuggestedAction: "Interest income reported in AIS but not in books — add to Other Sources",
		})
		return
	}
	delta := ais.InterestIncome.Sub(books.InterestIncome).Abs()
	if delta.GreaterThan(salaryThreshold) {
		r.Mismatches = append(r.Mismatches, AISMismatch{
			Category:        "INTEREST",
			Source:           "AIS vs Books",
			ExpectedValue:   ais.InterestIncome,
			ActualValue:     books.InterestIncome,
			Delta:           delta,
			Severity:        SeverityWarn,
			SuggestedAction: "Interest income differs between AIS and books — verify bank statements",
		})
	}
}

func checkDividend(ais AISSourceData, books BookData, r *AISReconcileResult) {
	if ais.DividendIncome.IsZero() {
		return
	}
	if books.DividendIncome.IsZero() {
		addMismatch(r, "DIVIDEND", "AIS", ais.DividendIncome, zero, SeverityError,
			"Dividend income reported in AIS but not in books — add to Other Sources")
		return
	}
	delta := ais.DividendIncome.Sub(books.DividendIncome).Abs()
	if delta.GreaterThan(salaryThreshold) {
		addMismatch(r, "DIVIDEND", "AIS vs Books", ais.DividendIncome, books.DividendIncome, SeverityWarn,
			"Dividend income differs between AIS and books — verify demat/broker statements")
	}
}

func checkSecurities(ais AISSourceData, books BookData, r *AISReconcileResult) {
	if ais.SecuritiesTrading.IsZero() {
		return
	}
	if books.SecuritiesInCG.IsZero() {
		addMismatch(r, "SECURITIES", "AIS", ais.SecuritiesTrading, zero, SeverityError,
			"Securities trading reported in AIS but not in Capital Gains schedule — add to Schedule CG")
		return
	}
	delta := ais.SecuritiesTrading.Sub(books.SecuritiesInCG).Abs()
	if delta.GreaterThan(salaryThreshold) {
		addMismatch(r, "SECURITIES", "AIS vs Books", ais.SecuritiesTrading, books.SecuritiesInCG, SeverityInfo,
			"Securities trading value differs — AIS reports gross, books may differ due to netting")
	}
}

func checkProperty(ais AISSourceData, books BookData, r *AISReconcileResult) {
	if ais.PropertyTxnValue.IsZero() || !books.PropertyInHP.IsZero() {
		return
	}
	addMismatch(r, "PROPERTY", "AIS", ais.PropertyTxnValue, zero, SeverityError,
		"High-value property transaction in AIS not in HP schedule — add to Schedule HP or CG")
}

func addMismatch(r *AISReconcileResult, cat, src string, expected, actual decimal.Decimal, sev MismatchSeverity, action string) {
	r.Mismatches = append(r.Mismatches, AISMismatch{
		Category: cat, Source: src, ExpectedValue: expected, ActualValue: actual,
		Delta: expected.Sub(actual).Abs(), Severity: sev, SuggestedAction: action,
	})
}
