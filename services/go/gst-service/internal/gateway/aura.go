package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/gst-service/internal/domain"
)

type AuraClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAuraClient(baseURL string) *AuraClient {
	return &AuraClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type auraInvoice struct {
	ID             uuid.UUID       `json:"id"`
	DocumentNumber string          `json:"document_number"`
	DocumentDate   string          `json:"document_date"`
	DocumentType   string          `json:"document_type"`
	SupplyType     string          `json:"supply_type"`
	ReverseCharge  bool            `json:"reverse_charge"`
	Supplier       auraParty       `json:"supplier"`
	Buyer          auraParty       `json:"buyer"`
	LineItems      []auraLineItem  `json:"line_items"`
	Totals         auraTotals      `json:"totals"`
	PlaceOfSupply  string          `json:"place_of_supply"`
	SourceSystem   string          `json:"source_system"`
}

type auraParty struct {
	GSTIN     string `json:"gstin"`
	Name      string `json:"name"`
	StateCode string `json:"state_code"`
}

type auraLineItem struct {
	HSN          string          `json:"hsn"`
	TaxableValue json.Number     `json:"taxable_value"`
	CGSTRate     json.Number     `json:"cgst_rate"`
	CGSTAmount   json.Number     `json:"cgst_amount"`
	SGSTRate     json.Number     `json:"sgst_rate"`
	SGSTAmount   json.Number     `json:"sgst_amount"`
	IGSTRate     json.Number     `json:"igst_rate"`
	IGSTAmount   json.Number     `json:"igst_amount"`
}

type auraTotals struct {
	TaxableValue json.Number `json:"taxable_value"`
	CGST         json.Number `json:"cgst"`
	SGST         json.Number `json:"sgst"`
	IGST         json.Number `json:"igst"`
	GrandTotal   json.Number `json:"grand_total"`
}

type auraListResponse struct {
	Invoices   []auraInvoice `json:"invoices"`
	TotalCount int           `json:"total_count"`
}

func (c *AuraClient) FetchARInvoices(ctx context.Context, tenantID uuid.UUID, gstin, period string) ([]domain.SalesRegisterEntry, error) {
	url := fmt.Sprintf("%s/v1/gateway/aura/invoices?gstin=%s&period=%s", c.baseURL, gstin, period)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Tenant-Id", tenantID.String())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch invoices: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("aura gateway returned %d", resp.StatusCode)
	}

	var wrapper httputil.SuccessResponse
	var listResp auraListResponse
	wrapper.Data = &listResp
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	entries := make([]domain.SalesRegisterEntry, 0, len(listResp.Invoices))
	for _, inv := range listResp.Invoices {
		e := mapAuraInvoiceToEntry(tenantID, gstin, period, inv)
		entries = append(entries, e)
	}

	return entries, nil
}

func mapAuraInvoiceToEntry(tenantID uuid.UUID, gstin, period string, inv auraInvoice) domain.SalesRegisterEntry {
	hsn := ""
	if len(inv.LineItems) > 0 {
		hsn = inv.LineItems[0].HSN
	}

	return domain.SalesRegisterEntry{
		ID:             uuid.New(),
		TenantID:       tenantID,
		GSTIN:          gstin,
		ReturnPeriod:   period,
		DocumentNumber: inv.DocumentNumber,
		DocumentDate:   inv.DocumentDate,
		DocumentType:   inv.DocumentType,
		SupplyType:     inv.SupplyType,
		ReverseCharge:  inv.ReverseCharge,
		SupplierGSTIN:  inv.Supplier.GSTIN,
		BuyerGSTIN:     inv.Buyer.GSTIN,
		BuyerName:      inv.Buyer.Name,
		BuyerState:     inv.Buyer.StateCode,
		PlaceOfSupply:  inv.PlaceOfSupply,
		HSN:            hsn,
		TaxableValue:   decimalFromJSON(inv.Totals.TaxableValue),
		CGSTRate:       firstLineItemDecimal(inv.LineItems, func(li auraLineItem) json.Number { return li.CGSTRate }),
		CGSTAmount:     decimalFromJSON(inv.Totals.CGST),
		SGSTRate:       firstLineItemDecimal(inv.LineItems, func(li auraLineItem) json.Number { return li.SGSTRate }),
		SGSTAmount:     decimalFromJSON(inv.Totals.SGST),
		IGSTRate:       firstLineItemDecimal(inv.LineItems, func(li auraLineItem) json.Number { return li.IGSTRate }),
		IGSTAmount:     decimalFromJSON(inv.Totals.IGST),
		GrandTotal:     decimalFromJSON(inv.Totals.GrandTotal),
		SourceSystem:   inv.SourceSystem,
		SourceID:       inv.ID.String(),
	}
}

func decimalFromJSON(n json.Number) decimal.Decimal {
	d, _ := decimal.NewFromString(n.String())
	return d
}

func firstLineItemDecimal(items []auraLineItem, fn func(auraLineItem) json.Number) decimal.Decimal {
	if len(items) == 0 {
		return decimal.Zero
	}
	return decimalFromJSON(fn(items[0]))
}
