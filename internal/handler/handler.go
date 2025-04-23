package handler

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/hibiken/asynq"
	"github.com/rhajizada/gazette/internal/repository"
	"golang.org/x/oauth2"
)

// Handler encapsulates dependencies for HTTP handlers.
type Handler struct {
	Repo     repository.Queries
	Client   *asynq.Client
	Secret   []byte
	Verifier *oidc.IDTokenVerifier
	OAuth    *oauth2.Config
}

// New creates a new Handler.
func New(r *repository.Queries, c *asynq.Client, s []byte, v *oidc.IDTokenVerifier, o *oauth2.Config) *Handler {
	return &Handler{
		Repo:     *r,
		Client:   c,
		Secret:   s,
		Verifier: v,
		OAuth:    o,
	}
}
