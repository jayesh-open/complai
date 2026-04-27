package domain

type GenerateEWBRequest struct {
	GSTIN         string    `json:"gstin"`
	SupplyType    string    `json:"supply_type"`
	SubSupplyType string    `json:"sub_supply_type"`
	DocType       string    `json:"doc_type"`
	DocNo         string    `json:"doc_no"`
	DocDate       string    `json:"doc_date"`
	FromGSTIN     string    `json:"from_gstin"`
	FromName      string    `json:"from_name"`
	FromPlace     string    `json:"from_place"`
	FromState     string    `json:"from_state"`
	FromPincode   string    `json:"from_pincode"`
	ToGSTIN       string    `json:"to_gstin"`
	ToName        string    `json:"to_name"`
	ToPlace       string    `json:"to_place"`
	ToState       string    `json:"to_state"`
	ToPincode     string    `json:"to_pincode"`
	TransportMode string    `json:"transport_mode"`
	VehicleNo     string    `json:"vehicle_no"`
	VehicleType   string    `json:"vehicle_type"`
	TransporterID string    `json:"transporter_id"`
	DistanceKM    int       `json:"distance_km"`
	TotalValue    float64   `json:"total_value"`
	TaxableValue  float64   `json:"taxable_value"`
	CGSTAmount    float64   `json:"cgst_amount"`
	SGSTAmount    float64   `json:"sgst_amount"`
	IGSTAmount    float64   `json:"igst_amount"`
	CessAmount    float64   `json:"cess_amount"`
	Items         []EWBItem `json:"items"`
}

type EWBItem struct {
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

type GenerateEWBResponse struct {
	EWBNumber  string `json:"ewb_no"`
	EWBDate    string `json:"ewb_date"`
	ValidUntil string `json:"valid_until"`
	Status     string `json:"status"`
}

type CancelEWBRequest struct {
	EWBNo  string `json:"ewb_no"`
	Reason string `json:"reason"`
	Remark string `json:"remark"`
}

type CancelEWBResponse struct {
	EWBNo      string `json:"ewb_no"`
	CancelDate string `json:"cancel_date"`
	Status     string `json:"status"`
}

type UpdateVehicleRequest struct {
	EWBNo         string `json:"ewb_no"`
	VehicleNo     string `json:"vehicle_no"`
	FromPlace     string `json:"from_place"`
	FromState     string `json:"from_state"`
	Reason        string `json:"reason"`
	TransportMode string `json:"transport_mode"`
	Remark        string `json:"remark"`
}

type UpdateVehicleResponse struct {
	EWBNo      string `json:"ewb_no"`
	VehicleNo  string `json:"vehicle_no"`
	ValidUntil string `json:"valid_until"`
	Status     string `json:"status"`
}

type ExtendValidityRequest struct {
	EWBNo             string `json:"ewb_no"`
	FromPlace         string `json:"from_place"`
	FromState         string `json:"from_state"`
	RemainingDistance int    `json:"remaining_distance"`
	ExtendReason      string `json:"extend_reason"`
	TransitType       string `json:"transit_type"`
	ConsignmentStatus string `json:"consignment_status"`
	Remark            string `json:"remark"`
}

type ExtendValidityResponse struct {
	EWBNo      string `json:"ewb_no"`
	ValidUntil string `json:"valid_until"`
	Status     string `json:"status"`
}

type ConsolidateEWBRequest struct {
	FromGSTIN     string   `json:"from_gstin"`
	FromPlace     string   `json:"from_place"`
	FromState     string   `json:"from_state"`
	ToPlace       string   `json:"to_place"`
	ToState       string   `json:"to_state"`
	VehicleNo     string   `json:"vehicle_no"`
	TransportMode string   `json:"transport_mode"`
	EWBNumbers    []string `json:"ewb_numbers"`
}

type ConsolidateEWBResponse struct {
	ConsolidatedEWBNo string `json:"consolidated_ewb_no"`
	Status            string `json:"status"`
}

type GetEWBResponse struct {
	EWBNumber  string  `json:"ewb_no"`
	EWBDate    string  `json:"ewb_date"`
	DocType    string  `json:"doc_type"`
	DocNo      string  `json:"doc_no"`
	DocDate    string  `json:"doc_date"`
	FromGSTIN  string  `json:"from_gstin"`
	FromName   string  `json:"from_name"`
	ToGSTIN    string  `json:"to_gstin"`
	ToName     string  `json:"to_name"`
	VehicleNo  string  `json:"vehicle_no"`
	Status     string  `json:"status"`
	ValidUntil string  `json:"valid_until"`
	DistanceKM int     `json:"distance_km"`
	TotalValue float64 `json:"total_value"`
}

type GatewayResponse struct {
	Data interface{}  `json:"data"`
	Meta ResponseMeta `json:"meta"`
}

type ResponseMeta struct {
	RequestID      string `json:"request_id"`
	LatencyMS      int64  `json:"latency_ms"`
	ProviderStatus string `json:"provider_status"`
}
