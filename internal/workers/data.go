package workers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/rhajizada/gazette/internal/repository"
)

func (h *Handler) HandleDataSync(ctx context.Context, t *asynq.Task) error {
	count, err := h.Repo.CountFeeds(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to list feeds: %v", err)
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
			return fmt.Errorf("failed to list feeds (offset %d): %v", offset, err)
		}

		for _, feed := range feeds {
			task, err := NewSyncFeedTask(feed.ID)
			if err != nil {
				return err
			}

			ti, err := h.Client.Enqueue(task, asynq.Queue("critical"))
			if err != nil {
				return err
			}
			log.Printf("queued sync task %s for feed %s", ti.ID, feed.ID)

		}

	}
	return nil
}
