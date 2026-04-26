package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/services/go/recon-service/internal/domain"
	"github.com/complai/complai/services/go/recon-service/internal/gateway"
)

// mockStore implements store.Repository for tests.
type mockStore struct {
	runs    map[uuid.UUID]*domain.ReconRun
	matches map[uuid.UUID]*domain.ReconMatch
	actions []domain.IMSAction
}

func newMockStore() *mockStore {
	return &mockStore{
		runs:    make(map[uuid.UUID]*domain.ReconRun),
		matches: make(map[uuid.UUID]*domain.ReconMatch),
	}
}

func (m *mockStore) CreateRun(_ context.Context, tenantID uuid.UUID, run *domain.ReconRun) error {
	run.ID = uuid.New()
	run.TenantID = tenantID
	run.RequestID = uuid.New()
	run.CreatedAt = time.Now().UTC()
	m.runs[run.ID] = run
	return nil
}

func (m *mockStore) GetRun(_ context.Context, _ uuid.UUID, runID uuid.UUID) (*domain.ReconRun, error) {
	r, ok := m.runs[runID]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return r, nil
}

func (m *mockStore) UpdateRun(_ context.Context, _ uuid.UUID, run *domain.ReconRun) error {
	m.runs[run.ID] = run
	return nil
}

func (m *mockStore) BulkInsertMatches(_ context.Context, tenantID uuid.UUID, matches []domain.ReconMatch) error {
	for i := range matches {
		matches[i].TenantID = tenantID
		matches[i].CreatedAt = time.Now().UTC()
		matches[i].UpdatedAt = time.Now().UTC()
		m.matches[matches[i].ID] = &matches[i]
	}
	return nil
}

func (m *mockStore) ListMatches(_ context.Context, _ uuid.UUID, runID uuid.UUID, matchType string, status string, limit, offset int) ([]domain.ReconMatch, error) {
	var results []domain.ReconMatch
	for _, match := range m.matches {
		if match.RunID != runID {
			continue
		}
		if matchType != "" && string(match.MatchType) != matchType {
			continue
		}
		if status != "" && string(match.Status) != status {
			continue
		}
		results = append(results, *match)
	}
	if offset >= len(results) {
		return nil, nil
	}
	end := offset + limit
	if end > len(results) {
		end = len(results)
	}
	return results[offset:end], nil
}

func (m *mockStore) GetMatch(_ context.Context, _ uuid.UUID, matchID uuid.UUID) (*domain.ReconMatch, error) {
	match, ok := m.matches[matchID]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return match, nil
}

func (m *mockStore) UpdateMatchStatus(_ context.Context, _ uuid.UUID, matchID uuid.UUID, status domain.MatchStatus, acceptedBy *uuid.UUID) error {
	match, ok := m.matches[matchID]
	if !ok {
		return fmt.Errorf("not found")
	}
	match.Status = status
	match.AcceptedBy = acceptedBy
	if status == domain.MatchStatusAccepted {
		now := time.Now().UTC()
		match.AcceptedAt = &now
	}
	return nil
}

func (m *mockStore) GetBucketSummary(_ context.Context, _ uuid.UUID, runID uuid.UUID) (*domain.BucketSummary, error) {
	summary := &domain.BucketSummary{}
	for _, match := range m.matches {
		if match.RunID != runID {
			continue
		}
		switch match.MatchType {
		case domain.MatchTypeDirect:
			summary.Matched++
		case domain.MatchTypeProbable:
			summary.Mismatch++
		case domain.MatchTypePartial:
			summary.Partial++
		case domain.MatchTypeMissing2B:
			summary.Missing2B++
		case domain.MatchTypeMissingPR:
			summary.MissingPR++
		case domain.MatchTypeDuplicate:
			summary.Duplicate++
		}
	}
	return summary, nil
}

func (m *mockStore) CreateIMSAction(_ context.Context, tenantID uuid.UUID, action *domain.IMSAction) error {
	action.ID = uuid.New()
	action.TenantID = tenantID
	action.CreatedAt = time.Now().UTC()
	m.actions = append(m.actions, *action)
	return nil
}

func (m *mockStore) ListIMSActions(_ context.Context, _ uuid.UUID, gstin, returnPeriod string) ([]domain.IMSAction, error) {
	var results []domain.IMSAction
	for _, a := range m.actions {
		if a.GSTIN == gstin && a.ReturnPeriod == returnPeriod {
			results = append(results, a)
		}
	}
	return results, nil
}

// failingStore wraps mockStore and selectively returns errors.
type failingStore struct {
	*mockStore
	failCreateRun        bool
	failBulkInsertMatches bool
}

func (f *failingStore) CreateRun(ctx context.Context, tenantID uuid.UUID, run *domain.ReconRun) error {
	if f.failCreateRun {
		return fmt.Errorf("db error: create run")
	}
	return f.mockStore.CreateRun(ctx, tenantID, run)
}

func (f *failingStore) BulkInsertMatches(ctx context.Context, tenantID uuid.UUID, matches []domain.ReconMatch) error {
	if f.failBulkInsertMatches {
		return fmt.Errorf("db error: bulk insert")
	}
	return f.mockStore.BulkInsertMatches(ctx, tenantID, matches)
}

// ---------------------------------------------------------------------------
// Tests: Health
// ---------------------------------------------------------------------------

func TestHealth(t *testing.T) {
	h := NewHandlers(newMockStore(), nil, nil)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	h.Health(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "recon-service")
}

// ---------------------------------------------------------------------------
// Tests: RunRecon — missing tenant
// ---------------------------------------------------------------------------

func TestRunRecon_MissingTenantID(t *testing.T) {
	h := NewHandlers(newMockStore(), nil, nil)
	body, _ := json.Marshal(domain.RunReconRequest{GSTIN: "27AABCU9603R1ZP", ReturnPeriod: "012024"})
	req := httptest.NewRequest(http.MethodPost, "/v1/recon/run", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.RunRecon(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: RunRecon — missing required fields
// ---------------------------------------------------------------------------

func TestRunRecon_MissingFields(t *testing.T) {
	h := NewHandlers(newMockStore(), nil, nil)
	body, _ := json.Marshal(domain.RunReconRequest{})
	req := httptest.NewRequest(http.MethodPost, "/v1/recon/run", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.RunRecon(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: RunRecon — invalid body
// ---------------------------------------------------------------------------

func TestRunRecon_InvalidBody(t *testing.T) {
	h := NewHandlers(newMockStore(), nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/v1/recon/run", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.RunRecon(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: AcceptMatch — missing user
// ---------------------------------------------------------------------------

func TestAcceptMatch_MissingUser(t *testing.T) {
	ms := newMockStore()
	matchID := uuid.New()
	ms.matches[matchID] = &domain.ReconMatch{
		ID:        matchID,
		MatchType: domain.MatchTypeDirect,
		Status:    domain.MatchStatusUnreviewed,
		PRAmount:  decimal.Zero,
		GSTR2BAmount: decimal.Zero,
	}
	h := NewHandlers(ms, nil, nil)

	router := NewRouter(h)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/v1/recon/matches/%s/accept", matchID.String()), nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: AcceptMatch — success
// ---------------------------------------------------------------------------

func TestAcceptMatch_Success(t *testing.T) {
	ms := newMockStore()
	matchID := uuid.New()
	ms.matches[matchID] = &domain.ReconMatch{
		ID:        matchID,
		MatchType: domain.MatchTypeDirect,
		Status:    domain.MatchStatusUnreviewed,
		PRAmount:  decimal.Zero,
		GSTR2BAmount: decimal.Zero,
	}
	h := NewHandlers(ms, nil, nil)

	router := NewRouter(h)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/v1/recon/matches/%s/accept", matchID.String()), nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	req.Header.Set("X-User-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, domain.MatchStatusAccepted, ms.matches[matchID].Status)
}

// ---------------------------------------------------------------------------
// Tests: BulkAccept — empty IDs
// ---------------------------------------------------------------------------

func TestBulkAccept_EmptyIDs(t *testing.T) {
	h := NewHandlers(newMockStore(), nil, nil)
	body, _ := json.Marshal(domain.BulkAcceptRequest{MatchIDs: []uuid.UUID{}})
	req := httptest.NewRequest(http.MethodPost, "/v1/recon/matches/bulk-accept", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	req.Header.Set("X-User-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.BulkAccept(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: ListMatches — missing run_id
// ---------------------------------------------------------------------------

func TestListMatches_MissingRunID(t *testing.T) {
	h := NewHandlers(newMockStore(), nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/recon/matches", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.ListMatches(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: GetRun — not found
// ---------------------------------------------------------------------------

func TestGetRun_NotFound(t *testing.T) {
	h := NewHandlers(newMockStore(), nil, nil)
	router := NewRouter(h)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/recon/run/%s", uuid.New().String()), nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: Router
// ---------------------------------------------------------------------------

func TestNewRouter(t *testing.T) {
	h := NewHandlers(newMockStore(), nil, nil)
	r := NewRouter(h)
	require.NotNil(t, r)

	// Health endpoint
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Ping (heartbeat) endpoint
	req = httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: IMSActionHandler — missing fields
// ---------------------------------------------------------------------------

func TestIMSAction_MissingFields(t *testing.T) {
	h := NewHandlers(newMockStore(), nil, nil)
	body, _ := json.Marshal(domain.IMSActionRequest{})
	req := httptest.NewRequest(http.MethodPost, "/v1/recon/ims/action", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.IMSActionHandler(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: GetIMSState — missing params
// ---------------------------------------------------------------------------

func TestGetIMSState_MissingParams(t *testing.T) {
	h := NewHandlers(newMockStore(), nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/recon/ims", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.GetIMSState(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: RunRecon — full success
// ---------------------------------------------------------------------------

func TestRunRecon_Success(t *testing.T) {
	// Mock Apex server (expects {"data": {"data": {"invoices": [...]}}})
	apexSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"data": map[string]interface{}{
					"invoices": []map[string]interface{}{
						{"id": "ap-1", "invoice_number": "INV-001", "invoice_date": "01/04/2026", "vendor_gstin": "27AABCB0001B1Z5", "taxable_value": 100000, "igst_amount": 18000, "total_amount": 118000, "hsn": "9988", "place_of_supply": "27"},
						{"id": "ap-2", "invoice_number": "INV-002", "invoice_date": "05/04/2026", "vendor_gstin": "27AABCB0002B1Z5", "taxable_value": 200000, "igst_amount": 36000, "total_amount": 236000, "hsn": "9988", "place_of_supply": "27"},
					},
					"total": 2,
				},
			},
		})
	}))
	defer apexSrv.Close()

	// Mock GSTN server (expects {"data": {"invoices": [...], "status": "..."}})
	gstnSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"gstin":       "29AABCA1234A1Z5",
				"ret_period":  "042026",
				"status":      "success",
				"total_count": 2,
				"invoices": []map[string]interface{}{
					{"supplier_gstin": "27AABCB0001B1Z5", "invoice_number": "INV-001", "invoice_date": "01/04/2026", "taxable_value": 100000, "igst_amount": 18000, "total_value": 118000, "hsn": "9988", "place_of_supply": "27"},
					{"supplier_gstin": "27AABCB0003B1Z5", "invoice_number": "INV-003", "invoice_date": "10/04/2026", "taxable_value": 50000, "igst_amount": 9000, "total_value": 59000, "hsn": "9988", "place_of_supply": "27"},
				},
			},
		})
	}))
	defer gstnSrv.Close()

	ms := newMockStore()
	apex := gateway.NewApexClient(apexSrv.URL)
	gstn := gateway.NewGSTNClient(gstnSrv.URL)
	h := NewHandlers(ms, apex, gstn)

	tenantID := uuid.New()
	body, _ := json.Marshal(domain.RunReconRequest{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026"})
	req := httptest.NewRequest(http.MethodPost, "/v1/recon/run", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()
	h.RunRecon(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	// Response is {"data": {"run_id": "...", "status": "..."}}
	var wrapper struct {
		Data domain.RunReconResponse `json:"data"`
	}
	json.NewDecoder(rec.Body).Decode(&wrapper)
	assert.Equal(t, "COMPLETED", wrapper.Data.Status)
	assert.NotEqual(t, uuid.Nil, wrapper.Data.RunID)

	// Verify matches were stored
	assert.True(t, len(ms.matches) > 0, "should have stored matches")
}

// ---------------------------------------------------------------------------
// Tests: RunRecon — CreateRun DB failure
// ---------------------------------------------------------------------------

func TestRunRecon_CreateRunFails(t *testing.T) {
	fs := &failingStore{mockStore: newMockStore(), failCreateRun: true}
	apex := gateway.NewApexClient("http://unused")
	gstn := gateway.NewGSTNClient("http://unused")
	h := NewHandlers(fs, apex, gstn)

	body, _ := json.Marshal(domain.RunReconRequest{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026"})
	req := httptest.NewRequest(http.MethodPost, "/v1/recon/run", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.RunRecon(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: RunRecon — Apex gateway failure
// ---------------------------------------------------------------------------

func TestRunRecon_ApexFails(t *testing.T) {
	apexSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer apexSrv.Close()

	ms := newMockStore()
	apex := gateway.NewApexClient(apexSrv.URL)
	gstn := gateway.NewGSTNClient("http://unused")
	h := NewHandlers(ms, apex, gstn)

	body, _ := json.Marshal(domain.RunReconRequest{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026"})
	req := httptest.NewRequest(http.MethodPost, "/v1/recon/run", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.RunRecon(rec, req)
	assert.Equal(t, http.StatusBadGateway, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: RunRecon — GSTN gateway failure
// ---------------------------------------------------------------------------

func TestRunRecon_GSTNFails(t *testing.T) {
	apexSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"data": map[string]interface{}{
					"invoices": []map[string]interface{}{
						{"id": "ap-1", "invoice_number": "INV-001", "invoice_date": "01/04/2026", "vendor_gstin": "27AABCB0001B1Z5", "taxable_value": 100000, "igst_amount": 18000, "total_amount": 118000, "hsn": "9988", "place_of_supply": "27"},
					},
					"total": 1,
				},
			},
		})
	}))
	defer apexSrv.Close()

	gstnSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer gstnSrv.Close()

	ms := newMockStore()
	apex := gateway.NewApexClient(apexSrv.URL)
	gstn := gateway.NewGSTNClient(gstnSrv.URL)
	h := NewHandlers(ms, apex, gstn)

	body, _ := json.Marshal(domain.RunReconRequest{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026"})
	req := httptest.NewRequest(http.MethodPost, "/v1/recon/run", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.RunRecon(rec, req)
	assert.Equal(t, http.StatusBadGateway, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: RunRecon — BulkInsertMatches failure
// ---------------------------------------------------------------------------

func TestRunRecon_BulkInsertFails(t *testing.T) {
	apexSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"data": map[string]interface{}{
					"invoices": []map[string]interface{}{
						{"id": "ap-1", "invoice_number": "INV-001", "invoice_date": "01/04/2026", "vendor_gstin": "27AABCB0001B1Z5", "taxable_value": 100000, "igst_amount": 18000, "total_amount": 118000, "hsn": "9988", "place_of_supply": "27"},
					},
					"total": 1,
				},
			},
		})
	}))
	defer apexSrv.Close()

	gstnSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"gstin":       "29AABCA1234A1Z5",
				"ret_period":  "042026",
				"status":      "success",
				"total_count": 1,
				"invoices": []map[string]interface{}{
					{"supplier_gstin": "27AABCB0001B1Z5", "invoice_number": "INV-001", "invoice_date": "01/04/2026", "taxable_value": 100000, "igst_amount": 18000, "total_value": 118000, "hsn": "9988", "place_of_supply": "27"},
				},
			},
		})
	}))
	defer gstnSrv.Close()

	fs := &failingStore{mockStore: newMockStore(), failBulkInsertMatches: true}
	apex := gateway.NewApexClient(apexSrv.URL)
	gstn := gateway.NewGSTNClient(gstnSrv.URL)
	h := NewHandlers(fs, apex, gstn)

	body, _ := json.Marshal(domain.RunReconRequest{GSTIN: "29AABCA1234A1Z5", ReturnPeriod: "042026"})
	req := httptest.NewRequest(http.MethodPost, "/v1/recon/run", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.RunRecon(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: IMSAction — success via GSTN gateway
// ---------------------------------------------------------------------------

func TestIMSAction_Success(t *testing.T) {
	gstnSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"invoice_id": "inv-001",
				"action":     "ACCEPT",
				"status":     "SUCCESS",
			},
		})
	}))
	defer gstnSrv.Close()

	ms := newMockStore()
	gstn := gateway.NewGSTNClient(gstnSrv.URL)
	h := NewHandlers(ms, nil, gstn)
	router := NewRouter(h)

	actionBody, _ := json.Marshal(map[string]string{
		"invoice_id": "inv-001",
		"action":     "ACCEPT",
		"reason":     "verified",
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/recon/ims/action?gstin=27TEST&return_period=042026", bytes.NewReader(actionBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	req.Header.Set("X-User-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var wrapper struct {
		Data struct {
			Status string `json:"status"`
		} `json:"data"`
	}
	json.NewDecoder(rec.Body).Decode(&wrapper)
	assert.Equal(t, "SUCCESS", wrapper.Data.Status)
}

// ---------------------------------------------------------------------------
// Tests: GetIMSState — success
// ---------------------------------------------------------------------------

func TestGetIMSState_Success(t *testing.T) {
	gstnSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"invoices": []map[string]interface{}{
					{"invoice_id": "inv-001", "action": "ACCEPT"},
				},
			},
		})
	}))
	defer gstnSrv.Close()

	ms := newMockStore()
	gstn := gateway.NewGSTNClient(gstnSrv.URL)
	h := NewHandlers(ms, nil, gstn)

	req := httptest.NewRequest(http.MethodGet, "/v1/recon/ims?gstin=27TEST&return_period=042026", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.GetIMSState(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: GetRun — success
// ---------------------------------------------------------------------------

func TestGetRun_Success(t *testing.T) {
	ms := newMockStore()
	run := &domain.ReconRun{
		GSTIN:        "29AABCA1234A1Z5",
		ReturnPeriod: "042026",
		Status:       "COMPLETED",
		StartedAt:    time.Now().UTC(),
	}
	_ = ms.CreateRun(context.Background(), uuid.New(), run)

	h := NewHandlers(ms, nil, nil)
	router := NewRouter(h)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/recon/run/%s", run.ID.String()), nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: GetRun — missing tenant
// ---------------------------------------------------------------------------

func TestGetRun_MissingTenant(t *testing.T) {
	h := NewHandlers(newMockStore(), nil, nil)
	router := NewRouter(h)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/recon/run/%s", uuid.New().String()), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: ListMatches — success
// ---------------------------------------------------------------------------

func TestListMatches_Success(t *testing.T) {
	ms := newMockStore()
	runID := uuid.New()
	matchID := uuid.New()
	ms.matches[matchID] = &domain.ReconMatch{
		ID:           matchID,
		RunID:        runID,
		MatchType:    domain.MatchTypeDirect,
		Status:       domain.MatchStatusUnreviewed,
		PRAmount:     decimal.NewFromInt(100000),
		GSTR2BAmount: decimal.NewFromInt(100000),
	}
	h := NewHandlers(ms, nil, nil)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/recon/matches?run_id=%s", runID.String()), nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.ListMatches(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: ListMatches — invalid run_id
// ---------------------------------------------------------------------------

func TestListMatches_InvalidRunID(t *testing.T) {
	h := NewHandlers(newMockStore(), nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/recon/matches?run_id=bad", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.ListMatches(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: ListMatches — with filters
// ---------------------------------------------------------------------------

func TestListMatches_WithFilters(t *testing.T) {
	ms := newMockStore()
	runID := uuid.New()
	ms.matches[uuid.New()] = &domain.ReconMatch{
		ID: uuid.New(), RunID: runID, MatchType: domain.MatchTypeDirect, Status: domain.MatchStatusUnreviewed,
		PRAmount: decimal.Zero, GSTR2BAmount: decimal.Zero,
	}
	h := NewHandlers(ms, nil, nil)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/recon/matches?run_id=%s&match_type=DIRECT&status=UNREVIEWED&limit=10&offset=0", runID.String()), nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.ListMatches(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: BulkAccept — success
// ---------------------------------------------------------------------------

func TestBulkAccept_Success(t *testing.T) {
	ms := newMockStore()
	m1, m2 := uuid.New(), uuid.New()
	ms.matches[m1] = &domain.ReconMatch{ID: m1, MatchType: domain.MatchTypeDirect, Status: domain.MatchStatusUnreviewed, PRAmount: decimal.Zero, GSTR2BAmount: decimal.Zero}
	ms.matches[m2] = &domain.ReconMatch{ID: m2, MatchType: domain.MatchTypeProbable, Status: domain.MatchStatusUnreviewed, PRAmount: decimal.Zero, GSTR2BAmount: decimal.Zero}
	h := NewHandlers(ms, nil, nil)

	body, _ := json.Marshal(domain.BulkAcceptRequest{MatchIDs: []uuid.UUID{m1, m2}})
	req := httptest.NewRequest(http.MethodPost, "/v1/recon/matches/bulk-accept", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	req.Header.Set("X-User-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.BulkAccept(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, domain.MatchStatusAccepted, ms.matches[m1].Status)
	assert.Equal(t, domain.MatchStatusAccepted, ms.matches[m2].Status)
}

// ---------------------------------------------------------------------------
// Tests: BulkAccept — missing user
// ---------------------------------------------------------------------------

func TestBulkAccept_MissingUser(t *testing.T) {
	h := NewHandlers(newMockStore(), nil, nil)
	body, _ := json.Marshal(domain.BulkAcceptRequest{MatchIDs: []uuid.UUID{uuid.New()}})
	req := httptest.NewRequest(http.MethodPost, "/v1/recon/matches/bulk-accept", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.BulkAccept(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: BulkAccept — invalid body
// ---------------------------------------------------------------------------

func TestBulkAccept_InvalidBody(t *testing.T) {
	h := NewHandlers(newMockStore(), nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/v1/recon/matches/bulk-accept", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	req.Header.Set("X-User-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.BulkAccept(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: IMSAction — invalid action value
// ---------------------------------------------------------------------------

func TestIMSAction_InvalidAction(t *testing.T) {
	h := NewHandlers(newMockStore(), nil, nil)
	body, _ := json.Marshal(domain.IMSActionRequest{InvoiceID: "inv-1", Action: "INVALID"})
	req := httptest.NewRequest(http.MethodPost, "/v1/recon/ims/action?gstin=X&return_period=042026", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.IMSActionHandler(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: IMSAction — missing query params
// ---------------------------------------------------------------------------

func TestIMSAction_MissingQueryParams(t *testing.T) {
	h := NewHandlers(newMockStore(), nil, nil)
	body, _ := json.Marshal(domain.IMSActionRequest{InvoiceID: "inv-1", Action: "ACCEPT"})
	req := httptest.NewRequest(http.MethodPost, "/v1/recon/ims/action", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.IMSActionHandler(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: IMSAction — invalid body
// ---------------------------------------------------------------------------

func TestIMSAction_InvalidBody(t *testing.T) {
	h := NewHandlers(newMockStore(), nil, nil)
	req := httptest.NewRequest(http.MethodPost, "/v1/recon/ims/action?gstin=X&return_period=042026", bytes.NewReader([]byte("bad")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	h.IMSActionHandler(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: GetIMSState — missing tenant
// ---------------------------------------------------------------------------

func TestGetIMSState_MissingTenant(t *testing.T) {
	h := NewHandlers(newMockStore(), nil, nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/recon/ims?gstin=X&return_period=042026", nil)
	rec := httptest.NewRecorder()
	h.GetIMSState(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: AcceptMatch — invalid match_id
// ---------------------------------------------------------------------------

func TestAcceptMatch_InvalidMatchID(t *testing.T) {
	h := NewHandlers(newMockStore(), nil, nil)
	router := NewRouter(h)
	req := httptest.NewRequest(http.MethodPost, "/v1/recon/matches/bad-id/accept", nil)
	req.Header.Set("X-Tenant-Id", uuid.New().String())
	req.Header.Set("X-User-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: AcceptMatch — missing tenant
// ---------------------------------------------------------------------------

func TestAcceptMatch_MissingTenant(t *testing.T) {
	h := NewHandlers(newMockStore(), nil, nil)
	router := NewRouter(h)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/v1/recon/matches/%s/accept", uuid.New().String()), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: IMS Round-Trip (mock gateway server)
// ---------------------------------------------------------------------------

func TestIMSRoundTrip(t *testing.T) {
	// 1. Set up a mock GSTN gateway that responds to IMS action + IMS get
	imsState := map[string]string{} // invoiceID -> action
	mockGW := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/gateway/adaequare/ims/action":
			var req struct {
				InvoiceID string `json:"invoice_id"`
				Action    string `json:"action"`
			}
			json.NewDecoder(r.Body).Decode(&req)
			imsState[req.InvoiceID] = req.Action
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]string{
					"invoice_id": req.InvoiceID,
					"action":     req.Action,
					"status":     "SUCCESS",
				},
			})
		case "/v1/gateway/adaequare/ims/get":
			invoices := []map[string]interface{}{}
			for id, action := range imsState {
				invoices = append(invoices, map[string]interface{}{
					"invoice_id": id, "action": action, "supplier_gstin": "29AAA",
					"invoice_number": "INV-001", "total_value": 10000,
				})
			}
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"gstin": "27TEST", "ret_period": "042026",
					"invoices": invoices, "total_count": len(invoices), "status": "SUCCESS",
				},
			})
		default:
			w.WriteHeader(404)
		}
	}))
	defer mockGW.Close()

	gstnClient := gateway.NewGSTNClient(mockGW.URL)
	ms := newMockStore()
	h := NewHandlers(ms, nil, gstnClient)
	router := NewRouter(h)
	tenantID := uuid.New()

	// 2. Send IMS ACCEPT action via recon-service API
	actionBody, _ := json.Marshal(map[string]string{
		"invoice_id": "inv-test-001",
		"action":     "ACCEPT",
		"reason":     "verified",
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/recon/ims/action?gstin=27TEST&return_period=042026", bytes.NewReader(actionBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.Header.Set("X-User-Id", uuid.New().String())
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	// httputil.JSON wraps response as {"data": ...}
	var actionWrapper struct {
		Data struct {
			InvoiceID string `json:"invoice_id"`
			Action    string `json:"action"`
			Status    string `json:"status"`
		} `json:"data"`
	}
	json.NewDecoder(rec.Body).Decode(&actionWrapper)
	assert.Equal(t, "ACCEPT", actionWrapper.Data.Action)
	assert.Equal(t, "SUCCESS", actionWrapper.Data.Status)

	// 3. Verify action was recorded in local store
	actions, err := ms.ListIMSActions(context.Background(), tenantID, "27TEST", "042026")
	require.NoError(t, err)
	require.Len(t, actions, 1)
	assert.Equal(t, "ACCEPT", actions[0].Action)
	assert.Equal(t, "inv-test-001", actions[0].InvoiceID)

	// 4. Fetch IMS state — should reflect the accepted invoice
	req2 := httptest.NewRequest(http.MethodGet, "/v1/recon/ims?gstin=27TEST&return_period=042026", nil)
	req2.Header.Set("X-Tenant-Id", tenantID.String())
	rec2 := httptest.NewRecorder()
	router.ServeHTTP(rec2, req2)

	require.Equal(t, http.StatusOK, rec2.Code)
	var imsWrapper struct {
		Data struct {
			Invoices []struct {
				InvoiceID string `json:"invoice_id"`
				Action    string `json:"action"`
			} `json:"invoices"`
		} `json:"data"`
	}
	json.NewDecoder(rec2.Body).Decode(&imsWrapper)
	require.Len(t, imsWrapper.Data.Invoices, 1)
	assert.Equal(t, "ACCEPT", imsWrapper.Data.Invoices[0].Action)
	assert.Equal(t, "inv-test-001", imsWrapper.Data.Invoices[0].InvoiceID)
}
