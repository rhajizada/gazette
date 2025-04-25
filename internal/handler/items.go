package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rhajizada/gazette/internal/middleware"
	"github.com/rhajizada/gazette/internal/repository"
)

// ListUserLikedItems returns items the user has liked.
// @Summary      List liked items
// @Description  Retrieves items liked by the user, paginated.
// @Tags         Items
// @Param        limit   query     int32  true   "Max number of items"
// @Param        offset  query     int32  true   "Number of items to skip"
// @Success      200     {object}  ListItemsResponse
// @Failure      400     {object}	 string
// @Failure      500     {object}	 string
// @Security     BearerAuth
// @Router       /api/items/ [get]
func (h *Handler) ListUserLikedItems(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	userID := claims.UserID

	// parse pagination parameters
	params, err := getPageParams(r.URL.Query())
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid paging: %v", err), http.StatusBadRequest)
		return
	}

	// total number of liked items
	total, err := h.Repo.CountLikedItems(r.Context(), userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed counting liked items: %v", err), http.StatusInternalServerError)
		return
	}

	// fetch liked items with their liked_at timestamps
	rows, err := h.Repo.ListUserLikedItems(r.Context(), repository.ListUserLikedItemsParams{
		UserID: userID,
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed listing liked items: %v", err), http.StatusInternalServerError)
		return
	}

	// map to API model
	items := make([]Item, len(rows))
	for i, row := range rows {
		authors := make([]Person, len(row.Authors))
		for i, v := range row.Authors {
			authors[i] = Person{
				Name:  v.Name,
				Email: v.Email,
			}
		}
		items[i] = Item{
			ID:              row.ID,
			FeedID:          row.FeedID,
			Title:           row.Title,
			Description:     row.Description,
			Content:         row.Content,
			Link:            row.Link,
			Links:           row.Links,
			UpdatedParsed:   row.UpdatedParsed,
			PublishedParsed: row.PublishedParsed,
			Authors:         Authors(authors),
			GUID:            row.Guid,
			Image:           row.Image,
			Categories:      row.Categories,
			Enclosures:      row.Enclosures,
			CreatedAt:       row.CreatedAt,
			UpdatedAt:       row.UpdatedAt,
			Liked:           row.Liked,
			LikedAt:         &row.LikedAt,
		}
	}

	resp := ListItemsResponse{
		Limit:      params.Limit,
		Offset:     params.Offset,
		TotalCount: total,
		Items:      items,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetItemByID returns a single item.
// @Summary      Get item
// @Description  Retrieves an item by ID, including like status for the user.
// @Tags         Items
// @Param        itemID  path      string  true  "Item UUID"
// @Success      200     {object}  Item
// @Failure      400     {object}	 string
// @Failure      404     {object}	 string
// @Failure      500     {object}	 string
// @Security     BearerAuth
// @Router       /api/items/{itemID} [get]
func (h *Handler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	userID := claims.UserID

	// parse item ID
	idStr := r.PathValue("itemID")
	itemID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid item ID", http.StatusBadRequest)
		return
	}

	// fetch the item
	row, err := h.Repo.GetItemByID(r.Context(), itemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "item not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("failed fetching item: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// check like status
	liked := false
	var likedAt *time.Time
	if like, err := h.Repo.GetUserLike(r.Context(), repository.GetUserLikeParams{
		UserID: userID,
		ItemID: itemID,
	}); err == nil {
		liked = true
		likedAt = &like.LikedAt
	}

	authors := make([]Person, len(row.Authors))
	for i, v := range row.Authors {
		authors[i] = Person{
			Name:  v.Name,
			Email: v.Email,
		}
	}

	// build the API model
	item := Item{
		ID:              row.ID,
		FeedID:          row.FeedID,
		Title:           row.Title,
		Description:     row.Description,
		Content:         row.Content,
		Link:            row.Link,
		Links:           row.Links,
		UpdatedParsed:   row.UpdatedParsed,
		PublishedParsed: row.PublishedParsed,
		Authors:         Authors(authors),
		GUID:            row.Guid,
		Image:           row.Image,
		Categories:      row.Categories,
		Enclosures:      row.Enclosures,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
		Liked:           liked,
		LikedAt:         likedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// LikeItem marks an item as liked by the user.
// @Summary      Like item
// @Description  Creates a like record for the current user on an item.
// @Tags         Items
// @Param        itemID  path      string  true  "Item UUID"
// @Success      200     {object}  map[string]time.Time  "liked_at"
// @Failure      400     {object}	 string
// @Failure      409     {object}	 string
// @Failure      500     {object}	 string
// @Security     BearerAuth
// @Router       /api/items/{itemID}/like [post]
func (h *Handler) LikeItem(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	userID := claims.UserID

	idStr := r.PathValue("itemID")
	itemID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid item ID", http.StatusBadRequest)
		return
	}

	like, err := h.Repo.CreateUserLike(r.Context(), repository.CreateUserLikeParams{
		UserID: userID,
		ItemID: itemID,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to like item: %v", err), http.StatusConflict)
		return
	}

	resp := map[string]time.Time{"liked_at": like.LikedAt}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// UnlikeItem removes a like from an item.
// @Summary      Unlike item
// @Description  Deletes the like record for the current user on an item.
// @Tags         Items
// @Param        itemID  path      string  true  "Item UUID"
// @Success      204     "No Content"
// @Failure      400     {object}	 string
// @Failure      500     {object}	 string
// @Security     BearerAuth
// @Router       /api/items/{itemID}/like [delete]
func (h *Handler) UnlikeItem(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	userID := claims.UserID

	idStr := r.PathValue("itemID")
	itemID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid item ID", http.StatusBadRequest)
		return
	}

	if err := h.Repo.DeleteUserLike(r.Context(), repository.DeleteUserLikeParams{
		UserID: userID,
		ItemID: itemID,
	}); err != nil {
		http.Error(w, fmt.Sprintf("failed to unlike item: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
