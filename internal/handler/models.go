package handler

import "github.com/rhajizada/gazette/internal/repository"

type CreateFeedRequest struct {
	FeedURL string `json:"feedURL"`
}

type ListFeedsResponse struct {
	Limit      int32             `json:"limit"`
	Offset     int32             `json:"offset"`
	TotalCount int64             `json:"totalCount"`
	Feeds      []repository.Feed `json:"feeds"`
}
