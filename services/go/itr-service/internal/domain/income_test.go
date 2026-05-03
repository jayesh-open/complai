package domain

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestComputeSalaryIncome_Basic(t *testing.T) {
	result := ComputeSalaryIncome(SalaryInput{
		GrossSalary:      d(1200000),
		AllowancesExempt: d(50000),
		Perquisites:      d(0),
	})
	assert.Equal(t, d(1150000), result.NetSalary)
	assert.Equal(t, "392", result.TDSSection)
}

func TestComputeSalaryIncome_WithPerquisites(t *testing.T) {
	result := ComputeSalaryIncome(SalaryInput{
		GrossSalary:      d(1000000),
		AllowancesExempt: d(100000),
		Perquisites:      d(50000),
		ProfitInLieu:     d(25000),
	})
	assert.Equal(t, d(975000), result.NetSalary)
}

func TestComputeSalaryIncome_NegativeNetClamped(t *testing.T) {
	result := ComputeSalaryIncome(SalaryInput{
		GrossSalary:      d(100000),
		AllowancesExempt: d(200000),
	})
	assert.True(t, result.NetSalary.IsZero())
}

func TestComputeHouseProperty_SelfOccupied(t *testing.T) {
	result := ComputeHousePropertyIncome(HousePropertyInput{
		IsSelfOccupied: true,
		InterestOnLoan: d(250000),
	})
	assert.Equal(t, d(200000), result.InterestDeduction)
	assert.Equal(t, d(-200000), result.IncomeFromHP)
}

func TestComputeHouseProperty_SelfOccupied_InterestBelowMax(t *testing.T) {
	result := ComputeHousePropertyIncome(HousePropertyInput{
		IsSelfOccupied: true,
		InterestOnLoan: d(100000),
	})
	assert.Equal(t, d(100000), result.InterestDeduction)
	assert.Equal(t, d(-100000), result.IncomeFromHP)
}

func TestComputeHouseProperty_LetOut(t *testing.T) {
	result := ComputeHousePropertyIncome(HousePropertyInput{
		GrossRent:        d(600000),
		MunicipalTaxPaid: d(20000),
		InterestOnLoan:   d(150000),
	})
	// NAV = 600K - 20K = 580K
	// StdDed = 30% of 580K = 174000
	// Income = 580K - 174K - 150K = 256000
	assert.Equal(t, d(580000), result.NetAnnualValue)
	assert.Equal(t, d(174000), result.StandardDeduction)
	assert.Equal(t, d(256000), result.IncomeFromHP)
}

func TestComputeCapitalGains_LTCG_and_STCG(t *testing.T) {
	entries := []CapitalGainEntry{
		{AssetType: "equity", GainType: LTCG, SaleAmount: d(500000), CostBasis: d(300000), Expenses: d(1000), IsSection112A: true},
		{AssetType: "debt_mf", GainType: STCG, SaleAmount: d(200000), CostBasis: d(150000), Expenses: d(500)},
	}
	result := ComputeCapitalGains(entries)
	assert.Equal(t, d(199000), result.LTCG)
	assert.Equal(t, d(49500), result.STCG)
	// 112A: 199K - 125K exemption = 74K taxable
	assert.Equal(t, d(74000), result.LTCG112A)
}

func TestComputeCapitalGains_LTCG112A_BelowExemption(t *testing.T) {
	entries := []CapitalGainEntry{
		{AssetType: "equity", GainType: LTCG, SaleAmount: d(200000), CostBasis: d(150000), IsSection112A: true},
	}
	result := ComputeCapitalGains(entries)
	// 50K < 125K exemption → taxable = 0
	assert.True(t, result.LTCG112A.IsZero())
}

func TestComputeCapitalGains_VDA(t *testing.T) {
	entries := []CapitalGainEntry{
		{AssetType: "bitcoin", GainType: STCG, SaleAmount: d(1000000), CostBasis: d(800000), IsVDA: true},
	}
	result := ComputeCapitalGains(entries)
	assert.Equal(t, d(200000), result.VDAGains)
	assert.True(t, result.Entries[0].IsVDA)
}

func TestComputeBusinessIncome_Normal(t *testing.T) {
	result := ComputeBusinessIncome(BusinessInput{
		GrossTurnover: d(5000000),
		Expenses:      d(3500000),
		Depreciation:  d(200000),
	})
	assert.Equal(t, d(1300000), result.NetIncome)
	assert.False(t, result.Is44BBC)
}

func TestComputeBusinessIncome_Section44BBC(t *testing.T) {
	result := ComputeBusinessIncome(BusinessInput{
		GrossReceipts:     d(10000000),
		Section44BBCApply: true,
	})
	// 10% of 1Cr = 10L
	assert.Equal(t, d(1000000), result.NetIncome)
	assert.True(t, result.Is44BBC)
}

func TestComputeOtherSources_WithFamilyPension(t *testing.T) {
	result := ComputeOtherSourcesIncome(OtherSourcesInput{
		Interest:      d(100000),
		Dividends:     d(50000),
		FamilyPension: d(60000),
	})
	// FP deduction = min(60K/3=20K, 15K) = 15K
	expected := d(100000 + 50000 + 60000 - 15000)
	assert.Equal(t, expected, result.TotalIncome)
	assert.Equal(t, d(15000), result.FamilyPensionDeduction)
}

func TestComputeOtherSources_FamilyPensionSmall(t *testing.T) {
	result := ComputeOtherSourcesIncome(OtherSourcesInput{
		FamilyPension: d(30000),
	})
	// FP deduction = min(30K/3=10K, 15K) = 10K
	assert.Equal(t, d(10000), result.FamilyPensionDeduction)
	assert.Equal(t, d(20000), result.TotalIncome)
}

func TestIncomeBreakdown_Total(t *testing.T) {
	ib := IncomeBreakdown{
		Salary:        d(1000000),
		HouseProperty: d(-200000),
		CapitalGains:  d(500000),
		Business:      d(0),
		OtherSources:  d(100000),
	}
	assert.Equal(t, decimal.NewFromInt(1400000), ib.Total())
}

func TestDeductionBreakdown_Total(t *testing.T) {
	db := DeductionBreakdown{
		StandardDeduction: d(75000),
		Section80C:        d(150000),
		Section80D:        d(25000),
		Section24b:        d(200000),
	}
	assert.Equal(t, d(450000), db.Total())
}
