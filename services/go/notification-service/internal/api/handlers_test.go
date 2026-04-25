package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/packages/shared-kernel-go/httputil"
	"github.com/complai/complai/services/go/notification-service/internal/domain"
)

// ---------------------------------------------------------------------------
// Mock store
// ---------------------------------------------------------------------------

type mockStore struct {
	createTemplateFn              func(ctx context.Context, tenantID uuid.UUID, t *domain.NotificationTemplate) error
	getTemplateFn                 func(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.NotificationTemplate, error)
	listTemplatesFn               func(ctx context.Context, tenantID uuid.UUID) ([]domain.NotificationTemplate, error)
	getPreferencesFn              func(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (*domain.NotificationPreference, error)
	upsertPreferencesFn           func(ctx context.Context, tenantID uuid.UUID, pref *domain.NotificationPreference) error
	createNotificationFn          func(ctx context.Context, tenantID uuid.UUID, n *domain.Notification) error
	getNotificationFn             func(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Notification, error)
	listNotificationsFn           func(ctx context.Context, tenantID uuid.UUID) ([]domain.Notification, error)
	getPendingDigestNotificationsFn func(ctx context.Context, tenantID uuid.UUID, digestGroup string, cutoffTime time.Time) (map[uuid.UUID][]domain.Notification, error)
	markNotificationsSentFn       func(ctx context.Context, tenantID uuid.UUID, ids []uuid.UUID, batchID uuid.UUID) error
	createBounceFn                func(ctx context.Context, tenantID uuid.UUID, b *domain.NotificationBounce) error
	markEmailInvalidFn            func(ctx context.Context, tenantID uuid.UUID, email string) error
}

func (m *mockStore) CreateTemplate(ctx context.Context, tenantID uuid.UUID, t *domain.NotificationTemplate) error {
	if m.createTemplateFn != nil {
		return m.createTemplateFn(ctx, tenantID, t)
	}
	t.ID = uuid.New()
	t.TenantID = tenantID
	t.Status = "active"
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return nil
}

func (m *mockStore) GetTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.NotificationTemplate, error) {
	if m.getTemplateFn != nil {
		return m.getTemplateFn(ctx, tenantID, id)
	}
	return nil, errors.New("not found")
}

func (m *mockStore) ListTemplates(ctx context.Context, tenantID uuid.UUID) ([]domain.NotificationTemplate, error) {
	if m.listTemplatesFn != nil {
		return m.listTemplatesFn(ctx, tenantID)
	}
	return nil, nil
}

func (m *mockStore) GetPreferences(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (*domain.NotificationPreference, error) {
	if m.getPreferencesFn != nil {
		return m.getPreferencesFn(ctx, tenantID, userID)
	}
	return nil, errors.New("not found")
}

func (m *mockStore) UpsertPreferences(ctx context.Context, tenantID uuid.UUID, pref *domain.NotificationPreference) error {
	if m.upsertPreferencesFn != nil {
		return m.upsertPreferencesFn(ctx, tenantID, pref)
	}
	pref.ID = uuid.New()
	pref.TenantID = tenantID
	pref.EmailValid = true
	pref.UnsubscribeToken = uuid.New()
	pref.CreatedAt = time.Now()
	pref.UpdatedAt = time.Now()
	return nil
}

func (m *mockStore) CreateNotification(ctx context.Context, tenantID uuid.UUID, n *domain.Notification) error {
	if m.createNotificationFn != nil {
		return m.createNotificationFn(ctx, tenantID, n)
	}
	n.ID = uuid.New()
	n.TenantID = tenantID
	n.CreatedAt = time.Now()
	return nil
}

func (m *mockStore) GetNotification(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Notification, error) {
	if m.getNotificationFn != nil {
		return m.getNotificationFn(ctx, tenantID, id)
	}
	return nil, errors.New("not found")
}

func (m *mockStore) ListNotifications(ctx context.Context, tenantID uuid.UUID) ([]domain.Notification, error) {
	if m.listNotificationsFn != nil {
		return m.listNotificationsFn(ctx, tenantID)
	}
	return nil, nil
}

func (m *mockStore) GetPendingDigestNotifications(ctx context.Context, tenantID uuid.UUID, digestGroup string, cutoffTime time.Time) (map[uuid.UUID][]domain.Notification, error) {
	if m.getPendingDigestNotificationsFn != nil {
		return m.getPendingDigestNotificationsFn(ctx, tenantID, digestGroup, cutoffTime)
	}
	return nil, nil
}

func (m *mockStore) MarkNotificationsSent(ctx context.Context, tenantID uuid.UUID, ids []uuid.UUID, batchID uuid.UUID) error {
	if m.markNotificationsSentFn != nil {
		return m.markNotificationsSentFn(ctx, tenantID, ids, batchID)
	}
	return nil
}

func (m *mockStore) CreateBounce(ctx context.Context, tenantID uuid.UUID, b *domain.NotificationBounce) error {
	if m.createBounceFn != nil {
		return m.createBounceFn(ctx, tenantID, b)
	}
	b.ID = uuid.New()
	b.TenantID = tenantID
	b.CreatedAt = time.Now()
	return nil
}

func (m *mockStore) MarkEmailInvalid(ctx context.Context, tenantID uuid.UUID, email string) error {
	if m.markEmailInvalidFn != nil {
		return m.markEmailInvalidFn(ctx, tenantID, email)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Mock email sender
// ---------------------------------------------------------------------------

type mockEmailSender struct {
	sendEmailFn func(ctx context.Context, to, subject, body string) error
	calls       []emailCall
}

type emailCall struct {
	To      string
	Subject string
	Body    string
}

func (m *mockEmailSender) SendEmail(ctx context.Context, to, subject, body string) error {
	m.calls = append(m.calls, emailCall{To: to, Subject: subject, Body: body})
	if m.sendEmailFn != nil {
		return m.sendEmailFn(ctx, to, subject, body)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func parseDataResponse(t *testing.T, body []byte, target interface{}) {
	t.Helper()
	var wrapper httputil.SuccessResponse
	wrapper.Data = target
	require.NoError(t, json.Unmarshal(body, &wrapper))
}

func newRequest(t *testing.T, method, path string, body interface{}, pathValues map[string]string) *http.Request {
	t.Helper()
	var req *http.Request
	if body != nil {
		b, err := json.Marshal(body)
		require.NoError(t, err)
		req = httptest.NewRequest(method, path, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	for k, v := range pathValues {
		req.SetPathValue(k, v)
	}
	return req
}

func reqWithTenant(req *http.Request, tenantID uuid.UUID) *http.Request {
	req.Header.Set("X-Tenant-Id", tenantID.String())
	return req
}

// ---------------------------------------------------------------------------
// Tests: Health
// ---------------------------------------------------------------------------

func TestHealth(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	h.Health(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var data map[string]string
	parseDataResponse(t, rec.Body.Bytes(), &data)
	assert.Equal(t, "ok", data["status"])
	assert.Equal(t, "notification-service", data["service"])
}

// ---------------------------------------------------------------------------
// Tests: SendNotification — creates notification and sends email
// ---------------------------------------------------------------------------

func TestSendNotification_Success(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	sender := &mockEmailSender{}

	ms := &mockStore{
		getPreferencesFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.NotificationPreference, error) {
			return &domain.NotificationPreference{
				EmailEnabled: true,
				EmailValid:   true,
			}, nil
		},
	}
	h := NewHandlers(ms, sender)

	body := domain.SendNotificationRequest{
		UserID:    userID,
		Channel:   "email",
		Subject:   "Test Subject",
		Body:      "<p>Hello</p>",
		Recipient: "user@example.com",
	}
	req := reqWithTenant(newRequest(t, http.MethodPost, "/v1/notifications/send", body, nil), tenantID)
	rec := httptest.NewRecorder()

	h.SendNotification(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var n domain.Notification
	parseDataResponse(t, rec.Body.Bytes(), &n)
	assert.Equal(t, "sent", n.Status)
	assert.NotEqual(t, uuid.Nil, n.ID)
	assert.Equal(t, tenantID, n.TenantID)
	assert.Equal(t, userID, n.UserID)

	// Verify email was sent
	require.Len(t, sender.calls, 1)
	assert.Equal(t, "user@example.com", sender.calls[0].To)
	assert.Equal(t, "Test Subject", sender.calls[0].Subject)
	assert.Equal(t, "<p>Hello</p>", sender.calls[0].Body)
}

// ---------------------------------------------------------------------------
// Tests: SendNotification — user with email_enabled=false
// ---------------------------------------------------------------------------

func TestSendNotification_EmailDisabled(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	sender := &mockEmailSender{}

	ms := &mockStore{
		getPreferencesFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.NotificationPreference, error) {
			return &domain.NotificationPreference{
				EmailEnabled: false,
				EmailValid:   true,
			}, nil
		},
	}
	h := NewHandlers(ms, sender)

	body := domain.SendNotificationRequest{
		UserID:    userID,
		Channel:   "email",
		Subject:   "Test Subject",
		Body:      "<p>Hello</p>",
		Recipient: "user@example.com",
	}
	req := reqWithTenant(newRequest(t, http.MethodPost, "/v1/notifications/send", body, nil), tenantID)
	rec := httptest.NewRecorder()

	h.SendNotification(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var n domain.Notification
	parseDataResponse(t, rec.Body.Bytes(), &n)
	assert.Equal(t, "queued", n.Status) // Not sent because email_enabled=false

	// Verify NO email was sent
	assert.Len(t, sender.calls, 0)
}

// ---------------------------------------------------------------------------
// Tests: SendNotification — invalid body
// ---------------------------------------------------------------------------

func TestSendNotification_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/notifications/send", bytes.NewReader([]byte("not json")))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.SendNotification(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: SendNotification — missing tenant header
// ---------------------------------------------------------------------------

func TestSendNotification_MissingTenant(t *testing.T) {
	h := NewHandlers(&mockStore{}, nil)

	body := domain.SendNotificationRequest{Channel: "email"}
	req := newRequest(t, http.MethodPost, "/v1/notifications/send", body, nil)
	rec := httptest.NewRecorder()

	h.SendNotification(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: SendNotification — store error
// ---------------------------------------------------------------------------

func TestSendNotification_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		createNotificationFn: func(_ context.Context, _ uuid.UUID, _ *domain.Notification) error {
			return errors.New("db error")
		},
	}
	h := NewHandlers(ms, nil)

	body := domain.SendNotificationRequest{
		UserID:    uuid.New(),
		Channel:   "email",
		Subject:   "Test",
		Body:      "Body",
		Recipient: "user@example.com",
	}
	req := reqWithTenant(newRequest(t, http.MethodPost, "/v1/notifications/send", body, nil), tenantID)
	rec := httptest.NewRecorder()

	h.SendNotification(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: SendNotification — email send failure
// ---------------------------------------------------------------------------

func TestSendNotification_EmailSendError(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	sender := &mockEmailSender{
		sendEmailFn: func(_ context.Context, _, _, _ string) error {
			return errors.New("smtp error")
		},
	}

	ms := &mockStore{
		getPreferencesFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.NotificationPreference, error) {
			return &domain.NotificationPreference{EmailEnabled: true, EmailValid: true}, nil
		},
	}
	h := NewHandlers(ms, sender)

	body := domain.SendNotificationRequest{
		UserID:    userID,
		Channel:   "email",
		Subject:   "Test",
		Body:      "Body",
		Recipient: "user@example.com",
	}
	req := reqWithTenant(newRequest(t, http.MethodPost, "/v1/notifications/send", body, nil), tenantID)
	rec := httptest.NewRecorder()

	h.SendNotification(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var n domain.Notification
	parseDataResponse(t, rec.Body.Bytes(), &n)
	assert.Equal(t, "failed", n.Status)
	require.NotNil(t, n.FailedReason)
	assert.Contains(t, *n.FailedReason, "smtp error")
}

// ---------------------------------------------------------------------------
// Tests: SendNotification — no preferences found defaults to sending
// ---------------------------------------------------------------------------

func TestSendNotification_NoPreferencesDefaultSend(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	sender := &mockEmailSender{}

	ms := &mockStore{
		getPreferencesFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.NotificationPreference, error) {
			return nil, errors.New("not found")
		},
	}
	h := NewHandlers(ms, sender)

	body := domain.SendNotificationRequest{
		UserID:    userID,
		Channel:   "email",
		Subject:   "Test",
		Body:      "Body",
		Recipient: "user@example.com",
	}
	req := reqWithTenant(newRequest(t, http.MethodPost, "/v1/notifications/send", body, nil), tenantID)
	rec := httptest.NewRecorder()

	h.SendNotification(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	// Should still send because default is email_enabled=true
	require.Len(t, sender.calls, 1)
}

// ---------------------------------------------------------------------------
// Tests: GetNotification
// ---------------------------------------------------------------------------

func TestGetNotification_Success(t *testing.T) {
	tenantID := uuid.New()
	notifID := uuid.New()
	ms := &mockStore{
		getNotificationFn: func(_ context.Context, _ uuid.UUID, id uuid.UUID) (*domain.Notification, error) {
			assert.Equal(t, notifID, id)
			subj := "Test"
			return &domain.Notification{
				ID:        notifID,
				TenantID:  tenantID,
				UserID:    uuid.New(),
				Channel:   "email",
				Subject:   &subj,
				Recipient: "user@example.com",
				Status:    "sent",
				Metadata:  "{}",
			}, nil
		},
	}
	h := NewHandlers(ms, nil)

	req := reqWithTenant(newRequest(t, http.MethodGet, "/v1/notifications/"+notifID.String(), nil, map[string]string{
		"notificationID": notifID.String(),
	}), tenantID)
	rec := httptest.NewRecorder()

	h.GetNotification(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var n domain.Notification
	parseDataResponse(t, rec.Body.Bytes(), &n)
	assert.Equal(t, notifID, n.ID)
}

func TestGetNotification_InvalidID(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, nil)

	req := reqWithTenant(newRequest(t, http.MethodGet, "/v1/notifications/bad", nil, map[string]string{
		"notificationID": "bad",
	}), tenantID)
	rec := httptest.NewRecorder()

	h.GetNotification(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetNotification_NotFound(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		getNotificationFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.Notification, error) {
			return nil, errors.New("not found")
		},
	}
	h := NewHandlers(ms, nil)

	req := reqWithTenant(newRequest(t, http.MethodGet, "/v1/notifications/"+uuid.New().String(), nil, map[string]string{
		"notificationID": uuid.New().String(),
	}), tenantID)
	rec := httptest.NewRecorder()

	h.GetNotification(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: ListNotifications
// ---------------------------------------------------------------------------

func TestListNotifications_Success(t *testing.T) {
	tenantID := uuid.New()
	subj := "Test"
	ms := &mockStore{
		listNotificationsFn: func(_ context.Context, _ uuid.UUID) ([]domain.Notification, error) {
			return []domain.Notification{
				{ID: uuid.New(), TenantID: tenantID, Channel: "email", Subject: &subj, Recipient: "user@example.com", Status: "sent", Metadata: "{}"},
			}, nil
		},
	}
	h := NewHandlers(ms, nil)

	req := reqWithTenant(httptest.NewRequest(http.MethodGet, "/v1/notifications", nil), tenantID)
	rec := httptest.NewRecorder()

	h.ListNotifications(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var notifications []domain.Notification
	parseDataResponse(t, rec.Body.Bytes(), &notifications)
	require.Len(t, notifications, 1)
}

func TestListNotifications_EmptyList(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listNotificationsFn: func(_ context.Context, _ uuid.UUID) ([]domain.Notification, error) {
			return nil, nil
		},
	}
	h := NewHandlers(ms, nil)

	req := reqWithTenant(httptest.NewRequest(http.MethodGet, "/v1/notifications", nil), tenantID)
	rec := httptest.NewRecorder()

	h.ListNotifications(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var notifications []domain.Notification
	parseDataResponse(t, rec.Body.Bytes(), &notifications)
	assert.NotNil(t, notifications)
	assert.Len(t, notifications, 0)
}

func TestListNotifications_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listNotificationsFn: func(_ context.Context, _ uuid.UUID) ([]domain.Notification, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewHandlers(ms, nil)

	req := reqWithTenant(httptest.NewRequest(http.MethodGet, "/v1/notifications", nil), tenantID)
	rec := httptest.NewRecorder()

	h.ListNotifications(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: GetPreferences
// ---------------------------------------------------------------------------

func TestGetPreferences_Success(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	ms := &mockStore{
		getPreferencesFn: func(_ context.Context, _ uuid.UUID, uid uuid.UUID) (*domain.NotificationPreference, error) {
			assert.Equal(t, userID, uid)
			return &domain.NotificationPreference{
				ID:           uuid.New(),
				TenantID:     tenantID,
				UserID:       userID,
				EmailEnabled: true,
				EmailValid:   true,
			}, nil
		},
	}
	h := NewHandlers(ms, nil)

	req := reqWithTenant(newRequest(t, http.MethodGet, "/v1/users/"+userID.String()+"/preferences", nil, map[string]string{
		"userID": userID.String(),
	}), tenantID)
	rec := httptest.NewRecorder()

	h.GetPreferences(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var pref domain.NotificationPreference
	parseDataResponse(t, rec.Body.Bytes(), &pref)
	assert.Equal(t, userID, pref.UserID)
	assert.True(t, pref.EmailEnabled)
}

func TestGetPreferences_InvalidUserID(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, nil)

	req := reqWithTenant(newRequest(t, http.MethodGet, "/v1/users/bad/preferences", nil, map[string]string{
		"userID": "bad",
	}), tenantID)
	rec := httptest.NewRecorder()

	h.GetPreferences(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetPreferences_NotFound(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	ms := &mockStore{
		getPreferencesFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.NotificationPreference, error) {
			return nil, errors.New("not found")
		},
	}
	h := NewHandlers(ms, nil)

	req := reqWithTenant(newRequest(t, http.MethodGet, "/v1/users/"+userID.String()+"/preferences", nil, map[string]string{
		"userID": userID.String(),
	}), tenantID)
	rec := httptest.NewRecorder()

	h.GetPreferences(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: UpdatePreferences
// ---------------------------------------------------------------------------

func TestUpdatePreferences_Success(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	ms := &mockStore{}
	h := NewHandlers(ms, nil)

	emailEnabled := true
	digestEnabled := true
	emailAddr := "user@example.com"
	body := domain.UpdatePreferencesRequest{
		EmailEnabled:  &emailEnabled,
		DigestEnabled: &digestEnabled,
		EmailAddress:  &emailAddr,
	}
	req := reqWithTenant(newRequest(t, http.MethodPut, "/v1/users/"+userID.String()+"/preferences", body, map[string]string{
		"userID": userID.String(),
	}), tenantID)
	rec := httptest.NewRecorder()

	h.UpdatePreferences(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var pref domain.NotificationPreference
	parseDataResponse(t, rec.Body.Bytes(), &pref)
	assert.Equal(t, userID, pref.UserID)
	assert.Equal(t, tenantID, pref.TenantID)
	assert.True(t, pref.EmailEnabled)
	assert.True(t, pref.DigestEnabled)
}

func TestUpdatePreferences_InvalidUserID(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, nil)

	req := reqWithTenant(newRequest(t, http.MethodPut, "/v1/users/bad/preferences", map[string]string{}, map[string]string{
		"userID": "bad",
	}), tenantID)
	rec := httptest.NewRecorder()

	h.UpdatePreferences(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdatePreferences_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	h := NewHandlers(&mockStore{}, nil)

	req := httptest.NewRequest(http.MethodPut, "/v1/users/"+userID.String()+"/preferences", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	req.SetPathValue("userID", userID.String())
	rec := httptest.NewRecorder()

	h.UpdatePreferences(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdatePreferences_StoreError(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	ms := &mockStore{
		upsertPreferencesFn: func(_ context.Context, _ uuid.UUID, _ *domain.NotificationPreference) error {
			return errors.New("db error")
		},
	}
	h := NewHandlers(ms, nil)

	emailEnabled := true
	body := domain.UpdatePreferencesRequest{EmailEnabled: &emailEnabled}
	req := reqWithTenant(newRequest(t, http.MethodPut, "/v1/users/"+userID.String()+"/preferences", body, map[string]string{
		"userID": userID.String(),
	}), tenantID)
	rec := httptest.NewRecorder()

	h.UpdatePreferences(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: CreateTemplate
// ---------------------------------------------------------------------------

func TestCreateTemplate_Success(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{}
	h := NewHandlers(ms, nil)

	subj := "Welcome {{name}}"
	body := domain.CreateTemplateRequest{
		Name:    "welcome",
		Channel: "email",
		Subject: &subj,
		Body:    "<p>Hello {{name}}</p>",
	}
	req := reqWithTenant(newRequest(t, http.MethodPost, "/v1/templates", body, nil), tenantID)
	rec := httptest.NewRecorder()

	h.CreateTemplate(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var tmpl domain.NotificationTemplate
	parseDataResponse(t, rec.Body.Bytes(), &tmpl)
	assert.Equal(t, "welcome", tmpl.Name)
	assert.Equal(t, "email", tmpl.Channel)
	assert.Equal(t, "active", tmpl.Status)
	assert.Equal(t, tenantID, tmpl.TenantID)
}

func TestCreateTemplate_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/templates", bytes.NewReader([]byte("bad json")))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.CreateTemplate(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateTemplate_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		createTemplateFn: func(_ context.Context, _ uuid.UUID, _ *domain.NotificationTemplate) error {
			return errors.New("db error")
		},
	}
	h := NewHandlers(ms, nil)

	body := domain.CreateTemplateRequest{Name: "test", Body: "<p>test</p>"}
	req := reqWithTenant(newRequest(t, http.MethodPost, "/v1/templates", body, nil), tenantID)
	rec := httptest.NewRecorder()

	h.CreateTemplate(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: ListTemplates
// ---------------------------------------------------------------------------

func TestListTemplates_Success(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listTemplatesFn: func(_ context.Context, _ uuid.UUID) ([]domain.NotificationTemplate, error) {
			return []domain.NotificationTemplate{
				{ID: uuid.New(), TenantID: tenantID, Name: "welcome", Channel: "email", Body: "<p>Hello</p>", Status: "active", Variables: "[]"},
			}, nil
		},
	}
	h := NewHandlers(ms, nil)

	req := reqWithTenant(httptest.NewRequest(http.MethodGet, "/v1/templates", nil), tenantID)
	rec := httptest.NewRecorder()

	h.ListTemplates(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var templates []domain.NotificationTemplate
	parseDataResponse(t, rec.Body.Bytes(), &templates)
	require.Len(t, templates, 1)
	assert.Equal(t, "welcome", templates[0].Name)
}

func TestListTemplates_EmptyList(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{}
	h := NewHandlers(ms, nil)

	req := reqWithTenant(httptest.NewRequest(http.MethodGet, "/v1/templates", nil), tenantID)
	rec := httptest.NewRecorder()

	h.ListTemplates(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var templates []domain.NotificationTemplate
	parseDataResponse(t, rec.Body.Bytes(), &templates)
	assert.NotNil(t, templates)
	assert.Len(t, templates, 0)
}

func TestListTemplates_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		listTemplatesFn: func(_ context.Context, _ uuid.UUID) ([]domain.NotificationTemplate, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewHandlers(ms, nil)

	req := reqWithTenant(httptest.NewRequest(http.MethodGet, "/v1/templates", nil), tenantID)
	rec := httptest.NewRecorder()

	h.ListTemplates(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: ProcessBounce — records bounce, marks email invalid
// ---------------------------------------------------------------------------

func TestProcessBounce_Success(t *testing.T) {
	tenantID := uuid.New()
	notifID := uuid.New()
	var capturedEmail string

	ms := &mockStore{
		markEmailInvalidFn: func(_ context.Context, _ uuid.UUID, email string) error {
			capturedEmail = email
			return nil
		},
	}
	h := NewHandlers(ms, nil)

	body := domain.ProcessBounceRequest{
		NotificationID: &notifID,
		BounceType:     "Permanent",
		BounceSubtype:  "General",
		EmailAddress:   "bounced@example.com",
		Diagnostic:     "550 5.1.1 user unknown",
	}
	req := reqWithTenant(newRequest(t, http.MethodPost, "/v1/notifications/bounce", body, nil), tenantID)
	rec := httptest.NewRecorder()

	h.ProcessBounce(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "bounced@example.com", capturedEmail)

	// Verify response contains audit event data
	var resp map[string]interface{}
	parseDataResponse(t, rec.Body.Bytes(), &resp)
	assert.Equal(t, "email_bounced", resp["audit_event"])
	assert.Equal(t, "bounced@example.com", resp["email"])
	assert.Equal(t, "email_marked_invalid", resp["action"])
}

func TestProcessBounce_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/notifications/bounce", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.ProcessBounce(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestProcessBounce_CreateBounceError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		createBounceFn: func(_ context.Context, _ uuid.UUID, _ *domain.NotificationBounce) error {
			return errors.New("db error")
		},
	}
	h := NewHandlers(ms, nil)

	body := domain.ProcessBounceRequest{
		BounceType:   "Permanent",
		EmailAddress: "bounced@example.com",
	}
	req := reqWithTenant(newRequest(t, http.MethodPost, "/v1/notifications/bounce", body, nil), tenantID)
	rec := httptest.NewRecorder()

	h.ProcessBounce(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestProcessBounce_MarkEmailInvalidError_NonFatal(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		markEmailInvalidFn: func(_ context.Context, _ uuid.UUID, _ string) error {
			return errors.New("db error")
		},
	}
	h := NewHandlers(ms, nil)

	body := domain.ProcessBounceRequest{
		BounceType:   "Permanent",
		EmailAddress: "bounced@example.com",
	}
	req := reqWithTenant(newRequest(t, http.MethodPost, "/v1/notifications/bounce", body, nil), tenantID)
	rec := httptest.NewRecorder()

	h.ProcessBounce(rec, req)

	// Should still return OK since bounce was recorded; MarkEmailInvalid failure is non-fatal
	assert.Equal(t, http.StatusOK, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: SendDigest — 5 queued notifications for same user -> single digest email
// ---------------------------------------------------------------------------

func TestSendDigest_Success(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	sender := &mockEmailSender{}
	var markedIDs []uuid.UUID

	notifIDs := make([]uuid.UUID, 5)
	notifications := make([]domain.Notification, 5)
	for i := 0; i < 5; i++ {
		notifIDs[i] = uuid.New()
		body := "<p>Notification " + string(rune('1'+i)) + "</p>"
		notifications[i] = domain.Notification{
			ID:        notifIDs[i],
			TenantID:  tenantID,
			UserID:    userID,
			Channel:   "email",
			Body:      &body,
			Recipient: "user@example.com",
			Status:    "queued",
			Metadata:  "{}",
		}
	}

	ms := &mockStore{
		getPendingDigestNotificationsFn: func(_ context.Context, _ uuid.UUID, dg string, _ time.Time) (map[uuid.UUID][]domain.Notification, error) {
			assert.Equal(t, "daily", dg)
			return map[uuid.UUID][]domain.Notification{
				userID: notifications,
			}, nil
		},
		getPreferencesFn: func(_ context.Context, _ uuid.UUID, uid uuid.UUID) (*domain.NotificationPreference, error) {
			assert.Equal(t, userID, uid)
			return &domain.NotificationPreference{
				DigestEnabled: true,
				EmailEnabled:  true,
				EmailValid:    true,
			}, nil
		},
		markNotificationsSentFn: func(_ context.Context, _ uuid.UUID, ids []uuid.UUID, _ uuid.UUID) error {
			markedIDs = ids
			return nil
		},
	}
	h := NewHandlers(ms, sender)

	body := domain.SendDigestRequest{
		DigestGroup: "daily",
		CutoffTime:  time.Now().Format(time.RFC3339),
	}
	req := reqWithTenant(newRequest(t, http.MethodPost, "/v1/notifications/digest", body, nil), tenantID)
	rec := httptest.NewRecorder()

	h.SendDigest(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var results []domain.DigestResult
	parseDataResponse(t, rec.Body.Bytes(), &results)
	require.Len(t, results, 1)
	assert.Equal(t, userID, results[0].UserID)
	assert.Equal(t, 5, results[0].NotificationCount)
	assert.True(t, results[0].DigestSent)
	require.NotNil(t, results[0].DigestID)

	// Verify only ONE email was sent (digest)
	require.Len(t, sender.calls, 1)
	assert.Equal(t, "user@example.com", sender.calls[0].To)
	assert.Contains(t, sender.calls[0].Subject, "Digest: daily")
	assert.Contains(t, sender.calls[0].Subject, "5 notifications")

	// Verify all 5 notifications were marked as sent
	require.Len(t, markedIDs, 5)
}

func TestSendDigest_DigestDisabled(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	sender := &mockEmailSender{}

	body1 := "<p>Notification 1</p>"
	ms := &mockStore{
		getPendingDigestNotificationsFn: func(_ context.Context, _ uuid.UUID, _ string, _ time.Time) (map[uuid.UUID][]domain.Notification, error) {
			return map[uuid.UUID][]domain.Notification{
				userID: {
					{ID: uuid.New(), TenantID: tenantID, UserID: userID, Body: &body1, Recipient: "user@example.com", Status: "queued", Metadata: "{}"},
				},
			}, nil
		},
		getPreferencesFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.NotificationPreference, error) {
			return &domain.NotificationPreference{
				DigestEnabled: false,
				EmailEnabled:  true,
				EmailValid:    true,
			}, nil
		},
	}
	h := NewHandlers(ms, sender)

	reqBody := domain.SendDigestRequest{
		DigestGroup: "daily",
		CutoffTime:  time.Now().Format(time.RFC3339),
	}
	req := reqWithTenant(newRequest(t, http.MethodPost, "/v1/notifications/digest", reqBody, nil), tenantID)
	rec := httptest.NewRecorder()

	h.SendDigest(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var results []domain.DigestResult
	parseDataResponse(t, rec.Body.Bytes(), &results)
	require.Len(t, results, 1)
	assert.False(t, results[0].DigestSent)

	// No email should have been sent
	assert.Len(t, sender.calls, 0)
}

func TestSendDigest_InvalidBody(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, nil)

	req := httptest.NewRequest(http.MethodPost, "/v1/notifications/digest", bytes.NewReader([]byte("bad")))
	req.Header.Set("X-Tenant-Id", tenantID.String())
	rec := httptest.NewRecorder()

	h.SendDigest(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSendDigest_InvalidCutoffTime(t *testing.T) {
	tenantID := uuid.New()
	h := NewHandlers(&mockStore{}, nil)

	body := domain.SendDigestRequest{
		DigestGroup: "daily",
		CutoffTime:  "not-a-time",
	}
	req := reqWithTenant(newRequest(t, http.MethodPost, "/v1/notifications/digest", body, nil), tenantID)
	rec := httptest.NewRecorder()

	h.SendDigest(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSendDigest_StoreError(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		getPendingDigestNotificationsFn: func(_ context.Context, _ uuid.UUID, _ string, _ time.Time) (map[uuid.UUID][]domain.Notification, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewHandlers(ms, nil)

	body := domain.SendDigestRequest{
		DigestGroup: "daily",
		CutoffTime:  time.Now().Format(time.RFC3339),
	}
	req := reqWithTenant(newRequest(t, http.MethodPost, "/v1/notifications/digest", body, nil), tenantID)
	rec := httptest.NewRecorder()

	h.SendDigest(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestSendDigest_EmptyResult(t *testing.T) {
	tenantID := uuid.New()
	ms := &mockStore{
		getPendingDigestNotificationsFn: func(_ context.Context, _ uuid.UUID, _ string, _ time.Time) (map[uuid.UUID][]domain.Notification, error) {
			return map[uuid.UUID][]domain.Notification{}, nil
		},
	}
	h := NewHandlers(ms, nil)

	body := domain.SendDigestRequest{
		DigestGroup: "daily",
		CutoffTime:  time.Now().Format(time.RFC3339),
	}
	req := reqWithTenant(newRequest(t, http.MethodPost, "/v1/notifications/digest", body, nil), tenantID)
	rec := httptest.NewRecorder()

	h.SendDigest(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var results []domain.DigestResult
	parseDataResponse(t, rec.Body.Bytes(), &results)
	assert.NotNil(t, results)
	assert.Len(t, results, 0)
}

func TestSendDigest_EmailSendError(t *testing.T) {
	tenantID := uuid.New()
	userID := uuid.New()
	sender := &mockEmailSender{
		sendEmailFn: func(_ context.Context, _, _, _ string) error {
			return errors.New("smtp error")
		},
	}

	body1 := "<p>Notification</p>"
	ms := &mockStore{
		getPendingDigestNotificationsFn: func(_ context.Context, _ uuid.UUID, _ string, _ time.Time) (map[uuid.UUID][]domain.Notification, error) {
			return map[uuid.UUID][]domain.Notification{
				userID: {
					{ID: uuid.New(), TenantID: tenantID, UserID: userID, Body: &body1, Recipient: "user@example.com", Status: "queued", Metadata: "{}"},
				},
			}, nil
		},
		getPreferencesFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID) (*domain.NotificationPreference, error) {
			return &domain.NotificationPreference{DigestEnabled: true, EmailEnabled: true, EmailValid: true}, nil
		},
	}
	h := NewHandlers(ms, sender)

	reqBody := domain.SendDigestRequest{
		DigestGroup: "daily",
		CutoffTime:  time.Now().Format(time.RFC3339),
	}
	req := reqWithTenant(newRequest(t, http.MethodPost, "/v1/notifications/digest", reqBody, nil), tenantID)
	rec := httptest.NewRecorder()

	h.SendDigest(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var results []domain.DigestResult
	parseDataResponse(t, rec.Body.Bytes(), &results)
	require.Len(t, results, 1)
	assert.False(t, results[0].DigestSent)
}

// ---------------------------------------------------------------------------
// Tests: tenantIDFromRequest
// ---------------------------------------------------------------------------

func TestTenantIDFromRequest_Valid(t *testing.T) {
	expected := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-Id", expected.String())

	got, err := tenantIDFromRequest(req)
	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestTenantIDFromRequest_Missing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	_, err := tenantIDFromRequest(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing X-Tenant-Id")
}

func TestTenantIDFromRequest_Invalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-Id", "not-uuid")

	_, err := tenantIDFromRequest(req)
	require.Error(t, err)
}

// ---------------------------------------------------------------------------
// Tests: strPtr helper
// ---------------------------------------------------------------------------

func TestStrPtr(t *testing.T) {
	s := strPtr("hello")
	require.NotNil(t, s)
	assert.Equal(t, "hello", *s)
}

// ---------------------------------------------------------------------------
// Tests: NewRouter
// ---------------------------------------------------------------------------

func TestNewRouter(t *testing.T) {
	ms := &mockStore{}
	r := NewRouter(ms, nil)
	require.NotNil(t, r)

	// Verify health endpoint
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify ping heartbeat endpoint
	req = httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// ---------------------------------------------------------------------------
// Tests: NewSMTPSender construction
// ---------------------------------------------------------------------------

func TestNewSMTPSender(t *testing.T) {
	s := NewSMTPSender("localhost", 1025, "noreply@complai.in")
	assert.Equal(t, "localhost", s.host)
	assert.Equal(t, 1025, s.port)
	assert.Equal(t, "noreply@complai.in", s.from)
}
