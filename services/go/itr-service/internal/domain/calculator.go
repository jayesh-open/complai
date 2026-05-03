package domain

import "github.com/shopspring/decimal"

var (
	zero    = decimal.NewFromInt(0)
	hundred = decimal.NewFromInt(100)

	cessRate    = decimal.NewFromFloat(0.04)
	rebate87ANewRegimeLimit = decimal.NewFromInt(700000)
)

type IncomeBreakdown struct {
	Salary        decimal.Decimal `json:"salary"`
	HouseProperty decimal.Decimal `json:"house_property"`
	CapitalGains  decimal.Decimal `json:"capital_gains"`
	Business      decimal.Decimal `json:"business"`
	OtherSources  decimal.Decimal `json:"other_sources"`
}

func (ib IncomeBreakdown) Total() decimal.Decimal {
	return ib.Salary.Add(ib.HouseProperty).Add(ib.CapitalGains).Add(ib.Business).Add(ib.OtherSources)
}

type DeductionBreakdown struct {
	StandardDeduction decimal.Decimal `json:"standard_deduction"`
	Section80C        decimal.Decimal `json:"section_80c,omitempty"`
	Section80D        decimal.Decimal `json:"section_80d,omitempty"`
	Section24b        decimal.Decimal `json:"section_24b,omitempty"`
	OtherVI_A         decimal.Decimal `json:"other_vi_a,omitempty"`
}

func (db DeductionBreakdown) Total() decimal.Decimal {
	return db.StandardDeduction.Add(db.Section80C).Add(db.Section80D).Add(db.Section24b).Add(db.OtherVI_A)
}

type TaxComputeInput struct {
	Income     IncomeBreakdown
	Deductions DeductionBreakdown
	Regime     RegimeType
	IsResident bool
	TDSCredits decimal.Decimal
	AdvanceTax decimal.Decimal
}

type TaxComputeResult struct {
	GrossIncome       decimal.Decimal `json:"gross_income"`
	StandardDeduction decimal.Decimal `json:"standard_deduction"`
	TotalDeductions   decimal.Decimal `json:"total_deductions"`
	TaxableIncome     decimal.Decimal `json:"taxable_income"`
	BaseTax           decimal.Decimal `json:"base_tax"`
	Surcharge         decimal.Decimal `json:"surcharge"`
	SurchargeRate     decimal.Decimal `json:"surcharge_rate"`
	HealthEdCess      decimal.Decimal `json:"health_ed_cess"`
	Rebate87A         decimal.Decimal `json:"rebate_87a"`
	GrossTaxPayable   decimal.Decimal `json:"gross_tax_payable"`
	TDSCredit         decimal.Decimal `json:"tds_credit"`
	AdvanceTax        decimal.Decimal `json:"advance_tax"`
	NetTaxPayable     decimal.Decimal `json:"net_tax_payable"`
	RefundDue         decimal.Decimal `json:"refund_due"`
	Regime            RegimeType      `json:"regime"`
	SlabBreakdown     []SlabDetail    `json:"slab_breakdown"`
}

type SlabDetail struct {
	From decimal.Decimal `json:"from"`
	To   decimal.Decimal `json:"to"`
	Rate decimal.Decimal `json:"rate"`
	Tax  decimal.Decimal `json:"tax"`
}

var newRegimeSlabs = []struct {
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

var newRegimeTopRate = decimal.NewFromFloat(0.30)
var newRegimeTopFloor = decimal.NewFromInt(2400000)

var oldRegimeSlabs = []struct {
	UpTo decimal.Decimal
	Rate decimal.Decimal
}{
	{decimal.NewFromInt(300000), zero},
	{decimal.NewFromInt(700000), decimal.NewFromFloat(0.05)},
	{decimal.NewFromInt(1000000), decimal.NewFromFloat(0.10)},
	{decimal.NewFromInt(1200000), decimal.NewFromFloat(0.15)},
	{decimal.NewFromInt(1500000), decimal.NewFromFloat(0.20)},
}

var oldRegimeTopRate = decimal.NewFromFloat(0.30)
var oldRegimeTopFloor = decimal.NewFromInt(1500000)

var (
	stdDeductionNew = decimal.NewFromInt(75000)
	stdDeductionOld = decimal.NewFromInt(50000)
)

var (
	surcharge10Threshold = decimal.NewFromInt(5000000)
	surcharge15Threshold = decimal.NewFromInt(10000000)
	surcharge25Threshold = decimal.NewFromInt(20000000)
	surcharge10Rate      = decimal.NewFromFloat(0.10)
	surcharge15Rate      = decimal.NewFromFloat(0.15)
	surcharge25Rate      = decimal.NewFromFloat(0.25)
)

func ComputeTax(in TaxComputeInput) TaxComputeResult {
	gross := in.Income.Total()

	var stdDed decimal.Decimal
	if in.Income.Salary.IsPositive() {
		if in.Regime == NewRegime {
			stdDed = stdDeductionNew
		} else {
			stdDed = stdDeductionOld
		}
	}

	var totalDed decimal.Decimal
	if in.Regime == NewRegime {
		totalDed = stdDed
	} else {
		totalDed = stdDed.Add(in.Deductions.Section80C).Add(in.Deductions.Section80D).Add(in.Deductions.Section24b).Add(in.Deductions.OtherVI_A)
	}

	taxable := gross.Sub(totalDed)
	if taxable.IsNegative() {
		taxable = zero
	}

	var baseTax decimal.Decimal
	var slabs []SlabDetail

	if in.Regime == NewRegime {
		baseTax, slabs = computeSlabTax(taxable, newRegimeSlabs, newRegimeTopRate, newRegimeTopFloor)
	} else {
		baseTax, slabs = computeSlabTax(taxable, oldRegimeSlabs, oldRegimeTopRate, oldRegimeTopFloor)
	}

	var rebate decimal.Decimal
	if in.Regime == NewRegime && in.IsResident && taxable.LessThanOrEqual(rebate87ANewRegimeLimit) {
		rebate = baseTax
	}

	afterRebate := baseTax.Sub(rebate)
	if afterRebate.IsNegative() {
		afterRebate = zero
	}

	surchargeRate, surcharge := computeSurcharge(afterRebate, taxable, in.Regime)

	taxPlusSurcharge := afterRebate.Add(surcharge)
	cess := taxPlusSurcharge.Mul(cessRate).Round(0)

	grossTax := taxPlusSurcharge.Add(cess)

	netPayable := grossTax.Sub(in.TDSCredits).Sub(in.AdvanceTax)
	var refund decimal.Decimal
	if netPayable.IsNegative() {
		refund = netPayable.Abs()
		netPayable = zero
	}

	return TaxComputeResult{
		GrossIncome:       gross,
		StandardDeduction: stdDed,
		TotalDeductions:   totalDed,
		TaxableIncome:     taxable,
		BaseTax:           baseTax,
		Surcharge:         surcharge,
		SurchargeRate:     surchargeRate,
		HealthEdCess:      cess,
		Rebate87A:         rebate,
		GrossTaxPayable:   grossTax,
		TDSCredit:         in.TDSCredits,
		AdvanceTax:        in.AdvanceTax,
		NetTaxPayable:     netPayable,
		RefundDue:         refund,
		Regime:            in.Regime,
		SlabBreakdown:     slabs,
	}
}

func computeSlabTax(taxable decimal.Decimal, slabs []struct {
	UpTo decimal.Decimal
	Rate decimal.Decimal
}, topRate, topFloor decimal.Decimal) (decimal.Decimal, []SlabDetail) {
	tax := zero
	prev := zero
	var details []SlabDetail

	for _, slab := range slabs {
		if taxable.LessThanOrEqual(prev) {
			break
		}
		bracket := decMin(taxable, slab.UpTo).Sub(prev)
		if bracket.IsPositive() {
			slabTax := bracket.Mul(slab.Rate).Round(0)
			tax = tax.Add(slabTax)
			details = append(details, SlabDetail{
				From: prev,
				To:   decMin(taxable, slab.UpTo),
				Rate: slab.Rate.Mul(hundred),
				Tax:  slabTax,
			})
		}
		prev = slab.UpTo
	}

	if taxable.GreaterThan(topFloor) {
		excess := taxable.Sub(topFloor)
		slabTax := excess.Mul(topRate).Round(0)
		tax = tax.Add(slabTax)
		details = append(details, SlabDetail{
			From: topFloor,
			To:   taxable,
			Rate: topRate.Mul(hundred),
			Tax:  slabTax,
		})
	}

	return tax, details
}

func computeSurcharge(taxAfterRebate, taxableIncome decimal.Decimal, regime RegimeType) (decimal.Decimal, decimal.Decimal) {
	if taxAfterRebate.IsZero() {
		return zero, zero
	}

	var rate decimal.Decimal
	if taxableIncome.GreaterThan(surcharge25Threshold) {
		rate = surcharge25Rate
		if regime == NewRegime {
			rate = surcharge25Rate
		}
	} else if taxableIncome.GreaterThan(surcharge15Threshold) {
		rate = surcharge15Rate
	} else if taxableIncome.GreaterThan(surcharge10Threshold) {
		rate = surcharge10Rate
	} else {
		return zero, zero
	}

	surcharge := taxAfterRebate.Mul(rate).Round(0)

	if regime == NewRegime && rate.GreaterThan(surcharge25Rate) {
		rate = surcharge25Rate
		surcharge = taxAfterRebate.Mul(rate).Round(0)
	}

	return rate, surcharge
}

func decMin(a, b decimal.Decimal) decimal.Decimal {
	if a.LessThan(b) {
		return a
	}
	return b
}
