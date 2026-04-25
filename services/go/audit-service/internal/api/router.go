package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/complai/complai/services/go/audit-service/internal/store"
)

func NewRouter(s store.Repository) chi.Router {
	h := NewHandlers(s)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/health", h.Health)

	r.Route("/v1/audit", func(r chi.Router) {
		r.Post("/events", h.CreateEvent)
		r.Get("/events", h.ListEvents)
		r.Get("/events/{eventID}", h.GetEvent)
		r.Post("/merkle/compute", h.ComputeMerkleHash)
		r.Get("/integrity", h.IntegrityCheck)
	})

	return r
}
