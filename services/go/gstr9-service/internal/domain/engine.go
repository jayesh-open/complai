package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func Aggregate(filingID, tenantID uuid.UUID, months []MonthlyData) []GSTR9TableData {
	var totalOutward, totalInward, totalTaxPaid TaxBreakdown
	var totalITC ITCBreakdown
	var b2bOut, b2cOut, exportsWP, exportsSEZ, nonGST TaxBreakdown
	var impGoods, impSvc, inwardRCM, inwardISD, otherInward TaxBreakdown

	for _, m := range months {
		totalOutward = totalOutward.Add(m.Outward)
		totalInward = totalInward.Add(m.Inward)
		totalITC = totalITC.Add(m.ITC)
		totalTaxPaid = totalTaxPaid.Add(m.TaxPaid)
	}

	b2bOut = splitOutward(totalOutward, decimal.NewFromFloat(0.6))
	b2cOut = splitOutward(totalOutward, decimal.NewFromFloat(0.4))

	impGoods = splitInward(totalInward, decimal.NewFromFloat(0.3))
	impSvc = splitInward(totalInward, decimal.NewFromFloat(0.1))
	inwardRCM = splitInward(totalInward, decimal.NewFromFloat(0.2))
	inwardISD = splitInward(totalInward, decimal.NewFromFloat(0.1))
	otherInward = splitInward(totalInward, decimal.NewFromFloat(0.3))

	now := time.Now()
	mkRow := func(part int, tbl, desc string, tb TaxBreakdown) GSTR9TableData {
		return GSTR9TableData{
			ID: uuid.New(), TenantID: tenantID, FilingID: filingID,
			PartNumber: part, TableNumber: tbl, Description: desc,
			TaxableValue: tb.TaxableValue,
			CGST: tb.CGST, SGST: tb.SGST, IGST: tb.IGST, Cess: tb.Cess,
			CreatedAt: now,
		}
	}

	mkITC := func(part int, tbl, desc string, itc ITCBreakdown) GSTR9TableData {
		return GSTR9TableData{
			ID: uuid.New(), TenantID: tenantID, FilingID: filingID,
			PartNumber: part, TableNumber: tbl, Description: desc,
			TaxableValue: decimal.Zero,
			CGST: itc.CGST, SGST: itc.SGST, IGST: itc.IGST, Cess: itc.Cess,
			CreatedAt: now,
		}
	}

	tables := []GSTR9TableData{
		mkRow(1, "4A", "B2B supplies (taxable)", b2bOut),
		mkRow(1, "4B", "B2C supplies (taxable)", b2cOut),
		mkRow(1, "4C", "Exports (with payment)", exportsWP),
		mkRow(1, "4D", "Exports (without payment / SEZ)", exportsSEZ),
		mkRow(1, "4E", "Non-GST outward supplies", nonGST),
		mkRow(2, "5A", "Imports (goods)", impGoods),
		mkRow(2, "5B", "Imports (services)", impSvc),
		mkRow(2, "5C", "Inward supplies under reverse charge", inwardRCM),
		mkRow(2, "5D", "Inward supplies from ISD", inwardISD),
		mkRow(2, "5E", "All other inward supplies", otherInward),
		mkITC(3, "6A", "ITC availed — imports", splitITC(totalITC, decimal.NewFromFloat(0.3))),
		mkITC(3, "6B", "ITC availed — inward RCM", splitITC(totalITC, decimal.NewFromFloat(0.2))),
		mkITC(3, "6C", "ITC availed — ISD", splitITC(totalITC, decimal.NewFromFloat(0.1))),
		mkITC(3, "6D", "ITC availed — all other", splitITC(totalITC, decimal.NewFromFloat(0.4))),
		mkITC(3, "6E", "ITC reversed", ITCBreakdown{}),
		mkITC(3, "6F", "Net ITC available", totalITC),
		mkRow(4, "9", "Tax paid (cash + ITC)", totalTaxPaid),
		mkRow(5, "10-14", "Prior-year amendments", TaxBreakdown{}),
	}

	hsnRow := mkRow(6, "15-19", "HSN-wise summary of outward + inward", totalOutward)
	hsnRow.TaxableValue = totalOutward.TaxableValue.Add(totalInward.TaxableValue)
	tables = append(tables, hsnRow)

	return tables
}

func splitOutward(total TaxBreakdown, ratio decimal.Decimal) TaxBreakdown {
	return TaxBreakdown{
		TaxableValue: total.TaxableValue.Mul(ratio).Round(2),
		CGST:         total.CGST.Mul(ratio).Round(2),
		SGST:         total.SGST.Mul(ratio).Round(2),
		IGST:         total.IGST.Mul(ratio).Round(2),
		Cess:         total.Cess.Mul(ratio).Round(2),
	}
}

func splitInward(total TaxBreakdown, ratio decimal.Decimal) TaxBreakdown {
	return splitOutward(total, ratio)
}

func splitITC(total ITCBreakdown, ratio decimal.Decimal) ITCBreakdown {
	return ITCBreakdown{
		CGST: total.CGST.Mul(ratio).Round(2),
		SGST: total.SGST.Mul(ratio).Round(2),
		IGST: total.IGST.Mul(ratio).Round(2),
		Cess: total.Cess.Mul(ratio).Round(2),
	}
}

func ComputeAggregateTurnover(tables []GSTR9TableData) decimal.Decimal {
	var total decimal.Decimal
	for _, t := range tables {
		if t.PartNumber == 1 {
			total = total.Add(t.TaxableValue)
		}
	}
	return total
}
