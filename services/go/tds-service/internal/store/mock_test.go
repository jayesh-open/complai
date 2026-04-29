package store

import (
	"context"
	"testing"
	"time"

	"github.com/complai/complai/services/go/tds-service/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	tenant1 = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	tenant2 = uuid.MustParse("22222222-2222-2222-2222-222222222222")
)

func makeDeductee(tenantID uuid.UUID) *domain.Deductee {
	return &domain.Deductee{
		ID:             uuid.New(),
		TenantID:       tenantID,
		VendorID:       uuid.New(),
		Name:           "Vendor",
		PAN:            "ABCCD1234E",
		PANVerified:    true,
		PANStatus:      "VALID",
		DeducteeType:   domain.DeducteeCompany,
		ResidentStatus: domain.Resident,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func TestMockStore_CreateAndGetDeductee(t *testing.T) {
	ms := NewMockStore()
	d := makeDeductee(tenant1)
	require.NoError(t, ms.CreateDeductee(context.Background(), tenant1, d))

	got, err := ms.GetDeductee(context.Background(), tenant1, d.ID)
	require.NoError(t, err)
	assert.Equal(t, d.ID, got.ID)
	assert.Equal(t, d.Name, got.Name)
}

func TestMockStore_GetDeductee_NotFound(t *testing.T) {
	ms := NewMockStore()
	_, err := ms.GetDeductee(context.Background(), tenant1, uuid.New())
	assert.Error(t, err)
}

func TestMockStore_GetDeductee_TenantIsolation(t *testing.T) {
	ms := NewMockStore()
	d := makeDeductee(tenant1)
	ms.CreateDeductee(context.Background(), tenant1, d)

	_, err := ms.GetDeductee(context.Background(), tenant2, d.ID)
	assert.Error(t, err, "should not return deductee belonging to different tenant")
}

func TestMockStore_GetDeducteeByVendor(t *testing.T) {
	ms := NewMockStore()
	d := makeDeductee(tenant1)
	ms.CreateDeductee(context.Background(), tenant1, d)

	got, err := ms.GetDeducteeByVendor(context.Background(), tenant1, d.VendorID)
	require.NoError(t, err)
	assert.Equal(t, d.ID, got.ID)
}

func TestMockStore_GetDeducteeByVendor_NotFound(t *testing.T) {
	ms := NewMockStore()
	_, err := ms.GetDeducteeByVendor(context.Background(), tenant1, uuid.New())
	assert.Error(t, err)
}

func TestMockStore_ListDeductees(t *testing.T) {
	ms := NewMockStore()
	for i := 0; i < 5; i++ {
		ms.CreateDeductee(context.Background(), tenant1, makeDeductee(tenant1))
	}
	ms.CreateDeductee(context.Background(), tenant2, makeDeductee(tenant2))

	list, total, err := ms.ListDeductees(context.Background(), tenant1, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, 5, total)
	assert.Len(t, list, 5)
}

func TestMockStore_ListDeductees_Pagination(t *testing.T) {
	ms := NewMockStore()
	for i := 0; i < 5; i++ {
		ms.CreateDeductee(context.Background(), tenant1, makeDeductee(tenant1))
	}

	list, total, err := ms.ListDeductees(context.Background(), tenant1, 2, 0)
	require.NoError(t, err)
	assert.Equal(t, 5, total)
	assert.Len(t, list, 2)

	list2, _, err := ms.ListDeductees(context.Background(), tenant1, 2, 4)
	require.NoError(t, err)
	assert.Len(t, list2, 1)
}

func TestMockStore_ListDeductees_OffsetBeyondTotal(t *testing.T) {
	ms := NewMockStore()
	ms.CreateDeductee(context.Background(), tenant1, makeDeductee(tenant1))

	list, total, err := ms.ListDeductees(context.Background(), tenant1, 10, 100)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Nil(t, list)
}

func TestMockStore_UpsertDeductee_Update(t *testing.T) {
	ms := NewMockStore()
	d := makeDeductee(tenant1)
	ms.CreateDeductee(context.Background(), tenant1, d)

	updated := &domain.Deductee{
		ID:           uuid.New(),
		TenantID:     tenant1,
		VendorID:     d.VendorID,
		Name:         "Updated Name",
		PAN:          "XYZCD5678F",
		PANStatus:    "VALID",
		DeducteeType: domain.DeducteeIndividual,
	}
	require.NoError(t, ms.UpsertDeductee(context.Background(), tenant1, updated))

	got, err := ms.GetDeductee(context.Background(), tenant1, d.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", got.Name)
	assert.Equal(t, "XYZCD5678F", got.PAN)
	assert.Equal(t, domain.DeducteeIndividual, got.DeducteeType)
}

func TestMockStore_UpsertDeductee_Insert(t *testing.T) {
	ms := NewMockStore()
	d := makeDeductee(tenant1)
	require.NoError(t, ms.UpsertDeductee(context.Background(), tenant1, d))

	got, err := ms.GetDeductee(context.Background(), tenant1, d.ID)
	require.NoError(t, err)
	assert.Equal(t, d.Name, got.Name)
}

func TestMockStore_CreateAndGetEntry(t *testing.T) {
	ms := NewMockStore()
	e := &domain.TDSEntry{
		ID: uuid.New(), TenantID: tenant1, DeducteeID: uuid.New(),
		Section: domain.Section393_1, PaymentCode: domain.CodeContractorOther,
		FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
		TransactionDate: time.Now(), GrossAmount: decimal.NewFromInt(100000),
		TDSRate: decimal.NewFromFloat(0.02), TDSAmount: decimal.NewFromInt(2000),
		TotalTax: decimal.NewFromInt(2000), Status: domain.StatusPending,
		NatureOfPayment: "contractor",
	}
	require.NoError(t, ms.CreateEntry(context.Background(), tenant1, e))

	got, err := ms.GetEntry(context.Background(), tenant1, e.ID)
	require.NoError(t, err)
	assert.Equal(t, e.ID, got.ID)
}

func TestMockStore_GetEntry_NotFound(t *testing.T) {
	ms := NewMockStore()
	_, err := ms.GetEntry(context.Background(), tenant1, uuid.New())
	assert.Error(t, err)
}

func TestMockStore_GetEntry_TenantIsolation(t *testing.T) {
	ms := NewMockStore()
	e := &domain.TDSEntry{
		ID: uuid.New(), TenantID: tenant1, DeducteeID: uuid.New(),
		Section: domain.Section393_1, PaymentCode: domain.CodeContractorOther,
		FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
		TransactionDate: time.Now(), GrossAmount: decimal.NewFromInt(100000),
		TDSRate: decimal.NewFromFloat(0.02), TDSAmount: decimal.NewFromInt(2000),
		TotalTax: decimal.NewFromInt(2000), Status: domain.StatusPending,
	}
	ms.CreateEntry(context.Background(), tenant1, e)

	_, err := ms.GetEntry(context.Background(), tenant2, e.ID)
	assert.Error(t, err)
}

func TestMockStore_ListEntries_Filters(t *testing.T) {
	ms := NewMockStore()
	e1 := &domain.TDSEntry{
		ID: uuid.New(), TenantID: tenant1,
		Section: domain.Section393_1, PaymentCode: domain.CodeContractorOther,
		FinancialYear: "2026-27", Quarter: "Q1",
		TransactionDate: time.Now(), GrossAmount: decimal.NewFromInt(50000),
		TDSRate: decimal.NewFromFloat(0.02), TDSAmount: decimal.NewFromInt(1000),
		TotalTax: decimal.NewFromInt(1000), Status: domain.StatusPending,
	}
	e2 := &domain.TDSEntry{
		ID: uuid.New(), TenantID: tenant1,
		Section: domain.Section393_1, PaymentCode: domain.CodeProfessional,
		FinancialYear: "2026-27", Quarter: "Q2",
		TransactionDate: time.Now(), GrossAmount: decimal.NewFromInt(80000),
		TDSRate: decimal.NewFromFloat(0.10), TDSAmount: decimal.NewFromInt(8000),
		TotalTax: decimal.NewFromInt(8000), Status: domain.StatusPending,
	}
	ms.CreateEntry(context.Background(), tenant1, e1)
	ms.CreateEntry(context.Background(), tenant1, e2)

	all, total, err := ms.ListEntries(context.Background(), tenant1, "", "", 10, 0)
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, all, 2)

	byFY, total, err := ms.ListEntries(context.Background(), tenant1, "2026-27", "", 10, 0)
	require.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, byFY, 2)

	byQ, total, err := ms.ListEntries(context.Background(), tenant1, "", "Q1", 10, 0)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, byQ, 1)
}

func TestMockStore_ListEntries_Pagination(t *testing.T) {
	ms := NewMockStore()
	for i := 0; i < 5; i++ {
		e := &domain.TDSEntry{
			ID: uuid.New(), TenantID: tenant1,
			Section: domain.Section393_1, PaymentCode: domain.CodeContractorOther,
			FinancialYear: "2026-27", Quarter: "Q1",
			TransactionDate: time.Now(), GrossAmount: decimal.NewFromInt(10000),
			TDSRate: decimal.NewFromFloat(0.02), TDSAmount: decimal.NewFromInt(200),
			TotalTax: decimal.NewFromInt(200), Status: domain.StatusPending,
		}
		ms.CreateEntry(context.Background(), tenant1, e)
	}

	list, total, err := ms.ListEntries(context.Background(), tenant1, "", "", 2, 0)
	require.NoError(t, err)
	assert.Equal(t, 5, total)
	assert.Len(t, list, 2)

	list2, _, _ := ms.ListEntries(context.Background(), tenant1, "", "", 10, 4)
	assert.Len(t, list2, 1)
}

func TestMockStore_GetAggregate_Empty(t *testing.T) {
	ms := NewMockStore()
	agg, err := ms.GetAggregate(context.Background(), tenant1, uuid.New(), domain.CodeContractorOther, "2026-27")
	require.NoError(t, err)
	assert.True(t, agg.TotalPaid.IsZero())
	assert.True(t, agg.TotalTDS.IsZero())
}

func TestMockStore_UpsertAndGetAggregate(t *testing.T) {
	ms := NewMockStore()
	did := uuid.New()
	agg := &domain.TDSAggregate{
		ID: uuid.New(), TenantID: tenant1, DeducteeID: did,
		PaymentCode: domain.CodeContractorOther, FinancialYear: "2026-27",
		TotalPaid: decimal.NewFromInt(100000), TotalTDS: decimal.NewFromInt(2000),
		TransactionCount: 1,
	}
	require.NoError(t, ms.UpsertAggregate(context.Background(), tenant1, agg))

	got, err := ms.GetAggregate(context.Background(), tenant1, did, domain.CodeContractorOther, "2026-27")
	require.NoError(t, err)
	assert.True(t, got.TotalPaid.Equal(decimal.NewFromInt(100000)))
	assert.Equal(t, 1, got.TransactionCount)
}

func TestMockStore_GetSummary(t *testing.T) {
	ms := NewMockStore()
	d := makeDeductee(tenant1)
	ms.CreateDeductee(context.Background(), tenant1, d)

	e := &domain.TDSEntry{
		ID: uuid.New(), TenantID: tenant1, DeducteeID: d.ID,
		Section: domain.Section393_1, PaymentCode: domain.CodeContractorOther,
		FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
		TransactionDate: time.Now(), GrossAmount: decimal.NewFromInt(100000),
		TDSRate: decimal.NewFromFloat(0.02), TDSAmount: decimal.NewFromInt(2000),
		TotalTax: decimal.NewFromInt(2000), Status: domain.StatusPending,
	}
	ms.CreateEntry(context.Background(), tenant1, e)

	dep := &domain.TDSEntry{
		ID: uuid.New(), TenantID: tenant1, DeducteeID: d.ID,
		Section: domain.Section393_1, PaymentCode: domain.CodeContractorOther,
		FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
		TransactionDate: time.Now(), GrossAmount: decimal.NewFromInt(50000),
		TDSRate: decimal.NewFromFloat(0.02), TDSAmount: decimal.NewFromInt(1000),
		TotalTax: decimal.NewFromInt(1000), Status: domain.StatusDeposited,
	}
	ms.CreateEntry(context.Background(), tenant1, dep)

	sum, err := ms.GetSummary(context.Background(), tenant1, "2026-27")
	require.NoError(t, err)
	assert.Equal(t, 1, sum.TotalDeductees)
	assert.Equal(t, 2, sum.TotalEntries)
	assert.True(t, sum.TotalTDSDeducted.Equal(decimal.NewFromInt(3000)))
	assert.True(t, sum.TotalTDSDeposited.Equal(decimal.NewFromInt(1000)))
	assert.True(t, sum.PendingDeposit.Equal(decimal.NewFromInt(2000)))
	assert.Equal(t, 2, sum.EntriesByPaymentCode[domain.CodeContractorOther])
	assert.Equal(t, 1, sum.EntriesByStatus[domain.StatusPending])
	assert.Equal(t, 1, sum.EntriesByStatus[domain.StatusDeposited])
}

func TestMockStore_GetSummary_Empty(t *testing.T) {
	ms := NewMockStore()
	sum, err := ms.GetSummary(context.Background(), tenant1, "2026-27")
	require.NoError(t, err)
	assert.Equal(t, 0, sum.TotalEntries)
	assert.True(t, sum.TotalTDSDeducted.IsZero())
}

func TestMockStore_AggKey(t *testing.T) {
	id := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	key := aggKey(id, domain.CodeContractorOther, "2026-27")
	assert.Contains(t, key, "1024")
	assert.Contains(t, key, "2026-27")
}

func TestMockStore_ListEntries_OffsetBeyondTotal(t *testing.T) {
	ms := NewMockStore()
	e := &domain.TDSEntry{
		ID: uuid.New(), TenantID: tenant1,
		Section: domain.Section393_1, PaymentCode: domain.CodeContractorOther,
		FinancialYear: "2026-27", Quarter: "Q1",
		TransactionDate: time.Now(), GrossAmount: decimal.NewFromInt(10000),
		TDSRate: decimal.NewFromFloat(0.02), TDSAmount: decimal.NewFromInt(200),
		TotalTax: decimal.NewFromInt(200), Status: domain.StatusPending,
	}
	ms.CreateEntry(context.Background(), tenant1, e)

	list, total, err := ms.ListEntries(context.Background(), tenant1, "", "", 10, 100)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Nil(t, list)
}

func TestMockStore_GetSummary_CrossTenantFiltering(t *testing.T) {
	ms := NewMockStore()
	ms.CreateDeductee(context.Background(), tenant1, makeDeductee(tenant1))
	ms.CreateDeductee(context.Background(), tenant2, makeDeductee(tenant2))

	e1 := &domain.TDSEntry{
		ID: uuid.New(), TenantID: tenant1,
		Section: domain.Section393_1, PaymentCode: domain.CodeContractorOther,
		FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
		TransactionDate: time.Now(), GrossAmount: decimal.NewFromInt(100000),
		TDSRate: decimal.NewFromFloat(0.02), TDSAmount: decimal.NewFromInt(2000),
		TotalTax: decimal.NewFromInt(2000), Status: domain.StatusPending,
	}
	e2 := &domain.TDSEntry{
		ID: uuid.New(), TenantID: tenant2,
		Section: domain.Section393_1, PaymentCode: domain.CodeProfessional,
		FinancialYear: "2026-27", TaxYear: "2026-27", Quarter: "Q1",
		TransactionDate: time.Now(), GrossAmount: decimal.NewFromInt(50000),
		TDSRate: decimal.NewFromFloat(0.10), TDSAmount: decimal.NewFromInt(5000),
		TotalTax: decimal.NewFromInt(5000), Status: domain.StatusPending,
	}
	eDiffFY := &domain.TDSEntry{
		ID: uuid.New(), TenantID: tenant1,
		Section: domain.Section393_1, PaymentCode: domain.CodeContractorOther,
		FinancialYear: "2025-26", TaxYear: "2025-26", Quarter: "Q4",
		TransactionDate: time.Now(), GrossAmount: decimal.NewFromInt(30000),
		TDSRate: decimal.NewFromFloat(0.02), TDSAmount: decimal.NewFromInt(600),
		TotalTax: decimal.NewFromInt(600), Status: domain.StatusDeposited,
	}
	ms.CreateEntry(context.Background(), tenant1, e1)
	ms.CreateEntry(context.Background(), tenant2, e2)
	ms.CreateEntry(context.Background(), tenant1, eDiffFY)

	sum, err := ms.GetSummary(context.Background(), tenant1, "2026-27")
	require.NoError(t, err)
	assert.Equal(t, 1, sum.TotalDeductees, "should only count tenant1 deductees")
	assert.Equal(t, 1, sum.TotalEntries, "should only count tenant1 entries for FY 2026-27")
	assert.True(t, sum.TotalTDSDeducted.Equal(decimal.NewFromInt(2000)))
}

func TestMockStore_ListEntries_TenantIsolation(t *testing.T) {
	ms := NewMockStore()
	e := &domain.TDSEntry{
		ID: uuid.New(), TenantID: tenant1,
		Section: domain.Section393_1, PaymentCode: domain.CodeContractorOther,
		FinancialYear: "2026-27", Quarter: "Q1",
		TransactionDate: time.Now(), GrossAmount: decimal.NewFromInt(10000),
		TDSRate: decimal.NewFromFloat(0.02), TDSAmount: decimal.NewFromInt(200),
		TotalTax: decimal.NewFromInt(200), Status: domain.StatusPending,
	}
	ms.CreateEntry(context.Background(), tenant1, e)

	list, total, err := ms.ListEntries(context.Background(), tenant2, "", "", 10, 0)
	require.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Empty(t, list)
}
