package repository

import (
	"database/sql"
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
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE UNIQUE INDEX IF NOT EXISTS amount_created_at_index ON transactions(amount, created_at);
	`)
	if err != nil {
		panic(err)
	}
}

func CreateTransaction(db *sql.DB, amount int64) (int64, error) {
	result, err := db.Exec("INSERT INTO transactions(amount) VALUES(?)", amount)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
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

func DeleteTransaction(db *sql.DB, transaction_id int64) error {
	_, err := db.Exec("DELETE FROM transactions WHERE id=?", transaction_id)
	if err != nil {
		return err
	}
	return nil
}

func GetTransaction(db *sql.DB, amount int64) (int64, error) {
	var trans_id int64
	err := db.QueryRow(`SELECT id FROM transactions WHERE 
		amount=? and status=1 and
		created_at BETWEEN datetime('now', '-`+TransactionTimeOutMinutes+` minutes') and datetime('now')`, amount).Scan(&trans_id)
	if err != nil {
		return 0, err
	}
	return trans_id, nil
}

func ConfirmTransaction(db *sql.DB, transaction_id int64) error {
	_, err := db.Exec("UPDATE transactions SET status=0 WHERE id=?", transaction_id)
	if err != nil {
		return err
	}
	return nil
}
