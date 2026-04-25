package httputil

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

const (
	DefaultRequestSizeLimit = 10 << 20 // 10 MB
	DefaultTimeout          = 30 * time.Second
)

func StandardMiddleware(logger zerolog.Logger) chi.Middlewares {
	return chi.Middlewares{
		RequestIDMiddleware,
		chimiddleware.RealIP,
		hlog.NewHandler(logger),
		hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Stringer("url", r.URL).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg("request")
		}),
		chimiddleware.Recoverer,
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Tenant-Id", "X-Request-Id"},
			ExposedHeaders:   []string{"X-Request-Id"},
			AllowCredentials: true,
			MaxAge:           300,
		}),
		RequestSizeLimitMiddleware(DefaultRequestSizeLimit),
		chimiddleware.Timeout(DefaultTimeout),
	}
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		w.Header().Set("X-Request-Id", requestID)

		ctx := r.Context()
		r = r.WithContext(ctx)
		r.Header.Set("X-Request-Id", requestID)

		next.ServeHTTP(w, r)
	})
}

func RequestSizeLimitMiddleware(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body != nil {
				r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			}
			next.ServeHTTP(w, r)
		})
	}
}
