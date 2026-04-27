package provider

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/complai/complai/services/go/ewb-gateway-service/internal/domain"
)

type mockEWB struct {
	EWBNumber  string
	DocType    string
	DocNo      string
	DocDate    string
	FromGSTIN  string
	FromName   string
	ToGSTIN    string
	ToName     string
	VehicleNo  string
	VehicleType string
	DistanceKM int
	TotalValue float64
	Status     string
	EWBDate    string
	ValidUntil time.Time
	GeneratedAt time.Time
}

type MockProvider struct {
	mu      sync.RWMutex
	ewbs    map[string]*mockEWB
	counter atomic.Int64
}

func NewMockProvider() *MockProvider {
	return &MockProvider{
		ewbs: make(map[string]*mockEWB),
	}
}

func (m *MockProvider) validityDays(distanceKm int, vehicleType string) int {
	if distanceKm <= 0 {
		return 1
	}
	divisor := 200
	if vehicleType == "O" {
		divisor = 20
	}
	days := (distanceKm + divisor - 1) / divisor
	if days < 1 {
		return 1
	}
	return days
}

func (m *MockProvider) GenerateEWB(_ context.Context, req *domain.GenerateEWBRequest) (*domain.GenerateEWBResponse, error) {
	if req.GSTIN == "" {
		return nil, fmt.Errorf("gstin is required")
	}
	if req.DocNo == "" {
		return nil, fmt.Errorf("doc_no is required")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, ewb := range m.ewbs {
		if ewb.DocNo == req.DocNo && ewb.FromGSTIN == req.GSTIN && ewb.Status == "ACT" {
			return &domain.GenerateEWBResponse{
				EWBNumber:  ewb.EWBNumber,
				EWBDate:    ewb.EWBDate,
				ValidUntil: ewb.ValidUntil.Format("02/01/2006 15:04:05"),
				Status:     ewb.Status,
			}, nil
		}
	}

	num := m.counter.Add(1)
	ewbNo := fmt.Sprintf("%012d", num)
	now := time.Now()
	days := m.validityDays(req.DistanceKM, req.VehicleType)
	validUntil := now.Add(time.Duration(days) * 24 * time.Hour)

	ewb := &mockEWB{
		EWBNumber:   ewbNo,
		DocType:     req.DocType,
		DocNo:       req.DocNo,
		DocDate:     req.DocDate,
		FromGSTIN:   req.GSTIN,
		FromName:    req.FromName,
		ToGSTIN:     req.ToGSTIN,
		ToName:      req.ToName,
		VehicleNo:   req.VehicleNo,
		VehicleType: req.VehicleType,
		DistanceKM:  req.DistanceKM,
		TotalValue:  req.TotalValue,
		Status:      "ACT",
		EWBDate:     now.Format("02/01/2006 15:04:05"),
		ValidUntil:  validUntil,
		GeneratedAt: now,
	}
	m.ewbs[ewbNo] = ewb

	return &domain.GenerateEWBResponse{
		EWBNumber:  ewbNo,
		EWBDate:    ewb.EWBDate,
		ValidUntil: validUntil.Format("02/01/2006 15:04:05"),
		Status:     "ACT",
	}, nil
}

func (m *MockProvider) CancelEWB(_ context.Context, req *domain.CancelEWBRequest) (*domain.CancelEWBResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ewb, ok := m.ewbs[req.EWBNo]
	if !ok {
		return nil, fmt.Errorf("EWB not found: %s", req.EWBNo)
	}
	if ewb.Status == "CNL" {
		return nil, fmt.Errorf("EWB already cancelled: %s", req.EWBNo)
	}
	if time.Since(ewb.GeneratedAt) > 24*time.Hour {
		return nil, fmt.Errorf("cancellation window expired for EWB: %s", req.EWBNo)
	}

	ewb.Status = "CNL"
	now := time.Now()

	return &domain.CancelEWBResponse{
		EWBNo:      req.EWBNo,
		CancelDate: now.Format("02/01/2006 15:04:05"),
		Status:     "CNL",
	}, nil
}

func (m *MockProvider) GetEWB(_ context.Context, ewbNo string) (*domain.GetEWBResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ewb, ok := m.ewbs[ewbNo]
	if !ok {
		return nil, fmt.Errorf("EWB not found: %s", ewbNo)
	}

	return &domain.GetEWBResponse{
		EWBNumber:  ewb.EWBNumber,
		EWBDate:    ewb.EWBDate,
		DocType:    ewb.DocType,
		DocNo:      ewb.DocNo,
		DocDate:    ewb.DocDate,
		FromGSTIN:  ewb.FromGSTIN,
		FromName:   ewb.FromName,
		ToGSTIN:    ewb.ToGSTIN,
		ToName:     ewb.ToName,
		VehicleNo:  ewb.VehicleNo,
		Status:     ewb.Status,
		ValidUntil: ewb.ValidUntil.Format("02/01/2006 15:04:05"),
		DistanceKM: ewb.DistanceKM,
		TotalValue: ewb.TotalValue,
	}, nil
}

func (m *MockProvider) UpdateVehicle(_ context.Context, req *domain.UpdateVehicleRequest) (*domain.UpdateVehicleResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ewb, ok := m.ewbs[req.EWBNo]
	if !ok {
		return nil, fmt.Errorf("EWB not found: %s", req.EWBNo)
	}
	if ewb.Status != "ACT" {
		return nil, fmt.Errorf("EWB not active: %s (status: %s)", req.EWBNo, ewb.Status)
	}

	ewb.VehicleNo = req.VehicleNo

	return &domain.UpdateVehicleResponse{
		EWBNo:      req.EWBNo,
		VehicleNo:  req.VehicleNo,
		ValidUntil: ewb.ValidUntil.Format("02/01/2006 15:04:05"),
		Status:     "ACT",
	}, nil
}

func (m *MockProvider) ExtendValidity(_ context.Context, req *domain.ExtendValidityRequest) (*domain.ExtendValidityResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ewb, ok := m.ewbs[req.EWBNo]
	if !ok {
		return nil, fmt.Errorf("EWB not found: %s", req.EWBNo)
	}
	if ewb.Status != "ACT" {
		return nil, fmt.Errorf("EWB not active: %s (status: %s)", req.EWBNo, ewb.Status)
	}

	days := m.validityDays(req.RemainingDistance, ewb.VehicleType)
	ewb.ValidUntil = ewb.ValidUntil.Add(time.Duration(days) * 24 * time.Hour)

	return &domain.ExtendValidityResponse{
		EWBNo:      req.EWBNo,
		ValidUntil: ewb.ValidUntil.Format("02/01/2006 15:04:05"),
		Status:     "ACT",
	}, nil
}

func (m *MockProvider) ConsolidateEWB(_ context.Context, req *domain.ConsolidateEWBRequest) (*domain.ConsolidateEWBResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(req.EWBNumbers) < 2 {
		return nil, fmt.Errorf("at least 2 EWBs required for consolidation")
	}

	for _, ewbNo := range req.EWBNumbers {
		ewb, ok := m.ewbs[ewbNo]
		if !ok {
			return nil, fmt.Errorf("EWB not found: %s", ewbNo)
		}
		if ewb.Status != "ACT" {
			return nil, fmt.Errorf("EWB not active: %s (status: %s)", ewbNo, ewb.Status)
		}
	}

	num := m.counter.Add(1)
	cewbNo := fmt.Sprintf("C%011d", num)

	return &domain.ConsolidateEWBResponse{
		ConsolidatedEWBNo: cewbNo,
		Status:            "ACT",
	}, nil
}

var _ EWBProvider = (*MockProvider)(nil)
