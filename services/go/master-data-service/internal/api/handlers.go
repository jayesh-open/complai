package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/master-data-service/internal/domain"
	"github.com/complai/complai/services/go/master-data-service/internal/store"
)

type Handlers struct {
	store store.Repository
}

func NewHandlers(s store.Repository) *Handlers {
	return &Handlers{store: s}
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "master-data-service"})
}

// ---------------------------------------------------------------------------
// Vendors
// ---------------------------------------------------------------------------

func (h *Handlers) CreateVendor(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.CreateVendorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	v := &domain.Vendor{
		Name:         req.Name,
		PAN:          req.PAN,
		GSTIN:        req.GSTIN,
		Email:        req.Email,
		Phone:        req.Phone,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		StateCode:    req.StateCode,
		Pincode:      req.Pincode,
		BankName:     req.BankName,
		BankAccount:  req.BankAccount,
		BankIFSC:     req.BankIFSC,
		Metadata:     "{}",
	}

	if err := h.store.CreateVendor(r.Context(), tenantID, v); err != nil {
		log.Error().Err(err).Msg("create vendor failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "create failed"})
		return
	}

	httputil.JSON(w, http.StatusCreated, v)
}

func (h *Handlers) GetVendor(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	vendorID, err := uuid.Parse(r.PathValue("vendorID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid vendor_id"})
		return
	}

	v, err := h.store.GetVendor(r.Context(), tenantID, vendorID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "vendor not found"})
		return
	}

	httputil.JSON(w, http.StatusOK, v)
}

func (h *Handlers) ListVendors(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	vendors, err := h.store.ListVendors(r.Context(), tenantID)
	if err != nil {
		log.Error().Err(err).Msg("list vendors failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if vendors == nil {
		vendors = []domain.Vendor{}
	}

	httputil.JSON(w, http.StatusOK, vendors)
}

func (h *Handlers) UpdateVendor(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	vendorID, err := uuid.Parse(r.PathValue("vendorID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid vendor_id"})
		return
	}

	var req domain.UpdateVendorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	v, err := h.store.UpdateVendor(r.Context(), tenantID, vendorID, &req)
	if err != nil {
		log.Error().Err(err).Msg("update vendor failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "update failed"})
		return
	}

	httputil.JSON(w, http.StatusOK, v)
}

// ---------------------------------------------------------------------------
// Customers
// ---------------------------------------------------------------------------

func (h *Handlers) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	c := &domain.Customer{
		Name:             req.Name,
		PAN:              req.PAN,
		GSTIN:            req.GSTIN,
		Email:            req.Email,
		Phone:            req.Phone,
		AddressLine1:     req.AddressLine1,
		City:             req.City,
		StateCode:        req.StateCode,
		Pincode:          req.Pincode,
		PaymentTermsDays: req.PaymentTermsDays,
		CreditLimit:      req.CreditLimit,
		Metadata:         "{}",
	}

	if err := h.store.CreateCustomer(r.Context(), tenantID, c); err != nil {
		log.Error().Err(err).Msg("create customer failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "create failed"})
		return
	}

	httputil.JSON(w, http.StatusCreated, c)
}

func (h *Handlers) GetCustomer(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	customerID, err := uuid.Parse(r.PathValue("customerID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid customer_id"})
		return
	}

	c, err := h.store.GetCustomer(r.Context(), tenantID, customerID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "customer not found"})
		return
	}

	httputil.JSON(w, http.StatusOK, c)
}

func (h *Handlers) ListCustomers(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	customers, err := h.store.ListCustomers(r.Context(), tenantID)
	if err != nil {
		log.Error().Err(err).Msg("list customers failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if customers == nil {
		customers = []domain.Customer{}
	}

	httputil.JSON(w, http.StatusOK, customers)
}

// ---------------------------------------------------------------------------
// Items
// ---------------------------------------------------------------------------

func (h *Handlers) CreateItem(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	uom := req.UnitOfMeasure
	if uom == "" {
		uom = "NOS"
	}

	i := &domain.Item{
		Name:          req.Name,
		Description:   req.Description,
		HSNCode:       req.HSNCode,
		UnitOfMeasure: uom,
		UnitPrice:     req.UnitPrice,
		GSTRate:       req.GSTRate,
		IsService:     req.IsService,
		Metadata:      "{}",
	}

	if err := h.store.CreateItem(r.Context(), tenantID, i); err != nil {
		log.Error().Err(err).Msg("create item failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "create failed"})
		return
	}

	httputil.JSON(w, http.StatusCreated, i)
}

func (h *Handlers) GetItem(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	itemID, err := uuid.Parse(r.PathValue("itemID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid item_id"})
		return
	}

	i, err := h.store.GetItem(r.Context(), tenantID, itemID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "item not found"})
		return
	}

	httputil.JSON(w, http.StatusOK, i)
}

func (h *Handlers) ListItems(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	items, err := h.store.ListItems(r.Context(), tenantID)
	if err != nil {
		log.Error().Err(err).Msg("list items failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if items == nil {
		items = []domain.Item{}
	}

	httputil.JSON(w, http.StatusOK, items)
}

// ---------------------------------------------------------------------------
// HSN Codes
// ---------------------------------------------------------------------------

func (h *Handlers) ListHSNCodes(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	codes, err := h.store.ListHSNCodes(r.Context(), tenantID)
	if err != nil {
		log.Error().Err(err).Msg("list hsn codes failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if codes == nil {
		codes = []domain.HSNCode{}
	}

	httputil.JSON(w, http.StatusOK, codes)
}

func (h *Handlers) CreateHSNCode(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.CreateHSNCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	h2 := &domain.HSNCode{
		Code:          req.Code,
		Description:   req.Description,
		GSTRate:       req.GSTRate,
		EffectiveFrom: req.EffectiveFrom,
	}

	if err := h.store.CreateHSNCode(r.Context(), tenantID, h2); err != nil {
		log.Error().Err(err).Msg("create hsn code failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "create failed"})
		return
	}

	httputil.JSON(w, http.StatusCreated, h2)
}

// ---------------------------------------------------------------------------
// State Codes
// ---------------------------------------------------------------------------

func (h *Handlers) ListStateCodes(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	codes, err := h.store.ListStateCodes(r.Context(), tenantID)
	if err != nil {
		log.Error().Err(err).Msg("list state codes failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if codes == nil {
		codes = []domain.StateCode{}
	}

	httputil.JSON(w, http.StatusOK, codes)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func tenantIDFromRequest(r *http.Request) (uuid.UUID, error) {
	h := r.Header.Get("X-Tenant-Id")
	if h == "" {
		return uuid.Nil, fmt.Errorf("missing X-Tenant-Id header")
	}
	return uuid.Parse(h)
}
