package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JscorpTech/paymento/database"
)

type Handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		DB: db,
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
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Amount        int64 `json:"amount"`
		TransactionID int64 `json:"transaction_id"`
	}{Amount: amount, TransactionID: transaction_id})

}
