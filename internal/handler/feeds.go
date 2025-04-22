package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	"github.com/rhajizada/gazette/internal/repository"
	"github.com/rhajizada/gazette/internal/tasks"
	"github.com/rhajizada/gazette/internal/typeext"
)

func (h *Handler) ListFeeds(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	params, err := getPageParams(v)
	if err != nil {
		msg := fmt.Sprintf("failed listing feeds: %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	var response ListFeedsResponse
	response.Limit = params.Limit
	response.Offset = response.Offset
	count, err := h.Repo.CountFeeds(r.Context())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		msg := fmt.Sprintf("failed listing feeds: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	response.TotalCount = count
	if count == 0 {
		response.Feeds = make([]repository.Feed, 0)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	data, err := h.Repo.ListFeeds(r.Context(),
		repository.ListFeedsParams{
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	if err != nil {
		msg := fmt.Sprintf("failed listing feeds: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	response.Feeds = data

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
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
	log.Printf("PublishedParsed: %s", feed.PublishedParsed)
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
	task, err := tasks.NewFeedSyncTask(data.ID)
	if err != nil {
		h.Repo.DeleteFeedByID(r.Context(), data.ID)
		msg := fmt.Sprintf("failed marshalling sync task: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	ti, err := h.Client.Enqueue(task)
	if err != nil {
		h.Repo.DeleteFeedByID(r.Context(), data.ID)
		msg := fmt.Sprintf("failed queuing sync task: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	log.Printf("queued sync task %s for feed %s", ti.ID, data.ID)

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
	response.Offset = response.Offset
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
		msg := fmt.Sprintf("failed listing feeds: %v", err)
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
