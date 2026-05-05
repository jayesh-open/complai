package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func makeBenchMonthlyData(period string, base float64) MonthlyData {
	return MonthlyData{
		ReturnPeriod: period,
		Outward: TaxBreakdown{
			TaxableValue: decimal.NewFromFloat(base),
			CGST:         decimal.NewFromFloat(base * 0.09),
			SGST:         decimal.NewFromFloat(base * 0.09),
			IGST:         decimal.NewFromFloat(base * 0.05),
			Cess:         decimal.NewFromFloat(base * 0.01),
		},
		Inward: TaxBreakdown{
			TaxableValue: decimal.NewFromFloat(base * 0.6),
			CGST:         decimal.NewFromFloat(base * 0.054),
			SGST:         decimal.NewFromFloat(base * 0.054),
			IGST:         decimal.NewFromFloat(base * 0.03),
			Cess:         decimal.NewFromFloat(base * 0.006),
		},
		ITC: ITCBreakdown{
			CGST: decimal.NewFromFloat(base * 0.054),
			SGST: decimal.NewFromFloat(base * 0.054),
			IGST: decimal.NewFromFloat(base * 0.03),
			Cess: decimal.NewFromFloat(base * 0.006),
		},
		TaxPaid: TaxBreakdown{
			CGST: decimal.NewFromFloat(base * 0.036),
			SGST: decimal.NewFromFloat(base * 0.036),
			IGST: decimal.NewFromFloat(base * 0.02),
			Cess: decimal.NewFromFloat(base * 0.004),
		},
	}
}

func BenchmarkAggregate50GSTINs(b *testing.B) {
	periods := []string{"202504", "202505", "202506", "202507", "202508", "202509", "202510", "202511", "202512", "202601", "202602", "202603"}
	gstinData := make([][]MonthlyData, 50)
	for g := range gstinData {
		months := make([]MonthlyData, 12)
		for i, p := range periods {
			base := float64(1000000 + g*100000 + i*50000)
			months[i] = makeBenchMonthlyData(p, base)
		}
		gstinData[g] = months
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for g := 0; g < 50; g++ {
			Aggregate(uuid.New(), uuid.New(), gstinData[g])
		}
	}
}

func BenchmarkReconcile100Mismatches(b *testing.B) {
	tenantID := uuid.New()
	filingID := uuid.New()
	gstr9cID := uuid.New()

	periods := []string{"202504", "202505", "202506", "202507", "202508", "202509", "202510", "202511", "202512", "202601", "202602", "202603"}
	months := make([]MonthlyData, 12)
	for i, p := range periods {
		months[i] = makeBenchMonthlyData(p, float64(5000000+i*500000))
	}

	tables := Aggregate(filingID, tenantID, months)

	filing := &GSTR9Filing{
		ID:                filingID,
		TenantID:          tenantID,
		AggregateTurnover: decimal.NewFromFloat(85000000),
	}

	audited := AuditedFinancials{
		Turnover:           decimal.NewFromFloat(87500000),
		TaxPayable:         TaxBreakdown{CGST: decimal.NewFromFloat(3000000), SGST: decimal.NewFromFloat(3000000), IGST: decimal.NewFromFloat(2000000), Cess: decimal.NewFromFloat(500000)},
		ITCClaimed:         ITCBreakdown{CGST: decimal.NewFromFloat(2500000), SGST: decimal.NewFromFloat(2500000), IGST: decimal.NewFromFloat(1500000), Cess: decimal.NewFromFloat(300000)},
		UnbilledRevenue:    decimal.NewFromFloat(1500000),
		UnadjustedAdvances: decimal.NewFromFloat(800000),
		DeemedSupply:       decimal.NewFromFloat(200000),
		CreditNotesAfterFY: decimal.NewFromFloat(500000),
		Sec15_3Adjustments: decimal.NewFromFloat(300000),
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		mismatches := Reconcile(filing, tables, audited, tenantID, gstr9cID)
		CanSubmit(mismatches)
	}
}
