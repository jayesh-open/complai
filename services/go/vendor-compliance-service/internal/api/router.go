package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *Handlers) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/health", h.Health)

	r.Route("/v1/vendor-compliance", func(r chi.Router) {
		r.Post("/sync", h.TriggerSync)
		r.Get("/vendors", h.ListVendors)
		r.Get("/vendors/{vendorId}/score", h.GetVendorScore)
		r.Get("/sync/status", h.GetSyncStatus)
		r.Get("/summary", h.GetScoreSummary)
	})

	return r
}
