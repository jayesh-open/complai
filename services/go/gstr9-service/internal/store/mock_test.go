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
