package domain

import "github.com/shopspring/decimal"

var (
	zero          = decimal.NewFromInt(0)
	cessRate      = decimal.NewFromFloat(0.04)
	noPANRate     = decimal.NewFromFloat(0.20)
	noPANPurchase = decimal.NewFromFloat(0.05)
)

type CalcInput struct {
	PaymentCode    PaymentCode
	GrossAmount    decimal.Decimal
	HasValidPAN    bool
	AggregateForFY decimal.Decimal
	AnnualSalary   decimal.Decimal
	DTAARate       *decimal.Decimal
	LowerCertRate  *decimal.Decimal
}

type CalcResult struct {
	Section      Section         `json:"section"`
	PaymentCode  PaymentCode     `json:"payment_code"`
	SubClause    string          `json:"sub_clause,omitempty"`
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
	switch in.PaymentCode {
	case CodeSalaryState, CodeSalaryPrivate, CodeSalaryCentral:
		return calcSalary(in)
	case CodeEPFWithdrawal:
		return calcEPF(in)
	case CodeRentPlant:
		return calcRent(in, rateRentPlant, "rent — plant/machinery")
	case CodeRentLand:
		return calcRent(in, rateRentLand, "rent — land/building")
	case CodeContractorIndiv:
		return calcContractor(in, rateContractorIndiv)
	case CodeContractorOther:
		return calcContractor(in, rateContractorOther)
	case CodeTechnical:
		return calcTechProf(in, rateTechnical, "technical services")
	case CodeProfessional:
		return calcTechProf(in, rateProfessional, "professional services")
	case CodeDirectorRem:
		return calcDirector(in)
	case CodePurchaseGoods:
		return calcPurchase(in)
	case CodeNonResident:
		return calcNonResident(in)
	default:
		return CalcResult{
			Section:     SectionForCode(in.PaymentCode),
			PaymentCode: in.PaymentCode,
			Explanation: "unknown payment code",
		}
	}
}

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
	section := SectionForCode(in.PaymentCode)
	taxable := in.AnnualSalary.Sub(stdDeduction)
	if taxable.LessThanOrEqual(zero) {
		return CalcResult{
			Section:     section,
			PaymentCode: in.PaymentCode,
			Explanation: "income below standard deduction",
		}
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
		Section:      section,
		PaymentCode:  in.PaymentCode,
		Rate:         effectiveRate,
		TDSAmount:    monthlyTDS,
		Cess:         monthlyCess,
		TotalTax:     monthlyTDS,
		ThresholdMet: true,
		Explanation:  "monthly TDS on projected annual salary (new tax regime)",
	}
}

var (
	rateEPF      = decimal.NewFromFloat(0.10)
	thresholdEPF = decimal.NewFromInt(50000)
)

func calcEPF(in CalcInput) CalcResult {
	section := SectionForCode(in.PaymentCode)
	if in.GrossAmount.LessThanOrEqual(thresholdEPF) {
		return CalcResult{
			Section:     section,
			PaymentCode: in.PaymentCode,
			Explanation: "below threshold ₹50K",
		}
	}
	return applyRate(in, rateEPF, "EPF withdrawal")
}

var (
	rateRentPlant = decimal.NewFromFloat(0.02)
	rateRentLand  = decimal.NewFromFloat(0.10)
	thresholdRent = decimal.NewFromInt(50000)
)

func calcRent(in CalcInput, baseRate decimal.Decimal, desc string) CalcResult {
	section := SectionForCode(in.PaymentCode)
	if in.GrossAmount.LessThanOrEqual(thresholdRent) {
		return CalcResult{
			Section:     section,
			PaymentCode: in.PaymentCode,
			Explanation: "below per-payment threshold ₹50K",
		}
	}
	return applyRate(in, baseRate, desc)
}

var (
	rateContractorIndiv        = decimal.NewFromFloat(0.01)
	rateContractorOther        = decimal.NewFromFloat(0.02)
	thresholdContractSingle    = decimal.NewFromInt(30000)
	thresholdContractAggregate = decimal.NewFromInt(100000)
)

func calcContractor(in CalcInput, baseRate decimal.Decimal) CalcResult {
	section := SectionForCode(in.PaymentCode)
	newAgg := in.AggregateForFY.Add(in.GrossAmount)
	if in.GrossAmount.LessThanOrEqual(thresholdContractSingle) && newAgg.LessThanOrEqual(thresholdContractAggregate) {
		return CalcResult{
			Section:     section,
			PaymentCode: in.PaymentCode,
			Explanation: "below threshold (single ≤₹30K and aggregate ≤₹1L)",
		}
	}
	return applyRate(in, baseRate, "contractor payment")
}

var (
	rateTechnical    = decimal.NewFromFloat(0.02)
	rateProfessional = decimal.NewFromFloat(0.10)
	thresholdTechProf = decimal.NewFromInt(50000)
)

func calcTechProf(in CalcInput, baseRate decimal.Decimal, desc string) CalcResult {
	section := SectionForCode(in.PaymentCode)
	newAgg := in.AggregateForFY.Add(in.GrossAmount)
	if newAgg.LessThanOrEqual(thresholdTechProf) {
		return CalcResult{
			Section:     section,
			PaymentCode: in.PaymentCode,
			Explanation: "below annual threshold ₹50K",
		}
	}
	return applyRate(in, baseRate, desc)
}

var rateDirector = decimal.NewFromFloat(0.10)

func calcDirector(in CalcInput) CalcResult {
	return applyRate(in, rateDirector, "director remuneration")
}

var (
	ratePurchase      = decimal.NewFromFloat(0.001)
	thresholdPurchase = decimal.NewFromInt(5000000)
)

func calcPurchase(in CalcInput) CalcResult {
	section := SectionForCode(in.PaymentCode)
	newAgg := in.AggregateForFY.Add(in.GrossAmount)
	if newAgg.LessThanOrEqual(thresholdPurchase) {
		return CalcResult{
			Section:     section,
			PaymentCode: in.PaymentCode,
			Explanation: "below annual threshold ₹50L",
		}
	}
	if !in.HasValidPAN {
		tds := in.GrossAmount.Mul(noPANPurchase).Round(0)
		return CalcResult{
			Section:      section,
			PaymentCode:  in.PaymentCode,
			SubClause:    SubClauseForCode(in.PaymentCode),
			Rate:         noPANPurchase,
			TDSAmount:    tds,
			TotalTax:     tds,
			NoPAN:        true,
			ThresholdMet: true,
			Explanation:  "purchase of goods — 5% (no valid PAN, s.397(2))",
		}
	}
	return applyRate(in, ratePurchase, "purchase of goods")
}

var rateNonResidentDefault = decimal.NewFromFloat(0.20)

func calcNonResident(in CalcInput) CalcResult {
	section := SectionForCode(in.PaymentCode)
	rate := rateNonResidentDefault
	if in.DTAARate != nil {
		rate = *in.DTAARate
	}
	tds := in.GrossAmount.Mul(rate).Round(0)
	cess := tds.Mul(cessRate).Round(0)
	total := tds.Add(cess)
	return CalcResult{
		Section:      section,
		PaymentCode:  in.PaymentCode,
		SubClause:    SubClauseForCode(in.PaymentCode),
		Rate:         rate,
		TDSAmount:    tds,
		Cess:         cess,
		TotalTax:     total,
		ThresholdMet: true,
		Explanation:  "non-resident payment (s.393(2))",
	}
}

func applyRate(in CalcInput, baseRate decimal.Decimal, desc string) CalcResult {
	section := SectionForCode(in.PaymentCode)
	rate := baseRate
	noPAN := false
	lowerCert := false

	if !in.HasValidPAN {
		rate = decMax(baseRate, noPANRate)
		noPAN = true
		desc += " — 20% (no valid PAN, s.397(2))"
	} else if in.LowerCertRate != nil {
		rate = *in.LowerCertRate
		lowerCert = true
		desc += " — lower deduction certificate applied"
	}

	tds := in.GrossAmount.Mul(rate).Round(0)
	return CalcResult{
		Section:      section,
		PaymentCode:  in.PaymentCode,
		SubClause:    SubClauseForCode(in.PaymentCode),
		Rate:         rate,
		TDSAmount:    tds,
		TotalTax:     tds,
		NoPAN:        noPAN,
		LowerCert:    lowerCert,
		ThresholdMet: true,
		Explanation:  desc,
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
