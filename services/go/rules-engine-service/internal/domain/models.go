package domain

import (
	"time"

	"github.com/google/uuid"
)

type Rule struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	Category      string     `json:"category"`
	Name          string     `json:"name"`
	Description   *string    `json:"description,omitempty"`
	Version       int        `json:"version"`
	Priority      int        `json:"priority"`
	Conditions    string     `json:"conditions"`
	Actions       string     `json:"actions"`
	EffectiveFrom time.Time  `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty"`
	Status        string     `json:"status"`
	CreatedBy     *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type RuleExecutionLog struct {
	ID              uuid.UUID  `json:"id"`
	TenantID        uuid.UUID  `json:"tenant_id"`
	RuleID          *uuid.UUID `json:"rule_id,omitempty"`
	InputData       string     `json:"input_data"`
	MatchedRules    *string    `json:"matched_rules,omitempty"`
	Output          string     `json:"output"`
	ExecutionTimeMs int        `json:"execution_time_ms"`
	TraceID         *string    `json:"trace_id,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type CreateRuleRequest struct {
	Category    string  `json:"category"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Priority    int     `json:"priority"`
	Conditions  string  `json:"conditions"`
	Actions     string  `json:"actions"`
}

// EvaluateRequest is the main input for rule evaluation.
type EvaluateRequest struct {
	Category string        `json:"category"`
	Input    EvaluateInput `json:"input"`
}

type EvaluateInput struct {
	// For tax_determination
	SupplyType    string  `json:"supply_type,omitempty"`
	PlaceOfSupply string  `json:"place_of_supply,omitempty"`
	SupplierState string  `json:"supplier_state,omitempty"`
	HSNCode       string  `json:"hsn_code,omitempty"`
	TaxableValue  float64 `json:"taxable_value,omitempty"`

	// For tds_applicability
	Section       string  `json:"section,omitempty"`
	PaymentAmount float64 `json:"payment_amount,omitempty"`
	DeducteeType  string  `json:"deductee_type,omitempty"`

	// For hsn_validation — uses HSNCode field above
}

type EvaluateResult struct {
	Category     string   `json:"category"`
	Matched      bool     `json:"matched"`
	MatchedRules []string `json:"matched_rules,omitempty"`
	Output       TaxOutput `json:"output"`
	ExecutionMs  int      `json:"execution_ms"`
}

type TaxOutput struct {
	// Tax determination output
	CGST    *float64 `json:"cgst,omitempty"`
	SGST    *float64 `json:"sgst,omitempty"`
	IGST    *float64 `json:"igst,omitempty"`
	GSTRate *float64 `json:"gst_rate,omitempty"`
	TaxType string   `json:"tax_type,omitempty"`

	// HSN validation output
	HSNValid       *bool   `json:"hsn_valid,omitempty"`
	HSNDescription *string `json:"hsn_description,omitempty"`

	// TDS output
	TDSApplicable  *bool    `json:"tds_applicable,omitempty"`
	TDSRate        *float64 `json:"tds_rate,omitempty"`
	TDSSection     *string  `json:"tds_section,omitempty"`
	TDSAmount      *float64 `json:"tds_amount,omitempty"`
	ThresholdLimit *float64 `json:"threshold_limit,omitempty"`
}
