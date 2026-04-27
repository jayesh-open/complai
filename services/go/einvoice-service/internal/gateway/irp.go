package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type IRPClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewIRPClient(baseURL string) *IRPClient {
	return &IRPClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type GenerateIRNGatewayRequest struct {
	GSTIN   string      `json:"gstin"`
	DocDtls DocDetails  `json:"doc_dtls"`
	SupDtls PartyDetail `json:"sup_dtls"`
	BuyDtls PartyDetail `json:"buy_dtls"`
	ValDtls ValDetails  `json:"val_dtls"`
}

type DocDetails struct {
	Typ string `json:"typ"`
	No  string `json:"no"`
	Dt  string `json:"dt"`
}

type PartyDetail struct {
	GSTIN string `json:"gstin"`
	LglNm string `json:"lgl_nm"`
	Stcd  string `json:"stcd,omitempty"`
	Pos   string `json:"pos,omitempty"`
}

type ValDetails struct {
	TaxableVal float64 `json:"taxable_val"`
	IGST       float64 `json:"igst"`
	CGST       float64 `json:"cgst"`
	SGST       float64 `json:"sgst"`
	CesVal     float64 `json:"ces_val"`
	TotInvVal  float64 `json:"tot_inv_val"`
}

type GenerateIRNGatewayResponse struct {
	IRN           string `json:"irn"`
	AckNo         string `json:"ack_no"`
	AckDt         string `json:"ack_dt"`
	SignedInvoice string `json:"signed_invoice"`
	SignedQRCode  string `json:"signed_qr_code"`
	Status        string `json:"status"`
}

type CancelIRNGatewayRequest struct {
	IRN    string `json:"irn"`
	CnlRsn string `json:"cnl_rsn"`
	CnlRem string `json:"cnl_rem"`
}

type CancelIRNGatewayResponse struct {
	IRN        string `json:"irn"`
	CancelDate string `json:"cancel_date"`
	Status     string `json:"status"`
}

type gatewayDataWrapper struct {
	Data json.RawMessage `json:"data"`
}

type gatewayResponseWrapper struct {
	Data gatewayDataWrapper `json:"data"`
}

func (c *IRPClient) GenerateIRN(ctx context.Context, tenantID uuid.UUID, req *GenerateIRNGatewayRequest) (*GenerateIRNGatewayResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/gateway/irp/invoice", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Tenant-Id", tenantID.String())
	httpReq.Header.Set("X-Idempotency-Key", uuid.New().String())

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("call irp gateway: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("irp gateway returned %d", resp.StatusCode)
	}

	var wrapper gatewayResponseWrapper
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	var innerData struct {
		Data GenerateIRNGatewayResponse `json:"data"`
	}
	if err := json.Unmarshal(wrapper.Data.Data, &innerData.Data); err != nil {
		log.Warn().Err(err).Msg("fallback: trying direct unmarshal")
		var direct GenerateIRNGatewayResponse
		if err2 := json.Unmarshal(wrapper.Data.Data, &direct); err2 != nil {
			return nil, fmt.Errorf("decode inner data: %w", err2)
		}
		return &direct, nil
	}

	return &innerData.Data, nil
}

func (c *IRPClient) CancelIRN(ctx context.Context, tenantID uuid.UUID, req *CancelIRNGatewayRequest) (*CancelIRNGatewayResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/gateway/irp/invoice/cancel", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Tenant-Id", tenantID.String())
	httpReq.Header.Set("X-Idempotency-Key", uuid.New().String())

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("call irp gateway cancel: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("irp gateway cancel returned %d", resp.StatusCode)
	}

	var wrapper gatewayResponseWrapper
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	var result CancelIRNGatewayResponse
	if err := json.Unmarshal(wrapper.Data.Data, &result); err != nil {
		return nil, fmt.Errorf("decode cancel data: %w", err)
	}

	return &result, nil
}
