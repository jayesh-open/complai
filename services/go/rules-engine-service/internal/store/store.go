package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/complai/complai/services/go/rules-engine-service/internal/domain"
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

func (s *Store) CreateRule(ctx context.Context, tenantID uuid.UUID, r *domain.Rule) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	r.ID = uuid.New()
	r.TenantID = tenantID
	r.Status = "active"

	err = tx.QueryRow(ctx,
		`INSERT INTO rules (id, tenant_id, category, name, description, version, priority, conditions, actions, effective_from, effective_to, status, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		 RETURNING created_at, updated_at`,
		r.ID, r.TenantID, r.Category, r.Name, r.Description, r.Version, r.Priority,
		r.Conditions, r.Actions, r.EffectiveFrom, r.EffectiveTo, r.Status, r.CreatedBy,
	).Scan(&r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert rule: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) GetRule(ctx context.Context, tenantID uuid.UUID, ruleID uuid.UUID) (*domain.Rule, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	var r domain.Rule
	err = tx.QueryRow(ctx,
		`SELECT id, tenant_id, category, name, description, version, priority, conditions, actions,
		        effective_from, effective_to, status, created_by, created_at, updated_at
		 FROM rules WHERE id = $1`, ruleID,
	).Scan(&r.ID, &r.TenantID, &r.Category, &r.Name, &r.Description, &r.Version, &r.Priority,
		&r.Conditions, &r.Actions, &r.EffectiveFrom, &r.EffectiveTo, &r.Status, &r.CreatedBy,
		&r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get rule: %w", err)
	}

	return &r, tx.Commit(ctx)
}

func (s *Store) ListRules(ctx context.Context, tenantID uuid.UUID, category string) ([]domain.Rule, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return nil, fmt.Errorf("set tenant: %w", err)
	}

	query := `SELECT id, tenant_id, category, name, description, version, priority, conditions, actions,
	                  effective_from, effective_to, status, created_by, created_at, updated_at
	           FROM rules WHERE status = 'active'`
	args := []interface{}{}

	if category != "" {
		query += " AND category = $1"
		args = append(args, category)
	}
	query += " ORDER BY priority ASC, created_at DESC"

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list rules: %w", err)
	}
	defer rows.Close()

	var rules []domain.Rule
	for rows.Next() {
		var r domain.Rule
		if err := rows.Scan(&r.ID, &r.TenantID, &r.Category, &r.Name, &r.Description, &r.Version,
			&r.Priority, &r.Conditions, &r.Actions, &r.EffectiveFrom, &r.EffectiveTo, &r.Status,
			&r.CreatedBy, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan rule: %w", err)
		}
		rules = append(rules, r)
	}

	return rules, tx.Commit(ctx)
}

func (s *Store) UpdateRule(ctx context.Context, tenantID uuid.UUID, r *domain.Rule) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	_, err = tx.Exec(ctx,
		`UPDATE rules SET category=$1, name=$2, description=$3, priority=$4, conditions=$5, actions=$6,
		        effective_from=$7, effective_to=$8, updated_at=now()
		 WHERE id=$9`,
		r.Category, r.Name, r.Description, r.Priority, r.Conditions, r.Actions,
		r.EffectiveFrom, r.EffectiveTo, r.ID)
	if err != nil {
		return fmt.Errorf("update rule: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) DeleteRule(ctx context.Context, tenantID uuid.UUID, ruleID uuid.UUID) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	_, err = tx.Exec(ctx,
		`UPDATE rules SET status='inactive', updated_at=now() WHERE id=$1`, ruleID)
	if err != nil {
		return fmt.Errorf("soft delete rule: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *Store) CreateExecutionLog(ctx context.Context, tenantID uuid.UUID, l *domain.RuleExecutionLog) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := setTenantID(ctx, tx, tenantID); err != nil {
		return fmt.Errorf("set tenant: %w", err)
	}

	l.ID = uuid.New()
	l.TenantID = tenantID

	err = tx.QueryRow(ctx,
		`INSERT INTO rule_execution_logs (id, tenant_id, rule_id, input_data, matched_rules, output, execution_time_ms, trace_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING created_at`,
		l.ID, l.TenantID, l.RuleID, l.InputData, l.MatchedRules, l.Output, l.ExecutionTimeMs, l.TraceID,
	).Scan(&l.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert execution log: %w", err)
	}

	return tx.Commit(ctx)
}
