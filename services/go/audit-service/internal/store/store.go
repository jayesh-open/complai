package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/complai/complai/services/go/audit-service/internal/domain"
)

// compile-time check
var _ Repository = (*Store)(nil)

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

func (s *Store) CreateEvent(ctx context.Context, tenantID uuid.UUID, e *domain.AuditEvent) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	e.ID = uuid.New()
	e.TenantID = tenantID

	err = tx.QueryRow(ctx,
		`INSERT INTO audit_events (id, tenant_id, user_id, resource_type, resource_id, action, old_value, new_value, status, error_message, ip_address, user_agent, trace_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		 RETURNING created_at`,
		e.ID, e.TenantID, e.UserID, e.ResourceType, e.ResourceID, e.Action,
		e.OldValue, e.NewValue, e.Status, e.ErrorMessage, e.IPAddress, e.UserAgent, e.TraceID,
	).Scan(&e.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert audit event: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) ListEvents(ctx context.Context, tenantID uuid.UUID, params domain.QueryParams) ([]domain.AuditEvent, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	query := `SELECT id, tenant_id, user_id, resource_type, resource_id, action, old_value, new_value, status, error_message, ip_address, user_agent, trace_id, created_at
		 FROM audit_events WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if params.ResourceType != "" {
		query += fmt.Sprintf(" AND resource_type = $%d", argIdx)
		args = append(args, params.ResourceType)
		argIdx++
	}
	if params.Action != "" {
		query += fmt.Sprintf(" AND action = $%d", argIdx)
		args = append(args, params.Action)
		argIdx++
	}
	if params.DateFrom != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argIdx)
		args = append(args, *params.DateFrom)
		argIdx++
	}
	if params.DateTo != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argIdx)
		args = append(args, *params.DateTo)
		argIdx++
	}

	query += " ORDER BY created_at DESC"

	if params.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, params.Limit)
		argIdx++
	}
	if params.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, params.Offset)
		argIdx++
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}
	defer rows.Close()

	var events []domain.AuditEvent
	for rows.Next() {
		var e domain.AuditEvent
		if err := rows.Scan(&e.ID, &e.TenantID, &e.UserID, &e.ResourceType, &e.ResourceID, &e.Action,
			&e.OldValue, &e.NewValue, &e.Status, &e.ErrorMessage, &e.IPAddress, &e.UserAgent, &e.TraceID, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}
		events = append(events, e)
	}

	return events, tx.Commit(ctx)
}

func (s *Store) GetEvent(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.AuditEvent, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var e domain.AuditEvent
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, user_id, resource_type, resource_id, action, old_value, new_value, status, error_message, ip_address, user_agent, trace_id, created_at
		 FROM audit_events WHERE id = $1`, id,
	).Scan(&e.ID, &e.TenantID, &e.UserID, &e.ResourceType, &e.ResourceID, &e.Action,
		&e.OldValue, &e.NewValue, &e.Status, &e.ErrorMessage, &e.IPAddress, &e.UserAgent, &e.TraceID, &e.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get event: %w", err)
	}

	return &e, tx.Commit(ctx)
}

func (s *Store) GetEventsForHour(ctx context.Context, tenantID uuid.UUID, hourBucket time.Time) ([]domain.AuditEvent, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	hourStart := hourBucket.Truncate(time.Hour)
	hourEnd := hourStart.Add(time.Hour)

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, user_id, resource_type, resource_id, action, old_value, new_value, status, error_message, ip_address, user_agent, trace_id, created_at
		 FROM audit_events WHERE created_at >= $1 AND created_at < $2 ORDER BY created_at`,
		hourStart, hourEnd,
	)
	if err != nil {
		return nil, fmt.Errorf("get events for hour: %w", err)
	}
	defer rows.Close()

	var events []domain.AuditEvent
	for rows.Next() {
		var e domain.AuditEvent
		if err := rows.Scan(&e.ID, &e.TenantID, &e.UserID, &e.ResourceType, &e.ResourceID, &e.Action,
			&e.OldValue, &e.NewValue, &e.Status, &e.ErrorMessage, &e.IPAddress, &e.UserAgent, &e.TraceID, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}
		events = append(events, e)
	}

	return events, tx.Commit(ctx)
}

func (s *Store) CreateMerkleChain(ctx context.Context, tenantID uuid.UUID, m *domain.MerkleChain) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	m.ID = uuid.New()
	m.TenantID = tenantID

	err = tx.QueryRow(ctx,
		`INSERT INTO merkle_chains (id, tenant_id, hour_bucket, event_count, hash_payload, previous_hash, computed_hash)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING created_at`,
		m.ID, m.TenantID, m.HourBucket, m.EventCount, m.HashPayload, m.PreviousHash, m.ComputedHash,
	).Scan(&m.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert merkle chain: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) GetMerkleChains(ctx context.Context, tenantID uuid.UUID, from, to time.Time) ([]domain.MerkleChain, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	rows, err := tx.Query(ctx,
		`SELECT id, tenant_id, hour_bucket, event_count, hash_payload, previous_hash, computed_hash, created_at
		 FROM merkle_chains WHERE hour_bucket >= $1 AND hour_bucket <= $2 ORDER BY hour_bucket`,
		from, to,
	)
	if err != nil {
		return nil, fmt.Errorf("get merkle chains: %w", err)
	}
	defer rows.Close()

	var chains []domain.MerkleChain
	for rows.Next() {
		var m domain.MerkleChain
		if err := rows.Scan(&m.ID, &m.TenantID, &m.HourBucket, &m.EventCount, &m.HashPayload, &m.PreviousHash, &m.ComputedHash, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan merkle chain: %w", err)
		}
		chains = append(chains, m)
	}

	return chains, tx.Commit(ctx)
}

func (s *Store) GetLatestMerkleChain(ctx context.Context, tenantID uuid.UUID) (*domain.MerkleChain, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var m domain.MerkleChain
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, hour_bucket, event_count, hash_payload, previous_hash, computed_hash, created_at
		 FROM merkle_chains ORDER BY hour_bucket DESC LIMIT 1`,
	).Scan(&m.ID, &m.TenantID, &m.HourBucket, &m.EventCount, &m.HashPayload, &m.PreviousHash, &m.ComputedHash, &m.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get latest merkle chain: %w", err)
	}

	return &m, tx.Commit(ctx)
}

func (s *Store) UpdateEventNewValue(ctx context.Context, tenantID uuid.UUID, eventID uuid.UUID, newValue string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	_, err = tx.Exec(ctx,
		`UPDATE audit_events SET new_value = $1 WHERE id = $2`,
		newValue, eventID)
	if err != nil {
		return fmt.Errorf("update event new_value: %w", err)
	}

	return tx.Commit(ctx)
}
