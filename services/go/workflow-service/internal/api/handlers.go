package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/workflow-service/internal/domain"
	"github.com/complai/complai/services/go/workflow-service/internal/store"
)

type Handlers struct {
	store  store.Repository
	engine WorkflowEngine
}

func NewHandlers(s store.Repository, engine WorkflowEngine) *Handlers {
	return &Handlers{store: s, engine: engine}
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "workflow-service"})
}

// ---------------------------------------------------------------------------
// Workflows
// ---------------------------------------------------------------------------

func (h *Handlers) StartWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.StartWorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	inst := &domain.WorkflowInstance{
		WorkflowType: req.WorkflowType,
		Input:        req.Input,
	}

	if err := h.store.CreateInstance(r.Context(), tenantID, inst); err != nil {
		log.Error().Err(err).Msg("create workflow instance failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "create failed"})
		return
	}

	if h.engine != nil {
		temporalWorkflowID := fmt.Sprintf("wf-%s", inst.ID.String())
		runID, err := h.engine.StartWorkflow(r.Context(), req.WorkflowType, temporalWorkflowID, req.Input)
		if err != nil {
			log.Error().Err(err).Msg("temporal start workflow failed")
		} else {
			twid := temporalWorkflowID
			trid := runID
			inst.TemporalWorkflowID = &twid
			inst.TemporalRunID = &trid
			if err := h.store.UpdateInstanceTemporalIDs(r.Context(), tenantID, inst.ID, temporalWorkflowID, runID); err != nil {
				log.Error().Err(err).Msg("update temporal IDs failed")
			}
		}
	}

	httputil.JSON(w, http.StatusCreated, inst)
}

func (h *Handlers) GetWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	instanceID, err := uuid.Parse(r.PathValue("instanceID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid instance_id"})
		return
	}

	inst, err := h.store.GetInstance(r.Context(), tenantID, instanceID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "workflow not found"})
		return
	}

	httputil.JSON(w, http.StatusOK, inst)
}

func (h *Handlers) ListWorkflows(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	instances, err := h.store.ListInstances(r.Context(), tenantID)
	if err != nil {
		log.Error().Err(err).Msg("list workflows failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if instances == nil {
		instances = []domain.WorkflowInstance{}
	}

	httputil.JSON(w, http.StatusOK, instances)
}

func (h *Handlers) SignalWorkflow(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	instanceID, err := uuid.Parse(r.PathValue("instanceID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid instance_id"})
		return
	}

	var req domain.SignalWorkflowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	inst, err := h.store.GetInstance(r.Context(), tenantID, instanceID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "workflow not found"})
		return
	}

	if h.engine != nil && inst.TemporalWorkflowID != nil {
		runID := ""
		if inst.TemporalRunID != nil {
			runID = *inst.TemporalRunID
		}
		if err := h.engine.SignalWorkflow(r.Context(), *inst.TemporalWorkflowID, runID, req.SignalName, req.Payload); err != nil {
			log.Error().Err(err).Msg("temporal signal failed")
			httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "signal failed"})
			return
		}
	}

	httputil.JSON(w, http.StatusOK, map[string]string{"status": "signaled"})
}

// ---------------------------------------------------------------------------
// Human Tasks
// ---------------------------------------------------------------------------

func (h *Handlers) ListHumanTasks(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	tasks, err := h.store.ListPendingTasks(r.Context(), tenantID)
	if err != nil {
		log.Error().Err(err).Msg("list tasks failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if tasks == nil {
		tasks = []domain.HumanTask{}
	}

	httputil.JSON(w, http.StatusOK, tasks)
}

func (h *Handlers) CompleteTask(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	taskID, err := uuid.Parse(r.PathValue("taskID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid task_id"})
		return
	}

	var req domain.CompleteTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	task, err := h.store.GetHumanTask(r.Context(), tenantID, taskID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}

	if err := h.store.CompleteHumanTask(r.Context(), tenantID, taskID, req.Output); err != nil {
		log.Error().Err(err).Msg("complete task failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "complete failed"})
		return
	}

	// Signal the workflow to resume
	if h.engine != nil {
		inst, err := h.store.GetInstance(r.Context(), tenantID, task.WorkflowInstanceID)
		if err == nil && inst.TemporalWorkflowID != nil {
			runID := ""
			if inst.TemporalRunID != nil {
				runID = *inst.TemporalRunID
			}
			if err := h.engine.SignalWorkflow(r.Context(), *inst.TemporalWorkflowID, runID, "human_task_completed", req.Output); err != nil {
				log.Error().Err(err).Msg("signal workflow after task completion failed")
			}
		}
	}

	httputil.JSON(w, http.StatusOK, map[string]string{"status": "completed"})
}

// ---------------------------------------------------------------------------
// Definitions
// ---------------------------------------------------------------------------

func (h *Handlers) CreateDefinition(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.CreateDefinitionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	d := &domain.WorkflowDefinition{
		WorkflowType: req.WorkflowType,
		Description:  req.Description,
		Version:      1,
		Config:       req.Config,
	}

	if err := h.store.CreateDefinition(r.Context(), tenantID, d); err != nil {
		log.Error().Err(err).Msg("create definition failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "create failed"})
		return
	}

	httputil.JSON(w, http.StatusCreated, d)
}

func (h *Handlers) ListDefinitions(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	defs, err := h.store.ListDefinitions(r.Context(), tenantID)
	if err != nil {
		log.Error().Err(err).Msg("list definitions failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if defs == nil {
		defs = []domain.WorkflowDefinition{}
	}

	httputil.JSON(w, http.StatusOK, defs)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func tenantIDFromRequest(r *http.Request) (uuid.UUID, error) {
	h := r.Header.Get("X-Tenant-Id")
	if h == "" {
		return uuid.Nil, fmt.Errorf("missing X-Tenant-Id header")
	}
	return uuid.Parse(h)
}
