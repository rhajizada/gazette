package tasks

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/mmcdole/gofeed"
	"github.com/rhajizada/gazette/internal/repository"
	"github.com/rhajizada/gazette/internal/typeext"
)

const MaxLimit = 100

type Handler struct {
	Repo   repository.Queries
	Client *asynq.Client
}

func NewHandler(r *repository.Queries, c *asynq.Client) *Handler {
	return &Handler{
		Repo:   *r,
		Client: c,
	}
}

func (h *Handler) HandleDataSync(ctx context.Context, t *asynq.Task) error {
	count, err := h.Repo.CountFeeds(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed listing feeds: %v", err)
	}
	if count == 0 {
		return nil
	}
	for offset := int64(0); offset < count; offset += MaxLimit {
		// calculate this batch’s size
		limit := MaxLimit
		if remaining := count - offset; remaining < MaxLimit {
			limit = int(remaining)
		}

		// ListFeedsByID should be your SQLC-generated method:
		// func (q *Queries) ListFeeds(ctx context.Context, limit, offset int64) ([]Feed, error)
		feeds, err := h.Repo.ListFeeds(ctx, repository.ListFeedsParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			return fmt.Errorf("failed listing feeds (offset %d): %v", offset, err)
		}

		for _, feed := range feeds {
			task, err := NewFeedSyncTask(feed.ID)
			if err != nil {
				return err
			}

			ti, err := h.Client.Enqueue(task)
			if err != nil {
				return err
			}
			log.Printf("queued sync task %s for feed %s", ti.ID, feed.ID)

		}

	}
	return nil
}

func (h *Handler) HandleFeedSync(ctx context.Context, t *asynq.Task) error {
	var p FeedSyncPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %w", asynq.SkipRetry)
	}
	feedID := p.FeedID

	data, err := h.Repo.GetFeedByID(ctx, feedID)
	if err != nil {
		return fmt.Errorf("failed fetching feed %q: %w", feedID, err)
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(data.FeedLink)
	if err != nil {
		return fmt.Errorf("failed parsing feed %q: %w", feedID, err)
	}

	lastItem, err := h.Repo.GetLastItem(ctx, feedID)
	var itemsToSync []*gofeed.Item
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			itemsToSync = feed.Items
		} else {
			return fmt.Errorf("failed fetching last item for feed %q: %w", feedID, err)
		}
	} else {
		cutoff := *lastItem.PublishedParsed
		for _, itm := range feed.Items {
			if itm.PublishedParsed.After(cutoff) {
				itemsToSync = append(itemsToSync, itm)
			}
		}
	}

	log.Printf("task %s - %d items to sync", t.ResultWriter().TaskID(), len(itemsToSync))

	for _, itm := range itemsToSync {
		_, err := h.Repo.CreateItem(ctx, repository.CreateItemParams{
			FeedID:          feedID,
			Title:           &itm.Title,
			Description:     &itm.Description,
			Content:         &itm.Content,
			Link:            itm.Link,
			Links:           itm.Links,
			PublishedParsed: itm.PublishedParsed,
			Authors:         typeext.Authors(itm.Authors),
			Guid:            &itm.GUID,
			Image:           data.Image,
			Categories:      itm.Categories,
			Enclosures:      typeext.Enclosures(itm.Enclosures),
		})
		if err != nil {
			return fmt.Errorf("failed creating item %q for feed %q: %w", itm.GUID, feedID, err)
		}
		log.Printf("synced item %s from feed %s", itm.GUID, feedID)
	}

	return nil
}
