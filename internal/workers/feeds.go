package workers

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

func (h *Handler) HandleFeedSync(ctx context.Context, t *asynq.Task) error {
	prefix := fmt.Sprintf(
		"task %s -",
		t.ResultWriter().TaskID(),
	)
	var p SyncFeedPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v", asynq.SkipRetry)
	}
	feedID := p.FeedID

	data, err := h.Repo.GetFeedByID(ctx, feedID)
	if err != nil {
		return fmt.Errorf("failed to get feed %q: %v", feedID, err)
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(data.FeedLink)
	if err != nil {
		return fmt.Errorf("failed to parse feed %q: %v", feedID, err)
	}

	lastItem, err := h.Repo.GetLastItem(ctx, feedID)
	var itemsToSync []*gofeed.Item
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			itemsToSync = feed.Items
		} else {
			return fmt.Errorf("failed to fetch last item from feed %q: %v", feedID, err)
		}
	} else {
		cutoff := *lastItem.PublishedParsed
		for _, itm := range feed.Items {
			if itm.PublishedParsed.After(cutoff) {
				itemsToSync = append(itemsToSync, itm)
			}
		}
	}

	log.Printf("%s %d items to sync for feed %q", prefix, len(itemsToSync), feedID)

	for _, itm := range itemsToSync {
		content := &itm.Content
		description := &itm.Description
		if len(itm.Description) > len(itm.Content) {
			content = &itm.Description
			description = &itm.Content
		}

		r, err := h.Repo.CreateItem(ctx, repository.CreateItemParams{
			FeedID:          feedID,
			Title:           &itm.Title,
			Description:     description,
			Content:         content,
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
			return fmt.Errorf("failed to create item %q for feed %q: %v", itm.GUID, feedID, err)
		}
		log.Printf("%s synced item %s from feed %s", prefix, itm.GUID, feedID)
		task, _ := NewEmbedItemTask(r.ID)
		tResp, err := h.Client.Enqueue(task, asynq.Queue("default"))
		if err != nil {
			return fmt.Errorf("failed to queue embedding task for item %s", r.ID)
		}
		log.Printf(
			"%s queued embdding task %s for item %s ",
			prefix, r.ID, tResp.ID,
		)
	}

	return nil
}
