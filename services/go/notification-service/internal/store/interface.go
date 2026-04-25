package store

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/notification-service/internal/domain"
)

// Repository defines the data-access contract for the notification service.
type Repository interface {
	CreateTemplate(ctx context.Context, tenantID uuid.UUID, t *domain.NotificationTemplate) error
	GetTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.NotificationTemplate, error)
	ListTemplates(ctx context.Context, tenantID uuid.UUID) ([]domain.NotificationTemplate, error)

	GetPreferences(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (*domain.NotificationPreference, error)
	UpsertPreferences(ctx context.Context, tenantID uuid.UUID, pref *domain.NotificationPreference) error

	CreateNotification(ctx context.Context, tenantID uuid.UUID, n *domain.Notification) error
	GetNotification(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Notification, error)
	ListNotifications(ctx context.Context, tenantID uuid.UUID) ([]domain.Notification, error)

	GetPendingDigestNotifications(ctx context.Context, tenantID uuid.UUID, digestGroup string, cutoffTime time.Time) (map[uuid.UUID][]domain.Notification, error)
	MarkNotificationsSent(ctx context.Context, tenantID uuid.UUID, ids []uuid.UUID, batchID uuid.UUID) error

	CreateBounce(ctx context.Context, tenantID uuid.UUID, b *domain.NotificationBounce) error
	MarkEmailInvalid(ctx context.Context, tenantID uuid.UUID, email string) error
}
