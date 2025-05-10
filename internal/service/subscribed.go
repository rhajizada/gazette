package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/rhajizada/gazette/internal/repository"
)

// ListUserSubscribedItems returns paginated list of recent items from subscribed user feeds
func (s *Service) ListUserSubscribedItems(ctx context.Context, r repository.ListSubscribedItemsByUserParams) (*ListItemsResponse, error) {
	total, err := s.Repo.CountSubscribedItemsByUser(ctx, r.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			total = 0
		} else {
			return nil, NewError(
				"failed to count subscribed items",
				http.StatusInternalServerError,
			)
		}
	}
	var rows []repository.ListSubscribedItemsByUserRow
	if total == 0 {
		rows = make([]repository.ListSubscribedItemsByUserRow, 0)
	} else {
		rows, err = s.Repo.ListSubscribedItemsByUser(ctx, r)
		if err != nil {
			return nil, NewError(
				fmt.Sprintf("failed to list subscribed items: %v", err),
				http.StatusInternalServerError,
			)
		}
	}

	items := make([]Item, total)
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
