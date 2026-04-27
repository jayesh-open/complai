package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/complai/complai/services/go/ewb-service/internal/gateway"
	"github.com/complai/complai/services/go/ewb-service/internal/store"
)

func NewRouter(s store.Repository, ewbClient *gateway.EWBClient, clock store.Clock) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}))

	h := NewHandlers(s, ewbClient, clock)

	r.Get("/health", h.Health)
	r.Get("/ping", h.Health)

	r.Route("/v1/ewb", func(r chi.Router) {
		r.Post("/generate", h.GenerateEWB)
		r.Get("/list", h.ListEWBs)
		r.Post("/consolidate", h.ConsolidateEWB)
		r.Get("/number/{ewbNo}", h.GetEWBByNumber)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetEWB)
			r.Post("/cancel", h.CancelEWB)
			r.Post("/vehicle", h.UpdateVehicle)
			r.Post("/extend", h.ExtendValidity)
			r.Get("/vehicles", h.GetVehicleHistory)
		})
	})

	return r
}
