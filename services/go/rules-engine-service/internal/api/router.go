package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/complai/complai/services/go/rules-engine-service/internal/store"
)

func NewRouter(s store.Repository) chi.Router {
	h := NewHandlers(s)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/health", h.Health)

	r.Route("/v1/rules", func(r chi.Router) {
		r.Post("/evaluate", h.Evaluate)
		r.Post("/", h.CreateRule)
		r.Get("/", h.ListRules)
		r.Get("/{ruleID}", h.GetRule)
	})

	return r
}
