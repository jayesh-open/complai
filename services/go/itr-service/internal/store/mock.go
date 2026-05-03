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
	bulkBatches      map[uuid.UUID]*domain.BulkFilingBatch
	bulkEmployees    map[uuid.UUID]*domain.BulkFilingEmployee
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
		bulkBatches:        make(map[uuid.UUID]*domain.BulkFilingBatch),
		bulkEmployees:      make(map[uuid.UUID]*domain.BulkFilingEmployee),
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

func (m *MockStore) CreateBulkBatch(_ context.Context, _ uuid.UUID, b *domain.BulkFilingBatch) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.bulkBatches[b.ID] = b
	return nil
}

func (m *MockStore) GetBulkBatch(_ context.Context, tenantID, id uuid.UUID) (*domain.BulkFilingBatch, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	b, ok := m.bulkBatches[id]
	if !ok || b.TenantID != tenantID {
		return nil, fmt.Errorf("batch not found")
	}
	return b, nil
}

func (m *MockStore) UpdateBulkBatchStatus(_ context.Context, tenantID, id uuid.UUID, status domain.BulkBatchStatus, processed, ready, mismatches int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	b, ok := m.bulkBatches[id]
	if !ok || b.TenantID != tenantID {
		return fmt.Errorf("batch not found")
	}
	b.Status = status
	b.Processed = processed
	b.Ready = ready
	b.WithMismatches = mismatches
	return nil
}

func (m *MockStore) ListBulkBatches(_ context.Context, tenantID uuid.UUID, limit, offset int) ([]domain.BulkFilingBatch, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var all []domain.BulkFilingBatch
	for _, b := range m.bulkBatches {
		if b.TenantID == tenantID {
			all = append(all, *b)
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

func (m *MockStore) CreateBulkEmployee(_ context.Context, _ uuid.UUID, e *domain.BulkFilingEmployee) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.bulkEmployees[e.ID] = e
	return nil
}

func (m *MockStore) GetBulkEmployee(_ context.Context, tenantID, id uuid.UUID) (*domain.BulkFilingEmployee, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.bulkEmployees[id]
	if !ok || e.TenantID != tenantID {
		return nil, fmt.Errorf("employee not found")
	}
	return e, nil
}

func (m *MockStore) ListBulkEmployees(_ context.Context, tenantID uuid.UUID, batchID uuid.UUID, limit, offset int) ([]domain.BulkFilingEmployee, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var all []domain.BulkFilingEmployee
	for _, e := range m.bulkEmployees {
		if e.TenantID == tenantID && e.BatchID == batchID {
			all = append(all, *e)
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

func (m *MockStore) UpdateBulkEmployeeStatus(_ context.Context, tenantID, id uuid.UUID, status domain.EmployeeFilingStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.bulkEmployees[id]
	if !ok || e.TenantID != tenantID {
		return fmt.Errorf("employee not found")
	}
	e.Status = status
	return nil
}
