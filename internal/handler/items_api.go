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

// ListUserLikedItems returns items the user has liked.
// @Summary      List liked items
// @Description  Retrieves items liked by the user, paginated.
// @Tags         Items
// @Param        limit   query     int32  true   "Max number of items"
// @Param        offset  query     int32  true   "Number of items to skip"
// @Success      200     {object}  service.ListItemsResponse
// @Failure      400     {object}  string
// @Failure      500     {object}  string
// @Security     BearerAuth
// @Router       /api/items [get]
func (h *Handler) ListUserLikedItems(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID

	params, err := getPageParams(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var resp *service.ListItemsResponse
	resp, err = h.Service.ListUserLikedItems(r.Context(),
		repository.ListUserLikedItemsParams{
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
			http.Error(w, "failed to list liked items", http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetItemByID returns a single item.
// @Summary      Get item
// @Description  Retrieves an item by ID, including like status for the user.
// @Tags         Items
// @Param        itemID  path      string  true  "Item UUID"
// @Success      200     {object}  service.Item
// @Failure      400     {object}  string
// @Failure      404     {object}  string
// @Failure      500     {object}  string
// @Security     BearerAuth
// @Router       /api/items/{itemID} [get]
func (h *Handler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID
	itemPath := r.PathValue("itemID")
	itemID, err := uuid.Parse(itemPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s is not a valid id", itemPath), http.StatusBadRequest)
		return
	}

	item, err := h.Service.GetItem(r.Context(), service.GetItemRequest{UserID: userID, ItemID: itemID})
	if err != nil {
		var serviceErr service.ServiceError
		if errors.As(err, &serviceErr) {
			http.Error(w, serviceErr.Error(), int(serviceErr.Code))
			return
		} else {
			http.Error(w, fmt.Sprintf("failed to fetch item %s", itemID), http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// ListSimiliarItemsByID returns paginated list of similiar items.
// @Summary      List similiar items
// @Description  Retrieves paginated list of similiar items.
// @Tags         Items
// @Param        itemID  path      string  true  "Item UUID"
// @Param        limit   query     int32   true  "Max number of items"
// @Param        offset  query     int32   true  "Number of items to skip"
// @Success      200     {object}  service.ListItemsResponse
// @Failure      400     {object}  string
// @Failure      404     {object}  string
// @Failure      500     {object}  string
// @Security     BearerAuth
// @Router       /api/items/{itemID}/similiar [get]
func (h *Handler) ListSimiliarItemsByID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID
	itemPath := r.PathValue("itemID")
	itemID, err := uuid.Parse(itemPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s is not a valid id", itemPath), http.StatusBadRequest)
		return
	}

	params, err := getPageParams(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	items, err := h.Service.ListSimiliarItemsByID(r.Context(), repository.ListSimilarItemsByItemIDForUserParams{
		ItemID: itemID,
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
			http.Error(w, fmt.Sprintf("failed to fetch similiar items %s", itemID), http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// LikeItem marks an item as liked by the user.
// @Summary      Like item
// @Description  Creates a like record for the current user on an item.
// @Tags         Items
// @Param        itemID  path      string  true  "Item UUID"
// @Success      200     {object}  service.LikeItemResponse
// @Failure      400     {object}  string
// @Failure      500     {object}  string
// @Security     BearerAuth
// @Router       /api/items/{itemID}/like [post]
func (h *Handler) LikeItem(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID
	itemPath := r.PathValue("itemID")
	itemID, err := uuid.Parse(itemPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s is not a valid id", itemPath), http.StatusBadRequest)
		return
	}

	likedAt, err := h.Service.LikeItem(r.Context(),
		repository.CreateUserLikeParams{
			UserID: userID,
			ItemID: itemID,
		})
	if err != nil {
		var serviceErr service.ServiceError
		if errors.As(err, &serviceErr) {
			http.Error(w, serviceErr.Error(), int(serviceErr.Code))
			return
		} else {
			http.Error(w, fmt.Sprintf("failed to like item %s", itemID), http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(likedAt)
}

// UnlikeItem removes a like from an item.
// @Summary      Unlike item
// @Description  Deletes the like record for the current user on an item.
// @Tags         Items
// @Param        itemID  path      string  true  "Item UUID"
// @Success      204     "No Content"
// @Failure      400     {object}  string
// @Failure      500     {object}  string
// @Security     BearerAuth
// @Router       /api/items/{itemID}/like [delete]
func (h *Handler) UnlikeItem(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID
	itemPath := r.PathValue("itemID")
	itemID, err := uuid.Parse(itemPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s is not a valid id", itemPath), http.StatusBadRequest)
		return
	}

	err = h.Service.UnlikeItem(r.Context(),
		repository.DeleteUserLikeParams{
			UserID: userID,
			ItemID: itemID,
		})
	if err != nil {
		var serviceErr service.ServiceError
		if errors.As(err, &serviceErr) {
			http.Error(w, serviceErr.Error(), int(serviceErr.Code))
			return
		} else {
			http.Error(w, fmt.Sprintf("failed to like item %s", itemID), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListItemCollections returns list of collections that item is in.
// @Summary      Get collections item is in.
// @Description  Retrieves list of collections that given item is in.
// @Tags         Items
// @Param        itemID  path      string  true  "Item UUID"
// @Param        limit   query     int32  true   "Max number of items"
// @Param        offset  query     int32  true   "Number of items to skip"
// @Success      200     {object}  service.ListCollectionsResponse
// @Failure      400     {object}  string
// @Failure      500     {object}  string
// @Security     BearerAuth
// @Router       /api/items/{itemID}/collections [get]
func (h *Handler) ListItemCollections(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID
	itemID, err := uuid.Parse(r.PathValue("itemID"))
	if err != nil {
		http.Error(w, "%s is not a valid id", http.StatusBadRequest)
		return
	}

	params, err := getPageParams(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var resp *service.ListCollectionsResponse

	resp, err = h.Service.ListItemCollections(r.Context(),
		repository.ListCollectionsByItemIDParams{
			UserID: userID,
			ItemID: itemID,
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	if err != nil {
		var serviceErr service.ServiceError
		if errors.As(err, &serviceErr) {
			http.Error(w, serviceErr.Error(), int(serviceErr.Code))
			return
		} else {
			http.Error(w, fmt.Sprintf("failed to list collections with item %s", itemID), http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
