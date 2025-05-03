package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rhajizada/gazette/internal/repository"
)

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
func (s *Service) ListCollections(ctx context.Context, r repository.ListCollectionsByUserParams) (*ListCollectionsResponse, error) {
	total, err := s.Repo.CountCollectionsByUserID(ctx, r.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, Err
		}
	}

	rows, err := s.Repo.ListCollectionsByUser(ctx, repository.ListCollectionsByUserParams{
		UserID: r.UserID,
		Limit:  r.Limit,
		Offset: r.Offset,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, Err
		}
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
func (s *Service) CreateCollection(ctx context.Context, r repository.CreateCollectionParams) (*Collection, error) {
	col, err := s.Repo.CreateCollection(ctx, r)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return nil, ErrAlreadyExists
		}
		return nil, Err
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, Err
		}
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
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		} else {
			return Err
		}
	}
	return nil
}

// AddItemToCollection adds an item to a collection.
func (s *Service) AddItemToCollection(ctx context.Context, r repository.AddItemToCollectionParams) (*AddItemToCollectionResponse, error) {
	rec, err := s.Repo.AddItemToCollection(ctx, r)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return nil, ErrBadInput
		} else {
			return nil, err
		}
	}
	return &AddItemToCollectionResponse{
		AddedAt: rec.AddedAt,
	}, nil
}

// RemoveItemFromCollection removes an item from a collection.
func (s *Service) RemoveItemFromCollection(ctx context.Context, r repository.RemoveItemFromCollectionParams) error {
	if err := s.Repo.RemoveItemFromCollection(ctx, r); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		} else {
			return Err
		}
	}
	return nil
}

// ListCollectionItems retrieves paginated items in a collection, including like status.
func (s *Service) ListCollectionItems(ctx context.Context, r ListCollectionItemsRequest) (*ListCollectionItemsResponse, error) {
	total, err := s.Repo.CountItemsInCollection(ctx, r.CollectionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, Err
		}
	}

	rows, err := s.Repo.ListItemsInCollection(ctx, repository.ListItemsInCollectionParams{
		CollectionID: r.CollectionID,
		Limit:        r.Limit,
		Offset:       r.Offset,
	})
	if err != nil && err != sql.ErrNoRows {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, Err
		}
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
