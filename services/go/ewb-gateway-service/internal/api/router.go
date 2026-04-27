package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/complai/complai/services/go/ewb-gateway-service/internal/provider"
)

func NewRouter(p provider.EWBProvider) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}))

	h := NewHandlers(p)

	r.Get("/health", h.Health)
	r.Get("/ping", h.Health)

	r.Route("/v1/gateway/ewb", func(r chi.Router) {
		r.Post("/generate", h.GenerateEWB)
		r.Post("/cancel", h.CancelEWB)
		r.Get("/", h.GetEWB)
		r.Post("/vehicle", h.UpdateVehicle)
		r.Post("/extend", h.ExtendValidity)
		r.Post("/consolidate", h.ConsolidateEWB)
	})

	return r
}
