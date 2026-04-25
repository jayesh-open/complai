package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/complai/complai/services/go/identity-service/internal/store"
)

func NewRouter(s store.Repository, keycloakURL, clientID, clientSec string) chi.Router {
	h := NewHandlers(s, keycloakURL, clientID, clientSec)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/health", h.Health)

	r.Route("/v1/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/refresh", h.Refresh)
		r.Post("/logout", h.Logout)
		r.Post("/step-up/check", h.StepUpCheck)
		r.Post("/step-up/verify", h.StepUpVerify)
		r.Post("/mfa/enroll", h.EnrollMFA)
	})

	r.Route("/v1/users", func(r chi.Router) {
		r.Get("/", h.ListUsers)
		r.Get("/{userID}", h.GetUser)
	})

	return r
}
