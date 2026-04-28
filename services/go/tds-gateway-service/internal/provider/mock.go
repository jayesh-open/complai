package provider

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/complai/complai/services/go/tds-gateway-service/internal/domain"
	"github.com/google/uuid"
)

type MockProvider struct {
	mu       sync.RWMutex
	challans map[string]*domain.ChallanResponse
	filings  map[string]*domain.FormFilingResponse
}

func NewMockProvider() *MockProvider {
	return &MockProvider{
		challans: make(map[string]*domain.ChallanResponse),
		filings:  make(map[string]*domain.FormFilingResponse),
	}
}

func (m *MockProvider) VerifyPAN(_ context.Context, req domain.PANVerifyRequest) (*domain.PANVerifyResponse, error) {
	pan := strings.ToUpper(req.PAN)
	if len(pan) != 10 {
		return &domain.PANVerifyResponse{PAN: pan, Status: "INVALID"}, nil
	}

	return &domain.PANVerifyResponse{
		PAN:      pan,
		Name:     nameForPAN(pan, req.Name),
		Status:   "VALID",
		Category: panCategory(pan[3]),
	}, nil
}

func panCategory(fourth byte) string {
	switch fourth {
	case 'P':
		return "INDIVIDUAL"
	case 'C':
		return "COMPANY"
	case 'H':
		return "HUF"
	case 'F':
		return "FIRM"
	case 'T':
		return "TRUST"
	case 'A':
		return "AOP"
	case 'L':
		return "LOCAL_AUTHORITY"
	case 'G':
		return "GOVERNMENT"
	default:
		return "INDIVIDUAL"
	}
}

func nameForPAN(pan, provided string) string {
	if provided != "" {
		return provided
	}
	return "Mock Entity " + pan[:5]
}

func (m *MockProvider) VerifyTAN(_ context.Context, req domain.TANVerifyRequest) (*domain.TANVerifyResponse, error) {
	tan := strings.ToUpper(req.TAN)
	if len(tan) != 10 {
		return &domain.TANVerifyResponse{TAN: tan, Status: "NOT_FOUND"}, nil
	}
	return &domain.TANVerifyResponse{
		TAN:    tan,
		Name:   "Mock Deductor " + tan[:4],
		Status: "ACTIVE",
	}, nil
}

func (m *MockProvider) GenerateChallan(_ context.Context, req domain.ChallanRequest) (*domain.ChallanResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	challan := &domain.ChallanResponse{
		ChallanNumber: fmt.Sprintf("CHL%s", uuid.New().String()[:8]),
		BSRCode:       "0001234",
		DepositDate:   time.Now().Format("2006-01-02"),
		Amount:        req.Amount + req.Surcharge + req.Cess + req.Interest + req.Penalty,
		Status:        "SUCCESS",
	}
	m.challans[challan.ChallanNumber] = challan
	return challan, nil
}

func (m *MockProvider) FileForm26Q(_ context.Context, req domain.Form26QRequest) (*domain.FormFilingResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	resp := &domain.FormFilingResponse{
		TokenNumber:           fmt.Sprintf("TKN26Q%s", uuid.New().String()[:8]),
		AcknowledgementNumber: fmt.Sprintf("ACK%s%s", req.Quarter, uuid.New().String()[:6]),
		FilingDate:            time.Now().Format("2006-01-02"),
		Status:                "ACCEPTED",
	}
	m.filings[resp.TokenNumber] = resp
	return resp, nil
}

func (m *MockProvider) FileForm24Q(_ context.Context, req domain.Form24QRequest) (*domain.FormFilingResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	resp := &domain.FormFilingResponse{
		TokenNumber:           fmt.Sprintf("TKN24Q%s", uuid.New().String()[:8]),
		AcknowledgementNumber: fmt.Sprintf("ACK%s%s", req.Quarter, uuid.New().String()[:6]),
		FilingDate:            time.Now().Format("2006-01-02"),
		Status:                "ACCEPTED",
	}
	m.filings[resp.TokenNumber] = resp
	return resp, nil
}
