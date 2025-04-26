package router

import (
	"net/http"

	_ "github.com/rhajizada/gazette/docs"
	"github.com/rhajizada/gazette/internal/handler"
)

func RegisterOAuthRoutes(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /callback", h.Callback)
	router.HandleFunc("GET /login", h.Login)
	return router
}
