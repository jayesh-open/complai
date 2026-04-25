package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	kmsTypes "github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/tenant-service/internal/domain"
)

// ---------------------------------------------------------------------------
// Mock store
// ---------------------------------------------------------------------------

type mockStore struct {
	createTenantFn      func(ctx context.Context, t *domain.Tenant) error
	getTenantFn         func(ctx context.Context, id uuid.UUID) (*domain.Tenant, error)
	listTenantsFn       func(ctx context.Context, id uuid.UUID) ([]domain.Tenant, error)
	updateTenantKMSFn   func(ctx context.Context, id uuid.UUID, arn string) error
	updateTenantStatusFn func(ctx context.Context, id uuid.UUID, status string) error
	createPANFn         func(ctx context.Context, tenantID uuid.UUID, p *domain.TenantPAN) error
	createGSTINFn       func(ctx context.Context, tenantID uuid.UUID, g *domain.TenantGSTIN) error
	createTANFn         func(ctx context.Context, tenantID uuid.UUID, t *domain.TenantTAN) error
	getHierarchyFn      func(ctx context.Context, id uuid.UUID) (*domain.TenantHierarchy, error)
}

func (m *mockStore) CreateTenant(ctx context.Context, t *domain.Tenant) error {
	if m.createTenantFn != nil {
		return m.createTenantFn(ctx, t)
	}
	t.ID = uuid.New()
	t.TenantID = t.ID
	t.Status = "active"
	t.Settings = "{}"
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return nil
}

func (m *mockStore) GetTenant(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	if m.getTenantFn != nil {
		return m.getTenantFn(ctx, id)
	}
	return nil, errors.New("not found")
}

func (m *mockStore) ListTenants(ctx context.Context, id uuid.UUID) ([]domain.Tenant, error) {
	if m.listTenantsFn != nil {
		return m.listTenantsFn(ctx, id)
	}
	return nil, nil
}

func (m *mockStore) UpdateTenantKMSKey(ctx context.Context, id uuid.UUID, arn string) error {
	if m.updateTenantKMSFn != nil {
		return m.updateTenantKMSFn(ctx, id, arn)
	}
	return nil
}

func (m *mockStore) UpdateTenantStatus(ctx context.Context, id uuid.UUID, status string) error {
	if m.updateTenantStatusFn != nil {
		return m.updateTenantStatusFn(ctx, id, status)
	}
	return nil
}

func (m *mockStore) CreatePAN(ctx context.Context, tenantID uuid.UUID, p *domain.TenantPAN) error {
	if m.createPANFn != nil {
		return m.createPANFn(ctx, tenantID, p)
	}
	p.ID = uuid.New()
	p.TenantID = tenantID
	p.Status = "active"
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

func (m *mockStore) CreateGSTIN(ctx context.Context, tenantID uuid.UUID, g *domain.TenantGSTIN) error {
	if m.createGSTINFn != nil {
		return m.createGSTINFn(ctx, tenantID, g)
	}
	g.ID = uuid.New()
	g.TenantID = tenantID
	g.Status = "active"
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()
	return nil
}

func (m *mockStore) CreateTAN(ctx context.Context, tenantID uuid.UUID, t *domain.TenantTAN) error {
	if m.createTANFn != nil {
		return m.createTANFn(ctx, tenantID, t)
	}
	t.ID = uuid.New()
	t.TenantID = tenantID
	t.Status = "active"
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return nil
}

func (m *mockStore) GetHierarchy(ctx context.Context, id uuid.UUID) (*domain.TenantHierarchy, error) {
	if m.getHierarchyFn != nil {
		return m.getHierarchyFn(ctx, id)
	}
	return &domain.TenantHierarchy{}, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// parseDataResponse parses {"data": ...} into target.
func parseDataResponse(t *testing.T, body []byte, target interface{}) {
	t.Helper()
	var wrapper httputil.SuccessResponse
	wrapper.Data = target
	require.NoError(t, json.Unmarshal(body, &wrapper))
}

// newRequest builds a request with optional JSON body and chi-style path values.
func newRequest(t *testing.T, method, path string, body interface{}, pathValues map[string]string) *http.Request {
	t.Helper()
	var req *http.Request
	if body != nil {
		b, err := json.Marshal(body)
		require.NoError(t, err)
		req = httptest.NewRequest(method, path, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	for k, v := range pathValues {
		req.SetPathValue(k, v)
	}
	return req
}

// ---------------------------------------------------------------------------
// Tests: Health
// ---------------------------------------------------------------------------

func TestHealth(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	h.Health(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "ok", data["status"])
	assert.Equal(t, "tenant-service", data["service"])
}

// ---------------------------------------------------------------------------
// Tests: CreateTenant
// ---------------------------------------------------------------------------

func TestCreateTenant_Success(t *testing.T) {
	ms := &mockStore{}
	h := NewHandlers(ms, nil) // nil KMS => skip KMS path

	body := domain.CreateTenantRequest{Name: "Acme Corp", Slug: "acme", Tier: "pooled"}
	req := newRequest(t, http.MethodPost, "/v1/tenants", body, nil)
	rec := httptest.NewRecorder()

	h.CreateTenant(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var tenant domain.Tenant
	parseDataResponse(t, rec.Body.Bytes(), &tenant)
	assert.Equal(t, "Acme Corp", tenant.Name)
	assert.Equal(t, "acme", tenant.Slug)
	assert.Equal(t, "pooled", tenant.Tier)
	assert.Equal(t, "active", tenant.Status)
	assert.NotEqual(t, uuid.Nil, tenant.ID)
}

func TestCreateTenant_InvalidBody(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil)
	req := httptest.NewRequest(http.MethodPost, "/v1/tenants", bytes.NewReader([]byte("not json")))
	rec := httptest.NewRecorder()

	h.CreateTenant(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid request", data["error"])
}

func TestCreateTenant_StoreError(t *testing.T) {
	ms := &mockStore{
		createTenantFn: func(_ context.Context, _ *domain.Tenant) error {
			return errors.New("db down")
		},
	}
	h := NewHandlers(ms, nil)

	body := domain.CreateTenantRequest{Name: "Acme", Slug: "acme", Tier: "pooled"}
	req := newRequest(t, http.MethodPost, "/v1/tenants", body, nil)
	rec := httptest.NewRecorder()

	h.CreateTenant(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "create failed", data["error"])
}

// ---------------------------------------------------------------------------
// Tests: GetTenant
// ---------------------------------------------------------------------------

func TestGetTenant_Success(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		getTenantFn: func(_ context.Context, id uuid.UUID) (*domain.Tenant, error) {
			assert.Equal(t, tenantID, id)
			return &domain.Tenant{
				ID:       tenantID,
				TenantID: tenantID,
				Name:     "Acme",
				Slug:     "acme",
				Tier:     "pooled",
				Status:   "active",
				Settings: "{}",
			}, nil
		},
	}
	h := NewHandlers(ms, nil)

	req := newRequest(t, http.MethodGet, "/v1/tenants/"+tenantID.String(), nil, map[string]string{
		"tenantID": tenantID.String(),
	})
	rec := httptest.NewRecorder()

	h.GetTenant(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var tenant domain.Tenant
	parseDataResponse(t, rec.Body.Bytes(), &tenant)
	assert.Equal(t, tenantID, tenant.ID)
	assert.Equal(t, "Acme", tenant.Name)
}

func TestGetTenant_InvalidID(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil)

	req := newRequest(t, http.MethodGet, "/v1/tenants/not-a-uuid", nil, map[string]string{
		"tenantID": "not-a-uuid",
	})
	rec := httptest.NewRecorder()

	h.GetTenant(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid tenant_id", data["error"])
}

func TestGetTenant_NotFound(t *testing.T) {
	ms := &mockStore{
		getTenantFn: func(_ context.Context, _ uuid.UUID) (*domain.Tenant, error) {
			return nil, errors.New("not found")
		},
	}
	h := NewHandlers(ms, nil)

	tenantID := uuid.New()
	req := newRequest(t, http.MethodGet, "/v1/tenants/"+tenantID.String(), nil, map[string]string{
		"tenantID": tenantID.String(),
	})
	rec := httptest.NewRecorder()

	h.GetTenant(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: ListTenants
// ---------------------------------------------------------------------------

func TestListTenants_Success(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listTenantsFn: func(_ context.Context, id uuid.UUID) ([]domain.Tenant, error) {
			assert.Equal(t, tenantID, id)
			return []domain.Tenant{
				{ID: tenantID, TenantID: tenantID, Name: "Acme", Slug: "acme", Tier: "pooled", Status: "active", Settings: "{}"},
			}, nil
		},
	}
	h := NewHandlers(ms, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/tenants", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListTenants(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var tenants []domain.Tenant
	parseDataResponse(t, rec.Body.Bytes(), &tenants)
	require.Len(t, tenants, 1)
	assert.Equal(t, "Acme", tenants[0].Name)
}

func TestListTenants_MissingHeader(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/tenants", nil)
	// No X-Tenant-Id header
	rec := httptest.NewRecorder()

	h.ListTenants(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Contains(t, data["error"], "missing X-Tenant-Id")
}

func TestListTenants_InvalidHeader(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/tenants", nil)
	req.Header.Set("X-Tenant-Id", "bad-uuid")
	rec := httptest.NewRecorder()

	h.ListTenants(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListTenants_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listTenantsFn: func(_ context.Context, _ uuid.UUID) ([]domain.Tenant, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewHandlers(ms, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/tenants", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListTenants(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "internal error", data["error"])
}

func TestListTenants_NilResultReturnsEmptyArray(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listTenantsFn: func(_ context.Context, _ uuid.UUID) ([]domain.Tenant, error) {
			return nil, nil
		},
	}
	h := NewHandlers(ms, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/tenants", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListTenants(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var tenants []domain.Tenant
	parseDataResponse(t, rec.Body.Bytes(), &tenants)
	assert.NotNil(t, tenants)
	assert.Len(t, tenants, 0)
}

// ---------------------------------------------------------------------------
// Tests: SuspendTenant
// ---------------------------------------------------------------------------

func TestSuspendTenant_Success(t *testing.T) {
	tenantID := uuid.New()
	var capturedStatus string
	ms := &mockStore{
		updateTenantStatusFn: func(_ context.Context, id uuid.UUID, status string) error {
			assert.Equal(t, tenantID, id)
			capturedStatus = status
			return nil
		},
	}
	h := NewHandlers(ms, nil)

	req := newRequest(t, http.MethodPost, "/v1/tenants/"+tenantID.String()+"/suspend", nil, map[string]string{
		"tenantID": tenantID.String(),
	})
	rec := httptest.NewRecorder()

	h.SuspendTenant(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "suspended", capturedStatus)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "suspended", data["status"])
}

func TestSuspendTenant_InvalidID(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil)

	req := newRequest(t, http.MethodPost, "/v1/tenants/bad/suspend", nil, map[string]string{
		"tenantID": "bad",
	})
	rec := httptest.NewRecorder()

	h.SuspendTenant(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSuspendTenant_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		updateTenantStatusFn: func(_ context.Context, _ uuid.UUID, _ string) error {
			return errors.New("db error")
		},
	}
	h := NewHandlers(ms, nil)

	req := newRequest(t, http.MethodPost, "/v1/tenants/"+tenantID.String()+"/suspend", nil, map[string]string{
		"tenantID": tenantID.String(),
	})
	rec := httptest.NewRecorder()

	h.SuspendTenant(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: ReactivateTenant
// ---------------------------------------------------------------------------

func TestReactivateTenant_Success(t *testing.T) {
	tenantID := uuid.New()
	var capturedStatus string
	ms := &mockStore{
		updateTenantStatusFn: func(_ context.Context, id uuid.UUID, status string) error {
			assert.Equal(t, tenantID, id)
			capturedStatus = status
			return nil
		},
	}
	h := NewHandlers(ms, nil)

	req := newRequest(t, http.MethodPost, "/v1/tenants/"+tenantID.String()+"/reactivate", nil, map[string]string{
		"tenantID": tenantID.String(),
	})
	rec := httptest.NewRecorder()

	h.ReactivateTenant(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "active", capturedStatus)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "active", data["status"])
}

func TestReactivateTenant_InvalidID(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil)

	req := newRequest(t, http.MethodPost, "/v1/tenants/bad/reactivate", nil, map[string]string{
		"tenantID": "bad",
	})
	rec := httptest.NewRecorder()

	h.ReactivateTenant(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestReactivateTenant_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		updateTenantStatusFn: func(_ context.Context, _ uuid.UUID, _ string) error {
			return errors.New("db error")
		},
	}
	h := NewHandlers(ms, nil)

	req := newRequest(t, http.MethodPost, "/v1/tenants/"+tenantID.String()+"/reactivate", nil, map[string]string{
		"tenantID": tenantID.String(),
	})
	rec := httptest.NewRecorder()

	h.ReactivateTenant(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: GetHierarchy
// ---------------------------------------------------------------------------

func TestGetHierarchy_Success(t *testing.T) {
	tenantID := uuid.New()
	panID := uuid.New()
	gstinID := uuid.New()
	tanID := uuid.New()

	ms := &mockStore{
		getHierarchyFn: func(_ context.Context, id uuid.UUID) (*domain.TenantHierarchy, error) {
			assert.Equal(t, tenantID, id)
			return &domain.TenantHierarchy{
				Tenant: domain.Tenant{
					ID: tenantID, TenantID: tenantID, Name: "Acme", Slug: "acme",
					Tier: "pooled", Status: "active", Settings: "{}",
				},
				PANs: []domain.PANWithSub{
					{
						TenantPAN: domain.TenantPAN{
							ID: panID, TenantID: tenantID, PAN: "ABCDE1234F",
							EntityName: "Acme Corp", PANType: "company", Status: "active",
						},
						GSTINs: []domain.TenantGSTIN{
							{ID: gstinID, TenantID: tenantID, PANID: panID, GSTIN: "29ABCDE1234F1Z5", StateCode: "29", RegistrationType: "regular", Status: "active"},
						},
						TANs: []domain.TenantTAN{
							{ID: tanID, TenantID: tenantID, PANID: panID, TAN: "BLRA12345F", DeductorName: "Acme Corp", Status: "active"},
						},
					},
				},
			}, nil
		},
	}
	h := NewHandlers(ms, nil)

	req := newRequest(t, http.MethodGet, "/v1/tenants/"+tenantID.String()+"/hierarchy", nil, map[string]string{
		"tenantID": tenantID.String(),
	})
	rec := httptest.NewRecorder()

	h.GetHierarchy(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var hierarchy domain.TenantHierarchy
	parseDataResponse(t, rec.Body.Bytes(), &hierarchy)
	assert.Equal(t, tenantID, hierarchy.Tenant.ID)
	require.Len(t, hierarchy.PANs, 1)
	assert.Equal(t, "ABCDE1234F", hierarchy.PANs[0].PAN)
	require.Len(t, hierarchy.PANs[0].GSTINs, 1)
	require.Len(t, hierarchy.PANs[0].TANs, 1)
}

func TestGetHierarchy_InvalidID(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil)

	req := newRequest(t, http.MethodGet, "/v1/tenants/bad/hierarchy", nil, map[string]string{
		"tenantID": "bad",
	})
	rec := httptest.NewRecorder()

	h.GetHierarchy(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetHierarchy_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		getHierarchyFn: func(_ context.Context, _ uuid.UUID) (*domain.TenantHierarchy, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewHandlers(ms, nil)

	req := newRequest(t, http.MethodGet, "/v1/tenants/"+tenantID.String()+"/hierarchy", nil, map[string]string{
		"tenantID": tenantID.String(),
	})
	rec := httptest.NewRecorder()

	h.GetHierarchy(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: CreatePAN
// ---------------------------------------------------------------------------

func TestCreatePAN_Success(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{}
	h := NewHandlers(ms, nil)

	body := domain.CreatePANRequest{PAN: "ABCDE1234F", EntityName: "Acme Corp", PANType: "company"}
	req := newRequest(t, http.MethodPost, "/v1/tenants/"+tenantID.String()+"/pans", body, map[string]string{
		"tenantID": tenantID.String(),
	})
	rec := httptest.NewRecorder()

	h.CreatePAN(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var pan domain.TenantPAN
	parseDataResponse(t, rec.Body.Bytes(), &pan)
	assert.Equal(t, "ABCDE1234F", pan.PAN)
	assert.Equal(t, "Acme Corp", pan.EntityName)
	assert.Equal(t, "company", pan.PANType)
	assert.Equal(t, "active", pan.Status)
	assert.Equal(t, tenantID, pan.TenantID)
}

func TestCreatePAN_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, nil)

	req := newRequest(t, http.MethodPost, "/v1/tenants/"+tenantID.String()+"/pans", nil, map[string]string{
		"tenantID": tenantID.String(),
	})
	req.Body = http.NoBody
	// Provide malformed JSON via replacing the body
	req = httptest.NewRequest(http.MethodPost, "/v1/tenants/"+tenantID.String()+"/pans", bytes.NewReader([]byte("bad json")))
	req.SetPathValue("tenantID", tenantID.String())
	rec := httptest.NewRecorder()

	h.CreatePAN(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid request", data["error"])
}

func TestCreatePAN_InvalidTenantID(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil)

	body := domain.CreatePANRequest{PAN: "ABCDE1234F", EntityName: "Acme", PANType: "company"}
	req := newRequest(t, http.MethodPost, "/v1/tenants/bad/pans", body, map[string]string{
		"tenantID": "bad",
	})
	rec := httptest.NewRecorder()

	h.CreatePAN(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreatePAN_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		createPANFn: func(_ context.Context, _ uuid.UUID, _ *domain.TenantPAN) error {
			return errors.New("db error")
		},
	}
	h := NewHandlers(ms, nil)

	body := domain.CreatePANRequest{PAN: "ABCDE1234F", EntityName: "Acme", PANType: "company"}
	req := newRequest(t, http.MethodPost, "/v1/tenants/"+tenantID.String()+"/pans", body, map[string]string{
		"tenantID": tenantID.String(),
	})
	rec := httptest.NewRecorder()

	h.CreatePAN(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: CreateGSTIN
// ---------------------------------------------------------------------------

func TestCreateGSTIN_Success(t *testing.T) {
	tenantID := uuid.New()
	panID := uuid.New()
	ms := &mockStore{}
	h := NewHandlers(ms, nil)

	tradeName := "Acme Trading"
	body := domain.CreateGSTINRequest{
		GSTIN: "29ABCDE1234F1Z5", TradeName: &tradeName,
		StateCode: "29", RegistrationType: "regular",
	}
	req := newRequest(t, http.MethodPost, "/v1/tenants/"+tenantID.String()+"/pans/"+panID.String()+"/gstins", body, map[string]string{
		"tenantID": tenantID.String(),
		"panID":    panID.String(),
	})
	rec := httptest.NewRecorder()

	h.CreateGSTIN(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var gstin domain.TenantGSTIN
	parseDataResponse(t, rec.Body.Bytes(), &gstin)
	assert.Equal(t, "29ABCDE1234F1Z5", gstin.GSTIN)
	assert.Equal(t, panID, gstin.PANID)
	assert.Equal(t, tenantID, gstin.TenantID)
	assert.Equal(t, "29", gstin.StateCode)
	assert.Equal(t, "regular", gstin.RegistrationType)
	assert.Equal(t, "active", gstin.Status)
	require.NotNil(t, gstin.TradeName)
	assert.Equal(t, "Acme Trading", *gstin.TradeName)
}

func TestCreateGSTIN_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	panID := uuid.New()
	h := NewHandlers(&mockStore{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/tenants/"+tenantID.String()+"/pans/"+panID.String()+"/gstins", bytes.NewReader([]byte("bad")))
	req.SetPathValue("tenantID", tenantID.String())
	req.SetPathValue("panID", panID.String())
	rec := httptest.NewRecorder()

	h.CreateGSTIN(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateGSTIN_InvalidTenantID(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/tenants/bad/pans/"+uuid.New().String()+"/gstins", bytes.NewReader([]byte("{}")))
	req.SetPathValue("tenantID", "bad")
	req.SetPathValue("panID", uuid.New().String())
	rec := httptest.NewRecorder()

	h.CreateGSTIN(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid tenant_id", data["error"])
}

func TestCreateGSTIN_InvalidPanID(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/tenants/"+tenantID.String()+"/pans/bad/gstins", bytes.NewReader([]byte("{}")))
	req.SetPathValue("tenantID", tenantID.String())
	req.SetPathValue("panID", "bad")
	rec := httptest.NewRecorder()

	h.CreateGSTIN(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid pan_id", data["error"])
}

func TestCreateGSTIN_StoreError(t *testing.T) {
	tenantID := uuid.New()
	panID := uuid.New()
	ms := &mockStore{
		createGSTINFn: func(_ context.Context, _ uuid.UUID, _ *domain.TenantGSTIN) error {
			return errors.New("db error")
		},
	}
	h := NewHandlers(ms, nil)

	body := domain.CreateGSTINRequest{
		GSTIN: "29ABCDE1234F1Z5", StateCode: "29", RegistrationType: "regular",
	}
	req := newRequest(t, http.MethodPost, "/v1/tenants/"+tenantID.String()+"/pans/"+panID.String()+"/gstins", body, map[string]string{
		"tenantID": tenantID.String(),
		"panID":    panID.String(),
	})
	rec := httptest.NewRecorder()

	h.CreateGSTIN(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: CreateTAN
// ---------------------------------------------------------------------------

func TestCreateTAN_Success(t *testing.T) {
	tenantID := uuid.New()
	panID := uuid.New()
	ms := &mockStore{}
	h := NewHandlers(ms, nil)

	body := domain.CreateTANRequest{TAN: "BLRA12345F", DeductorName: "Acme Corp"}
	req := newRequest(t, http.MethodPost, "/v1/tenants/"+tenantID.String()+"/pans/"+panID.String()+"/tans", body, map[string]string{
		"tenantID": tenantID.String(),
		"panID":    panID.String(),
	})
	rec := httptest.NewRecorder()

	h.CreateTAN(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var tan domain.TenantTAN
	parseDataResponse(t, rec.Body.Bytes(), &tan)
	assert.Equal(t, "BLRA12345F", tan.TAN)
	assert.Equal(t, "Acme Corp", tan.DeductorName)
	assert.Equal(t, panID, tan.PANID)
	assert.Equal(t, tenantID, tan.TenantID)
	assert.Equal(t, "active", tan.Status)
}

func TestCreateTAN_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	panID := uuid.New()
	h := NewHandlers(&mockStore{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/tenants/"+tenantID.String()+"/pans/"+panID.String()+"/tans", bytes.NewReader([]byte("bad")))
	req.SetPathValue("tenantID", tenantID.String())
	req.SetPathValue("panID", panID.String())
	rec := httptest.NewRecorder()

	h.CreateTAN(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateTAN_InvalidTenantID(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/tenants/bad/pans/"+uuid.New().String()+"/tans", bytes.NewReader([]byte("{}")))
	req.SetPathValue("tenantID", "bad")
	req.SetPathValue("panID", uuid.New().String())
	rec := httptest.NewRecorder()

	h.CreateTAN(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateTAN_InvalidPanID(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/tenants/"+tenantID.String()+"/pans/bad/tans", bytes.NewReader([]byte("{}")))
	req.SetPathValue("tenantID", tenantID.String())
	req.SetPathValue("panID", "bad")
	rec := httptest.NewRecorder()

	h.CreateTAN(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid pan_id", data["error"])
}

func TestCreateTAN_StoreError(t *testing.T) {
	tenantID := uuid.New()
	panID := uuid.New()
	ms := &mockStore{
		createTANFn: func(_ context.Context, _ uuid.UUID, _ *domain.TenantTAN) error {
			return errors.New("db error")
		},
	}
	h := NewHandlers(ms, nil)

	body := domain.CreateTANRequest{TAN: "BLRA12345F", DeductorName: "Acme Corp"}
	req := newRequest(t, http.MethodPost, "/v1/tenants/"+tenantID.String()+"/pans/"+panID.String()+"/tans", body, map[string]string{
		"tenantID": tenantID.String(),
		"panID":    panID.String(),
	})
	rec := httptest.NewRecorder()

	h.CreateTAN(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: tenantIDFromRequest (exercised via ListTenants)
// ---------------------------------------------------------------------------

func TestTenantIDFromRequest_Valid(t *testing.T) {
	expected := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-Id", expected.String())

	got, err := tenantIDFromRequest(req)
	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestTenantIDFromRequest_Missing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	_, err := tenantIDFromRequest(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing X-Tenant-Id")
}

func TestTenantIDFromRequest_Invalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-Id", "not-uuid")

	_, err := tenantIDFromRequest(req)
	require.Error(t, err)
}

// ---------------------------------------------------------------------------
// Tests: strPtr helper
// ---------------------------------------------------------------------------

func TestStrPtr(t *testing.T) {
	s := strPtr("hello")
	require.NotNil(t, s)
	assert.Equal(t, "hello", *s)
}

// ---------------------------------------------------------------------------
// Mock KMS client
// ---------------------------------------------------------------------------

type mockKMS struct {
	createKeyFn   func(ctx context.Context, params *kms.CreateKeyInput, optFns ...func(*kms.Options)) (*kms.CreateKeyOutput, error)
	createAliasFn func(ctx context.Context, params *kms.CreateAliasInput, optFns ...func(*kms.Options)) (*kms.CreateAliasOutput, error)
}

func (m *mockKMS) CreateKey(ctx context.Context, params *kms.CreateKeyInput, optFns ...func(*kms.Options)) (*kms.CreateKeyOutput, error) {
	if m.createKeyFn != nil {
		return m.createKeyFn(ctx, params, optFns...)
	}
	keyID := uuid.New().String()
	arn := "arn:aws:kms:ap-south-1:000000000000:key/" + keyID
	return &kms.CreateKeyOutput{
		KeyMetadata: &kmsTypes.KeyMetadata{
			KeyId: &keyID,
			Arn:   &arn,
		},
	}, nil
}

func (m *mockKMS) CreateAlias(ctx context.Context, params *kms.CreateAliasInput, optFns ...func(*kms.Options)) (*kms.CreateAliasOutput, error) {
	if m.createAliasFn != nil {
		return m.createAliasFn(ctx, params, optFns...)
	}
	return &kms.CreateAliasOutput{}, nil
}

// ---------------------------------------------------------------------------
// Tests: CreateTenant with KMS
// ---------------------------------------------------------------------------

func TestCreateTenant_WithKMS_Success(t *testing.T) {
	var capturedKMSKeyARN string
	ms := &mockStore{
		updateTenantKMSFn: func(_ context.Context, _ uuid.UUID, arn string) error {
			capturedKMSKeyARN = arn
			return nil
		},
	}
	mk := &mockKMS{}
	h := NewHandlers(ms, mk)

	body := domain.CreateTenantRequest{Name: "Acme Corp", Slug: "acme", Tier: "pooled"}
	req := newRequest(t, http.MethodPost, "/v1/tenants", body, nil)
	rec := httptest.NewRecorder()

	h.CreateTenant(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var tenant domain.Tenant
	parseDataResponse(t, rec.Body.Bytes(), &tenant)
	require.NotNil(t, tenant.KMSKeyARN)
	assert.Contains(t, *tenant.KMSKeyARN, "arn:aws:kms")
	assert.Equal(t, *tenant.KMSKeyARN, capturedKMSKeyARN)
}

func TestCreateTenant_WithKMS_CreateKeyError(t *testing.T) {
	ms := &mockStore{}
	mk := &mockKMS{
		createKeyFn: func(_ context.Context, _ *kms.CreateKeyInput, _ ...func(*kms.Options)) (*kms.CreateKeyOutput, error) {
			return nil, errors.New("kms unavailable")
		},
	}
	h := NewHandlers(ms, mk)

	body := domain.CreateTenantRequest{Name: "Acme", Slug: "acme", Tier: "pooled"}
	req := newRequest(t, http.MethodPost, "/v1/tenants", body, nil)
	rec := httptest.NewRecorder()

	h.CreateTenant(rec, req)

	// Tenant still created even if KMS fails
	assert.Equal(t, http.StatusCreated, rec.Code)
	var tenant domain.Tenant
	parseDataResponse(t, rec.Body.Bytes(), &tenant)
	assert.Nil(t, tenant.KMSKeyARN)
}

func TestCreateTenant_WithKMS_CreateAliasError(t *testing.T) {
	ms := &mockStore{}
	mk := &mockKMS{
		createAliasFn: func(_ context.Context, _ *kms.CreateAliasInput, _ ...func(*kms.Options)) (*kms.CreateAliasOutput, error) {
			return nil, errors.New("alias failed")
		},
	}
	h := NewHandlers(ms, mk)

	body := domain.CreateTenantRequest{Name: "Acme", Slug: "acme", Tier: "pooled"}
	req := newRequest(t, http.MethodPost, "/v1/tenants", body, nil)
	rec := httptest.NewRecorder()

	h.CreateTenant(rec, req)

	// Alias failure is non-fatal; KMS key ARN should still be set
	assert.Equal(t, http.StatusCreated, rec.Code)
	var tenant domain.Tenant
	parseDataResponse(t, rec.Body.Bytes(), &tenant)
	require.NotNil(t, tenant.KMSKeyARN)
	assert.Contains(t, *tenant.KMSKeyARN, "arn:aws:kms")
}

func TestCreateTenant_WithKMS_UpdateKMSKeyStoreError(t *testing.T) {
	ms := &mockStore{
		updateTenantKMSFn: func(_ context.Context, _ uuid.UUID, _ string) error {
			return errors.New("store error")
		},
	}
	mk := &mockKMS{}
	h := NewHandlers(ms, mk)

	body := domain.CreateTenantRequest{Name: "Acme", Slug: "acme", Tier: "pooled"}
	req := newRequest(t, http.MethodPost, "/v1/tenants", body, nil)
	rec := httptest.NewRecorder()

	h.CreateTenant(rec, req)

	// Still returns 201 even if UpdateKMSKey fails (it's logged but not fatal)
	assert.Equal(t, http.StatusCreated, rec.Code)
	var tenant domain.Tenant
	parseDataResponse(t, rec.Body.Bytes(), &tenant)
	require.NotNil(t, tenant.KMSKeyARN)
}

// ---------------------------------------------------------------------------
// Tests: NewRouter
// ---------------------------------------------------------------------------

func TestNewRouter(t *testing.T) {
	ms := &mockStore{}
	r := NewRouter(ms, nil)
	require.NotNil(t, r)

	// Verify health endpoint is reachable through the router
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify ping heartbeat endpoint
	req = httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}
