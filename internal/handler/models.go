package handler

type CreateFeedRequest struct {
	FeedURL string `json:"feed_url"`
}

type CreateCollectionRequest struct {
	Name string `json:"name"`
}
