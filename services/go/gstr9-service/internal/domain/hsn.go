package domain

import "github.com/shopspring/decimal"

type HSNEntry struct {
	HSNCode      string          `json:"hsn_code"`
	Description  string          `json:"description"`
	UQC          string          `json:"uqc"`
	Quantity     decimal.Decimal `json:"quantity"`
	TaxableValue decimal.Decimal `json:"taxable_value"`
	CGST         decimal.Decimal `json:"cgst"`
	SGST         decimal.Decimal `json:"sgst"`
	IGST         decimal.Decimal `json:"igst"`
	Cess         decimal.Decimal `json:"cess"`
}

type HSNSummary struct {
	Entries    []HSNEntry      `json:"entries"`
	TotalValue decimal.Decimal `json:"total_value"`
	TotalTax   decimal.Decimal `json:"total_tax"`
}

func BuildHSNSummary(entries []HSNEntry) HSNSummary {
	var totalValue, totalTax decimal.Decimal
	for _, e := range entries {
		totalValue = totalValue.Add(e.TaxableValue)
		totalTax = totalTax.Add(e.CGST).Add(e.SGST).Add(e.IGST).Add(e.Cess)
	}
	return HSNSummary{
		Entries:    entries,
		TotalValue: totalValue,
		TotalTax:   totalTax,
	}
}
