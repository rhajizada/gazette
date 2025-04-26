package service

import (
	"github.com/hibiken/asynq"
	"github.com/rhajizada/gazette/internal/repository"
)

type Service struct {
	Repo   repository.Queries
	Client *asynq.Client
}

func New(repo *repository.Queries, client *asynq.Client) *Service {
	return &Service{
		Repo:   *repo,
		Client: client,
	}
}
