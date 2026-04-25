package api

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/aura-gateway-service/internal/provider"
)

type Handlers struct {
	provider provider.AuraProvider
}

func NewHandlers(p provider.AuraProvider) *Handlers {
	return &Handlers{provider: p}
}

func (h *Handlers) Health(w http.ResponseWriter, _ *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "aura-gateway-service"})
}

func (h *Handlers) ListARInvoices(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	gstin := r.URL.Query().Get("gstin")
	if gstin == "" {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "missing gstin query parameter"})
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "042026"
	}

	resp, err := h.provider.ListARInvoices(r.Context(), tenantID, gstin, period)
	if err != nil {
		log.Error().Err(err).Str("tenant_id", tenantID.String()).Msg("list AR invoices failed")
		httputil.JSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	httputil.JSON(w, http.StatusOK, resp)
}

func tenantIDFromRequest(r *http.Request) (uuid.UUID, error) {
	h := r.Header.Get("X-Tenant-Id")
	if h == "" {
		return uuid.Nil, fmt.Errorf("missing X-Tenant-Id header")
	}
	return uuid.Parse(h)
}
