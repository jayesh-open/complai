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

	r.Route("/v1/gst/gstr1", func(r chi.Router) {
		r.Post("/ingest", h.Ingest)
		r.Post("/validate", h.Validate)
		r.Post("/approve", h.Approve)
		r.Post("/file", h.File)
		r.Get("/summary", h.Summary)
		r.Get("/entries", h.ListEntries)
		r.Get("/errors", h.ListErrors)
	})

	r.Route("/v1/gst/gstr3b", func(r chi.Router) {
		r.Post("/auto-fill", h.GSTR3BAutoFill)
		r.Get("/summary", h.GSTR3BSummary)
		r.Post("/approve", h.GSTR3BApprove)
		r.Post("/file", h.GSTR3BFile)
	})

	return r
}
