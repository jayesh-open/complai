package categorizer

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/complai/complai/services/go/gst-service/internal/domain"
)

func entry(docType, supplyType, buyerGSTIN string, cgst, sgst, igst int64) *domain.SalesRegisterEntry {
	return &domain.SalesRegisterEntry{
		DocumentType: docType,
		SupplyType:   supplyType,
		BuyerGSTIN:   buyerGSTIN,
		TaxableValue: decimal.NewFromInt(100000),
		CGSTRate:     decimal.NewFromInt(cgst),
		SGSTRate:     decimal.NewFromInt(sgst),
		IGSTRate:     decimal.NewFromInt(igst),
	}
}

func TestCategorize_B2B(t *testing.T) {
	e := entry("INV", "B2B", "29AABCB0001B1Z5", 9, 9, 0)
	assert.Equal(t, domain.SectionB2B, Categorize(e))
}

func TestCategorize_B2B_InterState(t *testing.T) {
	e := entry("INV", "B2B", "27AABCB0001B1Z5", 0, 0, 18)
	assert.Equal(t, domain.SectionB2B, Categorize(e))
}

func TestCategorize_B2CL(t *testing.T) {
	e := entry("INV", "B2CL", "URP", 0, 0, 18)
	assert.Equal(t, domain.SectionB2CL, Categorize(e))
}

func TestCategorize_B2CS(t *testing.T) {
	e := entry("INV", "B2CS", "URP", 9, 9, 0)
	assert.Equal(t, domain.SectionB2CS, Categorize(e))
}

func TestCategorize_B2CS_NoGSTIN(t *testing.T) {
	e := entry("INV", "B2CS", "", 9, 9, 0)
	assert.Equal(t, domain.SectionB2CS, Categorize(e))
}

func TestCategorize_Export(t *testing.T) {
	e := entry("INV", "EXP", "", 0, 0, 18)
	assert.Equal(t, domain.SectionEXP, Categorize(e))
}

func TestCategorize_CDNR_CreditNote(t *testing.T) {
	e := entry("CRN", "B2B", "29AABCB0001B1Z5", 9, 9, 0)
	assert.Equal(t, domain.SectionCDNR, Categorize(e))
}

func TestCategorize_CDNR_DebitNote(t *testing.T) {
	e := entry("DBN", "B2B", "29AABCB0001B1Z5", 0, 0, 18)
	assert.Equal(t, domain.SectionCDNR, Categorize(e))
}

func TestCategorize_CDNUR(t *testing.T) {
	e := entry("CRN", "B2CS", "URP", 9, 9, 0)
	assert.Equal(t, domain.SectionCDNUR, Categorize(e))
}

func TestCategorize_CDNUR_Empty(t *testing.T) {
	e := entry("CRN", "B2CS", "", 9, 9, 0)
	assert.Equal(t, domain.SectionCDNUR, Categorize(e))
}

func TestCategorize_NIL(t *testing.T) {
	e := entry("INV", "B2B", "29AABCB0001B1Z5", 0, 0, 0)
	assert.Equal(t, domain.SectionNIL, Categorize(e))
}

func TestCategorize_NIL_B2C(t *testing.T) {
	e := entry("INV", "B2CS", "URP", 0, 0, 0)
	assert.Equal(t, domain.SectionNIL, Categorize(e))
}

func TestCategorize_RCM_StillB2B(t *testing.T) {
	e := entry("INV", "B2B", "27AABCD0001D1Z5", 0, 0, 18)
	e.ReverseCharge = true
	assert.Equal(t, domain.SectionB2B, Categorize(e))
}

func TestCategorize_FallbackB2CS(t *testing.T) {
	e := entry("INV", "OTHER", "", 9, 9, 0)
	assert.Equal(t, domain.SectionB2CS, Categorize(e))
}
