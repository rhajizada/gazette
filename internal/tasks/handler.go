package tasks

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/mmcdole/gofeed"
	"github.com/rhajizada/gazette/internal/repository"
	"github.com/rhajizada/gazette/internal/typeext"
)

type Handler struct {
	Repo repository.Queries
}

func NewHandler(r *repository.Queries) *Handler {
	return &Handler{
		Repo: *r,
	}
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
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed fetching last item for feed %q: %w", feedID, err)
	}

	var cutoff *time.Time
	if err == nil {
		if lastItem.UpdatedParsed != nil {
			cutoff = lastItem.UpdatedParsed
		} else if lastItem.PublishedParsed != nil {
			cutoff = lastItem.PublishedParsed
		}
	}

	var itemsToSync []*gofeed.Item
	if cutoff == nil {
		itemsToSync = feed.Items
	} else {
		for _, itm := range feed.Items {
			var ts *time.Time
			if itm.UpdatedParsed != nil {
				ts = itm.UpdatedParsed
			} else if itm.PublishedParsed != nil {
				ts = itm.PublishedParsed
			}

			if ts != nil && ts.After(*cutoff) {
				itemsToSync = append(itemsToSync, itm)
			}
		}
	}

	for _, itm := range itemsToSync {
		i, err := h.Repo.CreateItem(ctx, repository.CreateItemParams{
			FeedID:          feedID,
			Title:           &itm.Title,
			Description:     &itm.Description,
			Content:         &itm.Content,
			Link:            itm.Link,
			Links:           itm.Links,
			UpdatedParsed:   itm.UpdatedParsed,
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
		log.Printf("synced item %s from feed %s", i.ID, feedID)
	}

	return nil
}
