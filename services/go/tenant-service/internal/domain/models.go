package domain

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID        uuid.UUID  `json:"id"`
	TenantID  uuid.UUID  `json:"tenant_id"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	Tier      string     `json:"tier"`
	Status    string     `json:"status"`
	KMSKeyARN *string    `json:"kms_key_arn,omitempty"`
	Settings  string     `json:"settings"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type TenantPAN struct {
	ID         uuid.UUID `json:"id"`
	TenantID   uuid.UUID `json:"tenant_id"`
	PAN        string    `json:"pan"`
	EntityName string    `json:"entity_name"`
	PANType    string    `json:"pan_type"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type TenantGSTIN struct {
	ID               uuid.UUID `json:"id"`
	TenantID         uuid.UUID `json:"tenant_id"`
	PANID            uuid.UUID `json:"pan_id"`
	GSTIN            string    `json:"gstin"`
	TradeName        *string   `json:"trade_name,omitempty"`
	StateCode        string    `json:"state_code"`
	RegistrationType string    `json:"registration_type"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type TenantTAN struct {
	ID            uuid.UUID `json:"id"`
	TenantID      uuid.UUID `json:"tenant_id"`
	PANID         uuid.UUID `json:"pan_id"`
	TAN           string    `json:"tan"`
	DeductorName  string    `json:"deductor_name"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateTenantRequest struct {
	Name string `json:"name" validate:"required"`
	Slug string `json:"slug" validate:"required"`
	Tier string `json:"tier" validate:"required,oneof=pooled bridge silo on_prem"`
}

type CreatePANRequest struct {
	PAN        string `json:"pan" validate:"required,len=10"`
	EntityName string `json:"entity_name" validate:"required"`
	PANType    string `json:"pan_type" validate:"required,oneof=company firm individual trust aop huf"`
}

type CreateGSTINRequest struct {
	GSTIN            string  `json:"gstin" validate:"required,len=15"`
	TradeName        *string `json:"trade_name"`
	StateCode        string  `json:"state_code" validate:"required,len=2"`
	RegistrationType string  `json:"registration_type" validate:"required,oneof=regular composition isd nrtp sez casual"`
}

type CreateTANRequest struct {
	TAN          string `json:"tan" validate:"required,len=10"`
	DeductorName string `json:"deductor_name" validate:"required"`
}

type TenantHierarchy struct {
	Tenant Tenant       `json:"tenant"`
	PANs   []PANWithSub `json:"pans"`
}

type PANWithSub struct {
	TenantPAN
	GSTINs []TenantGSTIN `json:"gstins"`
	TANs   []TenantTAN   `json:"tans"`
}
