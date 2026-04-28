package store

import (
	"context"
	"fmt"
	"sync"

	"github.com/complai/complai/services/go/tds-service/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type MockStore struct {
	mu         sync.RWMutex
	deductees  map[uuid.UUID]*domain.Deductee
	entries    map[uuid.UUID]*domain.TDSEntry
	aggregates map[string]*domain.TDSAggregate
}

func NewMockStore() *MockStore {
	return &MockStore{
		deductees:  make(map[uuid.UUID]*domain.Deductee),
		entries:    make(map[uuid.UUID]*domain.TDSEntry),
		aggregates: make(map[string]*domain.TDSAggregate),
	}
}

func aggKey(deducteeID uuid.UUID, section domain.Section, fy string) string {
	return fmt.Sprintf("%s:%s:%s", deducteeID, section, fy)
}

func (m *MockStore) CreateDeductee(_ context.Context, _ uuid.UUID, d *domain.Deductee) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.deductees[d.ID] = d
	return nil
}

func (m *MockStore) GetDeductee(_ context.Context, tenantID, id uuid.UUID) (*domain.Deductee, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	d, ok := m.deductees[id]
	if !ok || d.TenantID != tenantID {
		return nil, fmt.Errorf("deductee not found")
	}
	return d, nil
}

func (m *MockStore) GetDeducteeByVendor(_ context.Context, tenantID, vendorID uuid.UUID) (*domain.Deductee, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, d := range m.deductees {
		if d.VendorID == vendorID && d.TenantID == tenantID {
			return d, nil
		}
	}
	return nil, fmt.Errorf("deductee not found for vendor")
}

func (m *MockStore) ListDeductees(_ context.Context, tenantID uuid.UUID, limit, offset int) ([]domain.Deductee, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var all []domain.Deductee
	for _, d := range m.deductees {
		if d.TenantID == tenantID {
			all = append(all, *d)
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

func (m *MockStore) UpsertDeductee(_ context.Context, _ uuid.UUID, d *domain.Deductee) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, existing := range m.deductees {
		if existing.VendorID == d.VendorID && existing.TenantID == d.TenantID {
			existing.Name = d.Name
			existing.PAN = d.PAN
			existing.PANStatus = d.PANStatus
			existing.DeducteeType = d.DeducteeType
			return nil
		}
	}
	m.deductees[d.ID] = d
	return nil
}

func (m *MockStore) CreateEntry(_ context.Context, _ uuid.UUID, e *domain.TDSEntry) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[e.ID] = e
	return nil
}

func (m *MockStore) GetEntry(_ context.Context, tenantID, id uuid.UUID) (*domain.TDSEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.entries[id]
	if !ok || e.TenantID != tenantID {
		return nil, fmt.Errorf("entry not found")
	}
	return e, nil
}

func (m *MockStore) ListEntries(_ context.Context, tenantID uuid.UUID, fy, quarter string, limit, offset int) ([]domain.TDSEntry, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var all []domain.TDSEntry
	for _, e := range m.entries {
		if e.TenantID != tenantID {
			continue
		}
		if fy != "" && e.FinancialYear != fy {
			continue
		}
		if quarter != "" && e.Quarter != quarter {
			continue
		}
		all = append(all, *e)
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

func (m *MockStore) GetAggregate(_ context.Context, _ uuid.UUID, deducteeID uuid.UUID, section domain.Section, fy string) (*domain.TDSAggregate, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	a, ok := m.aggregates[aggKey(deducteeID, section, fy)]
	if !ok {
		return &domain.TDSAggregate{TotalPaid: decimal.Zero, TotalTDS: decimal.Zero}, nil
	}
	return a, nil
}

func (m *MockStore) UpsertAggregate(_ context.Context, _ uuid.UUID, agg *domain.TDSAggregate) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.aggregates[aggKey(agg.DeducteeID, agg.Section, agg.FinancialYear)] = agg
	return nil
}

func (m *MockStore) GetSummary(_ context.Context, tenantID uuid.UUID, fy string) (*domain.TDSSummary, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sum := &domain.TDSSummary{
		EntriesBySection: make(map[domain.Section]int),
		EntriesByStatus:  make(map[domain.EntryStatus]int),
		TotalTDSDeducted: decimal.Zero,
		TotalTDSDeposited: decimal.Zero,
		PendingDeposit:    decimal.Zero,
	}
	for _, d := range m.deductees {
		if d.TenantID == tenantID {
			sum.TotalDeductees++
		}
	}
	for _, e := range m.entries {
		if e.TenantID != tenantID || e.FinancialYear != fy {
			continue
		}
		sum.TotalEntries++
		sum.TotalTDSDeducted = sum.TotalTDSDeducted.Add(e.TotalTax)
		sum.EntriesBySection[e.Section]++
		sum.EntriesByStatus[e.Status]++
		if e.Status == domain.StatusDeposited {
			sum.TotalTDSDeposited = sum.TotalTDSDeposited.Add(e.TotalTax)
		}
	}
	sum.PendingDeposit = sum.TotalTDSDeducted.Sub(sum.TotalTDSDeposited)
	return sum, nil
}
