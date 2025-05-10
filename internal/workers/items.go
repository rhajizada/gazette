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
	"github.com/rhajizada/gazette/internal/textsplitter"
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

	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("%s embeddings for item %s already exist", prefix, itemID)
		return nil
	}

	toEmbed := item.Description
	if len(*item.Content) > len(*item.Description) {
		toEmbed = item.Content
	}
	if len(*item.Title) > len(*toEmbed) {
		toEmbed = item.Title
	}

	splitter := &textsplitter.RecursiveCharacter{
		Separators:   []string{"\n\n", "\n", " ", ""},
		ChunkSize:    500,
		ChunkOverlap: 50,
	}
	extracted := ExtractTextFromHTML(*toEmbed)

	chunks := splitter.Split(extracted)

	for i, c := range chunks {
		req := api.EmbeddingRequest{
			Model:  h.OllamaConfig.EmbeddingsModel,
			Prompt: c,
		}

		resp, err := client.Embeddings(ctx, &req)
		if err != nil {
			return fmt.Errorf("failed to generate embedding for item %s: %v", itemID, err)
		}
		embeddingValue := float64sToVector(resp.Embedding)

		h.Repo.CreateItemEmbedding(ctx, repository.CreateItemEmbeddingParams{
			ItemID:     itemID,
			ChunkIndex: int32(i),
			Embedding:  &embeddingValue,
		})

	}

	log.Printf(
		"%s synced embddings for item %s ",
		prefix, itemID,
	)

	return nil
}
