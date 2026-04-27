package provider

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/irp-gateway-service/internal/domain"
)

type MockProvider struct {
	mu   sync.RWMutex
	irns map[string]*storedIRN // keyed by IRN
	docs map[string]string     // "doctype|docnum|docdate" → IRN
}

type storedIRN struct {
	IRN           string
	AckNo         string
	AckDt         string
	DocType       string
	DocNo         string
	DocDate       string
	SupplierGSTIN string
	BuyerGSTIN    string
	TotalValue    float64
	SignedInvoice string
	SignedQRCode  string
	Status        string
	GeneratedAt   time.Time
	CancelledAt   *time.Time
}

func NewMockProvider() *MockProvider {
	return &MockProvider{
		irns: make(map[string]*storedIRN),
		docs: make(map[string]string),
	}
}

func (m *MockProvider) Authenticate(_ context.Context) (*domain.AuthResponse, error) {
	return &domain.AuthResponse{
		AccessToken: "mock-irp-token-" + uuid.New().String()[:8],
		TokenType:   "bearer",
		ExpiresIn:   86399,
		Scope:       "gsp",
	}, nil
}

func (m *MockProvider) GenerateIRN(_ context.Context, req *domain.GenerateIRNRequest) (*domain.GenerateIRNResponse, error) {
	if req.GSTIN == "" {
		return nil, fmt.Errorf("gstin is required")
	}
	if req.DocDtls.No == "" {
		return nil, fmt.Errorf("document number is required")
	}
	if req.DocDtls.Typ == "" {
		return nil, fmt.Errorf("document type is required")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	docKey := fmt.Sprintf("%s|%s|%s", req.DocDtls.Typ, req.DocDtls.No, req.DocDtls.Dt)
	if existingIRN, ok := m.docs[docKey]; ok {
		s := m.irns[existingIRN]
		if s.Status == "ACT" {
			return &domain.GenerateIRNResponse{
				IRN:           s.IRN,
				AckNo:         s.AckNo,
				AckDt:         s.AckDt,
				SignedInvoice: s.SignedInvoice,
				SignedQRCode:  s.SignedQRCode,
				Status:        "ACT",
				GeneratedAt:   s.GeneratedAt,
			}, nil
		}
	}

	now := time.Now()
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s%s%s%d", req.GSTIN, req.DocDtls.No, req.DocDtls.Dt, now.UnixNano())))
	irn := fmt.Sprintf("%x", hash)[:64]
	ackNo := fmt.Sprintf("%d", 100000000000+now.UnixNano()%900000000000)
	ackDt := now.Format("02/01/2006 15:04:05")

	signedInvoice := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`{"irn":"%s","gstin":"%s","doc_no":"%s"}`, irn, req.GSTIN, req.DocDtls.No)))
	signedQR := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("QR:%s:%s:%s:%.2f", irn, req.GSTIN, req.DocDtls.No, req.ValDtls.TotInvVal)))

	s := &storedIRN{
		IRN:           irn,
		AckNo:         ackNo,
		AckDt:         ackDt,
		DocType:       req.DocDtls.Typ,
		DocNo:         req.DocDtls.No,
		DocDate:       req.DocDtls.Dt,
		SupplierGSTIN: req.GSTIN,
		BuyerGSTIN:    req.BuyDtls.GSTIN,
		TotalValue:    req.ValDtls.TotInvVal,
		SignedInvoice: signedInvoice,
		SignedQRCode:  signedQR,
		Status:        "ACT",
		GeneratedAt:   now,
	}

	m.irns[irn] = s
	m.docs[docKey] = irn

	return &domain.GenerateIRNResponse{
		IRN:           irn,
		AckNo:         ackNo,
		AckDt:         ackDt,
		SignedInvoice: signedInvoice,
		SignedQRCode:  signedQR,
		Status:        "ACT",
		GeneratedAt:   now,
	}, nil
}

func (m *MockProvider) CancelIRN(_ context.Context, req *domain.CancelIRNRequest) (*domain.CancelIRNResponse, error) {
	if req.IRN == "" {
		return nil, fmt.Errorf("irn is required")
	}
	if req.CnlRsn == "" {
		return nil, fmt.Errorf("cancellation reason is required")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.irns[req.IRN]
	if !ok {
		return nil, fmt.Errorf("IRN not found: %s", req.IRN)
	}
	if s.Status == "CANC" {
		return nil, fmt.Errorf("IRN already cancelled: %s", req.IRN)
	}

	if time.Since(s.GeneratedAt) > 24*time.Hour {
		return nil, fmt.Errorf("cancellation window expired: IRN generated at %s, more than 24 hours ago", s.GeneratedAt.Format(time.RFC3339))
	}

	now := time.Now()
	s.Status = "CANC"
	s.CancelledAt = &now

	return &domain.CancelIRNResponse{
		IRN:         req.IRN,
		CancelDate:  now.Format("02/01/2006 15:04:05"),
		Status:      "CANC",
		CancelledAt: now,
	}, nil
}

func (m *MockProvider) GetIRNByIRN(_ context.Context, req *domain.GetIRNByIRNRequest) (*domain.GetIRNResponse, error) {
	if req.IRN == "" {
		return nil, fmt.Errorf("irn is required")
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	s, ok := m.irns[req.IRN]
	if !ok {
		return nil, fmt.Errorf("IRN not found: %s", req.IRN)
	}

	return toGetIRNResponse(s), nil
}

func (m *MockProvider) GetIRNByDoc(_ context.Context, req *domain.GetIRNByDocRequest) (*domain.GetIRNResponse, error) {
	if req.DocNum == "" || req.DocType == "" || req.DocDate == "" {
		return nil, fmt.Errorf("doc_type, doc_num, and doc_date are all required")
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	docKey := fmt.Sprintf("%s|%s|%s", strings.ToUpper(req.DocType), req.DocNum, req.DocDate)
	irn, ok := m.docs[docKey]
	if !ok {
		return nil, fmt.Errorf("no IRN found for document %s %s dated %s", req.DocType, req.DocNum, req.DocDate)
	}

	s := m.irns[irn]
	return toGetIRNResponse(s), nil
}

func (m *MockProvider) ValidateGSTIN(_ context.Context, req *domain.GSTINValidateRequest) (*domain.GSTINValidateResponse, error) {
	if req.GSTIN == "" || len(req.GSTIN) != 15 {
		return nil, fmt.Errorf("invalid GSTIN: must be 15 characters")
	}

	stateCode := req.GSTIN[:2]
	panPart := req.GSTIN[2:12]

	return &domain.GSTINValidateResponse{
		GSTIN:      req.GSTIN,
		LegalName:  fmt.Sprintf("Mock Entity %s", panPart),
		TradeName:  fmt.Sprintf("Trade %s", panPart[:5]),
		StateCode:  stateCode,
		Status:     "Active",
		EntityType: "Regular",
	}, nil
}

func toGetIRNResponse(s *storedIRN) *domain.GetIRNResponse {
	return &domain.GetIRNResponse{
		IRN:           s.IRN,
		AckNo:         s.AckNo,
		AckDt:         s.AckDt,
		Status:        s.Status,
		DocType:       s.DocType,
		DocNo:         s.DocNo,
		DocDate:       s.DocDate,
		SupplierGSTIN: s.SupplierGSTIN,
		BuyerGSTIN:    s.BuyerGSTIN,
		TotalValue:    s.TotalValue,
		SignedInvoice: s.SignedInvoice,
		SignedQRCode:  s.SignedQRCode,
		GeneratedAt:   s.GeneratedAt,
		CancelledAt:   s.CancelledAt,
	}
}
