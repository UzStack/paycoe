package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JscorpTech/paymento/internal/config"
	"github.com/JscorpTech/paymento/internal/domain"
	"github.com/JscorpTech/paymento/internal/repository"
	"go.uber.org/zap"
)

type Handler struct {
	DB    *sql.DB
	Log   *zap.Logger
	Tasks chan domain.Task
	Cfg   *config.Config
}

func NewHandler(db *sql.DB, log *zap.Logger, tasks chan domain.Task, cfg *config.Config) *Handler {
	return &Handler{
		DB:    db,
		Log:   log,
		Tasks: tasks,
		Cfg:   cfg,
	}
}

func (h *Handler) HandlerHome(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Amount int64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		fmt.Fprintln(w, err.Error())
	}
	var transaction_id int64
	amount := data.Amount
	for {
		if amount-data.Amount > 100 {
			json.NewEncoder(w).Encode(map[string]any{
				"status": 0,
			})
			return
		}
		status, err := repository.CheckTransaction(h.DB, amount)
		if err != nil {
			h.Log.Error(err.Error())
			amount += 1
			continue
		}
		if status {
			transaction_id, err = repository.CreateTransaction(h.DB, amount)
			if err != nil {
				fmt.Fprintln(w, err.Error())
			}
			break
		}
		amount += 1
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Amount        int64 `json:"amount"`
		TransactionID int64 `json:"transaction_id"`
	}{Amount: amount, TransactionID: transaction_id})

}
