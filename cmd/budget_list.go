package cmd

import (
	"finance-cli-manager/internal/db"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var budgetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all budgets",
	Run: func(cmd *cobra.Command, args []string) {

		database, err := db.ConnectDB()
		if err != nil {
			log.Fatal("Failed to connect to DB:", err)
		}
		defer database.Close()

		rows, err := database.Query("SELECT id, category, Amount FROM budgets ORDER BY category ASC")
		if err != nil {
			log.Fatal("Failed to fetch budgets:", err)
		}
		defer rows.Close()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tCATEGORY\tAMOUNT (€)")

		for rows.Next() {
			var id int
			var category string
			var amount float64

			err := rows.Scan(&id, &category, &amount)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Fprintf(w, "%d\t%s\t%.2f\n", id, category, amount)
		}

		w.Flush()
	},
}

func init() {
	budgetCmd.AddCommand(budgetListCmd)
}
