package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *Handlers) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1/tds", func(r chi.Router) {
		r.Get("/deductees", h.ListDeductees)
		r.Get("/deductees/{id}", h.GetDeductee)

		r.Post("/calculate", h.CalculateTDS)

		r.Post("/entries", h.CreateEntry)
		r.Get("/entries", h.ListEntries)
		r.Get("/entries/{id}", h.GetEntry)

		r.Get("/summary", h.GetSummary)
	})

	return r
}
