package store

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/ewb-service/internal/domain"
)

type Repository interface {
	CreateEWB(ctx context.Context, tenantID uuid.UUID, ewb *domain.EWayBill) error
	GetEWB(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.EWayBill, error)
	GetEWBByNumber(ctx context.Context, tenantID uuid.UUID, ewbNumber string) (*domain.EWayBill, error)
	ListEWBs(ctx context.Context, tenantID uuid.UUID, req *domain.ListEWBRequest) ([]domain.EWayBill, int, error)
	UpdateEWBGenerated(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, ewbNumber string, validFrom, validUntil time.Time) error
	UpdateEWBCancelled(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, reason string) error
	UpdateEWBStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status domain.EWBStatus) error
	UpdateEWBVehicle(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, vehicleNumber string) error
	UpdateEWBValidity(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, validUntil time.Time) error
	SetConsolidatedEWBID(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, consolidatedID uuid.UUID) error

	CreateItems(ctx context.Context, tenantID uuid.UUID, items []domain.EWBItem) error
	GetItems(ctx context.Context, tenantID uuid.UUID, ewbID uuid.UUID) ([]domain.EWBItem, error)

	CreateVehicleUpdate(ctx context.Context, tenantID uuid.UUID, update *domain.VehicleUpdate) error
	GetVehicleUpdates(ctx context.Context, tenantID uuid.UUID, ewbID uuid.UUID) ([]domain.VehicleUpdate, error)

	CreateConsolidation(ctx context.Context, tenantID uuid.UUID, consolidation *domain.Consolidation) error
}

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now() }

func CancellationWindowOpen(generatedAt *time.Time, clock Clock) bool {
	if generatedAt == nil {
		return false
	}
	return clock.Now().Sub(*generatedAt) < 24*time.Hour
}

func ValidityDays(distanceKm int, isODC bool) int {
	if distanceKm <= 0 {
		return 1
	}
	divisor := 200
	if isODC {
		divisor = 20
	}
	days := (distanceKm + divisor - 1) / divisor
	if days < 1 {
		return 1
	}
	return days
}
