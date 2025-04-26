package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/rhajizada/gazette/internal/middleware"
	"github.com/rhajizada/gazette/internal/repository"
	"github.com/rhajizada/gazette/internal/service"
)

// ListFeeds returns a paginated list of feeds.
// @Summary      List feeds
// @Description  Returns all feeds or only those the user is subscribed to.
// @Tags         Feeds
// @Param        subscribedOnly  query     bool   false  "Only subscribed feeds"
// @Param        limit           query     int32  true   "Max number of feeds to return"
// @Param        offset          query     int32  true   "Number of feeds to skip"
// @Success      200             {object}  service.ListFeedsResponse
// @Failure      400             {object}  string
// @Failure      500             {object}  string
// @Security     BearerAuth
// @Router       /api/feeds [get]
func (h *Handler) ListFeeds(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID

	params, err := getPageParams(r.URL.Query())
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid paging: %v", err), http.StatusBadRequest)
		return
	}

	subOnly := false
	if v := r.URL.Query().Get("subscribedOnly"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			subOnly = b
		}
	}

	resp, err := h.Service.ListFeeds(r.Context(), service.ListFeedsRequest{
		UserID:       userID,
		SubscbedOnly: subOnly,
		Limit:        params.Limit,
		Offset:       params.Offset,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed listing feeds: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// CreateFeed subscribes the user to a feed, creating it if necessary.
// @Summary      Create or subscribe feed
// @Description  Creates a new feed from URL or subscribes the user to it.
// @Tags         Feeds
// @Param        body  body      service.CreateFeedRequest  true  "Create feed payload"
// @Success      200   {object}  service.Feed
// @Failure      400   {object}  string
// @Failure      409   {object}  string
// @Failure      500   {object}  string
// @Security     BearerAuth
// @Router       /api/feeds [post]
func (h *Handler) CreateFeed(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID

	var req service.CreateFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("bad JSON: %v", err), http.StatusBadRequest)
		return
	}
	req.UserID = userID

	feed, err := h.Service.CreateFeed(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed creating feed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feed)
}

// GetFeedByID returns one feed with subscription info.
// @Summary      Get feed
// @Description  Retrieves a feed by ID, including the user’s subscription status.
// @Tags         Feeds
// @Param        feedID  path      string  true  "Feed UUID"
// @Success      200     {object}  service.Feed
// @Failure      400     {object}  string
// @Failure      404     {object}  string
// @Failure      500     {object}  string
// @Security     BearerAuth
// @Router       /api/feeds/{feedID} [get]
func (h *Handler) GetFeedByID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID
	feedID, err := uuid.Parse(r.PathValue("feedID"))
	if err != nil {
		http.Error(w, "invalid feed ID", http.StatusBadRequest)
		return
	}

	req := service.GetFeedRequest{
		GetUserFeedSubscriptionParams: repository.GetUserFeedSubscriptionParams{
			UserID: userID,
			FeedID: feedID,
		},
	}
	feed, err := h.Service.GetFeed(r.Context(), req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "feed not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("failed fetching feed: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feed)
}

// DeleteFeedByID deletes a feed entirely.
// @Summary      Delete feed
// @Description  Removes a feed and all its subscriptions/items.
// @Tags         Feeds
// @Param        feedID  path  string  true  "Feed UUID"
// @Success      204     "No Content"
// @Failure      400     {object}  string
// @Failure      500     {object}  string
// @Security     BearerAuth
// @Router       /api/feeds/{feedID} [delete]
func (h *Handler) DeleteFeedByID(w http.ResponseWriter, r *http.Request) {
	feedID, err := uuid.Parse(r.PathValue("feedID"))
	if err != nil {
		http.Error(w, "invalid feed ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteFeed(r.Context(), service.DeleteFeedRequest{FeedID: feedID}); err != nil {
		http.Error(w, fmt.Sprintf("failed deleting feed: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SubscribeToFeed subscribes the current user to a feed.
// @Summary      Subscribe to feed
// @Description  Subscribes the user to an existing feed.
// @Tags         Feeds
// @Param        feedID  path      string  true  "Feed UUID"
// @Success      200     {object}  service.SubscibeToFeedResponse
// @Failure      400     {object}  string
// @Failure      409     {object}  string
// @Failure      500     {object}  string
// @Security     BearerAuth
// @Router       /api/feeds/{feedID}/subscribe [put]
func (h *Handler) SubscribeToFeed(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID
	feedID, err := uuid.Parse(r.PathValue("feedID"))
	if err != nil {
		http.Error(w, "invalid feed ID", http.StatusBadRequest)
		return
	}

	respData, err := h.Service.SubscribeToFeed(r.Context(), service.SubscribeToFeedRequest{
		CreateUserFeedSubscriptionParams: repository.CreateUserFeedSubscriptionParams{
			UserID: userID,
			FeedID: feedID,
		},
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed subscribing: %v", err), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respData)
}

// UnsubscribeFromFeed unsubscribes the current user.
// @Summary      Unsubscribe from feed
// @Description  Removes the user’s subscription to a feed.
// @Tags         Feeds
// @Param        feedID  path  string  true  "Feed UUID"
// @Success      204     "No Content"
// @Failure      400     {object}  string
// @Failure      500     {object}  string
// @Security     BearerAuth
// @Router       /api/feeds/{feedID}/unsubscribe [put]
func (h *Handler) UnsubscribeFromFeed(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID
	feedID, err := uuid.Parse(r.PathValue("feedID"))
	if err != nil {
		http.Error(w, "invalid feed ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.UnsubscribeFromFeed(r.Context(), service.UnsubscribeFromFeedRequest{
		DeleteUserFeedSubscriptionParams: repository.DeleteUserFeedSubscriptionParams{
			UserID: userID,
			FeedID: feedID,
		},
	}); err != nil {
		http.Error(w, fmt.Sprintf("failed unsubscribing: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListItemsByFeedID returns paginated list of items.
// @Summary      List feed items
// @Description  Retrieves feed items.
// @Tags         Items
// @Param        feedID  path      string  true   "Feed UUID"
// @Param        limit   query     int32   true   "Max number of items"
// @Param        offset  query     int32   true   "Number of items to skip"
// @Success      200     {object}  service.ListItemsResponse
// @Failure      400     {object}  string
// @Failure      500     {object}  string
// @Security     BearerAuth
// @Router       /api/feeds/{feedID}/items [get]
func (h *Handler) ListItemsByFeedID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID
	feedID, err := uuid.Parse(r.PathValue("feedID"))
	if err != nil {
		http.Error(w, "invalid feed ID", http.StatusBadRequest)
		return
	}
	params, err := getPageParams(r.URL.Query())
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid paging: %v", err), http.StatusBadRequest)
		return
	}

	resp, err := h.Service.ListItemsByFeedID(r.Context(), service.ListItemsByFeedIDRequest{
		repository.ListItemsByFeedIDForUserParams{
			FeedID: feedID,
			UserID: userID,
			Limit:  params.Limit,
			Offset: params.Offset,
		},
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed listing items: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
