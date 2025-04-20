package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	"github.com/rhajizada/gazette/internal/repository"
	"github.com/rhajizada/gazette/internal/typeext"
)

// Handler encapsulates dependencies for HTTP handlers.
type Handler struct {
	Repo   repository.Queries
	Secret []byte
}

// New creates a new Handler.
func New(r *repository.Queries, secret []byte) *Handler {
	return &Handler{
		Repo:   *r,
		Secret: secret,
	}
}

func (h *Handler) ListFeeds(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	params, err := getListFeedsParams(v)
	if err != nil {
		msg := fmt.Sprintf("failed listing feeds: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	data, err := h.Repo.ListFeeds(r.Context(), params)
	if err != nil {
		msg := fmt.Sprintf("failed listing feeds: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateFeed(w http.ResponseWriter, r *http.Request) {
	var body CreateFeedRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		msg := fmt.Sprintf("error decoding JSON: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(body.FeedURL)
	if err != nil {
		msg := fmt.Sprintf("invalid feed URL: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	input := repository.CreateFeedParams{
		Title:           &feed.Title,
		Description:     &feed.Description,
		Link:            &feed.Link,
		FeedLink:        feed.FeedLink,
		Links:           feed.Links,
		UpdatedParsed:   feed.UpdatedParsed,
		PublishedParsed: feed.PublishedParsed,
		Authors:         typeext.Authors(feed.Authors),
		Language:        &feed.Language,
		Image:           feed.Image,
		Copyright:       &feed.Copyright,
		Generator:       &feed.Generator,
		Categories:      feed.Categories,
		FeedType:        &feed.FeedType,
		FeedVersion:     &feed.FeedVersion,
	}

	data, err := h.Repo.CreateFeed(r.Context(), input)
	if err != nil {
		msg := fmt.Sprintf("failed creating feed: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteFeedByID(w http.ResponseWriter, r *http.Request) {
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

	_, err = h.Repo.GetFeedByID(r.Context(), feedUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.Repo.DeleteFeedByID(r.Context(), feedUUID)
	msg := fmt.Sprintf("successfully deleted feed %v", feedID)
	w.Write([]byte(msg))
}
