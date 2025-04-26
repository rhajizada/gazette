package handler

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/rhajizada/gazette/internal/service"
	"golang.org/x/oauth2"
)

// Handler encapsulates dependencies for HTTP handlers.
type Handler struct {
	Service  *service.Service
	Secret   []byte
	Verifier *oidc.IDTokenVerifier
	OAuth    *oauth2.Config
}

// New creates a new Handler.
func New(service *service.Service, jwtSecret []byte, jwtVerifier *oidc.IDTokenVerifier, oauthConfig *oauth2.Config) *Handler {
	return &Handler{
		Service:  service,
		Secret:   jwtSecret,
		Verifier: jwtVerifier,
		OAuth:    oauthConfig,
	}
}
