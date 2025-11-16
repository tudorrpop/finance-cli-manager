package db

import (
	"database/sql"
	"fmt"
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
