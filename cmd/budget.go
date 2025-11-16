package cmd

import (
	"github.com/spf13/cobra"
)

var budgetCmd = &cobra.Command{
	Use:   "budget",
	Short: "Manage budgets",
	Long:  "Add, list, and update budgets.",
}

func init() {
	rootCmd.AddCommand(budgetCmd)
}
