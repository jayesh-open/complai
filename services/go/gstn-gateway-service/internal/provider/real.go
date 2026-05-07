package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/services/go/gstn-gateway-service/internal/domain"
)

var _ GSTNProvider = (*AdaequareProvider)(nil)

// TaxpayerCredentials holds per-GSTIN credentials for GSTN/IRP/EWB portals.
// In production these come from AWS Secrets Manager per tenant+GSTIN.
type TaxpayerCredentials struct {
	Username string
	Password string
}

type AdaequareProvider struct {
	baseURL      string
	clientID     string
	clientSecret string
	httpClient   *http.Client

	// Per-GSTIN taxpayer credentials (for testing: single set via env vars).
	// In production, replaced by a credential resolver.
	defaultCreds *TaxpayerCredentials

	mu          sync.RWMutex
	accessToken string
	tokenExpiry time.Time
}

type AdaequareOption func(*AdaequareProvider)

func WithTaxpayerCredentials(username, password string) AdaequareOption {
	return func(p *AdaequareProvider) {
		p.defaultCreds = &TaxpayerCredentials{Username: username, Password: password}
	}
}

func NewAdaequareProvider(baseURL, clientID, clientSecret string, opts ...AdaequareOption) *AdaequareProvider {
	p := &AdaequareProvider{
		baseURL:      strings.TrimRight(baseURL, "/"),
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, o := range opts {
		o(p)
	}
	return p
}

func (a *AdaequareProvider) Authenticate(ctx context.Context) (*domain.AuthResponse, error) {
	token, resp, err := a.authenticate(ctx)
	if err != nil {
		return nil, err
	}
	_ = token
	return resp, nil
}

func (a *AdaequareProvider) authenticate(ctx context.Context) (string, *domain.AuthResponse, error) {
	a.mu.RLock()
	if a.accessToken != "" && time.Now().Before(a.tokenExpiry) {
		token := a.accessToken
		a.mu.RUnlock()
		return token, &domain.AuthResponse{
			AccessToken: token,
			TokenType:   "bearer",
			ExpiresIn:   int(time.Until(a.tokenExpiry).Seconds()),
			Scope:       "gsp",
		}, nil
	}
	a.mu.RUnlock()

	a.mu.Lock()
	defer a.mu.Unlock()

	// Double-check after acquiring write lock
	if a.accessToken != "" && time.Now().Before(a.tokenExpiry) {
		return a.accessToken, &domain.AuthResponse{
			AccessToken: a.accessToken,
			TokenType:   "bearer",
			ExpiresIn:   int(time.Until(a.tokenExpiry).Seconds()),
			Scope:       "gsp",
		}, nil
	}

	authURL := a.gspURL("/gsp/authenticate?action=GSP&grant_type=token")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, authURL, nil)
	if err != nil {
		return "", nil, fmt.Errorf("adaequare auth: build request: %w", err)
	}
	req.Header.Set("gspappid", a.clientID)
	req.Header.Set("gspappsecret", a.clientSecret)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("adaequare auth: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("adaequare auth: read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("adaequare auth: HTTP %d: %s", resp.StatusCode, string(body))
	}

	var authResp domain.AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return "", nil, fmt.Errorf("adaequare auth: decode: %w", err)
	}

	if authResp.AccessToken == "" {
		return "", nil, fmt.Errorf("adaequare auth: empty access_token in response: %s", string(body))
	}

	// Cache with 1-hour safety margin
	a.accessToken = authResp.AccessToken
	a.tokenExpiry = time.Now().Add(time.Duration(authResp.ExpiresIn)*time.Second - time.Hour)

	log.Info().
		Int("expires_in", authResp.ExpiresIn).
		Str("jti", authResp.JTI).
		Msg("adaequare: authenticated")

	return a.accessToken, &authResp, nil
}

// gspURL builds the auth URL (no /test prefix — auth is environment-agnostic).
func (a *AdaequareProvider) gspURL(path string) string {
	base := a.baseURL
	if idx := strings.Index(base, "/test"); idx != -1 {
		base = base[:idx]
	}
	return base + path
}

// apiURL builds a data endpoint URL (uses full baseURL including /test for sandbox).
func (a *AdaequareProvider) apiURL(path string) string {
	return a.baseURL + path
}

type adaequareError struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Status    int    `json:"status"`
	ErrorCode string `json:"errorCode"`
}

func (a *AdaequareProvider) doJSON(ctx context.Context, method, url string, reqBody interface{}, gstin string) ([]byte, error) {
	token, _, err := a.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	var bodyReader io.Reader
	if reqBody != nil {
		b, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("adaequare: marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("adaequare: build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("requestid", uuid.New().String())
	if gstin != "" {
		req.Header.Set("gstin", gstin)
		// state-cd = first 2 chars of GSTIN (required by Returns/Common endpoints)
		if len(gstin) >= 2 {
			req.Header.Set("state-cd", gstin[:2])
		}
	}

	// GST Returns endpoints use "username" (no underscore)
	if a.defaultCreds != nil {
		req.Header.Set("username", a.defaultCreds.Username)
		req.Header.Set("password", a.defaultCreds.Password)
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("adaequare: %s %s: %w", method, url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("adaequare: read response: %w", err)
	}

	// Check for Adaequare-level errors (they return 200 with success:false)
	var ae adaequareError
	if json.Unmarshal(body, &ae) == nil && !ae.Success && ae.Message != "" {
		return nil, fmt.Errorf("adaequare: %s", ae.Message)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("adaequare: HTTP %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// --- GSTR-1 ---

func (a *AdaequareProvider) GSTR1Save(ctx context.Context, req *domain.GSTR1SaveRequest) (*domain.GSTR1SaveResponse, error) {
	url := a.apiURL("/enriched/returns?action=RETSAVE&gstin=" + req.GSTIN + "&ret_period=" + req.RetPeriod + "&rtn_typ=GSTR1")
	body, err := a.doJSON(ctx, http.MethodPost, url, req.Data, req.GSTIN)
	if err != nil {
		return nil, err
	}
	var resp domain.GSTR1SaveResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode gstr1 save: %w", err)
	}
	resp.RequestID = req.RequestID
	resp.SavedAt = time.Now()
	if resp.Status == "" {
		resp.Status = "saved"
	}
	return &resp, nil
}

func (a *AdaequareProvider) GSTR1Get(ctx context.Context, req *domain.GSTR1GetRequest) (*domain.GSTR1GetResponse, error) {
	url := a.apiURL("/enriched/returns?action=RETSUM&gstin=" + req.GSTIN + "&ret_period=" + req.RetPeriod + "&rtn_typ=GSTR1")
	if req.Section != "" {
		url += "&section=" + req.Section
	}
	body, err := a.doJSON(ctx, http.MethodGet, url, nil, req.GSTIN)
	if err != nil {
		return nil, err
	}
	var resp domain.GSTR1GetResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode gstr1 get: %w", err)
	}
	resp.RequestID = req.RequestID
	resp.GSTIN = req.GSTIN
	resp.RetPeriod = req.RetPeriod
	if resp.Status == "" {
		resp.Status = "success"
	}
	return &resp, nil
}

func (a *AdaequareProvider) GSTR1Reset(ctx context.Context, req *domain.GSTR1ResetRequest) (*domain.GSTR1ResetResponse, error) {
	url := a.apiURL("/enriched/returns?action=RESET&gstin=" + req.GSTIN + "&ret_period=" + req.RetPeriod + "&rtn_typ=GSTR1")
	body, err := a.doJSON(ctx, http.MethodPost, url, nil, req.GSTIN)
	if err != nil {
		return nil, err
	}
	var resp domain.GSTR1ResetResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode gstr1 reset: %w", err)
	}
	resp.RequestID = req.RequestID
	if resp.Status == "" {
		resp.Status = "success"
	}
	return &resp, nil
}

func (a *AdaequareProvider) GSTR1Submit(ctx context.Context, req *domain.GSTR1SubmitRequest) (*domain.GSTR1SubmitResponse, error) {
	url := a.apiURL("/enriched/returns?action=RETSUBMIT&gstin=" + req.GSTIN + "&ret_period=" + req.RetPeriod + "&rtn_typ=GSTR1")
	body, err := a.doJSON(ctx, http.MethodPost, url, nil, req.GSTIN)
	if err != nil {
		return nil, err
	}
	var resp domain.GSTR1SubmitResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode gstr1 submit: %w", err)
	}
	resp.RequestID = req.RequestID
	if resp.Status == "" {
		resp.Status = "submitted"
	}
	return &resp, nil
}

func (a *AdaequareProvider) GSTR1File(ctx context.Context, req *domain.GSTR1FileRequest) (*domain.GSTR1FileResponse, error) {
	url := a.apiURL("/enriched/returns?action=RETFILE&gstin=" + req.GSTIN + "&ret_period=" + req.RetPeriod + "&rtn_typ=GSTR1")
	payload := map[string]string{
		"sign_type": req.SignType,
		"pan":       req.PAN,
	}
	if req.EVOTP != "" {
		payload["ev_otp"] = req.EVOTP
	}
	body, err := a.doJSON(ctx, http.MethodPost, url, payload, req.GSTIN)
	if err != nil {
		return nil, err
	}
	var resp domain.GSTR1FileResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode gstr1 file: %w", err)
	}
	resp.RequestID = req.RequestID
	if resp.Status == "" {
		resp.Status = "filed"
	}
	if resp.FiledAt.IsZero() {
		resp.FiledAt = time.Now()
	}
	return &resp, nil
}

func (a *AdaequareProvider) GSTR1Status(ctx context.Context, req *domain.GSTR1StatusRequest) (*domain.GSTR1StatusResponse, error) {
	url := a.apiURL("/enriched/returns?action=FILEDET&gstin=" + req.GSTIN + "&ret_period=" + req.RetPeriod + "&rtn_typ=GSTR1")
	body, err := a.doJSON(ctx, http.MethodGet, url, nil, req.GSTIN)
	if err != nil {
		return nil, err
	}
	var resp domain.GSTR1StatusResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode gstr1 status: %w", err)
	}
	resp.RequestID = req.RequestID
	resp.GSTIN = req.GSTIN
	resp.RetPeriod = req.RetPeriod
	return &resp, nil
}

func (a *AdaequareProvider) GSTR1Summary(ctx context.Context, req *domain.GSTR1SummaryRequest) (*domain.GSTR1SummaryResponse, error) {
	url := a.apiURL("/enriched/returns?action=RETSUM&gstin=" + req.GSTIN + "&ret_period=" + req.RetPeriod + "&rtn_typ=GSTR1")
	body, err := a.doJSON(ctx, http.MethodGet, url, nil, req.GSTIN)
	if err != nil {
		return nil, err
	}
	var resp domain.GSTR1SummaryResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode gstr1 summary: %w", err)
	}
	resp.RequestID = req.RequestID
	resp.GSTIN = req.GSTIN
	resp.RetPeriod = req.RetPeriod
	if resp.Status == "" {
		resp.Status = "success"
	}
	return &resp, nil
}

// --- GSTR-2B ---

func (a *AdaequareProvider) GSTR2BGet(ctx context.Context, req *domain.GSTR2BGetRequest) (*domain.GSTR2BGetResponse, error) {
	url := a.apiURL("/enriched/returns?action=GET2B&gstin=" + req.GSTIN + "&ret_period=" + req.RetPeriod + "&rtn_typ=GSTR2B")
	body, err := a.doJSON(ctx, http.MethodGet, url, nil, req.GSTIN)
	if err != nil {
		return nil, err
	}
	var resp domain.GSTR2BGetResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode gstr2b: %w", err)
	}
	resp.RequestID = req.RequestID
	resp.GSTIN = req.GSTIN
	resp.RetPeriod = req.RetPeriod
	if resp.Status == "" {
		resp.Status = "success"
	}
	return &resp, nil
}

// --- GSTR-2A ---

func (a *AdaequareProvider) GSTR2AGet(ctx context.Context, req *domain.GSTR2AGetRequest) (*domain.GSTR2AGetResponse, error) {
	url := a.apiURL("/enriched/returns?action=GET2A&gstin=" + req.GSTIN + "&ret_period=" + req.RetPeriod + "&section=" + req.Section + "&rtn_typ=GSTR2A")
	body, err := a.doJSON(ctx, http.MethodGet, url, nil, req.GSTIN)
	if err != nil {
		return nil, err
	}
	var resp domain.GSTR2AGetResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode gstr2a: %w", err)
	}
	resp.RequestID = req.RequestID
	resp.GSTIN = req.GSTIN
	resp.RetPeriod = req.RetPeriod
	resp.Section = req.Section
	if resp.Status == "" {
		resp.Status = "success"
	}
	return &resp, nil
}

// --- IMS ---

func (a *AdaequareProvider) IMSGet(ctx context.Context, req *domain.IMSGetRequest) (*domain.IMSGetResponse, error) {
	url := a.apiURL("/enriched/returns?action=GETIMS&gstin=" + req.GSTIN + "&ret_period=" + req.RetPeriod + "&rtn_typ=IMS")
	body, err := a.doJSON(ctx, http.MethodGet, url, nil, req.GSTIN)
	if err != nil {
		return nil, err
	}
	var resp domain.IMSGetResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode ims: %w", err)
	}
	resp.RequestID = req.RequestID
	resp.GSTIN = req.GSTIN
	resp.RetPeriod = req.RetPeriod
	if resp.Status == "" {
		resp.Status = "success"
	}
	return &resp, nil
}

func (a *AdaequareProvider) IMSAction(ctx context.Context, req *domain.IMSActionRequest) (*domain.IMSActionResponse, error) {
	url := a.apiURL("/enriched/returns?action=IMSACTION&gstin=" + req.GSTIN + "&ret_period=" + req.RetPeriod + "&rtn_typ=IMS")
	payload := map[string]string{
		"invoice_id": req.InvoiceID,
		"action":     req.Action,
	}
	if req.Reason != "" {
		payload["reason"] = req.Reason
	}
	body, err := a.doJSON(ctx, http.MethodPost, url, payload, req.GSTIN)
	if err != nil {
		return nil, err
	}
	var resp domain.IMSActionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode ims action: %w", err)
	}
	resp.RequestID = req.RequestID
	if resp.Status == "" {
		resp.Status = "success"
	}
	return &resp, nil
}

func (a *AdaequareProvider) IMSBulkAction(ctx context.Context, req *domain.IMSBulkActionRequest) (*domain.IMSBulkActionResponse, error) {
	url := a.apiURL("/enriched/returns?action=IMSBULK&gstin=" + req.GSTIN + "&ret_period=" + req.RetPeriod + "&rtn_typ=IMS")
	payload := map[string]interface{}{
		"invoice_ids": req.InvoiceIDs,
		"action":      req.Action,
	}
	body, err := a.doJSON(ctx, http.MethodPost, url, payload, req.GSTIN)
	if err != nil {
		return nil, err
	}
	var resp domain.IMSBulkActionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode ims bulk: %w", err)
	}
	resp.RequestID = req.RequestID
	if resp.Status == "" {
		resp.Status = "success"
	}
	return &resp, nil
}

// --- GSTR-3B ---

func (a *AdaequareProvider) GSTR3BSave(ctx context.Context, req *domain.GSTR3BSaveRequest) (*domain.GSTR3BSaveResponse, error) {
	url := a.apiURL("/enriched/returns?action=RETSAVE&gstin=" + req.GSTIN + "&ret_period=" + req.RetPeriod + "&rtn_typ=GSTR3B")
	body, err := a.doJSON(ctx, http.MethodPost, url, req.Data, req.GSTIN)
	if err != nil {
		return nil, err
	}
	var resp domain.GSTR3BSaveResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode gstr3b save: %w", err)
	}
	resp.RequestID = req.RequestID
	if resp.Status == "" {
		resp.Status = "saved"
	}
	return &resp, nil
}

func (a *AdaequareProvider) GSTR3BSubmit(ctx context.Context, req *domain.GSTR3BSubmitRequest) (*domain.GSTR3BSubmitResponse, error) {
	url := a.apiURL("/enriched/returns?action=RETSUBMIT&gstin=" + req.GSTIN + "&ret_period=" + req.RetPeriod + "&rtn_typ=GSTR3B")
	body, err := a.doJSON(ctx, http.MethodPost, url, nil, req.GSTIN)
	if err != nil {
		return nil, err
	}
	var resp domain.GSTR3BSubmitResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode gstr3b submit: %w", err)
	}
	resp.RequestID = req.RequestID
	if resp.Status == "" {
		resp.Status = "submitted"
	}
	return &resp, nil
}

func (a *AdaequareProvider) GSTR3BFile(ctx context.Context, req *domain.GSTR3BFileRequest) (*domain.GSTR3BFileResponse, error) {
	url := a.apiURL("/enriched/returns?action=RETFILE&gstin=" + req.GSTIN + "&ret_period=" + req.RetPeriod + "&rtn_typ=GSTR3B")
	payload := map[string]string{
		"sign_type": req.SignType,
		"pan":       req.PAN,
	}
	if req.EVOTP != "" {
		payload["ev_otp"] = req.EVOTP
	}
	body, err := a.doJSON(ctx, http.MethodPost, url, payload, req.GSTIN)
	if err != nil {
		return nil, err
	}
	var resp domain.GSTR3BFileResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode gstr3b file: %w", err)
	}
	resp.RequestID = req.RequestID
	if resp.Status == "" {
		resp.Status = "filed"
	}
	if resp.FiledAt.IsZero() {
		resp.FiledAt = time.Now()
	}
	return &resp, nil
}

// --- GSTR-9 ---

func (a *AdaequareProvider) GSTR9Save(ctx context.Context, req *domain.GSTR9SaveRequest) (*domain.GSTR9SaveResponse, error) {
	url := a.apiURL("/enriched/returns?action=RETSAVE&gstin=" + req.GSTIN + "&ret_period=" + req.FinancialYear + "&rtn_typ=GSTR9")
	body, err := a.doJSON(ctx, http.MethodPost, url, req.Data, req.GSTIN)
	if err != nil {
		return nil, err
	}
	var resp domain.GSTR9SaveResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode gstr9 save: %w", err)
	}
	resp.RequestID = req.RequestID
	if resp.SavedAt.IsZero() {
		resp.SavedAt = time.Now()
	}
	if resp.Status == "" {
		resp.Status = "saved"
	}
	return &resp, nil
}

func (a *AdaequareProvider) GSTR9Submit(ctx context.Context, req *domain.GSTR9SubmitRequest) (*domain.GSTR9SubmitResponse, error) {
	url := a.apiURL("/enriched/returns?action=RETSUBMIT&gstin=" + req.GSTIN + "&ret_period=" + req.FinancialYear + "&rtn_typ=GSTR9")
	body, err := a.doJSON(ctx, http.MethodPost, url, nil, req.GSTIN)
	if err != nil {
		return nil, err
	}
	var resp domain.GSTR9SubmitResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode gstr9 submit: %w", err)
	}
	resp.RequestID = req.RequestID
	if resp.Status == "" {
		resp.Status = "submitted"
	}
	return &resp, nil
}

func (a *AdaequareProvider) GSTR9File(ctx context.Context, req *domain.GSTR9FileRequest) (*domain.GSTR9FileResponse, error) {
	url := a.apiURL("/enriched/returns?action=RETFILE&gstin=" + req.GSTIN + "&ret_period=" + req.FinancialYear + "&rtn_typ=GSTR9")
	payload := map[string]string{
		"sign_type": req.SignType,
		"pan":       req.PAN,
	}
	if req.EVOTP != "" {
		payload["ev_otp"] = req.EVOTP
	}
	body, err := a.doJSON(ctx, http.MethodPost, url, payload, req.GSTIN)
	if err != nil {
		return nil, err
	}
	var resp domain.GSTR9FileResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode gstr9 file: %w", err)
	}
	resp.RequestID = req.RequestID
	if resp.Status == "" {
		resp.Status = "filed"
	}
	if resp.FiledAt.IsZero() {
		resp.FiledAt = time.Now()
	}
	return &resp, nil
}

func (a *AdaequareProvider) GSTR9Status(ctx context.Context, req *domain.GSTR9StatusRequest) (*domain.GSTR9StatusResponse, error) {
	url := a.apiURL("/enriched/returns?action=FILEDET&reference=" + req.Reference + "&rtn_typ=GSTR9")
	body, err := a.doJSON(ctx, http.MethodGet, url, nil, "")
	if err != nil {
		return nil, err
	}
	var resp domain.GSTR9StatusResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode gstr9 status: %w", err)
	}
	resp.RequestID = req.RequestID
	return &resp, nil
}

// --- GSTR-9C ---

func (a *AdaequareProvider) GSTR9CSave(ctx context.Context, req *domain.GSTR9CSaveRequest) (*domain.GSTR9CSaveResponse, error) {
	url := a.apiURL("/enriched/returns?action=RETSAVE&gstin=" + req.GSTIN + "&ret_period=" + req.FinancialYear + "&rtn_typ=GSTR9C")
	body, err := a.doJSON(ctx, http.MethodPost, url, req.Data, req.GSTIN)
	if err != nil {
		return nil, err
	}
	var resp domain.GSTR9CSaveResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode gstr9c save: %w", err)
	}
	resp.RequestID = req.RequestID
	if resp.SavedAt.IsZero() {
		resp.SavedAt = time.Now()
	}
	if resp.Status == "" {
		resp.Status = "saved"
	}
	return &resp, nil
}

func (a *AdaequareProvider) GSTR9CFile(ctx context.Context, req *domain.GSTR9CFileRequest) (*domain.GSTR9CFileResponse, error) {
	url := a.apiURL("/enriched/returns?action=RETFILE&gstin=" + req.GSTIN + "&ret_period=" + req.FinancialYear + "&rtn_typ=GSTR9C")
	payload := map[string]string{
		"pan": req.PAN,
	}
	body, err := a.doJSON(ctx, http.MethodPost, url, payload, req.GSTIN)
	if err != nil {
		return nil, err
	}
	var resp domain.GSTR9CFileResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode gstr9c file: %w", err)
	}
	resp.RequestID = req.RequestID
	if resp.Status == "" {
		resp.Status = "filed"
	}
	if resp.FiledAt.IsZero() {
		resp.FiledAt = time.Now()
	}
	return &resp, nil
}

func (a *AdaequareProvider) GSTR9CStatus(ctx context.Context, req *domain.GSTR9CStatusRequest) (*domain.GSTR9CStatusResponse, error) {
	url := a.apiURL("/enriched/returns?action=FILEDET&reference=" + req.Reference + "&rtn_typ=GSTR9C")
	body, err := a.doJSON(ctx, http.MethodGet, url, nil, "")
	if err != nil {
		return nil, err
	}
	var resp domain.GSTR9CStatusResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("adaequare: decode gstr9c status: %w", err)
	}
	resp.RequestID = req.RequestID
	return &resp, nil
}
