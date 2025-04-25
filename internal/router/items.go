package router

import (
	"net/http"

	"github.com/rhajizada/gazette/internal/handler"
)

func RegisterItemRoutes(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /items/", h.ListUserLikedItems)
	router.HandleFunc("GET /items/{itemID}", h.GetItemByID)
	router.HandleFunc("PUT /items/{itemID}/like", h.LikeItem)
	router.HandleFunc("PUT /items/{itemID}/unlike", h.UnlikeItem)

	return router
}
