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

	r.Route("/api/v1/itr", func(r chi.Router) {
		r.Post("/taxpayers", h.CreateTaxpayer)
		r.Get("/taxpayers", h.ListTaxpayers)
		r.Get("/taxpayers/{id}", h.GetTaxpayer)

		r.Post("/filings", h.CreateFiling)
		r.Get("/filings", h.ListFilings)
		r.Get("/filings/{id}", h.GetFiling)

		r.Post("/compute-tax", h.ComputeTax)

		r.Post("/filings/{id}/income", h.AddIncomeEntry)
		r.Get("/filings/{id}/income", h.ListIncomeEntries)

		r.Post("/filings/{id}/deductions", h.AddDeduction)
		r.Get("/filings/{id}/deductions", h.ListDeductions)

		r.Get("/filings/{id}/computation", h.GetTaxComputation)

		r.Post("/filings/{id}/tds-credits", h.AddTDSCredit)
		r.Get("/filings/{id}/tds-credits", h.ListTDSCredits)

		r.Post("/reconcile-tds", h.ReconcileTDS)
		r.Post("/reconcile-ais", h.ReconcileAIS)

		r.Get("/eligibility/itr1", h.CheckITR1Eligibility)
		r.Get("/eligibility/itr2", h.CheckITR2Eligibility)
		r.Get("/eligibility/itr3", h.CheckITR3Eligibility)
		r.Get("/eligibility/itr4", h.CheckITR4Eligibility)
		r.Get("/eligibility/itr5", h.CheckITR5Eligibility)
		r.Get("/eligibility/itr6", h.CheckITR6Eligibility)
		r.Get("/eligibility/itr7", h.CheckITR7Eligibility)

		r.Post("/bulk/batches", h.CreateBulkBatch)
		r.Get("/bulk/batches", h.ListBulkBatches)
		r.Get("/bulk/batches/{batchId}", h.GetBulkBatch)
		r.Post("/bulk/batches/{batchId}/employees", h.AddBulkEmployee)
		r.Get("/bulk/batches/{batchId}/employees", h.ListBulkEmployees)
		r.Post("/bulk/batches/{batchId}/process", h.ProcessBulkBatch)
	})

	return r
}
