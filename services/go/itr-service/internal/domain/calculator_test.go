package domain

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func d(v int64) decimal.Decimal     { return decimal.NewFromInt(v) }
func df(v float64) decimal.Decimal  { return decimal.NewFromFloat(v) }

func TestComputeTax_NewRegime_ZeroIncome(t *testing.T) {
	result := ComputeTax(TaxComputeInput{
		Income: IncomeBreakdown{Salary: d(0)},
		Regime: NewRegime,
	})
	assert.True(t, result.BaseTax.IsZero())
	assert.True(t, result.GrossTaxPayable.IsZero())
	assert.Equal(t, NewRegime, result.Regime)
}

func TestComputeTax_NewRegime_BelowStdDeduction(t *testing.T) {
	result := ComputeTax(TaxComputeInput{
		Income: IncomeBreakdown{Salary: d(50000)},
		Regime: NewRegime,
	})
	assert.True(t, result.TaxableIncome.IsZero())
	assert.True(t, result.BaseTax.IsZero())
	assert.Equal(t, d(75000), result.StandardDeduction)
}

func TestComputeTax_NewRegime_500K_Salary(t *testing.T) {
	result := ComputeTax(TaxComputeInput{
		Income:     IncomeBreakdown{Salary: d(500000)},
		Regime:     NewRegime,
		IsResident: true,
	})
	// 500K - 75K std = 425K taxable
	// 0-4L = 0, 4L-4.25L = 25K * 5% = 1250
	assert.Equal(t, d(425000), result.TaxableIncome)
	assert.Equal(t, d(1250), result.BaseTax)
	// Under 7L → 87A rebate applies
	assert.Equal(t, d(1250), result.Rebate87A)
	assert.True(t, result.GrossTaxPayable.IsZero())
}

func TestComputeTax_NewRegime_700K_Rebate87A_Boundary(t *testing.T) {
	// 7L taxable = exactly at rebate boundary
	gross := d(775000) // 775K - 75K std = 700K taxable
	result := ComputeTax(TaxComputeInput{
		Income:     IncomeBreakdown{Salary: gross},
		Regime:     NewRegime,
		IsResident: true,
	})
	assert.Equal(t, d(700000), result.TaxableIncome)
	// 0-4L=0, 4-7L=3L*5%=15000, 7L-8L not reached
	// Actually: 0-4L=0, 4L-7L: 4L-8L slab at 5% → only up to 7L → 3L*5%=15000
	assert.Equal(t, d(15000), result.BaseTax)
	assert.Equal(t, d(15000), result.Rebate87A)
	assert.True(t, result.GrossTaxPayable.IsZero())
}

func TestComputeTax_NewRegime_700001_NoRebate(t *testing.T) {
	gross := d(775001) // 775001 - 75K = 700001 → just above rebate limit
	result := ComputeTax(TaxComputeInput{
		Income:     IncomeBreakdown{Salary: gross},
		Regime:     NewRegime,
		IsResident: true,
	})
	assert.Equal(t, d(700001), result.TaxableIncome)
	assert.True(t, result.Rebate87A.IsZero(), "should NOT get rebate above ₹7L")
	assert.True(t, result.GrossTaxPayable.IsPositive())
}

func TestComputeTax_NewRegime_1200K_AllSlabs(t *testing.T) {
	gross := d(1275000) // 1275K - 75K = 1200K taxable
	result := ComputeTax(TaxComputeInput{
		Income:     IncomeBreakdown{Salary: gross},
		Regime:     NewRegime,
		IsResident: true,
	})
	assert.Equal(t, d(1200000), result.TaxableIncome)
	// 0-4L=0, 4-8L=4L*5%=20000, 8-12L=4L*10%=40000
	assert.Equal(t, d(60000), result.BaseTax)
	assert.True(t, result.Rebate87A.IsZero())
	cess := d(60000).Mul(df(0.04)).Round(0)
	assert.Equal(t, cess, result.HealthEdCess)
}

func TestComputeTax_NewRegime_2500K_TopSlab(t *testing.T) {
	gross := d(2575000) // 2575K - 75K = 2500K
	result := ComputeTax(TaxComputeInput{
		Income:     IncomeBreakdown{Salary: gross},
		Regime:     NewRegime,
		IsResident: true,
	})
	assert.Equal(t, d(2500000), result.TaxableIncome)
	// 0-4L=0, 4-8L=20000, 8-12L=40000, 12-16L=60000, 16-20L=80000, 20-24L=100000, 24-25L=1L*30%=30000
	expected := d(20000 + 40000 + 60000 + 80000 + 100000 + 30000)
	assert.Equal(t, expected, result.BaseTax)
}

func TestComputeTax_NewRegime_HEC_4pct(t *testing.T) {
	result := ComputeTax(TaxComputeInput{
		Income:     IncomeBreakdown{Salary: d(1275000)},
		Regime:     NewRegime,
		IsResident: true,
	})
	assert.Equal(t, d(60000), result.BaseTax)
	assert.Equal(t, d(2400), result.HealthEdCess) // 60000*4%=2400
	assert.Equal(t, d(62400), result.GrossTaxPayable) // 60000+2400
}

func TestComputeTax_NewRegime_Surcharge10_Over50L(t *testing.T) {
	result := ComputeTax(TaxComputeInput{
		Income:     IncomeBreakdown{OtherSources: d(55000000)}, // 5.5Cr
		Regime:     NewRegime,
		IsResident: true,
	})
	assert.True(t, result.TaxableIncome.GreaterThan(surcharge10Threshold))
	assert.True(t, result.Surcharge.IsPositive())
}

func TestComputeTax_NewRegime_Surcharge15_Over1Cr(t *testing.T) {
	result := ComputeTax(TaxComputeInput{
		Income:     IncomeBreakdown{OtherSources: d(15000000)}, // 1.5Cr
		Regime:     NewRegime,
		IsResident: true,
	})
	assert.True(t, result.TaxableIncome.GreaterThan(surcharge15Threshold))
	assert.Equal(t, surcharge15Rate, result.SurchargeRate)
}

func TestComputeTax_NewRegime_Surcharge25_Over2Cr(t *testing.T) {
	result := ComputeTax(TaxComputeInput{
		Income:     IncomeBreakdown{OtherSources: d(25000000)}, // 2.5Cr
		Regime:     NewRegime,
		IsResident: true,
	})
	assert.True(t, result.TaxableIncome.GreaterThan(surcharge25Threshold))
	assert.Equal(t, surcharge25Rate, result.SurchargeRate)
}

func TestComputeTax_OldRegime_WithDeductions(t *testing.T) {
	result := ComputeTax(TaxComputeInput{
		Income: IncomeBreakdown{Salary: d(1500000)},
		Deductions: DeductionBreakdown{
			Section80C: d(150000),
			Section80D: d(25000),
			Section24b: d(200000),
		},
		Regime:     OldRegime,
		IsResident: true,
	})
	// 1500K - 50K std - 150K 80C - 25K 80D - 200K 24b = 1075K taxable
	assert.Equal(t, d(1075000), result.TaxableIncome)
	assert.Equal(t, d(50000), result.StandardDeduction)
	assert.Equal(t, d(425000), result.TotalDeductions)
}

func TestComputeTax_OldRegime_DeductionsIgnoredInNewRegime(t *testing.T) {
	result := ComputeTax(TaxComputeInput{
		Income: IncomeBreakdown{Salary: d(1500000)},
		Deductions: DeductionBreakdown{
			Section80C: d(150000),
			Section80D: d(25000),
		},
		Regime:     NewRegime,
		IsResident: true,
	})
	// New regime ignores 80C/80D. Only standard deduction (75K) applies.
	assert.Equal(t, d(1425000), result.TaxableIncome)
	assert.Equal(t, d(75000), result.TotalDeductions)
}

func TestComputeTax_TDS_Refund(t *testing.T) {
	result := ComputeTax(TaxComputeInput{
		Income:     IncomeBreakdown{Salary: d(1275000)},
		Regime:     NewRegime,
		IsResident: true,
		TDSCredits: d(100000),
	})
	// Tax = 62400, TDS = 100000 → refund 37600
	assert.Equal(t, d(37600), result.RefundDue)
	assert.True(t, result.NetTaxPayable.IsZero())
}

func TestComputeTax_TDS_BalancePayable(t *testing.T) {
	result := ComputeTax(TaxComputeInput{
		Income:     IncomeBreakdown{Salary: d(1275000)},
		Regime:     NewRegime,
		IsResident: true,
		TDSCredits: d(30000),
	})
	// Tax = 62400, TDS = 30000 → payable 32400
	assert.Equal(t, d(32400), result.NetTaxPayable)
	assert.True(t, result.RefundDue.IsZero())
}

func TestComputeTax_NoSalary_NoStdDeduction(t *testing.T) {
	result := ComputeTax(TaxComputeInput{
		Income:     IncomeBreakdown{OtherSources: d(800000)},
		Regime:     NewRegime,
		IsResident: true,
	})
	assert.True(t, result.StandardDeduction.IsZero())
	assert.Equal(t, d(800000), result.TaxableIncome)
}

func TestComputeTax_SlabBreakdown_Populated(t *testing.T) {
	result := ComputeTax(TaxComputeInput{
		Income:     IncomeBreakdown{Salary: d(2575000)},
		Regime:     NewRegime,
		IsResident: true,
	})
	require.NotEmpty(t, result.SlabBreakdown)
	assert.True(t, len(result.SlabBreakdown) >= 6)
}

func TestComputeTax_NonResident_NoRebate(t *testing.T) {
	result := ComputeTax(TaxComputeInput{
		Income:     IncomeBreakdown{Salary: d(775000)},
		Regime:     NewRegime,
		IsResident: false,
	})
	assert.Equal(t, d(700000), result.TaxableIncome)
	assert.True(t, result.Rebate87A.IsZero(), "NRIs don't get 87A rebate")
	assert.True(t, result.GrossTaxPayable.IsPositive())
}
