package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/complai/complai/services/go/user-role-service/internal/domain"
)

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) Pool() *pgxpool.Pool { return s.pool }

func setTenantID(ctx context.Context, tx pgx.Tx, tenantID uuid.UUID) error {
	_, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID.String()))
	return err
}

func (s *Store) CreateRole(ctx context.Context, tenantID uuid.UUID, r *domain.Role) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO roles (tenant_id, name, display_name, description, is_system)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`,
		tenantID, r.Name, r.DisplayName, r.Description, r.IsSystem,
	).Scan(&r.ID, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert role: %w", err)
	}
	r.TenantID = tenantID
	return tx.Commit(ctx)
}

func (s *Store) GetRole(ctx context.Context, tenantID, roleID uuid.UUID) (*domain.Role, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var r domain.Role
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, name, display_name, description, is_system, created_at, updated_at
		 FROM roles WHERE id = $1`, roleID,
	).Scan(&r.ID, &r.TenantID, &r.Name, &r.DisplayName, &r.Description, &r.IsSystem, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get role: %w", err)
	}
	return &r, tx.Commit(ctx)
}

func (s *Store) ListRoles(ctx context.Context, tenantID uuid.UUID) ([]domain.Role, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, name, display_name, description, is_system, created_at, updated_at
		 FROM roles ORDER BY created_at`)
	if err != nil {
		return nil, fmt.Errorf("list roles: %w", err)
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var r domain.Role
		if err := rows.Scan(&r.ID, &r.TenantID, &r.Name, &r.DisplayName, &r.Description, &r.IsSystem, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan role: %w", err)
		}
		roles = append(roles, r)
	}
	return roles, tx.Commit(ctx)
}

func (s *Store) CreatePermission(ctx context.Context, tenantID uuid.UUID, p *domain.Permission) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO permissions (tenant_id, resource, action, description)
		 VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		tenantID, p.Resource, p.Action, p.Description,
	).Scan(&p.ID, &p.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert permission: %w", err)
	}
	p.TenantID = tenantID
	return tx.Commit(ctx)
}

func (s *Store) AssignPermissionToRole(ctx context.Context, tenantID, roleID, permissionID uuid.UUID) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO role_permissions (tenant_id, role_id, permission_id)
		 VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
		tenantID, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("assign permission: %w", err)
	}
	return tx.Commit(ctx)
}

func (s *Store) AssignRoleToUser(ctx context.Context, tenantID, userID, roleID uuid.UUID, assignedBy *uuid.UUID) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO user_roles (tenant_id, user_id, role_id, assigned_by)
		 VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`,
		tenantID, userID, roleID, assignedBy)
	if err != nil {
		return fmt.Errorf("assign role: %w", err)
	}
	return tx.Commit(ctx)
}

func (s *Store) GetUserPermissions(ctx context.Context, tenantID, userID uuid.UUID) ([]domain.Permission, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT DISTINCT p.id, p.tenant_id, p.resource, p.action, p.description, p.created_at
		 FROM permissions p
		 JOIN role_permissions rp ON p.id = rp.permission_id
		 JOIN user_roles ur ON rp.role_id = ur.role_id
		 WHERE ur.user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("get permissions: %w", err)
	}
	defer rows.Close()

	var perms []domain.Permission
	for rows.Next() {
		var p domain.Permission
		if err := rows.Scan(&p.ID, &p.TenantID, &p.Resource, &p.Action, &p.Description, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan permission: %w", err)
		}
		perms = append(perms, p)
	}
	return perms, tx.Commit(ctx)
}

func (s *Store) CheckPolicy(ctx context.Context, tenantID, userID uuid.UUID, resource, action string) (*domain.PolicyCheckResponse, error) {
	perms, err := s.GetUserPermissions(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}

	for _, p := range perms {
		if p.Resource == resource && p.Action == action {
			return &domain.PolicyCheckResponse{Allow: true, Reasons: []string{"permission granted via role"}}, nil
		}
		if p.Resource == resource && p.Action == "*" {
			return &domain.PolicyCheckResponse{Allow: true, Reasons: []string{"wildcard action on resource"}}, nil
		}
		if p.Resource == "*" && p.Action == "*" {
			return &domain.PolicyCheckResponse{Allow: true, Reasons: []string{"superadmin wildcard"}}, nil
		}
	}

	return &domain.PolicyCheckResponse{Allow: false, Reasons: []string{"no matching permission found"}}, nil
}

func (s *Store) CreateRoleTemplate(ctx context.Context, rt *domain.RoleTemplate) error {
	permJSON, err := json.Marshal(rt.Permissions)
	if err != nil {
		return fmt.Errorf("marshal permissions: %w", err)
	}

	err = s.pool.QueryRow(ctx,
		`INSERT INTO role_templates (name, display_name, description, permissions)
		 VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		rt.Name, rt.DisplayName, rt.Description, permJSON,
	).Scan(&rt.ID, &rt.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert template: %w", err)
	}
	return nil
}

func (s *Store) GetRoleTemplates(ctx context.Context) ([]domain.RoleTemplate, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, name, display_name, description, permissions, created_at
		 FROM role_templates ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list templates: %w", err)
	}
	defer rows.Close()

	var templates []domain.RoleTemplate
	for rows.Next() {
		var rt domain.RoleTemplate
		var permJSON []byte
		if err := rows.Scan(&rt.ID, &rt.Name, &rt.DisplayName, &rt.Description, &permJSON, &rt.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan template: %w", err)
		}
		if err := json.Unmarshal(permJSON, &rt.Permissions); err != nil {
			return nil, fmt.Errorf("unmarshal permissions: %w", err)
		}
		templates = append(templates, rt)
	}
	return templates, nil
}

func (s *Store) CreateApproval(ctx context.Context, tenantID uuid.UUID, a *domain.ApprovalWorkflow) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO approval_workflows (tenant_id, resource_type, action_type, status, requested_by, payload)
		 VALUES ($1, $2, $3, 'pending_approval', $4, $5) RETURNING id, created_at`,
		tenantID, a.ResourceType, a.ActionType, a.RequestedBy, a.Payload,
	).Scan(&a.ID, &a.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert approval: %w", err)
	}
	a.TenantID = tenantID
	a.Status = "pending_approval"
	return tx.Commit(ctx)
}

func (s *Store) GetApproval(ctx context.Context, tenantID, approvalID uuid.UUID) (*domain.ApprovalWorkflow, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var a domain.ApprovalWorkflow
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, resource_type, action_type, status, requested_by, decided_by, payload, reason, created_at, decided_at
		 FROM approval_workflows WHERE id = $1`, approvalID,
	).Scan(&a.ID, &a.TenantID, &a.ResourceType, &a.ActionType, &a.Status, &a.RequestedBy, &a.DecidedBy, &a.Payload, &a.Reason, &a.CreatedAt, &a.DecidedAt)
	if err != nil {
		return nil, fmt.Errorf("get approval: %w", err)
	}
	return &a, tx.Commit(ctx)
}

func (s *Store) DecideApproval(ctx context.Context, tenantID, approvalID, decidedBy uuid.UUID, decision string, reason *string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	now := time.Now()
	tag, err := tx.Exec(ctx,
		`UPDATE approval_workflows
		 SET status = $1, decided_by = $2, reason = $3, decided_at = $4
		 WHERE id = $5 AND status = 'pending_approval'`,
		decision, decidedBy, reason, now, approvalID)
	if err != nil {
		return fmt.Errorf("decide approval: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("approval not found or already decided")
	}
	return tx.Commit(ctx)
}

func (s *Store) ListPendingApprovals(ctx context.Context, tenantID uuid.UUID) ([]domain.ApprovalWorkflow, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, resource_type, action_type, status, requested_by, decided_by, payload, reason, created_at, decided_at
		 FROM approval_workflows WHERE status = 'pending_approval' ORDER BY created_at`)
	if err != nil {
		return nil, fmt.Errorf("list approvals: %w", err)
	}
	defer rows.Close()

	var approvals []domain.ApprovalWorkflow
	for rows.Next() {
		var a domain.ApprovalWorkflow
		if err := rows.Scan(&a.ID, &a.TenantID, &a.ResourceType, &a.ActionType, &a.Status, &a.RequestedBy, &a.DecidedBy, &a.Payload, &a.Reason, &a.CreatedAt, &a.DecidedAt); err != nil {
			return nil, fmt.Errorf("scan approval: %w", err)
		}
		approvals = append(approvals, a)
	}
	return approvals, tx.Commit(ctx)
}
