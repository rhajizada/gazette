package router

import (
	"net/http"

	"github.com/rhajizada/gazette/internal/handler"
)

func RegisterAuthRoutes(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /callback", h.Callback)
	router.HandleFunc("GET /login", h.Login)
	return router
}
