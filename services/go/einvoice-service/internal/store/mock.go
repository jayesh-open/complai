package store

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/complai/complai/services/go/einvoice-service/internal/domain"
)

type MockStore struct {
	mu        sync.RWMutex
	invoices  map[uuid.UUID]*domain.EInvoice
	lineItems map[uuid.UUID][]domain.EInvoiceLineItem
}

func NewMockStore() *MockStore {
	return &MockStore{
		invoices:  make(map[uuid.UUID]*domain.EInvoice),
		lineItems: make(map[uuid.UUID][]domain.EInvoiceLineItem),
	}
}

func (m *MockStore) CreateEInvoice(_ context.Context, tenantID uuid.UUID, inv *domain.EInvoice) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	inv.ID = uuid.New()
	inv.TenantID = tenantID
	inv.RequestID = uuid.New()
	inv.Status = domain.IRNStatusPending
	now := time.Now()
	inv.CreatedAt = now
	inv.UpdatedAt = now

	copy := *inv
	m.invoices[inv.ID] = &copy
	return nil
}

func (m *MockStore) GetEInvoice(_ context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EInvoice, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	inv, ok := m.invoices[id]
	if !ok || inv.TenantID != tenantID {
		return nil, fmt.Errorf("einvoice not found: %s", id)
	}
	copy := *inv
	return &copy, nil
}

func (m *MockStore) GetEInvoiceByIRN(_ context.Context, tenantID uuid.UUID, irn string) (*domain.EInvoice, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, inv := range m.invoices {
		if inv.IRN == irn && inv.TenantID == tenantID {
			copy := *inv
			return &copy, nil
		}
	}
	return nil, fmt.Errorf("einvoice not found for IRN: %s", irn)
}

func (m *MockStore) GetEInvoiceByInvoiceNumber(_ context.Context, tenantID uuid.UUID, gstin, invoiceNumber string) (*domain.EInvoice, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, inv := range m.invoices {
		if inv.SupplierGSTIN == gstin && inv.InvoiceNumber == invoiceNumber && inv.TenantID == tenantID {
			copy := *inv
			return &copy, nil
		}
	}
	return nil, fmt.Errorf("einvoice not found for invoice number: %s", invoiceNumber)
}

func (m *MockStore) ListEInvoices(_ context.Context, tenantID uuid.UUID, req *domain.ListEInvoicesRequest) ([]domain.EInvoice, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var all []domain.EInvoice
	for _, inv := range m.invoices {
		if inv.TenantID == tenantID && inv.SupplierGSTIN == req.GSTIN {
			if req.Status != "" && string(inv.Status) != req.Status {
				continue
			}
			all = append(all, *inv)
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

func (m *MockStore) UpdateIRNGenerated(_ context.Context, tenantID uuid.UUID, id uuid.UUID, irn, ackNo, signedInvoice, signedQR string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	inv, ok := m.invoices[id]
	if !ok || inv.TenantID != tenantID {
		return fmt.Errorf("einvoice not found: %s", id)
	}
	now := time.Now()
	inv.IRN = irn
	inv.AckNo = ackNo
	inv.SignedInvoice = signedInvoice
	inv.SignedQRCode = signedQR
	inv.Status = domain.IRNStatusGenerated
	inv.IRNGeneratedAt = &now
	inv.UpdatedAt = now
	return nil
}

func (m *MockStore) UpdateIRNCancelled(_ context.Context, tenantID uuid.UUID, id uuid.UUID, reason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	inv, ok := m.invoices[id]
	if !ok || inv.TenantID != tenantID {
		return fmt.Errorf("einvoice not found: %s", id)
	}
	now := time.Now()
	inv.Status = domain.IRNStatusCancelled
	inv.IRNCancelledAt = &now
	inv.CancelReason = reason
	inv.UpdatedAt = now
	return nil
}

func (m *MockStore) UpdateIRNFailed(_ context.Context, tenantID uuid.UUID, id uuid.UUID, reason string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	inv, ok := m.invoices[id]
	if !ok || inv.TenantID != tenantID {
		return fmt.Errorf("einvoice not found: %s", id)
	}
	inv.Status = domain.IRNStatusFailed
	inv.CancelReason = reason
	inv.UpdatedAt = time.Now()
	return nil
}

func (m *MockStore) CreateLineItems(_ context.Context, tenantID uuid.UUID, items []domain.EInvoiceLineItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(items) == 0 {
		return nil
	}
	invoiceID := items[0].InvoiceID
	for i := range items {
		items[i].ID = uuid.New()
		items[i].TenantID = tenantID
		items[i].LineNumber = i + 1
		items[i].CreatedAt = time.Now()
	}
	m.lineItems[invoiceID] = append(m.lineItems[invoiceID], items...)
	return nil
}

func (m *MockStore) GetLineItems(_ context.Context, tenantID uuid.UUID, invoiceID uuid.UUID) ([]domain.EInvoiceLineItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	items := m.lineItems[invoiceID]
	result := make([]domain.EInvoiceLineItem, 0, len(items))
	for _, item := range items {
		if item.TenantID == tenantID {
			result = append(result, item)
		}
	}
	return result, nil
}

func (m *MockStore) GetSummary(_ context.Context, tenantID uuid.UUID, gstin string) (*domain.EInvoiceSummary, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	summary := &domain.EInvoiceSummary{TotalValue: decimal.Zero}
	for _, inv := range m.invoices {
		if inv.TenantID != tenantID || inv.SupplierGSTIN != gstin {
			continue
		}
		summary.TotalCount++
		summary.TotalValue = summary.TotalValue.Add(inv.TotalAmount)
		switch inv.Status {
		case domain.IRNStatusGenerated:
			summary.GeneratedCount++
		case domain.IRNStatusPending:
			summary.PendingCount++
		case domain.IRNStatusCancelled:
			summary.CancelledCount++
		case domain.IRNStatusFailed:
			summary.FailedCount++
		}
	}
	return summary, nil
}

func (m *MockStore) SetIRNGeneratedAt(id uuid.UUID, t *time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if inv, ok := m.invoices[id]; ok {
		inv.IRNGeneratedAt = t
	}
}

var _ Repository = (*MockStore)(nil)
