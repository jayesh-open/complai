package domain

import (
	"testing"

	"github.com/shopspring/decimal"
)

func BenchmarkComputeTax1000(b *testing.B) {
	inputs := make([]TaxComputeInput, 1000)
	for i := range inputs {
		salary := decimal.NewFromInt(int64(500000 + i*1000))
		inputs[i] = TaxComputeInput{
			Income: IncomeBreakdown{
				Salary:       salary,
				OtherSources: decimal.NewFromInt(50000),
			},
			Deductions: DeductionBreakdown{
				Section80C: decimal.NewFromInt(150000),
				Section80D: decimal.NewFromInt(25000),
			},
			Regime:     NewRegime,
			IsResident: true,
			TDSCredits: decimal.NewFromInt(80000),
		}
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, in := range inputs {
			ComputeTax(in)
		}
	}
}
