package provider

import (
	"context"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/services/go/apex-gateway-service/internal/domain"
)

func TestMockProvider_FetchVendors_ReturnsAll(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.FetchVendors(context.Background(), &domain.FetchVendorsRequest{
		TenantID:  uuid.New().String(),
		RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Equal(t, 50, resp.Total)
	assert.Len(t, resp.Vendors, 50)
}

func TestMockProvider_FetchVendors_WithLimit(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.FetchVendors(context.Background(), &domain.FetchVendorsRequest{
		TenantID:  uuid.New().String(),
		Limit:     10,
		RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Len(t, resp.Vendors, 10)
	assert.Equal(t, 50, resp.Total)
}

func TestMockProvider_FetchVendors_WithOffset(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.FetchVendors(context.Background(), &domain.FetchVendorsRequest{
		TenantID:  uuid.New().String(),
		Offset:    45,
		RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Len(t, resp.Vendors, 5)
	assert.Equal(t, 50, resp.Total)
	// First returned should be vendor 46
	assert.Equal(t, "VND-046", resp.Vendors[0].ID)
}

func TestMockProvider_FetchVendors_OffsetBeyondTotal(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.FetchVendors(context.Background(), &domain.FetchVendorsRequest{
		TenantID:  uuid.New().String(),
		Offset:    100,
		RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Empty(t, resp.Vendors)
	assert.Equal(t, 50, resp.Total)
}

func TestMockProvider_FetchVendors_LimitAndOffset(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.FetchVendors(context.Background(), &domain.FetchVendorsRequest{
		TenantID:  uuid.New().String(),
		Limit:     5,
		Offset:    10,
		RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Len(t, resp.Vendors, 5)
	assert.Equal(t, "VND-011", resp.Vendors[0].ID)
	assert.Equal(t, "VND-015", resp.Vendors[4].ID)
}

func TestMockProvider_FetchVendors_ValidGSTIN(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.FetchVendors(context.Background(), &domain.FetchVendorsRequest{
		TenantID:  uuid.New().String(),
		RequestID: uuid.New().String(),
	})
	require.NoError(t, err)

	// GSTIN format: 2-digit state code + 10-char PAN + 1 digit + Z + 1 alpha check
	gstinRegex := regexp.MustCompile(`^[0-9]{2}[A-Z]{5}[0-9]{4}[A-Z][0-9]Z[A-Z]$`)
	for _, v := range resp.Vendors {
		assert.Regexp(t, gstinRegex, v.GSTIN, "Invalid GSTIN for vendor %s: %s", v.Name, v.GSTIN)
	}
}

func TestMockProvider_FetchVendors_CategoryDistribution(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.FetchVendors(context.Background(), &domain.FetchVendorsRequest{
		TenantID:  uuid.New().String(),
		RequestID: uuid.New().String(),
	})
	require.NoError(t, err)

	categories := make(map[string]int)
	for _, v := range resp.Vendors {
		categories[v.Category]++
	}

	// Verify we have multiple categories
	assert.GreaterOrEqual(t, len(categories), 4, "Expected at least 4 vendor categories")
	assert.Greater(t, categories["Manufacturer"], 0)
	assert.Greater(t, categories["Service Provider"], 0)
	assert.Greater(t, categories["Trader"], 0)
	assert.Greater(t, categories["Logistics"], 0)
}

func TestMockProvider_FetchVendors_MSMEDistribution(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.FetchVendors(context.Background(), &domain.FetchVendorsRequest{
		TenantID:  uuid.New().String(),
		RequestID: uuid.New().String(),
	})
	require.NoError(t, err)

	msmeCount := 0
	for _, v := range resp.Vendors {
		if v.MSMERegistered {
			msmeCount++
		}
	}

	// Expect ~30% MSME (15 out of 50), allow some tolerance
	assert.GreaterOrEqual(t, msmeCount, 10, "Expected at least 10 MSME vendors")
	assert.LessOrEqual(t, msmeCount, 25, "Expected at most 25 MSME vendors")
}

func TestMockProvider_FetchVendors_SetsTenantID(t *testing.T) {
	p := NewMockProvider()
	tenantID := uuid.New().String()
	resp, err := p.FetchVendors(context.Background(), &domain.FetchVendorsRequest{
		TenantID:  tenantID,
		RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	for _, v := range resp.Vendors {
		assert.Equal(t, tenantID, v.TenantID)
	}
}

func TestMockProvider_FetchVendors_Deterministic(t *testing.T) {
	p1 := NewMockProvider()
	p2 := NewMockProvider()

	reqID := uuid.New().String()
	tenantID := uuid.New().String()

	resp1, err := p1.FetchVendors(context.Background(), &domain.FetchVendorsRequest{TenantID: tenantID, RequestID: reqID})
	require.NoError(t, err)
	resp2, err := p2.FetchVendors(context.Background(), &domain.FetchVendorsRequest{TenantID: tenantID, RequestID: reqID})
	require.NoError(t, err)

	require.Equal(t, len(resp1.Vendors), len(resp2.Vendors))
	for i := range resp1.Vendors {
		assert.Equal(t, resp1.Vendors[i].ID, resp2.Vendors[i].ID)
		assert.Equal(t, resp1.Vendors[i].Name, resp2.Vendors[i].Name)
		assert.Equal(t, resp1.Vendors[i].GSTIN, resp2.Vendors[i].GSTIN)
	}
}

func TestMockProvider_FetchAPInvoices_ReturnsAll(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.FetchAPInvoices(context.Background(), &domain.FetchAPInvoicesRequest{
		TenantID:  uuid.New().String(),
		RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	// Each vendor has 3-9 invoices, 50 vendors total => should be > 150
	assert.Greater(t, resp.Total, 150, "Expected more than 150 invoices total")
}

func TestMockProvider_FetchAPInvoices_FilterByVendorID(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.FetchAPInvoices(context.Background(), &domain.FetchAPInvoicesRequest{
		TenantID:  uuid.New().String(),
		VendorID:  "VND-001",
		RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Greater(t, resp.Total, 0)
	for _, inv := range resp.Invoices {
		assert.Equal(t, "VND-001", inv.VendorID)
	}
}

func TestMockProvider_FetchAPInvoices_FilterByVendorID_NoMatch(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.FetchAPInvoices(context.Background(), &domain.FetchAPInvoicesRequest{
		TenantID:  uuid.New().String(),
		VendorID:  "VND-999",
		RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Equal(t, 0, resp.Total)
	assert.Empty(t, resp.Invoices)
}

func TestMockProvider_FetchAPInvoices_SetsTenantID(t *testing.T) {
	p := NewMockProvider()
	tenantID := uuid.New().String()
	resp, err := p.FetchAPInvoices(context.Background(), &domain.FetchAPInvoicesRequest{
		TenantID:  tenantID,
		VendorID:  "VND-001",
		RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	for _, inv := range resp.Invoices {
		assert.Equal(t, tenantID, inv.TenantID)
	}
}

func TestMockProvider_FetchAPInvoices_InvoiceFields(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.FetchAPInvoices(context.Background(), &domain.FetchAPInvoicesRequest{
		TenantID:  uuid.New().String(),
		VendorID:  "VND-001",
		RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	require.Greater(t, len(resp.Invoices), 0)

	inv := resp.Invoices[0]
	assert.NotEmpty(t, inv.ID)
	assert.NotEmpty(t, inv.VendorID)
	assert.NotEmpty(t, inv.VendorGSTIN)
	assert.NotEmpty(t, inv.InvoiceNumber)
	assert.NotEmpty(t, inv.InvoiceDate)
	assert.NotEmpty(t, inv.DueDate)
	assert.Greater(t, inv.TaxableValue, float64(0))
	assert.Greater(t, inv.TotalAmount, inv.TaxableValue)

	// Verify GST amounts sum correctly
	gstTotal := inv.CGSTAmount + inv.SGSTAmount + inv.IGSTAmount
	assert.InDelta(t, inv.TotalAmount, inv.TaxableValue+gstTotal, 0.01)

	// Either IGST or CGST+SGST should be set, not both
	if inv.IGSTAmount > 0 {
		assert.Equal(t, float64(0), inv.CGSTAmount)
		assert.Equal(t, float64(0), inv.SGSTAmount)
	} else {
		assert.Equal(t, inv.CGSTAmount, inv.SGSTAmount)
	}

	// Payment status should be one of valid values
	validPaymentStatuses := map[string]bool{"paid": true, "unpaid": true, "overdue": true, "partial": true}
	assert.True(t, validPaymentStatuses[inv.PaymentStatus], "Invalid payment status: %s", inv.PaymentStatus)

	// GST filing status should be valid
	validFilingStatuses := map[string]bool{"filed": true, "pending": true, "late": true, "not_filed": true}
	assert.True(t, validFilingStatuses[inv.GSTFilingStatus], "Invalid filing status: %s", inv.GSTFilingStatus)

	// Mismatch status should be valid
	validMismatchStatuses := map[string]bool{"matched": true, "mismatched": true, "pending": true}
	assert.True(t, validMismatchStatuses[inv.MismatchStatus], "Invalid mismatch status: %s", inv.MismatchStatus)
}

func TestMockProvider_FetchAPInvoices_Deterministic(t *testing.T) {
	p1 := NewMockProvider()
	p2 := NewMockProvider()

	tenantID := uuid.New().String()
	reqID := uuid.New().String()

	resp1, err := p1.FetchAPInvoices(context.Background(), &domain.FetchAPInvoicesRequest{
		TenantID: tenantID, VendorID: "VND-005", RequestID: reqID,
	})
	require.NoError(t, err)
	resp2, err := p2.FetchAPInvoices(context.Background(), &domain.FetchAPInvoicesRequest{
		TenantID: tenantID, VendorID: "VND-005", RequestID: reqID,
	})
	require.NoError(t, err)

	require.Equal(t, resp1.Total, resp2.Total)
	for i := range resp1.Invoices {
		assert.Equal(t, resp1.Invoices[i].ID, resp2.Invoices[i].ID)
		assert.Equal(t, resp1.Invoices[i].TaxableValue, resp2.Invoices[i].TaxableValue)
		assert.Equal(t, resp1.Invoices[i].TotalAmount, resp2.Invoices[i].TotalAmount)
	}
}
