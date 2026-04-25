package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/complai/complai/services/go/workflow-service/internal/store"
)

func NewRouter(s store.Repository, engine WorkflowEngine) chi.Router {
	h := NewHandlers(s, engine)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/health", h.Health)

	r.Route("/v1", func(r chi.Router) {
		r.Post("/workflows/start", h.StartWorkflow)
		r.Get("/workflows", h.ListWorkflows)
		r.Get("/workflows/{instanceID}", h.GetWorkflow)
		r.Post("/workflows/{instanceID}/signal", h.SignalWorkflow)

		r.Get("/tasks", h.ListHumanTasks)
		r.Post("/tasks/{taskID}/complete", h.CompleteTask)

		r.Post("/definitions", h.CreateDefinition)
		r.Get("/definitions", h.ListDefinitions)
	})

	return r
}
