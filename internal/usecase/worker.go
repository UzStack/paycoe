package usecase

import (
	"context"
	"database/sql"

	"github.com/JscorpTech/paymento/internal/config"
	"github.com/JscorpTech/paymento/internal/domain"
	"github.com/JscorpTech/paymento/internal/repository"
	"go.uber.org/zap"
)

func InitWorker(ctx context.Context, log *zap.Logger, tasks <-chan domain.Task, cfg *config.Config, db *sql.DB) error {
	for range cfg.Workers {
		go Worker(ctx, tasks, log, cfg, db)
	}
	log.Info("Workers running: ", zap.Int("workers", cfg.Workers))
	return nil
}

func Worker(ctx context.Context, tasks <-chan domain.Task, log *zap.Logger, cfg *config.Config, db *sql.DB) error {
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
			if err := WebhookRequest(cfg.WebhookURL, map[string]any{
				"transaction_id": payload.TransID,
				"amount":         payload.Amount,
			}, log, 1); err != nil {
				log.Error(err.Error())
			}
			repository.ConfirmTransaction(db, payload.TransID)
		}
	}
}
