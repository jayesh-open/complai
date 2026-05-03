package store

import (
	"context"
	"sync"

	"github.com/complai/complai/services/go/gstr9-service/internal/domain"
	"github.com/google/uuid"
)

type MockStore struct {
	mu        sync.RWMutex
	filings   map[uuid.UUID]*domain.GSTR9Filing
	tableData map[uuid.UUID][]domain.GSTR9TableData
	auditLogs map[uuid.UUID][]domain.GSTR9AuditLog
}

func NewMockStore() *MockStore {
	return &MockStore{
		filings:   make(map[uuid.UUID]*domain.GSTR9Filing),
		tableData: make(map[uuid.UUID][]domain.GSTR9TableData),
		auditLogs: make(map[uuid.UUID][]domain.GSTR9AuditLog),
	}
}

func (m *MockStore) CreateFiling(_ context.Context, _ uuid.UUID, f *domain.GSTR9Filing) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, existing := range m.filings {
		if existing.TenantID == f.TenantID && existing.GSTIN == f.GSTIN && existing.FinancialYear == f.FinancialYear {
			return domain.ErrDuplicateFiling
		}
	}
	m.filings[f.ID] = f
	return nil
}

func (m *MockStore) GetFiling(_ context.Context, tenantID, id uuid.UUID) (*domain.GSTR9Filing, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f, ok := m.filings[id]
	if !ok || f.TenantID != tenantID {
		return nil, domain.ErrFilingNotFound
	}
	return f, nil
}

func (m *MockStore) UpdateFilingStatus(_ context.Context, tenantID, id uuid.UUID, status domain.FilingStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	f, ok := m.filings[id]
	if !ok || f.TenantID != tenantID {
		return domain.ErrFilingNotFound
	}
	f.Status = status
	return nil
}

func (m *MockStore) ListFilings(_ context.Context, tenantID uuid.UUID, gstin, fy string, limit, offset int) ([]domain.GSTR9Filing, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var all []domain.GSTR9Filing
	for _, f := range m.filings {
		if f.TenantID != tenantID {
			continue
		}
		if gstin != "" && f.GSTIN != gstin {
			continue
		}
		if fy != "" && f.FinancialYear != fy {
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

func (m *MockStore) CreateTableData(_ context.Context, _ uuid.UUID, td *domain.GSTR9TableData) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tableData[td.FilingID] = append(m.tableData[td.FilingID], *td)
	return nil
}

func (m *MockStore) ListTableData(_ context.Context, tenantID uuid.UUID, filingID uuid.UUID) ([]domain.GSTR9TableData, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []domain.GSTR9TableData
	for _, td := range m.tableData[filingID] {
		if td.TenantID == tenantID {
			out = append(out, td)
		}
	}
	return out, nil
}

func (m *MockStore) DeleteTableData(_ context.Context, _ uuid.UUID, filingID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.tableData, filingID)
	return nil
}

func (m *MockStore) CreateAuditLog(_ context.Context, _ uuid.UUID, log *domain.GSTR9AuditLog) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.auditLogs[log.FilingID] = append(m.auditLogs[log.FilingID], *log)
	return nil
}

func (m *MockStore) ListAuditLogs(_ context.Context, tenantID uuid.UUID, filingID uuid.UUID) ([]domain.GSTR9AuditLog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []domain.GSTR9AuditLog
	for _, l := range m.auditLogs[filingID] {
		if l.TenantID == tenantID {
			out = append(out, l)
		}
	}
	return out, nil
}

var _ Repository = (*MockStore)(nil)
