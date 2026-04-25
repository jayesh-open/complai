package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/vendor-compliance-service/internal/domain"
)

type Repository interface {
	UpsertVendorSnapshot(ctx context.Context, tenantID uuid.UUID, v *domain.VendorSnapshot) error
	ListVendorSnapshots(ctx context.Context, tenantID uuid.UUID) ([]domain.VendorSnapshot, error)
	GetVendorSnapshot(ctx context.Context, tenantID uuid.UUID, vendorID string) (*domain.VendorSnapshot, error)

	CreateComplianceScore(ctx context.Context, tenantID uuid.UUID, s *domain.ComplianceScore) error
	GetLatestScore(ctx context.Context, tenantID uuid.UUID, vendorID string) (*domain.ComplianceScore, error)
	ListLatestScores(ctx context.Context, tenantID uuid.UUID) ([]domain.ComplianceScore, error)
	GetScoreSummary(ctx context.Context, tenantID uuid.UUID) (*domain.ScoreSummary, error)

	CreateSyncStatus(ctx context.Context, tenantID uuid.UUID, s *domain.SyncStatus) error
	UpdateSyncStatus(ctx context.Context, tenantID uuid.UUID, syncID uuid.UUID, status string, vendorCount, scoredCount int, errMsg string) error
	GetLatestSyncStatus(ctx context.Context, tenantID uuid.UUID) (*domain.SyncStatus, error)
}
