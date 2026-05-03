package store

import (
	"context"
	"fmt"
	"sync"

	"github.com/complai/complai/services/go/itr-service/internal/domain"
	"github.com/google/uuid"
)

type MockStore struct {
	mu               sync.RWMutex
	taxpayers        map[uuid.UUID]*domain.Taxpayer
	filings          map[uuid.UUID]*domain.ITRFiling
	incomeEntries    map[uuid.UUID][]domain.IncomeEntry
	deductions       map[uuid.UUID][]domain.Deduction
	taxComputations  map[uuid.UUID]*domain.TaxComputation
	tdsCredits       map[uuid.UUID][]domain.TDSCredit
	aisReconciliations map[uuid.UUID][]domain.AISReconciliation
}

func NewMockStore() *MockStore {
	return &MockStore{
		taxpayers:          make(map[uuid.UUID]*domain.Taxpayer),
		filings:            make(map[uuid.UUID]*domain.ITRFiling),
		incomeEntries:      make(map[uuid.UUID][]domain.IncomeEntry),
		deductions:         make(map[uuid.UUID][]domain.Deduction),
		taxComputations:    make(map[uuid.UUID]*domain.TaxComputation),
		tdsCredits:         make(map[uuid.UUID][]domain.TDSCredit),
		aisReconciliations: make(map[uuid.UUID][]domain.AISReconciliation),
	}
}

func (m *MockStore) CreateTaxpayer(_ context.Context, _ uuid.UUID, t *domain.Taxpayer) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.taxpayers[t.ID] = t
	return nil
}

func (m *MockStore) GetTaxpayer(_ context.Context, tenantID, id uuid.UUID) (*domain.Taxpayer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.taxpayers[id]
	if !ok || t.TenantID != tenantID {
		return nil, fmt.Errorf("taxpayer not found")
	}
	return t, nil
}

func (m *MockStore) GetTaxpayerByPAN(_ context.Context, tenantID uuid.UUID, pan string) (*domain.Taxpayer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, t := range m.taxpayers {
		if t.PAN == pan && t.TenantID == tenantID {
			return t, nil
		}
	}
	return nil, fmt.Errorf("taxpayer not found")
}

func (m *MockStore) ListTaxpayers(_ context.Context, tenantID uuid.UUID, limit, offset int) ([]domain.Taxpayer, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var all []domain.Taxpayer
	for _, t := range m.taxpayers {
		if t.TenantID == tenantID {
			all = append(all, *t)
		}
	}
	total := len(all)
	if offset >= total {
		return nil, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return all[offset:end], total, nil
}

func (m *MockStore) CreateFiling(_ context.Context, _ uuid.UUID, f *domain.ITRFiling) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, existing := range m.filings {
		if existing.IdempotencyKey == f.IdempotencyKey {
			return fmt.Errorf("filing with idempotency key %s already exists", f.IdempotencyKey)
		}
	}
	m.filings[f.ID] = f
	return nil
}

func (m *MockStore) GetFiling(_ context.Context, tenantID, id uuid.UUID) (*domain.ITRFiling, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f, ok := m.filings[id]
	if !ok || f.TenantID != tenantID {
		return nil, fmt.Errorf("filing not found")
	}
	return f, nil
}

func (m *MockStore) GetFilingByIdempotencyKey(_ context.Context, tenantID uuid.UUID, key string) (*domain.ITRFiling, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, f := range m.filings {
		if f.IdempotencyKey == key && f.TenantID == tenantID {
			return f, nil
		}
	}
	return nil, fmt.Errorf("filing not found")
}

func (m *MockStore) UpdateFilingStatus(_ context.Context, tenantID, id uuid.UUID, status domain.FilingStatus, arn, ackNumber, errMsg string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	f, ok := m.filings[id]
	if !ok || f.TenantID != tenantID {
		return fmt.Errorf("filing not found")
	}
	f.Status = status
	f.ARN = arn
	f.AcknowledgementNumber = ackNumber
	f.ErrorMessage = errMsg
	return nil
}

func (m *MockStore) ListFilings(_ context.Context, tenantID uuid.UUID, taxYear string, limit, offset int) ([]domain.ITRFiling, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var all []domain.ITRFiling
	for _, f := range m.filings {
		if f.TenantID != tenantID {
			continue
		}
		if taxYear != "" && f.TaxYear != taxYear {
			continue
		}
		all = append(all, *f)
	}
	total := len(all)
	if offset >= total {
		return nil, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return all[offset:end], total, nil
}

func (m *MockStore) CreateIncomeEntry(_ context.Context, _ uuid.UUID, e *domain.IncomeEntry) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.incomeEntries[e.FilingID] = append(m.incomeEntries[e.FilingID], *e)
	return nil
}

func (m *MockStore) ListIncomeEntries(_ context.Context, _ uuid.UUID, filingID uuid.UUID) ([]domain.IncomeEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.incomeEntries[filingID], nil
}

func (m *MockStore) CreateDeduction(_ context.Context, _ uuid.UUID, d *domain.Deduction) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.deductions[d.FilingID] = append(m.deductions[d.FilingID], *d)
	return nil
}

func (m *MockStore) ListDeductions(_ context.Context, _ uuid.UUID, filingID uuid.UUID) ([]domain.Deduction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.deductions[filingID], nil
}

func (m *MockStore) SaveTaxComputation(_ context.Context, _ uuid.UUID, tc *domain.TaxComputation) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.taxComputations[tc.FilingID] = tc
	return nil
}

func (m *MockStore) GetTaxComputation(_ context.Context, _ uuid.UUID, filingID uuid.UUID) (*domain.TaxComputation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	tc, ok := m.taxComputations[filingID]
	if !ok {
		return nil, fmt.Errorf("tax computation not found")
	}
	return tc, nil
}

func (m *MockStore) CreateTDSCredit(_ context.Context, _ uuid.UUID, c *domain.TDSCredit) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tdsCredits[c.FilingID] = append(m.tdsCredits[c.FilingID], *c)
	return nil
}

func (m *MockStore) ListTDSCredits(_ context.Context, _ uuid.UUID, filingID uuid.UUID) ([]domain.TDSCredit, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tdsCredits[filingID], nil
}

func (m *MockStore) CreateAISReconciliation(_ context.Context, _ uuid.UUID, r *domain.AISReconciliation) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.aisReconciliations[r.FilingID] = append(m.aisReconciliations[r.FilingID], *r)
	return nil
}

func (m *MockStore) ListAISReconciliations(_ context.Context, _ uuid.UUID, filingID uuid.UUID) ([]domain.AISReconciliation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.aisReconciliations[filingID], nil
}
