package router

import (
	"net/http"

	_ "github.com/rhajizada/gazette/docs"
	"github.com/rhajizada/gazette/internal/handler"
)

func RegisterItemsAPI(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /items/", h.ListUserLikedItems)
	router.HandleFunc("GET /items/{itemID}", h.GetItemByID)
	router.HandleFunc("POST /items/{itemID}/like", h.LikeItem)
	router.HandleFunc("DELETE /items/{itemID}/like", h.UnlikeItem)

	return router
}
