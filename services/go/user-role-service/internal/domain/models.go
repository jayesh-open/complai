package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description *string   `json:"description,omitempty"`
	IsSystem    bool      `json:"is_system"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Permission struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type RolePermission struct {
	ID           uuid.UUID `json:"id"`
	TenantID     uuid.UUID `json:"tenant_id"`
	RoleID       uuid.UUID `json:"role_id"`
	PermissionID uuid.UUID `json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserRole struct {
	ID         uuid.UUID  `json:"id"`
	TenantID   uuid.UUID  `json:"tenant_id"`
	UserID     uuid.UUID  `json:"user_id"`
	RoleID     uuid.UUID  `json:"role_id"`
	AssignedBy *uuid.UUID `json:"assigned_by,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

type RoleTemplate struct {
	ID          uuid.UUID          `json:"id"`
	Name        string             `json:"name"`
	DisplayName string             `json:"display_name"`
	Description *string            `json:"description,omitempty"`
	Permissions []PermissionPair   `json:"permissions"`
	CreatedAt   time.Time          `json:"created_at"`
}

type PermissionPair struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

type ApprovalWorkflow struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	ResourceType string     `json:"resource_type"`
	ActionType   string     `json:"action_type"`
	Status       string     `json:"status"`
	RequestedBy  uuid.UUID  `json:"requested_by"`
	DecidedBy    *uuid.UUID `json:"decided_by,omitempty"`
	Payload      string     `json:"payload"`
	Reason       *string    `json:"reason,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	DecidedAt    *time.Time `json:"decided_at,omitempty"`
}

type CreateRoleRequest struct {
	Name        string  `json:"name" validate:"required"`
	DisplayName string  `json:"display_name" validate:"required"`
	Description *string `json:"description"`
}

type AssignPermissionsRequest struct {
	PermissionIDs []uuid.UUID `json:"permission_ids" validate:"required,min=1"`
}

type AssignRoleRequest struct {
	RoleID uuid.UUID `json:"role_id" validate:"required"`
}

type PolicyCheckRequest struct {
	UserID   uuid.UUID `json:"user_id" validate:"required"`
	Resource string    `json:"resource" validate:"required"`
	Action   string    `json:"action" validate:"required"`
}

type PolicyCheckResponse struct {
	Allow   bool     `json:"allow"`
	Reasons []string `json:"reasons"`
}

type CreateApprovalRequest struct {
	ResourceType string `json:"resource_type" validate:"required"`
	ActionType   string `json:"action_type" validate:"required"`
	Payload      string `json:"payload"`
}

type DecideApprovalRequest struct {
	Decision string  `json:"decision" validate:"required,oneof=approved rejected"`
	Reason   *string `json:"reason"`
}
