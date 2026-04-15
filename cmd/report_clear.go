package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all sync reports",
	Long: `Delete all stored sync reports from the local machine.

This command will permanently remove all report files located in ~/.ghsync/reports.
Use with caution, as this action cannot be undone.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Permanently delete all reports? (y/n): ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			fmt.Println("Aborting report clear.")
			return
		}

		err := clearReports()
		if err != nil {
			fmt.Printf("Error clearing reports: %v\n", err)
			return
		}

		fmt.Println("All sync reports cleared successfully.")
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
