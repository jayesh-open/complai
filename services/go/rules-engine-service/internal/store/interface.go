package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/rules-engine-service/internal/domain"
)

// Repository defines the data-access contract for the rules engine service.
type Repository interface {
	CreateRule(ctx context.Context, tenantID uuid.UUID, r *domain.Rule) error
	GetRule(ctx context.Context, tenantID uuid.UUID, ruleID uuid.UUID) (*domain.Rule, error)
	ListRules(ctx context.Context, tenantID uuid.UUID, category string) ([]domain.Rule, error)
	UpdateRule(ctx context.Context, tenantID uuid.UUID, r *domain.Rule) error
	DeleteRule(ctx context.Context, tenantID uuid.UUID, ruleID uuid.UUID) error
	CreateExecutionLog(ctx context.Context, tenantID uuid.UUID, l *domain.RuleExecutionLog) error
}
