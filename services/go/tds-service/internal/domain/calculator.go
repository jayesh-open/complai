package domain

import "github.com/shopspring/decimal"

var (
	zero      = decimal.NewFromInt(0)
	cessRate  = decimal.NewFromFloat(0.04)
	noPANRate = decimal.NewFromFloat(0.20)
)

type CalcInput struct {
	Section        Section
	GrossAmount    decimal.Decimal
	DeducteeType   DeducteeType
	HasValidPAN    bool
	ResidentStatus ResidentStatus
	RentType       RentType
	AggregateForFY decimal.Decimal
	AnnualSalary   decimal.Decimal
	DTAARate       *decimal.Decimal
	LowerCertRate  *decimal.Decimal
}

type CalcResult struct {
	Section      Section         `json:"section"`
	Rate         decimal.Decimal `json:"rate"`
	TDSAmount    decimal.Decimal `json:"tds_amount"`
	Surcharge    decimal.Decimal `json:"surcharge"`
	Cess         decimal.Decimal `json:"cess"`
	TotalTax     decimal.Decimal `json:"total_tax"`
	NoPAN        bool            `json:"no_pan"`
	LowerCert    bool            `json:"lower_cert"`
	ThresholdMet bool            `json:"threshold_met"`
	Explanation  string          `json:"explanation"`
}

func Calculate(in CalcInput) CalcResult {
	switch in.Section {
	case Section192:
		return calcSalary(in)
	case Section194C:
		return calcContractor(in)
	case Section194I:
		return calcRent(in)
	case Section194J:
		return calcProfessional(in)
	case Section194Q:
		return calcPurchase(in)
	case Section195:
		return calcNonResident(in)
	default:
		return CalcResult{Section: in.Section, Explanation: "unknown section"}
	}
}

// FY 2025-26 new tax regime slabs
var salarySlabs = []struct {
	UpTo decimal.Decimal
	Rate decimal.Decimal
}{
	{decimal.NewFromInt(400000), zero},
	{decimal.NewFromInt(800000), decimal.NewFromFloat(0.05)},
	{decimal.NewFromInt(1200000), decimal.NewFromFloat(0.10)},
	{decimal.NewFromInt(1600000), decimal.NewFromFloat(0.15)},
	{decimal.NewFromInt(2000000), decimal.NewFromFloat(0.20)},
	{decimal.NewFromInt(2400000), decimal.NewFromFloat(0.25)},
}

var (
	stdDeduction  = decimal.NewFromInt(75000)
	topRate       = decimal.NewFromFloat(0.30)
	topSlabFloor  = decimal.NewFromInt(2400000)
	rebateLimit   = decimal.NewFromInt(1200000)
	rebateMaximum = decimal.NewFromInt(60000)
	twelve        = decimal.NewFromInt(12)
)

func calcSalary(in CalcInput) CalcResult {
	taxable := in.AnnualSalary.Sub(stdDeduction)
	if taxable.LessThanOrEqual(zero) {
		return CalcResult{Section: Section192, Explanation: "income below standard deduction"}
	}

	tax := zero
	prev := zero
	for _, slab := range salarySlabs {
		if taxable.LessThanOrEqual(prev) {
			break
		}
		bracket := decMin(taxable, slab.UpTo).Sub(prev)
		if bracket.IsPositive() {
			tax = tax.Add(bracket.Mul(slab.Rate))
		}
		prev = slab.UpTo
	}
	if taxable.GreaterThan(topSlabFloor) {
		tax = tax.Add(taxable.Sub(topSlabFloor).Mul(topRate))
	}

	if taxable.LessThanOrEqual(rebateLimit) {
		tax = tax.Sub(decMin(tax, rebateMaximum))
	}

	cess := tax.Mul(cessRate).Round(0)
	total := tax.Add(cess)
	monthlyTDS := total.Div(twelve).Round(0)
	monthlyCess := cess.Div(twelve).Round(0)

	var effectiveRate decimal.Decimal
	if in.AnnualSalary.IsPositive() {
		effectiveRate = total.Div(in.AnnualSalary).Round(4)
	}

	return CalcResult{
		Section:      Section192,
		Rate:         effectiveRate,
		TDSAmount:    monthlyTDS,
		Cess:         monthlyCess,
		TotalTax:     monthlyTDS,
		ThresholdMet: true,
		Explanation:  "monthly TDS on projected annual salary (new tax regime FY 2025-26)",
	}
}

var (
	rate194CIndividual     = decimal.NewFromFloat(0.01)
	rate194COther          = decimal.NewFromFloat(0.02)
	threshold194CSingle    = decimal.NewFromInt(30000)
	threshold194CAggregate = decimal.NewFromInt(100000)
)

func calcContractor(in CalcInput) CalcResult {
	newAgg := in.AggregateForFY.Add(in.GrossAmount)
	if in.GrossAmount.LessThanOrEqual(threshold194CSingle) && newAgg.LessThanOrEqual(threshold194CAggregate) {
		return CalcResult{Section: Section194C, Explanation: "below threshold (single ≤₹30K and aggregate ≤₹1L)"}
	}
	rate := rate194COther
	if in.DeducteeType == DeducteeIndividual || in.DeducteeType == DeducteeHUF {
		rate = rate194CIndividual
	}
	return applyRate(in, rate, "contractor payment")
}

var (
	rate194ILand  = decimal.NewFromFloat(0.10)
	rate194IPlant = decimal.NewFromFloat(0.02)
	threshold194I = decimal.NewFromInt(240000)
)

func calcRent(in CalcInput) CalcResult {
	newAgg := in.AggregateForFY.Add(in.GrossAmount)
	if newAgg.LessThanOrEqual(threshold194I) {
		return CalcResult{Section: Section194I, Explanation: "below annual threshold ₹2.4L"}
	}
	rate := rate194ILand
	if in.RentType == RentPlantMachinery {
		rate = rate194IPlant
	}
	return applyRate(in, rate, "rent payment")
}

var (
	rate194J      = decimal.NewFromFloat(0.10)
	threshold194J = decimal.NewFromInt(30000)
)

func calcProfessional(in CalcInput) CalcResult {
	newAgg := in.AggregateForFY.Add(in.GrossAmount)
	if newAgg.LessThanOrEqual(threshold194J) {
		return CalcResult{Section: Section194J, Explanation: "below annual threshold ₹30K"}
	}
	return applyRate(in, rate194J, "professional fees")
}

var (
	rate194Q      = decimal.NewFromFloat(0.001)
	threshold194Q = decimal.NewFromInt(5000000)
	noPAN194Q     = decimal.NewFromFloat(0.05)
)

func calcPurchase(in CalcInput) CalcResult {
	newAgg := in.AggregateForFY.Add(in.GrossAmount)
	if newAgg.LessThanOrEqual(threshold194Q) {
		return CalcResult{Section: Section194Q, Explanation: "below annual threshold ₹50L"}
	}
	if !in.HasValidPAN {
		tds := in.GrossAmount.Mul(noPAN194Q).Round(0)
		return CalcResult{
			Section: Section194Q, Rate: noPAN194Q, TDSAmount: tds,
			TotalTax: tds, NoPAN: true, ThresholdMet: true,
			Explanation: "purchase of goods — 5% (no valid PAN)",
		}
	}
	return applyRate(in, rate194Q, "purchase of goods")
}

var rate195Default = decimal.NewFromFloat(0.20)

func calcNonResident(in CalcInput) CalcResult {
	rate := rate195Default
	if in.DTAARate != nil {
		rate = *in.DTAARate
	}
	tds := in.GrossAmount.Mul(rate).Round(0)
	cess := tds.Mul(cessRate).Round(0)
	total := tds.Add(cess)
	return CalcResult{
		Section: Section195, Rate: rate, TDSAmount: tds,
		Cess: cess, TotalTax: total, ThresholdMet: true,
		Explanation: "non-resident payment",
	}
}

func applyRate(in CalcInput, baseRate decimal.Decimal, desc string) CalcResult {
	rate := baseRate
	noPAN := false
	lowerCert := false

	if !in.HasValidPAN {
		rate = decMax(baseRate, noPANRate)
		noPAN = true
		desc += " — 20% (no valid PAN)"
	} else if in.LowerCertRate != nil {
		rate = *in.LowerCertRate
		lowerCert = true
		desc += " — lower deduction certificate applied"
	}

	tds := in.GrossAmount.Mul(rate).Round(0)
	return CalcResult{
		Section: in.Section, Rate: rate, TDSAmount: tds,
		TotalTax: tds, NoPAN: noPAN, LowerCert: lowerCert,
		ThresholdMet: true, Explanation: desc,
	}
}

func decMin(a, b decimal.Decimal) decimal.Decimal {
	if a.LessThan(b) {
		return a
	}
	return b
}

func decMax(a, b decimal.Decimal) decimal.Decimal {
	if a.GreaterThan(b) {
		return a
	}
	return b
}
