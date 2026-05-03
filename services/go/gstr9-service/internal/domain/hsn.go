package domain

import "github.com/shopspring/decimal"

var (
	HSN6DigitThreshold = decimal.NewFromInt(50000000)
	HSN4DigitThreshold = decimal.NewFromInt(15000000)
)

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
	DigitLevel int             `json:"digit_level"`
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

func HSNDigitLevel(turnover decimal.Decimal) int {
	if turnover.GreaterThan(HSN6DigitThreshold) {
		return 6
	}
	if turnover.GreaterThan(HSN4DigitThreshold) {
		return 4
	}
	return 2
}

func truncateHSN(code string, digits int) string {
	if len(code) <= digits {
		return code
	}
	return code[:digits]
}

func AggregateHSNByTurnover(entries []HSNEntry, turnover decimal.Decimal) HSNSummary {
	digits := HSNDigitLevel(turnover)

	grouped := make(map[string]*HSNEntry)
	for _, e := range entries {
		key := truncateHSN(e.HSNCode, digits)
		if existing, ok := grouped[key]; ok {
			existing.Quantity = existing.Quantity.Add(e.Quantity)
			existing.TaxableValue = existing.TaxableValue.Add(e.TaxableValue)
			existing.CGST = existing.CGST.Add(e.CGST)
			existing.SGST = existing.SGST.Add(e.SGST)
			existing.IGST = existing.IGST.Add(e.IGST)
			existing.Cess = existing.Cess.Add(e.Cess)
		} else {
			copy := e
			copy.HSNCode = key
			grouped[key] = &copy
		}
	}

	var result []HSNEntry
	var totalValue, totalTax decimal.Decimal
	for _, e := range grouped {
		result = append(result, *e)
		totalValue = totalValue.Add(e.TaxableValue)
		totalTax = totalTax.Add(e.CGST).Add(e.SGST).Add(e.IGST).Add(e.Cess)
	}

	return HSNSummary{
		Entries:    result,
		TotalValue: totalValue,
		TotalTax:   totalTax,
		DigitLevel: digits,
	}
}
