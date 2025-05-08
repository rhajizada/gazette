package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/rhajizada/gazette/internal/middleware"
	"github.com/rhajizada/gazette/internal/repository"
	"github.com/rhajizada/gazette/internal/service"
)

// ListCollections returns the user’s collections.
// @Summary      List collections
// @Description  Retrieves paginated collections for the current user.
// @Tags         Collections
// @Param        limit   query     int32  true   "Max number of collections"
// @Param        offset  query     int32  true   "Number of collections to skip"
// @Success      200     {object}  service.ListCollectionsResponse
// @Failure      400     {object}  string
// @Failure      500     {object}  string
// @Security     BearerAuth
// @Router       /api/collections [get]
func (h *Handler) ListCollections(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID

	params, err := getPageParams(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var resp *service.ListCollectionsResponse

	resp, err = h.Service.ListCollections(r.Context(),
		repository.ListCollectionsByUserParams{
			UserID: userID,
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	if err != nil {
		var serviceErr service.ServiceError
		if errors.As(err, &serviceErr) {
			http.Error(w, serviceErr.Error(), int(serviceErr.Code))
			return
		} else {
			http.Error(w, "failed to list collections", http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// CreateCollection creates a new collection.
// @Summary      Create collection
// @Description  Creates a named collection for the current user.FeedURL@Tags         Collections
// @Tags         Collections
// @Param        body    body      CreateCollectionRequest  true  "Collection name"
// @Success      200     {object}  service.Collection
// @Failure      400     {object}  string
// @Failure      409     {object}  string
// @Failure      500     {object}  string
// @Security     BearerAuth
// @Router       /api/collections [post]
func (h *Handler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID

	var req CreateCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	query := repository.CreateCollectionParams{
		UserID: userID,
		Name:   req.Name,
	}

	col, err := h.Service.CreateCollection(r.Context(), query)
	if err != nil {
		var serviceErr service.ServiceError
		if errors.As(err, &serviceErr) {
			http.Error(w, serviceErr.Error(), int(serviceErr.Code))
			return
		} else {
			http.Error(w, fmt.Sprintf("failed to create collection %s", req.Name), http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(col)
}

// GetCollectionByID retrieves a collection.
// @Summary      Get collection
// @Description  Retrieves a collection by ID.
// @Tags         Collections
// @Param        collectionID  path  string  true  "Collection UUID"
// @Success      200           {object}  service.Collection
// @Failure      400           {object}  string
// @Failure      404           {object}  string
// @Failure      500           {object}  string
// @Security     BearerAuth
// @Router       /api/collections/{collectionID} [get]
func (h *Handler) GetCollectionByID(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("collectionID")
	colID, err := uuid.Parse(path)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s is not a valid id", path), http.StatusBadRequest)
		return
	}

	col, err := h.Service.GetCollectionByID(r.Context(), colID)
	if err != nil {
		var serviceErr service.ServiceError
		if errors.As(err, &serviceErr) {
			http.Error(w, serviceErr.Error(), int(serviceErr.Code))
			return
		} else {
			http.Error(w, fmt.Sprintf("failed to fetch collection %s", colID), http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(col)
}

// DeleteCollectionByID deletes a collection.
// @Summary      Delete collection
// @Description  Deletes a collection by ID.
// @Tags         Collections
// @Param        collectionID  path  string  true  "Collection UUID"
// @Success      204  "No Content"
// @Failure      400  {object}  string
// @Failure      404  {object}  string
// @Failure      500  {object}  string
// @Security     BearerAuth
// @Router       /api/collections/{collectionID} [delete]
func (h *Handler) DeleteCollectionByID(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("collectionID")
	colID, err := uuid.Parse(path)
	if err != nil {
		http.Error(w, fmt.Sprintf("bad input: %s", path), http.StatusBadRequest)
		return
	}

	err = h.Service.DeleteCollection(r.Context(), colID)
	if err != nil {
		var serviceErr service.ServiceError
		if errors.As(err, &serviceErr) {
			http.Error(w, serviceErr.Error(), int(serviceErr.Code))
			return
		} else {
			http.Error(w, fmt.Sprintf("failed to delete collection %s", colID), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddItemToCollection adds an item to a collection.
// @Summary      Add item to collection
// @Description  Adds an item to the specified collection.
// @Tags         Collections
// @Param        collectionID  path  string  true  "Collection UUID"
// @Param        itemID        path  string  true  "Item UUID"
// @Success      200           {object}  service.AddItemToCollectionResponse
// @Failure      400  {object}  string
// @Failure      500  {object}  string
// @Security     BearerAuth
// @Router       /api/collections/{collectionID}/item/{itemID} [post]
func (h *Handler) AddItemToCollection(w http.ResponseWriter, r *http.Request) {
	colPath := r.PathValue("collectionID")
	colID, err := uuid.Parse(colPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s is not a valid id", colPath), http.StatusBadRequest)
		return
	}
	itemPath := r.PathValue("itemID")
	itemID, err := uuid.Parse(itemPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s is not a valid id", itemPath), http.StatusBadRequest)
		return
	}

	resp, err := h.Service.AddItemToCollection(r.Context(),
		repository.AddItemToCollectionParams{
			CollectionID: colID,
			ItemID:       itemID,
		})
	if err != nil {
		var serviceErr service.ServiceError
		if errors.As(err, &serviceErr) {
			http.Error(w, serviceErr.Error(), int(serviceErr.Code))
			return
		} else {
			http.Error(w, fmt.Sprintf("failed to add item %s to collection %s", itemID, colID), http.StatusBadRequest)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RemoveItemFromCollection removes an item from a collection.
// @Summary      Remove item from collection
// @Description  Removes the specified item from the collection.
// @Tags         Collections
// @Param        collectionID  path  string  true  "Collection UUID"
// @Param        itemID        path  string  true  "Item UUID"
// @Success      204  "No Content"
// @Failure      400  {object}  string
// @Failure      500  {object}  string
// @Security     BearerAuth
// @Router       /api/collections/{collectionID}/item/{itemID} [delete]
func (h *Handler) RemoveItemFromCollection(w http.ResponseWriter, r *http.Request) {
	colPath := r.PathValue("collectionID")
	colID, err := uuid.Parse(colPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s is not a valid id", colPath), http.StatusBadRequest)
		return
	}
	itemPath := r.PathValue("itemID")
	itemID, err := uuid.Parse(itemPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s is not a valid id", itemPath), http.StatusBadRequest)
		return
	}

	err = h.Service.RemoveItemFromCollection(r.Context(),
		repository.RemoveItemFromCollectionParams{
			CollectionID: colID,
			ItemID:       itemID,
		})
	if err != nil {
		var serviceErr service.ServiceError
		if errors.As(err, &serviceErr) {
			http.Error(w, serviceErr.Error(), int(serviceErr.Code))
			return
		} else {
			http.Error(w, fmt.Sprintf("failed to remove item %s from collection %s", itemID, colID), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListItemsByCollectionID returns paginated items in a collection.
// @Summary      List items in collection
// @Description  Retrieves items in the collection, including like status.
// @Tags         Collections
// @Param        collectionID  path      string  true   "Collection UUID"
// @Param        limit         query     int32   true   "Max number of items"
// @Param        offset        query     int32   true   "Number of items to skip"
// @Success      200           {object}  service.ListCollectionItemsResponse
// @Failure      400           {object}  string
// @Failure      500           {object}  string
// @Security     BearerAuth
// @Router       /api/collections/{collectionID}/items [get]
func (h *Handler) ListItemsByCollectionID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID
	colPath := r.PathValue("collectionID")
	colID, err := uuid.Parse(colPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s is not a valid id", colPath), http.StatusBadRequest)
		return
	}
	params, err := getPageParams(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var resp *service.ListCollectionItemsResponse

	resp, err = h.Service.ListCollectionItems(r.Context(), service.ListCollectionItemsRequest{
		UserID: userID,
		ListItemsInCollectionParams: repository.ListItemsInCollectionParams{
			CollectionID: colID,
			Limit:        params.Limit,
			Offset:       params.Offset,
		},
	})
	if err != nil {
		var serviceErr service.ServiceError
		if errors.As(err, &serviceErr) {
			http.Error(w, serviceErr.Error(), int(serviceErr.Code))
			return
		} else {
			http.Error(w, fmt.Sprintf("failed to list items in collection %s", colID), http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
