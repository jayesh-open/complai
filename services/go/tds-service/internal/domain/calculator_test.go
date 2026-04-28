package domain

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func d(v int64) decimal.Decimal   { return decimal.NewFromInt(v) }
func df(v float64) decimal.Decimal { return decimal.NewFromFloat(v) }
func dp(v float64) *decimal.Decimal { r := decimal.NewFromFloat(v); return &r }

func TestSalary_BelowStdDeduction(t *testing.T) {
	r := Calculate(CalcInput{Section: Section192, AnnualSalary: d(50000)})
	assert.False(t, r.ThresholdMet)
	assert.True(t, r.TDSAmount.IsZero())
}

func TestSalary_AtStdDeduction(t *testing.T) {
	r := Calculate(CalcInput{Section: Section192, AnnualSalary: d(75000)})
	assert.False(t, r.ThresholdMet)
}

func TestSalary_5L_WithinRebate(t *testing.T) {
	r := Calculate(CalcInput{Section: Section192, AnnualSalary: d(500000)})
	assert.True(t, r.ThresholdMet)
	assert.True(t, r.TDSAmount.IsZero(), "₹5L salary should be zero tax after rebate")
}

func TestSalary_10L_WithinRebate(t *testing.T) {
	r := Calculate(CalcInput{Section: Section192, AnnualSalary: d(1000000)})
	assert.True(t, r.ThresholdMet)
	assert.True(t, r.TDSAmount.IsZero(), "₹10L salary should be zero tax after rebate")
}

func TestSalary_12L_WithinRebate(t *testing.T) {
	r := Calculate(CalcInput{Section: Section192, AnnualSalary: d(1200000)})
	assert.True(t, r.TDSAmount.IsZero(), "₹12L → taxable ₹11.25L → within rebate limit")
}

func TestSalary_1275000_ExactRebateEdge(t *testing.T) {
	// Taxable = 1275000 - 75000 = 1200000 → exactly at rebate limit
	r := Calculate(CalcInput{Section: Section192, AnnualSalary: d(1275000)})
	assert.True(t, r.TDSAmount.IsZero(), "taxable income exactly at ₹12L → full rebate")
}

func TestSalary_13L_AboveRebate(t *testing.T) {
	// Taxable = 1300000 - 75000 = 1225000
	// 0-4L: 0, 4-8L: 20000, 8-12L: 40000, 12-12.25L: 3750 = 63750
	// Cess = 2550, Total = 66300, Monthly = 5525
	r := Calculate(CalcInput{Section: Section192, AnnualSalary: d(1300000)})
	assert.True(t, r.ThresholdMet)
	assert.Equal(t, "5525", r.TDSAmount.String())
}

func TestSalary_20L(t *testing.T) {
	// Taxable = 1925000
	// 0-4L:0, 4-8L:20000, 8-12L:40000, 12-16L:60000, 16-19.25L:65000 = 185000
	// Cess = 7400, Total = 192400, Monthly = 16033
	r := Calculate(CalcInput{Section: Section192, AnnualSalary: d(2000000)})
	assert.Equal(t, "16033", r.TDSAmount.String())
}

func TestSalary_30L_TopSlab(t *testing.T) {
	// Taxable = 2925000
	// 0-4L:0, 4-8L:20000, 8-12L:40000, 12-16L:60000, 16-20L:80000, 20-24L:100000, 24-29.25L:157500 = 457500
	// Cess = 18300, Total = 475800, Monthly = 39650
	r := Calculate(CalcInput{Section: Section192, AnnualSalary: d(3000000)})
	assert.Equal(t, "39650", r.TDSAmount.String())
}

func TestSalary_50L_HighSalary(t *testing.T) {
	// Taxable = 4925000
	// Slabs sum: 0+20000+40000+60000+80000+100000+757500 = 1057500
	// Cess = 42300, Total = 1099800, Monthly = 91650
	r := Calculate(CalcInput{Section: Section192, AnnualSalary: d(5000000)})
	assert.Equal(t, "91650", r.TDSAmount.String())
}

func TestContractor_Individual_BelowBothThresholds(t *testing.T) {
	r := Calculate(CalcInput{
		Section: Section194C, GrossAmount: d(25000), DeducteeType: DeducteeIndividual,
		HasValidPAN: true, AggregateForFY: d(50000),
	})
	assert.False(t, r.ThresholdMet, "25K single + 75K aggregate = below both thresholds")
}

func TestContractor_Individual_SingleThresholdExceeds(t *testing.T) {
	r := Calculate(CalcInput{
		Section: Section194C, GrossAmount: d(35000), DeducteeType: DeducteeIndividual,
		HasValidPAN: true, AggregateForFY: d(0),
	})
	assert.True(t, r.ThresholdMet)
	assert.Equal(t, "0.01", r.Rate.String())
	assert.Equal(t, "350", r.TDSAmount.String())
}

func TestContractor_Individual_AggregateThresholdExceeds(t *testing.T) {
	r := Calculate(CalcInput{
		Section: Section194C, GrossAmount: d(20000), DeducteeType: DeducteeIndividual,
		HasValidPAN: true, AggregateForFY: d(90000),
	})
	assert.True(t, r.ThresholdMet, "20K single but 110K aggregate")
	assert.Equal(t, "200", r.TDSAmount.String())
}

func TestContractor_Company_Rate(t *testing.T) {
	r := Calculate(CalcInput{
		Section: Section194C, GrossAmount: d(100000), DeducteeType: DeducteeCompany,
		HasValidPAN: true, AggregateForFY: d(0),
	})
	assert.Equal(t, "0.02", r.Rate.String())
	assert.Equal(t, "2000", r.TDSAmount.String())
}

func TestContractor_HUF_IndividualRate(t *testing.T) {
	r := Calculate(CalcInput{
		Section: Section194C, GrossAmount: d(50000), DeducteeType: DeducteeHUF,
		HasValidPAN: true,
	})
	assert.Equal(t, "0.01", r.Rate.String())
}

func TestContractor_NoPAN(t *testing.T) {
	r := Calculate(CalcInput{
		Section: Section194C, GrossAmount: d(50000), DeducteeType: DeducteeIndividual,
		HasValidPAN: false,
	})
	assert.True(t, r.NoPAN)
	assert.Equal(t, "0.2", r.Rate.String())
	assert.Equal(t, "10000", r.TDSAmount.String())
}

func TestContractor_LowerCert(t *testing.T) {
	r := Calculate(CalcInput{
		Section: Section194C, GrossAmount: d(100000), DeducteeType: DeducteeCompany,
		HasValidPAN: true, LowerCertRate: dp(0.005),
	})
	assert.True(t, r.LowerCert)
	assert.Equal(t, "0.005", r.Rate.String())
	assert.Equal(t, "500", r.TDSAmount.String())
}

func TestRent_LandBuilding_BelowThreshold(t *testing.T) {
	r := Calculate(CalcInput{
		Section: Section194I, GrossAmount: d(200000), RentType: RentLandBuilding,
		HasValidPAN: true, AggregateForFY: d(0),
	})
	assert.False(t, r.ThresholdMet)
}

func TestRent_LandBuilding_AboveThreshold(t *testing.T) {
	r := Calculate(CalcInput{
		Section: Section194I, GrossAmount: d(50000), RentType: RentLandBuilding,
		HasValidPAN: true, AggregateForFY: d(200000),
	})
	assert.True(t, r.ThresholdMet)
	assert.Equal(t, "0.1", r.Rate.String())
	assert.Equal(t, "5000", r.TDSAmount.String())
}

func TestRent_PlantMachinery(t *testing.T) {
	r := Calculate(CalcInput{
		Section: Section194I, GrossAmount: d(500000), RentType: RentPlantMachinery,
		HasValidPAN: true, AggregateForFY: d(0),
	})
	assert.Equal(t, "0.02", r.Rate.String())
	assert.Equal(t, "10000", r.TDSAmount.String())
}

func TestRent_NoPAN(t *testing.T) {
	r := Calculate(CalcInput{
		Section: Section194I, GrossAmount: d(300000), RentType: RentLandBuilding,
		HasValidPAN: false,
	})
	assert.True(t, r.NoPAN)
	assert.Equal(t, "0.2", r.Rate.String())
}

func TestProfessional_BelowThreshold(t *testing.T) {
	r := Calculate(CalcInput{
		Section: Section194J, GrossAmount: d(25000), HasValidPAN: true,
	})
	assert.False(t, r.ThresholdMet)
}

func TestProfessional_AboveThreshold(t *testing.T) {
	r := Calculate(CalcInput{
		Section: Section194J, GrossAmount: d(50000), HasValidPAN: true,
	})
	assert.True(t, r.ThresholdMet)
	assert.Equal(t, "0.1", r.Rate.String())
	assert.Equal(t, "5000", r.TDSAmount.String())
}

func TestProfessional_NoPAN(t *testing.T) {
	r := Calculate(CalcInput{
		Section: Section194J, GrossAmount: d(50000), HasValidPAN: false,
	})
	assert.True(t, r.NoPAN)
	assert.Equal(t, "0.2", r.Rate.String())
	assert.Equal(t, "10000", r.TDSAmount.String())
}

func TestPurchase_BelowThreshold(t *testing.T) {
	r := Calculate(CalcInput{
		Section: Section194Q, GrossAmount: d(1000000), HasValidPAN: true,
		AggregateForFY: d(3000000),
	})
	assert.False(t, r.ThresholdMet)
}

func TestPurchase_AboveThreshold(t *testing.T) {
	r := Calculate(CalcInput{
		Section: Section194Q, GrossAmount: d(1000000), HasValidPAN: true,
		AggregateForFY: d(4500000),
	})
	assert.True(t, r.ThresholdMet)
	assert.Equal(t, "0.001", r.Rate.String())
	assert.Equal(t, "1000", r.TDSAmount.String())
}

func TestPurchase_NoPAN_SpecialRate(t *testing.T) {
	r := Calculate(CalcInput{
		Section: Section194Q, GrossAmount: d(1000000), HasValidPAN: false,
		AggregateForFY: d(5000000),
	})
	assert.True(t, r.NoPAN)
	assert.Equal(t, "0.05", r.Rate.String())
	assert.Equal(t, "50000", r.TDSAmount.String())
}

func TestNonResident_DefaultRate(t *testing.T) {
	r := Calculate(CalcInput{
		Section: Section195, GrossAmount: d(1000000), ResidentStatus: NonResident,
		HasValidPAN: true,
	})
	assert.True(t, r.ThresholdMet)
	assert.Equal(t, "0.2", r.Rate.String())
	assert.Equal(t, "200000", r.TDSAmount.String())
	assert.Equal(t, "8000", r.Cess.String())
	assert.Equal(t, "208000", r.TotalTax.String())
}

func TestNonResident_DTAARate(t *testing.T) {
	r := Calculate(CalcInput{
		Section: Section195, GrossAmount: d(1000000), HasValidPAN: true,
		DTAARate: dp(0.10),
	})
	assert.Equal(t, "0.1", r.Rate.String())
	assert.Equal(t, "100000", r.TDSAmount.String())
	assert.Equal(t, "4000", r.Cess.String())
	assert.Equal(t, "104000", r.TotalTax.String())
}

func TestNonResident_ZeroDTAA(t *testing.T) {
	r := Calculate(CalcInput{
		Section: Section195, GrossAmount: d(500000), HasValidPAN: true,
		DTAARate: dp(0),
	})
	assert.True(t, r.TDSAmount.IsZero())
	assert.True(t, r.TotalTax.IsZero())
}

func TestUnknownSection(t *testing.T) {
	r := Calculate(CalcInput{Section: "999"})
	assert.Contains(t, r.Explanation, "unknown section")
}

func TestValidSection(t *testing.T) {
	assert.True(t, ValidSection(Section192))
	assert.True(t, ValidSection(Section194C))
	assert.True(t, ValidSection(Section194I))
	assert.True(t, ValidSection(Section194J))
	assert.True(t, ValidSection(Section194Q))
	assert.True(t, ValidSection(Section195))
	assert.False(t, ValidSection("999"))
	assert.False(t, ValidSection(""))
}

func TestZeroGrossAmount(t *testing.T) {
	r := Calculate(CalcInput{
		Section: Section194C, GrossAmount: d(0), DeducteeType: DeducteeIndividual,
		HasValidPAN: true, AggregateForFY: d(200000),
	})
	assert.True(t, r.TDSAmount.IsZero())
}
