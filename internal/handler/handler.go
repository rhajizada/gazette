package handler

import (
	"github.com/hibiken/asynq"
	"github.com/rhajizada/gazette/internal/repository"
)

// Handler encapsulates dependencies for HTTP handlers.
type Handler struct {
	Repo   repository.Queries
	Client *asynq.Client
	Secret []byte
}

// New creates a new Handler.
func New(r *repository.Queries, c *asynq.Client, secret []byte) *Handler {
	return &Handler{
		Repo:   *r,
		Client: c,
		Secret: secret,
	}
}
