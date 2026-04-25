package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/complai/complai/services/go/aura-gateway-service/internal/domain"
)

var _ AuraProvider = (*MockProvider)(nil)

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

type invoiceTemplate struct {
	docType       string
	supplyType    string
	reverseCharge bool
	buyerGSTIN    string
	buyerState    string
	hsn           string
	taxableValue  int64
}

var stateNames = map[string]string{
	"29": "Karnataka", "27": "Maharashtra", "07": "Delhi",
	"33": "Tamil Nadu", "36": "Telangana", "06": "Haryana",
	"09": "Uttar Pradesh", "21": "Odisha", "32": "Kerala",
	"24": "Gujarat",
}

func (m *MockProvider) ListARInvoices(_ context.Context, tenantID uuid.UUID, gstin, _ string) (*domain.InvoiceListResponse, error) {
	supplierState := "29" // Karnataka default
	if len(gstin) >= 2 {
		supplierState = gstin[:2]
	}

	templates := buildTemplates(supplierState)
	invoices := make([]domain.Invoice, 0, 100)
	summary := domain.InvoiceSummary{}
	baseDate := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)

	for i, tmpl := range templates {
		inv := buildInvoice(tenantID, gstin, supplierState, tmpl, i, baseDate)
		invoices = append(invoices, inv)
		updateSummary(&summary, tmpl, supplierState)
	}

	return &domain.InvoiceListResponse{
		Invoices:   invoices,
		TotalCount: len(invoices),
		Summary:    summary,
	}, nil
}

func buildTemplates(supplierState string) []invoiceTemplate {
	interStates := []string{"27", "07", "33", "36", "06"}
	if supplierState == "27" {
		interStates = []string{"29", "07", "33", "36", "06"}
	}

	templates := make([]invoiceTemplate, 0, 100)

	// 30 B2B intra-state
	for i := 0; i < 30; i++ {
		templates = append(templates, invoiceTemplate{
			docType: "INV", supplyType: "B2B", buyerGSTIN: fmt.Sprintf("%sAABCB%04dB1Z5", supplierState, i),
			buyerState: supplierState, hsn: hsnCodes[i%len(hsnCodes)], taxableValue: int64(50000 + i*1000),
		})
	}

	// 20 B2B inter-state
	for i := 0; i < 20; i++ {
		st := interStates[i%len(interStates)]
		templates = append(templates, invoiceTemplate{
			docType: "INV", supplyType: "B2B", buyerGSTIN: fmt.Sprintf("%sAABCC%04dC1Z5", st, i),
			buyerState: st, hsn: hsnCodes[(i+3)%len(hsnCodes)], taxableValue: int64(80000 + i*2000),
		})
	}

	// 15 B2CS (unregistered, ≤2 lakh)
	for i := 0; i < 15; i++ {
		templates = append(templates, invoiceTemplate{
			docType: "INV", supplyType: "B2CS", buyerGSTIN: "URP",
			buyerState: supplierState, hsn: hsnCodes[(i+1)%len(hsnCodes)], taxableValue: int64(10000 + i*5000),
		})
	}

	// 5 B2CL (unregistered, >2 lakh, inter-state)
	for i := 0; i < 5; i++ {
		st := interStates[i%len(interStates)]
		templates = append(templates, invoiceTemplate{
			docType: "INV", supplyType: "B2CL", buyerGSTIN: "URP",
			buyerState: st, hsn: hsnCodes[(i+2)%len(hsnCodes)], taxableValue: int64(250000 + i*50000),
		})
	}

	// 5 Exports
	for i := 0; i < 5; i++ {
		templates = append(templates, invoiceTemplate{
			docType: "INV", supplyType: "EXP", buyerGSTIN: "",
			buyerState: "96", hsn: hsnCodes[(i+4)%len(hsnCodes)], taxableValue: int64(500000 + i*100000),
		})
	}

	// 5 RCM (reverse charge)
	for i := 0; i < 5; i++ {
		st := interStates[i%len(interStates)]
		templates = append(templates, invoiceTemplate{
			docType: "INV", supplyType: "B2B", reverseCharge: true,
			buyerGSTIN: fmt.Sprintf("%sAABCD%04dD1Z5", st, i),
			buyerState: st, hsn: "9988", taxableValue: int64(100000 + i*10000),
		})
	}

	// 10 Credit Notes
	for i := 0; i < 10; i++ {
		st := supplierState
		bg := fmt.Sprintf("%sAABCB%04dB1Z5", supplierState, i)
		if i >= 5 {
			st = interStates[(i-5)%len(interStates)]
			bg = fmt.Sprintf("%sAABCC%04dC1Z5", st, i-5)
		}
		templates = append(templates, invoiceTemplate{
			docType: "CRN", supplyType: "B2B", buyerGSTIN: bg,
			buyerState: st, hsn: hsnCodes[i%len(hsnCodes)], taxableValue: int64(10000 + i*2000),
		})
	}

	// 5 Debit Notes
	for i := 0; i < 5; i++ {
		templates = append(templates, invoiceTemplate{
			docType: "DBN", supplyType: "B2B",
			buyerGSTIN: fmt.Sprintf("%sAABCB%04dB1Z5", supplierState, i+10),
			buyerState: supplierState, hsn: hsnCodes[(i+5)%len(hsnCodes)], taxableValue: int64(20000 + i*3000),
		})
	}

	// 5 NIL-rated
	for i := 0; i < 5; i++ {
		templates = append(templates, invoiceTemplate{
			docType: "INV", supplyType: "B2B",
			buyerGSTIN: fmt.Sprintf("%sAABCE%04dE1Z5", supplierState, i),
			buyerState: supplierState, hsn: "0101", taxableValue: int64(30000 + i*5000),
		})
	}

	return templates
}

var hsnCodes = []string{"9988", "8471", "9954", "9983", "9971", "8517", "7308", "3004", "6101", "9401"}

func buildInvoice(tenantID uuid.UUID, gstin, supplierState string, tmpl invoiceTemplate, idx int, baseDate time.Time) domain.Invoice {
	taxable := decimal.NewFromInt(tmpl.taxableValue)
	isIntra := tmpl.buyerState == supplierState

	var cgstRate, sgstRate, igstRate decimal.Decimal
	gstRate := decimal.NewFromInt(18)
	if tmpl.hsn == "0101" {
		gstRate = decimal.Zero
	} else if tmpl.hsn == "3004" || tmpl.hsn == "6101" {
		gstRate = decimal.NewFromInt(12)
	}

	if tmpl.supplyType == "EXP" {
		igstRate = gstRate
	} else if isIntra {
		half := gstRate.Div(decimal.NewFromInt(2))
		cgstRate = half
		sgstRate = half
	} else {
		igstRate = gstRate
	}

	cgstAmt := taxable.Mul(cgstRate).Div(decimal.NewFromInt(100))
	sgstAmt := taxable.Mul(sgstRate).Div(decimal.NewFromInt(100))
	igstAmt := taxable.Mul(igstRate).Div(decimal.NewFromInt(100))

	docDate := baseDate.Add(time.Duration(idx%28) * 24 * time.Hour)

	prefix := "INV"
	if tmpl.docType == "CRN" {
		prefix = "CRN"
	} else if tmpl.docType == "DBN" {
		prefix = "DBN"
	}

	buyerName := fmt.Sprintf("Buyer Corp %d", idx)
	if tmpl.buyerGSTIN == "URP" || tmpl.buyerGSTIN == "" {
		buyerName = fmt.Sprintf("Customer %d", idx)
	}

	stateName := stateNames[tmpl.buyerState]
	if stateName == "" {
		stateName = "Other"
	}

	return domain.Invoice{
		ID:             uuid.New(),
		TenantID:       tenantID,
		DocumentNumber: fmt.Sprintf("%s/2026/%04d", prefix, idx+1),
		DocumentDate:   docDate.Format("02/01/2006"),
		DocumentType:   tmpl.docType,
		SupplyType:     tmpl.supplyType,
		ReverseCharge:  tmpl.reverseCharge,
		Supplier: domain.Party{
			GSTIN: gstin, Name: "Complai Test Corp",
			Address: "100 MG Road, Bengaluru", StateCode: supplierState,
		},
		Buyer: domain.Party{
			GSTIN: tmpl.buyerGSTIN, Name: buyerName,
			Address: fmt.Sprintf("Street %d, %s", idx, stateName), StateCode: tmpl.buyerState,
		},
		LineItems: []domain.LineItem{{
			ItemID: fmt.Sprintf("ITEM-%04d", idx+1), Description: fmt.Sprintf("Service/Good %d", idx+1),
			HSN: tmpl.hsn, Unit: "NOS", Quantity: decimal.NewFromInt(1),
			UnitPrice: taxable, Discount: decimal.Zero, TaxableValue: taxable,
			CGSTRate: cgstRate, CGSTAmount: cgstAmt,
			SGSTRate: sgstRate, SGSTAmount: sgstAmt,
			IGSTRate: igstRate, IGSTAmount: igstAmt,
		}},
		Totals: domain.InvoiceTotals{
			TaxableValue: taxable, CGST: cgstAmt, SGST: sgstAmt, IGST: igstAmt,
			GrandTotal: taxable.Add(cgstAmt).Add(sgstAmt).Add(igstAmt),
		},
		PlaceOfSupply: tmpl.buyerState,
		SourceSystem:  "aura",
		CreatedAt:     time.Now().UTC(),
	}
}

func updateSummary(s *domain.InvoiceSummary, tmpl invoiceTemplate, supplierState string) {
	switch tmpl.docType {
	case "CRN":
		s.CreditNote++
		return
	case "DBN":
		s.DebitNote++
		return
	}

	switch tmpl.supplyType {
	case "EXP":
		s.ExportCount++
	case "B2CL":
		s.B2CLCount++
	case "B2CS":
		s.B2CSCount++
	case "B2B":
		if tmpl.reverseCharge {
			s.RCMCount++
		} else if tmpl.buyerState == supplierState {
			s.B2BIntraCount++
		} else {
			s.B2BInterCount++
		}
	}
}
