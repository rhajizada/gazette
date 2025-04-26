package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rhajizada/gazette/internal/repository"
)

// Collection represents a user's collection of items
type Collection struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}

// CreateCollectionRequest wraps parameters to create collection
type CreateCollectionRequest struct {
	repository.CreateCollectionParams
}

// ListCollectionsRequest wraps parameters to list collections
type ListCollectionsRequest struct {
	repository.ListCollectionsByUserParams
}

// AddItemToCollectionRequest wraps parameters to add item to specified collection
type AddItemToCollectionRequest struct {
	repository.AddItemToCollectionParams
}

// RemoveItemFromCollectionRequest wraps parameters to remove item from  specified collection
type RemoveItemFromCollectionRequest struct {
	repository.RemoveItemFromCollectionParams
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

// ListCollectionItemsResponse wraps a paginated list of items in a collection
type ListCollectionItemsResponse struct {
	Limit      int32  `json:"limit"`
	Offset     int32  `json:"offset"`
	TotalCount int64  `json:"total_count"`
	Items      []Item `json:"items"`
}

// AddItemToCollectionResponse
type AddItemToCollectionResponse struct {
	AddedAt time.Time `json:"added_at"`
}

// ListCollections retrieves a paginated list of collections for a user.
func (s *Service) ListCollections(ctx context.Context, r ListCollectionsRequest) (*ListCollectionsResponse, error) {
	total, err := s.Repo.CountCollectionsByUserID(ctx, r.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed counting collections: %w", err)
	}

	rows, err := s.Repo.ListCollectionsByUser(ctx, repository.ListCollectionsByUserParams{
		UserID: r.UserID,
		Limit:  r.Limit,
		Offset: r.Offset,
	})
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed listing collections: %w", err)
	}

	cols := make([]Collection, len(rows))
	for i, c := range rows {
		cols[i] = Collection{
			ID:          c.ID,
			Name:        c.Name,
			CreatedAt:   c.CreatedAt,
			LastUpdated: c.LastUpdated,
		}
	}

	return &ListCollectionsResponse{
		Limit:       r.Limit,
		Offset:      r.Offset,
		TotalCount:  total,
		Collections: cols,
	}, nil
}

// CreateCollection creates a new collection.
func (s *Service) CreateCollection(ctx context.Context, r CreateCollectionRequest) (*Collection, error) {
	col, err := s.Repo.CreateCollection(ctx, r.CreateCollectionParams)
	if err != nil {
		return nil, fmt.Errorf("failed creating collection: %w", err)
	}
	return &Collection{
		ID:          col.ID,
		Name:        col.Name,
		CreatedAt:   col.CreatedAt,
		LastUpdated: col.LastUpdated,
	}, nil
}

// GetCollection retrieves a single collection by ID.
func (s *Service) GetCollection(ctx context.Context, collectionID uuid.UUID) (*Collection, error) {
	col, err := s.Repo.GetCollectionByID(ctx, collectionID)
	if err != nil {
		return nil, fmt.Errorf("collection not found: %w", err)
	}
	return &Collection{
		ID:          col.ID,
		Name:        col.Name,
		CreatedAt:   col.CreatedAt,
		LastUpdated: col.LastUpdated,
	}, nil
}

// DeleteCollection deletes a collection by ID.
func (s *Service) DeleteCollection(ctx context.Context, collectionID uuid.UUID) error {
	if err := s.Repo.DeleteCollectionByID(ctx, collectionID); err != nil {
		return fmt.Errorf("failed deleting collection: %w", err)
	}
	return nil
}

// AddItemToCollection adds an item to a collection.
func (s *Service) AddItemToCollection(ctx context.Context, r AddItemToCollectionRequest) (*AddItemToCollectionResponse, error) {
	rec, err := s.Repo.AddItemToCollection(ctx, r.AddItemToCollectionParams)
	if err != nil {
		return nil, fmt.Errorf("failed adding item to collection: %w", err)
	}
	return &AddItemToCollectionResponse{
		AddedAt: rec.AddedAt,
	}, nil
}

// RemoveItemFromCollection removes an item from a collection.
func (s *Service) RemoveItemFromCollection(ctx context.Context, r RemoveItemFromCollectionRequest) error {
	if err := s.Repo.RemoveItemFromCollection(ctx, r.RemoveItemFromCollectionParams); err != nil {
		return fmt.Errorf("failed removing item from collection: %w", err)
	}
	return nil
}

// ListCollectionItems retrieves paginated items in a collection, including like status.
func (s *Service) ListCollectionItems(ctx context.Context, r ListCollectionItemsRequest) (*ListCollectionItemsResponse, error) {
	total, err := s.Repo.CountItemsInCollection(ctx, r.CollectionID)
	if err != nil {
		return nil, fmt.Errorf("failed counting items in collection: %w", err)
	}

	rows, err := s.Repo.ListItemsInCollection(ctx, repository.ListItemsInCollectionParams{
		CollectionID: r.CollectionID,
		Limit:        r.Limit,
		Offset:       r.Offset,
	})
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed listing items: %w", err)
	}

	items := make([]Item, len(rows))
	for i, row := range rows {
		// Determine like status
		liked := false
		var likedAt *time.Time
		if like, err := s.Repo.GetUserLike(ctx, repository.GetUserLikeParams{UserID: r.UserID, ItemID: row.ID}); err == nil {
			liked = true
			likedAt = &like.LikedAt
		}

		// Map authors
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
			Liked:           liked,
			LikedAt:         likedAt,
		}
	}

	return &ListCollectionItemsResponse{
		Limit:      r.Limit,
		Offset:     r.Offset,
		TotalCount: total,
		Items:      items,
	}, nil
}
