package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	"github.com/rhajizada/gazette/internal/middleware"
	"github.com/rhajizada/gazette/internal/repository"
	"github.com/rhajizada/gazette/internal/tasks"
	"github.com/rhajizada/gazette/internal/typeext"
)

// ListFeeds returns only the feeds subscribed by the current user
func (h *Handler) ListFeeds(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	userID := claims.UserID

	params, err := getPageParams(r.URL.Query())
	if err != nil {
		msg := fmt.Sprintf("invalid paging: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	total, err := h.Repo.CountFeedsByUserID(r.Context(), userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		msg := fmt.Sprintf("failed fethcing user feeds: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	resp := ListFeedsResponse{
		Limit:      params.Limit,
		Offset:     params.Offset,
		TotalCount: total,
	}
	if total == 0 {
		resp.Feeds = []repository.ListFeedsByUserIDRow{}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	feeds, err := h.Repo.ListFeedsByUserID(r.Context(), repository.ListFeedsByUserIDParams{
		UserID: userID,
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		msg := fmt.Sprintf("list failed: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	resp.Feeds = feeds

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// CreateFeed creates or binds a feed for the current user; errors if already subscribed
func (h *Handler) CreateFeed(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	userID := claims.UserID

	var req CreateFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := fmt.Sprintf("bad JSON: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	fp := gofeed.NewParser()
	remote, err := fp.ParseURL(req.FeedURL)
	if err != nil {
		msg := fmt.Sprintf("invalid feed URL: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// Lookup or create the feed record
	feed, err := h.Repo.GetFeedByFeedLink(r.Context(), remote.FeedLink)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		msg := fmt.Sprintf("lookup failed: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	if errors.Is(err, sql.ErrNoRows) {
		input := repository.CreateFeedParams{
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
		}
		feed, err = h.Repo.CreateFeed(r.Context(), input)
		if err != nil {
			msg := fmt.Sprintf("create failed: %v", err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		task, _ := tasks.NewFeedSyncTask(feed.ID)
		h.Client.Enqueue(task)
	}

	// Check existing user subscription
	_, err = h.Repo.GetUserFeedSubscription(r.Context(), repository.GetUserFeedSubscriptionParams{
		UserID: userID,
		FeedID: feed.ID,
	})
	if err == nil {
		msg := "already subscribed to this feed"
		http.Error(w, msg, http.StatusBadRequest)
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		msg := fmt.Sprintf("subscription lookup failed: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Bind feed to user
	_, err = h.Repo.CreateUserFeedSubscription(r.Context(), repository.CreateUserFeedSubscriptionParams{
		UserID: userID,
		FeedID: feed.ID,
	})
	if err != nil {
		msg := fmt.Sprintf("subscription failed: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feed)
}

// DeleteFeedByID unsubscribes the feed for the current user
func (h *Handler) DeleteFeedByID(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserClaims(r)
	userID := claims.UserID

	feedIDStr := r.PathValue("feedID")
	feedID, err := uuid.Parse(feedIDStr)
	if err != nil {
		msg := "invalid feed ID"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	err = h.Repo.DeleteUserFeedSubscription(r.Context(), repository.DeleteUserFeedSubscriptionParams{
		UserID: userID,
		FeedID: feedID,
	})
	if err != nil {
		msg := fmt.Sprintf("unsubscribe failed: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetFeedByID(w http.ResponseWriter, r *http.Request) {
	feedID := r.PathValue("feedID")
	if feedID == "" {
		http.Error(w, "missing 'id' parameter", http.StatusBadRequest)
		return
	}
	feedUUID, err := uuid.Parse(feedID)
	if err != nil {
		http.Error(w, "cannot parse 'id' parameter", http.StatusInternalServerError)
		return
	}
	data, err := h.Repo.GetFeedByID(r.Context(), feedUUID)
	if err != nil {
		msg := fmt.Sprintf("failed fetching feed: %v", err)
		http.Error(w, msg, http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ListItemsByFeedID(w http.ResponseWriter, r *http.Request) {
	feedID := r.PathValue("feedID")
	if feedID == "" {
		http.Error(w, "missing 'id' parameter", http.StatusBadRequest)
		return
	}
	feedUUID, err := uuid.Parse(feedID)
	if err != nil {
		http.Error(w, "cannot parse 'id' parameter", http.StatusInternalServerError)
		return
	}
	v := r.URL.Query()
	params, err := getPageParams(v)
	if err != nil {
		msg := fmt.Sprintf("failed listing feeds: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	var response ListItemsResponse
	response.Limit = params.Limit
	response.Offset = params.Offset
	count, err := h.Repo.CountItemsByFeedID(r.Context(), feedUUID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		msg := fmt.Sprintf("failed listing items: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	response.TotalCount = count
	if count == 0 {
		response.Items = make([]repository.Item, 0)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	data, err := h.Repo.ListItemsByFeedID(r.Context(),
		repository.ListItemsByFeedIDParams{
			FeedID: feedUUID,
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	if err != nil {
		msg := fmt.Sprintf("failed listing items: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	response.Items = data

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
