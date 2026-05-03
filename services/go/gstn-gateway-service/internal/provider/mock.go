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
	mu              sync.RWMutex
	filings         map[string]*domain.MockFiling         // key: gstin:ret_period
	gstr3bFilings   map[string]*domain.MockGSTR3BFiling   // key: gstin:ret_period
	gstr2bInvoices  map[string][]domain.GSTR2BInvoice     // key: gstin:ret_period
	gstr9Filings    map[string]*domain.MockGSTR9Filing     // key: reference
	gstr9cFilings   map[string]*domain.MockGSTR9CFiling    // key: reference
	requests        map[string]interface{}                 // idempotency: request_id → response
}

func NewMockProvider() *MockProvider {
	return &MockProvider{
		filings:        make(map[string]*domain.MockFiling),
		gstr3bFilings:  make(map[string]*domain.MockGSTR3BFiling),
		gstr2bInvoices: make(map[string][]domain.GSTR2BInvoice),
		gstr9Filings:   make(map[string]*domain.MockGSTR9Filing),
		gstr9cFilings:  make(map[string]*domain.MockGSTR9CFiling),
		requests:       make(map[string]interface{}),
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
	m.gstr3bFilings = make(map[string]*domain.MockGSTR3BFiling)
	m.gstr2bInvoices = make(map[string][]domain.GSTR2BInvoice)
	m.gstr9Filings = make(map[string]*domain.MockGSTR9Filing)
	m.gstr9cFilings = make(map[string]*domain.MockGSTR9CFiling)
	m.requests = make(map[string]interface{})
}

// generateMock2BInvoices creates 20 realistic GSTR-2B inward supply invoices.
func generateMock2BInvoices(retPeriod string) []domain.GSTR2BInvoice {
	suppliers := []struct {
		gstin string
		name  string
	}{
		{"29AABCA1234A1Z5", "Alpha Enterprises"},
		{"07BBBBB2222B1Z8", "Beta Technologies"},
		{"27CCCCC3333C1Z1", "Gamma Solutions"},
		{"33DDDDD4444D1Z4", "Delta Manufacturing"},
		{"06EEEEE5555E1Z7", "Epsilon Trading"},
		{"09FFFFF6666F1Z0", "Zeta Industries"},
		{"19GGGGG7777G1Z3", "Eta Logistics"},
		{"24HHHHH8888H1Z6", "Theta Imports"},
		{"32IIIII9999I1Z9", "Iota Services"},
		{"36JJJJJ0000J1Z2", "Kappa Supplies"},
	}

	hsnCodes := []string{"8471", "8523", "3004", "7308", "8544", "9403", "4819", "3926", "8504", "7210"}
	places := []string{"29", "07", "27", "33", "06", "09", "19", "24", "32", "36"}
	invTypes := []string{"R", "R", "R", "R", "R", "R", "R", "R", "SEWOP", "DE"}

	invoices := make([]domain.GSTR2BInvoice, 0, 20)
	for i := 0; i < 20; i++ {
		sup := suppliers[i%len(suppliers)]
		taxable := float64((i+1)*15000) + float64(i*731)
		cgst := taxable * 0.09
		sgst := taxable * 0.09
		igst := 0.0

		// For inter-state (different place of supply), use IGST instead
		if i%4 == 0 {
			igst = taxable * 0.18
			cgst = 0
			sgst = 0
		}

		day := (i%28) + 1
		invoices = append(invoices, domain.GSTR2BInvoice{
			SupplierGSTIN: sup.gstin,
			InvoiceNumber: fmt.Sprintf("INV-%s-%04d", retPeriod[:2], i+1),
			InvoiceDate:   fmt.Sprintf("%02d/%s/20%s", day, retPeriod[:2], retPeriod[4:]),
			InvoiceType:   invTypes[i%len(invTypes)],
			TaxableValue:  taxable,
			CGSTAmount:    cgst,
			SGSTAmount:    sgst,
			IGSTAmount:    igst,
			TotalValue:    taxable + cgst + sgst + igst,
			PlaceOfSupply: places[i%len(places)],
			ReverseCharge: i == 7 || i == 14,
			HSN:           hsnCodes[i%len(hsnCodes)],
			ITC:           "eligible",
			IMSAction:     "PENDING",
		})
	}
	return invoices
}

func (m *MockProvider) ensureGSTR2BData(gstin, retPeriod string) {
	key := filingKey(gstin, retPeriod)
	if _, ok := m.gstr2bInvoices[key]; !ok {
		m.gstr2bInvoices[key] = generateMock2BInvoices(retPeriod)
	}
}

func (m *MockProvider) GSTR2BGet(_ context.Context, req *domain.GSTR2BGetRequest) (*domain.GSTR2BGetResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ensureGSTR2BData(req.GSTIN, req.RetPeriod)

	key := filingKey(req.GSTIN, req.RetPeriod)
	invoices := m.gstr2bInvoices[key]

	return &domain.GSTR2BGetResponse{
		GSTIN:       req.GSTIN,
		RetPeriod:   req.RetPeriod,
		Invoices:    invoices,
		TotalCount:  len(invoices),
		GeneratedOn: time.Now().UTC().Format("02/01/2006 15:04:05"),
		Status:      "success",
		RequestID:   req.RequestID,
	}, nil
}

func (m *MockProvider) GSTR2AGet(_ context.Context, req *domain.GSTR2AGetRequest) (*domain.GSTR2AGetResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ensureGSTR2BData(req.GSTIN, req.RetPeriod)

	key := filingKey(req.GSTIN, req.RetPeriod)
	allInvoices := m.gstr2bInvoices[key]

	// For GSTR-2A, filter by section. Default to B2B.
	section := req.Section
	if section == "" {
		section = "B2B"
	}

	// The mock returns the same invoices for B2B section; other sections return empty.
	var filtered []domain.GSTR2BInvoice
	if section == "B2B" {
		filtered = allInvoices
	}

	return &domain.GSTR2AGetResponse{
		GSTIN:     req.GSTIN,
		RetPeriod: req.RetPeriod,
		Section:   section,
		Invoices:  filtered,
		Status:    "success",
		RequestID: req.RequestID,
	}, nil
}

func (m *MockProvider) IMSGet(_ context.Context, req *domain.IMSGetRequest) (*domain.IMSGetResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ensureGSTR2BData(req.GSTIN, req.RetPeriod)

	key := filingKey(req.GSTIN, req.RetPeriod)
	b2bInvoices := m.gstr2bInvoices[key]

	imsInvoices := make([]domain.IMSInvoice, 0, len(b2bInvoices))
	var summary domain.IMSSummary

	for i, inv := range b2bInvoices {
		imsInv := domain.IMSInvoice{
			InvoiceID:     fmt.Sprintf("IMS-%s-%04d", req.RetPeriod[:2], i+1),
			SupplierGSTIN: inv.SupplierGSTIN,
			InvoiceNumber: inv.InvoiceNumber,
			InvoiceDate:   inv.InvoiceDate,
			TaxableValue:  inv.TaxableValue,
			TotalValue:    inv.TotalValue,
			CGSTAmount:    inv.CGSTAmount,
			SGSTAmount:    inv.SGSTAmount,
			IGSTAmount:    inv.IGSTAmount,
			Action:        inv.IMSAction,
		}

		switch inv.IMSAction {
		case "ACCEPT":
			summary.Accepted++
			summary.AcceptedValue += inv.TotalValue
		case "REJECT":
			summary.Rejected++
			summary.RejectedValue += inv.TotalValue
		default:
			summary.Pending++
			summary.PendingValue += inv.TotalValue
		}

		imsInvoices = append(imsInvoices, imsInv)
	}

	return &domain.IMSGetResponse{
		GSTIN:      req.GSTIN,
		RetPeriod:  req.RetPeriod,
		Invoices:   imsInvoices,
		TotalCount: len(imsInvoices),
		Summary:    summary,
		Status:     "success",
		RequestID:  req.RequestID,
	}, nil
}

func (m *MockProvider) IMSAction(_ context.Context, req *domain.IMSActionRequest) (*domain.IMSActionResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if resp, ok := m.requests[req.RequestID]; ok {
		return resp.(*domain.IMSActionResponse), nil
	}

	if req.Action != "ACCEPT" && req.Action != "REJECT" && req.Action != "PENDING" {
		return nil, fmt.Errorf("invalid action: %s (must be ACCEPT, REJECT, or PENDING)", req.Action)
	}

	m.ensureGSTR2BData(req.GSTIN, req.RetPeriod)

	key := filingKey(req.GSTIN, req.RetPeriod)
	invoices := m.gstr2bInvoices[key]

	found := false
	for i := range invoices {
		invID := fmt.Sprintf("IMS-%s-%04d", req.RetPeriod[:2], i+1)
		if invID == req.InvoiceID {
			invoices[i].IMSAction = req.Action
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("invoice %s not found", req.InvoiceID)
	}

	m.gstr2bInvoices[key] = invoices

	now := time.Now().UTC()
	resp := &domain.IMSActionResponse{
		InvoiceID: req.InvoiceID,
		Action:    req.Action,
		Status:    "success",
		RequestID: req.RequestID,
		UpdatedAt: now.Format(time.RFC3339),
	}
	m.requests[req.RequestID] = resp
	return resp, nil
}

func (m *MockProvider) IMSBulkAction(_ context.Context, req *domain.IMSBulkActionRequest) (*domain.IMSBulkActionResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if resp, ok := m.requests[req.RequestID]; ok {
		return resp.(*domain.IMSBulkActionResponse), nil
	}

	if req.Action != "ACCEPT" && req.Action != "REJECT" && req.Action != "PENDING" {
		return nil, fmt.Errorf("invalid action: %s (must be ACCEPT, REJECT, or PENDING)", req.Action)
	}

	m.ensureGSTR2BData(req.GSTIN, req.RetPeriod)

	key := filingKey(req.GSTIN, req.RetPeriod)
	invoices := m.gstr2bInvoices[key]

	idSet := make(map[string]struct{}, len(req.InvoiceIDs))
	for _, id := range req.InvoiceIDs {
		idSet[id] = struct{}{}
	}

	updated := 0
	for i := range invoices {
		invID := fmt.Sprintf("IMS-%s-%04d", req.RetPeriod[:2], i+1)
		if _, ok := idSet[invID]; ok {
			invoices[i].IMSAction = req.Action
			updated++
		}
	}

	m.gstr2bInvoices[key] = invoices

	resp := &domain.IMSBulkActionResponse{
		OperationID:   uuid.New().String(),
		TotalInvoices: updated,
		Status:        "success",
		RequestID:     req.RequestID,
	}
	m.requests[req.RequestID] = resp
	return resp, nil
}

func (m *MockProvider) GSTR3BSave(_ context.Context, req *domain.GSTR3BSaveRequest) (*domain.GSTR3BSaveResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if resp, ok := m.requests[req.RequestID]; ok {
		return resp.(*domain.GSTR3BSaveResponse), nil
	}

	key := filingKey(req.GSTIN, req.RetPeriod)
	f, ok := m.gstr3bFilings[key]
	if !ok {
		f = &domain.MockGSTR3BFiling{
			GSTIN:     req.GSTIN,
			RetPeriod: req.RetPeriod,
			Status:    domain.StatusDraft,
		}
		m.gstr3bFilings[key] = f
	}

	if f.Status == domain.StatusSubmitted || f.Status == domain.StatusFiled {
		return nil, fmt.Errorf("cannot save: GSTR-3B is %s", f.Status)
	}

	f.Data = req.Data
	f.Status = domain.StatusSaved

	resp := &domain.GSTR3BSaveResponse{
		Status:    "success",
		RequestID: req.RequestID,
		Message:   "GSTR-3B data saved successfully",
	}
	m.requests[req.RequestID] = resp
	return resp, nil
}

func (m *MockProvider) GSTR3BSubmit(_ context.Context, req *domain.GSTR3BSubmitRequest) (*domain.GSTR3BSubmitResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if resp, ok := m.requests[req.RequestID]; ok {
		return resp.(*domain.GSTR3BSubmitResponse), nil
	}

	key := filingKey(req.GSTIN, req.RetPeriod)
	f, ok := m.gstr3bFilings[key]
	if !ok {
		return nil, fmt.Errorf("no GSTR-3B draft found for %s/%s", req.GSTIN, req.RetPeriod)
	}

	if f.Status == domain.StatusFiled {
		return nil, fmt.Errorf("GSTR-3B already filed")
	}
	if f.Status == domain.StatusSubmitted {
		return nil, fmt.Errorf("GSTR-3B already submitted")
	}
	if f.Data == nil {
		return nil, fmt.Errorf("no GSTR-3B data saved")
	}

	f.Status = domain.StatusSubmitted

	resp := &domain.GSTR3BSubmitResponse{
		Status:    "success",
		RequestID: req.RequestID,
		Message:   "GSTR-3B submitted successfully. Locked for filing.",
	}
	m.requests[req.RequestID] = resp
	return resp, nil
}

func (m *MockProvider) GSTR3BFile(_ context.Context, req *domain.GSTR3BFileRequest) (*domain.GSTR3BFileResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if resp, ok := m.requests[req.RequestID]; ok {
		return resp.(*domain.GSTR3BFileResponse), nil
	}

	key := filingKey(req.GSTIN, req.RetPeriod)
	f, ok := m.gstr3bFilings[key]
	if !ok {
		return nil, fmt.Errorf("no GSTR-3B filing found for %s/%s", req.GSTIN, req.RetPeriod)
	}

	if f.Status == domain.StatusFiled {
		return nil, fmt.Errorf("GSTR-3B already filed")
	}
	if f.Status != domain.StatusSubmitted {
		return nil, fmt.Errorf("must submit GSTR-3B before filing (current status: %s)", f.Status)
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
	f.ARN = fmt.Sprintf("AB%s%s%s", gstinPrefix, req.RetPeriod, uuid.New().String()[:8])
	f.FiledAt = &now

	resp := &domain.GSTR3BFileResponse{
		Status:    "success",
		ARN:       f.ARN,
		RequestID: req.RequestID,
		Message:   "GSTR-3B filed successfully",
		FiledAt:   now,
	}
	m.requests[req.RequestID] = resp
	return resp, nil
}

func (m *MockProvider) GSTR1Summary(_ context.Context, req *domain.GSTR1SummaryRequest) (*domain.GSTR1SummaryResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := filingKey(req.GSTIN, req.RetPeriod)
	f, ok := m.filings[key]
	if !ok {
		return &domain.GSTR1SummaryResponse{
			GSTIN:     req.GSTIN,
			RetPeriod: req.RetPeriod,
			Summary:   map[string]interface{}{},
			Status:    string(domain.StatusDraft),
			RequestID: req.RequestID,
		}, nil
	}

	summary := make(map[string]interface{})
	totalTaxable := 0.0
	totalTax := 0.0
	sectionCount := len(f.Sections)

	for section := range f.Sections {
		summary[section] = map[string]interface{}{
			"status": "saved",
		}
	}

	summary["total_sections"] = sectionCount
	summary["total_taxable_value"] = totalTaxable
	summary["total_tax"] = totalTax

	return &domain.GSTR1SummaryResponse{
		GSTIN:     req.GSTIN,
		RetPeriod: req.RetPeriod,
		Summary:   summary,
		Status:    string(f.Status),
		RequestID: req.RequestID,
	}, nil
}

func gstr9Ref(gstin, fy string) string {
	return fmt.Sprintf("GSTR9-%s-%s-%s", fy, gstin[:2], uuid.New().String()[:6])
}

func gstr9cRef(gstin, fy string) string {
	return fmt.Sprintf("GSTR9C-%s-%s-%s", fy, gstin[:2], uuid.New().String()[:6])
}

func (m *MockProvider) gstr9ByGSTINAndFY(gstin, fy string) *domain.MockGSTR9Filing {
	for _, f := range m.gstr9Filings {
		if f.GSTIN == gstin && f.FinancialYear == fy {
			return f
		}
	}
	return nil
}

func (m *MockProvider) gstr9cByGSTINAndFY(gstin, fy string) *domain.MockGSTR9CFiling {
	for _, f := range m.gstr9cFilings {
		if f.GSTIN == gstin && f.FinancialYear == fy {
			return f
		}
	}
	return nil
}

func (m *MockProvider) GSTR9Save(_ context.Context, req *domain.GSTR9SaveRequest) (*domain.GSTR9SaveResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if resp, ok := m.requests[req.RequestID]; ok {
		return resp.(*domain.GSTR9SaveResponse), nil
	}

	if len(req.GSTIN) < 15 {
		return nil, fmt.Errorf("invalid GSTIN format")
	}
	if len(req.FinancialYear) != 7 || req.FinancialYear[4] != '-' {
		return nil, fmt.Errorf("invalid financial_year format (expected YYYY-YY)")
	}

	existing := m.gstr9ByGSTINAndFY(req.GSTIN, req.FinancialYear)
	if existing != nil {
		if existing.Status == domain.StatusSubmitted || existing.Status == domain.StatusFiled {
			return nil, fmt.Errorf("cannot save: GSTR-9 is %s", existing.Status)
		}
		now := time.Now().UTC()
		if !existing.SavedAt.IsZero() && now.Sub(existing.SavedAt) > 24*time.Hour {
			return nil, fmt.Errorf("24-hour edit window expired; re-initiate the return")
		}
		existing.Data = req.Data
		existing.SavedAt = now
		resp := &domain.GSTR9SaveResponse{
			Status:    "success",
			Reference: existing.Reference,
			RequestID: req.RequestID,
			Message:   "GSTR-9 draft updated",
			SavedAt:   now,
		}
		m.requests[req.RequestID] = resp
		return resp, nil
	}

	now := time.Now().UTC()
	ref := gstr9Ref(req.GSTIN, req.FinancialYear)
	f := &domain.MockGSTR9Filing{
		GSTIN:         req.GSTIN,
		FinancialYear: req.FinancialYear,
		Status:        domain.StatusSaved,
		Reference:     ref,
		Data:          req.Data,
		SavedAt:       now,
	}
	m.gstr9Filings[ref] = f

	resp := &domain.GSTR9SaveResponse{
		Status:    "success",
		Reference: ref,
		RequestID: req.RequestID,
		Message:   "GSTR-9 draft saved",
		SavedAt:   now,
	}
	m.requests[req.RequestID] = resp
	return resp, nil
}

func (m *MockProvider) GSTR9Submit(_ context.Context, req *domain.GSTR9SubmitRequest) (*domain.GSTR9SubmitResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if resp, ok := m.requests[req.RequestID]; ok {
		return resp.(*domain.GSTR9SubmitResponse), nil
	}

	f, ok := m.gstr9Filings[req.Reference]
	if !ok {
		return nil, fmt.Errorf("no GSTR-9 draft found for reference %s", req.Reference)
	}
	if f.Status == domain.StatusFiled {
		return nil, fmt.Errorf("GSTR-9 already filed")
	}
	if f.Status == domain.StatusSubmitted {
		return nil, fmt.Errorf("GSTR-9 already submitted")
	}
	if f.Data == nil {
		return nil, fmt.Errorf("no GSTR-9 data saved")
	}

	f.Status = domain.StatusSubmitted

	resp := &domain.GSTR9SubmitResponse{
		Status:    "success",
		Reference: req.Reference,
		RequestID: req.RequestID,
		Message:   "GSTR-9 submitted. Locked for filing.",
	}
	m.requests[req.RequestID] = resp
	return resp, nil
}

func (m *MockProvider) GSTR9File(_ context.Context, req *domain.GSTR9FileRequest) (*domain.GSTR9FileResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if resp, ok := m.requests[req.RequestID]; ok {
		return resp.(*domain.GSTR9FileResponse), nil
	}

	f, ok := m.gstr9Filings[req.Reference]
	if !ok {
		return nil, fmt.Errorf("no GSTR-9 found for reference %s", req.Reference)
	}
	if f.Status == domain.StatusFiled {
		return nil, fmt.Errorf("GSTR-9 already filed")
	}
	if f.Status != domain.StatusSubmitted {
		return nil, fmt.Errorf("must submit GSTR-9 before filing (current status: %s)", f.Status)
	}

	if req.SignType != "DSC" && req.SignType != "EVC" {
		return nil, fmt.Errorf("invalid sign_type: %s (must be DSC or EVC)", req.SignType)
	}
	if req.SignType == "EVC" && req.EVOTP == "" {
		return nil, fmt.Errorf("EVC OTP required for EVC signing")
	}

	now := time.Now().UTC()
	f.Status = domain.StatusFiled
	f.ARN = fmt.Sprintf("AR%s%s%s", req.GSTIN[:2], req.FinancialYear[:4], uuid.New().String()[:8])
	f.FiledAt = &now

	resp := &domain.GSTR9FileResponse{
		Status:    "success",
		ARN:       f.ARN,
		RequestID: req.RequestID,
		Message:   "GSTR-9 filed successfully",
		FiledAt:   now,
	}
	m.requests[req.RequestID] = resp
	return resp, nil
}

func (m *MockProvider) GSTR9Status(_ context.Context, req *domain.GSTR9StatusRequest) (*domain.GSTR9StatusResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	f, ok := m.gstr9Filings[req.Reference]
	if !ok {
		return nil, fmt.Errorf("no GSTR-9 found for reference %s", req.Reference)
	}

	return &domain.GSTR9StatusResponse{
		Reference: req.Reference,
		Status:    string(f.Status),
		ARN:       f.ARN,
		FiledAt:   f.FiledAt,
		RequestID: req.RequestID,
	}, nil
}

func (m *MockProvider) GSTR9CSave(_ context.Context, req *domain.GSTR9CSaveRequest) (*domain.GSTR9CSaveResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if resp, ok := m.requests[req.RequestID]; ok {
		return resp.(*domain.GSTR9CSaveResponse), nil
	}

	if len(req.GSTIN) < 15 {
		return nil, fmt.Errorf("invalid GSTIN format")
	}
	if len(req.FinancialYear) != 7 || req.FinancialYear[4] != '-' {
		return nil, fmt.Errorf("invalid financial_year format (expected YYYY-YY)")
	}

	existing := m.gstr9cByGSTINAndFY(req.GSTIN, req.FinancialYear)
	if existing != nil {
		if existing.Status == domain.StatusFiled {
			return nil, fmt.Errorf("cannot save: GSTR-9C is already filed")
		}
		now := time.Now().UTC()
		if !existing.SavedAt.IsZero() && now.Sub(existing.SavedAt) > 24*time.Hour {
			return nil, fmt.Errorf("24-hour edit window expired; re-initiate the reconciliation")
		}
		existing.Data = req.Data
		existing.SavedAt = now
		resp := &domain.GSTR9CSaveResponse{
			Status:    "success",
			Reference: existing.Reference,
			RequestID: req.RequestID,
			Message:   "GSTR-9C draft updated",
			SavedAt:   now,
		}
		m.requests[req.RequestID] = resp
		return resp, nil
	}

	now := time.Now().UTC()
	ref := gstr9cRef(req.GSTIN, req.FinancialYear)
	f := &domain.MockGSTR9CFiling{
		GSTIN:         req.GSTIN,
		FinancialYear: req.FinancialYear,
		Status:        domain.StatusSaved,
		Reference:     ref,
		Data:          req.Data,
		SavedAt:       now,
	}
	m.gstr9cFilings[ref] = f

	resp := &domain.GSTR9CSaveResponse{
		Status:    "success",
		Reference: ref,
		RequestID: req.RequestID,
		Message:   "GSTR-9C draft saved",
		SavedAt:   now,
	}
	m.requests[req.RequestID] = resp
	return resp, nil
}

func (m *MockProvider) GSTR9CFile(_ context.Context, req *domain.GSTR9CFileRequest) (*domain.GSTR9CFileResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if resp, ok := m.requests[req.RequestID]; ok {
		return resp.(*domain.GSTR9CFileResponse), nil
	}

	f, ok := m.gstr9cFilings[req.Reference]
	if !ok {
		return nil, fmt.Errorf("no GSTR-9C found for reference %s", req.Reference)
	}
	if f.Status == domain.StatusFiled {
		return nil, fmt.Errorf("GSTR-9C already filed")
	}
	if f.Data == nil {
		return nil, fmt.Errorf("no GSTR-9C data saved")
	}

	now := time.Now().UTC()
	f.Status = domain.StatusFiled
	f.ARN = fmt.Sprintf("AC%s%s%s", req.GSTIN[:2], req.FinancialYear[:4], uuid.New().String()[:8])
	f.FiledAt = &now

	resp := &domain.GSTR9CFileResponse{
		Status:    "success",
		ARN:       f.ARN,
		RequestID: req.RequestID,
		Message:   "GSTR-9C filed successfully with DSC",
		FiledAt:   now,
	}
	m.requests[req.RequestID] = resp
	return resp, nil
}

func (m *MockProvider) GSTR9CStatus(_ context.Context, req *domain.GSTR9CStatusRequest) (*domain.GSTR9CStatusResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	f, ok := m.gstr9cFilings[req.Reference]
	if !ok {
		return nil, fmt.Errorf("no GSTR-9C found for reference %s", req.Reference)
	}

	return &domain.GSTR9CStatusResponse{
		Reference: req.Reference,
		Status:    string(f.Status),
		ARN:       f.ARN,
		FiledAt:   f.FiledAt,
		RequestID: req.RequestID,
	}, nil
}
