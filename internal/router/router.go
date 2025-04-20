package router

import (
	"net/http"

	"github.com/rhajizada/gazette/internal/handler"
)

func RegisterFeedRoutes(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /feeds", h.ListFeeds)
	router.HandleFunc("POST /feeds", h.CreateFeed)
	router.HandleFunc("GET /feeds/{feedID}", h.GetFeedByID)
	router.HandleFunc("DELETE /feeds/{feedID}", h.DeleteFeedByID)
	return router
}
