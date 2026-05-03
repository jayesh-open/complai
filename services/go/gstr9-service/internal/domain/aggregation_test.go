package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeMonthlyData(period string, outward, inward float64, itc float64, taxPaid float64) MonthlyData {
	return MonthlyData{
		ReturnPeriod: period,
		Outward: TaxBreakdown{
			TaxableValue: decimal.NewFromFloat(outward),
			CGST:         decimal.NewFromFloat(outward * 0.09),
			SGST:         decimal.NewFromFloat(outward * 0.09),
			IGST:         decimal.NewFromFloat(outward * 0.05),
			Cess:         decimal.NewFromFloat(outward * 0.01),
		},
		Inward: TaxBreakdown{
			TaxableValue: decimal.NewFromFloat(inward),
			CGST:         decimal.NewFromFloat(inward * 0.09),
			SGST:         decimal.NewFromFloat(inward * 0.09),
			IGST:         decimal.NewFromFloat(inward * 0.04),
			Cess:         decimal.NewFromFloat(inward * 0.005),
		},
		ITC: ITCBreakdown{
			CGST: decimal.NewFromFloat(itc * 0.35),
			SGST: decimal.NewFromFloat(itc * 0.35),
			IGST: decimal.NewFromFloat(itc * 0.2),
			Cess: decimal.NewFromFloat(itc * 0.1),
		},
		TaxPaid: TaxBreakdown{
			TaxableValue: decimal.Zero,
			CGST:         decimal.NewFromFloat(taxPaid * 0.4),
			SGST:         decimal.NewFromFloat(taxPaid * 0.4),
			IGST:         decimal.NewFromFloat(taxPaid * 0.15),
			Cess:         decimal.NewFromFloat(taxPaid * 0.05),
		},
	}
}

func make12Months() []MonthlyData {
	periods := ReturnPeriodsForFY("2025-26")
	months := make([]MonthlyData, 12)
	for i, p := range periods {
		outward := 500000.0 + float64(i)*10000
		inward := 300000.0 + float64(i)*5000
		itc := 50000.0 + float64(i)*2000
		taxPaid := 30000.0 + float64(i)*1000
		months[i] = makeMonthlyData(p, outward, inward, itc, taxPaid)
	}
	return months
}

func TestAggregate_12Months_ProducesAllTables(t *testing.T) {
	filingID := uuid.New()
	tenantID := uuid.New()
	months := make12Months()

	tables := Aggregate(filingID, tenantID, months)
	require.Equal(t, 27, len(tables), "expected 27 table rows")

	tableMap := make(map[string]GSTR9TableData)
	for _, td := range tables {
		tableMap[td.TableNumber] = td
	}

	expectedTables := []string{
		"4A", "4B", "4C", "4D", "4E",
		"5A", "5B", "5C", "5D", "5E",
		"6A", "6B", "6C", "6D", "6E", "6F", "6H", "8C",
		"9",
		"10", "11", "12", "13", "14",
		"17", "18", "19",
	}
	for _, tbl := range expectedTables {
		_, ok := tableMap[tbl]
		assert.True(t, ok, "missing table %s", tbl)
	}
}

func TestAggregate_Part1_OutwardSupplies(t *testing.T) {
	filingID := uuid.New()
	tenantID := uuid.New()
	months := make12Months()
	tables := Aggregate(filingID, tenantID, months)

	var totalOutwardTaxable decimal.Decimal
	for _, m := range months {
		totalOutwardTaxable = totalOutwardTaxable.Add(m.Outward.TaxableValue)
	}

	var part1Taxable decimal.Decimal
	for _, td := range tables {
		if td.PartNumber == 1 {
			part1Taxable = part1Taxable.Add(td.TaxableValue)
			assert.Equal(t, filingID, td.FilingID)
			assert.Equal(t, tenantID, td.TenantID)
		}
	}
	assert.True(t, part1Taxable.GreaterThan(decimal.Zero), "Part 1 must have positive taxable value")
	assert.True(t, totalOutwardTaxable.Equal(part1Taxable),
		"Part 1 taxable total should equal total outward taxable: got %s vs %s", part1Taxable, totalOutwardTaxable)
}

func TestAggregate_Part2_InwardSupplies(t *testing.T) {
	filingID := uuid.New()
	tenantID := uuid.New()
	months := make12Months()
	tables := Aggregate(filingID, tenantID, months)

	var totalInwardTaxable decimal.Decimal
	for _, m := range months {
		totalInwardTaxable = totalInwardTaxable.Add(m.Inward.TaxableValue)
	}

	var part2Taxable decimal.Decimal
	for _, td := range tables {
		if td.PartNumber == 2 {
			part2Taxable = part2Taxable.Add(td.TaxableValue)
		}
	}
	assert.True(t, part2Taxable.GreaterThan(decimal.Zero))
	assert.True(t, totalInwardTaxable.Equal(part2Taxable),
		"Part 2 taxable total should equal total inward taxable: got %s vs %s", part2Taxable, totalInwardTaxable)
}

func TestAggregate_Part3_ITCDetails(t *testing.T) {
	filingID := uuid.New()
	tenantID := uuid.New()
	months := make12Months()
	tables := Aggregate(filingID, tenantID, months)

	tableMap := make(map[string]GSTR9TableData)
	for _, td := range tables {
		tableMap[td.TableNumber] = td
	}

	t6a := tableMap["6A"]
	assert.True(t, t6a.CGST.GreaterThan(decimal.Zero), "6A CGST should be positive")

	t6f := tableMap["6F"]
	assert.True(t, t6f.CGST.GreaterThan(decimal.Zero), "6F net ITC CGST should be positive")

	t6e := tableMap["6E"]
	assert.True(t, t6e.CGST.IsZero(), "6E ITC reversed should be zero (no reversals)")

	itcSum := tableMap["6A"].CGST.Add(tableMap["6B"].CGST).Add(tableMap["6C"].CGST).Add(tableMap["6D"].CGST)
	var totalITCCGST decimal.Decimal
	for _, m := range months {
		totalITCCGST = totalITCCGST.Add(m.ITC.CGST)
	}
	assert.True(t, totalITCCGST.Equal(itcSum),
		"sum of 6A-6D CGST should equal total ITC CGST: got %s vs %s", itcSum, totalITCCGST)
}

func TestAggregate_Part4_TaxPaid(t *testing.T) {
	filingID := uuid.New()
	tenantID := uuid.New()
	months := make12Months()
	tables := Aggregate(filingID, tenantID, months)

	tableMap := make(map[string]GSTR9TableData)
	for _, td := range tables {
		tableMap[td.TableNumber] = td
	}

	t9 := tableMap["9"]
	var totalTaxCGST decimal.Decimal
	for _, m := range months {
		totalTaxCGST = totalTaxCGST.Add(m.TaxPaid.CGST)
	}
	assert.True(t, t9.CGST.Equal(totalTaxCGST))
}

func TestAggregate_Part5_PriorYearAmendments(t *testing.T) {
	filingID := uuid.New()
	tenantID := uuid.New()
	months := make12Months()
	tables := Aggregate(filingID, tenantID, months)

	for _, td := range tables {
		if td.PartNumber == 5 {
			assert.True(t, td.TaxableValue.IsZero(), "Part 5 table %s should be zero (no amendments)", td.TableNumber)
		}
	}
}

func TestAggregate_Part6_HSNSummary(t *testing.T) {
	filingID := uuid.New()
	tenantID := uuid.New()
	months := make12Months()
	tables := Aggregate(filingID, tenantID, months)

	tableMap := make(map[string]GSTR9TableData)
	for _, td := range tables {
		tableMap[td.TableNumber] = td
	}

	t17 := tableMap["17"]
	assert.True(t, t17.TaxableValue.GreaterThan(decimal.Zero), "Table 17 outward HSN should have value")

	t18 := tableMap["18"]
	assert.True(t, t18.TaxableValue.GreaterThan(decimal.Zero), "Table 18 inward HSN should have value")

	t19 := tableMap["19"]
	assert.True(t, t19.CGST.IsZero(), "Table 19 late fees should be zero")
}

func TestAggregate_LateITCReclaim_Table8C(t *testing.T) {
	filingID := uuid.New()
	tenantID := uuid.New()
	months := make12Months()

	months[0].LateITCReclaim = ITCBreakdown{
		CGST: decimal.NewFromInt(10000),
		SGST: decimal.NewFromInt(10000),
		IGST: decimal.NewFromInt(5000),
		Cess: decimal.NewFromInt(1000),
	}

	tables := Aggregate(filingID, tenantID, months)
	tableMap := make(map[string]GSTR9TableData)
	for _, td := range tables {
		tableMap[td.TableNumber] = td
	}

	t8c := tableMap["8C"]
	assert.True(t, t8c.CGST.Equal(decimal.NewFromInt(10000)), "8C CGST should be 10000")
	assert.True(t, t8c.SGST.Equal(decimal.NewFromInt(10000)), "8C SGST should be 10000")
	assert.True(t, t8c.IGST.Equal(decimal.NewFromInt(5000)), "8C IGST should be 5000")
	assert.True(t, t8c.Cess.Equal(decimal.NewFromInt(1000)), "8C Cess should be 1000")
}

func TestAggregate_Rule37Reclaim_Table6H(t *testing.T) {
	filingID := uuid.New()
	tenantID := uuid.New()
	months := make12Months()

	months[3].Rule37Reclaim = ITCBreakdown{
		CGST: decimal.NewFromInt(8000),
		SGST: decimal.NewFromInt(8000),
		IGST: decimal.NewFromInt(4000),
		Cess: decimal.NewFromInt(500),
	}

	tables := Aggregate(filingID, tenantID, months)
	tableMap := make(map[string]GSTR9TableData)
	for _, td := range tables {
		tableMap[td.TableNumber] = td
	}

	t6h := tableMap["6H"]
	assert.True(t, t6h.CGST.Equal(decimal.NewFromInt(8000)), "6H CGST should be 8000")
	assert.True(t, t6h.SGST.Equal(decimal.NewFromInt(8000)), "6H SGST should be 8000")
	assert.True(t, t6h.IGST.Equal(decimal.NewFromInt(4000)), "6H IGST should be 4000")
}

func TestAggregate_LateITC_NoDoubleCounting(t *testing.T) {
	filingID := uuid.New()
	tenantID := uuid.New()
	months := make12Months()

	lateITC := ITCBreakdown{
		CGST: decimal.NewFromInt(15000),
		SGST: decimal.NewFromInt(15000),
		IGST: decimal.NewFromInt(7000),
		Cess: decimal.NewFromInt(2000),
	}
	rule37 := ITCBreakdown{
		CGST: decimal.NewFromInt(5000),
		SGST: decimal.NewFromInt(5000),
		IGST: decimal.NewFromInt(3000),
		Cess: decimal.NewFromInt(500),
	}
	months[0].LateITCReclaim = lateITC
	months[6].Rule37Reclaim = rule37

	tables := Aggregate(filingID, tenantID, months)
	tableMap := make(map[string]GSTR9TableData)
	for _, td := range tables {
		tableMap[td.TableNumber] = td
	}

	t8c := tableMap["8C"]
	assert.True(t, t8c.CGST.Equal(lateITC.CGST), "8C should only contain late ITC, not Rule 37")

	t6h := tableMap["6H"]
	assert.True(t, t6h.CGST.Equal(rule37.CGST), "6H should only contain Rule 37 reclaim, not late ITC")

	t6f := tableMap["6F"]
	var totalRegularCGST decimal.Decimal
	for _, m := range months {
		totalRegularCGST = totalRegularCGST.Add(m.ITC.CGST)
	}
	expectedNetCGST := totalRegularCGST.Add(rule37.CGST)
	assert.True(t, t6f.CGST.Equal(expectedNetCGST),
		"6F net ITC should include regular ITC + Rule 37 but NOT late ITC (8C): got %s, expected %s",
		t6f.CGST, expectedNetCGST)
}

func TestAggregate_EmptyMonths(t *testing.T) {
	filingID := uuid.New()
	tenantID := uuid.New()
	tables := Aggregate(filingID, tenantID, []MonthlyData{})
	assert.Equal(t, 27, len(tables), "should still produce all table structure rows")

	for _, td := range tables {
		assert.True(t, td.TaxableValue.IsZero() || td.CGST.IsZero(),
			"table %s should be zero with no input", td.TableNumber)
	}
}

func TestAggregate_SingleMonth(t *testing.T) {
	filingID := uuid.New()
	tenantID := uuid.New()
	months := []MonthlyData{makeMonthlyData("202504", 1000000, 600000, 80000, 50000)}

	tables := Aggregate(filingID, tenantID, months)
	assert.Equal(t, 27, len(tables))

	turnover := ComputeAggregateTurnover(tables)
	assert.True(t, turnover.Equal(decimal.NewFromFloat(1000000)),
		"turnover with single month of 10L outward should be 10L: got %s", turnover)
}

func TestComputeAggregateTurnover(t *testing.T) {
	tables := []GSTR9TableData{
		{PartNumber: 1, TableNumber: "4A", TaxableValue: decimal.NewFromInt(600000)},
		{PartNumber: 1, TableNumber: "4B", TaxableValue: decimal.NewFromInt(400000)},
		{PartNumber: 2, TableNumber: "5A", TaxableValue: decimal.NewFromInt(200000)},
		{PartNumber: 4, TableNumber: "9", TaxableValue: decimal.Zero},
	}
	turnover := ComputeAggregateTurnover(tables)
	assert.True(t, turnover.Equal(decimal.NewFromInt(1000000)),
		"turnover should only sum Part 1: got %s", turnover)
}

func TestTaxBreakdown_Add(t *testing.T) {
	a := TaxBreakdown{
		TaxableValue: decimal.NewFromInt(100),
		CGST:         decimal.NewFromInt(9),
		SGST:         decimal.NewFromInt(9),
		IGST:         decimal.NewFromInt(18),
		Cess:         decimal.NewFromInt(1),
	}
	b := TaxBreakdown{
		TaxableValue: decimal.NewFromInt(200),
		CGST:         decimal.NewFromInt(18),
		SGST:         decimal.NewFromInt(18),
		IGST:         decimal.NewFromInt(36),
		Cess:         decimal.NewFromInt(2),
	}
	sum := a.Add(b)
	assert.True(t, sum.TaxableValue.Equal(decimal.NewFromInt(300)))
	assert.True(t, sum.CGST.Equal(decimal.NewFromInt(27)))
	assert.True(t, sum.TotalTax().Equal(decimal.NewFromInt(111)))
}

func TestITCBreakdown_Add(t *testing.T) {
	a := ITCBreakdown{CGST: decimal.NewFromInt(10), SGST: decimal.NewFromInt(10), IGST: decimal.NewFromInt(5), Cess: decimal.NewFromInt(1)}
	b := ITCBreakdown{CGST: decimal.NewFromInt(20), SGST: decimal.NewFromInt(20), IGST: decimal.NewFromInt(10), Cess: decimal.NewFromInt(2)}
	sum := a.Add(b)
	assert.True(t, sum.Total().Equal(decimal.NewFromInt(78)))
}

func TestReturnPeriodsForFY(t *testing.T) {
	periods := ReturnPeriodsForFY("2025-26")
	require.Equal(t, 12, len(periods))
	assert.Equal(t, "202504", periods[0])
	assert.Equal(t, "202603", periods[11])
}

func TestReturnPeriodsForFY_InvalidShort(t *testing.T) {
	periods := ReturnPeriodsForFY("25-26")
	assert.Empty(t, periods)
}

func TestReturnPeriodsForFY_Another(t *testing.T) {
	periods := ReturnPeriodsForFY("2024-25")
	require.Equal(t, 12, len(periods))
	assert.Equal(t, "202404", periods[0])
	assert.Equal(t, "202503", periods[11])
}
