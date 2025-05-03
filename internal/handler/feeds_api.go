package handler

import (
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subOnly := false
	if v := r.URL.Query().Get("subscribedOnly"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			subOnly = b
		}
	}

	var resp *service.ListFeedsResponse
	resp, err = h.Service.ListFeeds(r.Context(), service.ListFeedsRequest{
		UserID:       userID,
		SubscbedOnly: subOnly,
		Limit:        params.Limit,
		Offset:       params.Offset,
	})
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			resp = &service.ListFeedsResponse{
				Offset:     params.Offset,
				Limit:      params.Limit,
				TotalCount: 0,
				Feeds:      make([]service.Feed, 0),
			}
		} else {
			http.Error(w, "failed listing feeds", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// CreateFeed subscribes the user to a feed, creating it if necessary.
// @Summary      Create or subscribe feed
// @Description  Creates a new feed from URL or subscribes the user to it.
// @Tags         Feeds
// @Param        body  body      CreateFeedRequest  true  "Create feed payload"
// @Success      200   {object}  service.Feed
// @Failure      400   {object}  string
// @Failure      409   {object}  string
// @Failure      500   {object}  string
// @Security     BearerAuth
// @Router       /api/feeds [post]
func (h *Handler) CreateFeed(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID

	var req CreateFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := service.CreateFeedRequest{
		FeedURL: req.FeedURL,
		UserID:  userID,
	}

	feed, err := h.Service.CreateFeed(r.Context(), query)
	if err != nil {
		if errors.Is(err, service.ErrAlreadyExists) {
			http.Error(w, fmt.Sprintf("feed %s already exists", req.FeedURL), http.StatusConflict)
			return
		} else {
			http.Error(w, fmt.Sprintf("failed importing feed %s", req.FeedURL), http.StatusInternalServerError)
			return
		}
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
	path := r.PathValue("feedID")
	feedID, err := uuid.Parse(path)
	if err != nil {
		http.Error(w, fmt.Sprintf("bad input: %s", path), http.StatusBadRequest)
		return
	}

	req := repository.GetUserFeedSubscriptionParams{
		UserID: userID,
		FeedID: feedID,
	}
	feed, err := h.Service.GetFeed(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, fmt.Sprintf("feed %s not found", feedID), http.StatusNotFound)
			return
		} else {
			http.Error(w, fmt.Sprintf("failed fetching feed %s", feedID), http.StatusInternalServerError)
			return
		}
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
	path := r.PathValue("feedID")
	feedID, err := uuid.Parse(path)
	if err != nil {
		http.Error(w, fmt.Sprintf("bad input: %s", path), http.StatusBadRequest)
		return
	}

	err = h.Service.DeleteFeed(r.Context(), service.DeleteFeedRequest{FeedID: feedID})
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, fmt.Sprintf("feed %s not found", feedID), http.StatusNotFound)
			return
		} else {
			http.Error(w, fmt.Sprintf("failed deleting feed %s", feedID), http.StatusInternalServerError)
			return
		}
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
// @Failure      500     {object}  string
// @Security     BearerAuth
// @Router       /api/feeds/{feedID}/subscribe [put]
func (h *Handler) SubscribeToFeed(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID
	path := r.PathValue("feedID")
	feedID, err := uuid.Parse(path)
	if err != nil {
		http.Error(w, fmt.Sprintf("bad input: %s", path), http.StatusBadRequest)
		return
	}

	respData, err := h.Service.SubscribeToFeed(r.Context(),
		repository.CreateUserFeedSubscriptionParams{
			UserID: userID,
			FeedID: feedID,
		})
	if err != nil {
		if errors.Is(err, service.ErrBadInput) {
			http.Error(w, fmt.Sprintf("failed subscribing to feed %s", feedID), http.StatusBadRequest)
			return
		} else {
			http.Error(w, fmt.Sprintf("failed subscribing to feed %s", feedID), http.StatusInternalServerError)
			return
		}
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
// @Router       /api/feeds/{feedID}/subscribe [delete]
func (h *Handler) UnsubscribeFromFeed(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserClaims(r).UserID
	path := r.PathValue("feedID")
	feedID, err := uuid.Parse(path)
	if err != nil {
		http.Error(w, fmt.Sprintf("bad input: %s", path), http.StatusBadRequest)
		return
	}

	err = h.Service.UnsubscribeFromFeed(r.Context(),
		repository.DeleteUserFeedSubscriptionParams{
			UserID: userID,
			FeedID: feedID,
		})
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, fmt.Sprintf("user is not subscribed to feed %s", feedID), http.StatusNotFound)
			return
		} else {
			http.Error(w, fmt.Sprintf("failed unsubscribing from feed %s", feedID), http.StatusInternalServerError)
			return
		}
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
	path := r.PathValue("feedID")
	feedID, err := uuid.Parse(path)
	if err != nil {
		http.Error(w, fmt.Sprintf("bad input: %s", path), http.StatusBadRequest)
		return
	}
	params, err := getPageParams(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var resp *service.ListItemsResponse

	resp, err = h.Service.ListItemsByFeedID(r.Context(),
		repository.ListItemsByFeedIDForUserParams{
			FeedID: feedID,
			UserID: userID,
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			resp = &service.ListItemsResponse{
				Offset:     params.Offset,
				Limit:      params.Limit,
				TotalCount: 0,
				Items:      make([]service.Item, 0),
			}
		} else {
			http.Error(w, fmt.Sprintf("failed listing items in feed %s", feedID), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
