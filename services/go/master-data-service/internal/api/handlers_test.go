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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/master-data-service/internal/domain"
)

// ---------------------------------------------------------------------------
// Mock store
// ---------------------------------------------------------------------------

type mockStore struct {
	createVendorFn   func(ctx context.Context, tenantID uuid.UUID, v *domain.Vendor) error
	getVendorFn      func(ctx context.Context, tenantID uuid.UUID, vendorID uuid.UUID) (*domain.Vendor, error)
	listVendorsFn    func(ctx context.Context, tenantID uuid.UUID) ([]domain.Vendor, error)
	updateVendorFn   func(ctx context.Context, tenantID uuid.UUID, vendorID uuid.UUID, req *domain.UpdateVendorRequest) (*domain.Vendor, error)
	createCustomerFn func(ctx context.Context, tenantID uuid.UUID, c *domain.Customer) error
	getCustomerFn    func(ctx context.Context, tenantID uuid.UUID, customerID uuid.UUID) (*domain.Customer, error)
	listCustomersFn  func(ctx context.Context, tenantID uuid.UUID) ([]domain.Customer, error)
	createItemFn     func(ctx context.Context, tenantID uuid.UUID, i *domain.Item) error
	getItemFn        func(ctx context.Context, tenantID uuid.UUID, itemID uuid.UUID) (*domain.Item, error)
	listItemsFn      func(ctx context.Context, tenantID uuid.UUID) ([]domain.Item, error)
	listHSNCodesFn   func(ctx context.Context, tenantID uuid.UUID) ([]domain.HSNCode, error)
	getHSNCodeFn     func(ctx context.Context, tenantID uuid.UUID, hsnID uuid.UUID) (*domain.HSNCode, error)
	createHSNCodeFn  func(ctx context.Context, tenantID uuid.UUID, h *domain.HSNCode) error
	listStateCodesFn func(ctx context.Context, tenantID uuid.UUID) ([]domain.StateCode, error)
}

func (m *mockStore) CreateVendor(ctx context.Context, tenantID uuid.UUID, v *domain.Vendor) error {
	if m.createVendorFn != nil {
		return m.createVendorFn(ctx, tenantID, v)
	}
	v.ID = uuid.New()
	v.TenantID = tenantID
	v.KYCStatus = "pending"
	v.ComplianceScore = 0
	v.Status = "active"
	v.Metadata = "{}"
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()
	return nil
}

func (m *mockStore) GetVendor(ctx context.Context, tenantID uuid.UUID, vendorID uuid.UUID) (*domain.Vendor, error) {
	if m.getVendorFn != nil {
		return m.getVendorFn(ctx, tenantID, vendorID)
	}
	return nil, errors.New("not found")
}

func (m *mockStore) ListVendors(ctx context.Context, tenantID uuid.UUID) ([]domain.Vendor, error) {
	if m.listVendorsFn != nil {
		return m.listVendorsFn(ctx, tenantID)
	}
	return nil, nil
}

func (m *mockStore) UpdateVendor(ctx context.Context, tenantID uuid.UUID, vendorID uuid.UUID, req *domain.UpdateVendorRequest) (*domain.Vendor, error) {
	if m.updateVendorFn != nil {
		return m.updateVendorFn(ctx, tenantID, vendorID, req)
	}
	return &domain.Vendor{
		ID:        vendorID,
		TenantID:  tenantID,
		Name:      "Updated Vendor",
		Status:    "active",
		KYCStatus: "pending",
		Metadata:  "{}",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *mockStore) CreateCustomer(ctx context.Context, tenantID uuid.UUID, c *domain.Customer) error {
	if m.createCustomerFn != nil {
		return m.createCustomerFn(ctx, tenantID, c)
	}
	c.ID = uuid.New()
	c.TenantID = tenantID
	c.Status = "active"
	c.Metadata = "{}"
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

func (m *mockStore) GetCustomer(ctx context.Context, tenantID uuid.UUID, customerID uuid.UUID) (*domain.Customer, error) {
	if m.getCustomerFn != nil {
		return m.getCustomerFn(ctx, tenantID, customerID)
	}
	return nil, errors.New("not found")
}

func (m *mockStore) ListCustomers(ctx context.Context, tenantID uuid.UUID) ([]domain.Customer, error) {
	if m.listCustomersFn != nil {
		return m.listCustomersFn(ctx, tenantID)
	}
	return nil, nil
}

func (m *mockStore) CreateItem(ctx context.Context, tenantID uuid.UUID, i *domain.Item) error {
	if m.createItemFn != nil {
		return m.createItemFn(ctx, tenantID, i)
	}
	i.ID = uuid.New()
	i.TenantID = tenantID
	i.Status = "active"
	i.Metadata = "{}"
	i.CreatedAt = time.Now()
	i.UpdatedAt = time.Now()
	return nil
}

func (m *mockStore) GetItem(ctx context.Context, tenantID uuid.UUID, itemID uuid.UUID) (*domain.Item, error) {
	if m.getItemFn != nil {
		return m.getItemFn(ctx, tenantID, itemID)
	}
	return nil, errors.New("not found")
}

func (m *mockStore) ListItems(ctx context.Context, tenantID uuid.UUID) ([]domain.Item, error) {
	if m.listItemsFn != nil {
		return m.listItemsFn(ctx, tenantID)
	}
	return nil, nil
}

func (m *mockStore) ListHSNCodes(ctx context.Context, tenantID uuid.UUID) ([]domain.HSNCode, error) {
	if m.listHSNCodesFn != nil {
		return m.listHSNCodesFn(ctx, tenantID)
	}
	return nil, nil
}

func (m *mockStore) GetHSNCode(ctx context.Context, tenantID uuid.UUID, hsnID uuid.UUID) (*domain.HSNCode, error) {
	if m.getHSNCodeFn != nil {
		return m.getHSNCodeFn(ctx, tenantID, hsnID)
	}
	return nil, errors.New("not found")
}

func (m *mockStore) CreateHSNCode(ctx context.Context, tenantID uuid.UUID, h *domain.HSNCode) error {
	if m.createHSNCodeFn != nil {
		return m.createHSNCodeFn(ctx, tenantID, h)
	}
	h.ID = uuid.New()
	h.TenantID = tenantID
	h.EffectiveFrom = "2024-01-01"
	h.CreatedAt = time.Now()
	h.UpdatedAt = time.Now()
	return nil
}

func (m *mockStore) ListStateCodes(ctx context.Context, tenantID uuid.UUID) ([]domain.StateCode, error) {
	if m.listStateCodesFn != nil {
		return m.listStateCodesFn(ctx, tenantID)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func parseDataResponse(t *testing.T, body []byte, target interface{}) {
	t.Helper()
	var wrapper httputil.SuccessResponse
	wrapper.Data = target
	require.NoError(t, json.Unmarshal(body, &wrapper))
}

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

func reqWithTenant(t *testing.T, method, path string, tenantID uuid.UUID, body interface{}, pathValues map[string]string) *http.Request {
	t.Helper()
	req := newRequest(t, method, path, body, pathValues)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	return req
}

// ---------------------------------------------------------------------------
// Tests: Health
// ---------------------------------------------------------------------------

func TestHealth(t *testing.T) {
	h := NewHandlers(&mockStore{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	h.Health(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "ok", data["status"])
	assert.Equal(t, "master-data-service", data["service"])
}

// ---------------------------------------------------------------------------
// Tests: CreateVendor
// ---------------------------------------------------------------------------

func TestCreateVendor_Success(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{}
	h := NewHandlers(ms)

	body := domain.CreateVendorRequest{Name: "Vendor Corp"}
	req := reqWithTenant(t, http.MethodPost, "/v1/vendors", tenantID, body, nil)
	rec := httptest.NewRecorder()

	h.CreateVendor(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var vendor domain.Vendor
	parseDataResponse(t, rec.Body.Bytes(), &vendor)
	assert.Equal(t, "Vendor Corp", vendor.Name)
	assert.Equal(t, "active", vendor.Status)
	assert.Equal(t, "pending", vendor.KYCStatus)
	assert.Equal(t, tenantID, vendor.TenantID)
	assert.NotEqual(t, uuid.Nil, vendor.ID)
}

func TestCreateVendor_MissingTenantHeader(t *testing.T) {
	h := NewHandlers(&mockStore{})
	body := domain.CreateVendorRequest{Name: "Vendor Corp"}
	req := newRequest(t, http.MethodPost, "/v1/vendors", body, nil)
	rec := httptest.NewRecorder()

	h.CreateVendor(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Contains(t, data["error"], "missing X-Tenant-Id")
}

func TestCreateVendor_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})
	req := httptest.NewRequest(http.MethodPost, "/v1/vendors", bytes.NewReader([]byte("not json")))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.CreateVendor(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid request", data["error"])
}

func TestCreateVendor_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		createVendorFn: func(_ context.Context, _ uuid.UUID, _ *domain.Vendor) error {
			return errors.New("db down")
		},
	}
	h := NewHandlers(ms)

	body := domain.CreateVendorRequest{Name: "Vendor Corp"}
	req := reqWithTenant(t, http.MethodPost, "/v1/vendors", tenantID, body, nil)
	rec := httptest.NewRecorder()

	h.CreateVendor(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "create failed", data["error"])
}

// ---------------------------------------------------------------------------
// Tests: GetVendor
// ---------------------------------------------------------------------------

func TestGetVendor_Success(t *testing.T) {
	tenantID := uuid.New()
	vendorID := uuid.New()
	ms := &mockStore{
		getVendorFn: func(_ context.Context, tid uuid.UUID, vid uuid.UUID) (*domain.Vendor, error) {
			assert.Equal(t, tenantID, tid)
			assert.Equal(t, vendorID, vid)
			return &domain.Vendor{
				ID: vendorID, TenantID: tenantID, Name: "Vendor Corp",
				KYCStatus: "pending", Status: "active", Metadata: "{}",
			}, nil
		},
	}
	h := NewHandlers(ms)

	req := reqWithTenant(t, http.MethodGet, "/v1/vendors/"+vendorID.String(), tenantID, nil, map[string]string{
		"vendorID": vendorID.String(),
	})
	rec := httptest.NewRecorder()

	h.GetVendor(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var vendor domain.Vendor
	parseDataResponse(t, rec.Body.Bytes(), &vendor)
	assert.Equal(t, vendorID, vendor.ID)
	assert.Equal(t, "Vendor Corp", vendor.Name)
}

func TestGetVendor_InvalidID(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	req := reqWithTenant(t, http.MethodGet, "/v1/vendors/not-a-uuid", tenantID, nil, map[string]string{
		"vendorID": "not-a-uuid",
	})
	rec := httptest.NewRecorder()

	h.GetVendor(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid vendor_id", data["error"])
}

func TestGetVendor_NotFound(t *testing.T) {
	tenantID := uuid.New()
	vendorID := uuid.New()
	ms := &mockStore{
		getVendorFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.Vendor, error) {
			return nil, errors.New("not found")
		},
	}
	h := NewHandlers(ms)

	req := reqWithTenant(t, http.MethodGet, "/v1/vendors/"+vendorID.String(), tenantID, nil, map[string]string{
		"vendorID": vendorID.String(),
	})
	rec := httptest.NewRecorder()

	h.GetVendor(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetVendor_MissingTenantHeader(t *testing.T) {
	h := NewHandlers(&mockStore{})
	vendorID := uuid.New()
	req := newRequest(t, http.MethodGet, "/v1/vendors/"+vendorID.String(), nil, map[string]string{
		"vendorID": vendorID.String(),
	})
	rec := httptest.NewRecorder()

	h.GetVendor(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: ListVendors
// ---------------------------------------------------------------------------

func TestListVendors_Success(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listVendorsFn: func(_ context.Context, tid uuid.UUID) ([]domain.Vendor, error) {
			assert.Equal(t, tenantID, tid)
			return []domain.Vendor{
				{ID: uuid.New(), TenantID: tenantID, Name: "Vendor A", KYCStatus: "pending", Status: "active", Metadata: "{}"},
			}, nil
		},
	}
	h := NewHandlers(ms)

	req := reqWithTenant(t, http.MethodGet, "/v1/vendors", tenantID, nil, nil)
	rec := httptest.NewRecorder()

	h.ListVendors(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var vendors []domain.Vendor
	parseDataResponse(t, rec.Body.Bytes(), &vendors)
	require.Len(t, vendors, 1)
	assert.Equal(t, "Vendor A", vendors[0].Name)
}

func TestListVendors_MissingHeader(t *testing.T) {
	h := NewHandlers(&mockStore{})
	req := httptest.NewRequest(http.MethodGet, "/v1/vendors", nil)
	rec := httptest.NewRecorder()

	h.ListVendors(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListVendors_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listVendorsFn: func(_ context.Context, _ uuid.UUID) ([]domain.Vendor, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewHandlers(ms)

	req := reqWithTenant(t, http.MethodGet, "/v1/vendors", tenantID, nil, nil)
	rec := httptest.NewRecorder()

	h.ListVendors(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "internal error", data["error"])
}

func TestListVendors_NilResultReturnsEmptyArray(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listVendorsFn: func(_ context.Context, _ uuid.UUID) ([]domain.Vendor, error) {
			return nil, nil
		},
	}
	h := NewHandlers(ms)

	req := reqWithTenant(t, http.MethodGet, "/v1/vendors", tenantID, nil, nil)
	rec := httptest.NewRecorder()

	h.ListVendors(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var vendors []domain.Vendor
	parseDataResponse(t, rec.Body.Bytes(), &vendors)
	assert.NotNil(t, vendors)
	assert.Len(t, vendors, 0)
}

// ---------------------------------------------------------------------------
// Tests: UpdateVendor
// ---------------------------------------------------------------------------

func TestUpdateVendor_Success(t *testing.T) {
	tenantID := uuid.New()
	vendorID := uuid.New()
	ms := &mockStore{
		updateVendorFn: func(_ context.Context, tid uuid.UUID, vid uuid.UUID, req *domain.UpdateVendorRequest) (*domain.Vendor, error) {
			assert.Equal(t, tenantID, tid)
			assert.Equal(t, vendorID, vid)
			return &domain.Vendor{
				ID: vendorID, TenantID: tenantID, Name: *req.Name,
				Status: "active", KYCStatus: "pending", Metadata: "{}",
				CreatedAt: time.Now(), UpdatedAt: time.Now(),
			}, nil
		},
	}
	h := NewHandlers(ms)

	name := "Updated Vendor"
	body := domain.UpdateVendorRequest{Name: &name}
	req := reqWithTenant(t, http.MethodPut, "/v1/vendors/"+vendorID.String(), tenantID, body, map[string]string{
		"vendorID": vendorID.String(),
	})
	rec := httptest.NewRecorder()

	h.UpdateVendor(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var vendor domain.Vendor
	parseDataResponse(t, rec.Body.Bytes(), &vendor)
	assert.Equal(t, "Updated Vendor", vendor.Name)
}

func TestUpdateVendor_InvalidID(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	name := "Updated"
	body := domain.UpdateVendorRequest{Name: &name}
	req := reqWithTenant(t, http.MethodPut, "/v1/vendors/bad", tenantID, body, map[string]string{
		"vendorID": "bad",
	})
	rec := httptest.NewRecorder()

	h.UpdateVendor(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateVendor_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	vendorID := uuid.New()
	h := NewHandlers(&mockStore{})

	req := httptest.NewRequest(http.MethodPut, "/v1/vendors/"+vendorID.String(), bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.SetPathValue("vendorID", vendorID.String())
	rec := httptest.NewRecorder()

	h.UpdateVendor(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateVendor_StoreError(t *testing.T) {
	tenantID := uuid.New()
	vendorID := uuid.New()
	ms := &mockStore{
		updateVendorFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ *domain.UpdateVendorRequest) (*domain.Vendor, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewHandlers(ms)

	name := "Updated"
	body := domain.UpdateVendorRequest{Name: &name}
	req := reqWithTenant(t, http.MethodPut, "/v1/vendors/"+vendorID.String(), tenantID, body, map[string]string{
		"vendorID": vendorID.String(),
	})
	rec := httptest.NewRecorder()

	h.UpdateVendor(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "update failed", data["error"])
}

func TestUpdateVendor_MissingTenantHeader(t *testing.T) {
	h := NewHandlers(&mockStore{})
	vendorID := uuid.New()
	name := "Updated"
	body := domain.UpdateVendorRequest{Name: &name}
	req := newRequest(t, http.MethodPut, "/v1/vendors/"+vendorID.String(), body, map[string]string{
		"vendorID": vendorID.String(),
	})
	rec := httptest.NewRecorder()

	h.UpdateVendor(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: CreateCustomer
// ---------------------------------------------------------------------------

func TestCreateCustomer_Success(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{}
	h := NewHandlers(ms)

	body := domain.CreateCustomerRequest{Name: "Customer Corp"}
	req := reqWithTenant(t, http.MethodPost, "/v1/customers", tenantID, body, nil)
	rec := httptest.NewRecorder()

	h.CreateCustomer(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var customer domain.Customer
	parseDataResponse(t, rec.Body.Bytes(), &customer)
	assert.Equal(t, "Customer Corp", customer.Name)
	assert.Equal(t, "active", customer.Status)
	assert.Equal(t, tenantID, customer.TenantID)
	assert.NotEqual(t, uuid.Nil, customer.ID)
}

func TestCreateCustomer_MissingTenantHeader(t *testing.T) {
	h := NewHandlers(&mockStore{})
	body := domain.CreateCustomerRequest{Name: "Customer Corp"}
	req := newRequest(t, http.MethodPost, "/v1/customers", body, nil)
	rec := httptest.NewRecorder()

	h.CreateCustomer(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateCustomer_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})
	req := httptest.NewRequest(http.MethodPost, "/v1/customers", bytes.NewReader([]byte("not json")))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.CreateCustomer(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid request", data["error"])
}

func TestCreateCustomer_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		createCustomerFn: func(_ context.Context, _ uuid.UUID, _ *domain.Customer) error {
			return errors.New("db down")
		},
	}
	h := NewHandlers(ms)

	body := domain.CreateCustomerRequest{Name: "Customer Corp"}
	req := reqWithTenant(t, http.MethodPost, "/v1/customers", tenantID, body, nil)
	rec := httptest.NewRecorder()

	h.CreateCustomer(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "create failed", data["error"])
}

// ---------------------------------------------------------------------------
// Tests: GetCustomer
// ---------------------------------------------------------------------------

func TestGetCustomer_Success(t *testing.T) {
	tenantID := uuid.New()
	customerID := uuid.New()
	ms := &mockStore{
		getCustomerFn: func(_ context.Context, tid uuid.UUID, cid uuid.UUID) (*domain.Customer, error) {
			assert.Equal(t, tenantID, tid)
			assert.Equal(t, customerID, cid)
			return &domain.Customer{
				ID: customerID, TenantID: tenantID, Name: "Customer Corp",
				Status: "active", Metadata: "{}",
			}, nil
		},
	}
	h := NewHandlers(ms)

	req := reqWithTenant(t, http.MethodGet, "/v1/customers/"+customerID.String(), tenantID, nil, map[string]string{
		"customerID": customerID.String(),
	})
	rec := httptest.NewRecorder()

	h.GetCustomer(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var customer domain.Customer
	parseDataResponse(t, rec.Body.Bytes(), &customer)
	assert.Equal(t, customerID, customer.ID)
	assert.Equal(t, "Customer Corp", customer.Name)
}

func TestGetCustomer_InvalidID(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	req := reqWithTenant(t, http.MethodGet, "/v1/customers/not-a-uuid", tenantID, nil, map[string]string{
		"customerID": "not-a-uuid",
	})
	rec := httptest.NewRecorder()

	h.GetCustomer(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid customer_id", data["error"])
}

func TestGetCustomer_NotFound(t *testing.T) {
	tenantID := uuid.New()
	customerID := uuid.New()
	ms := &mockStore{
		getCustomerFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.Customer, error) {
			return nil, errors.New("not found")
		},
	}
	h := NewHandlers(ms)

	req := reqWithTenant(t, http.MethodGet, "/v1/customers/"+customerID.String(), tenantID, nil, map[string]string{
		"customerID": customerID.String(),
	})
	rec := httptest.NewRecorder()

	h.GetCustomer(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: ListCustomers
// ---------------------------------------------------------------------------

func TestListCustomers_Success(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listCustomersFn: func(_ context.Context, tid uuid.UUID) ([]domain.Customer, error) {
			assert.Equal(t, tenantID, tid)
			return []domain.Customer{
				{ID: uuid.New(), TenantID: tenantID, Name: "Customer A", Status: "active", Metadata: "{}"},
			}, nil
		},
	}
	h := NewHandlers(ms)

	req := reqWithTenant(t, http.MethodGet, "/v1/customers", tenantID, nil, nil)
	rec := httptest.NewRecorder()

	h.ListCustomers(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var customers []domain.Customer
	parseDataResponse(t, rec.Body.Bytes(), &customers)
	require.Len(t, customers, 1)
	assert.Equal(t, "Customer A", customers[0].Name)
}

func TestListCustomers_MissingHeader(t *testing.T) {
	h := NewHandlers(&mockStore{})
	req := httptest.NewRequest(http.MethodGet, "/v1/customers", nil)
	rec := httptest.NewRecorder()

	h.ListCustomers(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListCustomers_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listCustomersFn: func(_ context.Context, _ uuid.UUID) ([]domain.Customer, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewHandlers(ms)

	req := reqWithTenant(t, http.MethodGet, "/v1/customers", tenantID, nil, nil)
	rec := httptest.NewRecorder()

	h.ListCustomers(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "internal error", data["error"])
}

func TestListCustomers_NilResultReturnsEmptyArray(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listCustomersFn: func(_ context.Context, _ uuid.UUID) ([]domain.Customer, error) {
			return nil, nil
		},
	}
	h := NewHandlers(ms)

	req := reqWithTenant(t, http.MethodGet, "/v1/customers", tenantID, nil, nil)
	rec := httptest.NewRecorder()

	h.ListCustomers(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var customers []domain.Customer
	parseDataResponse(t, rec.Body.Bytes(), &customers)
	assert.NotNil(t, customers)
	assert.Len(t, customers, 0)
}

// ---------------------------------------------------------------------------
// Tests: CreateItem
// ---------------------------------------------------------------------------

func TestCreateItem_Success(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{}
	h := NewHandlers(ms)

	body := domain.CreateItemRequest{Name: "Widget", HSNCode: "8471", UnitOfMeasure: "NOS"}
	req := reqWithTenant(t, http.MethodPost, "/v1/items", tenantID, body, nil)
	rec := httptest.NewRecorder()

	h.CreateItem(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var item domain.Item
	parseDataResponse(t, rec.Body.Bytes(), &item)
	assert.Equal(t, "Widget", item.Name)
	assert.Equal(t, "8471", item.HSNCode)
	assert.Equal(t, "NOS", item.UnitOfMeasure)
	assert.Equal(t, "active", item.Status)
	assert.Equal(t, tenantID, item.TenantID)
}

func TestCreateItem_DefaultUOM(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{}
	h := NewHandlers(ms)

	body := domain.CreateItemRequest{Name: "Widget", HSNCode: "8471"}
	req := reqWithTenant(t, http.MethodPost, "/v1/items", tenantID, body, nil)
	rec := httptest.NewRecorder()

	h.CreateItem(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var item domain.Item
	parseDataResponse(t, rec.Body.Bytes(), &item)
	assert.Equal(t, "NOS", item.UnitOfMeasure)
}

func TestCreateItem_MissingTenantHeader(t *testing.T) {
	h := NewHandlers(&mockStore{})
	body := domain.CreateItemRequest{Name: "Widget", HSNCode: "8471"}
	req := newRequest(t, http.MethodPost, "/v1/items", body, nil)
	rec := httptest.NewRecorder()

	h.CreateItem(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateItem_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})
	req := httptest.NewRequest(http.MethodPost, "/v1/items", bytes.NewReader([]byte("not json")))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.CreateItem(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateItem_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		createItemFn: func(_ context.Context, _ uuid.UUID, _ *domain.Item) error {
			return errors.New("db down")
		},
	}
	h := NewHandlers(ms)

	body := domain.CreateItemRequest{Name: "Widget", HSNCode: "8471"}
	req := reqWithTenant(t, http.MethodPost, "/v1/items", tenantID, body, nil)
	rec := httptest.NewRecorder()

	h.CreateItem(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: GetItem
// ---------------------------------------------------------------------------

func TestGetItem_Success(t *testing.T) {
	tenantID := uuid.New()
	itemID := uuid.New()
	ms := &mockStore{
		getItemFn: func(_ context.Context, tid uuid.UUID, iid uuid.UUID) (*domain.Item, error) {
			assert.Equal(t, tenantID, tid)
			assert.Equal(t, itemID, iid)
			return &domain.Item{
				ID: itemID, TenantID: tenantID, Name: "Widget",
				HSNCode: "8471", UnitOfMeasure: "NOS", Status: "active", Metadata: "{}",
			}, nil
		},
	}
	h := NewHandlers(ms)

	req := reqWithTenant(t, http.MethodGet, "/v1/items/"+itemID.String(), tenantID, nil, map[string]string{
		"itemID": itemID.String(),
	})
	rec := httptest.NewRecorder()

	h.GetItem(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var item domain.Item
	parseDataResponse(t, rec.Body.Bytes(), &item)
	assert.Equal(t, itemID, item.ID)
	assert.Equal(t, "Widget", item.Name)
}

func TestGetItem_InvalidID(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	req := reqWithTenant(t, http.MethodGet, "/v1/items/not-a-uuid", tenantID, nil, map[string]string{
		"itemID": "not-a-uuid",
	})
	rec := httptest.NewRecorder()

	h.GetItem(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid item_id", data["error"])
}

func TestGetItem_NotFound(t *testing.T) {
	tenantID := uuid.New()
	itemID := uuid.New()
	ms := &mockStore{
		getItemFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.Item, error) {
			return nil, errors.New("not found")
		},
	}
	h := NewHandlers(ms)

	req := reqWithTenant(t, http.MethodGet, "/v1/items/"+itemID.String(), tenantID, nil, map[string]string{
		"itemID": itemID.String(),
	})
	rec := httptest.NewRecorder()

	h.GetItem(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: ListItems
// ---------------------------------------------------------------------------

func TestListItems_Success(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listItemsFn: func(_ context.Context, tid uuid.UUID) ([]domain.Item, error) {
			assert.Equal(t, tenantID, tid)
			return []domain.Item{
				{ID: uuid.New(), TenantID: tenantID, Name: "Widget", HSNCode: "8471", UnitOfMeasure: "NOS", Status: "active", Metadata: "{}"},
			}, nil
		},
	}
	h := NewHandlers(ms)

	req := reqWithTenant(t, http.MethodGet, "/v1/items", tenantID, nil, nil)
	rec := httptest.NewRecorder()

	h.ListItems(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var items []domain.Item
	parseDataResponse(t, rec.Body.Bytes(), &items)
	require.Len(t, items, 1)
	assert.Equal(t, "Widget", items[0].Name)
}

func TestListItems_MissingHeader(t *testing.T) {
	h := NewHandlers(&mockStore{})
	req := httptest.NewRequest(http.MethodGet, "/v1/items", nil)
	rec := httptest.NewRecorder()

	h.ListItems(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListItems_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listItemsFn: func(_ context.Context, _ uuid.UUID) ([]domain.Item, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewHandlers(ms)

	req := reqWithTenant(t, http.MethodGet, "/v1/items", tenantID, nil, nil)
	rec := httptest.NewRecorder()

	h.ListItems(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestListItems_NilResultReturnsEmptyArray(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listItemsFn: func(_ context.Context, _ uuid.UUID) ([]domain.Item, error) {
			return nil, nil
		},
	}
	h := NewHandlers(ms)

	req := reqWithTenant(t, http.MethodGet, "/v1/items", tenantID, nil, nil)
	rec := httptest.NewRecorder()

	h.ListItems(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var items []domain.Item
	parseDataResponse(t, rec.Body.Bytes(), &items)
	assert.NotNil(t, items)
	assert.Len(t, items, 0)
}

// ---------------------------------------------------------------------------
// Tests: ListHSNCodes
// ---------------------------------------------------------------------------

func TestListHSNCodes_Success(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listHSNCodesFn: func(_ context.Context, tid uuid.UUID) ([]domain.HSNCode, error) {
			assert.Equal(t, tenantID, tid)
			return []domain.HSNCode{
				{ID: uuid.New(), TenantID: tenantID, Code: "8471", Description: "Computers", GSTRate: 18.0, EffectiveFrom: "2024-01-01"},
			}, nil
		},
	}
	h := NewHandlers(ms)

	req := reqWithTenant(t, http.MethodGet, "/v1/hsn-codes", tenantID, nil, nil)
	rec := httptest.NewRecorder()

	h.ListHSNCodes(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var codes []domain.HSNCode
	parseDataResponse(t, rec.Body.Bytes(), &codes)
	require.Len(t, codes, 1)
	assert.Equal(t, "8471", codes[0].Code)
}

func TestListHSNCodes_MissingHeader(t *testing.T) {
	h := NewHandlers(&mockStore{})
	req := httptest.NewRequest(http.MethodGet, "/v1/hsn-codes", nil)
	rec := httptest.NewRecorder()

	h.ListHSNCodes(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListHSNCodes_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listHSNCodesFn: func(_ context.Context, _ uuid.UUID) ([]domain.HSNCode, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewHandlers(ms)

	req := reqWithTenant(t, http.MethodGet, "/v1/hsn-codes", tenantID, nil, nil)
	rec := httptest.NewRecorder()

	h.ListHSNCodes(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestListHSNCodes_NilResultReturnsEmptyArray(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listHSNCodesFn: func(_ context.Context, _ uuid.UUID) ([]domain.HSNCode, error) {
			return nil, nil
		},
	}
	h := NewHandlers(ms)

	req := reqWithTenant(t, http.MethodGet, "/v1/hsn-codes", tenantID, nil, nil)
	rec := httptest.NewRecorder()

	h.ListHSNCodes(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var codes []domain.HSNCode
	parseDataResponse(t, rec.Body.Bytes(), &codes)
	assert.NotNil(t, codes)
	assert.Len(t, codes, 0)
}

// ---------------------------------------------------------------------------
// Tests: CreateHSNCode
// ---------------------------------------------------------------------------

func TestCreateHSNCode_Success(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{}
	h := NewHandlers(ms)

	body := domain.CreateHSNCodeRequest{Code: "8471", Description: "Computers", GSTRate: 18.0, EffectiveFrom: "2024-01-01"}
	req := reqWithTenant(t, http.MethodPost, "/v1/hsn-codes", tenantID, body, nil)
	rec := httptest.NewRecorder()

	h.CreateHSNCode(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var hsn domain.HSNCode
	parseDataResponse(t, rec.Body.Bytes(), &hsn)
	assert.Equal(t, "8471", hsn.Code)
	assert.Equal(t, "Computers", hsn.Description)
	assert.Equal(t, 18.0, hsn.GSTRate)
	assert.Equal(t, tenantID, hsn.TenantID)
}

func TestCreateHSNCode_MissingTenantHeader(t *testing.T) {
	h := NewHandlers(&mockStore{})
	body := domain.CreateHSNCodeRequest{Code: "8471", Description: "Computers", GSTRate: 18.0}
	req := newRequest(t, http.MethodPost, "/v1/hsn-codes", body, nil)
	rec := httptest.NewRecorder()

	h.CreateHSNCode(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateHSNCode_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})
	req := httptest.NewRequest(http.MethodPost, "/v1/hsn-codes", bytes.NewReader([]byte("not json")))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.CreateHSNCode(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateHSNCode_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		createHSNCodeFn: func(_ context.Context, _ uuid.UUID, _ *domain.HSNCode) error {
			return errors.New("db down")
		},
	}
	h := NewHandlers(ms)

	body := domain.CreateHSNCodeRequest{Code: "8471", Description: "Computers", GSTRate: 18.0}
	req := reqWithTenant(t, http.MethodPost, "/v1/hsn-codes", tenantID, body, nil)
	rec := httptest.NewRecorder()

	h.CreateHSNCode(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: ListStateCodes
// ---------------------------------------------------------------------------

func TestListStateCodes_Success(t *testing.T) {
	tenantID := uuid.New()
	tinCode := "29"
	ms := &mockStore{
		listStateCodesFn: func(_ context.Context, tid uuid.UUID) ([]domain.StateCode, error) {
			assert.Equal(t, tenantID, tid)
			return []domain.StateCode{
				{ID: uuid.New(), TenantID: tenantID, Code: "29", Name: "Karnataka", TINCode: &tinCode},
			}, nil
		},
	}
	h := NewHandlers(ms)

	req := reqWithTenant(t, http.MethodGet, "/v1/state-codes", tenantID, nil, nil)
	rec := httptest.NewRecorder()

	h.ListStateCodes(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var codes []domain.StateCode
	parseDataResponse(t, rec.Body.Bytes(), &codes)
	require.Len(t, codes, 1)
	assert.Equal(t, "29", codes[0].Code)
	assert.Equal(t, "Karnataka", codes[0].Name)
}

func TestListStateCodes_MissingHeader(t *testing.T) {
	h := NewHandlers(&mockStore{})
	req := httptest.NewRequest(http.MethodGet, "/v1/state-codes", nil)
	rec := httptest.NewRecorder()

	h.ListStateCodes(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListStateCodes_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listStateCodesFn: func(_ context.Context, _ uuid.UUID) ([]domain.StateCode, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewHandlers(ms)

	req := reqWithTenant(t, http.MethodGet, "/v1/state-codes", tenantID, nil, nil)
	rec := httptest.NewRecorder()

	h.ListStateCodes(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestListStateCodes_NilResultReturnsEmptyArray(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listStateCodesFn: func(_ context.Context, _ uuid.UUID) ([]domain.StateCode, error) {
			return nil, nil
		},
	}
	h := NewHandlers(ms)

	req := reqWithTenant(t, http.MethodGet, "/v1/state-codes", tenantID, nil, nil)
	rec := httptest.NewRecorder()

	h.ListStateCodes(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var codes []domain.StateCode
	parseDataResponse(t, rec.Body.Bytes(), &codes)
	assert.NotNil(t, codes)
	assert.Len(t, codes, 0)
}

// ---------------------------------------------------------------------------
// Tests: tenantIDFromRequest
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
// Tests: NewRouter
// ---------------------------------------------------------------------------

func TestNewRouter(t *testing.T) {
	ms := &mockStore{}
	r := NewRouter(ms)
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
