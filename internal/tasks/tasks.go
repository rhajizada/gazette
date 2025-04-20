package tasks

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const (
	TypeFeedSync = "feed:sync"
)

type FeedSyncPayload struct {
	FeedID uuid.UUID
}

func NewFeedSyncTask(feedID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(FeedSyncPayload{FeedID: feedID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeFeedSync, payload), nil
}
