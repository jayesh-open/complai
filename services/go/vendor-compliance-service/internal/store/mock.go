package store

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/vendor-compliance-service/internal/domain"
)

var _ Repository = (*MockStore)(nil)

type MockStore struct {
	mu       sync.RWMutex
	vendors  map[string]*domain.VendorSnapshot // key: tenantID:vendorID
	scores   []domain.ComplianceScore
	syncs    []domain.SyncStatus
}

func NewMockStore() *MockStore {
	return &MockStore{
		vendors: make(map[string]*domain.VendorSnapshot),
	}
}

func vendorKey(tenantID uuid.UUID, vendorID string) string {
	return tenantID.String() + ":" + vendorID
}

func (m *MockStore) UpsertVendorSnapshot(_ context.Context, tenantID uuid.UUID, v *domain.VendorSnapshot) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	v.TenantID = tenantID
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}

	key := vendorKey(tenantID, v.VendorID)
	if existing, ok := m.vendors[key]; ok {
		v.ID = existing.ID
		v.CreatedAt = existing.CreatedAt
	} else {
		v.CreatedAt = time.Now().UTC()
	}
	v.UpdatedAt = time.Now().UTC()

	copy := *v
	m.vendors[key] = &copy
	return nil
}

func (m *MockStore) ListVendorSnapshots(_ context.Context, tenantID uuid.UUID) ([]domain.VendorSnapshot, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []domain.VendorSnapshot
	for _, v := range m.vendors {
		if v.TenantID == tenantID {
			result = append(result, *v)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, nil
}

func (m *MockStore) GetVendorSnapshot(_ context.Context, tenantID uuid.UUID, vendorID string) (*domain.VendorSnapshot, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := vendorKey(tenantID, vendorID)
	v, ok := m.vendors[key]
	if !ok {
		return nil, fmt.Errorf("vendor not found")
	}
	copy := *v
	return &copy, nil
}

func (m *MockStore) CreateComplianceScore(_ context.Context, tenantID uuid.UUID, s *domain.ComplianceScore) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	s.TenantID = tenantID
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	s.CreatedAt = time.Now().UTC()

	m.scores = append(m.scores, *s)
	return nil
}

func (m *MockStore) GetLatestScore(_ context.Context, tenantID uuid.UUID, vendorID string) (*domain.ComplianceScore, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var latest *domain.ComplianceScore
	for i := range m.scores {
		s := &m.scores[i]
		if s.TenantID == tenantID && s.VendorID == vendorID {
			if latest == nil || s.ScoredAt.After(latest.ScoredAt) {
				copy := *s
				latest = &copy
			}
		}
	}

	if latest == nil {
		return nil, fmt.Errorf("score not found")
	}
	return latest, nil
}

func (m *MockStore) ListLatestScores(_ context.Context, tenantID uuid.UUID) ([]domain.ComplianceScore, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	latestByVendor := make(map[string]*domain.ComplianceScore)
	for i := range m.scores {
		s := &m.scores[i]
		if s.TenantID == tenantID {
			if existing, ok := latestByVendor[s.VendorID]; !ok || s.ScoredAt.After(existing.ScoredAt) {
				copy := *s
				latestByVendor[s.VendorID] = &copy
			}
		}
	}

	var result []domain.ComplianceScore
	for _, s := range latestByVendor {
		result = append(result, *s)
	}
	return result, nil
}

func (m *MockStore) GetScoreSummary(_ context.Context, tenantID uuid.UUID) (*domain.ScoreSummary, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	latestByVendor := make(map[string]*domain.ComplianceScore)
	for i := range m.scores {
		s := &m.scores[i]
		if s.TenantID == tenantID {
			if existing, ok := latestByVendor[s.VendorID]; !ok || s.ScoredAt.After(existing.ScoredAt) {
				copy := *s
				latestByVendor[s.VendorID] = &copy
			}
		}
	}

	summary := &domain.ScoreSummary{}
	totalScore := 0
	for _, s := range latestByVendor {
		summary.Total++
		totalScore += s.TotalScore
		switch s.Category {
		case "A":
			summary.CatA++
		case "B":
			summary.CatB++
		case "C":
			summary.CatC++
		case "D":
			summary.CatD++
		}
	}
	if summary.Total > 0 {
		summary.AvgScore = totalScore / summary.Total
	}

	return summary, nil
}

func (m *MockStore) CreateSyncStatus(_ context.Context, tenantID uuid.UUID, s *domain.SyncStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	s.TenantID = tenantID
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	s.CreatedAt = time.Now().UTC()

	m.syncs = append(m.syncs, *s)
	return nil
}

func (m *MockStore) UpdateSyncStatus(_ context.Context, _ uuid.UUID, syncID uuid.UUID, status string, vendorCount, scoredCount int, errMsg string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := range m.syncs {
		if m.syncs[i].ID == syncID {
			m.syncs[i].Status = status
			m.syncs[i].VendorCount = vendorCount
			m.syncs[i].ScoredCount = scoredCount
			m.syncs[i].ErrorMessage = errMsg
			now := time.Now().UTC()
			m.syncs[i].CompletedAt = &now
			return nil
		}
	}
	return fmt.Errorf("sync not found")
}

func (m *MockStore) GetLatestSyncStatus(_ context.Context, tenantID uuid.UUID) (*domain.SyncStatus, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var latest *domain.SyncStatus
	for i := range m.syncs {
		s := &m.syncs[i]
		if s.TenantID == tenantID {
			if latest == nil || s.StartedAt.After(latest.StartedAt) {
				copy := *s
				latest = &copy
			}
		}
	}

	if latest == nil {
		return nil, fmt.Errorf("no sync status found")
	}
	return latest, nil
}
