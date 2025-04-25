package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	"github.com/rhajizada/gazette/internal/middleware"
	"github.com/rhajizada/gazette/internal/repository"
	"github.com/rhajizada/gazette/internal/tasks"
	"github.com/rhajizada/gazette/internal/typeext"
)

// ListFeeds returns feeds; if subscribedOnly=true, only subscribed ones
func (h *Handler) ListFeeds(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	userID := claims.UserID

	params, err := getPageParams(r.URL.Query())
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid paging: %v", err), http.StatusBadRequest)
		return
	}

	// optional subscribedOnly flag
	subOnly := false
	if v := r.URL.Query().Get("subscribedOnly"); v != "" {
		subOnly, _ = strconv.ParseBool(v)
	}

	var total int64
	if subOnly {
		total, err = h.Repo.CountFeedsByUserID(r.Context(), userID)
	} else {
		total, err = h.Repo.CountFeeds(r.Context())
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("failed counting feeds: %v", err), http.StatusInternalServerError)
		return
	}

	rows, err := h.Repo.ListFeedsByUserID(r.Context(), repository.ListFeedsByUserIDParams{
		UserID:  userID,
		Column2: subOnly,
		Limit:   params.Limit,
		Offset:  params.Offset,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed listing feeds: %v", err), http.StatusInternalServerError)
		return
	}

	feeds := make([]Feed, len(rows))
	for i, row := range rows {
		feeds[i] = Feed{
			ID:              row.ID,
			Title:           row.Title,
			Description:     row.Description,
			Link:            row.Link,
			FeedLink:        row.FeedLink,
			Links:           row.Links,
			UpdatedParsed:   row.UpdatedParsed,
			PublishedParsed: row.PublishedParsed,
			Authors:         row.Authors,
			Language:        row.Language,
			Image:           row.Image,
			Copyright:       row.Copyright,
			Generator:       row.Generator,
			Categories:      row.Categories,
			FeedType:        row.FeedType,
			FeedVersion:     row.FeedVersion,
			CreatedAt:       row.CreatedAt,
			LastUpdatedAt:   row.LastUpdatedAt,
			Subscribed:      row.SubscribedAt != nil,
			SubscribedAt:    row.SubscribedAt,
		}
	}

	resp := ListFeedsResponse{
		Limit:      params.Limit,
		Offset:     params.Offset,
		TotalCount: total,
		Feeds:      feeds,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// CreateFeed subscribes the current user to a feed (creating it if needed)
func (h *Handler) CreateFeed(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	userID := claims.UserID

	var req CreateFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("bad JSON: %v", err), http.StatusBadRequest)
		return
	}

	fp := gofeed.NewParser()
	remote, err := fp.ParseURL(req.FeedURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid feed URL: %v", err), http.StatusBadRequest)
		return
	}

	// lookup or create feed
	feed, err := h.Repo.GetFeedByFeedLink(r.Context(), remote.FeedLink)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		feed, err = h.Repo.CreateFeed(r.Context(), repository.CreateFeedParams{
			Title:           &remote.Title,
			Description:     &remote.Description,
			Link:            &remote.Link,
			FeedLink:        remote.FeedLink,
			Links:           remote.Links,
			UpdatedParsed:   remote.UpdatedParsed,
			PublishedParsed: remote.PublishedParsed,
			Authors:         typeext.Authors(remote.Authors),
			Language:        &remote.Language,
			Image:           remote.Image,
			Copyright:       &remote.Copyright,
			Generator:       &remote.Generator,
			Categories:      remote.Categories,
			FeedType:        &remote.FeedType,
			FeedVersion:     &remote.FeedVersion,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed creating feed: %v", err), http.StatusInternalServerError)
			return
		}
		task, _ := tasks.NewFeedSyncTask(feed.ID)
		h.Client.Enqueue(task)
	} else if err != nil {
		http.Error(w, fmt.Sprintf("lookup failed: %v", err), http.StatusInternalServerError)
		return
	}

	// subscribe user
	sub, err := h.Repo.CreateUserFeedSubscription(r.Context(), repository.CreateUserFeedSubscriptionParams{
		UserID: userID,
		FeedID: feed.ID,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed subscribing: %v", err), http.StatusConflict)
		return
	}

	resp := Feed{
		ID:              feed.ID,
		Title:           feed.Title,
		Description:     feed.Description,
		Link:            feed.Link,
		FeedLink:        feed.FeedLink,
		Links:           feed.Links,
		UpdatedParsed:   feed.UpdatedParsed,
		PublishedParsed: feed.PublishedParsed,
		Authors:         feed.Authors,
		Language:        feed.Language,
		Image:           feed.Image,
		Copyright:       feed.Copyright,
		Generator:       feed.Generator,
		Categories:      feed.Categories,
		FeedType:        feed.FeedType,
		FeedVersion:     feed.FeedVersion,
		CreatedAt:       feed.CreatedAt,
		LastUpdatedAt:   feed.LastUpdatedAt,
		Subscribed:      true,
		SubscribedAt:    &sub.SubscribedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetFeedByID returns one feed with subscription info
func (h *Handler) GetFeedByID(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	userID := claims.UserID

	idStr := r.PathValue("feedID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid feed ID", http.StatusBadRequest)
		return
	}

	feed, err := h.Repo.GetFeedByID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("feed not found: %v", err), http.StatusNotFound)
		return
	}

	subAt := (*time.Time)(nil)
	if uf, err := h.Repo.GetUserFeedSubscription(r.Context(), repository.GetUserFeedSubscriptionParams{UserID: userID, FeedID: id}); err == nil {
		subAt = &uf.SubscribedAt
	}

	resp := Feed{
		ID:              feed.ID,
		Title:           feed.Title,
		Description:     feed.Description,
		Link:            feed.Link,
		FeedLink:        feed.FeedLink,
		Links:           feed.Links,
		UpdatedParsed:   feed.UpdatedParsed,
		PublishedParsed: feed.PublishedParsed,
		Authors:         feed.Authors,
		Language:        feed.Language,
		Image:           feed.Image,
		Categories:      feed.Categories,
		FeedType:        feed.FeedType,
		FeedVersion:     feed.FeedVersion,
		CreatedAt:       feed.CreatedAt,
		LastUpdatedAt:   feed.LastUpdatedAt,
		Subscribed:      subAt != nil,
		SubscribedAt:    subAt,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteFeedByID removes a feed entirely
func (h *Handler) DeleteFeedByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("feedID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid feed ID", http.StatusBadRequest)
		return
	}

	err = h.Repo.DeleteFeedByID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed deleting feed: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SubscribeToFeed subscribes the current user to an existing feed
func (h *Handler) SubscribeToFeed(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	userID := claims.UserID

	idStr := r.PathValue("feedID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid feed ID", http.StatusBadRequest)
		return
	}

	sub, err := h.Repo.CreateUserFeedSubscription(r.Context(), repository.CreateUserFeedSubscriptionParams{UserID: userID, FeedID: id})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed subscribing: %v", err), http.StatusConflict)
		return
	}

	resp := map[string]any{"subscribed_at": sub.SubscribedAt}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// UnsubscribeFromFeed removes a user's subscription to a feed
func (h *Handler) UnsubscribeFromFeed(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	userID := claims.UserID

	idStr := r.PathValue("feedID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		msg := "invalid feed ID"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	err = h.Repo.DeleteUserFeedSubscription(r.Context(), repository.DeleteUserFeedSubscriptionParams{UserID: userID, FeedID: id})
	if err != nil {
		msg := fmt.Sprintf("failed unsubscribing: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListItemsByFeedID returns paginated items including like info
func (h *Handler) ListItemsByFeedID(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	userID := claims.UserID

	idStr := r.PathValue("feedID")
	feedID, err := uuid.Parse(idStr)
	if err != nil {
		msg := "invalid feed ID"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	params, err := getPageParams(r.URL.Query())
	if err != nil {
		msg := fmt.Sprintf("invalid paging: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// total count
	total, err := h.Repo.CountItemsByFeedID(r.Context(), feedID)
	if err != nil {
		msg := fmt.Sprintf("failed counting items: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// fetch items with like info
	rows, err := h.Repo.ListItemsByFeedIDForUser(r.Context(), repository.ListItemsByFeedIDForUserParams{
		FeedID: feedID,
		UserID: userID,
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		msg := fmt.Sprintf("failed listing items: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	items := make([]Item, len(rows))
	for i, row := range rows {
		likedAtPtr := row.LikedAt
		liked := row.LikedAt != nil

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

	resp := ListItemsResponse{
		Limit:      params.Limit,
		Offset:     params.Offset,
		TotalCount: total,
		Items:      items,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
