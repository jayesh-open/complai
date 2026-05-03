package store

import (
	"context"
	"sync"
	"time"

	"github.com/complai/complai/services/go/gstr9-service/internal/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type MockStore struct {
	mu          sync.RWMutex
	filings     map[uuid.UUID]*domain.GSTR9Filing
	tableData   map[uuid.UUID][]domain.GSTR9TableData
	auditLogs   map[uuid.UUID][]domain.GSTR9AuditLog
	gstr9cFiles map[uuid.UUID]*domain.GSTR9CFiling
	mismatches  map[uuid.UUID]*domain.GSTR9CMismatch

	ErrOnListFilings error
}

func NewMockStore() *MockStore {
	return &MockStore{
		filings:     make(map[uuid.UUID]*domain.GSTR9Filing),
		tableData:   make(map[uuid.UUID][]domain.GSTR9TableData),
		auditLogs:   make(map[uuid.UUID][]domain.GSTR9AuditLog),
		gstr9cFiles: make(map[uuid.UUID]*domain.GSTR9CFiling),
		mismatches:  make(map[uuid.UUID]*domain.GSTR9CMismatch),
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
	if m.ErrOnListFilings != nil {
		return nil, 0, m.ErrOnListFilings
	}
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

func (m *MockStore) CreateGSTR9CFiling(_ context.Context, _ uuid.UUID, f *domain.GSTR9CFiling) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, existing := range m.gstr9cFiles {
		if existing.TenantID == f.TenantID && existing.GSTR9FilingID == f.GSTR9FilingID {
			return domain.ErrGSTR9CDuplicate
		}
	}
	m.gstr9cFiles[f.ID] = f
	return nil
}

func (m *MockStore) GetGSTR9CFiling(_ context.Context, tenantID, id uuid.UUID) (*domain.GSTR9CFiling, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f, ok := m.gstr9cFiles[id]
	if !ok || f.TenantID != tenantID {
		return nil, domain.ErrGSTR9CNotFound
	}
	return f, nil
}

func (m *MockStore) GetGSTR9CFilingByGSTR9ID(_ context.Context, tenantID, gstr9FilingID uuid.UUID) (*domain.GSTR9CFiling, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, f := range m.gstr9cFiles {
		if f.TenantID == tenantID && f.GSTR9FilingID == gstr9FilingID {
			return f, nil
		}
	}
	return nil, domain.ErrGSTR9CNotFound
}

func (m *MockStore) UpdateGSTR9CStatus(_ context.Context, tenantID, id uuid.UUID, status domain.GSTR9CStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	f, ok := m.gstr9cFiles[id]
	if !ok || f.TenantID != tenantID {
		return domain.ErrGSTR9CNotFound
	}
	f.Status = status
	f.UpdatedAt = time.Now()
	return nil
}

func (m *MockStore) UpdateGSTR9CUnreconciled(_ context.Context, tenantID, id uuid.UUID, amount decimal.Decimal) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	f, ok := m.gstr9cFiles[id]
	if !ok || f.TenantID != tenantID {
		return domain.ErrGSTR9CNotFound
	}
	f.UnreconciledAmount = amount
	f.UpdatedAt = time.Now()
	return nil
}

func (m *MockStore) CertifyGSTR9C(_ context.Context, tenantID, id, certifiedBy uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	f, ok := m.gstr9cFiles[id]
	if !ok || f.TenantID != tenantID {
		return domain.ErrGSTR9CNotFound
	}
	now := time.Now()
	f.Status = domain.GSTR9CStatusCertified
	f.IsSelfCertified = true
	f.CertifiedAt = &now
	f.CertifiedBy = &certifiedBy
	f.UpdatedAt = now
	return nil
}

func (m *MockStore) CreateMismatch(_ context.Context, _ uuid.UUID, mm *domain.GSTR9CMismatch) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.mismatches[mm.ID] = mm
	return nil
}

func (m *MockStore) ListMismatches(_ context.Context, tenantID, gstr9cFilingID uuid.UUID) ([]domain.GSTR9CMismatch, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []domain.GSTR9CMismatch
	for _, mm := range m.mismatches {
		if mm.TenantID == tenantID && mm.GSTR9CFilingID == gstr9cFilingID {
			out = append(out, *mm)
		}
	}
	return out, nil
}

func (m *MockStore) GetMismatch(_ context.Context, tenantID, id uuid.UUID) (*domain.GSTR9CMismatch, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	mm, ok := m.mismatches[id]
	if !ok || mm.TenantID != tenantID {
		return nil, domain.ErrMismatchNotFound
	}
	return mm, nil
}

func (m *MockStore) ResolveMismatch(_ context.Context, tenantID, id uuid.UUID, reason string, resolvedBy uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	mm, ok := m.mismatches[id]
	if !ok || mm.TenantID != tenantID {
		return domain.ErrMismatchNotFound
	}
	now := time.Now()
	mm.Resolved = true
	mm.ResolvedReason = reason
	mm.ResolvedAt = &now
	mm.ResolvedBy = &resolvedBy
	return nil
}

func (m *MockStore) DeleteMismatches(_ context.Context, tenantID, gstr9cFilingID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for id, mm := range m.mismatches {
		if mm.TenantID == tenantID && mm.GSTR9CFilingID == gstr9cFilingID {
			delete(m.mismatches, id)
		}
	}
	return nil
}

var _ Repository = (*MockStore)(nil)
