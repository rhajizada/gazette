package service

import (
	"github.com/hibiken/asynq"
	"github.com/rhajizada/gazette/internal/repository"
)

type Service struct {
	Repo   repository.Queries
	Client *asynq.Client
}

func New(r *repository.Queries, c *asynq.Client) *Service {
	return &Service{
		Repo:   *r,
		Client: c,
	}
}
