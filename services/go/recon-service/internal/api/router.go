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

	r.Route("/v1/recon", func(r chi.Router) {
		r.Post("/run", h.RunRecon)
		r.Get("/run/{run_id}", h.GetRun)
		r.Get("/matches", h.ListMatches)
		r.Post("/matches/{match_id}/accept", h.AcceptMatch)
		r.Post("/matches/bulk-accept", h.BulkAccept)
		r.Post("/ims/action", h.IMSActionHandler)
		r.Get("/ims", h.GetIMSState)
	})

	return r
}
