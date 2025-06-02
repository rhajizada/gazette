package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rhajizada/gazette/internal/config"
)

type Cache struct {
	rdb *redis.Client
}

func New(conf *config.CacheConfig) (*Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &Cache{rdb: rdb}, nil
}

func (s *Cache) StoreUserSuggestions(
	ctx context.Context,
	userID uuid.UUID,
	items []SuggestedItem,
	ttl time.Duration,
) error {
	key := fmt.Sprintf("user:suggestions:%s", userID.String())
	if err := s.rdb.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis DEL %q: %w", key, err)
	}

	pipe := s.rdb.Pipeline()
	for _, it := range items {
		raw, err := json.Marshal(struct {
			Freshness  float64 `json:"freshness"`
			Similarity float64 `json:"similarity"`
			Score      float64 `json:"score"`
		}{
			Freshness:  it.Freshness,
			Similarity: it.Similarity,
			Score:      it.Score,
		})
		if err != nil {
			return fmt.Errorf("json marshal scores for item %q: %w", it.ID.String(), err)
		}
		pipe.HSet(ctx, key, it.ID.String(), raw)
	}
	if ttl > 0 {
		pipe.Expire(ctx, key, ttl)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("pipeline Exec for key %q: %w", key, err)
	}
	return nil
}

func (s *Cache) GetUserSuggestions(ctx context.Context, userID uuid.UUID) ([]SuggestedItem, error) {
	key := fmt.Sprintf("user:suggestions:%s", userID.String())
	rawMap, err := s.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("redis HGetAll for key %q: %w", key, err)
	}

	result := make([]SuggestedItem, 0, len(rawMap))
	for field, raw := range rawMap {
		itemUUID, err := uuid.Parse(field)
		if err != nil {
			log.Printf("warning: invalid UUID field %q in hash %q: %v\n", field, key, err)
			continue
		}
		var measures struct {
			Freshness  float64 `json:"freshness"`
			Similarity float64 `json:"similarity"`
			Score      float64 `json:"score"`
		}
		if err := json.Unmarshal([]byte(raw), &measures); err != nil {
			log.Printf("warning: could not unmarshal scores for item %q: %v\n", field, err)
			continue
		}
		result = append(result, SuggestedItem{
			ID:         itemUUID,
			Freshness:  measures.Freshness,
			Similarity: measures.Similarity,
			Score:      measures.Score,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})

	return result, nil
}
