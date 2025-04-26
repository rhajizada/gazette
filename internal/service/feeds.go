package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	"github.com/rhajizada/gazette/internal/repository"
	"github.com/rhajizada/gazette/internal/tasks"
	"github.com/rhajizada/gazette/internal/typeext"
)

// ListFeedsRequest wraps parameters for listing feeds.
// Embeds the repository-defined params, including UserID, Column2 (subscribed flag), Limit, and Offset.
type ListFeedsRequest struct {
	UserID       uuid.UUID
	SubscbedOnly bool
	Offset       int32
	Limit        int32
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
		return nil, fmt.Errorf("failed counting feeds: %w", err)
	}

	// fetch
	rows, err := s.Repo.ListFeedsByUserID(ctx, repository.ListFeedsByUserIDParams{
		UserID:  r.UserID,
		Column2: r.SubscbedOnly,
		Offset:  r.Offset,
		Limit:   r.Limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed listing feeds: %w", err)
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
	// parse remote feed
	parser := gofeed.NewParser()
	remote, err := parser.ParseURL(r.FeedURL)
	if err != nil {
		return nil, fmt.Errorf("invalid feed URL: %w", err)
	}

	// lookup or create feed record
	feed, err := s.Repo.GetFeedByFeedLink(ctx, remote.FeedLink)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			feed, err = s.Repo.CreateFeed(ctx, repository.CreateFeedParams{
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
				return nil, fmt.Errorf("failed creating feed: %w", err)
			}
			// enqueue sync
			task, _ := tasks.NewSyncFeedTask(feed.ID)
			s.Client.Enqueue(task)
		} else {
			return nil, fmt.Errorf("lookup failed: %w", err)
		}
	}

	// subscribe user
	sub, err := s.Repo.CreateUserFeedSubscription(ctx, repository.CreateUserFeedSubscriptionParams{UserID: r.UserID, FeedID: feed.ID})
	if err != nil {
		return nil, fmt.Errorf("failed subscribing: %w", err)
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
		return nil, fmt.Errorf("feed not found: %w", err)
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
	if err := s.Repo.DeleteFeedByID(ctx, r.FeedID); err != nil {
		return fmt.Errorf("failed deleting feed: %w", err)
	}
	return nil
}

// SubscribeToFeed subscribes a user to a feed.
func (s *Service) SubscribeToFeed(ctx context.Context, r repository.CreateUserFeedSubscriptionParams) (*SubscibeToFeedResponse, error) {
	sub, err := s.Repo.CreateUserFeedSubscription(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("failed subscribing: %w", err)
	}
	return &SubscibeToFeedResponse{
		SubscribedAt: &sub.SubscribedAt,
	}, nil
}

// UnsubscribeFromFeed removes a user's subscription.
func (s *Service) UnsubscribeFromFeed(ctx context.Context, r repository.DeleteUserFeedSubscriptionParams) error {
	if err := s.Repo.DeleteUserFeedSubscription(ctx, r); err != nil {
		return fmt.Errorf("failed unsubscribing: %w", err)
	}
	return nil
}

// ListItemsByFeedID returns paginated items from a feed, including per-user like status.
func (s *Service) ListItemsByFeedID(ctx context.Context, r repository.ListItemsByFeedIDForUserParams) (*ListItemsResponse, error) {
	// total count
	total, err := s.Repo.CountItemsByFeedID(ctx, r.FeedID)
	if err != nil {
		return nil, fmt.Errorf("failed counting items for feed %s: %w", r.FeedID, err)
	}

	// fetch rows with like info
	rows, err := s.Repo.ListItemsByFeedIDForUser(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("failed listing items for feed %s: %w", r.FeedID, err)
	}

	// map to service Item
	items := make([]Item, len(rows))
	for i, row := range rows {
		// determine like status
		liked := row.LikedAt != nil

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
