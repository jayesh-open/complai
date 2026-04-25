package tenant_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/packages/shared-kernel-go/tenant"
)

func TestWithTenantContext_RoundTrip(t *testing.T) {
	id := uuid.New()
	ctx := tenant.WithTenantContext(context.Background(), id)

	got, err := tenant.TenantIDFromContext(ctx)
	require.NoError(t, err)
	assert.Equal(t, id, got)
}

func TestTenantIDFromContext_Missing(t *testing.T) {
	_, err := tenant.TenantIDFromContext(context.Background())
	assert.ErrorIs(t, err, tenant.ErrMissingTenantID)
}

func TestTenantIDFromContext_NilUUID(t *testing.T) {
	ctx := context.WithValue(context.Background(), struct{}{}, uuid.Nil)
	_, err := tenant.TenantIDFromContext(ctx)
	assert.ErrorIs(t, err, tenant.ErrMissingTenantID)
}

func TestMustTenantIDFromContext_Panics(t *testing.T) {
	assert.Panics(t, func() {
		tenant.MustTenantIDFromContext(context.Background())
	})
}

func TestMustTenantIDFromContext_Success(t *testing.T) {
	id := uuid.New()
	ctx := tenant.WithTenantContext(context.Background(), id)
	assert.Equal(t, id, tenant.MustTenantIDFromContext(ctx))
}
