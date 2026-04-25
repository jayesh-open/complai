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
	"github.com/complai/complai/services/go/workflow-service/internal/domain"
)

// ---------------------------------------------------------------------------
// Mock store
// ---------------------------------------------------------------------------

type mockStore struct {
	createDefinitionFn        func(ctx context.Context, tenantID uuid.UUID, d *domain.WorkflowDefinition) error
	getDefinitionFn           func(ctx context.Context, tenantID, id uuid.UUID) (*domain.WorkflowDefinition, error)
	listDefinitionsFn         func(ctx context.Context, tenantID uuid.UUID) ([]domain.WorkflowDefinition, error)
	createInstanceFn          func(ctx context.Context, tenantID uuid.UUID, inst *domain.WorkflowInstance) error
	getInstanceFn             func(ctx context.Context, tenantID, id uuid.UUID) (*domain.WorkflowInstance, error)
	listInstancesFn           func(ctx context.Context, tenantID uuid.UUID) ([]domain.WorkflowInstance, error)
	updateInstanceStateFn     func(ctx context.Context, tenantID, id uuid.UUID, state string) error
	updateInstanceCompletedFn func(ctx context.Context, tenantID, id uuid.UUID, state string, output *string, errMsg *string) error
	updateInstanceTemporalFn  func(ctx context.Context, tenantID, id uuid.UUID, twid, trid string) error
	createHumanTaskFn         func(ctx context.Context, tenantID uuid.UUID, task *domain.HumanTask) error
	getHumanTaskFn            func(ctx context.Context, tenantID, id uuid.UUID) (*domain.HumanTask, error)
	listPendingTasksFn        func(ctx context.Context, tenantID uuid.UUID) ([]domain.HumanTask, error)
	completeHumanTaskFn       func(ctx context.Context, tenantID, id uuid.UUID, output string) error
}

func (m *mockStore) CreateDefinition(ctx context.Context, tenantID uuid.UUID, d *domain.WorkflowDefinition) error {
	if m.createDefinitionFn != nil {
		return m.createDefinitionFn(ctx, tenantID, d)
	}
	d.ID = uuid.New()
	d.TenantID = tenantID
	d.Status = "active"
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
	return nil
}

func (m *mockStore) GetDefinition(ctx context.Context, tenantID, id uuid.UUID) (*domain.WorkflowDefinition, error) {
	if m.getDefinitionFn != nil {
		return m.getDefinitionFn(ctx, tenantID, id)
	}
	return nil, errors.New("not found")
}

func (m *mockStore) ListDefinitions(ctx context.Context, tenantID uuid.UUID) ([]domain.WorkflowDefinition, error) {
	if m.listDefinitionsFn != nil {
		return m.listDefinitionsFn(ctx, tenantID)
	}
	return nil, nil
}

func (m *mockStore) CreateInstance(ctx context.Context, tenantID uuid.UUID, inst *domain.WorkflowInstance) error {
	if m.createInstanceFn != nil {
		return m.createInstanceFn(ctx, tenantID, inst)
	}
	inst.ID = uuid.New()
	inst.TenantID = tenantID
	inst.State = "running"
	inst.StartedAt = time.Now()
	return nil
}

func (m *mockStore) GetInstance(ctx context.Context, tenantID, id uuid.UUID) (*domain.WorkflowInstance, error) {
	if m.getInstanceFn != nil {
		return m.getInstanceFn(ctx, tenantID, id)
	}
	return nil, errors.New("not found")
}

func (m *mockStore) ListInstances(ctx context.Context, tenantID uuid.UUID) ([]domain.WorkflowInstance, error) {
	if m.listInstancesFn != nil {
		return m.listInstancesFn(ctx, tenantID)
	}
	return nil, nil
}

func (m *mockStore) UpdateInstanceState(ctx context.Context, tenantID, id uuid.UUID, state string) error {
	if m.updateInstanceStateFn != nil {
		return m.updateInstanceStateFn(ctx, tenantID, id, state)
	}
	return nil
}

func (m *mockStore) UpdateInstanceCompleted(ctx context.Context, tenantID, id uuid.UUID, state string, output *string, errMsg *string) error {
	if m.updateInstanceCompletedFn != nil {
		return m.updateInstanceCompletedFn(ctx, tenantID, id, state, output, errMsg)
	}
	return nil
}

func (m *mockStore) UpdateInstanceTemporalIDs(ctx context.Context, tenantID, id uuid.UUID, twid, trid string) error {
	if m.updateInstanceTemporalFn != nil {
		return m.updateInstanceTemporalFn(ctx, tenantID, id, twid, trid)
	}
	return nil
}

func (m *mockStore) CreateHumanTask(ctx context.Context, tenantID uuid.UUID, task *domain.HumanTask) error {
	if m.createHumanTaskFn != nil {
		return m.createHumanTaskFn(ctx, tenantID, task)
	}
	task.ID = uuid.New()
	task.TenantID = tenantID
	task.Status = "pending"
	task.CreatedAt = time.Now()
	return nil
}

func (m *mockStore) GetHumanTask(ctx context.Context, tenantID, id uuid.UUID) (*domain.HumanTask, error) {
	if m.getHumanTaskFn != nil {
		return m.getHumanTaskFn(ctx, tenantID, id)
	}
	return nil, errors.New("not found")
}

func (m *mockStore) ListPendingTasks(ctx context.Context, tenantID uuid.UUID) ([]domain.HumanTask, error) {
	if m.listPendingTasksFn != nil {
		return m.listPendingTasksFn(ctx, tenantID)
	}
	return nil, nil
}

func (m *mockStore) CompleteHumanTask(ctx context.Context, tenantID, id uuid.UUID, output string) error {
	if m.completeHumanTaskFn != nil {
		return m.completeHumanTaskFn(ctx, tenantID, id, output)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Mock workflow engine
// ---------------------------------------------------------------------------

type mockEngine struct {
	startWorkflowFn    func(ctx context.Context, workflowType, workflowID string, input interface{}) (string, error)
	signalWorkflowFn   func(ctx context.Context, workflowID, runID, signalName string, payload interface{}) error
	getWorkflowStatusFn func(ctx context.Context, workflowID, runID string) (string, error)
}

func (m *mockEngine) StartWorkflow(ctx context.Context, workflowType, workflowID string, input interface{}) (string, error) {
	if m.startWorkflowFn != nil {
		return m.startWorkflowFn(ctx, workflowType, workflowID, input)
	}
	return "mock-run-id", nil
}

func (m *mockEngine) SignalWorkflow(ctx context.Context, workflowID, runID, signalName string, payload interface{}) error {
	if m.signalWorkflowFn != nil {
		return m.signalWorkflowFn(ctx, workflowID, runID, signalName, payload)
	}
	return nil
}

func (m *mockEngine) GetWorkflowStatus(ctx context.Context, workflowID, runID string) (string, error) {
	if m.getWorkflowStatusFn != nil {
		return m.getWorkflowStatusFn(ctx, workflowID, runID)
	}
	return "RUNNING", nil
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

func withTenantHeader(req *http.Request, tenantID uuid.UUID) *http.Request {
	req.Header.Set("X-Tenant-Id", tenantID.String())
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
	assert.Equal(t, "workflow-service", data["service"])
}

// ---------------------------------------------------------------------------
// Tests: StartWorkflow
// ---------------------------------------------------------------------------

func TestStartWorkflow_Success(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{}
	me := &mockEngine{}
	h := NewHandlers(ms, me)

	body := domain.StartWorkflowRequest{WorkflowType: "sample_saga", Input: `{"key":"value"}`}
	req := withTenantHeader(newRequest(t, http.MethodPost, "/v1/workflows/start", body, nil), tenantID)
	rec := httptest.NewRecorder()

	h.StartWorkflow(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var inst domain.WorkflowInstance
	parseDataResponse(t, rec.Body.Bytes(), &inst)
	assert.Equal(t, "sample_saga", inst.WorkflowType)
	assert.Equal(t, "running", inst.State)
	assert.NotEqual(t, uuid.Nil, inst.ID)
	require.NotNil(t, inst.TemporalWorkflowID)
	assert.Contains(t, *inst.TemporalWorkflowID, "wf-")
	require.NotNil(t, inst.TemporalRunID)
	assert.Equal(t, "mock-run-id", *inst.TemporalRunID)
}

func TestStartWorkflow_NoEngine(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{}
	h := NewHandlers(ms, nil)

	body := domain.StartWorkflowRequest{WorkflowType: "sample_saga", Input: `{}`}
	req := withTenantHeader(newRequest(t, http.MethodPost, "/v1/workflows/start", body, nil), tenantID)
	rec := httptest.NewRecorder()

	h.StartWorkflow(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var inst domain.WorkflowInstance
	parseDataResponse(t, rec.Body.Bytes(), &inst)
	assert.Nil(t, inst.TemporalWorkflowID)
	assert.Nil(t, inst.TemporalRunID)
}

func TestStartWorkflow_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/workflows/start", bytes.NewReader([]byte("not json")))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.StartWorkflow(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid request", data["error"])
}

func TestStartWorkflow_MissingTenantHeader(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil)

	body := domain.StartWorkflowRequest{WorkflowType: "sample_saga", Input: `{}`}
	req := newRequest(t, http.MethodPost, "/v1/workflows/start", body, nil)
	rec := httptest.NewRecorder()

	h.StartWorkflow(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestStartWorkflow_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		createInstanceFn: func(_ context.Context, _ uuid.UUID, _ *domain.WorkflowInstance) error {
			return errors.New("db down")
		},
	}
	h := NewHandlers(ms, nil)

	body := domain.StartWorkflowRequest{WorkflowType: "sample_saga", Input: `{}`}
	req := withTenantHeader(newRequest(t, http.MethodPost, "/v1/workflows/start", body, nil), tenantID)
	rec := httptest.NewRecorder()

	h.StartWorkflow(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "create failed", data["error"])
}

func TestStartWorkflow_EngineError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{}
	me := &mockEngine{
		startWorkflowFn: func(_ context.Context, _, _ string, _ interface{}) (string, error) {
			return "", errors.New("temporal unavailable")
		},
	}
	h := NewHandlers(ms, me)

	body := domain.StartWorkflowRequest{WorkflowType: "sample_saga", Input: `{}`}
	req := withTenantHeader(newRequest(t, http.MethodPost, "/v1/workflows/start", body, nil), tenantID)
	rec := httptest.NewRecorder()

	h.StartWorkflow(rec, req)

	// Instance still created even if Temporal fails
	assert.Equal(t, http.StatusCreated, rec.Code)
	var inst domain.WorkflowInstance
	parseDataResponse(t, rec.Body.Bytes(), &inst)
	assert.Nil(t, inst.TemporalWorkflowID)
}

// ---------------------------------------------------------------------------
// Tests: GetWorkflow
// ---------------------------------------------------------------------------

func TestGetWorkflow_Success(t *testing.T) {
	tenantID := uuid.New()
	instanceID := uuid.New()
	ms := &mockStore{
		getInstanceFn: func(_ context.Context, tid, id uuid.UUID) (*domain.WorkflowInstance, error) {
			assert.Equal(t, tenantID, tid)
			assert.Equal(t, instanceID, id)
			return &domain.WorkflowInstance{
				ID:           instanceID,
				TenantID:     tenantID,
				WorkflowType: "sample_saga",
				State:        "running",
				Input:        `{}`,
			}, nil
		},
	}
	h := NewHandlers(ms, nil)

	req := withTenantHeader(newRequest(t, http.MethodGet, "/v1/workflows/"+instanceID.String(), nil, map[string]string{
		"instanceID": instanceID.String(),
	}), tenantID)
	rec := httptest.NewRecorder()

	h.GetWorkflow(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var inst domain.WorkflowInstance
	parseDataResponse(t, rec.Body.Bytes(), &inst)
	assert.Equal(t, instanceID, inst.ID)
	assert.Equal(t, "sample_saga", inst.WorkflowType)
}

func TestGetWorkflow_NotFound(t *testing.T) {
	tenantID := uuid.New()
	instanceID := uuid.New()
	ms := &mockStore{
		getInstanceFn: func(_ context.Context, _, _ uuid.UUID) (*domain.WorkflowInstance, error) {
			return nil, errors.New("not found")
		},
	}
	h := NewHandlers(ms, nil)

	req := withTenantHeader(newRequest(t, http.MethodGet, "/v1/workflows/"+instanceID.String(), nil, map[string]string{
		"instanceID": instanceID.String(),
	}), tenantID)
	rec := httptest.NewRecorder()

	h.GetWorkflow(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetWorkflow_InvalidID(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, nil)

	req := withTenantHeader(newRequest(t, http.MethodGet, "/v1/workflows/bad", nil, map[string]string{
		"instanceID": "bad",
	}), tenantID)
	rec := httptest.NewRecorder()

	h.GetWorkflow(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid instance_id", data["error"])
}

func TestGetWorkflow_MissingTenantHeader(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil)

	req := newRequest(t, http.MethodGet, "/v1/workflows/"+uuid.New().String(), nil, map[string]string{
		"instanceID": uuid.New().String(),
	})
	rec := httptest.NewRecorder()

	h.GetWorkflow(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: ListWorkflows
// ---------------------------------------------------------------------------

func TestListWorkflows_Success(t *testing.T) {
	tenantID := uuid.New()
	instanceID := uuid.New()
	ms := &mockStore{
		listInstancesFn: func(_ context.Context, tid uuid.UUID) ([]domain.WorkflowInstance, error) {
			assert.Equal(t, tenantID, tid)
			return []domain.WorkflowInstance{
				{ID: instanceID, TenantID: tenantID, WorkflowType: "sample_saga", State: "running", Input: `{}`},
			}, nil
		},
	}
	h := NewHandlers(ms, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/workflows", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListWorkflows(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var instances []domain.WorkflowInstance
	parseDataResponse(t, rec.Body.Bytes(), &instances)
	require.Len(t, instances, 1)
	assert.Equal(t, "sample_saga", instances[0].WorkflowType)
}

func TestListWorkflows_Empty(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listInstancesFn: func(_ context.Context, _ uuid.UUID) ([]domain.WorkflowInstance, error) {
			return nil, nil
		},
	}
	h := NewHandlers(ms, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/workflows", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListWorkflows(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var instances []domain.WorkflowInstance
	parseDataResponse(t, rec.Body.Bytes(), &instances)
	assert.NotNil(t, instances)
	assert.Len(t, instances, 0)
}

func TestListWorkflows_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listInstancesFn: func(_ context.Context, _ uuid.UUID) ([]domain.WorkflowInstance, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewHandlers(ms, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/workflows", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListWorkflows(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestListWorkflows_MissingTenantHeader(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/workflows", nil)
	rec := httptest.NewRecorder()

	h.ListWorkflows(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: SignalWorkflow
// ---------------------------------------------------------------------------

func TestSignalWorkflow_Success(t *testing.T) {
	tenantID := uuid.New()
	instanceID := uuid.New()
	twid := "wf-" + instanceID.String()
	trid := "run-123"

	var capturedSignal string
	ms := &mockStore{
		getInstanceFn: func(_ context.Context, _, id uuid.UUID) (*domain.WorkflowInstance, error) {
			return &domain.WorkflowInstance{
				ID:                  instanceID,
				TenantID:            tenantID,
				TemporalWorkflowID:  &twid,
				TemporalRunID:       &trid,
				State:               "running",
			}, nil
		},
	}
	me := &mockEngine{
		signalWorkflowFn: func(_ context.Context, wid, rid, signal string, _ interface{}) error {
			assert.Equal(t, twid, wid)
			assert.Equal(t, trid, rid)
			capturedSignal = signal
			return nil
		},
	}
	h := NewHandlers(ms, me)

	body := domain.SignalWorkflowRequest{SignalName: "human_task_completed", Payload: `{"approved":true}`}
	req := withTenantHeader(newRequest(t, http.MethodPost, "/v1/workflows/"+instanceID.String()+"/signal", body, map[string]string{
		"instanceID": instanceID.String(),
	}), tenantID)
	rec := httptest.NewRecorder()

	h.SignalWorkflow(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "human_task_completed", capturedSignal)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "signaled", data["status"])
}

func TestSignalWorkflow_NotFound(t *testing.T) {
	tenantID := uuid.New()
	instanceID := uuid.New()
	ms := &mockStore{
		getInstanceFn: func(_ context.Context, _, _ uuid.UUID) (*domain.WorkflowInstance, error) {
			return nil, errors.New("not found")
		},
	}
	h := NewHandlers(ms, nil)

	body := domain.SignalWorkflowRequest{SignalName: "test", Payload: `{}`}
	req := withTenantHeader(newRequest(t, http.MethodPost, "/v1/workflows/"+instanceID.String()+"/signal", body, map[string]string{
		"instanceID": instanceID.String(),
	}), tenantID)
	rec := httptest.NewRecorder()

	h.SignalWorkflow(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestSignalWorkflow_InvalidID(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, nil)

	body := domain.SignalWorkflowRequest{SignalName: "test", Payload: `{}`}
	req := withTenantHeader(newRequest(t, http.MethodPost, "/v1/workflows/bad/signal", body, map[string]string{
		"instanceID": "bad",
	}), tenantID)
	rec := httptest.NewRecorder()

	h.SignalWorkflow(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSignalWorkflow_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	instanceID := uuid.New()
	h := NewHandlers(&mockStore{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/workflows/"+instanceID.String()+"/signal", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.SetPathValue("instanceID", instanceID.String())
	rec := httptest.NewRecorder()

	h.SignalWorkflow(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSignalWorkflow_EngineError(t *testing.T) {
	tenantID := uuid.New()
	instanceID := uuid.New()
	twid := "wf-" + instanceID.String()
	trid := "run-123"

	ms := &mockStore{
		getInstanceFn: func(_ context.Context, _, _ uuid.UUID) (*domain.WorkflowInstance, error) {
			return &domain.WorkflowInstance{
				ID:                  instanceID,
				TenantID:            tenantID,
				TemporalWorkflowID:  &twid,
				TemporalRunID:       &trid,
				State:               "running",
			}, nil
		},
	}
	me := &mockEngine{
		signalWorkflowFn: func(_ context.Context, _, _, _ string, _ interface{}) error {
			return errors.New("temporal error")
		},
	}
	h := NewHandlers(ms, me)

	body := domain.SignalWorkflowRequest{SignalName: "test", Payload: `{}`}
	req := withTenantHeader(newRequest(t, http.MethodPost, "/v1/workflows/"+instanceID.String()+"/signal", body, map[string]string{
		"instanceID": instanceID.String(),
	}), tenantID)
	rec := httptest.NewRecorder()

	h.SignalWorkflow(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestSignalWorkflow_MissingTenantHeader(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil)

	body := domain.SignalWorkflowRequest{SignalName: "test", Payload: `{}`}
	req := newRequest(t, http.MethodPost, "/v1/workflows/"+uuid.New().String()+"/signal", body, map[string]string{
		"instanceID": uuid.New().String(),
	})
	rec := httptest.NewRecorder()

	h.SignalWorkflow(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: ListHumanTasks
// ---------------------------------------------------------------------------

func TestListHumanTasks_Success(t *testing.T) {
	tenantID := uuid.New()
	taskID := uuid.New()
	wfID := uuid.New()
	ms := &mockStore{
		listPendingTasksFn: func(_ context.Context, tid uuid.UUID) ([]domain.HumanTask, error) {
			assert.Equal(t, tenantID, tid)
			return []domain.HumanTask{
				{ID: taskID, TenantID: tenantID, WorkflowInstanceID: wfID, TaskType: "approval", Title: "Approve invoice", Status: "pending", Input: `{}`},
			}, nil
		},
	}
	h := NewHandlers(ms, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/tasks", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListHumanTasks(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var tasks []domain.HumanTask
	parseDataResponse(t, rec.Body.Bytes(), &tasks)
	require.Len(t, tasks, 1)
	assert.Equal(t, "Approve invoice", tasks[0].Title)
}

func TestListHumanTasks_Empty(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listPendingTasksFn: func(_ context.Context, _ uuid.UUID) ([]domain.HumanTask, error) {
			return nil, nil
		},
	}
	h := NewHandlers(ms, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/tasks", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListHumanTasks(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var tasks []domain.HumanTask
	parseDataResponse(t, rec.Body.Bytes(), &tasks)
	assert.NotNil(t, tasks)
	assert.Len(t, tasks, 0)
}

func TestListHumanTasks_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listPendingTasksFn: func(_ context.Context, _ uuid.UUID) ([]domain.HumanTask, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewHandlers(ms, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/tasks", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListHumanTasks(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestListHumanTasks_MissingTenantHeader(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/tasks", nil)
	rec := httptest.NewRecorder()

	h.ListHumanTasks(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: CompleteTask
// ---------------------------------------------------------------------------

func TestCompleteTask_Success(t *testing.T) {
	tenantID := uuid.New()
	taskID := uuid.New()
	wfInstanceID := uuid.New()
	twid := "wf-" + wfInstanceID.String()
	trid := "run-123"

	var signalCalled bool
	ms := &mockStore{
		getHumanTaskFn: func(_ context.Context, _, id uuid.UUID) (*domain.HumanTask, error) {
			return &domain.HumanTask{
				ID:                 taskID,
				TenantID:           tenantID,
				WorkflowInstanceID: wfInstanceID,
				TaskType:           "approval",
				Title:              "Approve",
				Status:             "pending",
				Input:              `{}`,
			}, nil
		},
		getInstanceFn: func(_ context.Context, _, id uuid.UUID) (*domain.WorkflowInstance, error) {
			return &domain.WorkflowInstance{
				ID:                  wfInstanceID,
				TenantID:            tenantID,
				TemporalWorkflowID:  &twid,
				TemporalRunID:       &trid,
				State:               "running",
			}, nil
		},
	}
	me := &mockEngine{
		signalWorkflowFn: func(_ context.Context, wid, rid, signal string, _ interface{}) error {
			assert.Equal(t, twid, wid)
			assert.Equal(t, trid, rid)
			assert.Equal(t, "human_task_completed", signal)
			signalCalled = true
			return nil
		},
	}
	h := NewHandlers(ms, me)

	body := domain.CompleteTaskRequest{Output: `{"approved":true}`}
	req := withTenantHeader(newRequest(t, http.MethodPost, "/v1/tasks/"+taskID.String()+"/complete", body, map[string]string{
		"taskID": taskID.String(),
	}), tenantID)
	rec := httptest.NewRecorder()

	h.CompleteTask(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, signalCalled)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "completed", data["status"])
}

func TestCompleteTask_TaskNotFound(t *testing.T) {
	tenantID := uuid.New()
	taskID := uuid.New()
	ms := &mockStore{
		getHumanTaskFn: func(_ context.Context, _, _ uuid.UUID) (*domain.HumanTask, error) {
			return nil, errors.New("not found")
		},
	}
	h := NewHandlers(ms, nil)

	body := domain.CompleteTaskRequest{Output: `{}`}
	req := withTenantHeader(newRequest(t, http.MethodPost, "/v1/tasks/"+taskID.String()+"/complete", body, map[string]string{
		"taskID": taskID.String(),
	}), tenantID)
	rec := httptest.NewRecorder()

	h.CompleteTask(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCompleteTask_InvalidID(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, nil)

	body := domain.CompleteTaskRequest{Output: `{}`}
	req := withTenantHeader(newRequest(t, http.MethodPost, "/v1/tasks/bad/complete", body, map[string]string{
		"taskID": "bad",
	}), tenantID)
	rec := httptest.NewRecorder()

	h.CompleteTask(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "invalid task_id", data["error"])
}

func TestCompleteTask_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	taskID := uuid.New()
	h := NewHandlers(&mockStore{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/tasks/"+taskID.String()+"/complete", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.SetPathValue("taskID", taskID.String())
	rec := httptest.NewRecorder()

	h.CompleteTask(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCompleteTask_StoreCompleteError(t *testing.T) {
	tenantID := uuid.New()
	taskID := uuid.New()
	wfInstanceID := uuid.New()

	ms := &mockStore{
		getHumanTaskFn: func(_ context.Context, _, _ uuid.UUID) (*domain.HumanTask, error) {
			return &domain.HumanTask{
				ID:                 taskID,
				TenantID:           tenantID,
				WorkflowInstanceID: wfInstanceID,
				Status:             "pending",
			}, nil
		},
		completeHumanTaskFn: func(_ context.Context, _, _ uuid.UUID, _ string) error {
			return errors.New("db error")
		},
	}
	h := NewHandlers(ms, nil)

	body := domain.CompleteTaskRequest{Output: `{}`}
	req := withTenantHeader(newRequest(t, http.MethodPost, "/v1/tasks/"+taskID.String()+"/complete", body, map[string]string{
		"taskID": taskID.String(),
	}), tenantID)
	rec := httptest.NewRecorder()

	h.CompleteTask(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestCompleteTask_MissingTenantHeader(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil)

	body := domain.CompleteTaskRequest{Output: `{}`}
	req := newRequest(t, http.MethodPost, "/v1/tasks/"+uuid.New().String()+"/complete", body, map[string]string{
		"taskID": uuid.New().String(),
	})
	rec := httptest.NewRecorder()

	h.CompleteTask(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: CreateDefinition
// ---------------------------------------------------------------------------

func TestCreateDefinition_Success(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{}
	h := NewHandlers(ms, nil)

	body := domain.CreateDefinitionRequest{WorkflowType: "sample_saga", Config: `{"retries":3}`}
	req := withTenantHeader(newRequest(t, http.MethodPost, "/v1/definitions", body, nil), tenantID)
	rec := httptest.NewRecorder()

	h.CreateDefinition(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var def domain.WorkflowDefinition
	parseDataResponse(t, rec.Body.Bytes(), &def)
	assert.Equal(t, "sample_saga", def.WorkflowType)
	assert.Equal(t, "active", def.Status)
	assert.Equal(t, 1, def.Version)
	assert.NotEqual(t, uuid.Nil, def.ID)
}

func TestCreateDefinition_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/definitions", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.CreateDefinition(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateDefinition_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		createDefinitionFn: func(_ context.Context, _ uuid.UUID, _ *domain.WorkflowDefinition) error {
			return errors.New("db error")
		},
	}
	h := NewHandlers(ms, nil)

	body := domain.CreateDefinitionRequest{WorkflowType: "sample_saga", Config: `{}`}
	req := withTenantHeader(newRequest(t, http.MethodPost, "/v1/definitions", body, nil), tenantID)
	rec := httptest.NewRecorder()

	h.CreateDefinition(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestCreateDefinition_MissingTenantHeader(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil)

	body := domain.CreateDefinitionRequest{WorkflowType: "sample_saga", Config: `{}`}
	req := newRequest(t, http.MethodPost, "/v1/definitions", body, nil)
	rec := httptest.NewRecorder()

	h.CreateDefinition(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: ListDefinitions
// ---------------------------------------------------------------------------

func TestListDefinitions_Success(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listDefinitionsFn: func(_ context.Context, tid uuid.UUID) ([]domain.WorkflowDefinition, error) {
			return []domain.WorkflowDefinition{
				{ID: uuid.New(), TenantID: tenantID, WorkflowType: "sample_saga", Version: 1, Status: "active", Config: `{}`},
			}, nil
		},
	}
	h := NewHandlers(ms, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/definitions", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListDefinitions(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var defs []domain.WorkflowDefinition
	parseDataResponse(t, rec.Body.Bytes(), &defs)
	require.Len(t, defs, 1)
	assert.Equal(t, "sample_saga", defs[0].WorkflowType)
}

func TestListDefinitions_Empty(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listDefinitionsFn: func(_ context.Context, _ uuid.UUID) ([]domain.WorkflowDefinition, error) {
			return nil, nil
		},
	}
	h := NewHandlers(ms, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/definitions", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListDefinitions(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var defs []domain.WorkflowDefinition
	parseDataResponse(t, rec.Body.Bytes(), &defs)
	assert.NotNil(t, defs)
	assert.Len(t, defs, 0)
}

func TestListDefinitions_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listDefinitionsFn: func(_ context.Context, _ uuid.UUID) ([]domain.WorkflowDefinition, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewHandlers(ms, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/definitions", nil)
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ListDefinitions(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestListDefinitions_MissingTenantHeader(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/v1/definitions", nil)
	rec := httptest.NewRecorder()

	h.ListDefinitions(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
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

// ---------------------------------------------------------------------------
// Tests: NoopEngine
// ---------------------------------------------------------------------------

func TestNoopEngine(t *testing.T) {
	e := NewNoopEngine()

	runID, err := e.StartWorkflow(context.Background(), "test", "wf-1", nil)
	require.NoError(t, err)
	assert.Equal(t, "", runID)

	err = e.SignalWorkflow(context.Background(), "wf-1", "", "signal", nil)
	require.NoError(t, err)

	status, err := e.GetWorkflowStatus(context.Background(), "wf-1", "")
	require.NoError(t, err)
	assert.Equal(t, "unknown", status)
}
