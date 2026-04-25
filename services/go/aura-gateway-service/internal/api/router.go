package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/complai/complai/services/go/aura-gateway-service/internal/provider"
)

func NewRouter(p provider.AuraProvider) chi.Router {
	h := NewHandlers(p)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/health", h.Health)

	r.Route("/v1/gateway/aura", func(r chi.Router) {
		r.Get("/invoices", h.ListARInvoices)
	})

	return r
}
