package tenant

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type contextKey struct{}

var ErrMissingTenantID = errors.New("tenant_id not found in context")

func WithTenantContext(ctx context.Context, tenantID uuid.UUID) context.Context {
	return context.WithValue(ctx, contextKey{}, tenantID)
}

func TenantIDFromContext(ctx context.Context) (uuid.UUID, error) {
	v, ok := ctx.Value(contextKey{}).(uuid.UUID)
	if !ok || v == uuid.Nil {
		return uuid.Nil, ErrMissingTenantID
	}
	return v, nil
}

func MustTenantIDFromContext(ctx context.Context) uuid.UUID {
	id, err := TenantIDFromContext(ctx)
	if err != nil {
		panic(err)
	}
	return id
}
