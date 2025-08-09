package cache

import "github.com/google/uuid"

// SuggestedItem now holds a uuid.UUID for ID.
type SuggestedItem struct {
	ID         uuid.UUID `json:"id"`
	Freshness  float64   `json:"freshness"`
	Similarity float64   `json:"similarity"`
	Score      float64   `json:"score"`
}
