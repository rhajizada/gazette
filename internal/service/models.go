package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/rhajizada/gazette/internal/repository"
)

// ListFeedsRequest wraps parameters for listing feeds.
type ListFeedsRequest struct {
	UserID       uuid.UUID
	SubscbedOnly bool
	Offset       int32
	Limit        int32
}

// ExportFeedsRequest wraps parameters for exporting feeds.
type ExportFeedsRequest struct {
	UserID       uuid.UUID
	SubscbedOnly bool
}

// CreateFeedRequest wraps parameters to create or subscribe to a feed.
type CreateFeedRequest struct {
	FeedURL string
	UserID  uuid.UUID
}

// GetFeedRequest wraps parameters to retrieve a feed and its subscription status.
type GetFeedRequest struct {
	repository.GetUserFeedSubscriptionParams
}

// DeleteFeedRequest wraps parameters to delete a feed entirely.
type DeleteFeedRequest struct {
	FeedID uuid.UUID
}

// SubscribeToFeedRequest wraps parameters to subscribe a user to a feed.
type SubscribeToFeedRequest struct {
	repository.CreateUserFeedSubscriptionParams
}

// UnsubscribeFromFeedRequest wraps parameters to remove a subscription.
type UnsubscribeFromFeedRequest struct {
	repository.DeleteUserFeedSubscriptionParams
}

// ListFeedItemsRequest wraps parameters for listing feed items with like info.
type ListFeedItemsRequest struct {
	repository.ListItemsByFeedIDForUserParams
}

// ListFeedsResponse wraps paginated feeds
type ListFeedsResponse struct {
	Limit      int32  `json:"limit"`
	Offset     int32  `json:"offset"`
	TotalCount int64  `json:"total_count"`
	Feeds      []Feed `json:"feeds"`
}

type SubscibeToFeedResponse struct {
	SubscribedAt *time.Time `json:"subscribed_at"`
}

type Feed struct {
	ID              uuid.UUID  `json:"id"`
	Title           *string    `json:"title,omitempty"`
	Description     *string    `json:"description,omitempty"`
	Link            *string    `json:"link,omitempty"`
	FeedLink        string     `json:"feed_link"`
	Links           []string   `json:"links,omitempty"`
	UpdatedParsed   *time.Time `json:"updated_parsed,omitempty"`
	PublishedParsed *time.Time `json:"published_parsed,omitempty"`
	Authors         Authors    `json:"authors,omitempty"`
	Language        *string    `json:"language,omitempty"`
	Image           any        `json:"image,omitempty"`
	Copyright       *string    `json:"copyright,omitempty"`
	Generator       *string    `json:"generator,omitempty"`
	Categories      []string   `json:"categories,omitempty"`
	FeedType        *string    `json:"feed_type,omitempty"`
	FeedVersion     *string    `json:"feed_version,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	LastUpdatedAt   time.Time  `json:"last_updated_at"`
	Subscribed      bool       `json:"subscribed"`
	SubscribedAt    *time.Time `json:"subscribed_at,omitempty"`
}

// ListItemsByFeedIDRequest wraps parameters for listing items from a feed with user-specific like info.
// Embeds FeedID, UserID, Limit, and Offset.
type ListItemsByFeedIDRequest struct {
	repository.ListItemsByFeedIDForUserParams
}

// GetItemRequest wraps parameters to retrieve a single item and its like status.
type GetItemRequest struct {
	UserID uuid.UUID
	ItemID uuid.UUID
}

// Item is the common model for items in API responses
type Item struct {
	ID              uuid.UUID  `json:"id"`
	FeedID          uuid.UUID  `json:"feed_id"`
	Title           *string    `json:"title,omitempty"`
	Description     *string    `json:"description,omitempty"`
	Content         *string    `json:"content,omitempty"`
	Link            string     `json:"link"`
	Links           []string   `json:"links,omitempty"`
	UpdatedParsed   *time.Time `json:"updated_parsed,omitempty"`
	PublishedParsed *time.Time `json:"published_parsed,omitempty"`
	Authors         Authors    `json:"authors,omitempty"`
	GUID            *string    `json:"guid,omitempty"`
	Image           any        `json:"image,omitempty"`
	Categories      []string   `json:"categories,omitempty"`
	Enclosures      any        `json:"enclosures,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Liked           bool       `json:"liked"`
	LikedAt         *time.Time `json:"liked_at,omitempty"`
}

// ListItemsResponse wraps paginated items
type ListItemsResponse struct {
	Limit      int32  `json:"limit"`
	Offset     int32  `json:"offset"`
	TotalCount int64  `json:"total_count"`
	Items      []Item `json:"items"`
}

// LikeItemResponse wraps response which inlude liked at time
type LikeItemResponse struct {
	LikedAt time.Time `json:"liked_at"`
}

// Collection represents a user's collection of items
type Collection struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}

// ListCollectionItemsRequest wraps parameters to list collection items
type ListCollectionItemsRequest struct {
	UserID uuid.UUID
	repository.ListItemsInCollectionParams
}

// ListCollectionsResponse wraps a paginated list of collections for a user
type ListCollectionsResponse struct {
	Limit       int32        `json:"limit"`
	Offset      int32        `json:"offset"`
	TotalCount  int64        `json:"total_count"`
	Collections []Collection `json:"collections"`
}

// AddItemToCollectionResponse
type AddItemToCollectionResponse struct {
	AddedAt time.Time `json:"added_at"`
}

// ListCategoriesResponse
type ListCategoriesResponse struct {
	Limit      int32    `json:"limit"`
	Offset     int32    `json:"offset"`
	TotalCount int64    `json:"total_count"`
	Categories []string `json:"categories"`
}

// Person represents an RSS‐feed author for documentation.
// swagger:model Person
type Person struct {
	// example: Jane Doe
	Name string `json:"name"`
	// example: jane@example.com
	Email string `json:"email,omitempty"`
}

// Authors is a list of Person.
// swagger:model Authors
type Authors []Person
