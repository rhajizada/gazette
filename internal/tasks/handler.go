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
	"github.com/ollama/ollama/api"
	"github.com/rhajizada/gazette/internal/config"
	"github.com/rhajizada/gazette/internal/repository"
	"github.com/rhajizada/gazette/internal/typeext"
)

const MaxLimit = 100

type Handler struct {
	Repo         repository.Queries
	Client       *asynq.Client
	OllamaConfig *config.OllamaConfig
}

func NewHandler(r *repository.Queries, c *asynq.Client, o *config.OllamaConfig) *Handler {
	return &Handler{
		Repo:         *r,
		Client:       c,
		OllamaConfig: o,
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
			task, err := NewSyncFeedTask(feed.ID)
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
	var p SyncFeedPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v", asynq.SkipRetry)
	}
	feedID := p.FeedID

	data, err := h.Repo.GetFeedByID(ctx, feedID)
	if err != nil {
		return fmt.Errorf("failed fetching feed %q: %v", feedID, err)
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(data.FeedLink)
	if err != nil {
		return fmt.Errorf("failed parsing feed %q: %v", feedID, err)
	}

	lastItem, err := h.Repo.GetLastItem(ctx, feedID)
	var itemsToSync []*gofeed.Item
	// TODO: revisit this
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			itemsToSync = feed.Items
		} else {
			return fmt.Errorf("failed fetching last item for feed %q: %v", feedID, err)
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
		r, err := h.Repo.CreateItem(ctx, repository.CreateItemParams{
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
			return fmt.Errorf("failed creating item %q for feed %q: %v", itm.GUID, feedID, err)
		}
		log.Printf("synced item %s from feed %s", itm.GUID, feedID)
		task, _ := NewEmbedItemTask(r.ID)
		tResp, err := h.Client.Enqueue(task)
		if err != nil {
			return fmt.Errorf("failed to queue embedding task for item %s", r.ID)
		}
		log.Printf(
			"task %s - queued embdding task %s for item %s ",
			t.ResultWriter().TaskID(), r.ID, tResp.ID,
		)
	}

	return nil
}

func (h *Handler) HandleEmbedItem(ctx context.Context, t *asynq.Task) error {
	var p EmbedItemPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v", asynq.SkipRetry)
	}
	itemID := p.ItemID

	client, err := GetOllamaClient(h.OllamaConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize ollama client: %v", err)
	}

	item, err := h.Repo.GetItemByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("failed to fetch item %s: %v", itemID, err)
	}

	toEmbed := item.Description
	if len(*item.Content) > len(*item.Description) {
		toEmbed = item.Content
	}

	extracted := ExtractTextFromHTML(*toEmbed)

	req := api.EmbeddingRequest{
		Model:  h.OllamaConfig.EmbeddingsModel,
		Prompt: extracted,
	}

	resp, err := client.Embeddings(ctx, &req)
	if err != nil {
		return fmt.Errorf("failed generating embedding for item %s: %v", itemID, err)
	}
	log.Printf(
		"task %s - generated embddings for item %s ",
		t.ResultWriter().TaskID(), itemID,
	)

	embeddingValue := vectorFromFloat64s(resp.Embedding)

	exists := true
	_, err = h.Repo.GetItemEmbeddingByID(ctx, itemID)
	if errors.Is(err, sql.ErrNoRows) {
		exists = false
	} else {
		return fmt.Errorf("failed embeddings for item %q: %v", itemID, err)
	}

	if exists {
		_, err = h.Repo.UpdateItemEmbeddingByID(ctx, repository.UpdateItemEmbeddingByIDParams{
			ItemID:    itemID,
			Embedding: &embeddingValue,
		})
		if err != nil {
			return fmt.Errorf("failed synching embeddings for item %s: %v", itemID, err)
		}
	} else {
		_, err = h.Repo.CreateItemEmbedding(ctx, repository.CreateItemEmbeddingParams{
			ItemID:    itemID,
			Embedding: &embeddingValue,
		})
		if err != nil {
			return fmt.Errorf("failed synching embeddings for item %s: %v", itemID, err)
		}
	}
	log.Printf(
		"task %s - synced embddings for item %s ",
		t.ResultWriter().TaskID(), itemID,
	)

	return nil
}
