package database

import "database/sql"

func InitTables(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS transactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			amount INTEGER NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE UNIQUE INDEX IF NOT EXISTS amount_index ON transactions(amount);
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
	result, err := db.Query("SELECT count(id) FROM transactions WHERE amount=?", amount)
	var count int
	defer result.Close()
	if err != nil {
		return false, err
	}
	if result.Next() {
		err = result.Scan(&count)
		if err != nil {
			return false, err
		}
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
