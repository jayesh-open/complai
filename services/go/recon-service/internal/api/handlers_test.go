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
