package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/complai/complai/services/go/kyc-gateway-service/internal/provider"
)

func NewRouter(p provider.KYCProvider) chi.Router {
	h := NewHandlers(p)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/health", h.Health)

	r.Route("/v1/gateway/kyc", func(r chi.Router) {
		r.Post("/pan/verify", h.VerifyPAN)
		r.Post("/gstin/verify", h.VerifyGSTIN)
		r.Post("/tan/verify", h.VerifyTAN)
		r.Post("/bank/verify", h.VerifyBank)
	})

	return r
}
