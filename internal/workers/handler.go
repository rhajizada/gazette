package workers

import (
	"github.com/hibiken/asynq"
	"github.com/rhajizada/gazette/internal/cache"
	"github.com/rhajizada/gazette/internal/config"
	"github.com/rhajizada/gazette/internal/repository"
)

const MaxLimit = 100

type Handler struct {
	Repo         repository.Queries
	Client       *asynq.Client
	Cache        *cache.Cache
	OllamaConfig *config.OllamaConfig
}

func NewHandler(repo *repository.Queries, client *asynq.Client, cache *cache.Cache, ollamaCfg *config.OllamaConfig) *Handler {
	return &Handler{
		Repo:         *repo,
		Client:       client,
		Cache:        cache,
		OllamaConfig: ollamaCfg,
	}
}
