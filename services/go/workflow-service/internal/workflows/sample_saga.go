package workflows

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type SagaInput struct {
	TenantID   string `json:"tenant_id"`
	WorkflowID string `json:"workflow_id"`
	TaskTitle  string `json:"task_title"`
}

type SagaOutput struct {
	Status          string `json:"status"`
	HumanTaskResult string `json:"human_task_result"`
}

func SampleSagaWorkflow(ctx workflow.Context, input SagaInput) (*SagaOutput, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy:         &temporal.RetryPolicy{MaximumAttempts: 3},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Step 1: Validate input
	var validateResult string
	err := workflow.ExecuteActivity(ctx, ValidateActivity, input).Get(ctx, &validateResult)
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	// Step 2: Create human task and wait for signal
	var taskResult string
	err = workflow.ExecuteActivity(ctx, CreateHumanTaskActivity, input).Get(ctx, &taskResult)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	// Wait for human task signal
	signalCh := workflow.GetSignalChannel(ctx, "human_task_completed")
	var humanTaskOutput string
	signalCh.Receive(ctx, &humanTaskOutput)

	// Step 3: Finalize
	var finalResult string
	err = workflow.ExecuteActivity(ctx, FinalizeActivity, humanTaskOutput).Get(ctx, &finalResult)
	if err != nil {
		return nil, fmt.Errorf("finalize: %w", err)
	}

	return &SagaOutput{
		Status:          "completed",
		HumanTaskResult: humanTaskOutput,
	}, nil
}

func ValidateActivity(ctx context.Context, input SagaInput) (string, error) {
	return "validated", nil
}

func CreateHumanTaskActivity(ctx context.Context, input SagaInput) (string, error) {
	return "task_created", nil
}

func FinalizeActivity(ctx context.Context, humanTaskOutput string) (string, error) {
	return fmt.Sprintf("finalized: %s", humanTaskOutput), nil
}
