package store

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/gst-service/internal/domain"
)

var _ Repository = (*MockStore)(nil)

type MockStore struct {
	mu       sync.RWMutex
	filings  map[uuid.UUID]*domain.GSTR1Filing
	entries  []domain.SalesRegisterEntry
	sections []domain.GSTR1Section
	valErrs  []domain.ValidationError
}

func NewMockStore() *MockStore {
	return &MockStore{
		filings: make(map[uuid.UUID]*domain.GSTR1Filing),
	}
}

func (m *MockStore) CreateFiling(_ context.Context, tenantID uuid.UUID, f *domain.GSTR1Filing) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	f.ID = uuid.New()
	f.TenantID = tenantID
	f.RequestID = uuid.New()
	f.CreatedAt = time.Now().UTC()
	f.UpdatedAt = f.CreatedAt

	copy := *f
	if f.CreatedBy != nil {
		cb := *f.CreatedBy
		copy.CreatedBy = &cb
	}
	m.filings[f.ID] = &copy
	return nil
}

func (m *MockStore) GetFiling(_ context.Context, _ uuid.UUID, filingID uuid.UUID) (*domain.GSTR1Filing, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	f, ok := m.filings[filingID]
	if !ok {
		return nil, fmt.Errorf("filing not found")
	}
	copy := *f
	return &copy, nil
}

func (m *MockStore) GetFilingByPeriod(_ context.Context, _ uuid.UUID, gstin, period string) (*domain.GSTR1Filing, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, f := range m.filings {
		if f.GSTIN == gstin && f.ReturnPeriod == period {
			copy := *f
			return &copy, nil
		}
	}
	return nil, fmt.Errorf("filing not found")
}

func (m *MockStore) UpdateFilingStatus(_ context.Context, _ uuid.UUID, filingID uuid.UUID, status domain.FilingStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	f, ok := m.filings[filingID]
	if !ok {
		return fmt.Errorf("filing not found")
	}
	f.Status = status
	f.UpdatedAt = time.Now().UTC()
	return nil
}

func (m *MockStore) UpdateFilingARN(_ context.Context, _ uuid.UUID, filingID uuid.UUID, arn string, filedBy uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	f, ok := m.filings[filingID]
	if !ok {
		return fmt.Errorf("filing not found")
	}
	f.Status = domain.FilingStatusFiled
	f.ARN = arn
	now := time.Now().UTC()
	f.FiledAt = &now
	f.FiledBy = &filedBy
	f.UpdatedAt = now
	return nil
}

func (m *MockStore) ApproveFiling(_ context.Context, _ uuid.UUID, filingID uuid.UUID, approvedBy uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	f, ok := m.filings[filingID]
	if !ok {
		return fmt.Errorf("filing not found")
	}
	f.Status = domain.FilingStatusApproved
	now := time.Now().UTC()
	f.ApprovedBy = &approvedBy
	f.ApprovedAt = &now
	f.UpdatedAt = now
	return nil
}

func (m *MockStore) BulkInsertEntries(_ context.Context, tenantID uuid.UUID, entries []domain.SalesRegisterEntry) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	inserted := 0
	for _, e := range entries {
		dup := false
		for _, existing := range m.entries {
			if existing.TenantID == tenantID && existing.GSTIN == e.GSTIN && existing.DocumentNumber == e.DocumentNumber {
				dup = true
				break
			}
		}
		if !dup {
			e.TenantID = tenantID
			e.CreatedAt = time.Now().UTC()
			e.UpdatedAt = e.CreatedAt
			m.entries = append(m.entries, e)
			inserted++
		}
	}
	return inserted, nil
}

func (m *MockStore) ListEntries(_ context.Context, _ uuid.UUID, filingID uuid.UUID, section string) ([]domain.SalesRegisterEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	f, ok := m.filings[filingID]
	if !ok {
		return nil, fmt.Errorf("filing not found")
	}

	var result []domain.SalesRegisterEntry
	for _, e := range m.entries {
		if e.GSTIN == f.GSTIN && e.ReturnPeriod == f.ReturnPeriod {
			if section == "" || e.Section == section {
				result = append(result, e)
			}
		}
	}
	return result, nil
}

func (m *MockStore) CountEntries(_ context.Context, _ uuid.UUID, gstin, period string) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, e := range m.entries {
		if e.GSTIN == gstin && e.ReturnPeriod == period {
			count++
		}
	}
	return count, nil
}

func (m *MockStore) CreateSections(_ context.Context, tenantID uuid.UUID, sections []domain.GSTR1Section) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, s := range sections {
		s.TenantID = tenantID
		s.CreatedAt = time.Now().UTC()
		m.sections = append(m.sections, s)
	}
	return nil
}

func (m *MockStore) ListSections(_ context.Context, _ uuid.UUID, filingID uuid.UUID) ([]domain.GSTR1Section, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []domain.GSTR1Section
	for _, s := range m.sections {
		if s.FilingID == filingID {
			result = append(result, s)
		}
	}
	return result, nil
}

func (m *MockStore) CreateValidationErrors(_ context.Context, tenantID uuid.UUID, errs []domain.ValidationError) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, e := range errs {
		e.TenantID = tenantID
		e.CreatedAt = time.Now().UTC()
		m.valErrs = append(m.valErrs, e)
	}
	return nil
}

func (m *MockStore) ListValidationErrors(_ context.Context, _ uuid.UUID, filingID uuid.UUID) ([]domain.ValidationError, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []domain.ValidationError
	for _, e := range m.valErrs {
		if e.FilingID == filingID {
			result = append(result, e)
		}
	}
	return result, nil
}

func (m *MockStore) CountValidationErrors(_ context.Context, _ uuid.UUID, filingID uuid.UUID) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, e := range m.valErrs {
		if e.FilingID == filingID {
			count++
		}
	}
	return count, nil
}
