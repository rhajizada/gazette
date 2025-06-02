package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/rhajizada/gazette/internal/repository"
)

// ListUserSuggestedItems returns paginated personalized items for the user,
// scored and ordered by similarity × freshness.
func (s *Service) ListUserSuggestedItems(ctx context.Context, r repository.ListSuggestedItemsByUserParams) (*ListItemsResponse, error) {
	cachedItems, err := s.Cache.GetUserSuggestions(ctx, r.UserID)
	var items []Item
	var total int64

	if err != nil || len(cachedItems) == 0 {
		total, err = s.Repo.CountSuggestedItemsByUser(ctx, r.UserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				total = 0
			} else {
				return nil, NewError(
					"failed to count suggested items",
					http.StatusInternalServerError,
				)
			}
		}

		var rows []repository.ListSuggestedItemsByUserRow
		if total == 0 {
			rows = make([]repository.ListSuggestedItemsByUserRow, 0)
		} else {
			rows, err = s.Repo.ListSuggestedItemsByUser(ctx, r)
			if err != nil {
				return nil, NewError(
					fmt.Sprintf("failed to list suggested items: %v", err),
					http.StatusInternalServerError,
				)
			}
		}

		items = make([]Item, len(rows))
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
				Liked:           false,
				LikedAt:         nil,
			}
		}
	} else {
		total = int64(len(cachedItems))
		start := int(r.Offset)
		limit := int(r.Limit)

		end := min(start+limit, int(total))

		pageSize := end - start
		items = make([]Item, pageSize)
		for i, j := range cachedItems[start:end] {
			t, err := s.GetItem(ctx, GetItemRequest{UserID: r.UserID, ItemID: j.ID})
			if err == nil {
				items[i] = *t
			}
		}
	}

	return &ListItemsResponse{
		Limit:      r.Limit,
		Offset:     r.Offset,
		TotalCount: total,
		Items:      items,
	}, nil
}
