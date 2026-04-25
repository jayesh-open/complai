package domain

import (
	"time"

	"github.com/google/uuid"
)

type WorkflowDefinition struct {
	ID           uuid.UUID `json:"id"`
	TenantID     uuid.UUID `json:"tenant_id"`
	WorkflowType string    `json:"workflow_type"`
	Description  *string   `json:"description,omitempty"`
	Version      int       `json:"version"`
	Status       string    `json:"status"`
	Config       string    `json:"config"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type WorkflowInstance struct {
	ID                 uuid.UUID  `json:"id"`
	TenantID           uuid.UUID  `json:"tenant_id"`
	WorkflowType       string     `json:"workflow_type"`
	TemporalWorkflowID *string    `json:"temporal_workflow_id,omitempty"`
	TemporalRunID      *string    `json:"temporal_run_id,omitempty"`
	State              string     `json:"state"`
	Input              string     `json:"input"`
	Output             *string    `json:"output,omitempty"`
	ErrorMessage       *string    `json:"error_message,omitempty"`
	StartedAt          time.Time  `json:"started_at"`
	CompletedAt        *time.Time `json:"completed_at,omitempty"`
	TraceID            *string    `json:"trace_id,omitempty"`
}

type HumanTask struct {
	ID                 uuid.UUID  `json:"id"`
	TenantID           uuid.UUID  `json:"tenant_id"`
	WorkflowInstanceID uuid.UUID  `json:"workflow_instance_id"`
	TaskType           string     `json:"task_type"`
	Title              string     `json:"title"`
	Description        *string    `json:"description,omitempty"`
	AssignedTo         *uuid.UUID `json:"assigned_to,omitempty"`
	Status             string     `json:"status"`
	Input              string     `json:"input"`
	Output             *string    `json:"output,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	CompletedAt        *time.Time `json:"completed_at,omitempty"`
}

type StartWorkflowRequest struct {
	WorkflowType string `json:"workflow_type"`
	Input        string `json:"input"` // JSON string
}

type SignalWorkflowRequest struct {
	SignalName string `json:"signal_name"`
	Payload    string `json:"payload"` // JSON string
}

type CompleteTaskRequest struct {
	Output string `json:"output"` // JSON string
}

type CreateDefinitionRequest struct {
	WorkflowType string  `json:"workflow_type"`
	Description  *string `json:"description,omitempty"`
	Config       string  `json:"config"`
}
