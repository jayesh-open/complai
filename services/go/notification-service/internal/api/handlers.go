package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/notification-service/internal/domain"
	"github.com/complai/complai/services/go/notification-service/internal/store"
)

// EmailSender is an interface for testability of email sending.
type EmailSender interface {
	SendEmail(ctx context.Context, to, subject, body string) error
}

// SMTPSender implements EmailSender using net/smtp (for Mailpit/SES).
type SMTPSender struct {
	host string
	port int
	from string
}

func NewSMTPSender(host string, port int, from string) *SMTPSender {
	return &SMTPSender{host: host, port: port, from: from}
}

func (s *SMTPSender) SendEmail(ctx context.Context, to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", s.from, to, subject, body)
	return smtp.SendMail(addr, nil, s.from, []string{to}, []byte(msg))
}

// Handlers holds the dependencies for HTTP handlers.
type Handlers struct {
	store       store.Repository
	emailSender EmailSender
}

func NewHandlers(s store.Repository, emailSender EmailSender) *Handlers {
	return &Handlers{store: s, emailSender: emailSender}
}

// ---------------------------------------------------------------------------
// Health
// ---------------------------------------------------------------------------

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	httputil.JSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "notification-service"})
}

// ---------------------------------------------------------------------------
// SendNotification
// ---------------------------------------------------------------------------

func (h *Handlers) SendNotification(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.SendNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	channel := req.Channel
	if channel == "" {
		channel = "email"
	}

	status := "queued"

	subject := req.Subject
	body := req.Body
	digestGroup := req.DigestGroup
	metadata := req.Metadata
	if metadata == "" {
		metadata = "{}"
	}

	n := &domain.Notification{
		UserID:     req.UserID,
		TemplateID: req.TemplateID,
		Channel:    channel,
		Subject:    strPtr(subject),
		Body:       strPtr(body),
		Recipient:  req.Recipient,
		Status:     status,
		Metadata:   metadata,
	}
	if digestGroup != "" {
		n.DigestGroup = strPtr(digestGroup)
	}

	if err := h.store.CreateNotification(r.Context(), tenantID, n); err != nil {
		log.Error().Err(err).Msg("create notification failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "create failed"})
		return
	}

	// If channel is email, check user preferences and send
	if channel == "email" && h.emailSender != nil {
		prefs, err := h.store.GetPreferences(r.Context(), tenantID, req.UserID)
		if err != nil {
			// No preferences found — default is email_enabled=true
			log.Debug().Err(err).Msg("no preferences found, defaulting to email enabled")
			prefs = &domain.NotificationPreference{EmailEnabled: true, EmailValid: true}
		}

		if prefs.EmailEnabled && prefs.EmailValid {
			if err := h.emailSender.SendEmail(r.Context(), req.Recipient, subject, body); err != nil {
				log.Error().Err(err).Msg("send email failed")
				n.Status = "failed"
				failReason := err.Error()
				n.FailedReason = &failReason
			} else {
				n.Status = "sent"
				now := time.Now()
				n.SentAt = &now
			}
		}
	}

	httputil.JSON(w, http.StatusCreated, n)
}

// ---------------------------------------------------------------------------
// GetNotification
// ---------------------------------------------------------------------------

func (h *Handlers) GetNotification(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	notificationID, err := uuid.Parse(r.PathValue("notificationID"))
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid notification_id"})
		return
	}

	n, err := h.store.GetNotification(r.Context(), tenantID, notificationID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "notification not found"})
		return
	}

	httputil.JSON(w, http.StatusOK, n)
}

// ---------------------------------------------------------------------------
// ListNotifications
// ---------------------------------------------------------------------------

func (h *Handlers) ListNotifications(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	notifications, err := h.store.ListNotifications(r.Context(), tenantID)
	if err != nil {
		log.Error().Err(err).Msg("list notifications failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if notifications == nil {
		notifications = []domain.Notification{}
	}

	httputil.JSON(w, http.StatusOK, notifications)
}

// ---------------------------------------------------------------------------
// GetPreferences
// ---------------------------------------------------------------------------

func (h *Handlers) GetPreferences(w http.ResponseWriter, r *http.Request) {
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

	pref, err := h.store.GetPreferences(r.Context(), tenantID, userID)
	if err != nil {
		httputil.JSON(w, http.StatusNotFound, map[string]string{"error": "preferences not found"})
		return
	}

	httputil.JSON(w, http.StatusOK, pref)
}

// ---------------------------------------------------------------------------
// UpdatePreferences
// ---------------------------------------------------------------------------

func (h *Handlers) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
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

	var req domain.UpdatePreferencesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	pref := &domain.NotificationPreference{
		UserID:       userID,
		EmailEnabled: true,
	}
	if req.EmailEnabled != nil {
		pref.EmailEnabled = *req.EmailEnabled
	}
	if req.DigestEnabled != nil {
		pref.DigestEnabled = *req.DigestEnabled
	}
	pref.QuietHoursStart = req.QuietHoursStart
	pref.QuietHoursEnd = req.QuietHoursEnd
	pref.EmailAddress = req.EmailAddress

	if err := h.store.UpsertPreferences(r.Context(), tenantID, pref); err != nil {
		log.Error().Err(err).Msg("upsert preferences failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "update failed"})
		return
	}

	httputil.JSON(w, http.StatusOK, pref)
}

// ---------------------------------------------------------------------------
// CreateTemplate
// ---------------------------------------------------------------------------

func (h *Handlers) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.CreateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	channel := req.Channel
	if channel == "" {
		channel = "email"
	}
	variables := req.Variables
	if variables == "" {
		variables = "[]"
	}

	t := &domain.NotificationTemplate{
		Name:      req.Name,
		Channel:   channel,
		Subject:   req.Subject,
		Body:      req.Body,
		Variables: variables,
	}

	if err := h.store.CreateTemplate(r.Context(), tenantID, t); err != nil {
		log.Error().Err(err).Msg("create template failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "create failed"})
		return
	}

	httputil.JSON(w, http.StatusCreated, t)
}

// ---------------------------------------------------------------------------
// ListTemplates
// ---------------------------------------------------------------------------

func (h *Handlers) ListTemplates(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	templates, err := h.store.ListTemplates(r.Context(), tenantID)
	if err != nil {
		log.Error().Err(err).Msg("list templates failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if templates == nil {
		templates = []domain.NotificationTemplate{}
	}

	httputil.JSON(w, http.StatusOK, templates)
}

// ---------------------------------------------------------------------------
// ProcessBounce
// ---------------------------------------------------------------------------

func (h *Handlers) ProcessBounce(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.ProcessBounceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	var bounceSubtype *string
	if req.BounceSubtype != "" {
		bounceSubtype = &req.BounceSubtype
	}
	var diagnostic *string
	if req.Diagnostic != "" {
		diagnostic = &req.Diagnostic
	}

	bounce := &domain.NotificationBounce{
		NotificationID: req.NotificationID,
		BounceType:     req.BounceType,
		BounceSubtype:  bounceSubtype,
		EmailAddress:   req.EmailAddress,
		Diagnostic:     diagnostic,
	}

	if err := h.store.CreateBounce(r.Context(), tenantID, bounce); err != nil {
		log.Error().Err(err).Msg("create bounce failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "create bounce failed"})
		return
	}

	// Mark email as invalid
	if err := h.store.MarkEmailInvalid(r.Context(), tenantID, req.EmailAddress); err != nil {
		log.Error().Err(err).Msg("mark email invalid failed")
		// Non-fatal: bounce is already recorded
	}

	httputil.JSON(w, http.StatusOK, map[string]interface{}{
		"bounce":      bounce,
		"audit_event": "email_bounced",
		"email":       req.EmailAddress,
		"action":      "email_marked_invalid",
	})
}

// ---------------------------------------------------------------------------
// SendDigest
// ---------------------------------------------------------------------------

func (h *Handlers) SendDigest(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromRequest(r)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	var req domain.SendDigestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	cutoffTime, err := time.Parse(time.RFC3339, req.CutoffTime)
	if err != nil {
		httputil.JSON(w, http.StatusBadRequest, map[string]string{"error": "invalid cutoff_time, expected RFC3339"})
		return
	}

	// Get pending notifications grouped by user
	grouped, err := h.store.GetPendingDigestNotifications(r.Context(), tenantID, req.DigestGroup, cutoffTime)
	if err != nil {
		log.Error().Err(err).Msg("get pending digest notifications failed")
		httputil.JSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	var results []domain.DigestResult

	for userID, notifications := range grouped {
		result := domain.DigestResult{
			UserID:            userID,
			NotificationCount: len(notifications),
		}

		// Check user preferences
		prefs, err := h.store.GetPreferences(r.Context(), tenantID, userID)
		if err != nil {
			log.Debug().Err(err).Msg("no preferences for digest user, defaulting to digest enabled")
			prefs = &domain.NotificationPreference{DigestEnabled: true, EmailEnabled: true, EmailValid: true}
		}

		if !prefs.DigestEnabled || !prefs.EmailEnabled || !prefs.EmailValid {
			result.DigestSent = false
			results = append(results, result)
			continue
		}

		// Build digest email body by combining all notification bodies
		digestBody := "<h2>Notification Digest</h2><hr>"
		for _, n := range notifications {
			if n.Body != nil {
				digestBody += fmt.Sprintf("<div style=\"margin-bottom:12px;\">%s</div>", *n.Body)
			}
		}

		// Determine recipient from first notification
		recipient := notifications[0].Recipient

		// Determine subject
		digestSubject := fmt.Sprintf("Digest: %s (%d notifications)", req.DigestGroup, len(notifications))

		// Send single digest email
		batchID := uuid.New()
		if h.emailSender != nil {
			if err := h.emailSender.SendEmail(r.Context(), recipient, digestSubject, digestBody); err != nil {
				log.Error().Err(err).Msg("send digest email failed")
				result.DigestSent = false
				results = append(results, result)
				continue
			}
		}

		// Mark all notifications as sent with batch ID
		ids := make([]uuid.UUID, len(notifications))
		for i, n := range notifications {
			ids[i] = n.ID
		}
		if err := h.store.MarkNotificationsSent(r.Context(), tenantID, ids, batchID); err != nil {
			log.Error().Err(err).Msg("mark notifications sent failed")
		}

		result.DigestSent = true
		result.DigestID = &batchID
		results = append(results, result)
	}

	if results == nil {
		results = []domain.DigestResult{}
	}

	httputil.JSON(w, http.StatusOK, results)
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

func strPtr(s string) *string { return &s }
