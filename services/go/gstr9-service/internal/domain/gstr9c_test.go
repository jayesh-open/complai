package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClassifySeverity_Info(t *testing.T) {
	assert.Equal(t, SeverityInfo, ClassifySeverity(decimal.NewFromInt(50)))
	assert.Equal(t, SeverityInfo, ClassifySeverity(decimal.NewFromInt(-50)))
	assert.Equal(t, SeverityInfo, ClassifySeverity(decimal.Zero))
}

func TestClassifySeverity_Warn(t *testing.T) {
	assert.Equal(t, SeverityWarn, ClassifySeverity(decimal.NewFromInt(100)))
	assert.Equal(t, SeverityWarn, ClassifySeverity(decimal.NewFromInt(5000)))
	assert.Equal(t, SeverityWarn, ClassifySeverity(decimal.NewFromInt(-9999)))
}

func TestClassifySeverity_Error(t *testing.T) {
	assert.Equal(t, SeverityError, ClassifySeverity(decimal.NewFromInt(10000)))
	assert.Equal(t, SeverityError, ClassifySeverity(decimal.NewFromInt(50000)))
	assert.Equal(t, SeverityError, ClassifySeverity(decimal.NewFromInt(-10000)))
}

func TestCanSubmit_NoMismatches(t *testing.T) {
	assert.True(t, CanSubmit(nil))
	assert.True(t, CanSubmit([]GSTR9CMismatch{}))
}

func TestCanSubmit_OnlyInfoAndWarn(t *testing.T) {
	mismatches := []GSTR9CMismatch{
		{Severity: SeverityInfo},
		{Severity: SeverityWarn},
	}
	assert.True(t, CanSubmit(mismatches))
}

func TestCanSubmit_UnresolvedErrorBlocks(t *testing.T) {
	mismatches := []GSTR9CMismatch{
		{Severity: SeverityInfo},
		{Severity: SeverityError, Resolved: false},
	}
	assert.False(t, CanSubmit(mismatches))
}

func TestCanSubmit_ResolvedErrorAllows(t *testing.T) {
	mismatches := []GSTR9CMismatch{
		{Severity: SeverityError, Resolved: true},
		{Severity: SeverityWarn},
	}
	assert.True(t, CanSubmit(mismatches))
}

func TestCanSubmit_MixedResolvedAndUnresolved(t *testing.T) {
	mismatches := []GSTR9CMismatch{
		{Severity: SeverityError, Resolved: true},
		{Severity: SeverityError, Resolved: false},
	}
	assert.False(t, CanSubmit(mismatches))
}

func buildTestTables() []GSTR9TableData {
	tenantID := uuid.New()
	filingID := uuid.New()
	months := make([]MonthlyData, 12)
	for i := range months {
		months[i] = MonthlyData{
			Outward: TaxBreakdown{TaxableValue: decimal.NewFromInt(500000), CGST: decimal.NewFromInt(45000), SGST: decimal.NewFromInt(45000), IGST: decimal.NewFromInt(25000), Cess: decimal.NewFromInt(5000)},
			Inward:  TaxBreakdown{TaxableValue: decimal.NewFromInt(300000), CGST: decimal.NewFromInt(27000), SGST: decimal.NewFromInt(27000), IGST: decimal.NewFromInt(12000), Cess: decimal.NewFromInt(1500)},
			ITC:     ITCBreakdown{CGST: decimal.NewFromInt(17500), SGST: decimal.NewFromInt(17500), IGST: decimal.NewFromInt(10000), Cess: decimal.NewFromInt(5000)},
			TaxPaid: TaxBreakdown{CGST: decimal.NewFromInt(12000), SGST: decimal.NewFromInt(12000), IGST: decimal.NewFromInt(4500), Cess: decimal.NewFromInt(1500)},
		}
	}
	return Aggregate(filingID, tenantID, months)
}

func matchingAuditedFinancials(tables []GSTR9TableData) AuditedFinancials {
	gstr9Turnover := ComputeAggregateTurnover(tables)
	var gstr9Tax TaxBreakdown
	var gstr9ITC ITCBreakdown
	for _, tbl := range tables {
		if tbl.TableNumber == "9" {
			gstr9Tax = TaxBreakdown{CGST: tbl.CGST, SGST: tbl.SGST, IGST: tbl.IGST, Cess: tbl.Cess}
		}
		if tbl.TableNumber == "6F" {
			gstr9ITC = ITCBreakdown{CGST: tbl.CGST, SGST: tbl.SGST, IGST: tbl.IGST, Cess: tbl.Cess}
		}
	}
	return AuditedFinancials{
		Turnover:   gstr9Turnover,
		TaxPayable: gstr9Tax,
		ITCClaimed: gstr9ITC,
	}
}

func TestReconcileTurnover_Match(t *testing.T) {
	tables := buildTestTables()
	audited := matchingAuditedFinancials(tables)

	filing := &GSTR9Filing{ID: uuid.New(), TenantID: uuid.New()}
	mismatches := Reconcile(filing, tables, audited, filing.TenantID, uuid.New())
	assert.Empty(t, mismatches, "matching financials should produce no mismatches")
}

func TestReconcileTurnover_Mismatch(t *testing.T) {
	tables := buildTestTables()
	audited := matchingAuditedFinancials(tables)
	audited.Turnover = audited.Turnover.Add(decimal.NewFromInt(50000))

	filing := &GSTR9Filing{ID: uuid.New(), TenantID: uuid.New()}
	mismatches := Reconcile(filing, tables, audited, filing.TenantID, uuid.New())

	var found bool
	for _, m := range mismatches {
		if m.Description == "Aggregate turnover mismatch" {
			found = true
			assert.Equal(t, "II", m.Section)
			assert.Equal(t, "turnover", m.Category)
			assert.Equal(t, SeverityError, m.Severity)
			assert.True(t, m.Difference.Equal(decimal.NewFromInt(50000)))
		}
	}
	assert.True(t, found, "should detect turnover mismatch")
}

func TestReconcileTurnover_UnbilledRevenue(t *testing.T) {
	tables := buildTestTables()
	audited := matchingAuditedFinancials(tables)
	audited.UnbilledRevenue = decimal.NewFromInt(5000)

	filing := &GSTR9Filing{ID: uuid.New(), TenantID: uuid.New()}
	mismatches := Reconcile(filing, tables, audited, filing.TenantID, uuid.New())

	var found bool
	for _, m := range mismatches {
		if m.Description == "Unbilled revenue not reflected in GSTR-9" {
			found = true
			assert.Equal(t, SeverityWarn, m.Severity)
		}
	}
	assert.True(t, found)
}

func TestReconcileTurnover_DeemedSupply(t *testing.T) {
	tables := buildTestTables()
	audited := matchingAuditedFinancials(tables)
	audited.DeemedSupply = decimal.NewFromInt(20000)

	filing := &GSTR9Filing{ID: uuid.New(), TenantID: uuid.New()}
	mismatches := Reconcile(filing, tables, audited, filing.TenantID, uuid.New())

	var found bool
	for _, m := range mismatches {
		if m.Description == "Deemed supply under Schedule I" {
			found = true
			assert.Equal(t, SeverityError, m.Severity)
		}
	}
	assert.True(t, found)
}

func TestReconcileTurnover_CreditNotesAfterFY(t *testing.T) {
	tables := buildTestTables()
	audited := matchingAuditedFinancials(tables)
	audited.CreditNotesAfterFY = decimal.NewFromInt(500)

	filing := &GSTR9Filing{ID: uuid.New(), TenantID: uuid.New()}
	mismatches := Reconcile(filing, tables, audited, filing.TenantID, uuid.New())

	var found bool
	for _, m := range mismatches {
		if m.Description == "Credit notes issued after FY end" {
			found = true
			assert.True(t, m.Difference.LessThan(decimal.Zero))
		}
	}
	assert.True(t, found)
}

func TestReconcileTax_Match(t *testing.T) {
	tables := buildTestTables()
	audited := matchingAuditedFinancials(tables)

	filing := &GSTR9Filing{ID: uuid.New(), TenantID: uuid.New()}
	mismatches := Reconcile(filing, tables, audited, filing.TenantID, uuid.New())
	assert.Empty(t, mismatches, "matching tax should produce no mismatches")
}

func TestReconcileTax_Mismatch(t *testing.T) {
	tables := buildTestTables()
	audited := matchingAuditedFinancials(tables)
	audited.TaxPayable.CGST = audited.TaxPayable.CGST.Add(decimal.NewFromInt(200))

	filing := &GSTR9Filing{ID: uuid.New(), TenantID: uuid.New()}
	mismatches := Reconcile(filing, tables, audited, filing.TenantID, uuid.New())

	var cgstMismatch bool
	for _, m := range mismatches {
		if m.Description == "CGST tax payable mismatch" {
			cgstMismatch = true
			assert.Equal(t, "III", m.Section)
			assert.Equal(t, "tax", m.Category)
			assert.True(t, m.Difference.Equal(decimal.NewFromInt(200)))
		}
	}
	assert.True(t, cgstMismatch)
}

func TestReconcileITC_Match(t *testing.T) {
	tables := buildTestTables()
	audited := matchingAuditedFinancials(tables)

	filing := &GSTR9Filing{ID: uuid.New(), TenantID: uuid.New()}
	mismatches := Reconcile(filing, tables, audited, filing.TenantID, uuid.New())
	assert.Empty(t, mismatches)
}

func TestReconcileITC_ExcessClaim(t *testing.T) {
	tables := buildTestTables()
	audited := matchingAuditedFinancials(tables)
	audited.ITCClaimed.CGST = audited.ITCClaimed.CGST.Sub(decimal.NewFromInt(20000))

	filing := &GSTR9Filing{ID: uuid.New(), TenantID: uuid.New()}
	mismatches := Reconcile(filing, tables, audited, filing.TenantID, uuid.New())

	var found bool
	for _, m := range mismatches {
		if m.Description == "CGST ITC mismatch" {
			found = true
			assert.Equal(t, "IV", m.Section)
			assert.Equal(t, "itc", m.Category)
			assert.True(t, m.Difference.LessThan(decimal.Zero))
			assert.Contains(t, m.SuggestedAction, "Excess ITC")
		}
	}
	assert.True(t, found)
}

func TestReconcileITC_ShortClaim(t *testing.T) {
	tables := buildTestTables()
	audited := matchingAuditedFinancials(tables)
	audited.ITCClaimed.CGST = audited.ITCClaimed.CGST.Add(decimal.NewFromInt(15000))

	filing := &GSTR9Filing{ID: uuid.New(), TenantID: uuid.New()}
	mismatches := Reconcile(filing, tables, audited, filing.TenantID, uuid.New())

	var found bool
	for _, m := range mismatches {
		if m.Description == "CGST ITC mismatch" {
			found = true
			assert.True(t, m.Difference.GreaterThan(decimal.Zero))
			assert.Contains(t, m.SuggestedAction, "Short ITC")
		}
	}
	assert.True(t, found)
}

func TestReconcile_AllSections(t *testing.T) {
	tables := buildTestTables()

	audited := AuditedFinancials{
		Turnover:           decimal.NewFromInt(99999999),
		TaxPayable:         TaxBreakdown{CGST: decimal.NewFromInt(1), SGST: decimal.NewFromInt(1), IGST: decimal.NewFromInt(1), Cess: decimal.NewFromInt(1)},
		ITCClaimed:         ITCBreakdown{CGST: decimal.NewFromInt(1), SGST: decimal.NewFromInt(1), IGST: decimal.NewFromInt(1), Cess: decimal.NewFromInt(1)},
		UnbilledRevenue:    decimal.NewFromInt(50000),
		UnadjustedAdvances: decimal.NewFromInt(30000),
		DeemedSupply:       decimal.NewFromInt(20000),
		CreditNotesAfterFY: decimal.NewFromInt(10000),
		Sec15_3Adjustments: decimal.NewFromInt(5000),
	}

	filing := &GSTR9Filing{ID: uuid.New(), TenantID: uuid.New()}
	mismatches := Reconcile(filing, tables, audited, filing.TenantID, uuid.New())
	require.NotEmpty(t, mismatches)

	sections := map[string]bool{}
	categories := map[string]bool{}
	for _, m := range mismatches {
		sections[m.Section] = true
		categories[m.Category] = true
		assert.NotEmpty(t, m.Reason)
		assert.NotEmpty(t, m.SuggestedAction)
		assert.NotEqual(t, uuid.Nil, m.ID)
	}
	assert.True(t, sections["II"], "should have Section II (turnover)")
	assert.True(t, sections["III"], "should have Section III (tax)")
	assert.True(t, sections["IV"], "should have Section IV (ITC)")
	assert.True(t, categories["turnover"])
	assert.True(t, categories["tax"])
	assert.True(t, categories["itc"])
}

func TestReconcile_Sec15_3Adjustments(t *testing.T) {
	tables := buildTestTables()
	audited := matchingAuditedFinancials(tables)
	audited.Sec15_3Adjustments = decimal.NewFromInt(500)

	filing := &GSTR9Filing{ID: uuid.New(), TenantID: uuid.New()}
	mismatches := Reconcile(filing, tables, audited, filing.TenantID, uuid.New())

	var found bool
	for _, m := range mismatches {
		if m.Description == "Adjustments per Section 15(3)" {
			found = true
			assert.Equal(t, SeverityWarn, m.Severity)
			assert.Contains(t, m.Reason, "Section 15(3)")
		}
	}
	assert.True(t, found)
}

func TestReconcile_UnadjustedAdvances(t *testing.T) {
	tables := buildTestTables()
	audited := matchingAuditedFinancials(tables)
	audited.UnadjustedAdvances = decimal.NewFromInt(50)

	filing := &GSTR9Filing{ID: uuid.New(), TenantID: uuid.New()}
	mismatches := Reconcile(filing, tables, audited, filing.TenantID, uuid.New())

	var found bool
	for _, m := range mismatches {
		if m.Description == "Unadjusted advances" {
			found = true
			assert.Equal(t, SeverityInfo, m.Severity)
		}
	}
	assert.True(t, found)
}

func TestSeverityBoundary_ExactThresholds(t *testing.T) {
	assert.Equal(t, SeverityInfo, ClassifySeverity(decimal.NewFromInt(99)))
	assert.Equal(t, SeverityWarn, ClassifySeverity(decimal.NewFromInt(100)))
	assert.Equal(t, SeverityWarn, ClassifySeverity(decimal.NewFromInt(9999)))
	assert.Equal(t, SeverityError, ClassifySeverity(decimal.NewFromInt(10000)))
}
