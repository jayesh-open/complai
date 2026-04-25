package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/document-service/internal/domain"
)

// Repository defines the data-access contract for the document service.
type Repository interface {
	CreateDocument(ctx context.Context, tenantID uuid.UUID, d *domain.Document) error
	GetDocument(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Document, error)
	ListDocuments(ctx context.Context, tenantID uuid.UUID) ([]domain.Document, error)
	UpdateOCRStatus(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, status string, result string) error
	DeleteDocument(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) error
}
