package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type EWBStatus string

const (
	EWBStatusPending        EWBStatus = "PENDING"
	EWBStatusActive         EWBStatus = "ACTIVE"
	EWBStatusVehicleUpdated EWBStatus = "VEHICLE_UPDATED"
	EWBStatusExtended       EWBStatus = "EXTENDED"
	EWBStatusCancelled      EWBStatus = "CANCELLED"
	EWBStatusConsolidated   EWBStatus = "CONSOLIDATED"
)

func CanTransitionTo(from, to EWBStatus) bool {
	switch to {
	case EWBStatusActive:
		return from == EWBStatusPending
	case EWBStatusVehicleUpdated:
		return from == EWBStatusActive || from == EWBStatusVehicleUpdated || from == EWBStatusExtended
	case EWBStatusExtended:
		return from == EWBStatusActive || from == EWBStatusVehicleUpdated || from == EWBStatusExtended
	case EWBStatusCancelled:
		return from == EWBStatusActive || from == EWBStatusVehicleUpdated || from == EWBStatusExtended
	case EWBStatusConsolidated:
		return from == EWBStatusActive || from == EWBStatusVehicleUpdated || from == EWBStatusExtended
	default:
		return false
	}
}

type EWayBill struct {
	ID                uuid.UUID       `json:"id"`
	TenantID          uuid.UUID       `json:"tenant_id"`
	EWBNumber         string          `json:"ewb_number"`
	EWBDate           string          `json:"ewb_date"`
	DocType           string          `json:"doc_type"`
	DocNumber         string          `json:"doc_number"`
	DocDate           string          `json:"doc_date"`
	SupplierGSTIN     string          `json:"supplier_gstin"`
	SupplierName      string          `json:"supplier_name"`
	BuyerGSTIN        string          `json:"buyer_gstin"`
	BuyerName         string          `json:"buyer_name"`
	SupplyType        string          `json:"supply_type"`
	SubSupplyType     string          `json:"sub_supply_type"`
	TransportMode     string          `json:"transport_mode"`
	VehicleNumber     string          `json:"vehicle_number"`
	VehicleType       string          `json:"vehicle_type"`
	TransporterID     string          `json:"transporter_id"`
	TransporterName   string          `json:"transporter_name"`
	FromPlace         string          `json:"from_place"`
	FromState         string          `json:"from_state"`
	FromPincode       string          `json:"from_pincode"`
	ToPlace           string          `json:"to_place"`
	ToState           string          `json:"to_state"`
	ToPincode         string          `json:"to_pincode"`
	DistanceKM        int             `json:"distance_km"`
	TaxableValue      decimal.Decimal `json:"taxable_value"`
	CGSTAmount        decimal.Decimal `json:"cgst_amount"`
	SGSTAmount        decimal.Decimal `json:"sgst_amount"`
	IGSTAmount        decimal.Decimal `json:"igst_amount"`
	CessAmount        decimal.Decimal `json:"cess_amount"`
	TotalValue        decimal.Decimal `json:"total_value"`
	Status            EWBStatus       `json:"status"`
	ValidFrom         *time.Time      `json:"valid_from,omitempty"`
	ValidUntil        *time.Time      `json:"valid_until,omitempty"`
	GeneratedAt       *time.Time      `json:"generated_at,omitempty"`
	CancelledAt       *time.Time      `json:"cancelled_at,omitempty"`
	CancelReason      string          `json:"cancel_reason,omitempty"`
	ConsolidatedEWBID *uuid.UUID      `json:"consolidated_ewb_id,omitempty"`
	RequestID         uuid.UUID       `json:"request_id"`
	SourceSystem      string          `json:"source_system"`
	SourceID          string          `json:"source_id,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

type EWBItem struct {
	ID           uuid.UUID       `json:"id"`
	EWBID        uuid.UUID       `json:"ewb_id"`
	TenantID     uuid.UUID       `json:"tenant_id"`
	ItemNumber   int             `json:"item_number"`
	ProductName  string          `json:"product_name"`
	ProductDesc  string          `json:"product_desc"`
	HSNCode      string          `json:"hsn_code"`
	Quantity     decimal.Decimal `json:"quantity"`
	Unit         string          `json:"unit"`
	TaxableValue decimal.Decimal `json:"taxable_value"`
	CGSTRate     decimal.Decimal `json:"cgst_rate"`
	SGSTRate     decimal.Decimal `json:"sgst_rate"`
	IGSTRate     decimal.Decimal `json:"igst_rate"`
	CessRate     decimal.Decimal `json:"cess_rate"`
	CreatedAt    time.Time       `json:"created_at"`
}

type VehicleUpdate struct {
	ID            uuid.UUID `json:"id"`
	EWBID         uuid.UUID `json:"ewb_id"`
	TenantID      uuid.UUID `json:"tenant_id"`
	VehicleNumber string    `json:"vehicle_number"`
	FromPlace     string    `json:"from_place"`
	FromState     string    `json:"from_state"`
	TransportMode string    `json:"transport_mode"`
	Reason        string    `json:"reason"`
	Remark        string    `json:"remark"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Consolidation struct {
	ID                    uuid.UUID  `json:"id"`
	TenantID              uuid.UUID  `json:"tenant_id"`
	ConsolidatedEWBNumber string     `json:"consolidated_ewb_number"`
	TripSheetNumber       string     `json:"trip_sheet_number"`
	VehicleNumber         string     `json:"vehicle_number"`
	FromPlace             string     `json:"from_place"`
	FromState             string     `json:"from_state"`
	ToPlace               string     `json:"to_place"`
	ToState               string     `json:"to_state"`
	TransportMode         string     `json:"transport_mode"`
	Status                string     `json:"status"`
	GeneratedAt           *time.Time `json:"generated_at,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

type GenerateEWBRequest struct {
	DocType       string            `json:"doc_type"`
	DocNumber     string            `json:"doc_number"`
	DocDate       string            `json:"doc_date"`
	SupplierGSTIN string            `json:"supplier_gstin"`
	SupplierName  string            `json:"supplier_name"`
	BuyerGSTIN    string            `json:"buyer_gstin"`
	BuyerName     string            `json:"buyer_name"`
	SupplyType    string            `json:"supply_type"`
	SubSupplyType string            `json:"sub_supply_type"`
	TransportMode string            `json:"transport_mode"`
	VehicleNumber string            `json:"vehicle_number"`
	VehicleType   string            `json:"vehicle_type"`
	TransporterID string            `json:"transporter_id"`
	FromPlace     string            `json:"from_place"`
	FromState     string            `json:"from_state"`
	FromPincode   string            `json:"from_pincode"`
	ToPlace       string            `json:"to_place"`
	ToState       string            `json:"to_state"`
	ToPincode     string            `json:"to_pincode"`
	DistanceKM    int               `json:"distance_km"`
	TaxableValue  decimal.Decimal   `json:"taxable_value"`
	CGSTAmount    decimal.Decimal   `json:"cgst_amount"`
	SGSTAmount    decimal.Decimal   `json:"sgst_amount"`
	IGSTAmount    decimal.Decimal   `json:"igst_amount"`
	CessAmount    decimal.Decimal   `json:"cess_amount"`
	TotalValue    decimal.Decimal   `json:"total_value"`
	SourceSystem  string            `json:"source_system"`
	SourceID      string            `json:"source_id"`
	Items         []ItemRequest     `json:"items"`
}

type ItemRequest struct {
	ProductName  string          `json:"product_name"`
	HSNCode      string          `json:"hsn_code"`
	Quantity     decimal.Decimal `json:"quantity"`
	Unit         string          `json:"unit"`
	TaxableValue decimal.Decimal `json:"taxable_value"`
	CGSTRate     decimal.Decimal `json:"cgst_rate"`
	SGSTRate     decimal.Decimal `json:"sgst_rate"`
	IGSTRate     decimal.Decimal `json:"igst_rate"`
	CessRate     decimal.Decimal `json:"cess_rate"`
}

type GenerateEWBResponse struct {
	ID         uuid.UUID `json:"id"`
	EWBNumber  string    `json:"ewb_number"`
	Status     EWBStatus `json:"status"`
	ValidFrom  time.Time `json:"valid_from"`
	ValidUntil time.Time `json:"valid_until"`
}

type CancelEWBRequest struct {
	Reason string `json:"reason"`
	Remark string `json:"remark"`
}

type CancelEWBResponse struct {
	ID          uuid.UUID `json:"id"`
	EWBNumber   string    `json:"ewb_number"`
	Status      EWBStatus `json:"status"`
	CancelledAt time.Time `json:"cancelled_at"`
}

type UpdateVehicleRequest struct {
	VehicleNumber string `json:"vehicle_number"`
	FromPlace     string `json:"from_place"`
	FromState     string `json:"from_state"`
	TransportMode string `json:"transport_mode"`
	Reason        string `json:"reason"`
	Remark        string `json:"remark"`
}

type UpdateVehicleResponse struct {
	ID            uuid.UUID `json:"id"`
	EWBNumber     string    `json:"ewb_number"`
	VehicleNumber string    `json:"vehicle_number"`
	Status        EWBStatus `json:"status"`
}

type ExtendValidityRequest struct {
	RemainingDistance int    `json:"remaining_distance"`
	FromPlace         string `json:"from_place"`
	FromState         string `json:"from_state"`
	ExtendReason      string `json:"extend_reason"`
	TransitType       string `json:"transit_type"`
	ConsignmentStatus string `json:"consignment_status"`
	Remark            string `json:"remark"`
}

type ExtendValidityResponse struct {
	ID         uuid.UUID `json:"id"`
	EWBNumber  string    `json:"ewb_number"`
	Status     EWBStatus `json:"status"`
	ValidUntil time.Time `json:"valid_until"`
}

type ConsolidateRequest struct {
	EWBIDS        []uuid.UUID `json:"ewb_ids"`
	VehicleNumber string      `json:"vehicle_number"`
	FromPlace     string      `json:"from_place"`
	FromState     string      `json:"from_state"`
	ToPlace       string      `json:"to_place"`
	ToState       string      `json:"to_state"`
	TransportMode string      `json:"transport_mode"`
}

type ConsolidateResponse struct {
	ConsolidationID       uuid.UUID `json:"consolidation_id"`
	ConsolidatedEWBNumber string    `json:"consolidated_ewb_number"`
	EWBCount              int       `json:"ewb_count"`
	Status                string    `json:"status"`
}

type ListEWBRequest struct {
	GSTIN      string `json:"gstin"`
	Status     string `json:"status"`
	PageSize   int    `json:"page_size"`
	PageOffset int    `json:"page_offset"`
}

type ListEWBResponse struct {
	EWayBills  []EWayBill `json:"eway_bills"`
	TotalCount int        `json:"total_count"`
	PageSize   int        `json:"page_size"`
	PageOffset int        `json:"page_offset"`
}
