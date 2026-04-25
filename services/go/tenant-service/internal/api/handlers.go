package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	kmsTypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/tenant-service/internal/domain"
	"github.com/complai/complai/services/go/tenant-service/internal/store"
)

// KMSClient is the subset of the AWS KMS client used by Handlers.
type KMSClient interface {
	CreateKey(ctx context.Context, params *kms.CreateKeyInput, optFns ...func(*kms.Options)) (*kms.CreateKeyOutput, error)
	CreateAlias(ctx context.Context, params *kms.CreateAliasInput, optFns ...func(*kms.Options)) (*kms.CreateAliasOutput, error)
}

type Handlers struct {
	store     store.Repository
	kmsClient KMSClient
}

func NewHandlers(s store.Repository, kmsClient KMSClient) *Handlers {
	return &Handlers{store: s, kmsClient: kmsClient}
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "tenant-service"})
}

func (h *Handlers) CreateTenant(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	t := &domain.Tenant{
		Name:     req.Name,
		Slug:     req.Slug,
		Tier:     req.Tier,
		Settings: "{}",
	}

	if err := h.store.CreateTenant(r.Context(), t); err != nil {
		log.Error().Err(err).Msg("create tenant failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "create failed"})
		return
	}

	if h.kmsClient != nil {
		keyARN, err := h.createKMSKey(r.Context(), t.ID)
		if err != nil {
			log.Error().Err(err).Msg("KMS key creation failed")
		} else {
			t.KMSKeyARN = &keyARN
			if err := h.store.UpdateTenantKMSKey(r.Context(), t.ID, keyARN); err != nil {
				log.Error().Err(err).Msg("store KMS ARN failed")
			}
		}
	}

	httputil.JSON(w, http.StatusCreated, t)
}

func (h *Handlers) GetTenant(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(r.PathValue("tenantID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid tenant_id"})
		return
	}

	t, err := h.store.GetTenant(r.Context(), tenantID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "tenant not found"})
		return
	}

	httputil.JSON(w, http.StatusOK, t)
}

func (h *Handlers) ListTenants(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	tenants, err := h.store.ListTenants(r.Context(), tenantID)
	if err != nil {
		log.Error().Err(err).Msg("list tenants failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if tenants == nil {
		tenants = []domain.Tenant{}
	}

	httputil.JSON(w, http.StatusOK, tenants)
}

func (h *Handlers) SuspendTenant(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(r.PathValue("tenantID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid tenant_id"})
		return
	}

	if err := h.store.UpdateTenantStatus(r.Context(), tenantID, "suspended"); err != nil {
		log.Error().Err(err).Msg("suspend tenant failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]string{"status": "suspended"})
}

func (h *Handlers) ReactivateTenant(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(r.PathValue("tenantID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid tenant_id"})
		return
	}

	if err := h.store.UpdateTenantStatus(r.Context(), tenantID, "active"); err != nil {
		log.Error().Err(err).Msg("reactivate failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]string{"status": "active"})
}

func (h *Handlers) GetHierarchy(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(r.PathValue("tenantID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid tenant_id"})
		return
	}

	hierarchy, err := h.store.GetHierarchy(r.Context(), tenantID)
	if err != nil {
		log.Error().Err(err).Msg("get hierarchy failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	httputil.JSON(w, http.StatusOK, hierarchy)
}

func (h *Handlers) CreatePAN(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(r.PathValue("tenantID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid tenant_id"})
		return
	}

	var req domain.CreatePANRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	p := &domain.TenantPAN{PAN: req.PAN, EntityName: req.EntityName, PANType: req.PANType}
	if err := h.store.CreatePAN(r.Context(), tenantID, p); err != nil {
		log.Error().Err(err).Msg("create PAN failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "create failed"})
		return
	}

	httputil.JSON(w, http.StatusCreated, p)
}

func (h *Handlers) CreateGSTIN(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(r.PathValue("tenantID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid tenant_id"})
		return
	}
	panID, err := uuid.Parse(r.PathValue("panID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid pan_id"})
		return
	}

	var req domain.CreateGSTINRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	g := &domain.TenantGSTIN{PANID: panID, GSTIN: req.GSTIN, TradeName: req.TradeName, StateCode: req.StateCode, RegistrationType: req.RegistrationType}
	if err := h.store.CreateGSTIN(r.Context(), tenantID, g); err != nil {
		log.Error().Err(err).Msg("create GSTIN failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "create failed"})
		return
	}

	httputil.JSON(w, http.StatusCreated, g)
}

func (h *Handlers) CreateTAN(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(r.PathValue("tenantID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid tenant_id"})
		return
	}
	panID, err := uuid.Parse(r.PathValue("panID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid pan_id"})
		return
	}

	var req domain.CreateTANRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	t := &domain.TenantTAN{PANID: panID, TAN: req.TAN, DeductorName: req.DeductorName}
	if err := h.store.CreateTAN(r.Context(), tenantID, t); err != nil {
		log.Error().Err(err).Msg("create TAN failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "create failed"})
		return
	}

	httputil.JSON(w, http.StatusCreated, t)
}

func (h *Handlers) createKMSKey(ctx context.Context, tenantID uuid.UUID) (string, error) {
	desc := fmt.Sprintf("Complai tenant CMK: %s", tenantID)
	out, err := h.kmsClient.CreateKey(ctx, &kms.CreateKeyInput{
		Description: &desc,
		Tags: []kmsTypes.Tag{
			{TagKey: strPtr("tenant_id"), TagValue: strPtr(tenantID.String())},
			{TagKey: strPtr("managed_by"), TagValue: strPtr("complai")},
		},
	})
	if err != nil {
		return "", fmt.Errorf("kms create key: %w", err)
	}

	alias := fmt.Sprintf("alias/complai-tenant-%s", tenantID)
	_, err = h.kmsClient.CreateAlias(ctx, &kms.CreateAliasInput{
		AliasName:   &alias,
		TargetKeyId: out.KeyMetadata.KeyId,
	})
	if err != nil {
		log.Warn().Err(err).Msg("kms create alias failed (non-fatal)")
	}

	return *out.KeyMetadata.Arn, nil
}

func strPtr(s string) *string { return &s }

func tenantIDFromRequest(r *http.Request) (uuid.UUID, error) {
	h := r.Header.Get("X-Tenant-Id")
	if h == "" {
		return uuid.Nil, fmt.Errorf("missing X-Tenant-Id header")
	}
	return uuid.Parse(h)
}
