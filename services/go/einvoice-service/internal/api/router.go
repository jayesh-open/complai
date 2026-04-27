package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/complai/complai/services/go/einvoice-service/internal/gateway"
	"github.com/complai/complai/services/go/einvoice-service/internal/store"
)

func NewRouter(s store.Repository, irp *gateway.IRPClient, clock store.Clock) chi.Router {
	h := NewHandlers(s, irp, clock)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/health", h.Health)

	r.Route("/v1/einvoice", func(r chi.Router) {
		r.Post("/generate", h.GenerateIRN)
		r.Get("/list", h.ListEInvoices)
		r.Get("/summary", h.GetSummary)
		r.Get("/irn/{irn}", h.GetEInvoiceByIRN)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetEInvoice)
			r.Post("/cancel", h.CancelIRN)
		})
	})

	return r
}
