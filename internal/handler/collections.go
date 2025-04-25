package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rhajizada/gazette/internal/middleware"
	"github.com/rhajizada/gazette/internal/repository"
)

// ListCollections returns paginated collections for the current user
func (h *Handler) ListCollections(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	userID := claims.UserID

	params, err := getPageParams(r.URL.Query())
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid paging: %v", err), http.StatusBadRequest)
		return
	}

	total, err := h.Repo.CountCollectionsByUserID(r.Context(), userID)
	if err != nil {
		msg := fmt.Sprintf("failed counting collections: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	rows, err := h.Repo.ListCollectionsByUser(r.Context(), repository.ListCollectionsByUserParams{
		UserID: userID,
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed listing collections: %v", err), http.StatusInternalServerError)
		return
	}

	cols := make([]Collection, len(rows))
	for i, c := range rows {
		cols[i] = Collection{
			ID:          c.ID,
			Name:        c.Name,
			CreatedAt:   c.CreatedAt,
			LastUpdated: c.LastUpdated,
		}
	}

	resp := ListCollectionsResponse{
		Limit:       params.Limit,
		Offset:      params.Offset,
		TotalCount:  total,
		Collections: cols,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// CreateCollection creates a new collection for the current user
func (h *Handler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	userID := claims.UserID

	var req CreateCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("bad JSON: %v", err), http.StatusBadRequest)
		return
	}

	col, err := h.Repo.CreateCollection(r.Context(), repository.CreateCollectionParams{
		UserID: userID,
		Name:   req.Name,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed creating collection: %v", err), http.StatusInternalServerError)
		return
	}

	resp := Collection{
		ID:          col.ID,
		Name:        col.Name,
		CreatedAt:   col.CreatedAt,
		LastUpdated: col.LastUpdated,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetCollectionByID returns one collection by ID
func (h *Handler) GetCollectionByID(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	userID := claims.UserID

	colID, err := uuid.Parse(r.PathValue("collectionID"))
	if err != nil {
		http.Error(w, "invalid collection ID", http.StatusBadRequest)
		return
	}

	col, err := h.Repo.GetCollectionByID(r.Context(), colID)
	if err != nil {
		http.Error(w, fmt.Sprintf("collection not found: %v", err), http.StatusNotFound)
		return
	}
	if col.UserID != userID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	resp := Collection{
		ID:          col.ID,
		Name:        col.Name,
		CreatedAt:   col.CreatedAt,
		LastUpdated: col.LastUpdated,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteCollectionByID deletes a collection by ID
func (h *Handler) DeleteCollectionByID(w http.ResponseWriter, r *http.Request) {
	colID, err := uuid.Parse(r.PathValue("collectionID"))
	if err != nil {
		http.Error(w, "invalid collection ID", http.StatusBadRequest)
		return
	}

	if err := h.Repo.DeleteCollectionByID(r.Context(), colID); err != nil {
		http.Error(w, fmt.Sprintf("failed deleting collection: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddItemToCollection adds an item to a collection
func (h *Handler) AddItemToCollection(w http.ResponseWriter, r *http.Request) {
	colID, err := uuid.Parse(r.PathValue("collectionID"))
	if err != nil {
		http.Error(w, "invalid collection ID", http.StatusBadRequest)
		return
	}
	itemID, err := uuid.Parse(r.PathValue("itemID"))
	if err != nil {
		http.Error(w, "invalid item ID", http.StatusBadRequest)
		return
	}

	rec, err := h.Repo.AddItemToCollection(r.Context(), repository.AddItemToCollectionParams{
		CollectionID: colID,
		ItemID:       itemID,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed adding item: %v", err), http.StatusInternalServerError)
		return
	}

	resp := struct {
		AddedAt time.Time `json:"added_at"`
	}{AddedAt: rec.AddedAt}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RemoveItemFromCollection removes an item from a collection
func (h *Handler) RemoveItemFromCollection(w http.ResponseWriter, r *http.Request) {
	colID, err := uuid.Parse(r.PathValue("collectionID"))
	if err != nil {
		http.Error(w, "invalid collection ID", http.StatusBadRequest)
		return
	}
	itemID, err := uuid.Parse(r.PathValue("itemID"))
	if err != nil {
		http.Error(w, "invalid item ID", http.StatusBadRequest)
		return
	}

	if err := h.Repo.RemoveItemFromCollection(r.Context(), repository.RemoveItemFromCollectionParams{
		CollectionID: colID,
		ItemID:       itemID,
	}); err != nil {
		http.Error(w, fmt.Sprintf("failed removing item: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListItemsByCollectionID returns paginated items in a collection
func (h *Handler) ListItemsByCollectionID(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	userID := claims.UserID

	colID, err := uuid.Parse(r.PathValue("collectionID"))
	if err != nil {
		http.Error(w, "invalid collection ID", http.StatusBadRequest)
		return
	}
	params, err := getPageParams(r.URL.Query())
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid paging: %v", err), http.StatusBadRequest)
		return
	}

	total, err := h.Repo.CountItemsInCollection(r.Context(), colID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed counting items: %v", err), http.StatusInternalServerError)
		return
	}

	rows, err := h.Repo.ListItemsInCollection(r.Context(), repository.ListItemsInCollectionParams{
		CollectionID: colID,
		Limit:        params.Limit,
		Offset:       params.Offset,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed listing items: %v", err), http.StatusInternalServerError)
		return
	}

	items := make([]Item, len(rows))
	for i, row := range rows {
		// convert sql.Null* to pointers

		// fetch like status per item
		liked := false
		var likedAtPtr *time.Time
		if like, err := h.Repo.GetUserLike(r.Context(), repository.GetUserLikeParams{UserID: userID, ItemID: row.ID}); err == nil {
			liked = true
			likedAtPtr = &like.LikedAt
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
			Authors:         row.Authors,
			GUID:            row.Guid,
			Image:           row.Image,
			Categories:      row.Categories,
			Enclosures:      row.Enclosures,
			CreatedAt:       row.CreatedAt,
			UpdatedAt:       row.UpdatedAt,
			Liked:           liked,
			LikedAt:         likedAtPtr,
		}
	}

	resp := ListCollectionItemsResponse{
		Limit:      params.Limit,
		Offset:     params.Offset,
		TotalCount: total,
		Items:      items,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
