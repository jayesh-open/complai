package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/complai/complai/services/go/gstn-gateway-service/internal/provider"
)

func NewRouter(p provider.GSTNProvider) chi.Router {
	h := NewHandlers(p)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/health", h.Health)

	r.Route("/v1/gateway/adaequare", func(r chi.Router) {
		r.Post("/authenticate", h.Authenticate)

		r.Route("/gstr1", func(r chi.Router) {
			r.Post("/save", h.GSTR1Save)
			r.Post("/get", h.GSTR1Get)
			r.Post("/reset", h.GSTR1Reset)
			r.Post("/submit", h.GSTR1Submit)
			r.Post("/file", h.GSTR1File)
			r.Post("/status", h.GSTR1Status)
			r.Post("/summary", h.GSTR1SummaryHandler)
		})

		r.Route("/gstr2b", func(r chi.Router) {
			r.Post("/get", h.GSTR2BGet)
		})

		r.Route("/gstr2a", func(r chi.Router) {
			r.Post("/get", h.GSTR2AGet)
		})

		r.Route("/ims", func(r chi.Router) {
			r.Post("/get", h.IMSGet)
			r.Post("/action", h.IMSAction)
			r.Post("/bulk-action", h.IMSBulkAction)
		})

		r.Route("/gstr3b", func(r chi.Router) {
			r.Post("/save", h.GSTR3BSave)
			r.Post("/submit", h.GSTR3BSubmit)
			r.Post("/file", h.GSTR3BFile)
		})

		r.Route("/gstr9", func(r chi.Router) {
			r.Post("/save", h.GSTR9Save)
			r.Post("/submit", h.GSTR9Submit)
			r.Post("/file", h.GSTR9File)
			r.Post("/status", h.GSTR9Status)
		})

		r.Route("/gstr9c", func(r chi.Router) {
			r.Post("/save", h.GSTR9CSave)
			r.Post("/file", h.GSTR9CFile)
			r.Post("/status", h.GSTR9CStatus)
		})
	})

	return r
}
