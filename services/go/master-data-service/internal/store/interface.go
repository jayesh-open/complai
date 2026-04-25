package store

import (
	"context"

	"github.com/google/uuid"

	"github.com/complai/complai/services/go/master-data-service/internal/domain"
)

// Repository defines the data-access contract for the master-data service.
type Repository interface {
	// Vendors
	CreateVendor(ctx context.Context, tenantID uuid.UUID, v *domain.Vendor) error
	GetVendor(ctx context.Context, tenantID uuid.UUID, vendorID uuid.UUID) (*domain.Vendor, error)
	ListVendors(ctx context.Context, tenantID uuid.UUID) ([]domain.Vendor, error)
	UpdateVendor(ctx context.Context, tenantID uuid.UUID, vendorID uuid.UUID, req *domain.UpdateVendorRequest) (*domain.Vendor, error)

	// Customers
	CreateCustomer(ctx context.Context, tenantID uuid.UUID, c *domain.Customer) error
	GetCustomer(ctx context.Context, tenantID uuid.UUID, customerID uuid.UUID) (*domain.Customer, error)
	ListCustomers(ctx context.Context, tenantID uuid.UUID) ([]domain.Customer, error)

	// Items
	CreateItem(ctx context.Context, tenantID uuid.UUID, i *domain.Item) error
	GetItem(ctx context.Context, tenantID uuid.UUID, itemID uuid.UUID) (*domain.Item, error)
	ListItems(ctx context.Context, tenantID uuid.UUID) ([]domain.Item, error)

	// HSN Codes
	ListHSNCodes(ctx context.Context, tenantID uuid.UUID) ([]domain.HSNCode, error)
	GetHSNCode(ctx context.Context, tenantID uuid.UUID, hsnID uuid.UUID) (*domain.HSNCode, error)
	CreateHSNCode(ctx context.Context, tenantID uuid.UUID, h *domain.HSNCode) error

	// State Codes
	ListStateCodes(ctx context.Context, tenantID uuid.UUID) ([]domain.StateCode, error)
}

// Compile-time check.
var _ Repository = (*Store)(nil)
