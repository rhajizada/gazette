package workers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"

	"github.com/rhajizada/gazette/internal/cache"
	"github.com/rhajizada/gazette/internal/repository"
)

const MaxItemsLimit int32 = 100

func (h *Handler) HandleCacheUser(ctx context.Context, t *asynq.Task) error {
	prefix := fmt.Sprintf("task %s -", t.ResultWriter().TaskID())

	var p CacheUserPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("invalid payload: %w", asynq.SkipRetry)
	}
	userID := p.UserID
	resp, err := h.Repo.ListSuggestedItemsForUserCache(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("%s - no suggested items found for user %s", prefix, userID)
			return nil
		} else {
			return fmt.Errorf(
				"failed to fetch suggested items for user %s", userID,
			)
		}
	}

	items := make([]cache.SuggestedItem, len(resp))
	for i, j := range resp {
		items[i] = cache.SuggestedItem{
			ID:         j.ID,
			Freshness:  j.Freshness,
			Similarity: j.Similarity,
			Score:      j.Score,
		}
	}

	err = h.Cache.StoreUserSuggestions(ctx, userID, items, time.Hour)
	if err != nil {
		return err
	}
	log.Printf("%s cached suggested items for user %s", prefix, userID)
	return nil
}

func (h *Handler) HandleEmbedUser(ctx context.Context, t *asynq.Task) error {
	prefix := fmt.Sprintf("task %s -", t.ResultWriter().TaskID())

	var p EmbedUserPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("invalid payload: %w", asynq.SkipRetry)
	}
	userID := p.UserID

	total, err := h.Repo.CountLikedItems(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			total = 0
		} else {
			return fmt.Errorf("counting liked items: %w", err)
		}
	}
	if total == 0 {
		log.Printf("%s user %s has no liked items; skipping clustering", prefix, userID)
		return nil
	}

	var likes []repository.ListUserLikedItemsRow
	for offset := int32(0); int64(offset) < total; offset += MaxItemsLimit {
		remaining := total - int64(offset)
		limit := MaxItemsLimit
		if remaining < int64(MaxItemsLimit) {
			limit = int32(remaining)
		}
		batch, err := h.Repo.ListUserLikedItems(ctx, repository.ListUserLikedItemsParams{
			UserID: userID,
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			return fmt.Errorf("listing liked items (offset %d): %w", offset, err)
		}
		likes = append(likes, batch...)
		if int64(len(batch)) < int64(limit) {
			break
		}
	}

	log.Printf("%s fetched %d/%d liked items for user %s", prefix, len(likes), total, userID)

	var points [][]float64
	for _, like := range likes {
		embs, err := h.Repo.GetItemEmbeddingsByItemID(ctx, like.ID)
		if err != nil {
			return fmt.Errorf("failed to fetch embeddings for item %s: %w", like.ID, err)
		}
		for _, e := range embs {
			vec := vectorToFloat64s(*e.Embedding)
			points = append(points, vec)
		}
	}

	if len(points) == 0 {
		return errors.New("no item embeddings available, skipping user profile embedding")
	}

	centroids, err := ClusterizeEmbeddings(points)
	if err != nil {
		return fmt.Errorf("failed to clusterize embeddings for user %s: %w", userID, err)
	}

	clusters, err := h.Repo.GetUserEmbeddingClustersByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("fetching clusters for cleanup: %w", err)
	}
	for _, c := range clusters {
		if err := h.Repo.DeleteUserEmbeddingCluster(ctx, repository.DeleteUserEmbeddingClusterParams{
			UserID:    c.UserID,
			ClusterID: c.ClusterID,
		}); err != nil {
			return fmt.Errorf("deleting cluster %d: %w", c.ClusterID, err)
		}
	}

	for i, cent := range centroids {
		r, err := h.Repo.CreateUserEmbeddingCluster(ctx, repository.CreateUserEmbeddingClusterParams{
			UserID:      userID,
			ClusterID:   int32(i),
			Centroid:    cent.Vector,
			MemberCount: cent.MemberCount,
		})
		if err != nil {
			return fmt.Errorf("failed to insert cluster %d for user %s", i, userID)
		}
		log.Printf("%s created user embedding cluster %d for user %s", prefix, r.ClusterID, userID)
	}

	log.Printf("%s synched embedding clusters for user %s", prefix, userID)

	task, _ := NewCacheUserTask(userID)
	tResp, err := h.Client.Enqueue(task, asynq.Queue("default"))
	if err != nil {
		return fmt.Errorf("failed to queue caching task for user %s", userID)
	}
	log.Printf(
		"%s queued caching task %s for user %s ",
		prefix, tResp.ID, userID,
	)
	return nil
}

func (h *Handler) HandleUserSync(ctx context.Context, t *asynq.Task) error {
	prefix := fmt.Sprintf(
		"task %s -",
		t.ResultWriter().TaskID(),
	)
	count, err := h.Repo.CountUsers(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to count users: %v", err)
	}
	if count == 0 {
		return nil
	}
	for offset := int64(0); offset < count; offset += MaxLimit {
		limit := MaxLimit
		if remaining := count - offset; remaining < MaxLimit {
			limit = int(remaining)
		}

		users, err := h.Repo.ListUsers(ctx, repository.ListUsersParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		if err != nil {
			return fmt.Errorf("failed to list users: %v", err)
		}

		for _, user := range users {
			task, err := NewCacheUserTask(user.ID)
			if err != nil {
				return err
			}

			ti, err := h.Client.Enqueue(task, asynq.Queue("critical"))
			if err != nil {
				return err
			}
			log.Printf("%s queued caching task %s for user %s", prefix, ti.ID, user.ID)

		}

	}
	return nil
}
