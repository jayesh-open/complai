package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/complai/complai/services/go/workflow-service/internal/domain"
)

type Store struct {
	pool *pgxpool.Pool
}

var _ Repository = (*Store)(nil)

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) Pool() *pgxpool.Pool { return s.pool }

func setTenantID(ctx context.Context, tx pgx.Tx, tenantID uuid.UUID) error {
	_, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID.String()))
	return err
}

// ---------------------------------------------------------------------------
// Definitions
// ---------------------------------------------------------------------------

func (s *Store) CreateDefinition(ctx context.Context, tenantID uuid.UUID, d *domain.WorkflowDefinition) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	d.ID = uuid.New()
	d.TenantID = tenantID

	err = tx.QueryRow(ctx,
		`INSERT INTO workflow_definitions (id, tenant_id, workflow_type, description, version, status, config)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING created_at, updated_at`,
		d.ID, d.TenantID, d.WorkflowType, d.Description, d.Version, "active", d.Config,
	).Scan(&d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert definition: %w", err)
	}
	d.Status = "active"
	return tx.Commit(ctx)
}

func (s *Store) GetDefinition(ctx context.Context, tenantID, id uuid.UUID) (*domain.WorkflowDefinition, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var d domain.WorkflowDefinition
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, workflow_type, description, version, status, config, created_at, updated_at
		 FROM workflow_definitions WHERE id = $1`, id,
	).Scan(&d.ID, &d.TenantID, &d.WorkflowType, &d.Description, &d.Version, &d.Status, &d.Config, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get definition: %w", err)
	}
	return &d, tx.Commit(ctx)
}

func (s *Store) ListDefinitions(ctx context.Context, tenantID uuid.UUID) ([]domain.WorkflowDefinition, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, workflow_type, description, version, status, config, created_at, updated_at
		 FROM workflow_definitions ORDER BY created_at`)
	if err != nil {
		return nil, fmt.Errorf("list definitions: %w", err)
	}
	defer rows.Close()

	var defs []domain.WorkflowDefinition
	for rows.Next() {
		var d domain.WorkflowDefinition
		if err := rows.Scan(&d.ID, &d.TenantID, &d.WorkflowType, &d.Description, &d.Version, &d.Status, &d.Config, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan definition: %w", err)
		}
		defs = append(defs, d)
	}
	return defs, tx.Commit(ctx)
}

// ---------------------------------------------------------------------------
// Instances
// ---------------------------------------------------------------------------

func (s *Store) CreateInstance(ctx context.Context, tenantID uuid.UUID, inst *domain.WorkflowInstance) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	inst.ID = uuid.New()
	inst.TenantID = tenantID

	err = tx.QueryRow(ctx,
		`INSERT INTO workflow_instances (id, tenant_id, workflow_type, state, input)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING started_at`,
		inst.ID, inst.TenantID, inst.WorkflowType, "running", inst.Input,
	).Scan(&inst.StartedAt)
	if err != nil {
		return fmt.Errorf("insert instance: %w", err)
	}
	inst.State = "running"
	return tx.Commit(ctx)
}

func (s *Store) GetInstance(ctx context.Context, tenantID, id uuid.UUID) (*domain.WorkflowInstance, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var inst domain.WorkflowInstance
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, workflow_type, temporal_workflow_id, temporal_run_id,
		        state, input, output, error_message, started_at, completed_at, trace_id
		 FROM workflow_instances WHERE id = $1`, id,
	).Scan(&inst.ID, &inst.TenantID, &inst.WorkflowType, &inst.TemporalWorkflowID, &inst.TemporalRunID,
		&inst.State, &inst.Input, &inst.Output, &inst.ErrorMessage, &inst.StartedAt, &inst.CompletedAt, &inst.TraceID)
	if err != nil {
		return nil, fmt.Errorf("get instance: %w", err)
	}
	return &inst, tx.Commit(ctx)
}

func (s *Store) ListInstances(ctx context.Context, tenantID uuid.UUID) ([]domain.WorkflowInstance, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, workflow_type, temporal_workflow_id, temporal_run_id,
		        state, input, output, error_message, started_at, completed_at, trace_id
		 FROM workflow_instances ORDER BY started_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("list instances: %w", err)
	}
	defer rows.Close()

	var instances []domain.WorkflowInstance
	for rows.Next() {
		var inst domain.WorkflowInstance
		if err := rows.Scan(&inst.ID, &inst.TenantID, &inst.WorkflowType, &inst.TemporalWorkflowID, &inst.TemporalRunID,
			&inst.State, &inst.Input, &inst.Output, &inst.ErrorMessage, &inst.StartedAt, &inst.CompletedAt, &inst.TraceID); err != nil {
			return nil, fmt.Errorf("scan instance: %w", err)
		}
		instances = append(instances, inst)
	}
	return instances, tx.Commit(ctx)
}

func (s *Store) UpdateInstanceState(ctx context.Context, tenantID, id uuid.UUID, state string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	_, err = tx.Exec(ctx,
		`UPDATE workflow_instances SET state = $1 WHERE id = $2`,
		state, id)
	if err != nil {
		return fmt.Errorf("update instance state: %w", err)
	}
	return tx.Commit(ctx)
}

func (s *Store) UpdateInstanceCompleted(ctx context.Context, tenantID, id uuid.UUID, state string, output *string, errMsg *string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	_, err = tx.Exec(ctx,
		`UPDATE workflow_instances SET state = $1, output = $2, error_message = $3, completed_at = now() WHERE id = $4`,
		state, output, errMsg, id)
	if err != nil {
		return fmt.Errorf("update instance completed: %w", err)
	}
	return tx.Commit(ctx)
}

func (s *Store) UpdateInstanceTemporalIDs(ctx context.Context, tenantID, id uuid.UUID, temporalWorkflowID, temporalRunID string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	_, err = tx.Exec(ctx,
		`UPDATE workflow_instances SET temporal_workflow_id = $1, temporal_run_id = $2 WHERE id = $3`,
		temporalWorkflowID, temporalRunID, id)
	if err != nil {
		return fmt.Errorf("update temporal IDs: %w", err)
	}
	return tx.Commit(ctx)
}

// ---------------------------------------------------------------------------
// Human Tasks
// ---------------------------------------------------------------------------

func (s *Store) CreateHumanTask(ctx context.Context, tenantID uuid.UUID, task *domain.HumanTask) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	task.ID = uuid.New()
	task.TenantID = tenantID

	err = tx.QueryRow(ctx,
		`INSERT INTO human_tasks (id, tenant_id, workflow_instance_id, task_type, title, description, assigned_to, status, input)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 RETURNING created_at`,
		task.ID, task.TenantID, task.WorkflowInstanceID, task.TaskType, task.Title, task.Description, task.AssignedTo, "pending", task.Input,
	).Scan(&task.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert human task: %w", err)
	}
	task.Status = "pending"
	return tx.Commit(ctx)
}

func (s *Store) GetHumanTask(ctx context.Context, tenantID, id uuid.UUID) (*domain.HumanTask, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var task domain.HumanTask
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, workflow_instance_id, task_type, title, description, assigned_to, status, input, output, created_at, completed_at
		 FROM human_tasks WHERE id = $1`, id,
	).Scan(&task.ID, &task.TenantID, &task.WorkflowInstanceID, &task.TaskType, &task.Title, &task.Description,
		&task.AssignedTo, &task.Status, &task.Input, &task.Output, &task.CreatedAt, &task.CompletedAt)
	if err != nil {
		return nil, fmt.Errorf("get human task: %w", err)
	}
	return &task, tx.Commit(ctx)
}

func (s *Store) ListPendingTasks(ctx context.Context, tenantID uuid.UUID) ([]domain.HumanTask, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, workflow_instance_id, task_type, title, description, assigned_to, status, input, output, created_at, completed_at
		 FROM human_tasks WHERE status = 'pending' ORDER BY created_at`)
	if err != nil {
		return nil, fmt.Errorf("list pending tasks: %w", err)
	}
	defer rows.Close()

	var tasks []domain.HumanTask
	for rows.Next() {
		var task domain.HumanTask
		if err := rows.Scan(&task.ID, &task.TenantID, &task.WorkflowInstanceID, &task.TaskType, &task.Title, &task.Description,
			&task.AssignedTo, &task.Status, &task.Input, &task.Output, &task.CreatedAt, &task.CompletedAt); err != nil {
			return nil, fmt.Errorf("scan human task: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, tx.Commit(ctx)
}

func (s *Store) CompleteHumanTask(ctx context.Context, tenantID, id uuid.UUID, output string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	_, err = tx.Exec(ctx,
		`UPDATE human_tasks SET status = 'completed', output = $1, completed_at = now() WHERE id = $2`,
		output, id)
	if err != nil {
		return fmt.Errorf("complete human task: %w", err)
	}
	return tx.Commit(ctx)
}
