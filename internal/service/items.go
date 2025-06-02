package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rhajizada/gazette/internal/repository"
	"github.com/rhajizada/gazette/internal/workers"
)

// ListUserLikedItems returns paginated items the user has liked, with liked timestamps.
func (s *Service) ListUserLikedItems(ctx context.Context, r repository.ListUserLikedItemsParams) (*ListItemsResponse, error) {
	total, err := s.Repo.CountLikedItems(ctx, r.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			total = 0
		} else {
			return nil, NewError(
				"failed to count liked items",
				http.StatusInternalServerError,
			)
		}
	}

	var rows []repository.ListUserLikedItemsRow
	if total == 0 {
		rows = make([]repository.ListUserLikedItemsRow, 0)
	} else {
		rows, err = s.Repo.ListUserLikedItems(ctx, r)
		if err != nil {
			return nil, NewError(
				"failed to list liked items",
				http.StatusInternalServerError,
			)
		}
	}

	items := make([]Item, len(rows))
	for i, row := range rows {
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
			return nil, NewError(
				fmt.Sprintf("item %s not found", r.ItemID),
				http.StatusNotFound,
			)
		} else {
			return nil, NewError(
				fmt.Sprintf("failed to fetch item %s", r.ItemID),
				http.StatusInternalServerError,
			)
		}
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

// ListSimiliarItemsByID retrieves top 10 most similiar items to given one.
func (s *Service) ListSimiliarItemsByID(ctx context.Context, r repository.ListSimilarItemsByItemIDForUserParams) (*ListItemsResponse, error) {
	_, err := s.GetItem(ctx, GetItemRequest{
		ItemID: r.ItemID,
		UserID: r.UserID,
	})
	if err != nil {
		return nil, err
	}

	total, err := s.Repo.CountSimilarItemsByItemID(ctx, r.ItemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			total = 0
		} else {
			return nil, NewError(
				"failed to count similiar items",
				http.StatusInternalServerError,
			)
		}
	}

	var rows []repository.ListSimilarItemsByItemIDForUserRow
	if total == 0 {
		rows = make([]repository.ListSimilarItemsByItemIDForUserRow, 0)
	} else {
		rows, err = s.Repo.ListSimilarItemsByItemIDForUser(ctx, r)
		if err != nil {
			return nil, NewError(
				"failed to fetch similiar items",
				http.StatusInternalServerError,
			)
		}
	}

	items := make([]Item, len(rows))
	for i, row := range rows {
		likedAt := row.LikedAt
		liked := likedAt != nil

		if like, err := s.Repo.GetUserLike(ctx, repository.GetUserLikeParams{UserID: r.UserID, ItemID: r.ItemID}); err == nil {
			liked = true
			likedAt = &like.LikedAt
		}
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

// LikeItem marks an item as liked by the user.
func (s *Service) LikeItem(ctx context.Context, r repository.CreateUserLikeParams) (*LikeItemResponse, error) {
	like, err := s.Repo.CreateUserLike(ctx, r)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return nil, NewError(
				fmt.Sprintf("already liked item %s", r.ItemID),
				http.StatusBadRequest,
			)
		}
		return nil, NewError(
			fmt.Sprintf("failed to like item %s", r.ItemID),
			http.StatusInternalServerError,
		)
	}

	task, _ := workers.NewEmbedUserTask(r.UserID)
	_, err = s.Client.Enqueue(task, asynq.Queue("critical"))
	if err != nil {
		log.Printf("failed to queue embed task for user %s: %v", r.UserID, err)
	}

	err = s.Cache.RemoveUserSuggestion(ctx, r.UserID, r.ItemID)
	if err != nil {
		log.Printf("warning: failed to remove %s from user suggested items cache", r.ItemID)
	}

	return &LikeItemResponse{LikedAt: like.LikedAt}, nil
}

// UnlikeItem removes a like from an item.
func (s *Service) UnlikeItem(ctx context.Context, r repository.DeleteUserLikeParams) error {
	if err := s.Repo.DeleteUserLike(ctx, r); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewError(
				fmt.Sprintf("like for item %s not found", r.ItemID),
				http.StatusNotFound,
			)
		}
		return NewError(
			fmt.Sprintf("failed to unlike item %s", r.ItemID),
			http.StatusInternalServerError,
		)
	}
	task, _ := workers.NewEmbedUserTask(r.UserID)
	_, err := s.Client.Enqueue(task, asynq.Queue("critical"))
	if err != nil {
		log.Printf("failed to queue embed task for user %s: %v", r.UserID, err)
	}

	return nil
}

// ListItemCollections returns list of collections that item is in.
func (s *Service) ListItemCollections(ctx context.Context, r repository.ListCollectionsByItemIDParams) (*ListCollectionsResponse, error) {
	total, err := s.Repo.CountCollectionsByItemID(ctx, repository.CountCollectionsByItemIDParams{
		ItemID: r.ItemID,
		UserID: r.UserID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			total = 0
		} else {
			return nil, NewError(
				fmt.Sprintf("failed to count collection with item %s", r.ItemID),
				http.StatusInternalServerError,
			)
		}
	}

	var rows []repository.Collection
	if total == 0 {
		rows = []repository.Collection{}
	} else {
		rows, err = s.Repo.ListCollectionsByItemID(ctx, r)
		if err != nil {
			return nil, NewError(
				fmt.Sprintf("failed to list collections for item %s", r.ItemID),
				http.StatusInternalServerError,
			)
		}
	}

	cols := make([]Collection, len(rows))
	for i, row := range rows {
		cols[i] = Collection{
			ID:          row.ID,
			Name:        row.Name,
			CreatedAt:   row.CreatedAt,
			LastUpdated: row.LastUpdated,
		}
	}

	return &ListCollectionsResponse{
		Limit:       r.Limit,
		Offset:      r.Offset,
		TotalCount:  total,
		Collections: cols,
	}, nil
}
