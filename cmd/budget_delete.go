package cmd

import (
	"finance-cli-manager/internal/db"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	deleteCategory string
)

var budgetDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a budget by category name",
	Run: func(cmd *cobra.Command, args []string) {

		if deleteCategory == "" {
			log.Fatal("You must provide a category name using --category")
		}

		database, err := db.ConnectDB()
		if err != nil {
			log.Fatal("Failed to connect to DB:", err)
		}
		defer database.Close()

		result, err := database.Exec("DELETE FROM budgets WHERE category = ?", deleteCategory)
		if err != nil {
			log.Fatal("Failed to delete budget:", err)
		}

		rows, _ := result.RowsAffected()
		if rows == 0 {
			fmt.Println("No budget found with that category.")
			return
		}

		fmt.Println("Budget deleted successfully!")
	},
}

func init() {
	budgetCmd.AddCommand(budgetDeleteCmd)
	budgetDeleteCmd.Flags().StringVar(&deleteCategory, "category", "", "Category of the budget to delete")
	budgetDeleteCmd.MarkFlagRequired("category")
}
