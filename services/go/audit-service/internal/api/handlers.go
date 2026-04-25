package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/audit-service/internal/domain"
	"github.com/complai/complai/services/go/audit-service/internal/store"
)

type Handlers struct {
	store store.Repository
}

func NewHandlers(s store.Repository) *Handlers {
	return &Handlers{store: s}
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "audit-service"})
}

func (h *Handlers) CreateEvent(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.CreateAuditEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	status := req.Status
	if status == "" {
		status = "success"
	}

	e := &domain.AuditEvent{
		UserID:       req.UserID,
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
		Action:       req.Action,
		OldValue:     req.OldValue,
		NewValue:     req.NewValue,
		Status:       status,
		ErrorMessage: req.ErrorMessage,
		IPAddress:    req.IPAddress,
		UserAgent:    req.UserAgent,
		TraceID:      req.TraceID,
	}

	if err := h.store.CreateEvent(r.Context(), tenantID, e); err != nil {
		log.Error().Err(err).Msg("create audit event failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "create failed"})
		return
	}

	httputil.JSON(w, http.StatusCreated, e)
}

func (h *Handlers) ListEvents(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	params := domain.QueryParams{
		ResourceType: r.URL.Query().Get("resource_type"),
		Action:       r.URL.Query().Get("action"),
	}

	if v := r.URL.Query().Get("date_from"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid date_from"})
			return
		}
		params.DateFrom = &t
	}
	if v := r.URL.Query().Get("date_to"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid date_to"})
			return
		}
		params.DateTo = &t
	}
	if v := r.URL.Query().Get("limit"); v != "" {
		l, err := strconv.Atoi(v)
		if err != nil {
			httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid limit"})
			return
		}
		params.Limit = l
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		o, err := strconv.Atoi(v)
		if err != nil {
			httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid offset"})
			return
		}
		params.Offset = o
	}

	events, err := h.store.ListEvents(r.Context(), tenantID, params)
	if err != nil {
		log.Error().Err(err).Msg("list audit events failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if events == nil {
		events = []domain.AuditEvent{}
	}

	httputil.JSON(w, http.StatusOK, events)
}

func (h *Handlers) GetEvent(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	eventID, err := uuid.Parse(r.PathValue("eventID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid event_id"})
		return
	}

	e, err := h.store.GetEvent(r.Context(), tenantID, eventID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "event not found"})
		return
	}

	httputil.JSON(w, http.StatusOK, e)
}

type computeMerkleRequest struct {
	HourBucket time.Time `json:"hour_bucket"`
}

func (h *Handlers) ComputeMerkleHash(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req computeMerkleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	hourBucket := req.HourBucket.Truncate(time.Hour)

	events, err := h.store.GetEventsForHour(r.Context(), tenantID, hourBucket)
	if err != nil {
		log.Error().Err(err).Msg("get events for hour failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	if len(events) == 0 {
		httputil.JSON(w, http.StatusOK, map[string]string{"message": "no events for this hour"})
		return
	}

	// Get previous hash from the latest merkle chain entry
	previousHash := ""
	latest, err := h.store.GetLatestMerkleChain(r.Context(), tenantID)
	if err == nil && latest != nil {
		previousHash = latest.ComputedHash
	}

	// Build hash payload
	payload := buildHashPayload(events)
	computedHash := computeHash(previousHash, events)

	chain := &domain.MerkleChain{
		HourBucket:   hourBucket,
		EventCount:   len(events),
		HashPayload:  payload,
		PreviousHash: previousHash,
		ComputedHash: computedHash,
	}

	if err := h.store.CreateMerkleChain(r.Context(), tenantID, chain); err != nil {
		log.Error().Err(err).Msg("create merkle chain failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "create merkle chain failed"})
		return
	}

	httputil.JSON(w, http.StatusCreated, chain)
}

func (h *Handlers) IntegrityCheck(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	dateFromStr := r.URL.Query().Get("date_from")
	dateToStr := r.URL.Query().Get("date_to")
	if dateFromStr == "" || dateToStr == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "date_from and date_to are required"})
		return
	}

	dateFrom, err := time.Parse(time.RFC3339, dateFromStr)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid date_from"})
		return
	}
	dateTo, err := time.Parse(time.RFC3339, dateToStr)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid date_to"})
		return
	}

	chains, err := h.store.GetMerkleChains(r.Context(), tenantID, dateFrom, dateTo)
	if err != nil {
		log.Error().Err(err).Msg("get merkle chains failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	result := domain.IntegrityCheckResult{
		Valid:       true,
		CheckedFrom: dateFrom,
		CheckedTo:   dateTo,
		ChainLength: len(chains),
		Message:     "all chain entries verified",
	}

	if len(chains) == 0 {
		result.Message = "no chain entries found in range"
		httputil.JSON(w, http.StatusOK, result)
		return
	}

	for i, chain := range chains {
		// Re-fetch events for this hour bucket and recompute hash
		events, err := h.store.GetEventsForHour(r.Context(), tenantID, chain.HourBucket)
		if err != nil {
			log.Error().Err(err).Msg("get events for integrity check failed")
			httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
			return
		}

		recomputedHash := computeHash(chain.PreviousHash, events)
		if recomputedHash != chain.ComputedHash {
			brokenAt := chain.HourBucket
			result.Valid = false
			result.BrokenAt = &brokenAt
			result.Message = fmt.Sprintf("hash mismatch at hour bucket %s", chain.HourBucket.Format(time.RFC3339))
			httputil.JSON(w, http.StatusOK, result)
			return
		}

		// Verify chain linkage (skip first entry — its previous_hash comes from before the range)
		if i > 0 {
			if chain.PreviousHash != chains[i-1].ComputedHash {
				brokenAt := chain.HourBucket
				result.Valid = false
				result.BrokenAt = &brokenAt
				result.Message = fmt.Sprintf("chain linkage broken at hour bucket %s", chain.HourBucket.Format(time.RFC3339))
				httputil.JSON(w, http.StatusOK, result)
				return
			}
		}
	}

	httputil.JSON(w, http.StatusOK, result)
}

func computeHash(previousHash string, events []domain.AuditEvent) string {
	sort.Slice(events, func(i, j int) bool {
		return events[i].CreatedAt.Before(events[j].CreatedAt)
	})
	var payload strings.Builder
	for _, e := range events {
		payload.WriteString(e.ID.String())
		payload.WriteString(e.CreatedAt.Format(time.RFC3339Nano))
	}
	h := sha256.New()
	h.Write([]byte(previousHash + payload.String()))
	return hex.EncodeToString(h.Sum(nil))
}

func buildHashPayload(events []domain.AuditEvent) string {
	sort.Slice(events, func(i, j int) bool {
		return events[i].CreatedAt.Before(events[j].CreatedAt)
	})
	var payload strings.Builder
	for _, e := range events {
		payload.WriteString(e.ID.String())
		payload.WriteString(e.CreatedAt.Format(time.RFC3339Nano))
	}
	return payload.String()
}

func tenantIDFromRequest(r *http.Request) (uuid.UUID, error) {
	h := r.Header.Get("X-Tenant-Id")
	if h == "" {
		return uuid.Nil, fmt.Errorf("missing X-Tenant-Id header")
	}
	return uuid.Parse(h)
}
