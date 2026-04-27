package domain

import "time"

type GatewayResponse struct {
	Data interface{}  `json:"data"`
	Meta ResponseMeta `json:"meta"`
}

type ResponseMeta struct {
	RequestID      string `json:"request_id"`
	LatencyMs      int    `json:"latency_ms"`
	ProviderStatus string `json:"provider_status"`
}

// --- Authentication ---

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// --- IRN Generation ---

type GenerateIRNRequest struct {
	GSTIN     string      `json:"gstin"`
	DocDtls   DocDetails  `json:"doc_dtls"`
	SupDtls   PartyDetail `json:"sup_dtls"`
	BuyDtls   PartyDetail `json:"buy_dtls"`
	ItemList  []LineItem  `json:"item_list"`
	ValDtls   ValDetails  `json:"val_dtls"`
	RequestID string      `json:"request_id"`
}

type DocDetails struct {
	Typ string `json:"typ"` // INV, CRN, DBN
	No  string `json:"no"`
	Dt  string `json:"dt"` // DD/MM/YYYY
}

type PartyDetail struct {
	GSTIN string `json:"gstin"`
	LglNm string `json:"lgl_nm"`
	TrdNm string `json:"trd_nm,omitempty"`
	Addr1 string `json:"addr1,omitempty"`
	Loc   string `json:"loc,omitempty"`
	Pin   int    `json:"pin,omitempty"`
	Stcd  string `json:"stcd,omitempty"`
	Pos   string `json:"pos,omitempty"`
}

type LineItem struct {
	SlNo       string  `json:"sl_no"`
	PrdDesc    string  `json:"prd_desc"`
	HsnCd      string  `json:"hsn_cd"`
	Qty        float64 `json:"qty"`
	Unit       string  `json:"unit"`
	UnitPrice  float64 `json:"unit_price"`
	Discount   float64 `json:"discount"`
	TaxableAmt float64 `json:"taxable_amt"`
	IgstRt     float64 `json:"igst_rt"`
	IgstAmt    float64 `json:"igst_amt"`
	CgstRt     float64 `json:"cgst_rt"`
	CgstAmt    float64 `json:"cgst_amt"`
	SgstRt     float64 `json:"sgst_rt"`
	SgstAmt    float64 `json:"sgst_amt"`
	CesRt      float64 `json:"ces_rt"`
	CesAmt     float64 `json:"ces_amt"`
}

type ValDetails struct {
	TaxableVal float64 `json:"taxable_val"`
	IGST       float64 `json:"igst"`
	CGST       float64 `json:"cgst"`
	SGST       float64 `json:"sgst"`
	CesVal     float64 `json:"ces_val"`
	Discount   float64 `json:"discount"`
	OthChrg    float64 `json:"oth_chrg"`
	TotInvVal  float64 `json:"tot_inv_val"`
}

type GenerateIRNResponse struct {
	IRN           string    `json:"irn"`
	AckNo         string    `json:"ack_no"`
	AckDt         string    `json:"ack_dt"`
	SignedInvoice string    `json:"signed_invoice"`
	SignedQRCode  string    `json:"signed_qr_code"`
	Status        string    `json:"status"`
	GeneratedAt   time.Time `json:"generated_at"`
}

// --- IRN Cancellation ---

type CancelIRNRequest struct {
	IRN        string `json:"irn"`
	CnlRsn    string `json:"cnl_rsn"`     // 1=Duplicate, 2=Data entry mistake, 3=Order cancelled, 4=Others
	CnlRem    string `json:"cnl_rem"`
	RequestID string `json:"request_id"`
}

type CancelIRNResponse struct {
	IRN          string    `json:"irn"`
	CancelDate   string    `json:"cancel_date"`
	Status       string    `json:"status"`
	CancelledAt  time.Time `json:"cancelled_at"`
}

// --- Get IRN ---

type GetIRNByIRNRequest struct {
	IRN       string `json:"irn"`
	RequestID string `json:"request_id"`
}

type GetIRNByDocRequest struct {
	DocType   string `json:"doc_type"` // INV, CRN, DBN
	DocNum    string `json:"doc_num"`
	DocDate   string `json:"doc_date"` // DD/MM/YYYY
	RequestID string `json:"request_id"`
}

type GetIRNResponse struct {
	IRN            string     `json:"irn"`
	AckNo          string     `json:"ack_no"`
	AckDt          string     `json:"ack_dt"`
	Status         string     `json:"status"` // ACT, CANC
	DocType        string     `json:"doc_type"`
	DocNo          string     `json:"doc_no"`
	DocDate        string     `json:"doc_date"`
	SupplierGSTIN  string     `json:"supplier_gstin"`
	BuyerGSTIN     string     `json:"buyer_gstin"`
	TotalValue     float64    `json:"total_value"`
	SignedInvoice  string     `json:"signed_invoice"`
	SignedQRCode   string     `json:"signed_qr_code"`
	GeneratedAt    time.Time  `json:"generated_at"`
	CancelledAt    *time.Time `json:"cancelled_at,omitempty"`
}

// --- GSTIN Validation ---

type GSTINValidateRequest struct {
	GSTIN     string `json:"gstin"`
	RequestID string `json:"request_id"`
}

type GSTINValidateResponse struct {
	GSTIN      string `json:"gstin"`
	LegalName  string `json:"legal_name"`
	TradeName  string `json:"trade_name"`
	StateCode  string `json:"state_code"`
	Status     string `json:"status"` // Active, Inactive, Cancelled
	EntityType string `json:"entity_type"`
}
