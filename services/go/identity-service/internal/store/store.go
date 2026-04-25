package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/complai/complai/services/go/identity-service/internal/domain"
)

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) Pool() *pgxpool.Pool { return s.pool }

func setTenantID(ctx context.Context, tx pgx.Tx, tenantID uuid.UUID) error {
	_, err := tx.Exec(ctx, fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID.String()))
	return err
}

func (s *Store) CreateUser(ctx context.Context, tenantID uuid.UUID, u *domain.User) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO users (tenant_id, external_id, email, email_verified, first_name, last_name, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, created_at, updated_at`,
		tenantID, u.ExternalID, u.Email, u.EmailVerified, u.FirstName, u.LastName, u.Status,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	u.TenantID = tenantID
	return tx.Commit(ctx)
}

func (s *Store) GetUserByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*domain.User, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var u domain.User
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, external_id, email, email_verified, first_name, last_name, status, created_at, updated_at
		 FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.TenantID, &u.ExternalID, &u.Email, &u.EmailVerified, &u.FirstName, &u.LastName, &u.Status, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &u, tx.Commit(ctx)
}

func (s *Store) GetUserByID(ctx context.Context, tenantID, userID uuid.UUID) (*domain.User, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var u domain.User
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, external_id, email, email_verified, first_name, last_name, status, created_at, updated_at
		 FROM users WHERE id = $1`, userID,
	).Scan(&u.ID, &u.TenantID, &u.ExternalID, &u.Email, &u.EmailVerified, &u.FirstName, &u.LastName, &u.Status, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &u, tx.Commit(ctx)
}

func (s *Store) ListUsers(ctx context.Context, tenantID uuid.UUID) ([]domain.User, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, external_id, email, email_verified, first_name, last_name, status, created_at, updated_at
		 FROM users ORDER BY created_at`)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.TenantID, &u.ExternalID, &u.Email, &u.EmailVerified, &u.FirstName, &u.LastName, &u.Status, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, u)
	}
	return users, tx.Commit(ctx)
}

func (s *Store) CreateSession(ctx context.Context, tenantID uuid.UUID, sess *domain.UserSession) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO user_sessions (tenant_id, user_id, device_info, ip_address, expires_at)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`,
		tenantID, sess.UserID, sess.DeviceInfo, sess.IPAddress, sess.ExpiresAt,
	).Scan(&sess.ID, &sess.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert session: %w", err)
	}
	sess.TenantID = tenantID
	return tx.Commit(ctx)
}

func (s *Store) CreateMFAFactor(ctx context.Context, tenantID uuid.UUID, f *domain.MFAFactor) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO mfa_factors (tenant_id, user_id, factor_type, secret_encrypted, phone_number, verified)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`,
		tenantID, f.UserID, f.FactorType, f.SecretEncrypted, f.PhoneNumber, f.Verified,
	).Scan(&f.ID, &f.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert mfa: %w", err)
	}
	f.TenantID = tenantID
	return tx.Commit(ctx)
}

func (s *Store) GetMFAFactors(ctx context.Context, tenantID, userID uuid.UUID) ([]domain.MFAFactor, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, user_id, factor_type, phone_number, verified, created_at
		 FROM mfa_factors WHERE user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("get mfa factors: %w", err)
	}
	defer rows.Close()

	var factors []domain.MFAFactor
	for rows.Next() {
		var f domain.MFAFactor
		if err := rows.Scan(&f.ID, &f.TenantID, &f.UserID, &f.FactorType, &f.PhoneNumber, &f.Verified, &f.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan mfa: %w", err)
		}
		factors = append(factors, f)
	}
	return factors, tx.Commit(ctx)
}

func (s *Store) CreateStepUpEvent(ctx context.Context, tenantID uuid.UUID, evt *domain.StepUpEvent) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	err = tx.QueryRow(ctx,
		`INSERT INTO step_up_events (tenant_id, user_id, session_id, action_class, verified_at, expires_at, mfa_method)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		tenantID, evt.UserID, evt.SessionID, evt.ActionClass, evt.VerifiedAt, evt.ExpiresAt, evt.MFAMethod,
	).Scan(&evt.ID)
	if err != nil {
		return fmt.Errorf("insert step up: %w", err)
	}
	evt.TenantID = tenantID
	return tx.Commit(ctx)
}

func (s *Store) HasValidStepUp(ctx context.Context, tenantID, userID, sessionID uuid.UUID, actionClass string) (bool, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return false, fmt.Errorf("set tenant: %w", err)
	}

	var count int
	err = tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM step_up_events
		 WHERE user_id = $1 AND session_id = $2 AND action_class = $3 AND expires_at > $4`,
		userID, sessionID, actionClass, time.Now(),
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("check step up: %w", err)
	}
	return count > 0, tx.Commit(ctx)
}
