package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	jwt.RegisteredClaims
	TenantID    uuid.UUID  `json:"tenant_id"`
	ActivePAN   string     `json:"active_pan,omitempty"`
	ActiveGSTIN string     `json:"active_gstin,omitempty"`
	ActiveTAN   string     `json:"active_tan,omitempty"`
	Roles       []string   `json:"roles,omitempty"`
	Permissions []string   `json:"permissions,omitempty"`
	StepUpAt    *time.Time `json:"step_up_at,omitempty"`
}

type claimsContextKey struct{}

func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey{}, claims)
}

func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	c, ok := ctx.Value(claimsContextKey{}).(*Claims)
	return c, ok
}

func (c *Claims) HasRole(role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func (c *Claims) HasPermission(perm string) bool {
	for _, p := range c.Permissions {
		if p == perm {
			return true
		}
	}
	return false
}

func (c *Claims) IsSteppedUp() bool {
	return c.StepUpAt != nil
}
