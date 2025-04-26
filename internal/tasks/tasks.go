package tasks

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const (
	TypeSyncData  = "sync:data"
	TypeSyncFeed  = "sync:feed"
	TypeEmbedItem = "embed:item"
)

type SyncFeedPayload struct {
	FeedID uuid.UUID
}

type EmbedItemPayload struct {
	ItemID uuid.UUID
}

func NewSyncDataTask() (*asynq.Task, error) {
	return asynq.NewTask(TypeSyncData, nil), nil
}

func NewSyncFeedTask(feedID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(SyncFeedPayload{FeedID: feedID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSyncFeed, payload), nil
}

func NewEmbedItemTask(itemID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(EmbedItemPayload{ItemID: itemID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeEmbedItem, payload), nil
}
