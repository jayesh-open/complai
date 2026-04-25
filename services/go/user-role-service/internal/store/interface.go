package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/user-role-service/internal/domain"
)

// Repository defines the interface for the user-role store.
// Handlers depend on this interface to allow mock-based testing.
type Repository interface {
	CreateRole(ctx context.Context, tenantID uuid.UUID, r *domain.Role) error
	GetRole(ctx context.Context, tenantID, roleID uuid.UUID) (*domain.Role, error)
	ListRoles(ctx context.Context, tenantID uuid.UUID) ([]domain.Role, error)
	CreatePermission(ctx context.Context, tenantID uuid.UUID, p *domain.Permission) error
	AssignPermissionToRole(ctx context.Context, tenantID, roleID, permissionID uuid.UUID) error
	AssignRoleToUser(ctx context.Context, tenantID, userID, roleID uuid.UUID, assignedBy *uuid.UUID) error
	GetUserPermissions(ctx context.Context, tenantID, userID uuid.UUID) ([]domain.Permission, error)
	CheckPolicy(ctx context.Context, tenantID, userID uuid.UUID, resource, action string) (*domain.PolicyCheckResponse, error)
	CreateRoleTemplate(ctx context.Context, rt *domain.RoleTemplate) error
	GetRoleTemplates(ctx context.Context) ([]domain.RoleTemplate, error)
	CreateApproval(ctx context.Context, tenantID uuid.UUID, a *domain.ApprovalWorkflow) error
	GetApproval(ctx context.Context, tenantID, approvalID uuid.UUID) (*domain.ApprovalWorkflow, error)
	DecideApproval(ctx context.Context, tenantID, approvalID, decidedBy uuid.UUID, decision string, reason *string) error
	ListPendingApprovals(ctx context.Context, tenantID uuid.UUID) ([]domain.ApprovalWorkflow, error)
}

// Compile-time check: *Store implements Repository.
var _ Repository = (*Store)(nil)
