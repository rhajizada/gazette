package workers

import (
	"github.com/hibiken/asynq"
	"github.com/rhajizada/gazette/internal/config"
)

// GetConfig generates asynq server configuration.
func GetConfig(cfg *config.QueuesConfig) *asynq.Config {
	total := cfg.Critical + cfg.Default + cfg.Low
	return &asynq.Config{
		Concurrency: total,
		Queues: map[string]int{
			"critical": cfg.Critical,
			"default":  cfg.Default,
			"low":      cfg.Low,
		},
		StrictPriority: true,
	}
}
