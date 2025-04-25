package handler

import (
	"time"

	"github.com/google/uuid"
	"github.com/rhajizada/gazette/internal/typeext"
)

type PageParams struct {
	Limit  int32
	Offset int32
}

type CreateFeedRequest struct {
	FeedURL string `json:"feedURL"`
}

type Feed struct {
	ID              uuid.UUID       `json:"id"`
	Title           *string         `json:"title,omitempty"`
	Description     *string         `json:"description,omitempty"`
	Link            *string         `json:"link,omitempty"`
	FeedLink        string          `json:"feed_link"`
	Links           []string        `json:"links,omitempty"`
	UpdatedParsed   *time.Time      `json:"updated_parsed,omitempty"`
	PublishedParsed *time.Time      `json:"published_parsed,omitempty"`
	Authors         typeext.Authors `json:"authors,omitempty"`
	Language        *string         `json:"language,omitempty"`
	Image           any             `json:"image,omitempty"`
	Copyright       *string         `json:"copyright,omitempty"`
	Generator       *string         `json:"generator,omitempty"`
	Categories      []string        `json:"categories,omitempty"`
	FeedType        *string         `json:"feed_type,omitempty"`
	FeedVersion     *string         `json:"feed_version,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	LastUpdatedAt   time.Time       `json:"last_updated_at"`
	Subscribed      bool            `json:"subscribed"`
	SubscribedAt    *time.Time      `json:"subscribed_at,omitempty"`
}

// Item is the common model for items in API responses
type Item struct {
	ID              uuid.UUID       `json:"id"`
	FeedID          uuid.UUID       `json:"feed_id"`
	Title           *string         `json:"title,omitempty"`
	Description     *string         `json:"description,omitempty"`
	Content         *string         `json:"content,omitempty"`
	Link            string          `json:"link"`
	Links           []string        `json:"links,omitempty"`
	UpdatedParsed   *time.Time      `json:"updated_parsed,omitempty"`
	PublishedParsed *time.Time      `json:"published_parsed,omitempty"`
	Authors         typeext.Authors `json:"authors,omitempty"`
	GUID            *string         `json:"guid,omitempty"`
	Image           any             `json:"image,omitempty"`
	Categories      []string        `json:"categories,omitempty"`
	Enclosures      any             `json:"enclosures,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	Liked           bool            `json:"liked"`
	LikedAt         *time.Time      `json:"liked_at,omitempty"`
}

// ListFeedsResponse wraps paginated feeds
type ListFeedsResponse struct {
	Limit      int32  `json:"limit"`
	Offset     int32  `json:"offset"`
	TotalCount int64  `json:"total_count"`
	Feeds      []Feed `json:"feeds"`
}

// ListItemsResponse wraps paginated items
type ListItemsResponse struct {
	Limit      int32  `json:"limit"`
	Offset     int32  `json:"offset"`
	TotalCount int64  `json:"total_count"`
	Items      []Item `json:"items"`
}
