package workers

import "go.uber.org/zap"

type Task interface {
	Paylod() any
}

type WebhookTask struct {
	OrderID int
	Amount  int
	Url     string
}

func (w WebhookTask) Paylod() any {
	return w
}

func InitWorker(log *zap.Logger, tasks <-chan Task) error {
	for i := range 10 {
		log.Info("Worker running", zap.Int("i", i))
		go Worker(tasks, log)
	}
	return nil
}

func Worker(tasks <-chan Task, log *zap.Logger) error {
	for task := range tasks {
		log.Info("Yangi task: ", zap.Any("task", task))
	}
	return nil
}
