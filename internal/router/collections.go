package router

import (
	"net/http"

	"github.com/rhajizada/gazette/internal/handler"
)

func RegisterCollectionRoutes(h *handler.Handler) *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /collections/", h.ListCollections)
	router.HandleFunc("POST /collections/", h.CreateCollection)
	router.HandleFunc("GET /collections/{collectionID}", h.GetCollectionByID)
	router.HandleFunc("DELETE /collections/{collectionID}", h.DeleteCollectionByID)
	router.HandleFunc("GET /collections/{collectionID}/items", h.ListItemsByCollectionID)
	router.HandleFunc("PUT /collections/{collectionID}/item/{itemID}", h.AddItemToCollection)
	router.HandleFunc("DELETE /collections/{collectionID}/item/{itemID}", h.RemoveItemFromCollection)
	return router
}
