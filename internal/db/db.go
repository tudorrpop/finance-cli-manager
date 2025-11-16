package db

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Budget struct {
	ID       int
	Category string
	Amount   float64
}

func AddBudget(db *sql.DB, category string, amount float64) error {
	query := `
		INSERT INTO budgets (category, amount, period, start_date)
		VALUES (?, ?, 'monthly', CURRENT_DATE())
	`
	_, err := db.Exec(query, category, amount)
	return err
}

func UpdateBudget(db *sql.DB, id int, category string, amount float64) error {
	query := `
		UPDATE budgets
		SET category = ?, amount = ?
		WHERE id = ?
	`
	res, err := db.Exec(query, category, amount, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("no budget found with ID %d", id)
	}
	return nil
}

func DeleteBudget(db *sql.DB, id int) error {
	query := `
		DELETE FROM budgets
		WHERE id = ?
	`
	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("no budget found with ID %d", id)
	}
	return nil
}

func ListBudgets(db *sql.DB) ([]Budget, error) {
	query := `
		SELECT id, category, amount
		FROM budgets
		ORDER BY category ASC
	`
	rows, err := db.Query(query)
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

func ProcessCSVTransactions(db *sql.DB, filePath string) (processed int, updated map[string]float64, skipped []string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return 0, nil, nil, err
	}

	updated = make(map[string]float64)
	skippedMap := make(map[string]bool)

	for _, rec := range records {
		if len(rec) < 4 {
			continue
		}
		category := rec[2]
		amountStr := rec[3]
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			continue
		}
		if amount < 0 {
			amount = -amount
		}

		var id int
		err = db.QueryRow("SELECT id FROM budgets WHERE category = ?", category).Scan(&id)
		if err == sql.ErrNoRows {
			skippedMap[category] = true
			continue
		} else if err != nil {
			return processed, nil, nil, err
		}

		_, err = db.Exec("UPDATE budgets SET amount = amount - ? WHERE id = ?", amount, id)
		if err != nil {
			return processed, nil, nil, err
		}

		updated[category] += amount
		processed++
	}

	for cat := range skippedMap {
		skipped = append(skipped, cat)
	}

	return processed, updated, skipped, nil
}
