package store

import (
	"context"
	"testing"
	"time"

	"github.com/complai/complai/services/go/gstr9-service/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testGSTR9FilingID = uuid.MustParse("33333333-3333-3333-3333-333333333333")

var testTenant = uuid.MustParse("11111111-1111-1111-1111-111111111111")

func TestMockStore_FilingCRUD(t *testing.T) {
	ms := NewMockStore()
	ctx := context.Background()

	filing := &domain.GSTR9Filing{
		ID: uuid.New(), TenantID: testTenant, GSTIN: "27AABCU9603R1ZM",
		FinancialYear: "2025-26", Status: domain.FilingStatusDraft,
		AggregateTurnover: decimal.Zero, RequestID: uuid.New(),
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}

	require.NoError(t, ms.CreateFiling(ctx, testTenant, filing))

	got, err := ms.GetFiling(ctx, testTenant, filing.ID)
	require.NoError(t, err)
	assert.Equal(t, filing.GSTIN, got.GSTIN)

	require.NoError(t, ms.UpdateFilingStatus(ctx, testTenant, filing.ID, domain.FilingStatusSaved))
	got, _ = ms.GetFiling(ctx, testTenant, filing.ID)
	assert.Equal(t, domain.FilingStatusSaved, got.Status)
}

func TestMockStore_DuplicateFiling(t *testing.T) {
	ms := NewMockStore()
	ctx := context.Background()
	f := &domain.GSTR9Filing{
		ID: uuid.New(), TenantID: testTenant, GSTIN: "27AABCU9603R1ZM",
		FinancialYear: "2025-26", Status: domain.FilingStatusDraft,
		RequestID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	require.NoError(t, ms.CreateFiling(ctx, testTenant, f))

	f2 := &domain.GSTR9Filing{
		ID: uuid.New(), TenantID: testTenant, GSTIN: "27AABCU9603R1ZM",
		FinancialYear: "2025-26", Status: domain.FilingStatusDraft,
		RequestID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	err := ms.CreateFiling(ctx, testTenant, f2)
	assert.ErrorIs(t, err, domain.ErrDuplicateFiling)
}

func TestMockStore_CrossTenantIsolation(t *testing.T) {
	ms := NewMockStore()
	ctx := context.Background()
	otherTenant := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	f := &domain.GSTR9Filing{
		ID: uuid.New(), TenantID: testTenant, GSTIN: "27AABCU9603R1ZM",
		FinancialYear: "2025-26", Status: domain.FilingStatusDraft,
		RequestID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	ms.CreateFiling(ctx, testTenant, f)

	_, err := ms.GetFiling(ctx, otherTenant, f.ID)
	assert.ErrorIs(t, err, domain.ErrFilingNotFound)

	err = ms.UpdateFilingStatus(ctx, otherTenant, f.ID, domain.FilingStatusSaved)
	assert.ErrorIs(t, err, domain.ErrFilingNotFound)
}

func TestMockStore_ListFilings(t *testing.T) {
	ms := NewMockStore()
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		f := &domain.GSTR9Filing{
			ID: uuid.New(), TenantID: testTenant,
			GSTIN: "27AABCU9603R1ZM", FinancialYear: "2025-26",
			Status: domain.FilingStatusDraft, RequestID: uuid.New(),
			CreatedAt: time.Now(), UpdatedAt: time.Now(),
		}
		if i > 0 {
			f.GSTIN = "29AABCU9603R1ZM"
		}
		if i > 1 {
			f.FinancialYear = "2024-25"
		}
		ms.CreateFiling(ctx, testTenant, f)
	}

	all, total, err := ms.ListFilings(ctx, testTenant, "", "", 50, 0)
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Equal(t, 3, len(all))

	filtered, total, _ := ms.ListFilings(ctx, testTenant, "27AABCU9603R1ZM", "", 50, 0)
	assert.Equal(t, 1, total)
	assert.Equal(t, 1, len(filtered))
}

func TestMockStore_TableDataCRUD(t *testing.T) {
	ms := NewMockStore()
	ctx := context.Background()
	filingID := uuid.New()

	td := &domain.GSTR9TableData{
		ID: uuid.New(), TenantID: testTenant, FilingID: filingID,
		PartNumber: 1, TableNumber: "4A", Description: "test",
		TaxableValue: decimal.NewFromInt(100), CGST: decimal.NewFromInt(9),
	}
	require.NoError(t, ms.CreateTableData(ctx, testTenant, td))

	list, err := ms.ListTableData(ctx, testTenant, filingID)
	require.NoError(t, err)
	assert.Equal(t, 1, len(list))

	require.NoError(t, ms.DeleteTableData(ctx, testTenant, filingID))
	list, _ = ms.ListTableData(ctx, testTenant, filingID)
	assert.Empty(t, list)
}

func TestMockStore_AuditLogCRUD(t *testing.T) {
	ms := NewMockStore()
	ctx := context.Background()
	filingID := uuid.New()

	log := &domain.GSTR9AuditLog{
		ID: uuid.New(), TenantID: testTenant, FilingID: filingID,
		Action: "created", Details: "test", ActorID: testTenant,
		CreatedAt: time.Now(),
	}
	require.NoError(t, ms.CreateAuditLog(ctx, testTenant, log))

	logs, err := ms.ListAuditLogs(ctx, testTenant, filingID)
	require.NoError(t, err)
	assert.Equal(t, 1, len(logs))
}

func TestMockStore_GSTR9CFilingCRUD(t *testing.T) {
	ms := NewMockStore()
	ctx := context.Background()

	filing := &domain.GSTR9CFiling{
		ID: uuid.New(), TenantID: testTenant, GSTR9FilingID: testGSTR9FilingID,
		Status: domain.GSTR9CStatusDraft, AuditedTurnover: decimal.NewFromInt(60000000),
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	require.NoError(t, ms.CreateGSTR9CFiling(ctx, testTenant, filing))

	got, err := ms.GetGSTR9CFiling(ctx, testTenant, filing.ID)
	require.NoError(t, err)
	assert.Equal(t, filing.GSTR9FilingID, got.GSTR9FilingID)

	gotByGSTR9, err := ms.GetGSTR9CFilingByGSTR9ID(ctx, testTenant, testGSTR9FilingID)
	require.NoError(t, err)
	assert.Equal(t, filing.ID, gotByGSTR9.ID)

	require.NoError(t, ms.UpdateGSTR9CStatus(ctx, testTenant, filing.ID, domain.GSTR9CStatusReconciled))
	got, _ = ms.GetGSTR9CFiling(ctx, testTenant, filing.ID)
	assert.Equal(t, domain.GSTR9CStatusReconciled, got.Status)

	require.NoError(t, ms.UpdateGSTR9CUnreconciled(ctx, testTenant, filing.ID, decimal.NewFromInt(5000)))
	got, _ = ms.GetGSTR9CFiling(ctx, testTenant, filing.ID)
	assert.True(t, got.UnreconciledAmount.Equal(decimal.NewFromInt(5000)))
}

func TestMockStore_GSTR9CDuplicate(t *testing.T) {
	ms := NewMockStore()
	ctx := context.Background()

	filing := &domain.GSTR9CFiling{
		ID: uuid.New(), TenantID: testTenant, GSTR9FilingID: testGSTR9FilingID,
		Status: domain.GSTR9CStatusDraft, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	require.NoError(t, ms.CreateGSTR9CFiling(ctx, testTenant, filing))

	dup := &domain.GSTR9CFiling{
		ID: uuid.New(), TenantID: testTenant, GSTR9FilingID: testGSTR9FilingID,
		Status: domain.GSTR9CStatusDraft, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	err := ms.CreateGSTR9CFiling(ctx, testTenant, dup)
	assert.ErrorIs(t, err, domain.ErrGSTR9CDuplicate)
}

func TestMockStore_GSTR9CCrossTenant(t *testing.T) {
	ms := NewMockStore()
	ctx := context.Background()
	otherTenant := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	filing := &domain.GSTR9CFiling{
		ID: uuid.New(), TenantID: testTenant, GSTR9FilingID: testGSTR9FilingID,
		Status: domain.GSTR9CStatusDraft, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	ms.CreateGSTR9CFiling(ctx, testTenant, filing)

	_, err := ms.GetGSTR9CFiling(ctx, otherTenant, filing.ID)
	assert.ErrorIs(t, err, domain.ErrGSTR9CNotFound)

	_, err = ms.GetGSTR9CFilingByGSTR9ID(ctx, otherTenant, testGSTR9FilingID)
	assert.ErrorIs(t, err, domain.ErrGSTR9CNotFound)
}

func TestMockStore_CertifyGSTR9C(t *testing.T) {
	ms := NewMockStore()
	ctx := context.Background()

	filing := &domain.GSTR9CFiling{
		ID: uuid.New(), TenantID: testTenant, GSTR9FilingID: testGSTR9FilingID,
		Status: domain.GSTR9CStatusReconciled, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	ms.CreateGSTR9CFiling(ctx, testTenant, filing)

	require.NoError(t, ms.CertifyGSTR9C(ctx, testTenant, filing.ID, testTenant))

	got, _ := ms.GetGSTR9CFiling(ctx, testTenant, filing.ID)
	assert.Equal(t, domain.GSTR9CStatusCertified, got.Status)
	assert.True(t, got.IsSelfCertified)
	assert.NotNil(t, got.CertifiedAt)
	assert.NotNil(t, got.CertifiedBy)
}

func TestMockStore_MismatchCRUD(t *testing.T) {
	ms := NewMockStore()
	ctx := context.Background()
	gstr9cID := uuid.New()

	mm := &domain.GSTR9CMismatch{
		ID: uuid.New(), TenantID: testTenant, GSTR9CFilingID: gstr9cID,
		Section: "II", Category: "turnover", Description: "test mismatch",
		BooksAmount: decimal.NewFromInt(100000), GSTR9Amount: decimal.NewFromInt(90000),
		Difference: decimal.NewFromInt(10000), Severity: domain.SeverityError,
		Reason: "test", SuggestedAction: "verify",
		CreatedAt: time.Now(),
	}
	require.NoError(t, ms.CreateMismatch(ctx, testTenant, mm))

	got, err := ms.GetMismatch(ctx, testTenant, mm.ID)
	require.NoError(t, err)
	assert.Equal(t, mm.Description, got.Description)

	list, err := ms.ListMismatches(ctx, testTenant, gstr9cID)
	require.NoError(t, err)
	assert.Equal(t, 1, len(list))

	require.NoError(t, ms.ResolveMismatch(ctx, testTenant, mm.ID, "verified ok", testTenant))
	got, _ = ms.GetMismatch(ctx, testTenant, mm.ID)
	assert.True(t, got.Resolved)
	assert.Equal(t, "verified ok", got.ResolvedReason)
	assert.NotNil(t, got.ResolvedAt)
}

func TestMockStore_MismatchCrossTenant(t *testing.T) {
	ms := NewMockStore()
	ctx := context.Background()
	otherTenant := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	mm := &domain.GSTR9CMismatch{
		ID: uuid.New(), TenantID: testTenant, GSTR9CFilingID: uuid.New(),
		Severity: domain.SeverityInfo, CreatedAt: time.Now(),
	}
	ms.CreateMismatch(ctx, testTenant, mm)

	_, err := ms.GetMismatch(ctx, otherTenant, mm.ID)
	assert.ErrorIs(t, err, domain.ErrMismatchNotFound)

	err = ms.ResolveMismatch(ctx, otherTenant, mm.ID, "reason", otherTenant)
	assert.ErrorIs(t, err, domain.ErrMismatchNotFound)
}

func TestMockStore_DeleteMismatches(t *testing.T) {
	ms := NewMockStore()
	ctx := context.Background()
	gstr9cID := uuid.New()

	for i := 0; i < 3; i++ {
		mm := &domain.GSTR9CMismatch{
			ID: uuid.New(), TenantID: testTenant, GSTR9CFilingID: gstr9cID,
			Severity: domain.SeverityWarn, CreatedAt: time.Now(),
		}
		ms.CreateMismatch(ctx, testTenant, mm)
	}

	list, _ := ms.ListMismatches(ctx, testTenant, gstr9cID)
	assert.Equal(t, 3, len(list))

	require.NoError(t, ms.DeleteMismatches(ctx, testTenant, gstr9cID))
	list, _ = ms.ListMismatches(ctx, testTenant, gstr9cID)
	assert.Empty(t, list)
}
