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

	r.Route("/v1/gateway/itr", func(r chi.Router) {
		r.Post("/pan-aadhaar/check", h.CheckPANAadhaarLink)
		r.Post("/ais/fetch", h.FetchAIS)
		r.Post("/submit", h.SubmitITR)
		r.Post("/itrv/generate", h.GenerateITRV)
		r.Post("/everify/check", h.CheckEVerification)
		r.Post("/refund/status", h.CheckRefundStatus)
	})

	return r
}
