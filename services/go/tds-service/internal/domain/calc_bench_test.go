package domain

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func BenchmarkCalculate1000(b *testing.B) {
	codes := []PaymentCode{CodeSalaryPrivate, CodeContractorOther, CodeProfessional, CodeNonResident, CodeRentPlant}
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			Calculate(CalcInput{
				PaymentCode: codes[j%len(codes)],
				GrossAmount: decimal.NewFromInt(int64(50000 + j*100)),
				HasValidPAN: true,
			})
		}
	}
}

func TestCalculation1000Under30s(t *testing.T) {
	codes := []PaymentCode{CodeSalaryPrivate, CodeContractorOther, CodeProfessional, CodeNonResident, CodeRentPlant}
	start := time.Now()
	for j := 0; j < 1000; j++ {
		Calculate(CalcInput{
			PaymentCode: codes[j%len(codes)],
			GrossAmount: decimal.NewFromInt(int64(50000 + j*100)),
			HasValidPAN: true,
		})
	}
	elapsed := time.Since(start)
	t.Logf("1000 calculations completed in %v", elapsed)
	assert.Less(t, elapsed, 30*time.Second, "1000 calculations must complete in under 30 seconds")
}
