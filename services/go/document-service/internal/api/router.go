package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/complai/complai/services/go/document-service/internal/store"
)

func NewRouter(s store.Repository, kmsClient KMSClient, s3Client S3Client, bucket string) chi.Router {
	h := NewHandlers(s, kmsClient, s3Client, bucket)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/health", h.Health)

	r.Route("/v1/documents", func(r chi.Router) {
		r.Post("/upload", h.Upload)
		r.Get("/", h.ListDocuments)
		r.Get("/{documentID}", h.GetDocument)
		r.Get("/{documentID}/download", h.Download)
		r.Post("/{documentID}/ocr", h.TriggerOCR)
	})

	return r
}
