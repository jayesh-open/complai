package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/complai/complai/services/go/master-data-service/internal/store"
)

func NewRouter(s store.Repository) chi.Router {
	h := NewHandlers(s)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/health", h.Health)

	r.Route("/v1", func(r chi.Router) {
		r.Post("/vendors", h.CreateVendor)
		r.Get("/vendors", h.ListVendors)
		r.Get("/vendors/{vendorID}", h.GetVendor)
		r.Put("/vendors/{vendorID}", h.UpdateVendor)

		r.Post("/customers", h.CreateCustomer)
		r.Get("/customers", h.ListCustomers)
		r.Get("/customers/{customerID}", h.GetCustomer)

		r.Post("/items", h.CreateItem)
		r.Get("/items", h.ListItems)
		r.Get("/items/{itemID}", h.GetItem)

		r.Get("/hsn-codes", h.ListHSNCodes)
		r.Post("/hsn-codes", h.CreateHSNCode)

		r.Get("/state-codes", h.ListStateCodes)
	})

	return r
}
