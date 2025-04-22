package router

import (
	"net/http"

	"github.com/rhajizada/gazette/internal/handler"
)

func RegisterItemRoutes(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /items/{itemID}", h.GetItemByID)
	return router
}
