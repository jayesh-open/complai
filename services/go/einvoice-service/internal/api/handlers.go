package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/einvoice-service/internal/domain"
	"github.com/complai/complai/services/go/einvoice-service/internal/gateway"
	"github.com/complai/complai/services/go/einvoice-service/internal/store"
)

type Handlers struct {
	store     store.Repository
	irpClient *gateway.IRPClient
	clock     store.Clock
}

func NewHandlers(s store.Repository, irp *gateway.IRPClient, clock store.Clock) *Handlers {
	if clock == nil {
		clock = store.RealClock{}
	}
	return &Handlers{store: s, irpClient: irp, clock: clock}
}

func (h *Handlers) Health(w http.ResponseWriter, _ *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "einvoice-service"})
}

func (h *Handlers) GenerateIRN(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GenerateIRNRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.InvoiceNumber == "" || req.SupplierGSTIN == "" || req.InvoiceType == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invoice_number, supplier_gstin, and invoice_type are required"})
		return
	}

	inv := &domain.EInvoice{
		InvoiceNumber: req.InvoiceNumber,
		InvoiceDate:   req.InvoiceDate,
		InvoiceType:   req.InvoiceType,
		SupplierGSTIN: req.SupplierGSTIN,
		SupplierName:  req.SupplierName,
		BuyerGSTIN:    req.BuyerGSTIN,
		BuyerName:     req.BuyerName,
		SupplyType:    req.SupplyType,
		PlaceOfSupply: req.PlaceOfSupply,
		ReverseCharge: req.ReverseCharge,
		TaxableValue:  req.TaxableValue,
		CGSTAmount:    req.CGSTAmount,
		SGSTAmount:    req.SGSTAmount,
		IGSTAmount:    req.IGSTAmount,
		CessAmount:    req.CessAmount,
		TotalAmount:   req.TotalAmount,
		SourceSystem:  req.SourceSystem,
		SourceID:      req.SourceID,
	}

	if err := h.store.CreateEInvoice(r.Context(), tenantID, inv); err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID.String()).Msg("create einvoice failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create e-invoice record"})
		return
	}

	if len(req.LineItems) > 0 {
		items := make([]domain.EInvoiceLineItem, len(req.LineItems))
		for i, li := range req.LineItems {
			items[i] = domain.EInvoiceLineItem{
				InvoiceID:    inv.ID,
				Description:  li.Description,
				HSNCode:      li.HSNCode,
				Quantity:     li.Quantity,
				Unit:         li.Unit,
				UnitPrice:    li.UnitPrice,
				Discount:     li.Discount,
				TaxableValue: li.TaxableValue,
				CGSTRate:     li.CGSTRate,
				CGSTAmount:   li.CGSTAmount,
				SGSTRate:     li.SGSTRate,
				SGSTAmount:   li.SGSTAmount,
				IGSTRate:     li.IGSTRate,
				IGSTAmount:   li.IGSTAmount,
				CessRate:     li.CessRate,
				CessAmount:   li.CessAmount,
			}
		}
		if err := h.store.CreateLineItems(r.Context(), tenantID, items); err != nil {
			log.Error().Err(err).Msg("create line items failed")
		}
	}

	irpReq := &gateway.GenerateIRNGatewayRequest{
		GSTIN:   req.SupplierGSTIN,
		DocDtls: gateway.DocDetails{Typ: string(req.InvoiceType), No: req.InvoiceNumber, Dt: req.InvoiceDate},
		SupDtls: gateway.PartyDetail{GSTIN: req.SupplierGSTIN, LglNm: req.SupplierName},
		BuyDtls: gateway.PartyDetail{GSTIN: req.BuyerGSTIN, LglNm: req.BuyerName, Pos: req.PlaceOfSupply},
		ValDtls: gateway.ValDetails{
			TaxableVal: req.TaxableValue.InexactFloat64(),
			IGST:       req.IGSTAmount.InexactFloat64(),
			CGST:       req.CGSTAmount.InexactFloat64(),
			SGST:       req.SGSTAmount.InexactFloat64(),
			CesVal:     req.CessAmount.InexactFloat64(),
			TotInvVal:  req.TotalAmount.InexactFloat64(),
		},
	}

	irpResp, err := h.irpClient.GenerateIRN(r.Context(), tenantID, irpReq)
	if err != nil {
		log.Error().Err(err).Msg("irp gateway generate IRN failed")
		_ = h.store.UpdateIRNFailed(r.Context(), tenantID, inv.ID, err.Error())
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": "IRN generation failed: " + err.Error()})
		return
	}

	if err := h.store.UpdateIRNGenerated(r.Context(), tenantID, inv.ID, irpResp.IRN, irpResp.AckNo, irpResp.SignedInvoice, irpResp.SignedQRCode); err != nil {
		log.Error().Err(err).Msg("update IRN generated failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to persist IRN"})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GenerateIRNResponse{
		ID:            inv.ID,
		IRN:           irpResp.IRN,
		AckNo:         irpResp.AckNo,
		Status:        domain.IRNStatusGenerated,
		SignedInvoice: irpResp.SignedInvoice,
		SignedQRCode:  irpResp.SignedQRCode,
		GeneratedAt:   time.Now(),
	})
}

func (h *Handlers) CancelIRN(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid invoice id"})
		return
	}

	var req domain.CancelIRNRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.Reason == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "reason is required"})
		return
	}

	inv, err := h.store.GetEInvoice(r.Context(), tenantID, id)
	if err != nil {
		log.Error().Err(err).Msg("get einvoice for cancel failed")
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "e-invoice not found"})
		return
	}

	if inv.Status != domain.IRNStatusGenerated {
		httputil.JSON(w, http.StatusUnprocessableEntity, map[string]string{"error": fmt.Sprintf("cannot cancel e-invoice in status %s", inv.Status)})
		return
	}

	if !store.CancellationWindowOpen(inv.IRNGeneratedAt, h.clock) {
		httputil.JSON(w, http.StatusUnprocessableEntity, map[string]string{
			"error": "cancellation window expired: IRN can only be cancelled within 24 hours of generation",
		})
		return
	}

	cancelReq := &gateway.CancelIRNGatewayRequest{
		IRN:    inv.IRN,
		CnlRsn: req.Reason,
		CnlRem: req.Remark,
	}
	_, err = h.irpClient.CancelIRN(r.Context(), tenantID, cancelReq)
	if err != nil {
		log.Error().Err(err).Msg("irp gateway cancel IRN failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": "IRN cancellation failed: " + err.Error()})
		return
	}

	if err := h.store.UpdateIRNCancelled(r.Context(), tenantID, id, req.Reason); err != nil {
		log.Error().Err(err).Msg("update IRN cancelled failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to persist cancellation"})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.CancelIRNResponse{
		ID:          id,
		IRN:         inv.IRN,
		Status:      domain.IRNStatusCancelled,
		CancelledAt: time.Now(),
	})
}

func (h *Handlers) GetEInvoice(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid invoice id"})
		return
	}

	inv, err := h.store.GetEInvoice(r.Context(), tenantID, id)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "e-invoice not found"})
		return
	}

	httputil.JSON(w, http.StatusOK, inv)
}

func (h *Handlers) GetEInvoiceByIRN(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	irn := chi.URLParam(r, "irn")
	if irn == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "irn is required"})
		return
	}

	inv, err := h.store.GetEInvoiceByIRN(r.Context(), tenantID, irn)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "e-invoice not found"})
		return
	}

	httputil.JSON(w, http.StatusOK, inv)
}

func (h *Handlers) ListEInvoices(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	gstin := r.URL.Query().Get("gstin")
	if gstin == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "gstin query parameter is required"})
		return
	}

	pageSize := 50
	pageOffset := 0

	req := &domain.ListEInvoicesRequest{
		GSTIN:      gstin,
		Status:     r.URL.Query().Get("status"),
		PageSize:   pageSize,
		PageOffset: pageOffset,
	}

	invoices, total, err := h.store.ListEInvoices(r.Context(), tenantID, req)
	if err != nil {
		log.Error().Err(err).Msg("list einvoices failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list e-invoices"})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.ListEInvoicesResponse{
		Invoices:   invoices,
		TotalCount: total,
		PageSize:   pageSize,
		PageOffset: pageOffset,
	})
}

func (h *Handlers) GetSummary(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	gstin := r.URL.Query().Get("gstin")
	if gstin == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "gstin query parameter is required"})
		return
	}

	summary, err := h.store.GetSummary(r.Context(), tenantID, gstin)
	if err != nil {
		log.Error().Err(err).Msg("get summary failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get summary"})
		return
	}

	httputil.JSON(w, http.StatusOK, summary)
}

func tenantIDFromRequest(r *http.Request) (uuid.UUID, error) {
	h := r.Header.Get("X-Tenant-Id")
	if h == "" {
		return uuid.Nil, fmt.Errorf("missing X-Tenant-Id header")
	}
	return uuid.Parse(h)
}
