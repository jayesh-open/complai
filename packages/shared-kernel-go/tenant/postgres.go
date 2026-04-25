package tenant

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// SetTenantID sets the tenant_id session variable on a Postgres transaction.
// This must be called within a transaction so that SET LOCAL scopes correctly.
// The RLS policy on every table references current_setting('app.tenant_id').
func SetTenantID(ctx context.Context, tx pgx.Tx, tenantID uuid.UUID) error {
	query := fmt.Sprintf("SET LOCAL app.tenant_id = '%s'", tenantID.String())
	_, err := tx.Exec(ctx, query)
	return err
}
