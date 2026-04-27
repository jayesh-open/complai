package store

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/ewb-service/internal/domain"
)

type MockStore struct {
	mu             sync.RWMutex
	ewbs           map[uuid.UUID]*domain.EWayBill
	items          map[uuid.UUID][]domain.EWBItem
	vehicleUpdates map[uuid.UUID][]domain.VehicleUpdate
	consolidations map[uuid.UUID]*domain.Consolidation
}

func NewMockStore() *MockStore {
	return &MockStore{
		ewbs:           make(map[uuid.UUID]*domain.EWayBill),
		items:          make(map[uuid.UUID][]domain.EWBItem),
		vehicleUpdates: make(map[uuid.UUID][]domain.VehicleUpdate),
		consolidations: make(map[uuid.UUID]*domain.Consolidation),
	}
}

func (m *MockStore) CreateEWB(_ context.Context, tenantID uuid.UUID, ewb *domain.EWayBill) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ewb.ID = uuid.New()
	ewb.TenantID = tenantID
	ewb.RequestID = uuid.New()
	ewb.Status = domain.EWBStatusPending
	now := time.Now()
	ewb.CreatedAt = now
	ewb.UpdatedAt = now

	cp := *ewb
	m.ewbs[ewb.ID] = &cp
	return nil
}

func (m *MockStore) GetEWB(_ context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EWayBill, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ewb, ok := m.ewbs[id]
	if !ok || ewb.TenantID != tenantID {
		return nil, fmt.Errorf("ewb not found: %s", id)
	}
	cp := *ewb
	return &cp, nil
}

func (m *MockStore) GetEWBByNumber(_ context.Context, tenantID uuid.UUID, ewbNumber string) (*domain.EWayBill, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, ewb := range m.ewbs {
		if ewb.EWBNumber == ewbNumber && ewb.TenantID == tenantID {
			cp := *ewb
			return &cp, nil
		}
	}
	return nil, fmt.Errorf("ewb not found for number: %s", ewbNumber)
}

func (m *MockStore) ListEWBs(_ context.Context, tenantID uuid.UUID, req *domain.ListEWBRequest) ([]domain.EWayBill, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var all []domain.EWayBill
	for _, ewb := range m.ewbs {
		if ewb.TenantID == tenantID && ewb.SupplierGSTIN == req.GSTIN {
			if req.Status != "" && string(ewb.Status) != req.Status {
				continue
			}
			all = append(all, *ewb)
		}
	}
	total := len(all)
	start := req.PageOffset
	if start > total {
		start = total
	}
	end := start + req.PageSize
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

func (m *MockStore) UpdateEWBGenerated(_ context.Context, tenantID uuid.UUID, id uuid.UUID, ewbNumber string, validFrom, validUntil time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ewb, ok := m.ewbs[id]
	if !ok || ewb.TenantID != tenantID {
		return fmt.Errorf("ewb not found: %s", id)
	}
	now := time.Now()
	ewb.EWBNumber = ewbNumber
	ewb.Status = domain.EWBStatusActive
	ewb.ValidFrom = &validFrom
	ewb.ValidUntil = &validUntil
	ewb.GeneratedAt = &now
	ewb.UpdatedAt = now
	return nil
}

func (m *MockStore) UpdateEWBCancelled(_ context.Context, tenantID uuid.UUID, id uuid.UUID, reason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ewb, ok := m.ewbs[id]
	if !ok || ewb.TenantID != tenantID {
		return fmt.Errorf("ewb not found: %s", id)
	}
	now := time.Now()
	ewb.Status = domain.EWBStatusCancelled
	ewb.CancelledAt = &now
	ewb.CancelReason = reason
	ewb.UpdatedAt = now
	return nil
}

func (m *MockStore) UpdateEWBStatus(_ context.Context, tenantID uuid.UUID, id uuid.UUID, status domain.EWBStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ewb, ok := m.ewbs[id]
	if !ok || ewb.TenantID != tenantID {
		return fmt.Errorf("ewb not found: %s", id)
	}
	ewb.Status = status
	ewb.UpdatedAt = time.Now()
	return nil
}

func (m *MockStore) UpdateEWBVehicle(_ context.Context, tenantID uuid.UUID, id uuid.UUID, vehicleNumber string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ewb, ok := m.ewbs[id]
	if !ok || ewb.TenantID != tenantID {
		return fmt.Errorf("ewb not found: %s", id)
	}
	ewb.VehicleNumber = vehicleNumber
	ewb.Status = domain.EWBStatusVehicleUpdated
	ewb.UpdatedAt = time.Now()
	return nil
}

func (m *MockStore) UpdateEWBValidity(_ context.Context, tenantID uuid.UUID, id uuid.UUID, validUntil time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ewb, ok := m.ewbs[id]
	if !ok || ewb.TenantID != tenantID {
		return fmt.Errorf("ewb not found: %s", id)
	}
	ewb.ValidUntil = &validUntil
	ewb.Status = domain.EWBStatusExtended
	ewb.UpdatedAt = time.Now()
	return nil
}

func (m *MockStore) SetConsolidatedEWBID(_ context.Context, tenantID uuid.UUID, id uuid.UUID, consolidatedID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ewb, ok := m.ewbs[id]
	if !ok || ewb.TenantID != tenantID {
		return fmt.Errorf("ewb not found: %s", id)
	}
	ewb.ConsolidatedEWBID = &consolidatedID
	ewb.Status = domain.EWBStatusConsolidated
	ewb.UpdatedAt = time.Now()
	return nil
}

func (m *MockStore) CreateItems(_ context.Context, tenantID uuid.UUID, items []domain.EWBItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(items) == 0 {
		return nil
	}
	ewbID := items[0].EWBID
	for i := range items {
		items[i].ID = uuid.New()
		items[i].TenantID = tenantID
		items[i].ItemNumber = i + 1
		items[i].CreatedAt = time.Now()
	}
	m.items[ewbID] = append(m.items[ewbID], items...)
	return nil
}

func (m *MockStore) GetItems(_ context.Context, tenantID uuid.UUID, ewbID uuid.UUID) ([]domain.EWBItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	items := m.items[ewbID]
	var result []domain.EWBItem
	for _, item := range items {
		if item.TenantID == tenantID {
			result = append(result, item)
		}
	}
	return result, nil
}

func (m *MockStore) CreateVehicleUpdate(_ context.Context, tenantID uuid.UUID, update *domain.VehicleUpdate) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	update.ID = uuid.New()
	update.TenantID = tenantID
	update.UpdatedAt = time.Now()
	cp := *update
	m.vehicleUpdates[update.EWBID] = append(m.vehicleUpdates[update.EWBID], cp)
	return nil
}

func (m *MockStore) GetVehicleUpdates(_ context.Context, tenantID uuid.UUID, ewbID uuid.UUID) ([]domain.VehicleUpdate, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	updates := m.vehicleUpdates[ewbID]
	var result []domain.VehicleUpdate
	for _, u := range updates {
		if u.TenantID == tenantID {
			result = append(result, u)
		}
	}
	return result, nil
}

func (m *MockStore) CreateConsolidation(_ context.Context, tenantID uuid.UUID, consolidation *domain.Consolidation) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	consolidation.ID = uuid.New()
	consolidation.TenantID = tenantID
	now := time.Now()
	consolidation.CreatedAt = now
	consolidation.UpdatedAt = now
	consolidation.GeneratedAt = &now
	cp := *consolidation
	m.consolidations[consolidation.ID] = &cp
	return nil
}

func (m *MockStore) SetGeneratedAt(id uuid.UUID, t *time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if ewb, ok := m.ewbs[id]; ok {
		ewb.GeneratedAt = t
	}
}

var _ Repository = (*MockStore)(nil)
