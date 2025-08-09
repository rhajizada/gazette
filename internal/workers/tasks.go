package workers

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const (
	TypeSyncData  = "sync:data"
	TypeSyncFeed  = "sync:feed"
	TypeEmbedItem = "embed:item"
	TypeCacheUser = "cache:user"
	TypeEmbedUser = "embed:user"
	TypeSyncUser  = "sync:user"
)

type SyncFeedPayload struct {
	FeedID uuid.UUID
}

type EmbedItemPayload struct {
	ItemID uuid.UUID
}

type CacheUserPayload struct {
	UserID uuid.UUID
}

type EmbedUserPayload struct {
	UserID uuid.UUID
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

func NewEmbedUserTask(userID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(EmbedUserPayload{UserID: userID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeEmbedUser, payload), nil
}

func NewCacheUserTask(userID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(CacheUserPayload{UserID: userID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeCacheUser, payload), nil
}

func NewSyncUserTask() (*asynq.Task, error) {
	return asynq.NewTask(TypeSyncUser, nil), nil
}
