package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rhajizada/gazette/internal/repository"
)

// ListCollections retrieves a paginated list of collections for a user.
func (s *Service) ListCollections(ctx context.Context, r repository.ListCollectionsByUserParams) (*ListCollectionsResponse, error) {
	total, err := s.Repo.CountCollectionsByUserID(ctx, r.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			total = 0
		} else {
			return nil, NewError(
				"failed to count collections",
				http.StatusInternalServerError,
			)
		}
	}

	var rows []repository.Collection

	if total == 0 {
		rows = make([]repository.Collection, 0)
	} else {
		rows, err = s.Repo.ListCollectionsByUser(ctx, repository.ListCollectionsByUserParams{
			UserID: r.UserID,
			Limit:  r.Limit,
			Offset: r.Offset,
		})
		if err != nil {
			return nil, NewError("failed to list collections", http.StatusInternalServerError)
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
			return nil, NewError(
				fmt.Sprintf("collection %s already exists", col.Name),
				http.StatusConflict,
			)
		}
		return nil, NewError(
			fmt.Sprintf("failed to create collection %s", col.Name),
			http.StatusInternalServerError,
		)
	}
	return &Collection{
		ID:          col.ID,
		Name:        col.Name,
		CreatedAt:   col.CreatedAt,
		LastUpdated: col.LastUpdated,
	}, nil
}

// GetCollectionByID retrieves a single collection by ID.
func (s *Service) GetCollectionByID(ctx context.Context, collectionID uuid.UUID) (*Collection, error) {
	col, err := s.Repo.GetCollectionByID(ctx, collectionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewError(
				fmt.Sprintf("collection %s not found", collectionID),
				http.StatusNotFound,
			)
		} else {
			return nil, NewError(
				fmt.Sprintf("failed to fetch collection %s", collectionID),
				http.StatusInternalServerError,
			)
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
			return NewError(
				fmt.Sprintf("collection %s not found", collectionID),
				http.StatusNotFound,
			)
		} else {
			return NewError(
				fmt.Sprintf("failed to delete collection %s", collectionID),
				http.StatusInternalServerError,
			)
		}
	}
	return nil
}

// AddItemToCollection adds an item to a collection.
func (s *Service) AddItemToCollection(ctx context.Context, r repository.AddItemToCollectionParams) (*AddItemToCollectionResponse, error) {
	rec, err := s.Repo.AddItemToCollection(ctx, r)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return nil, NewError(
					fmt.Sprintf("item %s is already in collection %s", r.ItemID, r.CollectionID),
					http.StatusBadRequest,
				)
			case pgerrcode.ForeignKeyViolation:
				switch pgErr.ConstraintName {
				case "collection_items_collection_id_fkey":
					return nil, NewError(
						fmt.Sprintf("collection %s not found", r.CollectionID),
						http.StatusBadRequest,
					)
				case "collection_items_item_id_fkey":
					return nil, NewError(
						fmt.Sprintf("item %s not found", r.CollectionID),
						http.StatusBadRequest,
					)
				default:
					return nil, NewError("invalid reference", http.StatusBadRequest)
				}
			}
		} else {
			return nil, NewError(
				fmt.Sprintf("failed to update collection %s", r.CollectionID),
				http.StatusInternalServerError,
			)
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
			return NewError(
				fmt.Sprintf("item %s is not in collection %s", r.ItemID, r.CollectionID),
				http.StatusBadRequest,
			)
		} else {
			return NewError(
				fmt.Sprintf("failed to remove item %s from collection %s", r.ItemID, r.CollectionID),
				http.StatusInternalServerError,
			)
		}
	}
	return nil
}

// ListCollectionItems retrieves paginated items in a collection, including like status.
func (s *Service) ListCollectionItems(ctx context.Context, r ListCollectionItemsRequest) (*ListItemsResponse, error) {
	total, err := s.Repo.CountItemsInCollection(ctx, r.CollectionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			total = 0
		} else {
			return nil, NewError(
				fmt.Sprintf("failed counting items in collection %s", r.CollectionID),
				http.StatusInternalServerError)
		}
	}

	var rows []repository.ListItemsInCollectionRow

	if total == 0 {
		rows = make([]repository.ListItemsInCollectionRow, 0)
	} else {
		rows, err = s.Repo.ListItemsInCollection(ctx, repository.ListItemsInCollectionParams{
			CollectionID: r.CollectionID,
			Limit:        r.Limit,
			Offset:       r.Offset,
		})
		if err != nil {
			return nil, NewError(
				fmt.Sprintf("failed to list items in colllection %s", r.CollectionID),
				http.StatusInternalServerError,
			)
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

	return &ListItemsResponse{
		Limit:      r.Limit,
		Offset:     r.Offset,
		TotalCount: total,
		Items:      items,
	}, nil
}
