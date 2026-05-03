package provider

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/complai/complai/services/go/itr-gateway-service/internal/domain"
	"github.com/google/uuid"
)

type MockProvider struct {
	mu       sync.RWMutex
	filings  map[string]*domain.ITRSubmitResponse
	itrvs    map[string]*domain.ITRVResponse
	everify  map[string]*domain.EVerifyResponse
}

func NewMockProvider() *MockProvider {
	return &MockProvider{
		filings: make(map[string]*domain.ITRSubmitResponse),
		itrvs:   make(map[string]*domain.ITRVResponse),
		everify: make(map[string]*domain.EVerifyResponse),
	}
}

func (m *MockProvider) CheckPANAadhaarLink(_ context.Context, req domain.PANAadhaarLinkRequest) (*domain.PANAadhaarLinkResponse, error) {
	pan := strings.ToUpper(req.PAN)
	if len(pan) != 10 {
		return &domain.PANAadhaarLinkResponse{PAN: pan, Linked: false}, nil
	}
	if pan[9] == 'Z' {
		return &domain.PANAadhaarLinkResponse{PAN: pan, Linked: false}, nil
	}
	return &domain.PANAadhaarLinkResponse{
		PAN:          pan,
		Linked:       true,
		LinkDate:     "2024-06-15",
		AadhaarLast4: "7890",
	}, nil
}

func (m *MockProvider) FetchAIS(_ context.Context, req domain.AISRequest) (*domain.AISResponse, error) {
	pan := strings.ToUpper(req.PAN)
	if len(pan) != 10 {
		return nil, fmt.Errorf("invalid PAN format")
	}
	return &domain.AISResponse{
		PAN:        pan,
		TaxYear:    req.TaxYear,
		Form168Ref: fmt.Sprintf("168-%s-%s", pan, req.TaxYear),
		TDSEntries: []domain.AISTDSEntry{
			{
				DeductorTAN:  "MUMB12345A",
				DeductorName: "Mock Employer Ltd",
				Section:      "392",
				Amount:       1200000,
				TDSAmount:    60000,
				Quarter:      "Q4",
			},
			{
				DeductorTAN:  "DELH67890B",
				DeductorName: "Mock Bank Ltd",
				Section:      "393(1)",
				Amount:       50000,
				TDSAmount:    5000,
				Quarter:      "Q4",
			},
		},
		InterestIncome:    50000,
		DividendIncome:    25000,
		SalaryIncome:      1200000,
		SecuritiesTrading: 150000,
	}, nil
}

func (m *MockProvider) SubmitITR(_ context.Context, req domain.ITRSubmitRequest) (*domain.ITRSubmitResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	pan := strings.ToUpper(req.PAN)
	if len(pan) != 10 {
		return nil, fmt.Errorf("invalid PAN format")
	}
	validForms := map[string]bool{"ITR-1": true, "ITR-2": true, "ITR-3": true, "ITR-4": true}
	if !validForms[req.FormType] {
		return nil, fmt.Errorf("unsupported form type: %s", req.FormType)
	}

	arn := fmt.Sprintf("ARN-%s-%s", req.TaxYear, uuid.New().String()[:8])
	resp := &domain.ITRSubmitResponse{
		ARN:               arn,
		AcknowledgementNo: fmt.Sprintf("ACK%s", uuid.New().String()[:10]),
		FilingDate:        time.Now().Format("2006-01-02"),
		Status:            "SUBMITTED",
	}
	m.filings[arn] = resp
	return resp, nil
}

func (m *MockProvider) GenerateITRV(_ context.Context, req domain.ITRVRequest) (*domain.ITRVResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.filings[req.ARN]; !ok {
		return nil, fmt.Errorf("filing not found for ARN: %s", req.ARN)
	}
	return &domain.ITRVResponse{
		ARN:               req.ARN,
		ITRVURL:           fmt.Sprintf("https://mock.incometax.gov.in/itrv/%s.pdf", req.ARN),
		AcknowledgementNo: m.filings[req.ARN].AcknowledgementNo,
	}, nil
}

func (m *MockProvider) CheckEVerification(_ context.Context, req domain.EVerifyRequest) (*domain.EVerifyResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.filings[req.ARN]; !ok {
		return nil, fmt.Errorf("filing not found for ARN: %s", req.ARN)
	}
	return &domain.EVerifyResponse{
		ARN:      req.ARN,
		Verified: true,
		Method:   "AADHAAR_OTP",
		Date:     time.Now().Format("2006-01-02"),
	}, nil
}

func (m *MockProvider) CheckRefundStatus(_ context.Context, req domain.RefundStatusRequest) (*domain.RefundStatusResponse, error) {
	pan := strings.ToUpper(req.PAN)
	if len(pan) != 10 {
		return nil, fmt.Errorf("invalid PAN format")
	}
	return &domain.RefundStatusResponse{
		PAN:     pan,
		TaxYear: req.TaxYear,
		Status:  "PROCESSED",
		Amount:  15000,
		BankRef: fmt.Sprintf("NEFT-%s", uuid.New().String()[:8]),
	}, nil
}
