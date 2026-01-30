package db

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// -- DATA STRUCTURES --

type Budget struct {
	ID       int
	Category string
	Amount   float64
}

type Transaction struct {
	ID       int
	Date     string
	Payee    string
	Category string
	Amount   float64
}

type BudgetReport struct {
	Category string
	Limit    float64
	Spent    float64
}

// -- CONNECTION & SETUP --

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./finance.db")
	if err != nil {
		return nil, err
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS budgets (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			category TEXT UNIQUE, 
			amount REAL
		)`,
		`CREATE TABLE IF NOT EXISTS transactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			date TEXT, 
			payee TEXT, 
			category TEXT, 
			amount REAL
		)`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return nil, fmt.Errorf("init db error: %v", err)
		}
	}

	return db, nil
}

// -- BUDGET FUNCTIONS --

func AddBudget(db *sql.DB, category string, amount float64) error {
	_, err := db.Exec("INSERT INTO budgets (category, amount) VALUES (?, ?)", category, amount)
	return err
}

func UpdateBudget(db *sql.DB, id int, category string, amount float64) error {
	_, err := db.Exec("UPDATE budgets SET category = ?, amount = ? WHERE id = ?", category, amount, id)
	return err
}

func DeleteBudget(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM budgets WHERE id = ?", id)
	return err
}

func ListBudgets(db *sql.DB) ([]Budget, error) {
	rows, err := db.Query("SELECT id, category, amount FROM budgets ORDER BY category ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []Budget
	for rows.Next() {
		var b Budget
		if err := rows.Scan(&b.ID, &b.Category, &b.Amount); err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}
	return budgets, nil
}

// -- TRANSACTION FUNCTIONS --

func ListTransactions(db *sql.DB, search string) ([]Transaction, error) {
	query := "SELECT id, date, payee, category, amount FROM transactions"
	args := []interface{}{}

	if search != "" {
		query += " WHERE payee LIKE ? OR category LIKE ?"
		args = append(args, "%"+search+"%", "%"+search+"%")
	}
	query += " ORDER BY id DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trans []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.Date, &t.Payee, &t.Category, &t.Amount); err != nil {
			return nil, err
		}
		trans = append(trans, t)
	}
	return trans, nil
}

func ProcessCSVTransactions(db *sql.DB, filePath string) (int, []string, float64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, nil, 0, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return 0, nil, 0, err
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, nil, 0, err
	}

	stmt, err := tx.Prepare("INSERT INTO transactions(date, payee, category, amount) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 0, nil, 0, err
	}
	defer stmt.Close()

	count := 0

	for i, row := range records {
		if i == 0 {
			continue
		}

		if len(row) < 4 {
			continue
		}

		date := row[0]
		payee := row[1]
		category := row[2]
		amountStr := row[3]

		if category == "" {
			category = "Uncategorized"
		}

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			continue
		}

		_, err = stmt.Exec(date, payee, category, amount)
		if err == nil {
			count++
		}
	}

	tx.Commit()

	return count, []string{}, 0, nil
}

// -- REPORT FUNCTIONS --

func GetBudgetReport(db *sql.DB) ([]BudgetReport, error) {
	query := `
	SELECT 
		b.category, 
		b.amount as 'limit', 
		COALESCE(SUM(t.amount), 0) as 'spent'
	FROM budgets b
	LEFT JOIN transactions t ON b.category = t.category
	GROUP BY b.id
	ORDER BY spent DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []BudgetReport
	for rows.Next() {
		var r BudgetReport
		if err := rows.Scan(&r.Category, &r.Limit, &r.Spent); err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}
	return reports, nil
}

func AddTransaction(db *sql.DB, date, payee, category string, amount float64) error {
	query := `INSERT INTO transactions (date, payee, category, amount) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, date, payee, category, amount)
	return err
}

func DeleteTransaction(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM transactions WHERE id = ?", id)
	return err
}
