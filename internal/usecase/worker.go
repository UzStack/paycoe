package usecase

import (
	"github.com/JscorpTech/paymento/internal/domain"
	"go.uber.org/zap"
)

func InitWorker(log *zap.Logger, tasks <-chan domain.Task) error {
	for i := range 10 {
		log.Info("Worker running", zap.Int("i", i))
		go Worker(tasks, log)
	}
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
