package tenant

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tenantID uuid.UUID

		if claims, ok := claimsFromRequest(r); ok {
			if tid, err := extractTenantIDFromClaims(claims); err == nil {
				tenantID = tid
			}
		}

		if tenantID == uuid.Nil {
			if header := r.Header.Get("X-Tenant-Id"); header != "" {
				if parsed, err := uuid.Parse(header); err == nil {
					tenantID = parsed
				}
			}
		}

		if tenantID == uuid.Nil {
			http.Error(w, `{"error":"missing tenant_id"}`, http.StatusBadRequest)
			return
		}

		ctx := WithTenantContext(r.Context(), tenantID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func claimsFromRequest(r *http.Request) (jwt.MapClaims, bool) {
	tokenStr := r.Header.Get("Authorization")
	if len(tokenStr) <= 7 || tokenStr[:7] != "Bearer " {
		return nil, false
	}
	tokenStr = tokenStr[7:]

	// Parse without validation — the auth middleware upstream handles full validation.
	// This middleware only needs to extract the tenant_id claim.
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	token, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return nil, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	return claims, ok
}

func extractTenantIDFromClaims(claims jwt.MapClaims) (uuid.UUID, error) {
	v, ok := claims["tenant_id"]
	if !ok {
		return uuid.Nil, ErrMissingTenantID
	}
	s, ok := v.(string)
	if !ok {
		return uuid.Nil, ErrMissingTenantID
	}
	return uuid.Parse(s)
}
