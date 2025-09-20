package usecase

import (
	"context"

	"github.com/JscorpTech/paymento/internal/config"
	"github.com/JscorpTech/paymento/internal/domain"
	"go.uber.org/zap"
)

func InitWorker(ctx context.Context, log *zap.Logger, tasks <-chan domain.Task, cfg *config.Config) error {
	for range cfg.Workers {
		go Worker(ctx, tasks, log)
	}
	log.Info("Workers running: ", zap.Int("workers", cfg.Workers))
	return nil
}

func Worker(ctx context.Context, tasks <-chan domain.Task, log *zap.Logger) error {
	for {
		select {
		case <-ctx.Done():
			log.Info("Worker stop ")
			return nil
		case task, ok := <-tasks:
			if !ok {
				continue
			}
			payload := task.Paylod().(domain.WebhookTask)
			WebhookRequest("https://vesbini.felixits.uz/api/basket/", map[string]any{
				"order_id": payload.OrderID,
			}, log, 1)
		}
	}
}
