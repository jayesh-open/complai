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

	r.Route("/v1/gateway/tds", func(r chi.Router) {
		r.Post("/pan/verify", h.VerifyPAN)
		r.Post("/tan/verify", h.VerifyTAN)
		r.Post("/challan/generate", h.GenerateChallan)
		r.Post("/form26q/file", h.FileForm26Q)
		r.Post("/form24q/file", h.FileForm24Q)
	})

	return r
}
