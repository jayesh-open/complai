package provider

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockProvider_ListARInvoices(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.ListARInvoices(context.Background(), uuid.New(), "29AABCA1234A1Z5", "042026")
	require.NoError(t, err)
	assert.Equal(t, 100, resp.TotalCount)
	assert.Len(t, resp.Invoices, 100)
}

func TestMockProvider_InvoiceTypeMix(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.ListARInvoices(context.Background(), uuid.New(), "29AABCA1234A1Z5", "042026")
	require.NoError(t, err)

	s := resp.Summary
	assert.Equal(t, 35, s.B2BIntraCount)
	assert.Equal(t, 20, s.B2BInterCount)
	assert.Equal(t, 15, s.B2CSCount)
	assert.Equal(t, 5, s.B2CLCount)
	assert.Equal(t, 5, s.ExportCount)
	assert.Equal(t, 5, s.RCMCount)
	assert.Equal(t, 10, s.CreditNote)
	assert.Equal(t, 5, s.DebitNote)
}

func TestMockProvider_InterStateTax(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.ListARInvoices(context.Background(), uuid.New(), "29AABCA1234A1Z5", "042026")
	require.NoError(t, err)

	for _, inv := range resp.Invoices {
		if inv.PlaceOfSupply != "29" && inv.SupplyType != "EXP" && inv.DocumentType == "INV" {
			assert.True(t, inv.Totals.IGST.IsPositive(), "inter-state should have IGST")
			assert.True(t, inv.Totals.CGST.IsZero(), "inter-state should have zero CGST")
		}
	}
}

func TestMockProvider_IntraStateTax(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.ListARInvoices(context.Background(), uuid.New(), "29AABCA1234A1Z5", "042026")
	require.NoError(t, err)

	for _, inv := range resp.Invoices {
		li := inv.LineItems[0]
		if inv.PlaceOfSupply == "29" && inv.DocumentType == "INV" && li.HSN != "0101" {
			assert.True(t, inv.Totals.CGST.IsPositive(), "intra-state should have CGST")
			assert.True(t, inv.Totals.SGST.IsPositive(), "intra-state should have SGST")
			assert.True(t, inv.Totals.IGST.IsZero(), "intra-state should have zero IGST")
		}
	}
}

func TestMockProvider_NILRatedZeroTax(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.ListARInvoices(context.Background(), uuid.New(), "29AABCA1234A1Z5", "042026")
	require.NoError(t, err)

	var nilCount int
	for _, inv := range resp.Invoices {
		if len(inv.LineItems) > 0 && inv.LineItems[0].HSN == "0101" {
			nilCount++
			assert.True(t, inv.Totals.CGST.IsZero())
			assert.True(t, inv.Totals.SGST.IsZero())
			assert.True(t, inv.Totals.IGST.IsZero())
			assert.True(t, inv.Totals.TaxableValue.Equal(inv.Totals.GrandTotal), "grand total should equal taxable for NIL-rated")
		}
	}
	assert.Equal(t, 5, nilCount)
}

func TestMockProvider_DifferentSupplierState(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.ListARInvoices(context.Background(), uuid.New(), "27AABCA1234A1Z5", "042026")
	require.NoError(t, err)
	assert.Equal(t, 100, resp.TotalCount)
	assert.Equal(t, "27AABCA1234A1Z5", resp.Invoices[0].Supplier.GSTIN)
	assert.Equal(t, "27", resp.Invoices[0].Supplier.StateCode)
}

func TestMockProvider_ExportInvoices(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.ListARInvoices(context.Background(), uuid.New(), "29AABCA1234A1Z5", "042026")
	require.NoError(t, err)

	var expCount int
	for _, inv := range resp.Invoices {
		if inv.SupplyType == "EXP" {
			expCount++
			assert.Equal(t, "96", inv.PlaceOfSupply)
			assert.Empty(t, inv.Buyer.GSTIN)
		}
	}
	assert.Equal(t, 5, expCount)
}

func TestMockProvider_SourceSystem(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.ListARInvoices(context.Background(), uuid.New(), "29AABCA1234A1Z5", "042026")
	require.NoError(t, err)

	for _, inv := range resp.Invoices {
		assert.Equal(t, "aura", inv.SourceSystem)
	}
}
