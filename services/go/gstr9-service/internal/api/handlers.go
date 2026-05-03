package api

import (
	"github.com/complai/complai/services/go/gstr9-service/internal/store"
)

type Handlers struct {
	store     store.Repository
	gstSvcURL string
}

func NewHandlers(s store.Repository, gstSvcURL string) *Handlers {
	return &Handlers{store: s, gstSvcURL: gstSvcURL}
}
