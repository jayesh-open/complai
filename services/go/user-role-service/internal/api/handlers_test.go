package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/user-role-service/internal/domain"
	"github.com/complai/complai/services/go/user-role-service/internal/store"
)

// Ensure the mock implements the interface at compile time.
var _ store.Repository = (*mockStore)(nil)

// mockStore is a test double implementing store.Repository.
type mockStore struct {
	createRoleFn            func(ctx context.Context, tenantID uuid.UUID, r *domain.Role) error
	getRoleFn               func(ctx context.Context, tenantID, roleID uuid.UUID) (*domain.Role, error)
	listRolesFn             func(ctx context.Context, tenantID uuid.UUID) ([]domain.Role, error)
	createPermissionFn      func(ctx context.Context, tenantID uuid.UUID, p *domain.Permission) error
	assignPermissionToRoleFn func(ctx context.Context, tenantID, roleID, permissionID uuid.UUID) error
	assignRoleToUserFn      func(ctx context.Context, tenantID, userID, roleID uuid.UUID, assignedBy *uuid.UUID) error
	getUserPermissionsFn    func(ctx context.Context, tenantID, userID uuid.UUID) ([]domain.Permission, error)
	checkPolicyFn           func(ctx context.Context, tenantID, userID uuid.UUID, resource, action string) (*domain.PolicyCheckResponse, error)
	createRoleTemplateFn    func(ctx context.Context, rt *domain.RoleTemplate) error
	getRoleTemplatesFn      func(ctx context.Context) ([]domain.RoleTemplate, error)
	createApprovalFn        func(ctx context.Context, tenantID uuid.UUID, a *domain.ApprovalWorkflow) error
	getApprovalFn           func(ctx context.Context, tenantID, approvalID uuid.UUID) (*domain.ApprovalWorkflow, error)
	decideApprovalFn        func(ctx context.Context, tenantID, approvalID, decidedBy uuid.UUID, decision string, reason *string) error
	listPendingApprovalsFn  func(ctx context.Context, tenantID uuid.UUID) ([]domain.ApprovalWorkflow, error)
}

func (m *mockStore) CreateRole(ctx context.Context, tenantID uuid.UUID, r *domain.Role) error {
	if m.createRoleFn != nil {
		return m.createRoleFn(ctx, tenantID, r)
	}
	return nil
}

func (m *mockStore) GetRole(ctx context.Context, tenantID, roleID uuid.UUID) (*domain.Role, error) {
	if m.getRoleFn != nil {
		return m.getRoleFn(ctx, tenantID, roleID)
	}
	return nil, nil
}

func (m *mockStore) ListRoles(ctx context.Context, tenantID uuid.UUID) ([]domain.Role, error) {
	if m.listRolesFn != nil {
		return m.listRolesFn(ctx, tenantID)
	}
	return nil, nil
}

func (m *mockStore) CreatePermission(ctx context.Context, tenantID uuid.UUID, p *domain.Permission) error {
	if m.createPermissionFn != nil {
		return m.createPermissionFn(ctx, tenantID, p)
	}
	return nil
}

func (m *mockStore) AssignPermissionToRole(ctx context.Context, tenantID, roleID, permissionID uuid.UUID) error {
	if m.assignPermissionToRoleFn != nil {
		return m.assignPermissionToRoleFn(ctx, tenantID, roleID, permissionID)
	}
	return nil
}

func (m *mockStore) AssignRoleToUser(ctx context.Context, tenantID, userID, roleID uuid.UUID, assignedBy *uuid.UUID) error {
	if m.assignRoleToUserFn != nil {
		return m.assignRoleToUserFn(ctx, tenantID, userID, roleID, assignedBy)
	}
	return nil
}

func (m *mockStore) GetUserPermissions(ctx context.Context, tenantID, userID uuid.UUID) ([]domain.Permission, error) {
	if m.getUserPermissionsFn != nil {
		return m.getUserPermissionsFn(ctx, tenantID, userID)
	}
	return nil, nil
}

func (m *mockStore) CheckPolicy(ctx context.Context, tenantID, userID uuid.UUID, resource, action string) (*domain.PolicyCheckResponse, error) {
	if m.checkPolicyFn != nil {
		return m.checkPolicyFn(ctx, tenantID, userID, resource, action)
	}
	return &domain.PolicyCheckResponse{Allow: false, Reasons: []string{"no mock"}}, nil
}

func (m *mockStore) CreateRoleTemplate(ctx context.Context, rt *domain.RoleTemplate) error {
	if m.createRoleTemplateFn != nil {
		return m.createRoleTemplateFn(ctx, rt)
	}
	return nil
}

func (m *mockStore) GetRoleTemplates(ctx context.Context) ([]domain.RoleTemplate, error) {
	if m.getRoleTemplatesFn != nil {
		return m.getRoleTemplatesFn(ctx)
	}
	return nil, nil
}

func (m *mockStore) CreateApproval(ctx context.Context, tenantID uuid.UUID, a *domain.ApprovalWorkflow) error {
	if m.createApprovalFn != nil {
		return m.createApprovalFn(ctx, tenantID, a)
	}
	return nil
}

func (m *mockStore) GetApproval(ctx context.Context, tenantID, approvalID uuid.UUID) (*domain.ApprovalWorkflow, error) {
	if m.getApprovalFn != nil {
		return m.getApprovalFn(ctx, tenantID, approvalID)
	}
	return nil, nil
}

func (m *mockStore) DecideApproval(ctx context.Context, tenantID, approvalID, decidedBy uuid.UUID, decision string, reason *string) error {
	if m.decideApprovalFn != nil {
		return m.decideApprovalFn(ctx, tenantID, approvalID, decidedBy, decision, reason)
	}
	return nil
}

func (m *mockStore) ListPendingApprovals(ctx context.Context, tenantID uuid.UUID) ([]domain.ApprovalWorkflow, error) {
	if m.listPendingApprovalsFn != nil {
		return m.listPendingApprovalsFn(ctx, tenantID)
	}
	return nil, nil
}

// ---------- helpers ----------

func newTestHandlers(ms *mockStore) *Handlers {
	return NewHandlers(ms)
}

// withChiURLParams sets path parameters using Go 1.22's native SetPathValue,
// which is what chi v5.2.1 uses internally.
func withChiURLParams(r *http.Request, params map[string]string) *http.Request {
	for k, v := range params {
		r.SetPathValue(k, v)
	}
	return r
}

// parseSuccessResp parses an httputil.SuccessResponse envelope, returning the raw "data" JSON.
func parseSuccessResp(t *testing.T, body []byte) json.RawMessage {
	t.Helper()
	var env httputil.SuccessResponse
	err := json.Unmarshal(body, &env)
	require.NoError(t, err, "failed to unmarshal SuccessResponse envelope")
	raw, err := json.Marshal(env.Data)
	require.NoError(t, err)
	return raw
}

// ---------- Health ----------

func TestHealth(t *testing.T) {
	h := newTestHandlers(&mockStore{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	h.Health(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	data := parseSuccessResp(t, rec.Body.Bytes())
	var m map[string]string
	require.NoError(t, json.Unmarshal(data, &m))
	assert.Equal(t, "ok", m["status"])
	assert.Equal(t, "user-role-service", m["service"])
}

// ---------- ListRoles ----------

func TestListRoles_Success(t *testing.T) {
	tenantID := uuid.New()
	roles := []domain.Role{
		{ID: uuid.New(), TenantID: tenantID, Name: "admin", DisplayName: "Admin", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	ms := &mockStore{
		listRolesFn: func(_ context.Context, tid uuid.UUID) ([]domain.Role, error) {
			assert.Equal(t, tenantID, tid)
			return roles, nil
		},
	}
	h := newTestHandlers(ms)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/roles", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())

	h.ListRoles(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	data := parseSuccessResp(t, rec.Body.Bytes())
	var got []domain.Role
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Len(t, got, 1)
	assert.Equal(t, "admin", got[0].Name)
}

func TestListRoles_MissingTenantHeader(t *testing.T) {
	h := newTestHandlers(&mockStore{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/roles", nil)
	// no X-Tenant-Id header

	h.ListRoles(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListRoles_EmptyList(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listRolesFn: func(_ context.Context, _ uuid.UUID) ([]domain.Role, error) {
			return nil, nil // handler converts nil to []
		},
	}
	h := newTestHandlers(ms)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/roles", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())

	h.ListRoles(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	data := parseSuccessResp(t, rec.Body.Bytes())
	var got []domain.Role
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Len(t, got, 0)
}

func TestListRoles_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listRolesFn: func(_ context.Context, _ uuid.UUID) ([]domain.Role, error) {
			return nil, fmt.Errorf("db error")
		},
	}
	h := newTestHandlers(ms)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/roles", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())

	h.ListRoles(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------- CreateRole ----------

func TestCreateRole_Success(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		createRoleFn: func(_ context.Context, tid uuid.UUID, r *domain.Role) error {
			assert.Equal(t, tenantID, tid)
			assert.Equal(t, "viewer", r.Name)
			r.ID = uuid.New()
			r.TenantID = tid
			r.CreatedAt = time.Now()
			r.UpdatedAt = time.Now()
			return nil
		},
	}
	h := newTestHandlers(ms)

	body := `{"name":"viewer","display_name":"Viewer"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/roles", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.Header.Set("Content-Type", "application/json")

	h.CreateRole(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	data := parseSuccessResp(t, rec.Body.Bytes())
	var got domain.Role
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, "viewer", got.Name)
	assert.Equal(t, "Viewer", got.DisplayName)
}

func TestCreateRole_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	h := newTestHandlers(&mockStore{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/roles", bytes.NewBufferString("not json"))
	req.Header.Set("X-Tenant-Id", tenantID.String())

	h.CreateRole(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateRole_MissingTenantHeader(t *testing.T) {
	h := newTestHandlers(&mockStore{})
	body := `{"name":"viewer","display_name":"Viewer"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/roles", bytes.NewBufferString(body))

	h.CreateRole(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateRole_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		createRoleFn: func(_ context.Context, _ uuid.UUID, _ *domain.Role) error {
			return fmt.Errorf("db error")
		},
	}
	h := newTestHandlers(ms)
	body := `{"name":"viewer","display_name":"Viewer"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/roles", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())

	h.CreateRole(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------- AssignPermissions ----------

func TestAssignPermissions_Success(t *testing.T) {
	tenantID := uuid.New()
	roleID := uuid.New()
	perm1 := uuid.New()
	perm2 := uuid.New()

	var assignedIDs []uuid.UUID
	ms := &mockStore{
		assignPermissionToRoleFn: func(_ context.Context, tid, rid, pid uuid.UUID) error {
			assert.Equal(t, tenantID, tid)
			assert.Equal(t, roleID, rid)
			assignedIDs = append(assignedIDs, pid)
			return nil
		},
	}
	h := newTestHandlers(ms)

	body := fmt.Sprintf(`{"permission_ids":["%s","%s"]}`, perm1, perm2)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/roles/"+roleID.String()+"/permissions", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req = withChiURLParams(req, map[string]string{"roleID": roleID.String()})

	h.AssignPermissions(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Len(t, assignedIDs, 2)
	assert.Contains(t, assignedIDs, perm1)
	assert.Contains(t, assignedIDs, perm2)
}

func TestAssignPermissions_InvalidRoleID(t *testing.T) {
	tenantID := uuid.New()
	h := newTestHandlers(&mockStore{})
	body := `{"permission_ids":["` + uuid.New().String() + `"]}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/roles/not-a-uuid/permissions", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req = withChiURLParams(req, map[string]string{"roleID": "not-a-uuid"})

	h.AssignPermissions(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAssignPermissions_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	roleID := uuid.New()
	h := newTestHandlers(&mockStore{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/roles/"+roleID.String()+"/permissions", bytes.NewBufferString("bad"))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req = withChiURLParams(req, map[string]string{"roleID": roleID.String()})

	h.AssignPermissions(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAssignPermissions_StoreError(t *testing.T) {
	tenantID := uuid.New()
	roleID := uuid.New()
	ms := &mockStore{
		assignPermissionToRoleFn: func(_ context.Context, _, _, _ uuid.UUID) error {
			return fmt.Errorf("db error")
		},
	}
	h := newTestHandlers(ms)
	body := `{"permission_ids":["` + uuid.New().String() + `"]}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/roles/"+roleID.String()+"/permissions", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req = withChiURLParams(req, map[string]string{"roleID": roleID.String()})

	h.AssignPermissions(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAssignPermissions_MissingTenantHeader(t *testing.T) {
	roleID := uuid.New()
	h := newTestHandlers(&mockStore{})
	body := `{"permission_ids":["` + uuid.New().String() + `"]}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/roles/"+roleID.String()+"/permissions", bytes.NewBufferString(body))
	req = withChiURLParams(req, map[string]string{"roleID": roleID.String()})

	h.AssignPermissions(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------- AssignRole ----------

func TestAssignRole_Success(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	roleID := uuid.New()

	ms := &mockStore{
		assignRoleToUserFn: func(_ context.Context, tid, uid, rid uuid.UUID, ab *uuid.UUID) error {
			assert.Equal(t, tenantID, tid)
			assert.Equal(t, userID, uid)
			assert.Equal(t, roleID, rid)
			assert.Nil(t, ab)
			return nil
		},
	}
	h := newTestHandlers(ms)

	body := fmt.Sprintf(`{"role_id":"%s"}`, roleID)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/users/"+userID.String()+"/roles", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req = withChiURLParams(req, map[string]string{"userID": userID.String()})

	h.AssignRole(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	data := parseSuccessResp(t, rec.Body.Bytes())
	var m map[string]string
	require.NoError(t, json.Unmarshal(data, &m))
	assert.Equal(t, "role_assigned", m["status"])
}

func TestAssignRole_InvalidUserID(t *testing.T) {
	tenantID := uuid.New()
	h := newTestHandlers(&mockStore{})
	body := `{"role_id":"` + uuid.New().String() + `"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/users/bad-id/roles", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req = withChiURLParams(req, map[string]string{"userID": "bad-id"})

	h.AssignRole(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAssignRole_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	h := newTestHandlers(&mockStore{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/users/"+userID.String()+"/roles", bytes.NewBufferString("nope"))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req = withChiURLParams(req, map[string]string{"userID": userID.String()})

	h.AssignRole(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAssignRole_StoreError(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	roleID := uuid.New()
	ms := &mockStore{
		assignRoleToUserFn: func(_ context.Context, _, _, _ uuid.UUID, _ *uuid.UUID) error {
			return fmt.Errorf("db error")
		},
	}
	h := newTestHandlers(ms)
	body := fmt.Sprintf(`{"role_id":"%s"}`, roleID)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/users/"+userID.String()+"/roles", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req = withChiURLParams(req, map[string]string{"userID": userID.String()})

	h.AssignRole(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAssignRole_MissingTenantHeader(t *testing.T) {
	userID := uuid.New()
	h := newTestHandlers(&mockStore{})
	body := `{"role_id":"` + uuid.New().String() + `"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/users/"+userID.String()+"/roles", bytes.NewBufferString(body))
	req = withChiURLParams(req, map[string]string{"userID": userID.String()})

	h.AssignRole(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------- PolicyCheck ----------

func TestPolicyCheck_Success_Allow(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	ms := &mockStore{
		checkPolicyFn: func(_ context.Context, tid, uid uuid.UUID, resource, action string) (*domain.PolicyCheckResponse, error) {
			assert.Equal(t, tenantID, tid)
			assert.Equal(t, userID, uid)
			assert.Equal(t, "invoices", resource)
			assert.Equal(t, "read", action)
			return &domain.PolicyCheckResponse{Allow: true, Reasons: []string{"permission granted via role"}}, nil
		},
	}
	h := newTestHandlers(ms)

	body := fmt.Sprintf(`{"user_id":"%s","resource":"invoices","action":"read"}`, userID)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/policy/check", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())

	h.PolicyCheck(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	data := parseSuccessResp(t, rec.Body.Bytes())
	var resp domain.PolicyCheckResponse
	require.NoError(t, json.Unmarshal(data, &resp))
	assert.True(t, resp.Allow)
}

func TestPolicyCheck_Success_Deny(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	ms := &mockStore{
		checkPolicyFn: func(_ context.Context, _, _ uuid.UUID, _, _ string) (*domain.PolicyCheckResponse, error) {
			return &domain.PolicyCheckResponse{Allow: false, Reasons: []string{"no matching permission found"}}, nil
		},
	}
	h := newTestHandlers(ms)

	body := fmt.Sprintf(`{"user_id":"%s","resource":"invoices","action":"delete"}`, userID)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/policy/check", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())

	h.PolicyCheck(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	data := parseSuccessResp(t, rec.Body.Bytes())
	var resp domain.PolicyCheckResponse
	require.NoError(t, json.Unmarshal(data, &resp))
	assert.False(t, resp.Allow)
}

func TestPolicyCheck_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	h := newTestHandlers(&mockStore{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/policy/check", bytes.NewBufferString("{bad"))
	req.Header.Set("X-Tenant-Id", tenantID.String())

	h.PolicyCheck(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPolicyCheck_MissingTenantHeader(t *testing.T) {
	h := newTestHandlers(&mockStore{})
	body := `{"user_id":"` + uuid.New().String() + `","resource":"x","action":"y"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/policy/check", bytes.NewBufferString(body))

	h.PolicyCheck(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPolicyCheck_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		checkPolicyFn: func(_ context.Context, _, _ uuid.UUID, _, _ string) (*domain.PolicyCheckResponse, error) {
			return nil, fmt.Errorf("db error")
		},
	}
	h := newTestHandlers(ms)
	body := `{"user_id":"` + uuid.New().String() + `","resource":"x","action":"y"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/policy/check", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())

	h.PolicyCheck(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------- CreateApproval ----------

func TestCreateApproval_Success(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	ms := &mockStore{
		createApprovalFn: func(_ context.Context, tid uuid.UUID, a *domain.ApprovalWorkflow) error {
			assert.Equal(t, tenantID, tid)
			assert.Equal(t, userID, a.RequestedBy)
			assert.Equal(t, "role_change", a.ResourceType)
			assert.Equal(t, "promote", a.ActionType)
			a.ID = uuid.New()
			a.TenantID = tid
			a.Status = "pending_approval"
			a.CreatedAt = time.Now()
			return nil
		},
	}
	h := newTestHandlers(ms)

	body := `{"resource_type":"role_change","action_type":"promote","payload":"{\"role\":\"admin\"}"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/approvals", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.Header.Set("X-User-Id", userID.String())

	h.CreateApproval(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	data := parseSuccessResp(t, rec.Body.Bytes())
	var got domain.ApprovalWorkflow
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, "pending_approval", got.Status)
}

func TestCreateApproval_EmptyPayloadDefaultsToJSON(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	var capturedPayload string
	ms := &mockStore{
		createApprovalFn: func(_ context.Context, _ uuid.UUID, a *domain.ApprovalWorkflow) error {
			capturedPayload = a.Payload
			a.ID = uuid.New()
			a.Status = "pending_approval"
			a.CreatedAt = time.Now()
			return nil
		},
	}
	h := newTestHandlers(ms)

	body := `{"resource_type":"role_change","action_type":"promote"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/approvals", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.Header.Set("X-User-Id", userID.String())

	h.CreateApproval(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "{}", capturedPayload)
}

func TestCreateApproval_MissingTenantHeader(t *testing.T) {
	userID := uuid.New()
	h := newTestHandlers(&mockStore{})
	body := `{"resource_type":"role_change","action_type":"promote"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/approvals", bytes.NewBufferString(body))
	req.Header.Set("X-User-Id", userID.String())

	h.CreateApproval(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateApproval_MissingUserHeader(t *testing.T) {
	tenantID := uuid.New()
	h := newTestHandlers(&mockStore{})
	body := `{"resource_type":"role_change","action_type":"promote"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/approvals", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())

	h.CreateApproval(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateApproval_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	h := newTestHandlers(&mockStore{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/approvals", bytes.NewBufferString("{bad"))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.Header.Set("X-User-Id", userID.String())

	h.CreateApproval(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateApproval_StoreError(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	ms := &mockStore{
		createApprovalFn: func(_ context.Context, _ uuid.UUID, _ *domain.ApprovalWorkflow) error {
			return fmt.Errorf("db error")
		},
	}
	h := newTestHandlers(ms)
	body := `{"resource_type":"role_change","action_type":"promote"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/approvals", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.Header.Set("X-User-Id", userID.String())

	h.CreateApproval(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------- DecideApproval ----------

func TestDecideApproval_ManagerApprovesAnalystRequest(t *testing.T) {
	tenantID := uuid.New()
	analystID := uuid.New()  // the person who requested the approval
	managerID := uuid.New()  // the person deciding (different person)
	approvalID := uuid.New()

	ms := &mockStore{
		getApprovalFn: func(_ context.Context, tid, aid uuid.UUID) (*domain.ApprovalWorkflow, error) {
			assert.Equal(t, tenantID, tid)
			assert.Equal(t, approvalID, aid)
			return &domain.ApprovalWorkflow{
				ID:           approvalID,
				TenantID:     tenantID,
				ResourceType: "role_change",
				ActionType:   "promote",
				Status:       "pending_approval",
				RequestedBy:  analystID, // requested by analyst
				Payload:      "{}",
				CreatedAt:    time.Now(),
			}, nil
		},
		decideApprovalFn: func(_ context.Context, tid, aid, decidedBy uuid.UUID, decision string, reason *string) error {
			assert.Equal(t, tenantID, tid)
			assert.Equal(t, approvalID, aid)
			assert.Equal(t, managerID, decidedBy)
			assert.Equal(t, "approved", decision)
			return nil
		},
	}
	h := newTestHandlers(ms)

	body := `{"decision":"approved","reason":"looks good"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/v1/approvals/"+approvalID.String(), bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.Header.Set("X-User-Id", managerID.String()) // manager decides
	req = withChiURLParams(req, map[string]string{"approvalID": approvalID.String()})

	h.DecideApproval(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	data := parseSuccessResp(t, rec.Body.Bytes())
	var m map[string]string
	require.NoError(t, json.Unmarshal(data, &m))
	assert.Equal(t, "approved", m["status"])
	assert.Equal(t, approvalID.String(), m["approval_id"])
}

func TestDecideApproval_SelfApprovalDenied(t *testing.T) {
	tenantID := uuid.New()
	analystID := uuid.New() // same person requests AND tries to approve
	approvalID := uuid.New()

	ms := &mockStore{
		getApprovalFn: func(_ context.Context, _, _ uuid.UUID) (*domain.ApprovalWorkflow, error) {
			return &domain.ApprovalWorkflow{
				ID:           approvalID,
				TenantID:     tenantID,
				ResourceType: "role_change",
				ActionType:   "promote",
				Status:       "pending_approval",
				RequestedBy:  analystID, // same as the decider
				Payload:      "{}",
				CreatedAt:    time.Now(),
			}, nil
		},
	}
	h := newTestHandlers(ms)

	body := `{"decision":"approved"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/v1/approvals/"+approvalID.String(), bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.Header.Set("X-User-Id", analystID.String()) // same person tries to approve
	req = withChiURLParams(req, map[string]string{"approvalID": approvalID.String()})

	h.DecideApproval(rec, req)

	// The self-approval 403 is NOT wrapped in httputil.JSON — it uses direct json.Encoder
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Equal(t, "application/json; charset=utf-8", rec.Header().Get("Content-Type"))

	var resp map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "self_approval_denied", resp["error"])
	assert.Equal(t, "Cannot approve your own request (maker-checker)", resp["message"])
}

func TestDecideApproval_InvalidApprovalID(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	h := newTestHandlers(&mockStore{})
	body := `{"decision":"approved"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/v1/approvals/not-a-uuid", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.Header.Set("X-User-Id", userID.String())
	req = withChiURLParams(req, map[string]string{"approvalID": "not-a-uuid"})

	h.DecideApproval(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDecideApproval_ApprovalNotFound(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	approvalID := uuid.New()
	ms := &mockStore{
		getApprovalFn: func(_ context.Context, _, _ uuid.UUID) (*domain.ApprovalWorkflow, error) {
			return nil, fmt.Errorf("not found")
		},
	}
	h := newTestHandlers(ms)
	body := `{"decision":"approved"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/v1/approvals/"+approvalID.String(), bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.Header.Set("X-User-Id", userID.String())
	req = withChiURLParams(req, map[string]string{"approvalID": approvalID.String()})

	h.DecideApproval(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDecideApproval_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	managerID := uuid.New()
	analystID := uuid.New()
	approvalID := uuid.New()

	ms := &mockStore{
		getApprovalFn: func(_ context.Context, _, _ uuid.UUID) (*domain.ApprovalWorkflow, error) {
			return &domain.ApprovalWorkflow{
				ID:          approvalID,
				TenantID:    tenantID,
				Status:      "pending_approval",
				RequestedBy: analystID,
				Payload:     "{}",
				CreatedAt:   time.Now(),
			}, nil
		},
	}
	h := newTestHandlers(ms)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/v1/approvals/"+approvalID.String(), bytes.NewBufferString("{bad"))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.Header.Set("X-User-Id", managerID.String())
	req = withChiURLParams(req, map[string]string{"approvalID": approvalID.String()})

	h.DecideApproval(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDecideApproval_StoreDecideError(t *testing.T) {
	tenantID := uuid.New()
	analystID := uuid.New()
	managerID := uuid.New()
	approvalID := uuid.New()

	ms := &mockStore{
		getApprovalFn: func(_ context.Context, _, _ uuid.UUID) (*domain.ApprovalWorkflow, error) {
			return &domain.ApprovalWorkflow{
				ID:          approvalID,
				TenantID:    tenantID,
				Status:      "pending_approval",
				RequestedBy: analystID,
				Payload:     "{}",
				CreatedAt:   time.Now(),
			}, nil
		},
		decideApprovalFn: func(_ context.Context, _, _, _ uuid.UUID, _ string, _ *string) error {
			return fmt.Errorf("db error")
		},
	}
	h := newTestHandlers(ms)

	body := `{"decision":"rejected","reason":"nope"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/v1/approvals/"+approvalID.String(), bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.Header.Set("X-User-Id", managerID.String())
	req = withChiURLParams(req, map[string]string{"approvalID": approvalID.String()})

	h.DecideApproval(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestDecideApproval_MissingTenantHeader(t *testing.T) {
	userID := uuid.New()
	approvalID := uuid.New()
	h := newTestHandlers(&mockStore{})
	body := `{"decision":"approved"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/v1/approvals/"+approvalID.String(), bytes.NewBufferString(body))
	req.Header.Set("X-User-Id", userID.String())
	req = withChiURLParams(req, map[string]string{"approvalID": approvalID.String()})

	h.DecideApproval(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDecideApproval_MissingUserHeader(t *testing.T) {
	tenantID := uuid.New()
	approvalID := uuid.New()
	h := newTestHandlers(&mockStore{})
	body := `{"decision":"approved"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/v1/approvals/"+approvalID.String(), bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req = withChiURLParams(req, map[string]string{"approvalID": approvalID.String()})

	h.DecideApproval(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDecideApproval_RejectionSuccess(t *testing.T) {
	tenantID := uuid.New()
	analystID := uuid.New()
	managerID := uuid.New()
	approvalID := uuid.New()

	reason := "not justified"
	ms := &mockStore{
		getApprovalFn: func(_ context.Context, _, _ uuid.UUID) (*domain.ApprovalWorkflow, error) {
			return &domain.ApprovalWorkflow{
				ID:          approvalID,
				TenantID:    tenantID,
				Status:      "pending_approval",
				RequestedBy: analystID,
				Payload:     "{}",
				CreatedAt:   time.Now(),
			}, nil
		},
		decideApprovalFn: func(_ context.Context, _, _, _ uuid.UUID, decision string, r *string) error {
			assert.Equal(t, "rejected", decision)
			assert.NotNil(t, r)
			assert.Equal(t, reason, *r)
			return nil
		},
	}
	h := newTestHandlers(ms)

	body := fmt.Sprintf(`{"decision":"rejected","reason":"%s"}`, reason)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/v1/approvals/"+approvalID.String(), bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.Header.Set("X-User-Id", managerID.String())
	req = withChiURLParams(req, map[string]string{"approvalID": approvalID.String()})

	h.DecideApproval(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	data := parseSuccessResp(t, rec.Body.Bytes())
	var m map[string]string
	require.NoError(t, json.Unmarshal(data, &m))
	assert.Equal(t, "rejected", m["status"])
}

// ---------- ListApprovals ----------

func TestListApprovals_Success(t *testing.T) {
	tenantID := uuid.New()
	approvals := []domain.ApprovalWorkflow{
		{
			ID:           uuid.New(),
			TenantID:     tenantID,
			ResourceType: "role_change",
			ActionType:   "promote",
			Status:       "pending_approval",
			RequestedBy:  uuid.New(),
			Payload:      "{}",
			CreatedAt:    time.Now(),
		},
	}
	ms := &mockStore{
		listPendingApprovalsFn: func(_ context.Context, tid uuid.UUID) ([]domain.ApprovalWorkflow, error) {
			assert.Equal(t, tenantID, tid)
			return approvals, nil
		},
	}
	h := newTestHandlers(ms)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/approvals", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())

	h.ListApprovals(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	data := parseSuccessResp(t, rec.Body.Bytes())
	var got []domain.ApprovalWorkflow
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Len(t, got, 1)
}

func TestListApprovals_EmptyList(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listPendingApprovalsFn: func(_ context.Context, _ uuid.UUID) ([]domain.ApprovalWorkflow, error) {
			return nil, nil
		},
	}
	h := newTestHandlers(ms)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/approvals", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())

	h.ListApprovals(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	data := parseSuccessResp(t, rec.Body.Bytes())
	var got []domain.ApprovalWorkflow
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Len(t, got, 0)
}

func TestListApprovals_MissingTenantHeader(t *testing.T) {
	h := newTestHandlers(&mockStore{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/approvals", nil)

	h.ListApprovals(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListApprovals_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listPendingApprovalsFn: func(_ context.Context, _ uuid.UUID) ([]domain.ApprovalWorkflow, error) {
			return nil, fmt.Errorf("db error")
		},
	}
	h := newTestHandlers(ms)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/approvals", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())

	h.ListApprovals(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------- ListTemplates ----------

func TestListTemplates_Success(t *testing.T) {
	templates := []domain.RoleTemplate{
		{
			ID:          uuid.New(),
			Name:        "admin_template",
			DisplayName: "Admin Template",
			Permissions: []domain.PermissionPair{{Resource: "all", Action: "*"}},
			CreatedAt:   time.Now(),
		},
	}
	ms := &mockStore{
		getRoleTemplatesFn: func(_ context.Context) ([]domain.RoleTemplate, error) {
			return templates, nil
		},
	}
	h := newTestHandlers(ms)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/templates", nil)

	h.ListTemplates(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	data := parseSuccessResp(t, rec.Body.Bytes())
	var got []domain.RoleTemplate
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Len(t, got, 1)
	assert.Equal(t, "admin_template", got[0].Name)
}

func TestListTemplates_EmptyList(t *testing.T) {
	ms := &mockStore{
		getRoleTemplatesFn: func(_ context.Context) ([]domain.RoleTemplate, error) {
			return nil, nil
		},
	}
	h := newTestHandlers(ms)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/templates", nil)

	h.ListTemplates(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	data := parseSuccessResp(t, rec.Body.Bytes())
	var got []domain.RoleTemplate
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Len(t, got, 0)
}

func TestListTemplates_StoreError(t *testing.T) {
	ms := &mockStore{
		getRoleTemplatesFn: func(_ context.Context) ([]domain.RoleTemplate, error) {
			return nil, fmt.Errorf("db error")
		},
	}
	h := newTestHandlers(ms)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/templates", nil)

	h.ListTemplates(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------- Router integration test ----------

func TestNewRouter_ReturnsNonNil(t *testing.T) {
	ms := &mockStore{}
	r := NewRouter(ms)
	assert.NotNil(t, r)
}

// ---------- tenantIDFromRequest / userIDFromRequest helpers ----------

func TestTenantIDFromRequest_InvalidUUID(t *testing.T) {
	h := newTestHandlers(&mockStore{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/roles", nil)
	req.Header.Set("X-Tenant-Id", "not-a-uuid")

	h.ListRoles(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserIDFromRequest_InvalidUUID(t *testing.T) {
	tenantID := uuid.New()
	h := newTestHandlers(&mockStore{})
	body := `{"resource_type":"role_change","action_type":"promote"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/approvals", bytes.NewBufferString(body))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.Header.Set("X-User-Id", "not-a-uuid")

	h.CreateApproval(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
