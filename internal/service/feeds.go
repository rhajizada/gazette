package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mmcdole/gofeed"
	"github.com/rhajizada/gazette/internal/repository"
	"github.com/rhajizada/gazette/internal/typeext"
	"github.com/rhajizada/gazette/internal/workers"
)

// ListFeeds retrieves a paginated list of feeds, optionally only subscribed.
func (s *Service) ListFeeds(ctx context.Context, r ListFeedsRequest) (*ListFeedsResponse, error) {
	// count
	var total int64
	var err error
	if r.SubscbedOnly {
		total, err = s.Repo.CountFeedsByUserID(ctx, r.UserID)
	} else {
		total, err = s.Repo.CountFeeds(ctx)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			total = 0
		} else {
			return nil, NewError(
				"failed to count feeds",
				http.StatusInternalServerError,
			)
		}
	}

	var rows []repository.ListFeedsByUserIDRow

	if total == 0 {
		rows = make([]repository.ListFeedsByUserIDRow, 0)
	} else {
		rows, err = s.Repo.ListFeedsByUserID(ctx, repository.ListFeedsByUserIDParams{
			UserID:  r.UserID,
			Column2: r.SubscbedOnly,
			Offset:  r.Offset,
			Limit:   r.Limit,
		})
		if err != nil {
			return nil, NewError(
				"failed to list feeds",
				http.StatusInternalServerError,
			)
		}
	}

	feeds := make([]Feed, len(rows))
	for i, row := range rows {
		auths := make(Authors, len(row.Authors))
		for j, a := range row.Authors {
			auths[j] = Person{Name: a.Name, Email: a.Email}
		}

		feeds[i] = Feed{
			ID:              row.ID,
			Title:           row.Title,
			Description:     row.Description,
			Link:            row.Link,
			FeedLink:        row.FeedLink,
			Links:           row.Links,
			UpdatedParsed:   row.UpdatedParsed,
			PublishedParsed: row.PublishedParsed,
			Authors:         auths,
			Language:        row.Language,
			Image:           row.Image,
			Copyright:       row.Copyright,
			Generator:       row.Generator,
			Categories:      row.Categories,
			FeedType:        row.FeedType,
			FeedVersion:     row.FeedVersion,
			CreatedAt:       row.CreatedAt,
			LastUpdatedAt:   row.LastUpdatedAt,
			Subscribed:      row.SubscribedAt != nil,
			SubscribedAt:    row.SubscribedAt,
		}
	}

	return &ListFeedsResponse{
		Limit:      r.Limit,
		Offset:     r.Offset,
		TotalCount: total,
		Feeds:      feeds,
	}, nil
}

// CreateFeed creates a feed if needed, enqueues a sync task, and subscribes the user.
func (s *Service) CreateFeed(ctx context.Context, r CreateFeedRequest) (*Feed, error) {
	parser := gofeed.NewParser()
	remote, err := parser.ParseURL(r.FeedURL)
	if err != nil {
		return nil, NewError(
			fmt.Sprintf("invalid feed URL %s", r.FeedURL),
			http.StatusBadRequest,
		)
	}

	feed, err := s.Repo.CreateFeed(ctx, repository.CreateFeedParams{
		Title:           &remote.Title,
		Description:     &remote.Description,
		Link:            &remote.Link,
		FeedLink:        remote.FeedLink,
		Links:           remote.Links,
		UpdatedParsed:   remote.UpdatedParsed,
		PublishedParsed: remote.PublishedParsed,
		Authors:         typeext.Authors(remote.Authors),
		Language:        &remote.Language,
		Image:           remote.Image,
		Copyright:       &remote.Copyright,
		Generator:       &remote.Generator,
		Categories:      remote.Categories,
		FeedType:        &remote.FeedType,
		FeedVersion:     &remote.FeedVersion,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return nil, NewError(
				fmt.Sprintf("feed %s already exists", r.FeedURL),
				http.StatusConflict,
			)
		}
		return nil, NewError(
			fmt.Sprintf("failed to load create feed %s", r.FeedURL),
			http.StatusInternalServerError,
		)
	}

	task, _ := workers.NewSyncFeedTask(feed.ID)
	s.Client.Enqueue(task, asynq.Queue("critical"))

	sub, err := s.Repo.CreateUserFeedSubscription(ctx, repository.CreateUserFeedSubscriptionParams{UserID: r.UserID, FeedID: feed.ID})
	if err != nil {
		return nil, NewError(
			fmt.Sprintf("failed to subscribe to feed %s", feed.ID),
			http.StatusInternalServerError,
		)
	}

	auths := make(Authors, len(feed.Authors))
	for i, a := range feed.Authors {
		auths[i] = Person{Name: a.Name, Email: a.Email}
	}

	return &Feed{
		ID:              feed.ID,
		Title:           feed.Title,
		Description:     feed.Description,
		Link:            feed.Link,
		FeedLink:        feed.FeedLink,
		Links:           feed.Links,
		UpdatedParsed:   feed.UpdatedParsed,
		PublishedParsed: feed.PublishedParsed,
		Authors:         auths,
		Language:        feed.Language,
		Image:           feed.Image,
		Copyright:       feed.Copyright,
		Generator:       feed.Generator,
		Categories:      feed.Categories,
		FeedType:        feed.FeedType,
		FeedVersion:     feed.FeedVersion,
		CreatedAt:       feed.CreatedAt,
		LastUpdatedAt:   feed.LastUpdatedAt,
		Subscribed:      true,
		SubscribedAt:    &sub.SubscribedAt,
	}, nil
}

// GetFeed retrieves a feed and the user's subscription status.
func (s *Service) GetFeed(ctx context.Context, r repository.GetUserFeedSubscriptionParams) (*Feed, error) {
	feed, err := s.Repo.GetFeedByID(ctx, r.FeedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewError(
				fmt.Sprintf("feed %s not found", r.FeedID),
				http.StatusNotFound,
			)
		} else {
			return nil, NewError(
				fmt.Sprintf("failed to fetch feed %s", r.FeedID),
				http.StatusInternalServerError,
			)
		}
	}

	// check subscription
	subAt := (*time.Time)(nil)
	if uf, err := s.Repo.GetUserFeedSubscription(ctx, r); err == nil {
		subAt = &uf.SubscribedAt
	}

	// map authors
	auths := make(Authors, len(feed.Authors))
	for i, a := range feed.Authors {
		auths[i] = Person{Name: a.Name, Email: a.Email}
	}

	return &Feed{
		ID:              feed.ID,
		Title:           feed.Title,
		Description:     feed.Description,
		Link:            feed.Link,
		FeedLink:        feed.FeedLink,
		Links:           feed.Links,
		UpdatedParsed:   feed.UpdatedParsed,
		PublishedParsed: feed.PublishedParsed,
		Authors:         auths,
		Language:        feed.Language,
		Image:           feed.Image,
		Categories:      feed.Categories,
		FeedType:        feed.FeedType,
		FeedVersion:     feed.FeedVersion,
		CreatedAt:       feed.CreatedAt,
		LastUpdatedAt:   feed.LastUpdatedAt,
		Subscribed:      subAt != nil,
		SubscribedAt:    subAt,
	}, nil
}

// DeleteFeed deletes a feed entirely.
func (s *Service) DeleteFeed(ctx context.Context, r DeleteFeedRequest) error {
	err := s.Repo.DeleteFeedByID(ctx, r.FeedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewError(
				fmt.Sprintf("feed %s not found", r.FeedID),
				http.StatusNotFound,
			)
		} else {
			return NewError(
				fmt.Sprintf("failed to delete feed %s", r.FeedID),
				http.StatusInternalServerError,
			)
		}
	}
	return nil
}

// SubscribeToFeed subscribes a user to a feed.
func (s *Service) SubscribeToFeed(ctx context.Context, r repository.CreateUserFeedSubscriptionParams) (*SubscibeToFeedResponse, error) {
	sub, err := s.Repo.CreateUserFeedSubscription(ctx, r)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return nil, NewError(
					fmt.Sprintf("already subscribed to feed %s", r.FeedID),
					http.StatusBadRequest,
				)
			case pgerrcode.ForeignKeyViolation:
				return nil, NewError(
					fmt.Sprintf("feed %s not found", r.FeedID),
					http.StatusBadRequest,
				)
			}
		}
		return nil, NewError(
			fmt.Sprintf("failed to subsribe to feed %s", r.FeedID),
			http.StatusInternalServerError,
		)
	}

	return &SubscibeToFeedResponse{
		SubscribedAt: &sub.SubscribedAt,
	}, nil
}

// UnsubscribeFromFeed removes a user's subscription.
func (s *Service) UnsubscribeFromFeed(ctx context.Context, r repository.DeleteUserFeedSubscriptionParams) error {
	err := s.Repo.DeleteUserFeedSubscription(ctx, r)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewError(
				fmt.Sprintf("user not subscribed to feed %s", r.FeedID),
				http.StatusBadRequest,
			)
		} else {
			return NewError(
				fmt.Sprintf("failed to unsubscrive from feed %s", r.FeedID),
				http.StatusInternalServerError,
			)
		}
	}
	return nil
}

// ListItemsByFeedID returns paginated items from a feed, including per-user like status.
func (s *Service) ListItemsByFeedID(ctx context.Context, r repository.ListItemsByFeedIDForUserParams) (*ListItemsResponse, error) {
	total, err := s.Repo.CountItemsByFeedID(ctx, r.FeedID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			total = 0
		} else {
			return nil, NewError(
				fmt.Sprintf("failed to list items in feed %s", r.FeedID),
				http.StatusInternalServerError,
			)
		}
	}

	var rows []repository.ListItemsByFeedIDForUserRow
	if total == 0 {
		rows = make([]repository.ListItemsByFeedIDForUserRow, 0)
	} else {
		rows, err = s.Repo.ListItemsByFeedIDForUser(ctx, r)
		if err != nil {
			return nil, NewError(
				fmt.Sprintf("failed to list items in feed %s", r.FeedID),
				http.StatusInternalServerError,
			)
		}
	}

	items := make([]Item, len(rows))
	for i, row := range rows {
		liked := row.LikedAt != nil

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
