package domain

import (
	"time"

	"github.com/google/uuid"
)

// HSNCode represents a Harmonized System of Nomenclature code with GST rate.
type HSNCode struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	Code          string     `json:"code"`
	Description   string     `json:"description"`
	GSTRate       float64    `json:"gst_rate"`
	EffectiveFrom string     `json:"effective_from"`
	EffectiveTo   *string    `json:"effective_to,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// StateCode represents an Indian state code for GST.
type StateCode struct {
	ID       uuid.UUID `json:"id"`
	TenantID uuid.UUID `json:"tenant_id"`
	Code     string    `json:"code"`
	Name     string    `json:"name"`
	TINCode  *string   `json:"tin_code,omitempty"`
}

// Vendor represents a supplier/vendor master record.
type Vendor struct {
	ID              uuid.UUID  `json:"id"`
	TenantID        uuid.UUID  `json:"tenant_id"`
	Name            string     `json:"name"`
	PAN             *string    `json:"pan,omitempty"`
	GSTIN           *string    `json:"gstin,omitempty"`
	Email           *string    `json:"email,omitempty"`
	Phone           *string    `json:"phone,omitempty"`
	AddressLine1    *string    `json:"address_line1,omitempty"`
	AddressLine2    *string    `json:"address_line2,omitempty"`
	City            *string    `json:"city,omitempty"`
	StateCode       *string    `json:"state_code,omitempty"`
	Pincode         *string    `json:"pincode,omitempty"`
	BankName        *string    `json:"bank_name,omitempty"`
	BankAccount     *string    `json:"bank_account,omitempty"`
	BankIFSC        *string    `json:"bank_ifsc,omitempty"`
	KYCStatus       string     `json:"kyc_status"`
	ComplianceScore float64    `json:"compliance_score"`
	Status          string     `json:"status"`
	Metadata        string     `json:"metadata"`
	CreatedBy       *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Customer represents a customer master record.
type Customer struct {
	ID               uuid.UUID `json:"id"`
	TenantID         uuid.UUID `json:"tenant_id"`
	Name             string    `json:"name"`
	PAN              *string   `json:"pan,omitempty"`
	GSTIN            *string   `json:"gstin,omitempty"`
	Email            *string   `json:"email,omitempty"`
	Phone            *string   `json:"phone,omitempty"`
	AddressLine1     *string   `json:"address_line1,omitempty"`
	City             *string   `json:"city,omitempty"`
	StateCode        *string   `json:"state_code,omitempty"`
	Pincode          *string   `json:"pincode,omitempty"`
	PaymentTermsDays *int      `json:"payment_terms_days,omitempty"`
	CreditLimit      *float64  `json:"credit_limit,omitempty"`
	Status           string    `json:"status"`
	Metadata         string    `json:"metadata"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Item represents a product or service in the catalog.
type Item struct {
	ID            uuid.UUID `json:"id"`
	TenantID      uuid.UUID `json:"tenant_id"`
	Name          string    `json:"name"`
	Description   *string   `json:"description,omitempty"`
	HSNCode       string    `json:"hsn_code"`
	UnitOfMeasure string    `json:"unit_of_measure"`
	UnitPrice     *float64  `json:"unit_price,omitempty"`
	GSTRate       *float64  `json:"gst_rate,omitempty"`
	IsService     bool      `json:"is_service"`
	Status        string    `json:"status"`
	Metadata      string    `json:"metadata"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Request types

type CreateVendorRequest struct {
	Name         string  `json:"name" validate:"required"`
	PAN          *string `json:"pan"`
	GSTIN        *string `json:"gstin"`
	Email        *string `json:"email"`
	Phone        *string `json:"phone"`
	AddressLine1 *string `json:"address_line1"`
	AddressLine2 *string `json:"address_line2"`
	City         *string `json:"city"`
	StateCode    *string `json:"state_code"`
	Pincode      *string `json:"pincode"`
	BankName     *string `json:"bank_name"`
	BankAccount  *string `json:"bank_account"`
	BankIFSC     *string `json:"bank_ifsc"`
}

type UpdateVendorRequest struct {
	Name         *string `json:"name"`
	PAN          *string `json:"pan"`
	GSTIN        *string `json:"gstin"`
	Email        *string `json:"email"`
	Phone        *string `json:"phone"`
	AddressLine1 *string `json:"address_line1"`
	AddressLine2 *string `json:"address_line2"`
	City         *string `json:"city"`
	StateCode    *string `json:"state_code"`
	Pincode      *string `json:"pincode"`
	BankName     *string `json:"bank_name"`
	BankAccount  *string `json:"bank_account"`
	BankIFSC     *string `json:"bank_ifsc"`
	Status       *string `json:"status"`
}

type CreateCustomerRequest struct {
	Name             string   `json:"name" validate:"required"`
	PAN              *string  `json:"pan"`
	GSTIN            *string  `json:"gstin"`
	Email            *string  `json:"email"`
	Phone            *string  `json:"phone"`
	AddressLine1     *string  `json:"address_line1"`
	City             *string  `json:"city"`
	StateCode        *string  `json:"state_code"`
	Pincode          *string  `json:"pincode"`
	PaymentTermsDays *int     `json:"payment_terms_days"`
	CreditLimit      *float64 `json:"credit_limit"`
}

type CreateItemRequest struct {
	Name          string   `json:"name" validate:"required"`
	Description   *string  `json:"description"`
	HSNCode       string   `json:"hsn_code" validate:"required"`
	UnitOfMeasure string   `json:"unit_of_measure"`
	UnitPrice     *float64 `json:"unit_price"`
	GSTRate       *float64 `json:"gst_rate"`
	IsService     bool     `json:"is_service"`
}

type CreateHSNCodeRequest struct {
	Code          string  `json:"code" validate:"required"`
	Description   string  `json:"description" validate:"required"`
	GSTRate       float64 `json:"gst_rate" validate:"required"`
	EffectiveFrom string  `json:"effective_from"`
}
