package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/identity-service/internal/domain"
	"github.com/complai/complai/services/go/identity-service/internal/store"
)

type Handlers struct {
	store       store.Repository
	keycloakURL string
	clientID    string
	clientSec   string
}

func NewHandlers(s store.Repository, keycloakURL, clientID, clientSec string) *Handlers {
	return &Handlers{store: s, keycloakURL: keycloakURL, clientID: clientID, clientSec: clientSec}
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "identity-service"})
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	tokenResp, err := h.keycloakLogin(req.Username, req.Password)
	if err != nil {
		log.Error().Err(err).Str("username", req.Username).Msg("keycloak login failed")
		httputil.JSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	httputil.JSON(w, http.StatusOK, tokenResp)
}

func (h *Handlers) Refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	tokenResp, err := h.keycloakRefresh(body.RefreshToken)
	if err != nil {
		httputil.JSON(w, http.StatusUnauthorized, map[string]string{"error": "refresh failed"})
		return
	}

	httputil.JSON(w, http.StatusOK, tokenResp)
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	if err := h.keycloakLogout(body.RefreshToken); err != nil {
		log.Error().Err(err).Msg("keycloak logout failed")
	}

	httputil.JSON(w, http.StatusOK, map[string]string{"status": "logged_out"})
}

func (h *Handlers) ListUsers(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	users, err := h.store.ListUsers(r.Context(), tenantID)
	if err != nil {
		log.Error().Err(err).Msg("list users failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if users == nil {
		users = []domain.User{}
	}

	httputil.JSON(w, http.StatusOK, users)
}

func (h *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(r.PathValue("userID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
		return
	}

	user, err := h.store.GetUserByID(r.Context(), tenantID, userID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}

	httputil.JSON(w, http.StatusOK, user)
}

func (h *Handlers) StepUpCheck(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	userID, err := userIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	sessionID, err := sessionIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.StepUpCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	valid, err := h.store.HasValidStepUp(r.Context(), tenantID, userID, sessionID, req.ActionClass)
	if err != nil {
		log.Error().Err(err).Msg("step-up check failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	if !valid {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":    "step_up_required",
			"message":  "This action requires step-up authentication",
			"action":   req.ActionClass,
			"step_up_url": "/v1/auth/step-up",
		})
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]string{"status": "step_up_valid"})
}

func (h *Handlers) StepUpVerify(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	userID, err := userIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	sessionID, err := sessionIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.StepUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	// In production, validate MFA code against enrolled factor.
	// For dev, accept code "123456" as valid.
	if req.MFACode != "123456" {
		factors, _ := h.store.GetMFAFactors(r.Context(), tenantID, userID)
		if len(factors) == 0 {
			httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "no MFA factor enrolled"})
			return
		}
		httputil.JSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid MFA code"})
		return
	}

	now := time.Now()
	evt := &domain.StepUpEvent{
		UserID:      userID,
		SessionID:   sessionID,
		ActionClass: req.ActionClass,
		VerifiedAt:  now,
		ExpiresAt:   now.Add(5 * time.Minute),
		MFAMethod:   "totp",
	}

	if err := h.store.CreateStepUpEvent(r.Context(), tenantID, evt); err != nil {
		log.Error().Err(err).Msg("create step-up event failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	httputil.JSON(w, http.StatusOK, map[string]interface{}{
		"status":     "step_up_verified",
		"action":     req.ActionClass,
		"expires_at": evt.ExpiresAt,
	})
}

func (h *Handlers) EnrollMFA(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	userID, err := userIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var body struct {
		FactorType string `json:"factor_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	secret := "dev-totp-secret-base32"
	f := &domain.MFAFactor{
		UserID:          userID,
		FactorType:      body.FactorType,
		SecretEncrypted: &secret,
		Verified:        true,
	}

	if err := h.store.CreateMFAFactor(r.Context(), tenantID, f); err != nil {
		log.Error().Err(err).Msg("enroll MFA failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	httputil.JSON(w, http.StatusCreated, map[string]interface{}{
		"factor_id":   f.ID,
		"factor_type": f.FactorType,
		"verified":    f.Verified,
	})
}

func (h *Handlers) keycloakLogin(username, password string) (*domain.TokenResponse, error) {
	resp, err := http.PostForm(h.keycloakURL+"/realms/complai/protocol/openid-connect/token",
		map[string][]string{
			"grant_type":    {"password"},
			"client_id":     {h.clientID},
			"client_secret": {h.clientSec},
			"username":      {username},
			"password":      {password},
		})
	if err != nil {
		return nil, fmt.Errorf("keycloak request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("keycloak returned %d: %s", resp.StatusCode, string(body))
	}

	var kcResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &kcResp); err != nil {
		return nil, fmt.Errorf("parse keycloak response: %w", err)
	}

	return &domain.TokenResponse{
		AccessToken:  kcResp.AccessToken,
		RefreshToken: kcResp.RefreshToken,
		TokenType:    kcResp.TokenType,
		ExpiresIn:    kcResp.ExpiresIn,
	}, nil
}

func (h *Handlers) keycloakRefresh(refreshToken string) (*domain.TokenResponse, error) {
	resp, err := http.PostForm(h.keycloakURL+"/realms/complai/protocol/openid-connect/token",
		map[string][]string{
			"grant_type":    {"refresh_token"},
			"client_id":     {h.clientID},
			"client_secret": {h.clientSec},
			"refresh_token": {refreshToken},
		})
	if err != nil {
		return nil, fmt.Errorf("keycloak refresh: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("keycloak returned %d: %s", resp.StatusCode, string(body))
	}

	var kcResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &kcResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &domain.TokenResponse{
		AccessToken:  kcResp.AccessToken,
		RefreshToken: kcResp.RefreshToken,
		TokenType:    kcResp.TokenType,
		ExpiresIn:    kcResp.ExpiresIn,
	}, nil
}

func (h *Handlers) keycloakLogout(refreshToken string) error {
	resp, err := http.PostForm(h.keycloakURL+"/realms/complai/protocol/openid-connect/logout",
		map[string][]string{
			"client_id":     {h.clientID},
			"client_secret": {h.clientSec},
			"refresh_token": {refreshToken},
		})
	if err != nil {
		return fmt.Errorf("keycloak logout: %w", err)
	}
	resp.Body.Close()
	return nil
}

func tenantIDFromRequest(r *http.Request) (uuid.UUID, error) {
	h := r.Header.Get("X-Tenant-Id")
	if h == "" {
		return uuid.Nil, fmt.Errorf("missing X-Tenant-Id header")
	}
	return uuid.Parse(h)
}

func userIDFromRequest(r *http.Request) (uuid.UUID, error) {
	h := r.Header.Get("X-User-Id")
	if h == "" {
		return uuid.Nil, fmt.Errorf("missing X-User-Id header")
	}
	return uuid.Parse(h)
}

func sessionIDFromRequest(r *http.Request) (uuid.UUID, error) {
	h := r.Header.Get("X-Session-Id")
	if h == "" {
		return uuid.Nil, fmt.Errorf("missing X-Session-Id header")
	}
	return uuid.Parse(h)
}
