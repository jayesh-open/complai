package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/complai/complai/services/go/notification-service/internal/store"
)

func NewRouter(s store.Repository, emailSender EmailSender) chi.Router {
	h := NewHandlers(s, emailSender)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/health", h.Health)

	r.Post("/v1/notifications/send", h.SendNotification)
	r.Get("/v1/notifications", h.ListNotifications)
	r.Get("/v1/notifications/{notificationID}", h.GetNotification)
	r.Post("/v1/notifications/bounce", h.ProcessBounce)
	r.Post("/v1/notifications/digest", h.SendDigest)

	r.Get("/v1/users/{userID}/preferences", h.GetPreferences)
	r.Put("/v1/users/{userID}/preferences", h.UpdatePreferences)

	r.Post("/v1/templates", h.CreateTemplate)
	r.Get("/v1/templates", h.ListTemplates)

	return r
}
