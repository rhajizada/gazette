package workers

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

func NewHandler(repo *repository.Queries, client *asynq.Client, ollamaCfg *config.OllamaConfig) *Handler {
	return &Handler{
		Repo:         *repo,
		Client:       client,
		OllamaConfig: ollamaCfg,
	}
}
