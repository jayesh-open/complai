package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/complai/complai/services/go/tenant-service/internal/store"
)

func NewRouter(s store.Repository, kmsClient KMSClient) chi.Router {
	h := NewHandlers(s, kmsClient)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/health", h.Health)

	r.Route("/v1/tenants", func(r chi.Router) {
		r.Post("/", h.CreateTenant)
		r.Get("/", h.ListTenants)
		r.Get("/{tenantID}", h.GetTenant)
		r.Post("/{tenantID}/suspend", h.SuspendTenant)
		r.Post("/{tenantID}/reactivate", h.ReactivateTenant)
		r.Get("/{tenantID}/hierarchy", h.GetHierarchy)
		r.Post("/{tenantID}/pans", h.CreatePAN)
		r.Post("/{tenantID}/pans/{panID}/gstins", h.CreateGSTIN)
		r.Post("/{tenantID}/pans/{panID}/tans", h.CreateTAN)
	})

	return r
}
