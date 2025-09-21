package repository

import (
	"database/sql"
	"github.com/google/uuid"
)

var (
	TransactionTimeOutMinutes = "30"
)

func InitTables(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS transactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			amount INTEGER NOT NULL,
			status BOOLEAN DEFAULT 1,
			webhook_status BOOLEAN DEFAULT 0,
			transaction_id TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS amount_created_at_index ON transactions(amount, created_at);
		CREATE UNIQUE INDEX IF NOT EXISTS transaction_id_index ON transactions(transaction_id);
	`)
	if err != nil {
		panic(err)
	}
}

func CreateTransaction(db *sql.DB, amount int64) (string, error) {
	transaction_id := uuid.New().String()
	_, err := db.Exec("INSERT INTO transactions(amount,transaction_id) VALUES(?, ?)", amount, transaction_id)
	if err != nil {
		return "", err
	}
	return transaction_id, nil
}

func CheckTransaction(db *sql.DB, amount int64) (bool, error) {
	var count int
	err := db.QueryRow(`SELECT count(id) FROM transactions WHERE 
		amount=? and status=1 and 
		created_at BETWEEN datetime('now', '-`+TransactionTimeOutMinutes+` minutes') and 
		datetime('now')`, amount).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func DeleteTransaction(db *sql.DB, transaction_id string) error {
	_, err := db.Exec("DELETE FROM transactions WHERE transaction_id=?", transaction_id)
	if err != nil {
		return err
	}
	return nil
}

func GetTransaction(db *sql.DB, amount int64) (string, error) {
	var trans_id string
	err := db.QueryRow(`SELECT transaction_id FROM transactions WHERE 
		amount=? and status=1 and
		created_at BETWEEN datetime('now', '-`+TransactionTimeOutMinutes+` minutes') and datetime('now')`, amount).Scan(&trans_id)
	if err != nil {
		return "", err
	}
	return trans_id, nil
}

func GetOldTransactions(db *sql.DB) ([]map[string]any, error) {
	result, err := db.Query(`SELECT transaction_id, amount FROM transactions WHERE status=1 and created_at <= datetime('now', '-` + TransactionTimeOutMinutes + ` minutes')`)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	var transactions []map[string]any
	for result.Next() {
		var transaction_id string
		var amount int64
		if err := result.Scan(&transaction_id, &amount); err != nil {
			return nil, err
		}
		transactions = append(transactions, map[string]any{
			"amount":         amount,
			"transaction_id": transaction_id,
		})
	}
	return transactions, nil
}

func ConfirmTransaction(db *sql.DB, transactionID string, webhookStatus bool) error {
	_, err := db.Exec("UPDATE transactions SET status=0, webhook_status=? WHERE transaction_id=?", webhookStatus, transactionID)
	if err != nil {
		return err
	}
	return nil
}
