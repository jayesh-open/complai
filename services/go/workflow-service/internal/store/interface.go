package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/workflow-service/internal/domain"
)

// Repository defines the data-access contract for the workflow service.
type Repository interface {
	CreateDefinition(ctx context.Context, tenantID uuid.UUID, d *domain.WorkflowDefinition) error
	GetDefinition(ctx context.Context, tenantID, id uuid.UUID) (*domain.WorkflowDefinition, error)
	ListDefinitions(ctx context.Context, tenantID uuid.UUID) ([]domain.WorkflowDefinition, error)

	CreateInstance(ctx context.Context, tenantID uuid.UUID, inst *domain.WorkflowInstance) error
	GetInstance(ctx context.Context, tenantID, id uuid.UUID) (*domain.WorkflowInstance, error)
	ListInstances(ctx context.Context, tenantID uuid.UUID) ([]domain.WorkflowInstance, error)
	UpdateInstanceState(ctx context.Context, tenantID, id uuid.UUID, state string) error
	UpdateInstanceCompleted(ctx context.Context, tenantID, id uuid.UUID, state string, output *string, errMsg *string) error
	UpdateInstanceTemporalIDs(ctx context.Context, tenantID, id uuid.UUID, temporalWorkflowID, temporalRunID string) error

	CreateHumanTask(ctx context.Context, tenantID uuid.UUID, task *domain.HumanTask) error
	GetHumanTask(ctx context.Context, tenantID, id uuid.UUID) (*domain.HumanTask, error)
	ListPendingTasks(ctx context.Context, tenantID uuid.UUID) ([]domain.HumanTask, error)
	CompleteHumanTask(ctx context.Context, tenantID, id uuid.UUID, output string) error
}
