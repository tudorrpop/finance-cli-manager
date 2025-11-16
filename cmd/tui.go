package cmd

import (
	"finance-cli-manager/internal/tui"
	"log"

	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive TUI",
	Long:  "Start the interactive Terminal User Interface to manage budgets, transactions, and reports.",
	Run: func(cmd *cobra.Command, args []string) {
		err := tui.RunTUI()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
