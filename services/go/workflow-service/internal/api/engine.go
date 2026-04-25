package api

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/client"

	"github.com/complai/complai/services/go/workflow-service/internal/workflows"
)

// WorkflowEngine abstracts Temporal for testability.
type WorkflowEngine interface {
	StartWorkflow(ctx context.Context, workflowType, workflowID string, input interface{}) (runID string, err error)
	SignalWorkflow(ctx context.Context, workflowID, runID, signalName string, payload interface{}) error
	GetWorkflowStatus(ctx context.Context, workflowID, runID string) (string, error)
}

// ---------------------------------------------------------------------------
// TemporalEngine — real Temporal client wrapper
// ---------------------------------------------------------------------------

type TemporalEngine struct {
	client client.Client
}

func NewTemporalEngine(c client.Client) *TemporalEngine {
	return &TemporalEngine{client: c}
}

func (e *TemporalEngine) StartWorkflow(ctx context.Context, workflowType, workflowID string, input interface{}) (string, error) {
	opts := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "complai-default",
	}

	var run client.WorkflowRun
	var err error
	switch workflowType {
	case "sample_saga":
		run, err = e.client.ExecuteWorkflow(ctx, opts, workflows.SampleSagaWorkflow, input)
	default:
		return "", fmt.Errorf("unknown workflow type: %s", workflowType)
	}
	if err != nil {
		return "", fmt.Errorf("execute workflow: %w", err)
	}
	return run.GetRunID(), nil
}

func (e *TemporalEngine) SignalWorkflow(ctx context.Context, workflowID, runID, signalName string, payload interface{}) error {
	return e.client.SignalWorkflow(ctx, workflowID, runID, signalName, payload)
}

func (e *TemporalEngine) GetWorkflowStatus(ctx context.Context, workflowID, runID string) (string, error) {
	desc, err := e.client.DescribeWorkflowExecution(ctx, workflowID, runID)
	if err != nil {
		return "", fmt.Errorf("describe workflow: %w", err)
	}
	return desc.WorkflowExecutionInfo.Status.String(), nil
}

// ---------------------------------------------------------------------------
// NoopEngine — used when Temporal is unavailable
// ---------------------------------------------------------------------------

type NoopEngine struct{}

func NewNoopEngine() *NoopEngine {
	return &NoopEngine{}
}

func (e *NoopEngine) StartWorkflow(_ context.Context, _, _ string, _ interface{}) (string, error) {
	return "", nil
}

func (e *NoopEngine) SignalWorkflow(_ context.Context, _, _, _ string, _ interface{}) error {
	return nil
}

func (e *NoopEngine) GetWorkflowStatus(_ context.Context, _, _ string) (string, error) {
	return "unknown", nil
}
