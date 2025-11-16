package cmd

import (
	"finance-cli-manager/internal/db"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	updateCategory    string
	updateNewCategory string
	updateAmount      float64
)

var budgetUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing budget",
	Run: func(cmd *cobra.Command, args []string) {

		if updateCategory == "" {
			log.Fatal("Error: --category is required")
		}

		if updateNewCategory == "" && updateAmount == 0 {
			log.Fatal("Error: at least one of --new-category or --amount must be provided")
		}

		database, err := db.ConnectDB()
		if err != nil {
			log.Fatal("Failed to connect to DB:", err)
		}
		defer database.Close()

		query := "UPDATE budgets SET "
		params := []interface{}{}

		if updateNewCategory != "" {
			query += "category = ?, "
			params = append(params, updateNewCategory)
		}

		if updateAmount > 0 {
			query += "amount = ?, "
			params = append(params, updateAmount)
		}

		query = query[:len(query)-2]

		query += " WHERE category = ?"
		params = append(params, updateCategory)

		result, err := database.Exec(query, params...)
		if err != nil {
			log.Fatal("Update failed:", err)
		}

		rows, _ := result.RowsAffected()
		if rows == 0 {
			fmt.Println("No budget found with category:", updateCategory)
			return
		}

		fmt.Println("Budget was successfully updated!")
	},
}

func init() {
	budgetCmd.AddCommand(budgetUpdateCmd)

	budgetUpdateCmd.Flags().StringVar(&updateCategory, "category", "", "Existing category to update (required)")
	budgetUpdateCmd.Flags().StringVar(&updateNewCategory, "new-category", "", "New category name")
	budgetUpdateCmd.Flags().Float64Var(&updateAmount, "amount", 0, "New budget amount")
}
