package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/recon-service/internal/domain"
	"github.com/complai/complai/services/go/recon-service/internal/gateway"
	"github.com/complai/complai/services/go/recon-service/internal/matcher"
	"github.com/complai/complai/services/go/recon-service/internal/store"
)

type Handlers struct {
	store      store.Repository
	apexClient *gateway.ApexClient
	gstnClient *gateway.GSTNClient
}

func NewHandlers(s store.Repository, apex *gateway.ApexClient, gstn *gateway.GSTNClient) *Handlers {
	return &Handlers{store: s, apexClient: apex, gstnClient: gstn}
}

func (h *Handlers) Health(w http.ResponseWriter, _ *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "recon-service"})
}

func (h *Handlers) RunRecon(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.RunReconRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.GSTIN == "" || req.ReturnPeriod == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "gstin and return_period are required"})
		return
	}

	// Create the run record
	run := &domain.ReconRun{
		GSTIN:        req.GSTIN,
		ReturnPeriod: req.ReturnPeriod,
		Status:       "RUNNING",
		StartedAt:    time.Now().UTC(),
	}
	if err := h.store.CreateRun(r.Context(), tenantID, run); err != nil {
		log.Error().Err(err).Msg("create recon run failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create recon run"})
		return
	}

	// Fetch PR from Apex gateway
	prEntries, err := h.apexClient.FetchAPInvoices(r.Context(), tenantID, req.GSTIN, req.ReturnPeriod)
	if err != nil {
		log.Error().Err(err).Msg("fetch AP invoices from Apex failed")
		run.Status = "FAILED"
		now := time.Now().UTC()
		run.CompletedAt = &now
		_ = h.store.UpdateRun(r.Context(), tenantID, run)
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": "failed to fetch purchase register from Apex"})
		return
	}

	// Fetch GSTR-2B from GSTN gateway
	gstr2bEntries, err := h.gstnClient.FetchGSTR2B(r.Context(), tenantID, req.GSTIN, req.ReturnPeriod)
	if err != nil {
		log.Error().Err(err).Msg("fetch GSTR-2B from GSTN failed")
		run.Status = "FAILED"
		now := time.Now().UTC()
		run.CompletedAt = &now
		_ = h.store.UpdateRun(r.Context(), tenantID, run)
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": "failed to fetch GSTR-2B from GSTN"})
		return
	}

	// Run the 5-stage match pipeline
	matches := matcher.Run(prEntries, gstr2bEntries, req.GSTIN, req.ReturnPeriod, run.ID)

	// Set tenant_id on all matches
	for i := range matches {
		matches[i].TenantID = tenantID
	}

	// Store matches
	if len(matches) > 0 {
		if err := h.store.BulkInsertMatches(r.Context(), tenantID, matches); err != nil {
			log.Error().Err(err).Msg("bulk insert matches failed")
			run.Status = "FAILED"
			now := time.Now().UTC()
			run.CompletedAt = &now
			_ = h.store.UpdateRun(r.Context(), tenantID, run)
			httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to store matches"})
			return
		}
	}

	// Compute summary counts
	var matched, mismatch, partial, missing2b, missingPR, duplicate int
	for _, m := range matches {
		switch m.MatchType {
		case domain.MatchTypeDirect:
			matched++
		case domain.MatchTypeProbable:
			mismatch++
		case domain.MatchTypePartial:
			partial++
		case domain.MatchTypeMissing2B:
			missing2b++
		case domain.MatchTypeMissingPR:
			missingPR++
		case domain.MatchTypeDuplicate:
			duplicate++
		}
	}

	// Update the run
	now := time.Now().UTC()
	run.Status = "COMPLETED"
	run.PRCount = len(prEntries)
	run.GSTR2BCount = len(gstr2bEntries)
	run.Matched = matched
	run.Mismatch = mismatch
	run.Partial = partial
	run.Missing2B = missing2b
	run.MissingPR = missingPR
	run.Duplicate = duplicate
	run.CompletedAt = &now
	if err := h.store.UpdateRun(r.Context(), tenantID, run); err != nil {
		log.Error().Err(err).Msg("update recon run failed")
	}

	httputil.JSON(w, http.StatusOK, domain.RunReconResponse{
		RunID:  run.ID,
		Status: run.Status,
	})
}

func (h *Handlers) GetRun(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	runID, err := uuid.Parse(chi.URLParam(r, "run_id"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid run_id"})
		return
	}

	run, err := h.store.GetRun(r.Context(), tenantID, runID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "run not found"})
		return
	}

	summary, _ := h.store.GetBucketSummary(r.Context(), tenantID, runID)

	httputil.JSON(w, http.StatusOK, map[string]interface{}{
		"run":     run,
		"summary": summary,
	})
}

func (h *Handlers) ListMatches(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	runIDStr := r.URL.Query().Get("run_id")
	if runIDStr == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "run_id is required"})
		return
	}
	runID, err := uuid.Parse(runIDStr)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid run_id"})
		return
	}

	matchType := r.URL.Query().Get("match_type")
	status := r.URL.Query().Get("status")

	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	matches, err := h.store.ListMatches(r.Context(), tenantID, runID, matchType, status, limit, offset)
	if err != nil {
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list matches"})
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]interface{}{
		"matches":     matches,
		"total_count": len(matches),
	})
}

func (h *Handlers) AcceptMatch(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	matchID, err := uuid.Parse(chi.URLParam(r, "match_id"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid match_id"})
		return
	}

	userID := userIDFromRequest(r)
	if userID == uuid.Nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "X-User-Id header is required"})
		return
	}

	if err := h.store.UpdateMatchStatus(r.Context(), tenantID, matchID, domain.MatchStatusAccepted, &userID); err != nil {
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to accept match"})
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]string{"status": "accepted", "match_id": matchID.String()})
}

func (h *Handlers) BulkAccept(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.BulkAcceptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if len(req.MatchIDs) == 0 {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "match_ids is required"})
		return
	}

	userID := userIDFromRequest(r)
	if userID == uuid.Nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "X-User-Id header is required"})
		return
	}

	accepted := 0
	for _, id := range req.MatchIDs {
		if err := h.store.UpdateMatchStatus(r.Context(), tenantID, id, domain.MatchStatusAccepted, &userID); err != nil {
			log.Error().Err(err).Str("match_id", id.String()).Msg("accept match failed")
			continue
		}
		accepted++
	}

	httputil.JSON(w, http.StatusOK, map[string]interface{}{
		"accepted": accepted,
		"total":    len(req.MatchIDs),
	})
}

func (h *Handlers) IMSActionHandler(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.IMSActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.InvoiceID == "" || req.Action == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invoice_id and action are required"})
		return
	}

	if req.Action != "ACCEPT" && req.Action != "REJECT" && req.Action != "PENDING" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "action must be ACCEPT, REJECT, or PENDING"})
		return
	}

	gstin := r.URL.Query().Get("gstin")
	returnPeriod := r.URL.Query().Get("return_period")
	if gstin == "" || returnPeriod == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "gstin and return_period query params are required"})
		return
	}

	// Proxy to GSTN gateway
	actionResp, err := h.gstnClient.SendIMSAction(r.Context(), tenantID, gstin, returnPeriod, req.InvoiceID, req.Action, req.Reason)
	if err != nil {
		log.Error().Err(err).Msg("IMS action via GSTN gateway failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": "failed to send IMS action to GSTN"})
		return
	}

	// Record the action locally
	userID := userIDFromRequest(r)
	action := &domain.IMSAction{
		GSTIN:        gstin,
		ReturnPeriod: returnPeriod,
		InvoiceID:    req.InvoiceID,
		Action:       req.Action,
		Reason:       req.Reason,
		CreatedBy:    userID,
	}
	if err := h.store.CreateIMSAction(r.Context(), tenantID, action); err != nil {
		log.Error().Err(err).Msg("store ims action failed")
	}

	httputil.JSON(w, http.StatusOK, actionResp)
}

func (h *Handlers) GetIMSState(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	gstin := r.URL.Query().Get("gstin")
	returnPeriod := r.URL.Query().Get("return_period")
	if gstin == "" || returnPeriod == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "gstin and return_period query params are required"})
		return
	}

	imsState, err := h.gstnClient.FetchIMSState(r.Context(), tenantID, gstin, returnPeriod)
	if err != nil {
		log.Error().Err(err).Msg("fetch IMS state failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": "failed to fetch IMS state from GSTN"})
		return
	}

	httputil.JSON(w, http.StatusOK, imsState)
}

func tenantIDFromRequest(r *http.Request) (uuid.UUID, error) {
	h := r.Header.Get("X-Tenant-Id")
	if h == "" {
		return uuid.Nil, fmt.Errorf("missing X-Tenant-Id header")
	}
	return uuid.Parse(h)
}

func userIDFromRequest(r *http.Request) uuid.UUID {
	id, err := uuid.Parse(r.Header.Get("X-User-Id"))
	if err != nil {
		return uuid.Nil
	}
	return id
}
