package handler

import "github.com/rhajizada/gazette/internal/repository"

type PageParams struct {
	Limit  int32
	Offset int32
}

type CreateFeedRequest struct {
	FeedURL string `json:"feedURL"`
}

type ListFeedsResponse struct {
	Limit      int32             `json:"limit"`
	Offset     int32             `json:"offset"`
	TotalCount int64             `json:"totalCount"`
	Feeds      []repository.Feed `json:"feeds"`
}

type ListItemsResponse struct {
	Limit      int32             `json:"limit"`
	Offset     int32             `json:"offset"`
	TotalCount int64             `json:"totalCount"`
	Items      []repository.Item `json:"items"`
}
