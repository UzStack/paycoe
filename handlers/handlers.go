package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JscorpTech/paymento/database"
	"github.com/JscorpTech/paymento/workers"
	"go.uber.org/zap"
)

type Handler struct {
	DB    *sql.DB
	Log   *zap.Logger
	Tasks chan workers.Task
}

func NewHandler(db *sql.DB, log *zap.Logger, tasks chan workers.Task) *Handler {
	return &Handler{
		DB:    db,
		Log:   log,
		Tasks: tasks,
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
		status, err := database.CheckTransaction(h.DB, amount)
		if err != nil {
			amount += 1
			continue
		}
		if status {
			transaction_id, err = database.CreateTransaction(h.DB, amount)
			if err != nil {
				fmt.Fprintln(w, err.Error())
			}
			break
		}
		amount += 1
	}

	h.Tasks <- workers.WebhookTask{
		Url:     "https://example.com",
		OrderID: 121,
		Amount:  1212,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Amount        int64 `json:"amount"`
		TransactionID int64 `json:"transaction_id"`
	}{Amount: amount, TransactionID: transaction_id})

}
