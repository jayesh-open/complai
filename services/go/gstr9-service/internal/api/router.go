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

	r.Route("/api/v1/gstr9", func(r chi.Router) {
		r.Post("/annual-return", h.CreateAnnualReturn)
		r.Get("/annual-return", h.ListAnnualReturns)
		r.Get("/annual-return/{id}", h.GetAnnualReturn)
		r.Put("/annual-return/{id}/save", h.SaveAnnualReturn)
		r.Post("/annual-return/{id}/aggregate", h.AggregateAnnualReturn)
		r.Get("/annual-return/{id}/table/{table}", h.GetTableData)

		r.Post("/reconciliation/{gstr9Id}", h.InitiateReconciliation)
		r.Get("/reconciliation/{id}", h.GetReconciliation)
		r.Get("/reconciliation/{id}/mismatches", h.ListReconciliationMismatches)
		r.Put("/reconciliation/{id}/mismatch/{mismatchId}/resolve", h.ResolveMismatch)
		r.Put("/reconciliation/{id}/certify", h.CertifyReconciliation)
	})

	return r
}
