package store

import (
	"context"
	"testing"
	"time"

	"github.com/complai/complai/services/go/itr-service/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockStore_ImplementsRepository(t *testing.T) {
	var _ Repository = (*MockStore)(nil)
}

func TestMockStore_TaxpayerCRUD(t *testing.T) {
	s := NewMockStore()
	ctx := context.Background()
	tenantID := uuid.New()

	tp := &domain.Taxpayer{
		ID:              uuid.New(),
		TenantID:        tenantID,
		PAN:             "ABCDE1234F",
		Name:            "Test User",
		DateOfBirth:     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		AssesseeType:    domain.AssesseeIndividual,
		ResidencyStatus: domain.Resident,
	}

	require.NoError(t, s.CreateTaxpayer(ctx, tenantID, tp))

	got, err := s.GetTaxpayer(ctx, tenantID, tp.ID)
	require.NoError(t, err)
	assert.Equal(t, "ABCDE1234F", got.PAN)

	got2, err := s.GetTaxpayerByPAN(ctx, tenantID, "ABCDE1234F")
	require.NoError(t, err)
	assert.Equal(t, tp.ID, got2.ID)

	list, total, err := s.ListTaxpayers(ctx, tenantID, 50, 0)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, list, 1)

	_, err = s.GetTaxpayer(ctx, uuid.New(), tp.ID)
	assert.Error(t, err, "should fail for wrong tenant")
}

func TestMockStore_FilingCRUD(t *testing.T) {
	s := NewMockStore()
	ctx := context.Background()
	tenantID := uuid.New()

	f := &domain.ITRFiling{
		ID:             uuid.New(),
		TenantID:       tenantID,
		TaxpayerID:     uuid.New(),
		PAN:            "ABCDE1234F",
		TaxYear:        "2026-27",
		FormType:       domain.FormITR1,
		RegimeSelected: domain.NewRegime,
		Status:         domain.StatusDraft,
		GrossIncome:    decimal.NewFromInt(1000000),
		IdempotencyKey: tenantID.String() + ":ABCDE1234F:2026-27:ITR-1",
	}

	require.NoError(t, s.CreateFiling(ctx, tenantID, f))

	got, err := s.GetFiling(ctx, tenantID, f.ID)
	require.NoError(t, err)
	assert.Equal(t, domain.StatusDraft, got.Status)

	require.NoError(t, s.UpdateFilingStatus(ctx, tenantID, f.ID, domain.StatusFiled, "ARN123", "ACK456", ""))
	got2, _ := s.GetFiling(ctx, tenantID, f.ID)
	assert.Equal(t, domain.StatusFiled, got2.Status)
	assert.Equal(t, "ARN123", got2.ARN)

	list, total, err := s.ListFilings(ctx, tenantID, "2026-27", 50, 0)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, list, 1)

	err = s.CreateFiling(ctx, tenantID, f)
	assert.Error(t, err, "duplicate idempotency key should fail")
}

func TestMockStore_IncomeEntries(t *testing.T) {
	s := NewMockStore()
	ctx := context.Background()
	tenantID := uuid.New()
	filingID := uuid.New()

	e := &domain.IncomeEntry{
		ID:          uuid.New(),
		TenantID:    tenantID,
		FilingID:    filingID,
		Head:        domain.HeadSalary,
		Section:     "392",
		Description: "salary",
		Amount:      decimal.NewFromInt(1200000),
	}
	require.NoError(t, s.CreateIncomeEntry(ctx, tenantID, e))

	entries, err := s.ListIncomeEntries(ctx, tenantID, filingID)
	require.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, domain.HeadSalary, entries[0].Head)
}

func TestMockStore_Deductions(t *testing.T) {
	s := NewMockStore()
	ctx := context.Background()
	tenantID := uuid.New()
	filingID := uuid.New()

	d := &domain.Deduction{
		ID:       uuid.New(),
		TenantID: tenantID,
		FilingID: filingID,
		Section:  "VI-A",
		Label:    "Life insurance",
		Claimed:  decimal.NewFromInt(200000),
		Allowed:  decimal.NewFromInt(150000),
		MaxLimit: decimal.NewFromInt(150000),
	}
	require.NoError(t, s.CreateDeduction(ctx, tenantID, d))

	deds, err := s.ListDeductions(ctx, tenantID, filingID)
	require.NoError(t, err)
	assert.Len(t, deds, 1)
}

func TestMockStore_TaxComputation(t *testing.T) {
	s := NewMockStore()
	ctx := context.Background()
	tenantID := uuid.New()
	filingID := uuid.New()

	tc := &domain.TaxComputation{
		ID:            uuid.New(),
		TenantID:      tenantID,
		FilingID:      filingID,
		RegimeType:    domain.NewRegime,
		GrossIncome:   decimal.NewFromInt(1200000),
		TaxableIncome: decimal.NewFromInt(1125000),
		BaseTax:       decimal.NewFromInt(52500),
	}
	require.NoError(t, s.SaveTaxComputation(ctx, tenantID, tc))

	got, err := s.GetTaxComputation(ctx, tenantID, filingID)
	require.NoError(t, err)
	assert.Equal(t, domain.NewRegime, got.RegimeType)

	_, err = s.GetTaxComputation(ctx, tenantID, uuid.New())
	assert.Error(t, err, "should fail for unknown filing")
}

func TestMockStore_TDSCredits(t *testing.T) {
	s := NewMockStore()
	ctx := context.Background()
	tenantID := uuid.New()
	filingID := uuid.New()

	c := &domain.TDSCredit{
		ID:           uuid.New(),
		TenantID:     tenantID,
		FilingID:     filingID,
		DeductorTAN:  "MUMB12345A",
		DeductorName: "Infosys Ltd",
		Section:      "392",
		TDSAmount:    decimal.NewFromInt(50000),
		GrossPayment: decimal.NewFromInt(500000),
		TaxYear:      "2026-27",
	}
	require.NoError(t, s.CreateTDSCredit(ctx, tenantID, c))

	credits, err := s.ListTDSCredits(ctx, tenantID, filingID)
	require.NoError(t, err)
	assert.Len(t, credits, 1)
}

func TestMockStore_AISReconciliation(t *testing.T) {
	s := NewMockStore()
	ctx := context.Background()
	tenantID := uuid.New()
	filingID := uuid.New()

	r := &domain.AISReconciliation{
		ID:             uuid.New(),
		TenantID:       tenantID,
		FilingID:       filingID,
		PAN:            "ABCDE1234F",
		TaxYear:        "2026-27",
		SourceType:     "TDS",
		ReportedAmount: decimal.NewFromInt(500000),
		AISAmount:      decimal.NewFromInt(500000),
		Discrepancy:    decimal.Zero,
		Status:         "MATCHED",
	}
	require.NoError(t, s.CreateAISReconciliation(ctx, tenantID, r))

	recons, err := s.ListAISReconciliations(ctx, tenantID, filingID)
	require.NoError(t, err)
	assert.Len(t, recons, 1)
}

func TestMockStore_FilingByIdempotencyKey(t *testing.T) {
	s := NewMockStore()
	ctx := context.Background()
	tenantID := uuid.New()

	f := &domain.ITRFiling{
		ID:             uuid.New(),
		TenantID:       tenantID,
		TaxpayerID:     uuid.New(),
		PAN:            "ABCDE1234F",
		TaxYear:        "2026-27",
		FormType:       domain.FormITR1,
		RegimeSelected: domain.NewRegime,
		Status:         domain.StatusDraft,
		IdempotencyKey: "test-key-123",
	}
	require.NoError(t, s.CreateFiling(ctx, tenantID, f))

	got, err := s.GetFilingByIdempotencyKey(ctx, tenantID, "test-key-123")
	require.NoError(t, err)
	assert.Equal(t, f.ID, got.ID)

	_, err = s.GetFilingByIdempotencyKey(ctx, tenantID, "nonexistent")
	assert.Error(t, err)
}

func TestMockStore_ListTaxpayers_Pagination(t *testing.T) {
	s := NewMockStore()
	ctx := context.Background()
	tenantID := uuid.New()

	for i := 0; i < 5; i++ {
		tp := &domain.Taxpayer{
			ID:       uuid.New(),
			TenantID: tenantID,
			PAN:      "ABCDE123" + string(rune('0'+i)) + "F",
			Name:     "User",
		}
		require.NoError(t, s.CreateTaxpayer(ctx, tenantID, tp))
	}

	list, total, err := s.ListTaxpayers(ctx, tenantID, 2, 0)
	require.NoError(t, err)
	assert.Equal(t, 5, total)
	assert.Len(t, list, 2)

	list2, _, _ := s.ListTaxpayers(ctx, tenantID, 50, 10)
	assert.Nil(t, list2)
}

func TestMockStore_UpdateFilingStatus_NotFound(t *testing.T) {
	s := NewMockStore()
	ctx := context.Background()
	err := s.UpdateFilingStatus(ctx, uuid.New(), uuid.New(), domain.StatusFiled, "", "", "")
	assert.Error(t, err)
}

func TestMockStore_UpdateFilingStatus_WrongTenant(t *testing.T) {
	s := NewMockStore()
	ctx := context.Background()
	tenantID := uuid.New()

	f := &domain.ITRFiling{
		ID:             uuid.New(),
		TenantID:       tenantID,
		TaxYear:        "2026-27",
		FormType:       domain.FormITR1,
		RegimeSelected: domain.NewRegime,
		Status:         domain.StatusDraft,
		IdempotencyKey: "upd-test",
	}
	require.NoError(t, s.CreateFiling(ctx, tenantID, f))

	err := s.UpdateFilingStatus(ctx, uuid.New(), f.ID, domain.StatusFiled, "", "", "")
	assert.Error(t, err, "should fail for wrong tenant")
}

func TestMockStore_GetTaxpayerByPAN_NotFound(t *testing.T) {
	s := NewMockStore()
	ctx := context.Background()
	_, err := s.GetTaxpayerByPAN(ctx, uuid.New(), "ZZZZZ9999Z")
	assert.Error(t, err)
}

func TestMockStore_ListFilings_FilterByTaxYear(t *testing.T) {
	s := NewMockStore()
	ctx := context.Background()
	tenantID := uuid.New()

	for _, ty := range []string{"2026-27", "2027-28"} {
		f := &domain.ITRFiling{
			ID:             uuid.New(),
			TenantID:       tenantID,
			TaxYear:        ty,
			FormType:       domain.FormITR1,
			RegimeSelected: domain.NewRegime,
			Status:         domain.StatusDraft,
			IdempotencyKey: tenantID.String() + ":" + ty,
		}
		require.NoError(t, s.CreateFiling(ctx, tenantID, f))
	}

	list, total, err := s.ListFilings(ctx, tenantID, "2026-27", 50, 0)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, list, 1)

	listAll, totalAll, _ := s.ListFilings(ctx, tenantID, "", 50, 0)
	assert.Equal(t, 2, totalAll)
	assert.Len(t, listAll, 2)
}
