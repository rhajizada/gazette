package workers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/ollama/ollama/api"
	"github.com/rhajizada/gazette/internal/repository"
)

func (h *Handler) HandleEmbedItem(ctx context.Context, t *asynq.Task) error {
	prefix := fmt.Sprintf(
		"task %s -",
		t.ResultWriter().TaskID(),
	)
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
	if len(*item.Title) > len(*toEmbed) {
		toEmbed = item.Title
	}

	extracted := ExtractTextFromHTML(*toEmbed)

	req := api.EmbeddingRequest{
		Model:  h.OllamaConfig.EmbeddingsModel,
		Prompt: extracted,
	}

	resp, err := client.Embeddings(ctx, &req)
	if err != nil {
		return fmt.Errorf("failed to generate embedding for item %s: %v", itemID, err)
	}
	log.Printf(
		"%s generated embddings for item %s ",
		prefix, itemID,
	)

	embeddingValue := vectorFromFloat64s(resp.Embedding)

	exists := true
	_, err = h.Repo.GetItemEmbeddingByID(ctx, itemID)
	if errors.Is(err, sql.ErrNoRows) {
		exists = false
	} else {
		return fmt.Errorf("failed to generate embeddings for item %q: %v", itemID, err)
	}

	if exists {
		_, err = h.Repo.UpdateItemEmbeddingByID(ctx, repository.UpdateItemEmbeddingByIDParams{
			ItemID:    itemID,
			Embedding: &embeddingValue,
		})
		if err != nil {
			return fmt.Errorf("failed to sync embeddings for item %s: %v", itemID, err)
		}
	} else {
		_, err = h.Repo.CreateItemEmbedding(ctx, repository.CreateItemEmbeddingParams{
			ItemID:    itemID,
			Embedding: &embeddingValue,
		})
		if err != nil {
			return fmt.Errorf("failed to sync embeddings for item %s: %v", itemID, err)
		}
	}
	log.Printf(
		"%s synced embddings for item %s ",
		prefix, itemID,
	)

	return nil
}
