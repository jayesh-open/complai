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
	r := Calculate(CalcInput{PaymentCode: CodeSalaryPrivate, AnnualSalary: d(50000)})
	assert.False(t, r.ThresholdMet)
	assert.True(t, r.TDSAmount.IsZero())
}

func TestSalary_AtStdDeduction(t *testing.T) {
	r := Calculate(CalcInput{PaymentCode: CodeSalaryPrivate, AnnualSalary: d(75000)})
	assert.False(t, r.ThresholdMet)
}

func TestSalary_5L_WithinRebate(t *testing.T) {
	r := Calculate(CalcInput{PaymentCode: CodeSalaryPrivate, AnnualSalary: d(500000)})
	assert.True(t, r.ThresholdMet)
	assert.True(t, r.TDSAmount.IsZero(), "₹5L salary should be zero tax after rebate")
}

func TestSalary_10L_WithinRebate(t *testing.T) {
	r := Calculate(CalcInput{PaymentCode: CodeSalaryPrivate, AnnualSalary: d(1000000)})
	assert.True(t, r.ThresholdMet)
	assert.True(t, r.TDSAmount.IsZero(), "₹10L salary should be zero tax after rebate")
}

func TestSalary_12L_WithinRebate(t *testing.T) {
	r := Calculate(CalcInput{PaymentCode: CodeSalaryPrivate, AnnualSalary: d(1200000)})
	assert.True(t, r.TDSAmount.IsZero(), "₹12L → taxable ₹11.25L → within rebate limit")
}

func TestSalary_1275000_ExactRebateEdge(t *testing.T) {
	r := Calculate(CalcInput{PaymentCode: CodeSalaryPrivate, AnnualSalary: d(1275000)})
	assert.True(t, r.TDSAmount.IsZero(), "taxable income exactly at ₹12L → full rebate")
}

func TestSalary_13L_AboveRebate(t *testing.T) {
	r := Calculate(CalcInput{PaymentCode: CodeSalaryPrivate, AnnualSalary: d(1300000)})
	assert.True(t, r.ThresholdMet)
	assert.Equal(t, "5525", r.TDSAmount.String())
}

func TestSalary_20L(t *testing.T) {
	r := Calculate(CalcInput{PaymentCode: CodeSalaryPrivate, AnnualSalary: d(2000000)})
	assert.Equal(t, "16033", r.TDSAmount.String())
}

func TestSalary_30L_TopSlab(t *testing.T) {
	r := Calculate(CalcInput{PaymentCode: CodeSalaryPrivate, AnnualSalary: d(3000000)})
	assert.Equal(t, "39650", r.TDSAmount.String())
}

func TestSalary_50L_HighSalary(t *testing.T) {
	r := Calculate(CalcInput{PaymentCode: CodeSalaryPrivate, AnnualSalary: d(5000000)})
	assert.Equal(t, "91650", r.TDSAmount.String())
}

func TestSalary_StateEmployer(t *testing.T) {
	r := Calculate(CalcInput{PaymentCode: CodeSalaryState, AnnualSalary: d(2000000)})
	assert.Equal(t, Section392, r.Section)
	assert.Equal(t, "16033", r.TDSAmount.String())
}

func TestSalary_CentralEmployer(t *testing.T) {
	r := Calculate(CalcInput{PaymentCode: CodeSalaryCentral, AnnualSalary: d(2000000)})
	assert.Equal(t, Section392, r.Section)
	assert.Equal(t, "16033", r.TDSAmount.String())
}

func TestEPF_BelowThreshold(t *testing.T) {
	r := Calculate(CalcInput{PaymentCode: CodeEPFWithdrawal, GrossAmount: d(40000), HasValidPAN: true})
	assert.False(t, r.ThresholdMet)
	assert.True(t, r.TDSAmount.IsZero())
}

func TestEPF_AboveThreshold(t *testing.T) {
	r := Calculate(CalcInput{PaymentCode: CodeEPFWithdrawal, GrossAmount: d(100000), HasValidPAN: true})
	assert.True(t, r.ThresholdMet)
	assert.Equal(t, "0.1", r.Rate.String())
	assert.Equal(t, "10000", r.TDSAmount.String())
}

func TestContractorIndiv_BelowBothThresholds(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodeContractorIndiv, GrossAmount: d(25000),
		HasValidPAN: true, AggregateForFY: d(50000),
	})
	assert.False(t, r.ThresholdMet, "25K single + 75K aggregate = below both thresholds")
}

func TestContractorIndiv_SingleThresholdExceeds(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodeContractorIndiv, GrossAmount: d(35000),
		HasValidPAN: true, AggregateForFY: d(0),
	})
	assert.True(t, r.ThresholdMet)
	assert.Equal(t, "0.01", r.Rate.String())
	assert.Equal(t, "350", r.TDSAmount.String())
}

func TestContractorIndiv_AggregateThresholdExceeds(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodeContractorIndiv, GrossAmount: d(20000),
		HasValidPAN: true, AggregateForFY: d(90000),
	})
	assert.True(t, r.ThresholdMet, "20K single but 110K aggregate")
	assert.Equal(t, "200", r.TDSAmount.String())
}

func TestContractorOther_Rate(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodeContractorOther, GrossAmount: d(100000),
		HasValidPAN: true, AggregateForFY: d(0),
	})
	assert.Equal(t, "0.02", r.Rate.String())
	assert.Equal(t, "2000", r.TDSAmount.String())
}

func TestContractor_NoPAN(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodeContractorIndiv, GrossAmount: d(50000),
		HasValidPAN: false,
	})
	assert.True(t, r.NoPAN)
	assert.Equal(t, "0.2", r.Rate.String())
	assert.Equal(t, "10000", r.TDSAmount.String())
}

func TestContractor_LowerCert(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodeContractorOther, GrossAmount: d(100000),
		HasValidPAN: true, LowerCertRate: dp(0.005),
	})
	assert.True(t, r.LowerCert)
	assert.Equal(t, "0.005", r.Rate.String())
	assert.Equal(t, "500", r.TDSAmount.String())
}

func TestRentLand_BelowThreshold(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodeRentLand, GrossAmount: d(40000),
		HasValidPAN: true,
	})
	assert.False(t, r.ThresholdMet, "below per-payment ₹50K threshold")
}

func TestRentLand_AboveThreshold(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodeRentLand, GrossAmount: d(60000),
		HasValidPAN: true,
	})
	assert.True(t, r.ThresholdMet)
	assert.Equal(t, "0.1", r.Rate.String())
	assert.Equal(t, "6000", r.TDSAmount.String())
}

func TestRentPlant(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodeRentPlant, GrossAmount: d(100000),
		HasValidPAN: true,
	})
	assert.True(t, r.ThresholdMet)
	assert.Equal(t, "0.02", r.Rate.String())
	assert.Equal(t, "2000", r.TDSAmount.String())
}

func TestRent_NoPAN(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodeRentLand, GrossAmount: d(60000),
		HasValidPAN: false,
	})
	assert.True(t, r.NoPAN)
	assert.Equal(t, "0.2", r.Rate.String())
}

func TestTechnical_BelowThreshold(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodeTechnical, GrossAmount: d(25000), HasValidPAN: true,
		AggregateForFY: d(0),
	})
	assert.False(t, r.ThresholdMet)
}

func TestTechnical_AboveThreshold(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodeTechnical, GrossAmount: d(60000), HasValidPAN: true,
		AggregateForFY: d(0),
	})
	assert.True(t, r.ThresholdMet)
	assert.Equal(t, "0.02", r.Rate.String())
	assert.Equal(t, "1200", r.TDSAmount.String())
}

func TestProfessional_BelowThreshold(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodeProfessional, GrossAmount: d(25000), HasValidPAN: true,
		AggregateForFY: d(0),
	})
	assert.False(t, r.ThresholdMet)
}

func TestProfessional_AboveThreshold(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodeProfessional, GrossAmount: d(60000), HasValidPAN: true,
		AggregateForFY: d(0),
	})
	assert.True(t, r.ThresholdMet)
	assert.Equal(t, "0.1", r.Rate.String())
	assert.Equal(t, "6000", r.TDSAmount.String())
}

func TestProfessional_NoPAN(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodeProfessional, GrossAmount: d(60000), HasValidPAN: false,
		AggregateForFY: d(0),
	})
	assert.True(t, r.NoPAN)
	assert.Equal(t, "0.2", r.Rate.String())
	assert.Equal(t, "12000", r.TDSAmount.String())
}

func TestDirector_NoThreshold(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodeDirectorRem, GrossAmount: d(500000), HasValidPAN: true,
	})
	assert.True(t, r.ThresholdMet)
	assert.Equal(t, "0.1", r.Rate.String())
	assert.Equal(t, "50000", r.TDSAmount.String())
}

func TestDirector_SmallAmount(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodeDirectorRem, GrossAmount: d(1000), HasValidPAN: true,
	})
	assert.True(t, r.ThresholdMet, "director remuneration has no threshold")
	assert.Equal(t, "100", r.TDSAmount.String())
}

func TestPurchase_BelowThreshold(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodePurchaseGoods, GrossAmount: d(1000000), HasValidPAN: true,
		AggregateForFY: d(3000000),
	})
	assert.False(t, r.ThresholdMet)
}

func TestPurchase_AboveThreshold(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodePurchaseGoods, GrossAmount: d(1000000), HasValidPAN: true,
		AggregateForFY: d(4500000),
	})
	assert.True(t, r.ThresholdMet)
	assert.Equal(t, "0.001", r.Rate.String())
	assert.Equal(t, "1000", r.TDSAmount.String())
}

func TestPurchase_NoPAN_SpecialRate(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodePurchaseGoods, GrossAmount: d(1000000), HasValidPAN: false,
		AggregateForFY: d(5000000),
	})
	assert.True(t, r.NoPAN)
	assert.Equal(t, "0.05", r.Rate.String())
	assert.Equal(t, "50000", r.TDSAmount.String())
}

func TestNonResident_DefaultRate(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodeNonResident, GrossAmount: d(1000000),
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
		PaymentCode: CodeNonResident, GrossAmount: d(1000000), HasValidPAN: true,
		DTAARate: dp(0.10),
	})
	assert.Equal(t, "0.1", r.Rate.String())
	assert.Equal(t, "100000", r.TDSAmount.String())
	assert.Equal(t, "4000", r.Cess.String())
	assert.Equal(t, "104000", r.TotalTax.String())
}

func TestNonResident_ZeroDTAA(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodeNonResident, GrossAmount: d(500000), HasValidPAN: true,
		DTAARate: dp(0),
	})
	assert.True(t, r.TDSAmount.IsZero())
	assert.True(t, r.TotalTax.IsZero())
}

func TestUnknownPaymentCode(t *testing.T) {
	r := Calculate(CalcInput{PaymentCode: "9999"})
	assert.Contains(t, r.Explanation, "unknown payment code")
}

func TestValidPaymentCode(t *testing.T) {
	assert.True(t, ValidPaymentCode(CodeSalaryPrivate))
	assert.True(t, ValidPaymentCode(CodeContractorIndiv))
	assert.True(t, ValidPaymentCode(CodeContractorOther))
	assert.True(t, ValidPaymentCode(CodeRentPlant))
	assert.True(t, ValidPaymentCode(CodeRentLand))
	assert.True(t, ValidPaymentCode(CodeTechnical))
	assert.True(t, ValidPaymentCode(CodeProfessional))
	assert.True(t, ValidPaymentCode(CodeDirectorRem))
	assert.True(t, ValidPaymentCode(CodePurchaseGoods))
	assert.True(t, ValidPaymentCode(CodeNonResident))
	assert.False(t, ValidPaymentCode("9999"))
	assert.False(t, ValidPaymentCode(""))
}

func TestValidSection(t *testing.T) {
	assert.True(t, ValidSection(Section392))
	assert.True(t, ValidSection(Section393_1))
	assert.True(t, ValidSection(Section393_2))
	assert.True(t, ValidSection(Section393_3))
	assert.False(t, ValidSection("999"))
	assert.False(t, ValidSection(""))
}

func TestRejectITA1961Sections(t *testing.T) {
	ita1961 := []string{"192", "194C", "194J", "194A", "194H", "194I", "195", "194Q"}
	for _, s := range ita1961 {
		assert.False(t, ValidSection(Section(s)), "ITA 1961 section %q must be rejected", s)
	}
}

func TestSectionForCode(t *testing.T) {
	assert.Equal(t, Section392, SectionForCode(CodeSalaryPrivate))
	assert.Equal(t, Section393_1, SectionForCode(CodeContractorOther))
	assert.Equal(t, Section393_2, SectionForCode(CodeNonResident))
}

func TestZeroGrossAmount(t *testing.T) {
	r := Calculate(CalcInput{
		PaymentCode: CodeContractorIndiv, GrossAmount: d(0),
		HasValidPAN: true, AggregateForFY: d(200000),
	})
	assert.True(t, r.TDSAmount.IsZero())
}
