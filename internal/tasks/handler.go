package tasks

import (
	"github.com/hibiken/asynq"
	"github.com/rhajizada/gazette/internal/config"
	"github.com/rhajizada/gazette/internal/repository"
)

const MaxLimit = 100

type Handler struct {
	Repo         repository.Queries
	Client       *asynq.Client
	OllamaConfig *config.OllamaConfig
}

func NewHandler(r *repository.Queries, c *asynq.Client, o *config.OllamaConfig) *Handler {
	return &Handler{
		Repo:         *r,
		Client:       c,
		OllamaConfig: o,
	}
}
