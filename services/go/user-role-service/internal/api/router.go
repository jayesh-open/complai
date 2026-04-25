package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/complai/complai/services/go/user-role-service/internal/store"
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
		r.Get("/roles", h.ListRoles)
		r.Post("/roles", h.CreateRole)
		r.Post("/roles/{roleID}/permissions", h.AssignPermissions)
		r.Post("/users/{userID}/roles", h.AssignRole)
		r.Post("/policy/check", h.PolicyCheck)
		r.Get("/templates", h.ListTemplates)
		r.Post("/approvals", h.CreateApproval)
		r.Get("/approvals", h.ListApprovals)
		r.Patch("/approvals/{approvalID}", h.DecideApproval)
	})

	return r
}
