package handler

import (
	"database/sql"
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
		http.Error(w, fmt.Sprintf("invalid paging: %v", err), http.StatusBadRequest)
		return
	}

	resp, err := h.Service.ListUserLikedItems(r.Context(),
		repository.ListUserLikedItemsParams{
			UserID: userID,
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed listing liked items: %v", err), http.StatusInternalServerError)
		return
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
	itemID, err := uuid.Parse(r.PathValue("itemID"))
	if err != nil {
		http.Error(w, "invalid item ID", http.StatusBadRequest)
		return
	}

	item, err := h.Service.GetItem(r.Context(), service.GetItemRequest{UserID: userID, ItemID: itemID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "item not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("failed fetching item: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// LikeItem marks an item as liked by the user.
// @Summary      Like item
// @Description  Creates a like record for the current user on an item.
// @Tags         Items
// @Param        itemID  path      string  true  "Item UUID"
// @Success      200     {object}  service.LikeItemResponse
// @Failure      400     {object}  string
// @Failure      409     {object}  string
// @Failure      500     {object}  string
// @Security     BearerAuth
// @Router       /api/items/{itemID}/like [post]
func (h *Handler) LikeItem(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID
	itemID, err := uuid.Parse(r.PathValue("itemID"))
	if err != nil {
		http.Error(w, "invalid item ID", http.StatusBadRequest)
		return
	}

	likedAt, err := h.Service.LikeItem(r.Context(),
		repository.CreateUserLikeParams{
			UserID: userID,
			ItemID: itemID,
		})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to like item: %v", err), http.StatusConflict)
		return
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
	itemID, err := uuid.Parse(r.PathValue("itemID"))
	if err != nil {
		http.Error(w, "invalid item ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.UnlikeItem(r.Context(),
		repository.DeleteUserLikeParams{
			UserID: userID,
			ItemID: itemID,
		}); err != nil {
		http.Error(w, fmt.Sprintf("failed to unlike item: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
