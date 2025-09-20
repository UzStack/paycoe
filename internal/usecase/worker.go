package usecase

import (
	"github.com/JscorpTech/paymento/internal/config"
	"github.com/JscorpTech/paymento/internal/domain"
	"go.uber.org/zap"
)

func InitWorker(log *zap.Logger, tasks <-chan domain.Task, cfg *config.Config) error {
	for range cfg.Workers {
		go Worker(tasks, log)
	}
	log.Info("Workers running: ", zap.Int("workers", cfg.Workers))
	return nil
}

func Worker(tasks <-chan domain.Task, log *zap.Logger) error {
	for task := range tasks {
		if payload, err := task.Paylod(); err == nil {
			log.Info("Yangi task: ", zap.Any("payload", payload))
		}
	}
	return nil
}
