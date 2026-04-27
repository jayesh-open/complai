package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/complai/complai/services/go/irp-gateway-service/internal/provider"
)

func NewRouter(p provider.IRPProvider) chi.Router {
	h := NewHandlers(p)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/health", h.Health)

	r.Route("/v1/gateway/irp", func(r chi.Router) {
		r.Post("/authenticate", h.Authenticate)

		r.Route("/invoice", func(r chi.Router) {
			r.Post("/", h.GenerateIRN)
			r.Post("/cancel", h.CancelIRN)
			r.Get("/irn", h.GetIRNByIRN)
			r.Get("/irnbydocdetails", h.GetIRNByDoc)
		})

		r.Route("/master", func(r chi.Router) {
			r.Get("/gstin", h.ValidateGSTIN)
		})
	})

	return r
}
