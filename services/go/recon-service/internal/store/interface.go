package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/recon-service/internal/domain"
)

type Repository interface {
	CreateRun(ctx context.Context, tenantID uuid.UUID, run *domain.ReconRun) error
	GetRun(ctx context.Context, tenantID uuid.UUID, runID uuid.UUID) (*domain.ReconRun, error)
	UpdateRun(ctx context.Context, tenantID uuid.UUID, run *domain.ReconRun) error

	BulkInsertMatches(ctx context.Context, tenantID uuid.UUID, matches []domain.ReconMatch) error
	ListMatches(ctx context.Context, tenantID uuid.UUID, runID uuid.UUID, matchType string, status string, limit, offset int) ([]domain.ReconMatch, error)
	GetMatch(ctx context.Context, tenantID uuid.UUID, matchID uuid.UUID) (*domain.ReconMatch, error)
	UpdateMatchStatus(ctx context.Context, tenantID uuid.UUID, matchID uuid.UUID, status domain.MatchStatus, acceptedBy *uuid.UUID) error

	GetBucketSummary(ctx context.Context, tenantID uuid.UUID, runID uuid.UUID) (*domain.BucketSummary, error)

	CreateIMSAction(ctx context.Context, tenantID uuid.UUID, action *domain.IMSAction) error
	ListIMSActions(ctx context.Context, tenantID uuid.UUID, gstin, returnPeriod string) ([]domain.IMSAction, error)
}
