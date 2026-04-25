package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/complai/complai/services/go/notification-service/internal/domain"
)

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) Pool() *pgxpool.Pool { return s.pool }

// Compile-time interface check.
var _ Repository = (*Store)(nil)

func setTenantID(ctx context.Context, tx pgx.Tx, tenantID uuid.UUID) error {
	_, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID.String()))
	return err
}

// ---------------------------------------------------------------------------
// Templates
// ---------------------------------------------------------------------------

func (s *Store) CreateTemplate(ctx context.Context, tenantID uuid.UUID, t *domain.NotificationTemplate) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	t.ID = uuid.New()
	t.TenantID = tenantID

	variables := t.Variables
	if variables == "" {
		variables = "[]"
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO notification_templates (id, tenant_id, name, channel, subject, body, variables, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING created_at, updated_at`,
		t.ID, tenantID, t.Name, t.Channel, t.Subject, t.Body, variables, "active",
	).Scan(&t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert template: %w", err)
	}
	t.Status = "active"
	t.Variables = variables
	return tx.Commit(ctx)
}

func (s *Store) GetTemplate(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.NotificationTemplate, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var t domain.NotificationTemplate
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, name, channel, subject, body, variables, status, created_at, updated_at
		 FROM notification_templates WHERE id = $1`, id,
	).Scan(&t.ID, &t.TenantID, &t.Name, &t.Channel, &t.Subject, &t.Body, &t.Variables, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get template: %w", err)
	}
	return &t, tx.Commit(ctx)
}

func (s *Store) ListTemplates(ctx context.Context, tenantID uuid.UUID) ([]domain.NotificationTemplate, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, name, channel, subject, body, variables, status, created_at, updated_at
		 FROM notification_templates ORDER BY created_at`)
	if err != nil {
		return nil, fmt.Errorf("list templates: %w", err)
	}
	defer rows.Close()

	var templates []domain.NotificationTemplate
	for rows.Next() {
		var t domain.NotificationTemplate
		if err := rows.Scan(&t.ID, &t.TenantID, &t.Name, &t.Channel, &t.Subject, &t.Body, &t.Variables, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan template: %w", err)
		}
		templates = append(templates, t)
	}
	return templates, tx.Commit(ctx)
}

// ---------------------------------------------------------------------------
// Preferences
// ---------------------------------------------------------------------------

func (s *Store) GetPreferences(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (*domain.NotificationPreference, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var p domain.NotificationPreference
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, user_id, email_enabled, digest_enabled, quiet_hours_start, quiet_hours_end,
		        email_address, email_valid, bounce_count, unsubscribe_token, created_at, updated_at
		 FROM notification_preferences WHERE user_id = $1`, userID,
	).Scan(&p.ID, &p.TenantID, &p.UserID, &p.EmailEnabled, &p.DigestEnabled, &p.QuietHoursStart, &p.QuietHoursEnd,
		&p.EmailAddress, &p.EmailValid, &p.BounceCount, &p.UnsubscribeToken, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get preferences: %w", err)
	}
	return &p, tx.Commit(ctx)
}

func (s *Store) UpsertPreferences(ctx context.Context, tenantID uuid.UUID, pref *domain.NotificationPreference) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO notification_preferences (tenant_id, user_id, email_enabled, digest_enabled, quiet_hours_start, quiet_hours_end, email_address)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (tenant_id, user_id) DO UPDATE SET
		   email_enabled = EXCLUDED.email_enabled,
		   digest_enabled = EXCLUDED.digest_enabled,
		   quiet_hours_start = EXCLUDED.quiet_hours_start,
		   quiet_hours_end = EXCLUDED.quiet_hours_end,
		   email_address = EXCLUDED.email_address,
		   updated_at = now()
		 RETURNING id, unsubscribe_token, email_valid, bounce_count, created_at, updated_at`,
		tenantID, pref.UserID, pref.EmailEnabled, pref.DigestEnabled,
		pref.QuietHoursStart, pref.QuietHoursEnd, pref.EmailAddress,
	).Scan(&pref.ID, &pref.UnsubscribeToken, &pref.EmailValid, &pref.BounceCount, &pref.CreatedAt, &pref.UpdatedAt)
	if err != nil {
		return fmt.Errorf("upsert preferences: %w", err)
	}
	pref.TenantID = tenantID
	return tx.Commit(ctx)
}

// ---------------------------------------------------------------------------
// Notifications
// ---------------------------------------------------------------------------

func (s *Store) CreateNotification(ctx context.Context, tenantID uuid.UUID, n *domain.Notification) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	n.ID = uuid.New()
	n.TenantID = tenantID

	metadata := n.Metadata
	if metadata == "" {
		metadata = "{}"
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO notifications (id, tenant_id, user_id, template_id, channel, subject, body, recipient, status, digest_group, metadata)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		 RETURNING created_at`,
		n.ID, tenantID, n.UserID, n.TemplateID, n.Channel, n.Subject, n.Body, n.Recipient, n.Status, n.DigestGroup, metadata,
	).Scan(&n.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert notification: %w", err)
	}
	n.Metadata = metadata
	return tx.Commit(ctx)
}

func (s *Store) GetNotification(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Notification, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var n domain.Notification
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, user_id, template_id, channel, subject, body, recipient, status,
		        sent_at, failed_reason, digest_group, digest_batch_id, metadata, created_at
		 FROM notifications WHERE id = $1`, id,
	).Scan(&n.ID, &n.TenantID, &n.UserID, &n.TemplateID, &n.Channel, &n.Subject, &n.Body, &n.Recipient, &n.Status,
		&n.SentAt, &n.FailedReason, &n.DigestGroup, &n.DigestBatchID, &n.Metadata, &n.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get notification: %w", err)
	}
	return &n, tx.Commit(ctx)
}

func (s *Store) ListNotifications(ctx context.Context, tenantID uuid.UUID) ([]domain.Notification, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, user_id, template_id, channel, subject, body, recipient, status,
		        sent_at, failed_reason, digest_group, digest_batch_id, metadata, created_at
		 FROM notifications ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("list notifications: %w", err)
	}
	defer rows.Close()

	var notifications []domain.Notification
	for rows.Next() {
		var n domain.Notification
		if err := rows.Scan(&n.ID, &n.TenantID, &n.UserID, &n.TemplateID, &n.Channel, &n.Subject, &n.Body, &n.Recipient, &n.Status,
			&n.SentAt, &n.FailedReason, &n.DigestGroup, &n.DigestBatchID, &n.Metadata, &n.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan notification: %w", err)
		}
		notifications = append(notifications, n)
	}
	return notifications, tx.Commit(ctx)
}

// ---------------------------------------------------------------------------
// Digest
// ---------------------------------------------------------------------------

func (s *Store) GetPendingDigestNotifications(ctx context.Context, tenantID uuid.UUID, digestGroup string, cutoffTime time.Time) (map[uuid.UUID][]domain.Notification, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, user_id, template_id, channel, subject, body, recipient, status,
		        sent_at, failed_reason, digest_group, digest_batch_id, metadata, created_at
		 FROM notifications
		 WHERE digest_group = $1 AND status = 'queued' AND created_at <= $2
		 ORDER BY user_id, created_at`, digestGroup, cutoffTime)
	if err != nil {
		return nil, fmt.Errorf("get pending digest: %w", err)
	}
	defer rows.Close()

	grouped := make(map[uuid.UUID][]domain.Notification)
	for rows.Next() {
		var n domain.Notification
		if err := rows.Scan(&n.ID, &n.TenantID, &n.UserID, &n.TemplateID, &n.Channel, &n.Subject, &n.Body, &n.Recipient, &n.Status,
			&n.SentAt, &n.FailedReason, &n.DigestGroup, &n.DigestBatchID, &n.Metadata, &n.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan digest notification: %w", err)
		}
		grouped[n.UserID] = append(grouped[n.UserID], n)
	}
	return grouped, tx.Commit(ctx)
}

func (s *Store) MarkNotificationsSent(ctx context.Context, tenantID uuid.UUID, ids []uuid.UUID, batchID uuid.UUID) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	for _, id := range ids {
		_, err := tx.Exec(ctx,
			`UPDATE notifications SET status = 'sent', sent_at = now(), digest_batch_id = $1 WHERE id = $2`,
			batchID, id)
		if err != nil {
			return fmt.Errorf("mark sent %s: %w", id, err)
		}
	}
	return tx.Commit(ctx)
}

// ---------------------------------------------------------------------------
// Bounces
// ---------------------------------------------------------------------------

func (s *Store) CreateBounce(ctx context.Context, tenantID uuid.UUID, b *domain.NotificationBounce) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	b.ID = uuid.New()
	b.TenantID = tenantID

	err = tx.QueryRow(ctx,
		`INSERT INTO notification_bounces (id, tenant_id, notification_id, bounce_type, bounce_subtype, email_address, diagnostic)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING created_at`,
		b.ID, tenantID, b.NotificationID, b.BounceType, b.BounceSubtype, b.EmailAddress, b.Diagnostic,
	).Scan(&b.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert bounce: %w", err)
	}
	return tx.Commit(ctx)
}

func (s *Store) MarkEmailInvalid(ctx context.Context, tenantID uuid.UUID, email string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	_, err = tx.Exec(ctx,
		`UPDATE notification_preferences SET email_valid = false, bounce_count = bounce_count + 1, updated_at = now()
		 WHERE email_address = $1`, email)
	if err != nil {
		return fmt.Errorf("mark email invalid: %w", err)
	}
	return tx.Commit(ctx)
}
