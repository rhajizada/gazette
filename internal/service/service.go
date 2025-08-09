package service

import (
	"github.com/hibiken/asynq"
	"github.com/rhajizada/gazette/internal/cache"
	"github.com/rhajizada/gazette/internal/repository"
)

type Service struct {
	Repo   repository.Queries
	Client *asynq.Client
	Cache  *cache.Cache
}

func New(repo *repository.Queries, client *asynq.Client, cache *cache.Cache) *Service {
	return &Service{
		Repo:   *repo,
		Client: client,
		Cache:  cache,
	}
}
