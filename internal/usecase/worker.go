package usecase

import (
	"context"
	"database/sql"
	"time"

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
	go CloseTransactionWorker(ctx, db, log)
	return nil
}

func CloseTransactionWorker(ctx context.Context, db *sql.DB, log *zap.Logger) {
	for {
		select {
		case <-ctx.Done():
			log.Info("Worker stop transaction")
		default:
			transactions, err := repository.GetOldTransactions(db)
			if err != nil {
				log.Error("old transactions close error", zap.Error(err))
			}
			for _, transaction := range transactions {
				log.Info("transaction", zap.Any("amount", transaction["amount"]), zap.Any("transaction_id", transaction["transaction_id"]))
			}
			time.Sleep(1 * time.Minute)
		}
	}
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
