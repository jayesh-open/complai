package outbox

import (
	"context"

	"github.com/jackc/pgx/v5"
)

const insertQuery = `
INSERT INTO outbox (id, aggregate_type, aggregate_id, event_type, payload, target_queue, request_id, status, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`

func PublishInTx(ctx context.Context, tx pgx.Tx, row OutboxRow) error {
	_, err := tx.Exec(ctx, insertQuery,
		row.ID,
		row.AggregateType,
		row.AggregateID,
		row.EventType,
		row.Payload,
		row.TargetQueue,
		row.RequestID,
		row.Status,
		row.CreatedAt,
	)
	return err
}
