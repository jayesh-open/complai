package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/user-role-service/internal/domain"
	"github.com/complai/complai/services/go/user-role-service/internal/store"
)

type Handlers struct {
	store store.Repository
}

func NewHandlers(s store.Repository) *Handlers {
	return &Handlers{store: s}
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "user-role-service"})
}

func (h *Handlers) ListRoles(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	roles, err := h.store.ListRoles(r.Context(), tenantID)
	if err != nil {
		log.Error().Err(err).Msg("list roles failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if roles == nil {
		roles = []domain.Role{}
	}
	httputil.JSON(w, http.StatusOK, roles)
}

func (h *Handlers) CreateRole(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	role := &domain.Role{Name: req.Name, DisplayName: req.DisplayName, Description: req.Description}
	if err := h.store.CreateRole(r.Context(), tenantID, role); err != nil {
		log.Error().Err(err).Msg("create role failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "create failed"})
		return
	}
	httputil.JSON(w, http.StatusCreated, role)
}

func (h *Handlers) AssignPermissions(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	roleID, err := uuid.Parse(r.PathValue("roleID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid role_id"})
		return
	}

	var req domain.AssignPermissionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	for _, pid := range req.PermissionIDs {
		if err := h.store.AssignPermissionToRole(r.Context(), tenantID, roleID, pid); err != nil {
			log.Error().Err(err).Msg("assign permission failed")
			httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "assign failed"})
			return
		}
	}
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "permissions_assigned"})
}

func (h *Handlers) AssignRole(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	userID, err := uuid.Parse(r.PathValue("userID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
		return
	}

	var req domain.AssignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	if err := h.store.AssignRoleToUser(r.Context(), tenantID, userID, req.RoleID, nil); err != nil {
		log.Error().Err(err).Msg("assign role failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "assign failed"})
		return
	}
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "role_assigned"})
}

func (h *Handlers) PolicyCheck(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.PolicyCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	resp, err := h.store.CheckPolicy(r.Context(), tenantID, req.UserID, req.Resource, req.Action)
	if err != nil {
		log.Error().Err(err).Msg("policy check failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	httputil.JSON(w, http.StatusOK, resp)
}

func (h *Handlers) CreateApproval(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	requestedBy, err := userIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.CreateApprovalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	a := &domain.ApprovalWorkflow{
		ResourceType: req.ResourceType,
		ActionType:   req.ActionType,
		RequestedBy:  requestedBy,
		Payload:      req.Payload,
	}
	if a.Payload == "" {
		a.Payload = "{}"
	}

	if err := h.store.CreateApproval(r.Context(), tenantID, a); err != nil {
		log.Error().Err(err).Msg("create approval failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "create failed"})
		return
	}
	httputil.JSON(w, http.StatusCreated, a)
}

func (h *Handlers) DecideApproval(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	decidedBy, err := userIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	approvalID, err := uuid.Parse(r.PathValue("approvalID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid approval_id"})
		return
	}

	approval, err := h.store.GetApproval(r.Context(), tenantID, approvalID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "approval not found"})
		return
	}

	if approval.RequestedBy == decidedBy {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "self_approval_denied",
			"message": "Cannot approve your own request (maker-checker)",
		})
		return
	}

	var req domain.DecideApprovalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	if err := h.store.DecideApproval(r.Context(), tenantID, approvalID, decidedBy, req.Decision, req.Reason); err != nil {
		log.Error().Err(err).Msg("decide approval failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "decide failed"})
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]string{"status": req.Decision, "approval_id": approvalID.String()})
}

func (h *Handlers) ListApprovals(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	approvals, err := h.store.ListPendingApprovals(r.Context(), tenantID)
	if err != nil {
		log.Error().Err(err).Msg("list approvals failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if approvals == nil {
		approvals = []domain.ApprovalWorkflow{}
	}
	httputil.JSON(w, http.StatusOK, approvals)
}

func (h *Handlers) ListTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := h.store.GetRoleTemplates(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("list templates failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if templates == nil {
		templates = []domain.RoleTemplate{}
	}
	httputil.JSON(w, http.StatusOK, templates)
}

func tenantIDFromRequest(r *http.Request) (uuid.UUID, error) {
	h := r.Header.Get("X-Tenant-Id")
	if h == "" {
		return uuid.Nil, fmt.Errorf("missing X-Tenant-Id header")
	}
	return uuid.Parse(h)
}

func userIDFromRequest(r *http.Request) (uuid.UUID, error) {
	h := r.Header.Get("X-User-Id")
	if h == "" {
		return uuid.Nil, fmt.Errorf("missing X-User-Id header")
	}
	return uuid.Parse(h)
}
