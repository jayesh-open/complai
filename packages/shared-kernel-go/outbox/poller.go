package outbox

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/complai/complai/packages/shared-kernel-go/messagebus"
)

const (
	DefaultPollInterval = 500 * time.Millisecond
	DefaultBatchSize    = 100
)

const selectPendingQuery = `
SELECT id, aggregate_type, aggregate_id, event_type, payload, target_queue, request_id, status, created_at
FROM outbox
WHERE status = 'pending'
ORDER BY created_at ASC
LIMIT $1
FOR UPDATE SKIP LOCKED
`

const markPublishedQuery = `
UPDATE outbox SET status = 'published', published_at = $1 WHERE id = $2
`

const markFailedQuery = `
UPDATE outbox SET status = 'failed' WHERE id = $1
`

type PollerConfig struct {
	PollInterval time.Duration
	BatchSize    int
}

type OutboxPoller struct {
	pool   *pgxpool.Pool
	bus    messagebus.MessageBus
	cfg    PollerConfig
	logger zerolog.Logger
}

func NewOutboxPoller(pool *pgxpool.Pool, bus messagebus.MessageBus, cfg PollerConfig, logger zerolog.Logger) *OutboxPoller {
	if cfg.PollInterval == 0 {
		cfg.PollInterval = DefaultPollInterval
	}
	if cfg.BatchSize == 0 {
		cfg.BatchSize = DefaultBatchSize
	}
	return &OutboxPoller{
		pool:   pool,
		bus:    bus,
		cfg:    cfg,
		logger: logger.With().Str("component", "outbox-poller").Logger(),
	}
}

func (p *OutboxPoller) Start(ctx context.Context) {
	ticker := time.NewTicker(p.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			p.logger.Info().Msg("outbox poller stopping")
			return
		case <-ticker.C:
			if err := p.poll(ctx); err != nil {
				p.logger.Error().Err(err).Msg("outbox poll cycle failed")
			}
		}
	}
}

func (p *OutboxPoller) poll(ctx context.Context) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	rows, err := tx.Query(ctx, selectPendingQuery, p.cfg.BatchSize)
	if err != nil {
		return err
	}
	defer rows.Close()

	var pending []OutboxRow
	for rows.Next() {
		var row OutboxRow
		if err := rows.Scan(
			&row.ID, &row.AggregateType, &row.AggregateID,
			&row.EventType, &row.Payload, &row.TargetQueue,
			&row.RequestID, &row.Status, &row.CreatedAt,
		); err != nil {
			return err
		}
		pending = append(pending, row)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	now := time.Now().UTC()
	for _, row := range pending {
		event := messagebus.Event{
			ID:        row.ID,
			Type:      row.EventType,
			Source:    row.AggregateType,
			TenantID:  uuid.Nil, // tenant_id is embedded in the payload
			Payload:   json.RawMessage(row.Payload),
			Timestamp: row.CreatedAt,
			RequestID: row.RequestID,
		}

		if err := p.bus.Publish(ctx, row.TargetQueue, event); err != nil {
			p.logger.Error().Err(err).
				Str("outbox_id", row.ID.String()).
				Str("event_type", row.EventType).
				Msg("failed to publish outbox event")
			if _, markErr := tx.Exec(ctx, markFailedQuery, row.ID); markErr != nil {
				p.logger.Error().Err(markErr).Str("outbox_id", row.ID.String()).Msg("failed to mark outbox row as failed")
			}
			continue
		}

		if _, err := tx.Exec(ctx, markPublishedQuery, now, row.ID); err != nil {
			p.logger.Error().Err(err).Str("outbox_id", row.ID.String()).Msg("failed to mark outbox row as published")
		}
	}

	return tx.Commit(ctx)
}
