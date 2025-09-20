package domain

type Task interface {
	Paylod() (map[string]any, error)
}

type WebhookTask struct {
	OrderID int
	Amount  int
	Url     string
}

func (w WebhookTask) Paylod() (map[string]any, error) {
	return map[string]any{
		"order_id": w.OrderID,
		"amount":   w.Amount,
		"url":      w.Url,
	}, nil
}
