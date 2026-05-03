package domain

import "github.com/shopspring/decimal"

type AISEntry struct {
	SourceType     string          `json:"source_type"`
	DeductorTAN    string          `json:"deductor_tan"`
	DeductorName   string          `json:"deductor_name"`
	Section        string          `json:"section"`
	Amount         decimal.Decimal `json:"amount"`
	TDSAmount      decimal.Decimal `json:"tds_amount"`
}

type TDSCreditEntry struct {
	DeductorTAN  string          `json:"deductor_tan"`
	DeductorName string          `json:"deductor_name"`
	Section      string          `json:"section"`
	GrossPayment decimal.Decimal `json:"gross_payment"`
	TDSAmount    decimal.Decimal `json:"tds_amount"`
}

type ReconciliationResult struct {
	Matched    []ReconciliationMatch `json:"matched"`
	Unmatched  []ReconciliationGap   `json:"unmatched"`
	TotalAIS   decimal.Decimal       `json:"total_ais_tds"`
	TotalClaim decimal.Decimal       `json:"total_claimed_tds"`
	Difference decimal.Decimal       `json:"difference"`
}

type ReconciliationMatch struct {
	DeductorTAN string          `json:"deductor_tan"`
	Section     string          `json:"section"`
	AISAmount   decimal.Decimal `json:"ais_amount"`
	ClaimAmount decimal.Decimal `json:"claim_amount"`
	Discrepancy decimal.Decimal `json:"discrepancy"`
	Status      string          `json:"status"`
}

type ReconciliationGap struct {
	Source      string          `json:"source"`
	DeductorTAN string         `json:"deductor_tan"`
	Section     string          `json:"section"`
	Amount      decimal.Decimal `json:"amount"`
	Issue       string          `json:"issue"`
}

func ReconcileTDS(aisEntries []AISEntry, tdsClaims []TDSCreditEntry) ReconciliationResult {
	result := ReconciliationResult{}

	type key struct {
		TAN     string
		Section string
	}

	aisMap := make(map[key]decimal.Decimal)
	for _, a := range aisEntries {
		k := key{TAN: a.DeductorTAN, Section: a.Section}
		aisMap[k] = aisMap[k].Add(a.TDSAmount)
		result.TotalAIS = result.TotalAIS.Add(a.TDSAmount)
	}

	claimMap := make(map[key]decimal.Decimal)
	for _, c := range tdsClaims {
		k := key{TAN: c.DeductorTAN, Section: c.Section}
		claimMap[k] = claimMap[k].Add(c.TDSAmount)
		result.TotalClaim = result.TotalClaim.Add(c.TDSAmount)
	}

	seen := make(map[key]bool)
	for k, aisAmt := range aisMap {
		seen[k] = true
		claimAmt, found := claimMap[k]
		if !found {
			result.Unmatched = append(result.Unmatched, ReconciliationGap{
				Source:      "AIS",
				DeductorTAN: k.TAN,
				Section:     k.Section,
				Amount:      aisAmt,
				Issue:       "present in AIS (Form 168) but not claimed",
			})
			continue
		}
		disc := aisAmt.Sub(claimAmt)
		status := "MATCHED"
		if !disc.IsZero() {
			status = "DISCREPANCY"
		}
		result.Matched = append(result.Matched, ReconciliationMatch{
			DeductorTAN: k.TAN,
			Section:     k.Section,
			AISAmount:   aisAmt,
			ClaimAmount: claimAmt,
			Discrepancy: disc,
			Status:      status,
		})
	}

	for k, claimAmt := range claimMap {
		if seen[k] {
			continue
		}
		result.Unmatched = append(result.Unmatched, ReconciliationGap{
			Source:      "CLAIM",
			DeductorTAN: k.TAN,
			Section:     k.Section,
			Amount:      claimAmt,
			Issue:       "claimed but not found in AIS (Form 168)",
		})
	}

	result.Difference = result.TotalAIS.Sub(result.TotalClaim)
	return result
}
