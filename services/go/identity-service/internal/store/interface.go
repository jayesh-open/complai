package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/identity-service/internal/domain"
)

// Repository defines the data-access contract used by the API handlers.
// The concrete *Store satisfies this interface; tests can supply a mock.
type Repository interface {
	CreateUser(ctx context.Context, tenantID uuid.UUID, u *domain.User) error
	GetUserByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, tenantID, userID uuid.UUID) (*domain.User, error)
	ListUsers(ctx context.Context, tenantID uuid.UUID) ([]domain.User, error)
	CreateSession(ctx context.Context, tenantID uuid.UUID, sess *domain.UserSession) error
	CreateMFAFactor(ctx context.Context, tenantID uuid.UUID, f *domain.MFAFactor) error
	GetMFAFactors(ctx context.Context, tenantID, userID uuid.UUID) ([]domain.MFAFactor, error)
	CreateStepUpEvent(ctx context.Context, tenantID uuid.UUID, evt *domain.StepUpEvent) error
	HasValidStepUp(ctx context.Context, tenantID, userID, sessionID uuid.UUID, actionClass string) (bool, error)
}
