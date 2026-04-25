package api

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/audit-service/internal/domain"
)

// ---------------------------------------------------------------------------
// Mock store
// ---------------------------------------------------------------------------

type mockStore struct {
	createEventFn         func(ctx context.Context, tenantID uuid.UUID, e *domain.AuditEvent) error
	listEventsFn          func(ctx context.Context, tenantID uuid.UUID, params domain.QueryParams) ([]domain.AuditEvent, error)
	getEventFn            func(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AuditEvent, error)
	getEventsForHourFn    func(ctx context.Context, tenantID uuid.UUID, hourBucket time.Time) ([]domain.AuditEvent, error)
	createMerkleChainFn   func(ctx context.Context, tenantID uuid.UUID, m *domain.MerkleChain) error
	getMerkleChainsFn     func(ctx context.Context, tenantID uuid.UUID, from, to time.Time) ([]domain.MerkleChain, error)
	getLatestMerkleChainFn func(ctx context.Context, tenantID uuid.UUID) (*domain.MerkleChain, error)
	updateEventNewValueFn func(ctx context.Context, tenantID uuid.UUID, eventID uuid.UUID, newValue string) error
}

func (m *mockStore) CreateEvent(ctx context.Context, tenantID uuid.UUID, e *domain.AuditEvent) error {
	if m.createEventFn != nil {
		return m.createEventFn(ctx, tenantID, e)
	}
	e.ID = uuid.New()
	e.TenantID = tenantID
	e.CreatedAt = time.Now()
	if e.Status == "" {
		e.Status = "success"
	}
	return nil
}

func (m *mockStore) ListEvents(ctx context.Context, tenantID uuid.UUID, params domain.QueryParams) ([]domain.AuditEvent, error) {
	if m.listEventsFn != nil {
		return m.listEventsFn(ctx, tenantID, params)
	}
	return nil, nil
}

func (m *mockStore) GetEvent(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AuditEvent, error) {
	if m.getEventFn != nil {
		return m.getEventFn(ctx, tenantID, id)
	}
	return nil, errors.New("not found")
}

func (m *mockStore) GetEventsForHour(ctx context.Context, tenantID uuid.UUID, hourBucket time.Time) ([]domain.AuditEvent, error) {
	if m.getEventsForHourFn != nil {
		return m.getEventsForHourFn(ctx, tenantID, hourBucket)
	}
	return nil, nil
}

func (m *mockStore) CreateMerkleChain(ctx context.Context, tenantID uuid.UUID, mc *domain.MerkleChain) error {
	if m.createMerkleChainFn != nil {
		return m.createMerkleChainFn(ctx, tenantID, mc)
	}
	mc.ID = uuid.New()
	mc.TenantID = tenantID
	mc.CreatedAt = time.Now()
	return nil
}

func (m *mockStore) GetMerkleChains(ctx context.Context, tenantID uuid.UUID, from, to time.Time) ([]domain.MerkleChain, error) {
	if m.getMerkleChainsFn != nil {
		return m.getMerkleChainsFn(ctx, tenantID, from, to)
	}
	return nil, nil
}

func (m *mockStore) GetLatestMerkleChain(ctx context.Context, tenantID uuid.UUID) (*domain.MerkleChain, error) {
	if m.getLatestMerkleChainFn != nil {
		return m.getLatestMerkleChainFn(ctx, tenantID)
	}
	return nil, errors.New("no chain")
}

func (m *mockStore) UpdateEventNewValue(ctx context.Context, tenantID uuid.UUID, eventID uuid.UUID, newValue string) error {
	if m.updateEventNewValueFn != nil {
		return m.updateEventNewValueFn(ctx, tenantID, eventID, newValue)
	}
	return nil
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

func testComputeHash(previousHash string, events []domain.AuditEvent) string {
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
	assert.Equal(t, "audit-service", data["service"])
}

// ---------------------------------------------------------------------------
// Tests: CreateEvent
// ---------------------------------------------------------------------------

func TestCreateEvent_Success(t *testing.T) {
	ms := &mockStore{}
	h := NewHandlers(ms)

	tenantID := uuid.New()
	body := domain.CreateAuditEventRequest{
		ResourceType: "tenant",
		Action:       "create",
		Status:       "success",
	}
	req := newRequest(t, http.MethodPost, "/v1/audit/events", body, nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.CreateEvent(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var event domain.AuditEvent
	parseDataResponse(t, rec.Body.Bytes(), &event)
	assert.Equal(t, tenantID, event.TenantID)
	assert.Equal(t, "tenant", event.ResourceType)
	assert.Equal(t, "create", event.Action)
	assert.Equal(t, "success", event.Status)
	assert.NotEqual(t, uuid.Nil, event.ID)
}

func TestCreateEvent_DefaultStatus(t *testing.T) {
	ms := &mockStore{}
	h := NewHandlers(ms)

	tenantID := uuid.New()
	body := domain.CreateAuditEventRequest{
		ResourceType: "tenant",
		Action:       "create",
	}
	req := newRequest(t, http.MethodPost, "/v1/audit/events", body, nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.CreateEvent(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var event domain.AuditEvent
	parseDataResponse(t, rec.Body.Bytes(), &event)
	assert.Equal(t, "success", event.Status)
}

func TestCreateEvent_InvalidBody(t *testing.T) {
	h := NewHandlers(&mockStore{})

	tenantID := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/v1/audit/events", bytes.NewReader([]byte("not json")))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.CreateEvent(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid request", data["error"])
}

func TestCreateEvent_MissingTenantID(t *testing.T) {
	h := NewHandlers(&mockStore{})

	body := domain.CreateAuditEventRequest{ResourceType: "tenant", Action: "create"}
	req := newRequest(t, http.MethodPost, "/v1/audit/events", body, nil)
	// No X-Tenant-Id header
	rec := httptest.NewRecorder()

	h.CreateEvent(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Contains(t, data["error"], "missing X-Tenant-Id")
}

func TestCreateEvent_StoreError(t *testing.T) {
	ms := &mockStore{
		createEventFn: func(_ context.Context, _ uuid.UUID, _ *domain.AuditEvent) error {
			return errors.New("db down")
		},
	}
	h := NewHandlers(ms)

	tenantID := uuid.New()
	body := domain.CreateAuditEventRequest{ResourceType: "tenant", Action: "create"}
	req := newRequest(t, http.MethodPost, "/v1/audit/events", body, nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.CreateEvent(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "create failed", data["error"])
}

// ---------------------------------------------------------------------------
// Tests: ListEvents
// ---------------------------------------------------------------------------

func TestListEvents_Success(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listEventsFn: func(_ context.Context, id uuid.UUID, params domain.QueryParams) ([]domain.AuditEvent, error) {
			assert.Equal(t, tenantID, id)
			return []domain.AuditEvent{
				{ID: uuid.New(), TenantID: tenantID, ResourceType: "tenant", Action: "create", Status: "success", CreatedAt: time.Now()},
			}, nil
		},
	}
	h := NewHandlers(ms)

	req := httptest.NewRequest(http.MethodGet, "/v1/audit/events", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListEvents(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var events []domain.AuditEvent
	parseDataResponse(t, rec.Body.Bytes(), &events)
	require.Len(t, events, 1)
	assert.Equal(t, "tenant", events[0].ResourceType)
}

func TestListEvents_WithQueryParams(t *testing.T) {
	tenantID := uuid.New()
	var capturedParams domain.QueryParams
	ms := &mockStore{
		listEventsFn: func(_ context.Context, _ uuid.UUID, params domain.QueryParams) ([]domain.AuditEvent, error) {
			capturedParams = params
			return []domain.AuditEvent{}, nil
		},
	}
	h := NewHandlers(ms)

	req := httptest.NewRequest(http.MethodGet, "/v1/audit/events?resource_type=tenant&action=create&limit=10&offset=5", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListEvents(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "tenant", capturedParams.ResourceType)
	assert.Equal(t, "create", capturedParams.Action)
	assert.Equal(t, 10, capturedParams.Limit)
	assert.Equal(t, 5, capturedParams.Offset)
}

func TestListEvents_MissingTenantID(t *testing.T) {
	h := NewHandlers(&mockStore{})

	req := httptest.NewRequest(http.MethodGet, "/v1/audit/events", nil)
	rec := httptest.NewRecorder()

	h.ListEvents(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListEvents_InvalidDateFrom(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	req := httptest.NewRequest(http.MethodGet, "/v1/audit/events?date_from=bad-date", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListEvents(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid date_from", data["error"])
}

func TestListEvents_InvalidDateTo(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	req := httptest.NewRequest(http.MethodGet, "/v1/audit/events?date_to=bad-date", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListEvents(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid date_to", data["error"])
}

func TestListEvents_InvalidLimit(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	req := httptest.NewRequest(http.MethodGet, "/v1/audit/events?limit=abc", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListEvents(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid limit", data["error"])
}

func TestListEvents_InvalidOffset(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	req := httptest.NewRequest(http.MethodGet, "/v1/audit/events?offset=abc", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListEvents(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid offset", data["error"])
}

func TestListEvents_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listEventsFn: func(_ context.Context, _ uuid.UUID, _ domain.QueryParams) ([]domain.AuditEvent, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewHandlers(ms)

	req := httptest.NewRequest(http.MethodGet, "/v1/audit/events", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListEvents(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "internal error", data["error"])
}

func TestListEvents_NilResultReturnsEmptyArray(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listEventsFn: func(_ context.Context, _ uuid.UUID, _ domain.QueryParams) ([]domain.AuditEvent, error) {
			return nil, nil
		},
	}
	h := NewHandlers(ms)

	req := httptest.NewRequest(http.MethodGet, "/v1/audit/events", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListEvents(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var events []domain.AuditEvent
	parseDataResponse(t, rec.Body.Bytes(), &events)
	assert.NotNil(t, events)
	assert.Len(t, events, 0)
}

// ---------------------------------------------------------------------------
// Tests: GetEvent
// ---------------------------------------------------------------------------

func TestGetEvent_Success(t *testing.T) {
	tenantID := uuid.New()
	eventID := uuid.New()
	ms := &mockStore{
		getEventFn: func(_ context.Context, tid uuid.UUID, id uuid.UUID) (*domain.AuditEvent, error) {
			assert.Equal(t, tenantID, tid)
			assert.Equal(t, eventID, id)
			return &domain.AuditEvent{
				ID:           eventID,
				TenantID:     tenantID,
				ResourceType: "tenant",
				Action:       "create",
				Status:       "success",
				CreatedAt:    time.Now(),
			}, nil
		},
	}
	h := NewHandlers(ms)

	req := newRequest(t, http.MethodGet, "/v1/audit/events/"+eventID.String(), nil, map[string]string{
		"eventID": eventID.String(),
	})
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.GetEvent(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var event domain.AuditEvent
	parseDataResponse(t, rec.Body.Bytes(), &event)
	assert.Equal(t, eventID, event.ID)
	assert.Equal(t, "tenant", event.ResourceType)
}

func TestGetEvent_InvalidID(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	req := newRequest(t, http.MethodGet, "/v1/audit/events/not-a-uuid", nil, map[string]string{
		"eventID": "not-a-uuid",
	})
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.GetEvent(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid event_id", data["error"])
}

func TestGetEvent_NotFound(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		getEventFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.AuditEvent, error) {
			return nil, errors.New("not found")
		},
	}
	h := NewHandlers(ms)

	eventID := uuid.New()
	req := newRequest(t, http.MethodGet, "/v1/audit/events/"+eventID.String(), nil, map[string]string{
		"eventID": eventID.String(),
	})
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.GetEvent(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetEvent_MissingTenantID(t *testing.T) {
	h := NewHandlers(&mockStore{})

	eventID := uuid.New()
	req := newRequest(t, http.MethodGet, "/v1/audit/events/"+eventID.String(), nil, map[string]string{
		"eventID": eventID.String(),
	})
	// No X-Tenant-Id header
	rec := httptest.NewRecorder()

	h.GetEvent(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: ComputeMerkleHash
// ---------------------------------------------------------------------------

func TestComputeMerkleHash_Success(t *testing.T) {
	tenantID := uuid.New()
	hourBucket := time.Date(2026, 4, 25, 14, 0, 0, 0, time.UTC)

	event1 := domain.AuditEvent{
		ID: uuid.New(), TenantID: tenantID, ResourceType: "tenant", Action: "create",
		Status: "success", CreatedAt: hourBucket.Add(5 * time.Minute),
	}
	event2 := domain.AuditEvent{
		ID: uuid.New(), TenantID: tenantID, ResourceType: "user", Action: "update",
		Status: "success", CreatedAt: hourBucket.Add(30 * time.Minute),
	}

	var capturedChain *domain.MerkleChain
	ms := &mockStore{
		getEventsForHourFn: func(_ context.Context, _ uuid.UUID, _ time.Time) ([]domain.AuditEvent, error) {
			return []domain.AuditEvent{event1, event2}, nil
		},
		getLatestMerkleChainFn: func(_ context.Context, _ uuid.UUID) (*domain.MerkleChain, error) {
			return nil, errors.New("no chain")
		},
		createMerkleChainFn: func(_ context.Context, tid uuid.UUID, mc *domain.MerkleChain) error {
			mc.ID = uuid.New()
			mc.TenantID = tid
			mc.CreatedAt = time.Now()
			capturedChain = mc
			return nil
		},
	}
	h := NewHandlers(ms)

	body := map[string]string{"hour_bucket": "2026-04-25T14:00:00Z"}
	req := newRequest(t, http.MethodPost, "/v1/audit/merkle/compute", body, nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ComputeMerkleHash(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	require.NotNil(t, capturedChain)
	assert.Equal(t, 2, capturedChain.EventCount)
	assert.Equal(t, hourBucket, capturedChain.HourBucket)
	assert.Equal(t, "", capturedChain.PreviousHash)

	// Verify the hash is correct
	expectedHash := testComputeHash("", []domain.AuditEvent{event1, event2})
	assert.Equal(t, expectedHash, capturedChain.ComputedHash)
}

func TestComputeMerkleHash_WithPreviousChain(t *testing.T) {
	tenantID := uuid.New()
	hourBucket := time.Date(2026, 4, 25, 15, 0, 0, 0, time.UTC)
	previousHash := "abc123def456"

	event1 := domain.AuditEvent{
		ID: uuid.New(), TenantID: tenantID, ResourceType: "tenant", Action: "create",
		Status: "success", CreatedAt: hourBucket.Add(10 * time.Minute),
	}

	ms := &mockStore{
		getEventsForHourFn: func(_ context.Context, _ uuid.UUID, _ time.Time) ([]domain.AuditEvent, error) {
			return []domain.AuditEvent{event1}, nil
		},
		getLatestMerkleChainFn: func(_ context.Context, _ uuid.UUID) (*domain.MerkleChain, error) {
			return &domain.MerkleChain{ComputedHash: previousHash}, nil
		},
		createMerkleChainFn: func(_ context.Context, tid uuid.UUID, mc *domain.MerkleChain) error {
			mc.ID = uuid.New()
			mc.TenantID = tid
			mc.CreatedAt = time.Now()
			return nil
		},
	}
	h := NewHandlers(ms)

	body := map[string]string{"hour_bucket": "2026-04-25T15:00:00Z"}
	req := newRequest(t, http.MethodPost, "/v1/audit/merkle/compute", body, nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ComputeMerkleHash(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var chain domain.MerkleChain
	parseDataResponse(t, rec.Body.Bytes(), &chain)
	assert.Equal(t, previousHash, chain.PreviousHash)

	expectedHash := testComputeHash(previousHash, []domain.AuditEvent{event1})
	assert.Equal(t, expectedHash, chain.ComputedHash)
}

func TestComputeMerkleHash_NoEvents(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		getEventsForHourFn: func(_ context.Context, _ uuid.UUID, _ time.Time) ([]domain.AuditEvent, error) {
			return []domain.AuditEvent{}, nil
		},
	}
	h := NewHandlers(ms)

	body := map[string]string{"hour_bucket": "2026-04-25T14:00:00Z"}
	req := newRequest(t, http.MethodPost, "/v1/audit/merkle/compute", body, nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ComputeMerkleHash(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "no events for this hour", data["message"])
}

func TestComputeMerkleHash_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	req := httptest.NewRequest(http.MethodPost, "/v1/audit/merkle/compute", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ComputeMerkleHash(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestComputeMerkleHash_MissingTenantID(t *testing.T) {
	h := NewHandlers(&mockStore{})

	body := map[string]string{"hour_bucket": "2026-04-25T14:00:00Z"}
	req := newRequest(t, http.MethodPost, "/v1/audit/merkle/compute", body, nil)
	rec := httptest.NewRecorder()

	h.ComputeMerkleHash(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestComputeMerkleHash_GetEventsError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		getEventsForHourFn: func(_ context.Context, _ uuid.UUID, _ time.Time) ([]domain.AuditEvent, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewHandlers(ms)

	body := map[string]string{"hour_bucket": "2026-04-25T14:00:00Z"}
	req := newRequest(t, http.MethodPost, "/v1/audit/merkle/compute", body, nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ComputeMerkleHash(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestComputeMerkleHash_CreateChainError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		getEventsForHourFn: func(_ context.Context, _ uuid.UUID, _ time.Time) ([]domain.AuditEvent, error) {
			return []domain.AuditEvent{
				{ID: uuid.New(), TenantID: tenantID, ResourceType: "t", Action: "c", Status: "success", CreatedAt: time.Now()},
			}, nil
		},
		getLatestMerkleChainFn: func(_ context.Context, _ uuid.UUID) (*domain.MerkleChain, error) {
			return nil, errors.New("no chain")
		},
		createMerkleChainFn: func(_ context.Context, _ uuid.UUID, _ *domain.MerkleChain) error {
			return errors.New("db error")
		},
	}
	h := NewHandlers(ms)

	body := map[string]string{"hour_bucket": "2026-04-25T14:00:00Z"}
	req := newRequest(t, http.MethodPost, "/v1/audit/merkle/compute", body, nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ComputeMerkleHash(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: IntegrityCheck
// ---------------------------------------------------------------------------

func TestIntegrityCheck_ValidChain(t *testing.T) {
	tenantID := uuid.New()
	hour1 := time.Date(2026, 4, 25, 14, 0, 0, 0, time.UTC)
	hour2 := time.Date(2026, 4, 25, 15, 0, 0, 0, time.UTC)

	events1 := []domain.AuditEvent{
		{ID: uuid.New(), TenantID: tenantID, ResourceType: "tenant", Action: "create", Status: "success", CreatedAt: hour1.Add(5 * time.Minute)},
	}
	events2 := []domain.AuditEvent{
		{ID: uuid.New(), TenantID: tenantID, ResourceType: "user", Action: "update", Status: "success", CreatedAt: hour2.Add(10 * time.Minute)},
	}

	hash1 := testComputeHash("", events1)
	hash2 := testComputeHash(hash1, events2)

	chains := []domain.MerkleChain{
		{ID: uuid.New(), TenantID: tenantID, HourBucket: hour1, EventCount: 1, PreviousHash: "", ComputedHash: hash1},
		{ID: uuid.New(), TenantID: tenantID, HourBucket: hour2, EventCount: 1, PreviousHash: hash1, ComputedHash: hash2},
	}

	ms := &mockStore{
		getMerkleChainsFn: func(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]domain.MerkleChain, error) {
			return chains, nil
		},
		getEventsForHourFn: func(_ context.Context, _ uuid.UUID, hb time.Time) ([]domain.AuditEvent, error) {
			if hb.Equal(hour1) {
				return events1, nil
			}
			return events2, nil
		},
	}
	h := NewHandlers(ms)

	req := httptest.NewRequest(http.MethodGet, "/v1/audit/integrity?date_from=2026-04-25T14:00:00Z&date_to=2026-04-25T16:00:00Z", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.IntegrityCheck(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result domain.IntegrityCheckResult
	parseDataResponse(t, rec.Body.Bytes(), &result)
	assert.True(t, result.Valid)
	assert.Equal(t, 2, result.ChainLength)
	assert.Equal(t, "all chain entries verified", result.Message)
	assert.Nil(t, result.BrokenAt)
}

func TestIntegrityCheck_TamperDetection(t *testing.T) {
	tenantID := uuid.New()
	hour1 := time.Date(2026, 4, 25, 14, 0, 0, 0, time.UTC)

	originalEvents := []domain.AuditEvent{
		{ID: uuid.New(), TenantID: tenantID, ResourceType: "tenant", Action: "create", Status: "success", CreatedAt: hour1.Add(5 * time.Minute)},
	}

	// Compute the original hash
	originalHash := testComputeHash("", originalEvents)

	chains := []domain.MerkleChain{
		{ID: uuid.New(), TenantID: tenantID, HourBucket: hour1, EventCount: 1, PreviousHash: "", ComputedHash: originalHash},
	}

	// Simulate tampering: modify an event's ID (which changes the hash)
	tamperedEvents := []domain.AuditEvent{
		{ID: uuid.New(), TenantID: tenantID, ResourceType: "tenant", Action: "create", Status: "success", CreatedAt: hour1.Add(5 * time.Minute)},
	}

	ms := &mockStore{
		getMerkleChainsFn: func(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]domain.MerkleChain, error) {
			return chains, nil
		},
		getEventsForHourFn: func(_ context.Context, _ uuid.UUID, _ time.Time) ([]domain.AuditEvent, error) {
			// Return tampered events (different IDs)
			return tamperedEvents, nil
		},
	}
	h := NewHandlers(ms)

	req := httptest.NewRequest(http.MethodGet, "/v1/audit/integrity?date_from=2026-04-25T14:00:00Z&date_to=2026-04-25T15:00:00Z", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.IntegrityCheck(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result domain.IntegrityCheckResult
	parseDataResponse(t, rec.Body.Bytes(), &result)
	assert.False(t, result.Valid)
	assert.NotNil(t, result.BrokenAt)
	assert.Contains(t, result.Message, "hash mismatch")
}

func TestIntegrityCheck_BrokenChainLinkage(t *testing.T) {
	tenantID := uuid.New()
	hour1 := time.Date(2026, 4, 25, 14, 0, 0, 0, time.UTC)
	hour2 := time.Date(2026, 4, 25, 15, 0, 0, 0, time.UTC)

	events1 := []domain.AuditEvent{
		{ID: uuid.New(), TenantID: tenantID, ResourceType: "tenant", Action: "create", Status: "success", CreatedAt: hour1.Add(5 * time.Minute)},
	}
	events2 := []domain.AuditEvent{
		{ID: uuid.New(), TenantID: tenantID, ResourceType: "user", Action: "update", Status: "success", CreatedAt: hour2.Add(10 * time.Minute)},
	}

	hash1 := testComputeHash("", events1)
	// Intentionally use wrong previous hash for chain linkage check
	wrongPrevHash := "wrong_hash"
	hash2 := testComputeHash(wrongPrevHash, events2)

	chains := []domain.MerkleChain{
		{ID: uuid.New(), TenantID: tenantID, HourBucket: hour1, EventCount: 1, PreviousHash: "", ComputedHash: hash1},
		{ID: uuid.New(), TenantID: tenantID, HourBucket: hour2, EventCount: 1, PreviousHash: wrongPrevHash, ComputedHash: hash2},
	}

	ms := &mockStore{
		getMerkleChainsFn: func(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]domain.MerkleChain, error) {
			return chains, nil
		},
		getEventsForHourFn: func(_ context.Context, _ uuid.UUID, hb time.Time) ([]domain.AuditEvent, error) {
			if hb.Equal(hour1) {
				return events1, nil
			}
			return events2, nil
		},
	}
	h := NewHandlers(ms)

	req := httptest.NewRequest(http.MethodGet, "/v1/audit/integrity?date_from=2026-04-25T14:00:00Z&date_to=2026-04-25T16:00:00Z", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.IntegrityCheck(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result domain.IntegrityCheckResult
	parseDataResponse(t, rec.Body.Bytes(), &result)
	assert.False(t, result.Valid)
	assert.NotNil(t, result.BrokenAt)
	assert.Contains(t, result.Message, "chain linkage broken")
}

func TestIntegrityCheck_NoChains(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		getMerkleChainsFn: func(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]domain.MerkleChain, error) {
			return []domain.MerkleChain{}, nil
		},
	}
	h := NewHandlers(ms)

	req := httptest.NewRequest(http.MethodGet, "/v1/audit/integrity?date_from=2026-04-25T14:00:00Z&date_to=2026-04-25T16:00:00Z", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.IntegrityCheck(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var result domain.IntegrityCheckResult
	parseDataResponse(t, rec.Body.Bytes(), &result)
	assert.True(t, result.Valid)
	assert.Equal(t, 0, result.ChainLength)
	assert.Equal(t, "no chain entries found in range", result.Message)
}

func TestIntegrityCheck_MissingDateFrom(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	req := httptest.NewRequest(http.MethodGet, "/v1/audit/integrity?date_to=2026-04-25T16:00:00Z", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.IntegrityCheck(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "date_from and date_to are required", data["error"])
}

func TestIntegrityCheck_MissingDateTo(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	req := httptest.NewRequest(http.MethodGet, "/v1/audit/integrity?date_from=2026-04-25T14:00:00Z", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.IntegrityCheck(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestIntegrityCheck_InvalidDateFrom(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	req := httptest.NewRequest(http.MethodGet, "/v1/audit/integrity?date_from=bad&date_to=2026-04-25T16:00:00Z", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.IntegrityCheck(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid date_from", data["error"])
}

func TestIntegrityCheck_InvalidDateTo(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{})

	req := httptest.NewRequest(http.MethodGet, "/v1/audit/integrity?date_from=2026-04-25T14:00:00Z&date_to=bad", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.IntegrityCheck(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid date_to", data["error"])
}

func TestIntegrityCheck_MissingTenantID(t *testing.T) {
	h := NewHandlers(&mockStore{})

	req := httptest.NewRequest(http.MethodGet, "/v1/audit/integrity?date_from=2026-04-25T14:00:00Z&date_to=2026-04-25T16:00:00Z", nil)
	rec := httptest.NewRecorder()

	h.IntegrityCheck(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestIntegrityCheck_GetChainsError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		getMerkleChainsFn: func(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]domain.MerkleChain, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewHandlers(ms)

	req := httptest.NewRequest(http.MethodGet, "/v1/audit/integrity?date_from=2026-04-25T14:00:00Z&date_to=2026-04-25T16:00:00Z", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.IntegrityCheck(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestIntegrityCheck_GetEventsError(t *testing.T) {
	tenantID := uuid.New()
	hour1 := time.Date(2026, 4, 25, 14, 0, 0, 0, time.UTC)

	chains := []domain.MerkleChain{
		{ID: uuid.New(), TenantID: tenantID, HourBucket: hour1, EventCount: 1, PreviousHash: "", ComputedHash: "abc"},
	}

	ms := &mockStore{
		getMerkleChainsFn: func(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]domain.MerkleChain, error) {
			return chains, nil
		},
		getEventsForHourFn: func(_ context.Context, _ uuid.UUID, _ time.Time) ([]domain.AuditEvent, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewHandlers(ms)

	req := httptest.NewRequest(http.MethodGet, "/v1/audit/integrity?date_from=2026-04-25T14:00:00Z&date_to=2026-04-25T16:00:00Z", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.IntegrityCheck(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
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

// ---------------------------------------------------------------------------
// Tests: computeHash / buildHashPayload
// ---------------------------------------------------------------------------

func TestComputeHash_Deterministic(t *testing.T) {
	events := []domain.AuditEvent{
		{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111"), CreatedAt: time.Date(2026, 4, 25, 14, 5, 0, 0, time.UTC)},
		{ID: uuid.MustParse("22222222-2222-2222-2222-222222222222"), CreatedAt: time.Date(2026, 4, 25, 14, 10, 0, 0, time.UTC)},
	}

	hash1 := computeHash("prev", events)
	hash2 := computeHash("prev", events)
	assert.Equal(t, hash1, hash2)

	// Different previous hash should produce different result
	hash3 := computeHash("other", events)
	assert.NotEqual(t, hash1, hash3)
}

func TestComputeHash_OrderIndependent(t *testing.T) {
	e1 := domain.AuditEvent{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111"), CreatedAt: time.Date(2026, 4, 25, 14, 5, 0, 0, time.UTC)}
	e2 := domain.AuditEvent{ID: uuid.MustParse("22222222-2222-2222-2222-222222222222"), CreatedAt: time.Date(2026, 4, 25, 14, 10, 0, 0, time.UTC)}

	// Pass events in different order - should produce same hash due to sorting
	hash1 := computeHash("", []domain.AuditEvent{e1, e2})
	hash2 := computeHash("", []domain.AuditEvent{e2, e1})
	assert.Equal(t, hash1, hash2)
}

func TestBuildHashPayload(t *testing.T) {
	events := []domain.AuditEvent{
		{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111"), CreatedAt: time.Date(2026, 4, 25, 14, 5, 0, 0, time.UTC)},
	}

	payload := buildHashPayload(events)
	assert.Contains(t, payload, "11111111-1111-1111-1111-111111111111")
	assert.Contains(t, payload, "2026-04-25T14:05:00Z")
}
