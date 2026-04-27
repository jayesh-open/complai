package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type EWBClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewEWBClient(baseURL string) *EWBClient {
	return &EWBClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type GenerateEWBGatewayRequest struct {
	GSTIN         string          `json:"gstin"`
	SupplyType    string          `json:"supply_type"`
	SubSupplyType string          `json:"sub_supply_type"`
	DocType       string          `json:"doc_type"`
	DocNo         string          `json:"doc_no"`
	DocDate       string          `json:"doc_date"`
	FromGSTIN     string          `json:"from_gstin"`
	FromName      string          `json:"from_name"`
	FromPlace     string          `json:"from_place"`
	FromState     string          `json:"from_state"`
	FromPincode   string          `json:"from_pincode"`
	ToGSTIN       string          `json:"to_gstin"`
	ToName        string          `json:"to_name"`
	ToPlace       string          `json:"to_place"`
	ToState       string          `json:"to_state"`
	ToPincode     string          `json:"to_pincode"`
	TransportMode string          `json:"transport_mode"`
	VehicleNo     string          `json:"vehicle_no"`
	VehicleType   string          `json:"vehicle_type"`
	TransporterID string          `json:"transporter_id"`
	DistanceKM    int             `json:"distance_km"`
	TotalValue    float64         `json:"total_value"`
	TaxableValue  float64         `json:"taxable_value"`
	CGSTAmount    float64         `json:"cgst_amount"`
	SGSTAmount    float64         `json:"sgst_amount"`
	IGSTAmount    float64         `json:"igst_amount"`
	CessAmount    float64         `json:"cess_amount"`
	Items         []GatewayItem   `json:"items"`
}

type GatewayItem struct {
	ProductName string  `json:"product_name"`
	HSNCode     string  `json:"hsn_code"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit"`
	TaxableVal  float64 `json:"taxable_value"`
	CGSTRate    float64 `json:"cgst_rate"`
	SGSTRate    float64 `json:"sgst_rate"`
	IGSTRate    float64 `json:"igst_rate"`
	CessRate    float64 `json:"cess_rate"`
}

type GenerateEWBGatewayResponse struct {
	EWBNumber  string `json:"ewb_no"`
	EWBDate    string `json:"ewb_date"`
	ValidUntil string `json:"valid_until"`
	Status     string `json:"status"`
}

type CancelEWBGatewayRequest struct {
	EWBNo  string `json:"ewb_no"`
	Reason string `json:"reason"`
	Remark string `json:"remark"`
}

type CancelEWBGatewayResponse struct {
	EWBNo      string `json:"ewb_no"`
	CancelDate string `json:"cancel_date"`
	Status     string `json:"status"`
}

type UpdateVehicleGatewayRequest struct {
	EWBNo         string `json:"ewb_no"`
	VehicleNo     string `json:"vehicle_no"`
	FromPlace     string `json:"from_place"`
	FromState     string `json:"from_state"`
	Reason        string `json:"reason"`
	TransportMode string `json:"transport_mode"`
	Remark        string `json:"remark"`
}

type UpdateVehicleGatewayResponse struct {
	EWBNo      string `json:"ewb_no"`
	VehicleNo  string `json:"vehicle_no"`
	ValidUntil string `json:"valid_until"`
	Status     string `json:"status"`
}

type ExtendValidityGatewayRequest struct {
	EWBNo             string `json:"ewb_no"`
	FromPlace         string `json:"from_place"`
	FromState         string `json:"from_state"`
	RemainingDistance  int    `json:"remaining_distance"`
	ExtendReason      string `json:"extend_reason"`
	TransitType       string `json:"transit_type"`
	ConsignmentStatus string `json:"consignment_status"`
	Remark            string `json:"remark"`
}

type ExtendValidityGatewayResponse struct {
	EWBNo      string `json:"ewb_no"`
	ValidUntil string `json:"valid_until"`
	Status     string `json:"status"`
}

type ConsolidateGatewayRequest struct {
	FromGSTIN     string   `json:"from_gstin"`
	FromPlace     string   `json:"from_place"`
	FromState     string   `json:"from_state"`
	ToPlace       string   `json:"to_place"`
	ToState       string   `json:"to_state"`
	VehicleNo     string   `json:"vehicle_no"`
	TransportMode string   `json:"transport_mode"`
	EWBNumbers    []string `json:"ewb_numbers"`
}

type ConsolidateGatewayResponse struct {
	ConsolidatedEWBNo string `json:"consolidated_ewb_no"`
	Status            string `json:"status"`
}

func (c *EWBClient) GenerateEWB(ctx context.Context, tenantID uuid.UUID, req *GenerateEWBGatewayRequest) (*GenerateEWBGatewayResponse, error) {
	var resp GenerateEWBGatewayResponse
	if err := c.post(ctx, "/v1/gateway/ewb/generate", tenantID, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *EWBClient) CancelEWB(ctx context.Context, tenantID uuid.UUID, req *CancelEWBGatewayRequest) (*CancelEWBGatewayResponse, error) {
	var resp CancelEWBGatewayResponse
	if err := c.post(ctx, "/v1/gateway/ewb/cancel", tenantID, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *EWBClient) UpdateVehicle(ctx context.Context, tenantID uuid.UUID, req *UpdateVehicleGatewayRequest) (*UpdateVehicleGatewayResponse, error) {
	var resp UpdateVehicleGatewayResponse
	if err := c.post(ctx, "/v1/gateway/ewb/vehicle", tenantID, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *EWBClient) ExtendValidity(ctx context.Context, tenantID uuid.UUID, req *ExtendValidityGatewayRequest) (*ExtendValidityGatewayResponse, error) {
	var resp ExtendValidityGatewayResponse
	if err := c.post(ctx, "/v1/gateway/ewb/extend", tenantID, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *EWBClient) Consolidate(ctx context.Context, tenantID uuid.UUID, req *ConsolidateGatewayRequest) (*ConsolidateGatewayResponse, error) {
	var resp ConsolidateGatewayResponse
	if err := c.post(ctx, "/v1/gateway/ewb/consolidate", tenantID, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *EWBClient) post(ctx context.Context, path string, tenantID uuid.UUID, reqBody, respTarget interface{}) error {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("gateway error (status %d): %s", resp.StatusCode, string(respBody))
	}

	// httputil.JSON wraps in {"data": ...}, gateway also wraps in {"data": ..., "meta": ...}
	var outer struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(respBody, &outer); err != nil {
		return fmt.Errorf("decode outer wrapper: %w", err)
	}

	var inner struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(outer.Data, &inner); err != nil {
		return fmt.Errorf("decode inner wrapper: %w", err)
	}

	if err := json.Unmarshal(inner.Data, respTarget); err != nil {
		return fmt.Errorf("decode response data: %w", err)
	}
	return nil
}
