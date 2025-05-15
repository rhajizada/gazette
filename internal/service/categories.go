package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/rhajizada/gazette/internal/repository"
)

// ListCategories retrieves a paginated list of collections for a user.
func (s *Service) ListCategories(ctx context.Context, r repository.ListDistinctCategoriesParams) (*ListCategoriesResponse, error) {
	total, err := s.Repo.CountDistinctCategories(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			total = 0
		} else {
			return nil, NewError(
				"failed to count categories",
				http.StatusInternalServerError,
			)
		}
	}

	var categories []string
	if total == 0 {
		categories = make([]string, 0)
	} else {
		preprocessed, err := s.Repo.ListDistinctCategories(ctx, r)
		if err != nil {
			return nil, NewError("failed to list collections", http.StatusInternalServerError)
		}
		categories = make([]string, len(preprocessed))
		for i, j := range preprocessed {
			categories[i] = j.(string)
		}
	}
	return &ListCategoriesResponse{
		Limit:      r.Limit,
		Offset:     r.Offset,
		TotalCount: total,
		Categories: categories,
	}, nil
}

// ListCategoryItems retrieves a paginated list of items that are in given category
func (s *Service) ListCategoryItems(ctx context.Context, r repository.ListItemsByCategoryForUserParams) (*ListItemsResponse, error) {
	total, err := s.Repo.CountItemsByCategoryForUser(ctx, r.Categories)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			total = 0
		} else {
			return nil, NewError(
				"failed to count items in categories",
				http.StatusInternalServerError,
			)
		}
	}
	var rows []repository.ListItemsByCategoryForUserRow
	if total == 0 {
		rows = make([]repository.ListItemsByCategoryForUserRow, 0)
	} else {
		rows, err = s.Repo.ListItemsByCategoryForUser(ctx, r)
		if err != nil {
			return nil, NewError(
				fmt.Sprintf("failed to list items in categories: %v", err),
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

		liked := row.LikedAt != nil

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
			LikedAt:         row.LikedAt,
		}
	}

	return &ListItemsResponse{
		Limit:      r.Limit,
		Offset:     r.Offset,
		TotalCount: total,
		Items:      items,
	}, nil
}
