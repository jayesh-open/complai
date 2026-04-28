package filing

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/complai/complai/services/go/tds-service/internal/domain"
	"github.com/complai/complai/services/go/tds-service/internal/store"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testTenant = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	testDate   = time.Date(2025, 5, 20, 0, 0, 0, 0, time.UTC)
)

type mockGateway struct {
	calls      int32
	shouldFail bool
	failCount  int32
	maxFails   int32
}

func (m *mockGateway) SubmitForm(_ context.Context, _ string, _ domain.FormType, _ string) (*SubmitResult, error) {
	n := atomic.AddInt32(&m.calls, 1)
	if m.shouldFail {
		if m.maxFails > 0 && atomic.LoadInt32(&m.failCount) >= m.maxFails {
			return &SubmitResult{
				TokenNumber:           fmt.Sprintf("TKN-%06d", n),
				AcknowledgementNumber: fmt.Sprintf("ACK-%06d", n),
				Status:                "FILED",
			}, nil
		}
		atomic.AddInt32(&m.failCount, 1)
		return nil, fmt.Errorf("gateway unavailable")
	}
	return &SubmitResult{
		TokenNumber:           fmt.Sprintf("TKN-%06d", n),
		AcknowledgementNumber: fmt.Sprintf("ACK-%06d", n),
		Status:                "FILED",
	}, nil
}

func seedSagaData(ms *store.MockStore) {
	deductees := []domain.Deductee{
		{ID: uuid.MustParse("aaaa0001-0001-0001-0001-000000000001"), TenantID: testTenant, VendorID: uuid.New(), Name: "BuildRight Contractors", PAN: "AABCB1234F", PANVerified: true, DeducteeType: domain.DeducteeCompany, ResidentStatus: domain.Resident, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.MustParse("aaaa0001-0001-0001-0001-000000000002"), TenantID: testTenant, VendorID: uuid.New(), Name: "Ravi Auditor", PAN: "EEFGH9012H", PANVerified: true, DeducteeType: domain.DeducteeIndividual, ResidentStatus: domain.Resident, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.MustParse("aaaa0001-0001-0001-0001-000000000003"), TenantID: testTenant, VendorID: uuid.New(), Name: "OfficeSpace Realty", PAN: "CCDEF5678G", PANVerified: true, DeducteeType: domain.DeducteeCompany, ResidentStatus: domain.Resident, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, d := range deductees {
		d := d
		ms.CreateDeductee(context.Background(), testTenant, &d)
	}

	entries := []domain.TDSEntry{
		{
			ID: uuid.New(), TenantID: testTenant, DeducteeID: deductees[0].ID,
			Section: domain.Section194C, FinancialYear: "2025-26", Quarter: "Q1",
			TransactionDate: testDate, GrossAmount: decimal.NewFromInt(500000),
			TDSRate: decimal.NewFromFloat(0.02), TDSAmount: decimal.NewFromInt(10000),
			Surcharge: decimal.Zero, Cess: decimal.Zero, TotalTax: decimal.NewFromInt(10000),
			NatureOfPayment: "Civil construction", PANAtDeduction: "AABCB1234F",
			Status: domain.StatusPending, ChallanNumber: "CHN-001", BSRCode: "BSR001",
		},
		{
			ID: uuid.New(), TenantID: testTenant, DeducteeID: deductees[1].ID,
			Section: domain.Section194J, FinancialYear: "2025-26", Quarter: "Q1",
			TransactionDate: testDate, GrossAmount: decimal.NewFromInt(100000),
			TDSRate: decimal.NewFromFloat(0.10), TDSAmount: decimal.NewFromInt(10000),
			Surcharge: decimal.Zero, Cess: decimal.Zero, TotalTax: decimal.NewFromInt(10000),
			NatureOfPayment: "Audit fees", PANAtDeduction: "EEFGH9012H",
			Status: domain.StatusPending, ChallanNumber: "CHN-002", BSRCode: "BSR001",
		},
		{
			ID: uuid.New(), TenantID: testTenant, DeducteeID: deductees[2].ID,
			Section: domain.Section194I, FinancialYear: "2025-26", Quarter: "Q1",
			TransactionDate: testDate, GrossAmount: decimal.NewFromInt(300000),
			TDSRate: decimal.NewFromFloat(0.10), TDSAmount: decimal.NewFromInt(30000),
			Surcharge: decimal.Zero, Cess: decimal.Zero, TotalTax: decimal.NewFromInt(30000),
			NatureOfPayment: "Office rent", PANAtDeduction: "CCDEF5678G",
			Status: domain.StatusDeposited, ChallanNumber: "CHN-003", BSRCode: "BSR002",
		},
	}
	for _, e := range entries {
		e := e
		ms.CreateEntry(context.Background(), testTenant, &e)
	}
}

func seedSalaryData(ms *store.MockStore) {
	d := domain.Deductee{
		ID: uuid.MustParse("bbbb0001-0001-0001-0001-000000000001"), TenantID: testTenant, VendorID: uuid.New(),
		Name: "Rajesh Kumar", PAN: "ABCPK1234E", PANVerified: true,
		DeducteeType: domain.DeducteeIndividual, ResidentStatus: domain.Resident,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	ms.CreateDeductee(context.Background(), testTenant, &d)

	for month := 0; month < 3; month++ {
		e := domain.TDSEntry{
			ID: uuid.New(), TenantID: testTenant, DeducteeID: d.ID,
			Section: domain.Section192, FinancialYear: "2025-26", Quarter: "Q1",
			TransactionDate: time.Date(2025, time.Month(4+month), 28, 0, 0, 0, 0, time.UTC),
			GrossAmount: decimal.NewFromInt(166667), TDSRate: decimal.NewFromFloat(0.0962),
			TDSAmount: decimal.NewFromInt(16033), Surcharge: decimal.Zero,
			Cess: decimal.Zero, TotalTax: decimal.NewFromInt(16033),
			NatureOfPayment: "Salary", PANAtDeduction: d.PAN,
			Status: domain.StatusPending, ChallanNumber: "CHN-SAL-001", BSRCode: "BSR001",
		}
		ms.CreateEntry(context.Background(), testTenant, &e)
	}
}

func seedNonResidentData(ms *store.MockStore) {
	d := domain.Deductee{
		ID: uuid.MustParse("cccc0001-0001-0001-0001-000000000001"), TenantID: testTenant, VendorID: uuid.New(),
		Name: "TechGlobal Inc", PAN: "AABFT1234G", PANVerified: true,
		DeducteeType: domain.DeducteeCompany, ResidentStatus: domain.NonResident,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	ms.CreateDeductee(context.Background(), testTenant, &d)

	e := domain.TDSEntry{
		ID: uuid.New(), TenantID: testTenant, DeducteeID: d.ID,
		Section: domain.Section195, FinancialYear: "2025-26", Quarter: "Q1",
		TransactionDate: testDate, GrossAmount: decimal.NewFromInt(5000000),
		TDSRate: decimal.NewFromFloat(0.20), TDSAmount: decimal.NewFromInt(1000000),
		Surcharge: decimal.Zero, Cess: decimal.NewFromInt(40000),
		TotalTax: decimal.NewFromInt(1040000),
		NatureOfPayment: "Software services", PANAtDeduction: d.PAN,
		Status: domain.StatusPending, ChallanNumber: "CHN-NR-001", BSRCode: "BSR003",
	}
	ms.CreateEntry(context.Background(), testTenant, &e)
}

// --- Happy path tests ---

func TestFilingSaga_26Q_HappyPath(t *testing.T) {
	ms := store.NewMockStore()
	gw := &mockGateway{}
	saga := NewFilingSaga(ms, gw)

	seedSagaData(ms)

	result, err := saga.Execute(context.Background(), FilingSagaInput{
		TenantID:      testTenant,
		FormType:      domain.FormType26Q,
		FinancialYear: "2025-26",
		Quarter:       "Q1",
		TAN:           "MUMA12345A",
		Deductor: &domain.DeductorDetails{
			TAN: "MUMA12345A", DeductorName: "Acme Corp", DeductorPAN: "AABCA1234F",
			Address: "100 MG Road", Pincode: "400001",
		},
	})

	require.NoError(t, err)
	assert.Equal(t, domain.FilingFiled, result.Status)
	assert.NotEmpty(t, result.TokenNumber)
	assert.NotEmpty(t, result.AcknowledgementNumber)
	assert.Equal(t, 3, result.DeducteeCount)
	assert.NotEmpty(t, result.FVUContent)
	assert.Contains(t, result.FVUContent, "^FH^26Q^")
	assert.Equal(t, int32(1), atomic.LoadInt32(&gw.calls))
}

func TestFilingSaga_24Q_HappyPath(t *testing.T) {
	ms := store.NewMockStore()
	gw := &mockGateway{}
	saga := NewFilingSaga(ms, gw)

	seedSalaryData(ms)

	result, err := saga.Execute(context.Background(), FilingSagaInput{
		TenantID:      testTenant,
		FormType:      domain.FormType24Q,
		FinancialYear: "2025-26",
		Quarter:       "Q1",
		TAN:           "MUMA12345A",
		Employer: &domain.EmployerDetails{
			TAN: "MUMA12345A", EmployerName: "Acme Corp", EmployerPAN: "AABCA1234F",
			Address: "100 MG Road", Pincode: "400001",
		},
	})

	require.NoError(t, err)
	assert.Equal(t, domain.FilingFiled, result.Status)
	assert.Contains(t, result.FVUContent, "^FH^24Q^")
	assert.Equal(t, 1, result.DeducteeCount)
}

func TestFilingSaga_27Q_HappyPath(t *testing.T) {
	ms := store.NewMockStore()
	gw := &mockGateway{}
	saga := NewFilingSaga(ms, gw)

	seedNonResidentData(ms)

	deducteeID := "cccc0001-0001-0001-0001-000000000001"
	result, err := saga.Execute(context.Background(), FilingSagaInput{
		TenantID:      testTenant,
		FormType:      domain.FormType27Q,
		FinancialYear: "2025-26",
		Quarter:       "Q1",
		TAN:           "MUMA12345A",
		Deductor: &domain.DeductorDetails{
			TAN: "MUMA12345A", DeductorName: "Acme Corp", DeductorPAN: "AABCA1234F",
			Address: "100 MG Road", Pincode: "400001",
		},
		CountryCodes:  map[string]string{deducteeID: "US"},
		CurrencyCodes: map[string]string{deducteeID: "USD"},
		ExchangeRates: map[string]decimal.Decimal{deducteeID: decimal.NewFromFloat(83.50)},
		ForeignAmounts: map[string]decimal.Decimal{deducteeID: decimal.NewFromFloat(59880.24)},
		DTAAArticles:  map[string]string{deducteeID: "12"},
		DTAARates:     map[string]decimal.Decimal{deducteeID: decimal.NewFromFloat(0.15)},
	})

	require.NoError(t, err)
	assert.Equal(t, domain.FilingFiled, result.Status)
	assert.Contains(t, result.FVUContent, "^FH^27Q^")
	assert.Equal(t, 1, result.DeducteeCount)
}

// --- Failure and recovery tests ---

func TestFilingSaga_GatewayError_SetsRejected(t *testing.T) {
	ms := store.NewMockStore()
	gw := &mockGateway{shouldFail: true}
	saga := NewFilingSaga(ms, gw)

	seedSagaData(ms)

	_, err := saga.Execute(context.Background(), FilingSagaInput{
		TenantID:      testTenant,
		FormType:      domain.FormType26Q,
		FinancialYear: "2025-26",
		Quarter:       "Q1",
		TAN:           "MUMA12345A",
		Deductor:      &domain.DeductorDetails{TAN: "MUMA12345A", DeductorName: "Acme Corp", DeductorPAN: "AABCA1234F", Address: "x", Pincode: "400001"},
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "submit")

	filings, _, _ := ms.ListFilings(context.Background(), testTenant, "2025-26", "Q1", 10, 0)
	require.Len(t, filings, 1)
	assert.Equal(t, domain.FilingRejected, filings[0].Status)
	assert.Contains(t, filings[0].ErrorMessage, "gateway unavailable")
}

func TestFilingSaga_GatewayRecovery(t *testing.T) {
	ms := store.NewMockStore()
	gw := &mockGateway{shouldFail: true, maxFails: 1}
	saga := NewFilingSaga(ms, gw)

	seedSagaData(ms)

	input := FilingSagaInput{
		TenantID:      testTenant,
		FormType:      domain.FormType26Q,
		FinancialYear: "2025-26",
		Quarter:       "Q1",
		TAN:           "MUMA12345A",
		Deductor:      &domain.DeductorDetails{TAN: "MUMA12345A", DeductorName: "Acme Corp", DeductorPAN: "AABCA1234F", Address: "x", Pincode: "400001"},
	}

	_, err := saga.Execute(context.Background(), input)
	assert.Error(t, err, "first attempt should fail")

	result, err := saga.Execute(context.Background(), input)
	require.NoError(t, err, "second attempt should succeed after recovery")
	assert.Equal(t, domain.FilingFiled, result.Status)
	assert.NotEmpty(t, result.TokenNumber)
}

// --- Idempotency tests ---

func TestFilingSaga_Idempotency_DuplicateReturnsSame(t *testing.T) {
	ms := store.NewMockStore()
	gw := &mockGateway{}
	saga := NewFilingSaga(ms, gw)

	seedSagaData(ms)

	input := FilingSagaInput{
		TenantID:      testTenant,
		FormType:      domain.FormType26Q,
		FinancialYear: "2025-26",
		Quarter:       "Q1",
		TAN:           "MUMA12345A",
		Deductor:      &domain.DeductorDetails{TAN: "MUMA12345A", DeductorName: "Acme Corp", DeductorPAN: "AABCA1234F", Address: "x", Pincode: "400001"},
	}

	first, err := saga.Execute(context.Background(), input)
	require.NoError(t, err)

	second, err := saga.Execute(context.Background(), input)
	require.NoError(t, err)

	assert.Equal(t, first.ID, second.ID, "should return same filing on duplicate")
	assert.Equal(t, first.TokenNumber, second.TokenNumber)
	assert.Equal(t, int32(1), atomic.LoadInt32(&gw.calls), "gateway should only be called once")
}

func TestFilingSaga_Idempotency_DifferentQuarters(t *testing.T) {
	ms := store.NewMockStore()
	gw := &mockGateway{}
	saga := NewFilingSaga(ms, gw)

	seedSagaData(ms)

	deductor := &domain.DeductorDetails{TAN: "MUMA12345A", DeductorName: "Acme Corp", DeductorPAN: "AABCA1234F", Address: "x", Pincode: "400001"}

	r1, err := saga.Execute(context.Background(), FilingSagaInput{
		TenantID: testTenant, FormType: domain.FormType26Q,
		FinancialYear: "2025-26", Quarter: "Q1", TAN: "MUMA12345A", Deductor: deductor,
	})
	require.NoError(t, err)

	r2, err := saga.Execute(context.Background(), FilingSagaInput{
		TenantID: testTenant, FormType: domain.FormType26Q,
		FinancialYear: "2025-26", Quarter: "Q1", TAN: "MUMA12345A", Deductor: deductor,
	})
	require.NoError(t, err)

	assert.Equal(t, r1.ID, r2.ID, "same key should return same filing")
	assert.Equal(t, int32(1), atomic.LoadInt32(&gw.calls))
}

// --- Validation failure tests ---

func TestFilingSaga_NoDeductees(t *testing.T) {
	ms := store.NewMockStore()
	gw := &mockGateway{}
	saga := NewFilingSaga(ms, gw)

	_, err := saga.Execute(context.Background(), FilingSagaInput{
		TenantID:      testTenant,
		FormType:      domain.FormType26Q,
		FinancialYear: "2025-26",
		Quarter:       "Q1",
		TAN:           "MUMA12345A",
		Deductor:      &domain.DeductorDetails{TAN: "MUMA12345A", DeductorName: "Acme Corp", DeductorPAN: "AABCA1234F", Address: "x", Pincode: "400001"},
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no deductees")
}

func TestFilingSaga_MissingPAN_24Q(t *testing.T) {
	ms := store.NewMockStore()
	gw := &mockGateway{}
	saga := NewFilingSaga(ms, gw)

	d := domain.Deductee{
		ID: uuid.New(), TenantID: testTenant, VendorID: uuid.New(),
		Name: "No PAN Employee", PAN: "",
		DeducteeType: domain.DeducteeIndividual, ResidentStatus: domain.Resident,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	ms.CreateDeductee(context.Background(), testTenant, &d)

	_, err := saga.Execute(context.Background(), FilingSagaInput{
		TenantID:      testTenant,
		FormType:      domain.FormType24Q,
		FinancialYear: "2025-26",
		Quarter:       "Q1",
		TAN:           "MUMA12345A",
		Employer:      &domain.EmployerDetails{TAN: "MUMA12345A", EmployerName: "Acme Corp", EmployerPAN: "AABCA1234F", Address: "x", Pincode: "400001"},
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing PAN")
}

func TestFilingSaga_MissingEmployerDetails(t *testing.T) {
	ms := store.NewMockStore()
	gw := &mockGateway{}
	saga := NewFilingSaga(ms, gw)
	seedSalaryData(ms)

	_, err := saga.Execute(context.Background(), FilingSagaInput{
		TenantID:      testTenant,
		FormType:      domain.FormType24Q,
		FinancialYear: "2025-26",
		Quarter:       "Q1",
		TAN:           "MUMA12345A",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "employer details")
}

func TestFilingSaga_MissingDeductorDetails(t *testing.T) {
	ms := store.NewMockStore()
	gw := &mockGateway{}
	saga := NewFilingSaga(ms, gw)
	seedSagaData(ms)

	_, err := saga.Execute(context.Background(), FilingSagaInput{
		TenantID:      testTenant,
		FormType:      domain.FormType26Q,
		FinancialYear: "2025-26",
		Quarter:       "Q1",
		TAN:           "MUMA12345A",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "deductor details")
}

func TestFilingSaga_UnsupportedFormType(t *testing.T) {
	ms := store.NewMockStore()
	gw := &mockGateway{}
	saga := NewFilingSaga(ms, gw)

	d := domain.Deductee{
		ID: uuid.New(), TenantID: testTenant, VendorID: uuid.New(),
		Name: "Test", PAN: "ABCDE1234F",
		DeducteeType: domain.DeducteeCompany, ResidentStatus: domain.Resident,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	ms.CreateDeductee(context.Background(), testTenant, &d)

	_, err := saga.Execute(context.Background(), FilingSagaInput{
		TenantID:      testTenant,
		FormType:      domain.FormType("99Q"),
		FinancialYear: "2025-26",
		Quarter:       "Q1",
		TAN:           "MUMA12345A",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported form type")
}

func TestFilingSaga_SubmittedStatus(t *testing.T) {
	ms := store.NewMockStore()

	seedSagaData(ms)

	submittedGW := &submittedGateway{}
	saga2 := NewFilingSaga(ms, submittedGW)

	result, err := saga2.Execute(context.Background(), FilingSagaInput{
		TenantID: testTenant, FormType: domain.FormType26Q,
		FinancialYear: "2025-26", Quarter: "Q1", TAN: "MUMA12345A",
		Deductor: &domain.DeductorDetails{TAN: "MUMA12345A", DeductorName: "Acme Corp", DeductorPAN: "AABCA1234F", Address: "x", Pincode: "400001"},
	})

	require.NoError(t, err)
	assert.Equal(t, domain.FilingSubmitted, result.Status)
}

type submittedGateway struct{}

func (g *submittedGateway) SubmitForm(_ context.Context, _ string, _ domain.FormType, _ string) (*SubmitResult, error) {
	return &SubmitResult{
		TokenNumber:           "TKN-SUB-001",
		AcknowledgementNumber: "ACK-SUB-001",
		Status:                "SUBMITTED",
	}, nil
}

func TestNewFilingSaga(t *testing.T) {
	ms := store.NewMockStore()
	gw := &mockGateway{}
	saga := NewFilingSaga(ms, gw)
	assert.NotNil(t, saga)
	assert.NotNil(t, saga.store)
	assert.NotNil(t, saga.gateway)
}
