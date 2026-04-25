package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/services/go/identity-service/internal/domain"
)

// ---------------------------------------------------------------------------
// Mock repository
// ---------------------------------------------------------------------------

type mockRepo struct {
	createUserFn       func(ctx context.Context, tenantID uuid.UUID, u *domain.User) error
	getUserByEmailFn   func(ctx context.Context, tenantID uuid.UUID, email string) (*domain.User, error)
	getUserByIDFn      func(ctx context.Context, tenantID, userID uuid.UUID) (*domain.User, error)
	listUsersFn        func(ctx context.Context, tenantID uuid.UUID) ([]domain.User, error)
	createSessionFn    func(ctx context.Context, tenantID uuid.UUID, sess *domain.UserSession) error
	createMFAFactorFn  func(ctx context.Context, tenantID uuid.UUID, f *domain.MFAFactor) error
	getMFAFactorsFn    func(ctx context.Context, tenantID, userID uuid.UUID) ([]domain.MFAFactor, error)
	createStepUpFn     func(ctx context.Context, tenantID uuid.UUID, evt *domain.StepUpEvent) error
	hasValidStepUpFn   func(ctx context.Context, tenantID, userID, sessionID uuid.UUID, actionClass string) (bool, error)
}

func (m *mockRepo) CreateUser(ctx context.Context, tenantID uuid.UUID, u *domain.User) error {
	if m.createUserFn != nil {
		return m.createUserFn(ctx, tenantID, u)
	}
	return nil
}

func (m *mockRepo) GetUserByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*domain.User, error) {
	if m.getUserByEmailFn != nil {
		return m.getUserByEmailFn(ctx, tenantID, email)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) GetUserByID(ctx context.Context, tenantID, userID uuid.UUID) (*domain.User, error) {
	if m.getUserByIDFn != nil {
		return m.getUserByIDFn(ctx, tenantID, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) ListUsers(ctx context.Context, tenantID uuid.UUID) ([]domain.User, error) {
	if m.listUsersFn != nil {
		return m.listUsersFn(ctx, tenantID)
	}
	return nil, nil
}

func (m *mockRepo) CreateSession(ctx context.Context, tenantID uuid.UUID, sess *domain.UserSession) error {
	if m.createSessionFn != nil {
		return m.createSessionFn(ctx, tenantID, sess)
	}
	return nil
}

func (m *mockRepo) CreateMFAFactor(ctx context.Context, tenantID uuid.UUID, f *domain.MFAFactor) error {
	if m.createMFAFactorFn != nil {
		return m.createMFAFactorFn(ctx, tenantID, f)
	}
	return nil
}

func (m *mockRepo) GetMFAFactors(ctx context.Context, tenantID, userID uuid.UUID) ([]domain.MFAFactor, error) {
	if m.getMFAFactorsFn != nil {
		return m.getMFAFactorsFn(ctx, tenantID, userID)
	}
	return nil, nil
}

func (m *mockRepo) CreateStepUpEvent(ctx context.Context, tenantID uuid.UUID, evt *domain.StepUpEvent) error {
	if m.createStepUpFn != nil {
		return m.createStepUpFn(ctx, tenantID, evt)
	}
	return nil
}

func (m *mockRepo) HasValidStepUp(ctx context.Context, tenantID, userID, sessionID uuid.UUID, actionClass string) (bool, error) {
	if m.hasValidStepUpFn != nil {
		return m.hasValidStepUpFn(ctx, tenantID, userID, sessionID, actionClass)
	}
	return false, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

var (
	testTenantID  = uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	testUserID    = uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	testSessionID = uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")
)

func newHandlers(repo *mockRepo, keycloakURL string) *Handlers {
	return NewHandlers(repo, keycloakURL, "test-client", "test-secret")
}

// jsonBody is a shorthand for encoding a value as a request body.
func jsonBody(t *testing.T, v interface{}) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return bytes.NewBuffer(b)
}

// dataEnvelope matches the httputil.JSON wrapper: {"data": ...}
type dataEnvelope struct {
	Data json.RawMessage `json:"data"`
}

func parseData(t *testing.T, body []byte) json.RawMessage {
	t.Helper()
	var env dataEnvelope
	require.NoError(t, json.Unmarshal(body, &env), "response should be a data envelope")
	return env.Data
}

func parseDataMap(t *testing.T, body []byte) map[string]interface{} {
	t.Helper()
	raw := parseData(t, body)
	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(raw, &m))
	return m
}

func addTenantHeader(r *http.Request) {
	r.Header.Set("X-Tenant-Id", testTenantID.String())
}

func addAllHeaders(r *http.Request) {
	r.Header.Set("X-Tenant-Id", testTenantID.String())
	r.Header.Set("X-User-Id", testUserID.String())
	r.Header.Set("X-Session-Id", testSessionID.String())
}

// withChiURLParam sets a URL path parameter so that r.PathValue(key) returns
// the given value. Go 1.22+ supports SetPathValue on *http.Request directly.
func withChiURLParam(r *http.Request, key, value string) *http.Request {
	r.SetPathValue(key, value)
	return r
}

// ---------------------------------------------------------------------------
// Health
// ---------------------------------------------------------------------------

func TestHealth(t *testing.T) {
	h := newHandlers(&mockRepo{}, "")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	h.Health(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	data := parseDataMap(t, rr.Body.Bytes())
	assert.Equal(t, "ok", data["status"])
	assert.Equal(t, "identity-service", data["service"])
}

// ---------------------------------------------------------------------------
// Login
// ---------------------------------------------------------------------------

func fakeKeycloak(t *testing.T, statusCode int, respBody interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(respBody)
	}))
}

func TestLogin_Success(t *testing.T) {
	kcResp := map[string]interface{}{
		"access_token":  "test-access-token",
		"refresh_token": "test-refresh-token",
		"token_type":    "Bearer",
		"expires_in":    float64(300),
	}
	kc := fakeKeycloak(t, http.StatusOK, kcResp)
	defer kc.Close()

	h := newHandlers(&mockRepo{}, kc.URL)
	rr := httptest.NewRecorder()
	body := jsonBody(t, domain.LoginRequest{Username: "user@test.com", Password: "pass"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", body)
	req.Header.Set("Content-Type", "application/json")

	h.Login(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	data := parseDataMap(t, rr.Body.Bytes())
	assert.Equal(t, "test-access-token", data["access_token"])
	assert.Equal(t, "test-refresh-token", data["refresh_token"])
	assert.Equal(t, "Bearer", data["token_type"])
	assert.Equal(t, float64(300), data["expires_in"])
}

func TestLogin_InvalidBody(t *testing.T) {
	h := newHandlers(&mockRepo{}, "")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")

	h.Login(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	data := parseDataMap(t, rr.Body.Bytes())
	assert.Contains(t, data["error"], "invalid request body")
}

func TestLogin_KeycloakReject(t *testing.T) {
	kc := fakeKeycloak(t, http.StatusUnauthorized, map[string]string{"error": "invalid_grant"})
	defer kc.Close()

	h := newHandlers(&mockRepo{}, kc.URL)
	rr := httptest.NewRecorder()
	body := jsonBody(t, domain.LoginRequest{Username: "bad@test.com", Password: "wrong"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", body)
	req.Header.Set("Content-Type", "application/json")

	h.Login(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	data := parseDataMap(t, rr.Body.Bytes())
	assert.Equal(t, "invalid credentials", data["error"])
}

func TestLogin_KeycloakUnreachable(t *testing.T) {
	// Point at a URL that is guaranteed to fail fast.
	h := newHandlers(&mockRepo{}, "http://127.0.0.1:1")
	rr := httptest.NewRecorder()
	body := jsonBody(t, domain.LoginRequest{Username: "u", Password: "p"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", body)
	req.Header.Set("Content-Type", "application/json")

	h.Login(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

// ---------------------------------------------------------------------------
// Refresh
// ---------------------------------------------------------------------------

func TestRefresh_Success(t *testing.T) {
	kcResp := map[string]interface{}{
		"access_token":  "new-access",
		"refresh_token": "new-refresh",
		"token_type":    "Bearer",
		"expires_in":    float64(300),
	}
	kc := fakeKeycloak(t, http.StatusOK, kcResp)
	defer kc.Close()

	h := newHandlers(&mockRepo{}, kc.URL)
	rr := httptest.NewRecorder()
	body := jsonBody(t, map[string]string{"refresh_token": "old-refresh"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/refresh", body)
	req.Header.Set("Content-Type", "application/json")

	h.Refresh(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	data := parseDataMap(t, rr.Body.Bytes())
	assert.Equal(t, "new-access", data["access_token"])
}

func TestRefresh_InvalidBody(t *testing.T) {
	h := newHandlers(&mockRepo{}, "")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/refresh", bytes.NewBufferString("{bad"))
	req.Header.Set("Content-Type", "application/json")

	h.Refresh(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestRefresh_KeycloakFail(t *testing.T) {
	kc := fakeKeycloak(t, http.StatusUnauthorized, map[string]string{"error": "invalid_grant"})
	defer kc.Close()

	h := newHandlers(&mockRepo{}, kc.URL)
	rr := httptest.NewRecorder()
	body := jsonBody(t, map[string]string{"refresh_token": "expired"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/refresh", body)
	req.Header.Set("Content-Type", "application/json")

	h.Refresh(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	data := parseDataMap(t, rr.Body.Bytes())
	assert.Equal(t, "refresh failed", data["error"])
}

// ---------------------------------------------------------------------------
// Logout
// ---------------------------------------------------------------------------

func TestLogout_Success(t *testing.T) {
	kc := fakeKeycloak(t, http.StatusNoContent, nil)
	defer kc.Close()

	h := newHandlers(&mockRepo{}, kc.URL)
	rr := httptest.NewRecorder()
	body := jsonBody(t, map[string]string{"refresh_token": "tok"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/logout", body)
	req.Header.Set("Content-Type", "application/json")

	h.Logout(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	data := parseDataMap(t, rr.Body.Bytes())
	assert.Equal(t, "logged_out", data["status"])
}

func TestLogout_InvalidBody(t *testing.T) {
	h := newHandlers(&mockRepo{}, "")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/logout", bytes.NewBufferString("bad"))
	req.Header.Set("Content-Type", "application/json")

	h.Logout(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLogout_KeycloakError(t *testing.T) {
	// Even if keycloak fails, logout returns 200 (handler logs but doesn't fail).
	h := newHandlers(&mockRepo{}, "http://127.0.0.1:1")
	rr := httptest.NewRecorder()
	body := jsonBody(t, map[string]string{"refresh_token": "tok"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/logout", body)
	req.Header.Set("Content-Type", "application/json")

	h.Logout(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

// ---------------------------------------------------------------------------
// ListUsers
// ---------------------------------------------------------------------------

func TestListUsers_Success(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	repo := &mockRepo{
		listUsersFn: func(_ context.Context, tid uuid.UUID) ([]domain.User, error) {
			assert.Equal(t, testTenantID, tid)
			return []domain.User{
				{ID: testUserID, TenantID: tid, Email: "a@b.com", FirstName: "A", LastName: "B", Status: "active", CreatedAt: now, UpdatedAt: now},
			}, nil
		},
	}
	h := newHandlers(repo, "")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
	addTenantHeader(req)

	h.ListUsers(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	raw := parseData(t, rr.Body.Bytes())
	var users []domain.User
	require.NoError(t, json.Unmarshal(raw, &users))
	require.Len(t, users, 1)
	assert.Equal(t, "a@b.com", users[0].Email)
}

func TestListUsers_EmptyList(t *testing.T) {
	repo := &mockRepo{
		listUsersFn: func(_ context.Context, _ uuid.UUID) ([]domain.User, error) {
			return nil, nil
		},
	}
	h := newHandlers(repo, "")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
	addTenantHeader(req)

	h.ListUsers(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	raw := parseData(t, rr.Body.Bytes())
	var users []domain.User
	require.NoError(t, json.Unmarshal(raw, &users))
	assert.Empty(t, users)
}

func TestListUsers_MissingTenantHeader(t *testing.T) {
	h := newHandlers(&mockRepo{}, "")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)

	h.ListUsers(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	data := parseDataMap(t, rr.Body.Bytes())
	assert.Contains(t, data["error"], "X-Tenant-Id")
}

func TestListUsers_StoreError(t *testing.T) {
	repo := &mockRepo{
		listUsersFn: func(_ context.Context, _ uuid.UUID) ([]domain.User, error) {
			return nil, errors.New("db down")
		},
	}
	h := newHandlers(repo, "")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
	addTenantHeader(req)

	h.ListUsers(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	data := parseDataMap(t, rr.Body.Bytes())
	assert.Equal(t, "internal error", data["error"])
}

// ---------------------------------------------------------------------------
// GetUser
// ---------------------------------------------------------------------------

func TestGetUser_Success(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	repo := &mockRepo{
		getUserByIDFn: func(_ context.Context, tid, uid uuid.UUID) (*domain.User, error) {
			assert.Equal(t, testTenantID, tid)
			assert.Equal(t, testUserID, uid)
			return &domain.User{
				ID: uid, TenantID: tid, Email: "u@u.com",
				FirstName: "U", LastName: "U", Status: "active",
				CreatedAt: now, UpdatedAt: now,
			}, nil
		},
	}
	h := newHandlers(repo, "")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/users/"+testUserID.String(), nil)
	addTenantHeader(req)
	req = withChiURLParam(req, "userID", testUserID.String())

	h.GetUser(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	data := parseDataMap(t, rr.Body.Bytes())
	assert.Equal(t, testUserID.String(), data["id"])
}

func TestGetUser_MissingTenantHeader(t *testing.T) {
	h := newHandlers(&mockRepo{}, "")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/users/"+testUserID.String(), nil)
	req = withChiURLParam(req, "userID", testUserID.String())

	h.GetUser(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetUser_InvalidUserID(t *testing.T) {
	h := newHandlers(&mockRepo{}, "")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/users/not-a-uuid", nil)
	addTenantHeader(req)
	req = withChiURLParam(req, "userID", "not-a-uuid")

	h.GetUser(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	data := parseDataMap(t, rr.Body.Bytes())
	assert.Equal(t, "invalid user_id", data["error"])
}

func TestGetUser_NotFound(t *testing.T) {
	repo := &mockRepo{
		getUserByIDFn: func(_ context.Context, _, _ uuid.UUID) (*domain.User, error) {
			return nil, errors.New("get user: no rows")
		},
	}
	h := newHandlers(repo, "")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/users/"+testUserID.String(), nil)
	addTenantHeader(req)
	req = withChiURLParam(req, "userID", testUserID.String())

	h.GetUser(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	data := parseDataMap(t, rr.Body.Bytes())
	assert.Equal(t, "user not found", data["error"])
}

// ---------------------------------------------------------------------------
// StepUpCheck
// ---------------------------------------------------------------------------

func TestStepUpCheck_Valid(t *testing.T) {
	repo := &mockRepo{
		hasValidStepUpFn: func(_ context.Context, tid, uid, sid uuid.UUID, ac string) (bool, error) {
			assert.Equal(t, testTenantID, tid)
			assert.Equal(t, testUserID, uid)
			assert.Equal(t, testSessionID, sid)
			assert.Equal(t, "invoice_approve", ac)
			return true, nil
		},
	}
	h := newHandlers(repo, "")
	rr := httptest.NewRecorder()
	body := jsonBody(t, domain.StepUpCheckRequest{ActionClass: "invoice_approve"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/step-up/check", body)
	addAllHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	h.StepUpCheck(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	data := parseDataMap(t, rr.Body.Bytes())
	assert.Equal(t, "step_up_valid", data["status"])
}

func TestStepUpCheck_Required(t *testing.T) {
	repo := &mockRepo{
		hasValidStepUpFn: func(_ context.Context, _, _, _ uuid.UUID, _ string) (bool, error) {
			return false, nil
		},
	}
	h := newHandlers(repo, "")
	rr := httptest.NewRecorder()
	body := jsonBody(t, domain.StepUpCheckRequest{ActionClass: "invoice_approve"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/step-up/check", body)
	addAllHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	h.StepUpCheck(rr, req)

	assert.Equal(t, http.StatusForbidden, rr.Code)

	// StepUpCheck writes 403 directly (not via httputil.JSON), so no data envelope.
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, "step_up_required", resp["error"])
	assert.Equal(t, "invoice_approve", resp["action"])
	assert.Equal(t, "/v1/auth/step-up", resp["step_up_url"])
}

func TestStepUpCheck_MissingTenantHeader(t *testing.T) {
	h := newHandlers(&mockRepo{}, "")
	rr := httptest.NewRecorder()
	body := jsonBody(t, domain.StepUpCheckRequest{ActionClass: "a"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/step-up/check", body)
	// no headers
	h.StepUpCheck(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestStepUpCheck_MissingUserHeader(t *testing.T) {
	h := newHandlers(&mockRepo{}, "")
	rr := httptest.NewRecorder()
	body := jsonBody(t, domain.StepUpCheckRequest{ActionClass: "a"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/step-up/check", body)
	req.Header.Set("X-Tenant-Id", testTenantID.String())
	// missing X-User-Id, X-Session-Id
	h.StepUpCheck(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestStepUpCheck_MissingSessionHeader(t *testing.T) {
	h := newHandlers(&mockRepo{}, "")
	rr := httptest.NewRecorder()
	body := jsonBody(t, domain.StepUpCheckRequest{ActionClass: "a"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/step-up/check", body)
	req.Header.Set("X-Tenant-Id", testTenantID.String())
	req.Header.Set("X-User-Id", testUserID.String())
	// missing X-Session-Id
	h.StepUpCheck(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestStepUpCheck_InvalidBody(t *testing.T) {
	h := newHandlers(&mockRepo{}, "")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/step-up/check", bytes.NewBufferString("bad"))
	addAllHeaders(req)
	h.StepUpCheck(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestStepUpCheck_StoreError(t *testing.T) {
	repo := &mockRepo{
		hasValidStepUpFn: func(_ context.Context, _, _, _ uuid.UUID, _ string) (bool, error) {
			return false, errors.New("db err")
		},
	}
	h := newHandlers(repo, "")
	rr := httptest.NewRecorder()
	body := jsonBody(t, domain.StepUpCheckRequest{ActionClass: "a"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/step-up/check", body)
	addAllHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	h.StepUpCheck(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// ---------------------------------------------------------------------------
// StepUpVerify
// ---------------------------------------------------------------------------

func TestStepUpVerify_ValidCode(t *testing.T) {
	repo := &mockRepo{
		createStepUpFn: func(_ context.Context, tid uuid.UUID, evt *domain.StepUpEvent) error {
			assert.Equal(t, testTenantID, tid)
			assert.Equal(t, testUserID, evt.UserID)
			assert.Equal(t, testSessionID, evt.SessionID)
			assert.Equal(t, "invoice_approve", evt.ActionClass)
			assert.Equal(t, "totp", evt.MFAMethod)
			// Set a fake ID to prove the handler returns data from the event.
			evt.ID = uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd")
			return nil
		},
	}
	h := newHandlers(repo, "")
	rr := httptest.NewRecorder()
	body := jsonBody(t, domain.StepUpRequest{ActionClass: "invoice_approve", MFACode: "123456"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/step-up/verify", body)
	addAllHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	h.StepUpVerify(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	data := parseDataMap(t, rr.Body.Bytes())
	assert.Equal(t, "step_up_verified", data["status"])
	assert.Equal(t, "invoice_approve", data["action"])
	assert.NotEmpty(t, data["expires_at"])
}

func TestStepUpVerify_InvalidCode_WithFactors(t *testing.T) {
	repo := &mockRepo{
		getMFAFactorsFn: func(_ context.Context, _, _ uuid.UUID) ([]domain.MFAFactor, error) {
			return []domain.MFAFactor{{ID: uuid.New(), FactorType: "totp", Verified: true}}, nil
		},
	}
	h := newHandlers(repo, "")
	rr := httptest.NewRecorder()
	body := jsonBody(t, domain.StepUpRequest{ActionClass: "a", MFACode: "000000"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/step-up/verify", body)
	addAllHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	h.StepUpVerify(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	data := parseDataMap(t, rr.Body.Bytes())
	assert.Equal(t, "invalid MFA code", data["error"])
}

func TestStepUpVerify_InvalidCode_NoFactors(t *testing.T) {
	repo := &mockRepo{
		getMFAFactorsFn: func(_ context.Context, _, _ uuid.UUID) ([]domain.MFAFactor, error) {
			return nil, nil
		},
	}
	h := newHandlers(repo, "")
	rr := httptest.NewRecorder()
	body := jsonBody(t, domain.StepUpRequest{ActionClass: "a", MFACode: "999999"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/step-up/verify", body)
	addAllHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	h.StepUpVerify(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	data := parseDataMap(t, rr.Body.Bytes())
	assert.Equal(t, "no MFA factor enrolled", data["error"])
}

func TestStepUpVerify_InvalidBody(t *testing.T) {
	h := newHandlers(&mockRepo{}, "")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/step-up/verify", bytes.NewBufferString("bad"))
	addAllHeaders(req)

	h.StepUpVerify(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestStepUpVerify_MissingHeaders(t *testing.T) {
	h := newHandlers(&mockRepo{}, "")

	t.Run("no tenant", func(t *testing.T) {
		rr := httptest.NewRecorder()
		body := jsonBody(t, domain.StepUpRequest{ActionClass: "a", MFACode: "123456"})
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/step-up/verify", body)
		h.StepUpVerify(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("no user", func(t *testing.T) {
		rr := httptest.NewRecorder()
		body := jsonBody(t, domain.StepUpRequest{ActionClass: "a", MFACode: "123456"})
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/step-up/verify", body)
		req.Header.Set("X-Tenant-Id", testTenantID.String())
		h.StepUpVerify(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("no session", func(t *testing.T) {
		rr := httptest.NewRecorder()
		body := jsonBody(t, domain.StepUpRequest{ActionClass: "a", MFACode: "123456"})
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/step-up/verify", body)
		req.Header.Set("X-Tenant-Id", testTenantID.String())
		req.Header.Set("X-User-Id", testUserID.String())
		h.StepUpVerify(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestStepUpVerify_CreateEventError(t *testing.T) {
	repo := &mockRepo{
		createStepUpFn: func(_ context.Context, _ uuid.UUID, _ *domain.StepUpEvent) error {
			return errors.New("db down")
		},
	}
	h := newHandlers(repo, "")
	rr := httptest.NewRecorder()
	body := jsonBody(t, domain.StepUpRequest{ActionClass: "a", MFACode: "123456"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/step-up/verify", body)
	addAllHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	h.StepUpVerify(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// ---------------------------------------------------------------------------
// EnrollMFA
// ---------------------------------------------------------------------------

func TestEnrollMFA_Success(t *testing.T) {
	factorID := uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee")
	repo := &mockRepo{
		createMFAFactorFn: func(_ context.Context, tid uuid.UUID, f *domain.MFAFactor) error {
			assert.Equal(t, testTenantID, tid)
			assert.Equal(t, testUserID, f.UserID)
			assert.Equal(t, "totp", f.FactorType)
			assert.True(t, f.Verified)
			f.ID = factorID
			return nil
		},
	}
	h := newHandlers(repo, "")
	rr := httptest.NewRecorder()
	body := jsonBody(t, map[string]string{"factor_type": "totp"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/mfa/enroll", body)
	addAllHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	h.EnrollMFA(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	data := parseDataMap(t, rr.Body.Bytes())
	assert.Equal(t, factorID.String(), data["factor_id"])
	assert.Equal(t, "totp", data["factor_type"])
	assert.Equal(t, true, data["verified"])
}

func TestEnrollMFA_InvalidBody(t *testing.T) {
	h := newHandlers(&mockRepo{}, "")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/mfa/enroll", bytes.NewBufferString("bad"))
	addAllHeaders(req)
	h.EnrollMFA(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestEnrollMFA_MissingTenantHeader(t *testing.T) {
	h := newHandlers(&mockRepo{}, "")
	rr := httptest.NewRecorder()
	body := jsonBody(t, map[string]string{"factor_type": "totp"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/mfa/enroll", body)
	// no headers
	h.EnrollMFA(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestEnrollMFA_MissingUserHeader(t *testing.T) {
	h := newHandlers(&mockRepo{}, "")
	rr := httptest.NewRecorder()
	body := jsonBody(t, map[string]string{"factor_type": "totp"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/mfa/enroll", body)
	req.Header.Set("X-Tenant-Id", testTenantID.String())
	// no X-User-Id
	h.EnrollMFA(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestEnrollMFA_StoreError(t *testing.T) {
	repo := &mockRepo{
		createMFAFactorFn: func(_ context.Context, _ uuid.UUID, _ *domain.MFAFactor) error {
			return errors.New("db boom")
		},
	}
	h := newHandlers(repo, "")
	rr := httptest.NewRecorder()
	body := jsonBody(t, map[string]string{"factor_type": "totp"})
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/mfa/enroll", body)
	addAllHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	h.EnrollMFA(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// ---------------------------------------------------------------------------
// Router integration — ensures routes are wired correctly
// ---------------------------------------------------------------------------

func TestRouter_Integration(t *testing.T) {
	kc := fakeKeycloak(t, http.StatusOK, map[string]interface{}{
		"access_token":  "tok",
		"refresh_token": "rtok",
		"token_type":    "Bearer",
		"expires_in":    float64(300),
	})
	defer kc.Close()

	repo := &mockRepo{
		listUsersFn: func(_ context.Context, _ uuid.UUID) ([]domain.User, error) {
			return []domain.User{}, nil
		},
		getUserByIDFn: func(_ context.Context, _, uid uuid.UUID) (*domain.User, error) {
			return &domain.User{ID: uid, Email: "x@x.com", Status: "active"}, nil
		},
		hasValidStepUpFn: func(_ context.Context, _, _, _ uuid.UUID, _ string) (bool, error) {
			return true, nil
		},
		createStepUpFn: func(_ context.Context, _ uuid.UUID, evt *domain.StepUpEvent) error {
			return nil
		},
		createMFAFactorFn: func(_ context.Context, _ uuid.UUID, f *domain.MFAFactor) error {
			f.ID = uuid.New()
			return nil
		},
	}

	router := NewRouter(repo, kc.URL, "cid", "csec")
	ts := httptest.NewServer(router)
	defer ts.Close()

	cases := []struct {
		name   string
		method string
		path   string
		body   interface{}
		code   int
	}{
		{"health", http.MethodGet, "/health", nil, http.StatusOK},
		{"ping", http.MethodGet, "/ping", nil, http.StatusOK},
		{"login", http.MethodPost, "/v1/auth/login", domain.LoginRequest{Username: "u", Password: "p"}, http.StatusOK},
		{"refresh", http.MethodPost, "/v1/auth/refresh", map[string]string{"refresh_token": "r"}, http.StatusOK},
		{"logout", http.MethodPost, "/v1/auth/logout", map[string]string{"refresh_token": "r"}, http.StatusOK},
	}

	client := ts.Client()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			if tc.body != nil {
				b, _ := json.Marshal(tc.body)
				req, _ = http.NewRequest(tc.method, ts.URL+tc.path, bytes.NewBuffer(b))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, _ = http.NewRequest(tc.method, ts.URL+tc.path, nil)
			}
			addTenantHeader(req)
			req.Header.Set("X-User-Id", testUserID.String())
			req.Header.Set("X-Session-Id", testSessionID.String())

			resp, err := client.Do(req)
			require.NoError(t, err)
			resp.Body.Close()
			assert.Equal(t, tc.code, resp.StatusCode, tc.name)
		})
	}

	// Test user routes specifically
	t.Run("list_users", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/v1/users", nil)
		addTenantHeader(req)
		resp, err := client.Do(req)
		require.NoError(t, err)
		resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("get_user", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/users/%s", ts.URL, testUserID), nil)
		addTenantHeader(req)
		resp, err := client.Do(req)
		require.NoError(t, err)
		resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("step_up_check", func(t *testing.T) {
		b, _ := json.Marshal(domain.StepUpCheckRequest{ActionClass: "x"})
		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/auth/step-up/check", bytes.NewBuffer(b))
		addAllHeaders(req)
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		require.NoError(t, err)
		resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("step_up_verify", func(t *testing.T) {
		b, _ := json.Marshal(domain.StepUpRequest{ActionClass: "x", MFACode: "123456"})
		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/auth/step-up/verify", bytes.NewBuffer(b))
		addAllHeaders(req)
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		require.NoError(t, err)
		resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("enroll_mfa", func(t *testing.T) {
		b, _ := json.Marshal(map[string]string{"factor_type": "totp"})
		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/v1/auth/mfa/enroll", bytes.NewBuffer(b))
		addAllHeaders(req)
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		require.NoError(t, err)
		resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})
}

// ---------------------------------------------------------------------------
// Helper function tests (tenantIDFromRequest, etc.)
// ---------------------------------------------------------------------------

func TestTenantIDFromRequest_InvalidUUID(t *testing.T) {
	h := newHandlers(&mockRepo{}, "")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
	req.Header.Set("X-Tenant-Id", "not-a-uuid")

	h.ListUsers(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
