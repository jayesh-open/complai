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
	"github.com/complai/complai/services/go/ewb-service/internal/domain"
	"github.com/complai/complai/services/go/ewb-service/internal/gateway"
	"github.com/complai/complai/services/go/ewb-service/internal/store"
)

type Handlers struct {
	store     store.Repository
	ewbClient *gateway.EWBClient
	clock     store.Clock
}

func NewHandlers(s store.Repository, ewb *gateway.EWBClient, clock store.Clock) *Handlers {
	if clock == nil {
		clock = store.RealClock{}
	}
	return &Handlers{store: s, ewbClient: ewb, clock: clock}
}

func (h *Handlers) Health(w http.ResponseWriter, _ *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "ewb-service"})
}

func (h *Handlers) GenerateEWB(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.GenerateEWBRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.DocNumber == "" || req.SupplierGSTIN == "" || req.DocType == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "doc_number, supplier_gstin, and doc_type are required"})
		return
	}

	ewb := &domain.EWayBill{
		DocType:       req.DocType,
		DocNumber:     req.DocNumber,
		DocDate:       req.DocDate,
		SupplierGSTIN: req.SupplierGSTIN,
		SupplierName:  req.SupplierName,
		BuyerGSTIN:    req.BuyerGSTIN,
		BuyerName:     req.BuyerName,
		SupplyType:    req.SupplyType,
		SubSupplyType: req.SubSupplyType,
		TransportMode: req.TransportMode,
		VehicleNumber: req.VehicleNumber,
		VehicleType:   req.VehicleType,
		TransporterID: req.TransporterID,
		FromPlace:     req.FromPlace,
		FromState:     req.FromState,
		FromPincode:   req.FromPincode,
		ToPlace:       req.ToPlace,
		ToState:       req.ToState,
		ToPincode:     req.ToPincode,
		DistanceKM:    req.DistanceKM,
		TaxableValue:  req.TaxableValue,
		CGSTAmount:    req.CGSTAmount,
		SGSTAmount:    req.SGSTAmount,
		IGSTAmount:    req.IGSTAmount,
		CessAmount:    req.CessAmount,
		TotalValue:    req.TotalValue,
		SourceSystem:  req.SourceSystem,
		SourceID:      req.SourceID,
	}

	if err := h.store.CreateEWB(r.Context(), tenantID, ewb); err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID.String()).Msg("create ewb failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create e-way bill record"})
		return
	}

	if len(req.Items) > 0 {
		items := make([]domain.EWBItem, len(req.Items))
		for i, li := range req.Items {
			items[i] = domain.EWBItem{
				EWBID:        ewb.ID,
				ProductName:  li.ProductName,
				HSNCode:      li.HSNCode,
				Quantity:     li.Quantity,
				Unit:         li.Unit,
				TaxableValue: li.TaxableValue,
				CGSTRate:     li.CGSTRate,
				SGSTRate:     li.SGSTRate,
				IGSTRate:     li.IGSTRate,
				CessRate:     li.CessRate,
			}
		}
		if err := h.store.CreateItems(r.Context(), tenantID, items); err != nil {
			log.Error().Err(err).Msg("create ewb items failed")
		}
	}

	isODC := req.VehicleType == "O"
	gwReq := &gateway.GenerateEWBGatewayRequest{
		GSTIN:         req.SupplierGSTIN,
		SupplyType:    req.SupplyType,
		SubSupplyType: req.SubSupplyType,
		DocType:       req.DocType,
		DocNo:         req.DocNumber,
		DocDate:       req.DocDate,
		FromGSTIN:     req.SupplierGSTIN,
		FromName:      req.SupplierName,
		FromPlace:     req.FromPlace,
		FromState:     req.FromState,
		FromPincode:   req.FromPincode,
		ToGSTIN:       req.BuyerGSTIN,
		ToName:        req.BuyerName,
		ToPlace:       req.ToPlace,
		ToState:       req.ToState,
		ToPincode:     req.ToPincode,
		TransportMode: req.TransportMode,
		VehicleNo:     req.VehicleNumber,
		VehicleType:   req.VehicleType,
		TransporterID: req.TransporterID,
		DistanceKM:    req.DistanceKM,
		TotalValue:    req.TotalValue.InexactFloat64(),
		TaxableValue:  req.TaxableValue.InexactFloat64(),
		CGSTAmount:    req.CGSTAmount.InexactFloat64(),
		SGSTAmount:    req.SGSTAmount.InexactFloat64(),
		IGSTAmount:    req.IGSTAmount.InexactFloat64(),
		CessAmount:    req.CessAmount.InexactFloat64(),
	}

	gwResp, err := h.ewbClient.GenerateEWB(r.Context(), tenantID, gwReq)
	if err != nil {
		log.Error().Err(err).Msg("ewb gateway generate failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": "EWB generation failed: " + err.Error()})
		return
	}

	now := h.clock.Now()
	days := store.ValidityDays(req.DistanceKM, isODC)
	validUntil := now.Add(time.Duration(days) * 24 * time.Hour)

	if err := h.store.UpdateEWBGenerated(r.Context(), tenantID, ewb.ID, gwResp.EWBNumber, now, validUntil); err != nil {
		log.Error().Err(err).Msg("update ewb generated failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to persist EWB"})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.GenerateEWBResponse{
		ID:         ewb.ID,
		EWBNumber:  gwResp.EWBNumber,
		Status:     domain.EWBStatusActive,
		ValidFrom:  now,
		ValidUntil: validUntil,
	})
}

func (h *Handlers) CancelEWB(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid ewb id"})
		return
	}

	var req domain.CancelEWBRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Reason == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "reason is required"})
		return
	}

	ewb, err := h.store.GetEWB(r.Context(), tenantID, id)
	if err != nil {
		log.Error().Err(err).Msg("get ewb for cancel failed")
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "e-way bill not found"})
		return
	}

	if !domain.CanTransitionTo(ewb.Status, domain.EWBStatusCancelled) {
		httputil.JSON(w, http.StatusUnprocessableEntity, map[string]string{
			"error": fmt.Sprintf("cannot cancel e-way bill in status %s", ewb.Status),
		})
		return
	}

	if !store.CancellationWindowOpen(ewb.GeneratedAt, h.clock) {
		httputil.JSON(w, http.StatusUnprocessableEntity, map[string]string{
			"error": "cancellation window expired: EWB can only be cancelled within 24 hours of generation",
		})
		return
	}

	gwReq := &gateway.CancelEWBGatewayRequest{
		EWBNo:  ewb.EWBNumber,
		Reason: req.Reason,
		Remark: req.Remark,
	}
	_, err = h.ewbClient.CancelEWB(r.Context(), tenantID, gwReq)
	if err != nil {
		log.Error().Err(err).Msg("ewb gateway cancel failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": "EWB cancellation failed: " + err.Error()})
		return
	}

	if err := h.store.UpdateEWBCancelled(r.Context(), tenantID, id, req.Reason); err != nil {
		log.Error().Err(err).Msg("update ewb cancelled failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to persist cancellation"})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.CancelEWBResponse{
		ID:          id,
		EWBNumber:   ewb.EWBNumber,
		Status:      domain.EWBStatusCancelled,
		CancelledAt: h.clock.Now(),
	})
}

func (h *Handlers) UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid ewb id"})
		return
	}

	var req domain.UpdateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.VehicleNumber == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "vehicle_number is required"})
		return
	}

	ewb, err := h.store.GetEWB(r.Context(), tenantID, id)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "e-way bill not found"})
		return
	}

	if !domain.CanTransitionTo(ewb.Status, domain.EWBStatusVehicleUpdated) {
		httputil.JSON(w, http.StatusUnprocessableEntity, map[string]string{
			"error": fmt.Sprintf("cannot update vehicle for e-way bill in status %s", ewb.Status),
		})
		return
	}

	gwReq := &gateway.UpdateVehicleGatewayRequest{
		EWBNo:         ewb.EWBNumber,
		VehicleNo:     req.VehicleNumber,
		FromPlace:     req.FromPlace,
		FromState:     req.FromState,
		Reason:        req.Reason,
		TransportMode: req.TransportMode,
		Remark:        req.Remark,
	}
	_, err = h.ewbClient.UpdateVehicle(r.Context(), tenantID, gwReq)
	if err != nil {
		log.Error().Err(err).Msg("ewb gateway update vehicle failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": "vehicle update failed: " + err.Error()})
		return
	}

	vu := &domain.VehicleUpdate{
		EWBID:         id,
		VehicleNumber: req.VehicleNumber,
		FromPlace:     req.FromPlace,
		FromState:     req.FromState,
		TransportMode: req.TransportMode,
		Reason:        req.Reason,
		Remark:        req.Remark,
	}
	if err := h.store.CreateVehicleUpdate(r.Context(), tenantID, vu); err != nil {
		log.Error().Err(err).Msg("create vehicle update record failed")
	}

	if err := h.store.UpdateEWBVehicle(r.Context(), tenantID, id, req.VehicleNumber); err != nil {
		log.Error().Err(err).Msg("update ewb vehicle failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to persist vehicle update"})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.UpdateVehicleResponse{
		ID:            id,
		EWBNumber:     ewb.EWBNumber,
		VehicleNumber: req.VehicleNumber,
		Status:        domain.EWBStatusVehicleUpdated,
	})
}

func (h *Handlers) ExtendValidity(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid ewb id"})
		return
	}

	var req domain.ExtendValidityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.RemainingDistance <= 0 {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "remaining_distance must be positive"})
		return
	}

	ewb, err := h.store.GetEWB(r.Context(), tenantID, id)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "e-way bill not found"})
		return
	}

	if !domain.CanTransitionTo(ewb.Status, domain.EWBStatusExtended) {
		httputil.JSON(w, http.StatusUnprocessableEntity, map[string]string{
			"error": fmt.Sprintf("cannot extend e-way bill in status %s", ewb.Status),
		})
		return
	}

	gwReq := &gateway.ExtendValidityGatewayRequest{
		EWBNo:             ewb.EWBNumber,
		FromPlace:         req.FromPlace,
		FromState:         req.FromState,
		RemainingDistance:  req.RemainingDistance,
		ExtendReason:      req.ExtendReason,
		TransitType:       req.TransitType,
		ConsignmentStatus: req.ConsignmentStatus,
		Remark:            req.Remark,
	}
	_, err = h.ewbClient.ExtendValidity(r.Context(), tenantID, gwReq)
	if err != nil {
		log.Error().Err(err).Msg("ewb gateway extend validity failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": "validity extension failed: " + err.Error()})
		return
	}

	isODC := ewb.VehicleType == "O"
	days := store.ValidityDays(req.RemainingDistance, isODC)
	base := h.clock.Now()
	if ewb.ValidUntil != nil && ewb.ValidUntil.After(base) {
		base = *ewb.ValidUntil
	}
	newValidUntil := base.Add(time.Duration(days) * 24 * time.Hour)

	if err := h.store.UpdateEWBValidity(r.Context(), tenantID, id, newValidUntil); err != nil {
		log.Error().Err(err).Msg("update ewb validity failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to persist validity extension"})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.ExtendValidityResponse{
		ID:         id,
		EWBNumber:  ewb.EWBNumber,
		Status:     domain.EWBStatusExtended,
		ValidUntil: newValidUntil,
	})
}

func (h *Handlers) ConsolidateEWB(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.ConsolidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if len(req.EWBIDS) < 2 {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "at least 2 EWBs required for consolidation"})
		return
	}
	if req.VehicleNumber == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "vehicle_number is required"})
		return
	}

	var ewbNumbers []string
	for _, ewbID := range req.EWBIDS {
		ewb, err := h.store.GetEWB(r.Context(), tenantID, ewbID)
		if err != nil {
			httputil.JSON(w, http.StatusNotFound, map[string]string{"error": fmt.Sprintf("e-way bill not found: %s", ewbID)})
			return
		}
		if !domain.CanTransitionTo(ewb.Status, domain.EWBStatusConsolidated) {
			httputil.JSON(w, http.StatusUnprocessableEntity, map[string]string{
				"error": fmt.Sprintf("cannot consolidate e-way bill %s in status %s", ewbID, ewb.Status),
			})
			return
		}
		ewbNumbers = append(ewbNumbers, ewb.EWBNumber)
	}

	gwReq := &gateway.ConsolidateGatewayRequest{
		FromGSTIN:     req.FromPlace,
		FromPlace:     req.FromPlace,
		FromState:     req.FromState,
		ToPlace:       req.ToPlace,
		ToState:       req.ToState,
		VehicleNo:     req.VehicleNumber,
		TransportMode: req.TransportMode,
		EWBNumbers:    ewbNumbers,
	}
	gwResp, err := h.ewbClient.Consolidate(r.Context(), tenantID, gwReq)
	if err != nil {
		log.Error().Err(err).Msg("ewb gateway consolidate failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": "consolidation failed: " + err.Error()})
		return
	}

	consolidation := &domain.Consolidation{
		ConsolidatedEWBNumber: gwResp.ConsolidatedEWBNo,
		VehicleNumber:         req.VehicleNumber,
		FromPlace:             req.FromPlace,
		FromState:             req.FromState,
		ToPlace:               req.ToPlace,
		ToState:               req.ToState,
		TransportMode:         req.TransportMode,
		Status:                "ACTIVE",
	}
	if err := h.store.CreateConsolidation(r.Context(), tenantID, consolidation); err != nil {
		log.Error().Err(err).Msg("create consolidation record failed")
	}

	for _, ewbID := range req.EWBIDS {
		if err := h.store.SetConsolidatedEWBID(r.Context(), tenantID, ewbID, consolidation.ID); err != nil {
			log.Error().Err(err).Str("ewb_id", ewbID.String()).Msg("set consolidated ewb id failed")
		}
	}

	httputil.JSON(w, http.StatusOK, domain.ConsolidateResponse{
		ConsolidationID:       consolidation.ID,
		ConsolidatedEWBNumber: gwResp.ConsolidatedEWBNo,
		EWBCount:              len(req.EWBIDS),
		Status:                "ACTIVE",
	})
}

func (h *Handlers) GetEWB(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid ewb id"})
		return
	}

	ewb, err := h.store.GetEWB(r.Context(), tenantID, id)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "e-way bill not found"})
		return
	}

	httputil.JSON(w, http.StatusOK, ewb)
}

func (h *Handlers) GetEWBByNumber(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	ewbNo := chi.URLParam(r, "ewbNo")
	if ewbNo == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "ewb number is required"})
		return
	}

	ewb, err := h.store.GetEWBByNumber(r.Context(), tenantID, ewbNo)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "e-way bill not found"})
		return
	}

	httputil.JSON(w, http.StatusOK, ewb)
}

func (h *Handlers) GetVehicleHistory(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid ewb id"})
		return
	}

	updates, err := h.store.GetVehicleUpdates(r.Context(), tenantID, id)
	if err != nil {
		log.Error().Err(err).Msg("get vehicle history failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get vehicle history"})
		return
	}

	httputil.JSON(w, http.StatusOK, updates)
}

func (h *Handlers) ListEWBs(w http.ResponseWriter, r *http.Request) {
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

	req := &domain.ListEWBRequest{
		GSTIN:      gstin,
		Status:     r.URL.Query().Get("status"),
		PageSize:   50,
		PageOffset: 0,
	}

	ewbs, total, err := h.store.ListEWBs(r.Context(), tenantID, req)
	if err != nil {
		log.Error().Err(err).Msg("list ewbs failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list e-way bills"})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.ListEWBResponse{
		EWayBills:  ewbs,
		TotalCount: total,
		PageSize:   req.PageSize,
		PageOffset: req.PageOffset,
	})
}

func tenantIDFromRequest(r *http.Request) (uuid.UUID, error) {
	h := r.Header.Get("X-Tenant-Id")
	if h == "" {
		return uuid.Nil, fmt.Errorf("missing X-Tenant-Id header")
	}
	return uuid.Parse(h)
}
