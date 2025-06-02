package router

import (
	"net/http"

	_ "github.com/rhajizada/gazette/docs"
	"github.com/rhajizada/gazette/internal/handler"
)

func RegisterAPI(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /categories", h.ListCategories)
	router.HandleFunc("GET /categories/items", h.ListItemsByCategories)
	router.HandleFunc("GET /collections", h.ListCollections)
	router.HandleFunc("POST /collections", h.CreateCollection)
	router.HandleFunc("GET /collections/{collectionID}", h.GetCollectionByID)
	router.HandleFunc("DELETE /collections/{collectionID}", h.DeleteCollectionByID)
	router.HandleFunc("GET /collections/{collectionID}/items", h.ListItemsByCollectionID)
	router.HandleFunc("POST /collections/{collectionID}/item/{itemID}", h.AddItemToCollection)
	router.HandleFunc("DELETE /collections/{collectionID}/item/{itemID}", h.RemoveItemFromCollection)
	router.HandleFunc("GET /items", h.ListUserLikedItems)
	router.HandleFunc("GET /items/{itemID}", h.GetItemByID)
	router.HandleFunc("GET /items/{itemID}/similiar", h.ListSimiliarItemsByID)
	router.HandleFunc("POST /items/{itemID}/like", h.LikeItem)
	router.HandleFunc("DELETE /items/{itemID}/like", h.UnlikeItem)
	router.HandleFunc("GET /items/{itemID}/collections", h.ListItemCollections)
	router.HandleFunc("GET /feeds", h.ListFeeds)
	router.HandleFunc("POST /feeds", h.CreateFeed)
	router.HandleFunc("GET /feeds/export", h.ExportFeeds)
	router.HandleFunc("GET /feeds/{feedID}", h.GetFeedByID)
	router.HandleFunc("DELETE /feeds/{feedID}", h.DeleteFeedByID)
	router.HandleFunc("PUT /feeds/{feedID}/subscribe", h.SubscribeToFeed)
	router.HandleFunc("DELETE /feeds/{feedID}/subscribe", h.UnsubscribeFromFeed)
	router.HandleFunc("GET /feeds/{feedID}/items", h.ListItemsByFeedID)
	router.HandleFunc("GET /subscribed", h.ListUserSubscribedItems)
	router.HandleFunc("GET /suggested", h.ListUserSuggestedItems)
	router.HandleFunc("GET /user", h.GetUser)
	return router
}
