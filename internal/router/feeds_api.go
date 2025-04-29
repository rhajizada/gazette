package router

import (
	"net/http"

	_ "github.com/rhajizada/gazette/docs"
	"github.com/rhajizada/gazette/internal/handler"
)

func RegisterFeedsAPI(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /feeds", h.ListFeeds)
	router.HandleFunc("POST /feeds/", h.CreateFeed)
	router.HandleFunc("GET /feeds/{feedID}", h.GetFeedByID)
	router.HandleFunc("DELETE /feeds/{feedID}", h.DeleteFeedByID)
	router.HandleFunc("PUT /feeds/{feedID}/subscribe", h.SubscribeToFeed)
	router.HandleFunc("DELETE /feeds/{feedID}/subscribe", h.UnsubscribeFromFeed)
	router.HandleFunc("GET /feeds/{feedID}/items", h.ListItemsByFeedID)
	router.HandleFunc("GET /items/{itemID}", h.GetItemByID)
	return router
}
