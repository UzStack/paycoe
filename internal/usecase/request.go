package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func WebhookRequest(url string, data map[string]any, log *zap.Logger, retry int) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	response, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Error("webhook error:", zap.Int("status", response.StatusCode))
		if retry < 3 {
			time.Sleep(5 * time.Second)
			return WebhookRequest(url, data, log, retry+1)
		}
		return errors.New("failed send webhook " + response.Status)
	}
	log.Info("webhook success:", zap.Int("status", response.StatusCode))
	return nil
}

