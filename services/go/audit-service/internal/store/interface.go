package store

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/audit-service/internal/domain"
)

// Repository defines the data-access contract for the audit service.
type Repository interface {
	CreateEvent(ctx context.Context, tenantID uuid.UUID, e *domain.AuditEvent) error
	ListEvents(ctx context.Context, tenantID uuid.UUID, params domain.QueryParams) ([]domain.AuditEvent, error)
	GetEvent(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AuditEvent, error)
	GetEventsForHour(ctx context.Context, tenantID uuid.UUID, hourBucket time.Time) ([]domain.AuditEvent, error)
	CreateMerkleChain(ctx context.Context, tenantID uuid.UUID, m *domain.MerkleChain) error
	GetMerkleChains(ctx context.Context, tenantID uuid.UUID, from, to time.Time) ([]domain.MerkleChain, error)
	GetLatestMerkleChain(ctx context.Context, tenantID uuid.UUID) (*domain.MerkleChain, error)
	UpdateEventNewValue(ctx context.Context, tenantID uuid.UUID, eventID uuid.UUID, newValue string) error
}
