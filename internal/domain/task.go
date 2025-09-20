package domain

type Task interface {
	Paylod() any
}

type WebhookTask struct {
	TransID int64
	Amount  int64
	Url     string
}

func (w WebhookTask) Paylod() any {
	return w
}
