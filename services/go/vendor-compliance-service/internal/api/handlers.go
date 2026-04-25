package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/vendor-compliance-service/internal/domain"
	"github.com/complai/complai/services/go/vendor-compliance-service/internal/gateway"
	"github.com/complai/complai/services/go/vendor-compliance-service/internal/scorer"
	"github.com/complai/complai/services/go/vendor-compliance-service/internal/store"
)

type Handlers struct {
	store      store.Repository
	apexClient *gateway.ApexClient
}

func NewHandlers(s store.Repository, apex *gateway.ApexClient) *Handlers {
	return &Handlers{store: s, apexClient: apex}
}

func (h *Handlers) Health(w http.ResponseWriter, _ *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "vendor-compliance-service"})
}

func (h *Handlers) TriggerSync(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Create sync status record
	syncStatus := &domain.SyncStatus{
		SyncType:  "full",
		Status:    "running",
		StartedAt: time.Now().UTC(),
	}
	if err := h.store.CreateSyncStatus(r.Context(), tenantID, syncStatus); err != nil {
		log.Error().Err(err).Msg("create sync status failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create sync status"})
		return
	}

	// Fetch vendors from Apex gateway
	vendors, err := h.apexClient.FetchVendors(r.Context(), tenantID)
	if err != nil {
		log.Error().Err(err).Msg("fetch vendors from apex failed")
		_ = h.store.UpdateSyncStatus(r.Context(), tenantID, syncStatus.ID, "failed", 0, 0, err.Error())
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": "failed to fetch vendors from Apex"})
		return
	}

	// Upsert vendor snapshots
	for i := range vendors {
		if err := h.store.UpsertVendorSnapshot(r.Context(), tenantID, &vendors[i]); err != nil {
			log.Error().Err(err).Str("vendor_id", vendors[i].VendorID).Msg("upsert vendor failed")
		}
	}

	// Score each vendor
	scoredCount := 0
	for _, vendor := range vendors {
		invoices, err := h.apexClient.FetchAPInvoices(r.Context(), tenantID, vendor.VendorID)
		if err != nil {
			log.Error().Err(err).Str("vendor_id", vendor.VendorID).Msg("fetch invoices failed")
			continue
		}

		score := scorer.Score(vendor, invoices)
		if err := h.store.CreateComplianceScore(r.Context(), tenantID, &score); err != nil {
			log.Error().Err(err).Str("vendor_id", vendor.VendorID).Msg("create score failed")
			continue
		}
		scoredCount++
	}

	// Update sync status
	_ = h.store.UpdateSyncStatus(r.Context(), tenantID, syncStatus.ID, "completed", len(vendors), scoredCount, "")

	httputil.JSON(w, http.StatusOK, domain.SyncResponse{
		SyncID:      syncStatus.ID,
		VendorCount: len(vendors),
		ScoredCount: scoredCount,
		Status:      "completed",
	})
}

func (h *Handlers) ListVendors(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	vendors, err := h.store.ListVendorSnapshots(r.Context(), tenantID)
	if err != nil {
		log.Error().Err(err).Msg("list vendors failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list vendors"})
		return
	}

	scores, err := h.store.ListLatestScores(r.Context(), tenantID)
	if err != nil {
		log.Error().Err(err).Msg("list scores failed")
	}

	scoreMap := make(map[string]domain.ComplianceScore)
	for _, s := range scores {
		scoreMap[s.VendorID] = s
	}

	var vendorScores []domain.VendorScoreResponse
	for _, v := range vendors {
		vs := domain.VendorScoreResponse{Vendor: v}
		if s, ok := scoreMap[v.VendorID]; ok {
			vs.Score = s
		}
		vendorScores = append(vendorScores, vs)
	}

	summary, _ := h.store.GetScoreSummary(r.Context(), tenantID)
	if summary == nil {
		summary = &domain.ScoreSummary{}
	}

	httputil.JSON(w, http.StatusOK, domain.VendorListResponse{
		Vendors: vendorScores,
		Total:   len(vendorScores),
		Summary: *summary,
	})
}

func (h *Handlers) GetVendorScore(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	vendorID := chi.URLParam(r, "vendorId")
	if vendorID == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "vendorId is required"})
		return
	}

	vendor, err := h.store.GetVendorSnapshot(r.Context(), tenantID, vendorID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "vendor not found"})
		return
	}

	score, err := h.store.GetLatestScore(r.Context(), tenantID, vendorID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "score not found for vendor"})
		return
	}

	httputil.JSON(w, http.StatusOK, domain.VendorScoreResponse{
		Vendor: *vendor,
		Score:  *score,
	})
}

func (h *Handlers) GetSyncStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	status, err := h.store.GetLatestSyncStatus(r.Context(), tenantID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "no sync status found"})
		return
	}

	httputil.JSON(w, http.StatusOK, status)
}

func (h *Handlers) GetScoreSummary(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	summary, err := h.store.GetScoreSummary(r.Context(), tenantID)
	if err != nil {
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get score summary"})
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
