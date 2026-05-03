package domain

import "github.com/shopspring/decimal"

type SalaryInput struct {
	GrossSalary       decimal.Decimal `json:"gross_salary"`
	AllowancesExempt  decimal.Decimal `json:"allowances_exempt"`
	Perquisites       decimal.Decimal `json:"perquisites"`
	ProfitInLieu      decimal.Decimal `json:"profit_in_lieu"`
	TDSByEmployer     decimal.Decimal `json:"tds_by_employer"`
	EmployerTAN       string          `json:"employer_tan"`
}

type SalaryResult struct {
	GrossSalary  decimal.Decimal `json:"gross_salary"`
	Exempt       decimal.Decimal `json:"exempt"`
	NetSalary    decimal.Decimal `json:"net_salary"`
	TDSSection   string          `json:"tds_section"`
}

func ComputeSalaryIncome(in SalaryInput) SalaryResult {
	net := in.GrossSalary.Sub(in.AllowancesExempt).Add(in.Perquisites).Add(in.ProfitInLieu)
	if net.IsNegative() {
		net = zero
	}
	return SalaryResult{
		GrossSalary: in.GrossSalary,
		Exempt:      in.AllowancesExempt,
		NetSalary:   net,
		TDSSection:  "392",
	}
}

type HousePropertyInput struct {
	PropertyCount       int             `json:"property_count"`
	GrossRent           decimal.Decimal `json:"gross_rent"`
	MunicipalTaxPaid    decimal.Decimal `json:"municipal_tax_paid"`
	InterestOnLoan      decimal.Decimal `json:"interest_on_loan"`
	IsSelfOccupied      bool            `json:"is_self_occupied"`
}

type HousePropertyResult struct {
	GrossAnnualValue  decimal.Decimal `json:"gross_annual_value"`
	NetAnnualValue    decimal.Decimal `json:"net_annual_value"`
	StandardDeduction decimal.Decimal `json:"standard_deduction_30pct"`
	InterestDeduction decimal.Decimal `json:"interest_deduction"`
	IncomeFromHP      decimal.Decimal `json:"income_from_hp"`
}

var (
	hpStdDeductionRate     = decimal.NewFromFloat(0.30)
	selfOccupiedInterestMax = decimal.NewFromInt(200000)
)

func ComputeHousePropertyIncome(in HousePropertyInput) HousePropertyResult {
	if in.IsSelfOccupied {
		interestDed := decMin(in.InterestOnLoan, selfOccupiedInterestMax)
		return HousePropertyResult{
			InterestDeduction: interestDed,
			IncomeFromHP:      interestDed.Neg(),
		}
	}

	nav := in.GrossRent.Sub(in.MunicipalTaxPaid)
	if nav.IsNegative() {
		nav = zero
	}
	stdDed := nav.Mul(hpStdDeductionRate).Round(0)
	income := nav.Sub(stdDed).Sub(in.InterestOnLoan)

	return HousePropertyResult{
		GrossAnnualValue:  in.GrossRent,
		NetAnnualValue:    nav,
		StandardDeduction: stdDed,
		InterestDeduction: in.InterestOnLoan,
		IncomeFromHP:      income,
	}
}

type CapitalGainEntry struct {
	AssetType     string          `json:"asset_type"`
	GainType      CapitalGainType `json:"gain_type"`
	SaleAmount    decimal.Decimal `json:"sale_amount"`
	CostBasis     decimal.Decimal `json:"cost_basis"`
	Expenses      decimal.Decimal `json:"expenses"`
	IsSection112A bool            `json:"is_section_112a"`
	IsVDA         bool            `json:"is_vda"`
}

type CapitalGainsResult struct {
	LTCG          decimal.Decimal       `json:"ltcg"`
	STCG          decimal.Decimal       `json:"stcg"`
	LTCG112A      decimal.Decimal       `json:"ltcg_112a"`
	VDAGains      decimal.Decimal       `json:"vda_gains"`
	TotalGains    decimal.Decimal       `json:"total_gains"`
	Entries       []CapitalGainComputed `json:"entries"`
}

type CapitalGainComputed struct {
	AssetType string          `json:"asset_type"`
	GainType  CapitalGainType `json:"gain_type"`
	Gain      decimal.Decimal `json:"gain"`
	Is112A    bool            `json:"is_112a"`
	IsVDA     bool            `json:"is_vda"`
}

var ltcg112AExemption = decimal.NewFromInt(125000)

func ComputeCapitalGains(entries []CapitalGainEntry) CapitalGainsResult {
	result := CapitalGainsResult{}
	for _, e := range entries {
		gain := e.SaleAmount.Sub(e.CostBasis).Sub(e.Expenses)
		computed := CapitalGainComputed{
			AssetType: e.AssetType,
			GainType:  e.GainType,
			Gain:      gain,
			Is112A:    e.IsSection112A,
			IsVDA:     e.IsVDA,
		}
		result.Entries = append(result.Entries, computed)

		if e.IsVDA {
			result.VDAGains = result.VDAGains.Add(gain)
		}
		if e.GainType == LTCG {
			result.LTCG = result.LTCG.Add(gain)
			if e.IsSection112A {
				result.LTCG112A = result.LTCG112A.Add(gain)
			}
		} else {
			result.STCG = result.STCG.Add(gain)
		}
	}

	if result.LTCG112A.GreaterThan(ltcg112AExemption) {
		result.LTCG112A = result.LTCG112A.Sub(ltcg112AExemption)
	} else {
		result.LTCG112A = zero
	}

	result.TotalGains = result.LTCG.Add(result.STCG)
	return result
}

type BusinessInput struct {
	GrossTurnover     decimal.Decimal `json:"gross_turnover"`
	GrossReceipts     decimal.Decimal `json:"gross_receipts"`
	Expenses          decimal.Decimal `json:"expenses"`
	Depreciation      decimal.Decimal `json:"depreciation"`
	Section44BBCApply bool            `json:"section_44bbc_apply"`
}

type BusinessResult struct {
	GrossProfitLoss decimal.Decimal `json:"gross_profit_loss"`
	NetIncome       decimal.Decimal `json:"net_income"`
	Is44BBC         bool            `json:"is_44bbc"`
}

var section44BBCRate = decimal.NewFromFloat(0.10)

func ComputeBusinessIncome(in BusinessInput) BusinessResult {
	if in.Section44BBCApply {
		net := in.GrossReceipts.Mul(section44BBCRate).Round(0)
		return BusinessResult{
			GrossProfitLoss: in.GrossReceipts,
			NetIncome:       net,
			Is44BBC:         true,
		}
	}
	profit := in.GrossTurnover.Add(in.GrossReceipts).Sub(in.Expenses).Sub(in.Depreciation)
	return BusinessResult{
		GrossProfitLoss: in.GrossTurnover.Add(in.GrossReceipts),
		NetIncome:       profit,
	}
}

type OtherSourcesInput struct {
	Interest    decimal.Decimal `json:"interest"`
	Dividends   decimal.Decimal `json:"dividends"`
	FamilyPension decimal.Decimal `json:"family_pension"`
	OtherIncome decimal.Decimal `json:"other_income"`
}

type OtherSourcesResult struct {
	TotalIncome         decimal.Decimal `json:"total_income"`
	FamilyPensionDeduction decimal.Decimal `json:"family_pension_deduction"`
}

var familyPensionDeductionRate = decimal.NewFromFloat(1.0 / 3.0)
var familyPensionDeductionMax = decimal.NewFromInt(15000)

func ComputeOtherSourcesIncome(in OtherSourcesInput) OtherSourcesResult {
	fpDed := zero
	if in.FamilyPension.IsPositive() {
		fpDed = decMin(in.FamilyPension.Mul(familyPensionDeductionRate).Round(0), familyPensionDeductionMax)
	}
	total := in.Interest.Add(in.Dividends).Add(in.FamilyPension).Sub(fpDed).Add(in.OtherIncome)
	return OtherSourcesResult{
		TotalIncome:            total,
		FamilyPensionDeduction: fpDed,
	}
}
