package provider

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/gstn-gateway-service/internal/domain"
)

var _ GSTNProvider = (*MockProvider)(nil)

type MockProvider struct {
	mu       sync.RWMutex
	filings  map[string]*domain.MockFiling // key: gstin:ret_period
	requests map[string]interface{}        // idempotency: request_id → response
}

func NewMockProvider() *MockProvider {
	return &MockProvider{
		filings:  make(map[string]*domain.MockFiling),
		requests: make(map[string]interface{}),
	}
}

func filingKey(gstin, retPeriod string) string {
	return gstin + ":" + retPeriod
}

func (m *MockProvider) Authenticate(_ context.Context) (*domain.AuthResponse, error) {
	return &domain.AuthResponse{
		AccessToken: "mock-gsp-token-" + uuid.New().String()[:8],
		TokenType:   "bearer",
		ExpiresIn:   86399,
		Scope:       "gsp",
		JTI:         uuid.New().String(),
	}, nil
}

func (m *MockProvider) GSTR1Save(_ context.Context, req *domain.GSTR1SaveRequest) (*domain.GSTR1SaveResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if resp, ok := m.requests[req.RequestID]; ok {
		return resp.(*domain.GSTR1SaveResponse), nil
	}

	key := filingKey(req.GSTIN, req.RetPeriod)
	f, ok := m.filings[key]
	if !ok {
		f = &domain.MockFiling{
			GSTIN:     req.GSTIN,
			RetPeriod: req.RetPeriod,
			Status:    domain.StatusDraft,
			Sections:  make(map[string]interface{}),
			Token:     uuid.New().String()[:16],
		}
		m.filings[key] = f
	}

	if f.Status == domain.StatusSubmitted || f.Status == domain.StatusFiled {
		return nil, fmt.Errorf("cannot save: filing is %s", f.Status)
	}

	f.Sections[req.Section] = req.Data
	f.Status = domain.StatusSaved

	resp := &domain.GSTR1SaveResponse{
		Status:    "success",
		RequestID: req.RequestID,
		Token:     f.Token,
		Message:   fmt.Sprintf("Section %s saved successfully", req.Section),
		SavedAt:   time.Now().UTC(),
	}

	m.requests[req.RequestID] = resp
	return resp, nil
}

func (m *MockProvider) GSTR1Get(_ context.Context, req *domain.GSTR1GetRequest) (*domain.GSTR1GetResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := filingKey(req.GSTIN, req.RetPeriod)
	f, ok := m.filings[key]
	if !ok {
		return &domain.GSTR1GetResponse{
			GSTIN:     req.GSTIN,
			RetPeriod: req.RetPeriod,
			Data:      map[string]interface{}{},
			Status:    string(domain.StatusDraft),
			RequestID: req.RequestID,
		}, nil
	}

	data := make(map[string]interface{})
	if req.Section != "" {
		if sectionData, exists := f.Sections[req.Section]; exists {
			data[req.Section] = sectionData
		}
	} else {
		for k, v := range f.Sections {
			data[k] = v
		}
	}

	return &domain.GSTR1GetResponse{
		GSTIN:     req.GSTIN,
		RetPeriod: req.RetPeriod,
		Data:      data,
		Status:    string(f.Status),
		RequestID: req.RequestID,
	}, nil
}

func (m *MockProvider) GSTR1Reset(_ context.Context, req *domain.GSTR1ResetRequest) (*domain.GSTR1ResetResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if resp, ok := m.requests[req.RequestID]; ok {
		return resp.(*domain.GSTR1ResetResponse), nil
	}

	key := filingKey(req.GSTIN, req.RetPeriod)
	f, ok := m.filings[key]
	if !ok {
		resp := &domain.GSTR1ResetResponse{
			Status:    "success",
			RequestID: req.RequestID,
			Message:   "Nothing to reset",
		}
		m.requests[req.RequestID] = resp
		return resp, nil
	}

	if f.Status == domain.StatusFiled {
		return nil, fmt.Errorf("cannot reset: filing is already filed")
	}

	f.Sections = make(map[string]interface{})
	f.Status = domain.StatusDraft
	f.Token = uuid.New().String()[:16]

	resp := &domain.GSTR1ResetResponse{
		Status:    "success",
		RequestID: req.RequestID,
		Message:   "Draft reset successfully",
	}
	m.requests[req.RequestID] = resp
	return resp, nil
}

func (m *MockProvider) GSTR1Submit(_ context.Context, req *domain.GSTR1SubmitRequest) (*domain.GSTR1SubmitResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if resp, ok := m.requests[req.RequestID]; ok {
		return resp.(*domain.GSTR1SubmitResponse), nil
	}

	key := filingKey(req.GSTIN, req.RetPeriod)
	f, ok := m.filings[key]
	if !ok {
		return nil, fmt.Errorf("no draft found for %s/%s", req.GSTIN, req.RetPeriod)
	}

	if f.Status == domain.StatusFiled {
		return nil, fmt.Errorf("already filed")
	}
	if f.Status == domain.StatusSubmitted {
		return nil, fmt.Errorf("already submitted")
	}
	if len(f.Sections) == 0 {
		return nil, fmt.Errorf("no sections saved")
	}

	f.Status = domain.StatusSubmitted

	resp := &domain.GSTR1SubmitResponse{
		Status:    "success",
		RequestID: req.RequestID,
		Token:     f.Token,
		Message:   "GSTR-1 submitted successfully. Locked for filing.",
	}
	m.requests[req.RequestID] = resp
	return resp, nil
}

func (m *MockProvider) GSTR1File(_ context.Context, req *domain.GSTR1FileRequest) (*domain.GSTR1FileResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if resp, ok := m.requests[req.RequestID]; ok {
		return resp.(*domain.GSTR1FileResponse), nil
	}

	key := filingKey(req.GSTIN, req.RetPeriod)
	f, ok := m.filings[key]
	if !ok {
		return nil, fmt.Errorf("no filing found for %s/%s", req.GSTIN, req.RetPeriod)
	}

	if f.Status == domain.StatusFiled {
		return nil, fmt.Errorf("already filed")
	}
	if f.Status != domain.StatusSubmitted {
		return nil, fmt.Errorf("must submit before filing (current status: %s)", f.Status)
	}

	if req.SignType != "DSC" && req.SignType != "EVC" {
		return nil, fmt.Errorf("invalid sign_type: %s (must be DSC or EVC)", req.SignType)
	}

	if req.SignType == "EVC" && req.EVOTP == "" {
		return nil, fmt.Errorf("EVC OTP required for EVC signing")
	}

	now := time.Now().UTC()
	f.Status = domain.StatusFiled
	gstinPrefix := req.GSTIN
	if len(gstinPrefix) > 2 {
		gstinPrefix = gstinPrefix[:2]
	}
	f.ARN = fmt.Sprintf("AA%s%s%s", gstinPrefix, req.RetPeriod, uuid.New().String()[:8])
	f.FiledAt = &now

	resp := &domain.GSTR1FileResponse{
		Status:    "success",
		ARN:       f.ARN,
		RequestID: req.RequestID,
		Message:   "GSTR-1 filed successfully",
		FiledAt:   now,
	}
	m.requests[req.RequestID] = resp
	return resp, nil
}

func (m *MockProvider) GSTR1Status(_ context.Context, req *domain.GSTR1StatusRequest) (*domain.GSTR1StatusResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := filingKey(req.GSTIN, req.RetPeriod)
	f, ok := m.filings[key]
	if !ok {
		return &domain.GSTR1StatusResponse{
			GSTIN:     req.GSTIN,
			RetPeriod: req.RetPeriod,
			Status:    string(domain.StatusDraft),
			RequestID: req.RequestID,
		}, nil
	}

	return &domain.GSTR1StatusResponse{
		GSTIN:     req.GSTIN,
		RetPeriod: req.RetPeriod,
		Status:    string(f.Status),
		ARN:       f.ARN,
		FiledAt:   f.FiledAt,
		RequestID: req.RequestID,
	}, nil
}

func (m *MockProvider) ResetState() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.filings = make(map[string]*domain.MockFiling)
	m.requests = make(map[string]interface{})
}
