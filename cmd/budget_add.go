package cmd

import (
	"fmt"
	"log"
	"time"

	"finance-cli-manager/internal/db"

	"github.com/spf13/cobra"
)

var (
	category string
	amount   float64
	period   string
	start    string
)

var budgetAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new budget",
	Run: func(cmd *cobra.Command, args []string) {

		startDate, err := time.Parse("2006-01-02", start)
		if err != nil {
			log.Fatalf("Invalid date format (use YYYY-MM-DD): %v", err)
		}

		database, err := db.ConnectDB()
		if err != nil {
			log.Fatal("Failed to connect to DB:", err)
		}
		defer database.Close()

		query := `
            INSERT INTO budgets (category, amount, period, start_date)
            VALUES (?, ?, ?, ?)
        `
		_, err = database.Exec(query, category, amount, period, startDate)
		if err != nil {
			log.Fatal("Failed to insert budget:", err)
		}

		fmt.Println("Budget added successfully!")
	},
}

func init() {
	budgetCmd.AddCommand(budgetAddCmd)

	budgetAddCmd.Flags().StringVar(&category, "category", "", "Category name")
	budgetAddCmd.Flags().Float64Var(&amount, "amount", 0, "Budget amount")
	budgetAddCmd.Flags().StringVar(&period, "period", "monthly", "Budget period (monthly, weekly, yearly)")
	budgetAddCmd.Flags().StringVar(&start, "start", "2025-01-01", "Start date (YYYY-MM-DD)")

	budgetAddCmd.MarkFlagRequired("category")
	budgetAddCmd.MarkFlagRequired("amount")
}
