package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/luigiiamatore/ghsync/internal/ui"
	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all sync reports",
	Long: `Delete all stored sync reports from the local machine.

This command will permanently remove all report files located in ~/.ghsync/reports.
Use with caution, as this action cannot be undone.`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			return
		}

		reportsDir := filepath.Join(home, ".ghsync", "reports")
		entries, err := os.ReadDir(reportsDir)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("\n⚠ No reports found to clear.")
				return
			}
			fmt.Printf("Error reading reports: %v\n", err)
			return
		}

		if len(entries) == 0 {
			fmt.Println("\n⚠ No reports found to clear.")
			return
		}

		// Show warning
		fmt.Println()
		ui.PrintWarning(
			fmt.Sprintf("You are about to delete %d report(s).", len(entries)),
			"This action cannot be undone!",
		)
		fmt.Println()

		fmt.Printf("  Permanently delete all reports? (y/n): ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			fmt.Println("  ✓ Clear operation cancelled.")
			return
		}

		err = clearReports()
		if err != nil {
			fmt.Printf("Error clearing reports: %v\n", err)
			return
		}

		fmt.Println()
		ui.PrintSuccess(fmt.Sprintf("Cleared %d report(s).", len(entries)))
		fmt.Println()
	},
}

func clearReports() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	reportsDir := filepath.Join(home, ".ghsync", "reports")
	err = os.RemoveAll(reportsDir)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	reportCmd.AddCommand(clearCmd)
}
