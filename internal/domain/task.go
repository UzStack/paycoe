package domain

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
