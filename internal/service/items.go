package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rhajizada/gazette/internal/repository"
)

// ListUserLikedItemsRequest wraps parameters for listing items liked by a user.
type ListUserLikedItemsRequest struct {
	repository.ListUserLikedItemsParams
}

// GetItemRequest wraps parameters to retrieve a single item and its like status.
type GetItemRequest struct {
	UserID uuid.UUID
	ItemID uuid.UUID
}

// LikeItemRequest wraps parameters to like an item.
type LikeItemRequest struct {
	repository.CreateUserLikeParams
}

// UnlikeItemRequest wraps parameters to unlike an item.
type UnlikeItemRequest struct {
	repository.DeleteUserLikeParams
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

// ListUserLikedItems returns paginated items the user has liked, with liked timestamps.
func (s *Service) ListUserLikedItems(ctx context.Context, r ListUserLikedItemsRequest) (*ListItemsResponse, error) {
	// count total liked items
	total, err := s.Repo.CountLikedItems(ctx, r.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed counting liked items: %w", err)
	}

	// fetch liked items
	rows, err := s.Repo.ListUserLikedItems(ctx, r.ListUserLikedItemsParams)
	if err != nil {
		return nil, fmt.Errorf("failed listing liked items: %w", err)
	}

	// map to service Item
	items := make([]Item, len(rows))
	for i, row := range rows {
		// map authors
		auths := make(Authors, len(row.Authors))
		for j, a := range row.Authors {
			auths[j] = Person{Name: a.Name, Email: a.Email}
		}

		items[i] = Item{
			ID:              row.ID,
			FeedID:          row.FeedID,
			Title:           row.Title,
			Description:     row.Description,
			Content:         row.Content,
			Link:            row.Link,
			Links:           row.Links,
			UpdatedParsed:   row.UpdatedParsed,
			PublishedParsed: row.PublishedParsed,
			Authors:         auths,
			GUID:            row.Guid,
			Image:           row.Image,
			Categories:      row.Categories,
			Enclosures:      row.Enclosures,
			CreatedAt:       row.CreatedAt,
			UpdatedAt:       row.UpdatedAt,
			Liked:           row.Liked,
			LikedAt:         &row.LikedAt,
		}
	}

	return &ListItemsResponse{
		Limit:      r.Limit,
		Offset:     r.Offset,
		TotalCount: total,
		Items:      items,
	}, nil
}

// GetItem retrieves a single item and its like status for the user.
func (s *Service) GetItem(ctx context.Context, r GetItemRequest) (*Item, error) {
	// fetch the item
	row, err := s.Repo.GetItemByID(ctx, r.ItemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("item not found: %w", err)
		}
		return nil, fmt.Errorf("failed fetching item: %w", err)
	}

	// determine like status
	liked := false
	var likedAt *time.Time
	if like, err := s.Repo.GetUserLike(ctx, repository.GetUserLikeParams{UserID: r.UserID, ItemID: r.ItemID}); err == nil {
		liked = true
		likedAt = &like.LikedAt
	}

	// map authors
	auths := make(Authors, len(row.Authors))
	for j, a := range row.Authors {
		auths[j] = Person{Name: a.Name, Email: a.Email}
	}

	item := Item{
		ID:              row.ID,
		FeedID:          row.FeedID,
		Title:           row.Title,
		Description:     row.Description,
		Content:         row.Content,
		Link:            row.Link,
		Links:           row.Links,
		UpdatedParsed:   row.UpdatedParsed,
		PublishedParsed: row.PublishedParsed,
		Authors:         auths,
		GUID:            row.Guid,
		Image:           row.Image,
		Categories:      row.Categories,
		Enclosures:      row.Enclosures,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
		Liked:           liked,
		LikedAt:         likedAt,
	}
	return &item, nil
}

// LikeItem marks an item as liked by the user.
func (s *Service) LikeItem(ctx context.Context, r LikeItemRequest) (*time.Time, error) {
	like, err := s.Repo.CreateUserLike(ctx, r.CreateUserLikeParams)
	if err != nil {
		return nil, fmt.Errorf("failed to like item: %w", err)
	}
	return &like.LikedAt, nil
}

// UnlikeItem removes a like from an item.
func (s *Service) UnlikeItem(ctx context.Context, r UnlikeItemRequest) error {
	if err := s.Repo.DeleteUserLike(ctx, r.DeleteUserLikeParams); err != nil {
		return fmt.Errorf("failed to unlike item: %w", err)
	}
	return nil
}
