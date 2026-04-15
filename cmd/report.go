package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/luigiiamatore/ghsync/internal/report"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "View sync reports",
	Long: `Display the results of your repository sync operations.

By default, shows the most recent report. Use --all to view all reports.

EXAMPLES:
  ghsync report
  ghsync report --all`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			return
		}

		reportDir := filepath.Join(home, ".ghsync", "reports")
		entries, err := os.ReadDir(reportDir)
		if err != nil {
			fmt.Printf("Error reading reports: %v\n", err)
			return
		}

		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Name() > entries[j].Name()
		})

		showAll, _ := cmd.Flags().GetBool("all")
		if !showAll && len(entries) > 0 {
			reportPath := filepath.Join(reportDir, entries[0].Name())
			syncReport, err := readSyncReport(reportPath)
			if err != nil {
				fmt.Printf("Error reading report file: %v\n", err)
				return
			}

			fmt.Println()
			printSyncReport(syncReport)
			fmt.Println()
		} else if showAll {
			fmt.Println()
			fmt.Println("╭─ All Sync Reports ────────────────────────────╮")
			fmt.Printf("  Found %d report(s)\n", len(entries))
			fmt.Println("╰────────────────────────────────────────────────╯")
			fmt.Println()

			for i, entry := range entries {
				reportPath := filepath.Join(reportDir, entry.Name())
				syncReport, err := readSyncReport(reportPath)
				if err != nil {
					fmt.Printf("Error reading report file: %v\n", err)
					continue
				}

				if i > 0 {
					fmt.Println()
				}
				printSyncReport(syncReport)
			}

			fmt.Println()
		} else {
			fmt.Println("\n⚠ No reports found. Run 'ghsync pull' to generate a report.\n")
		}
	},
}

func readSyncReport(reportPath string) (*report.SyncReport, error) {
	reportData, err := os.ReadFile(reportPath)
	if err != nil {
		return nil, err
	}

	var syncReport report.SyncReport
	err = json.Unmarshal(reportData, &syncReport)
	if err != nil {
		return nil, err
	}

	return &syncReport, nil
}

func printSyncReport(syncReport *report.SyncReport) {
	fmt.Println("╭─ Sync Report ─────────────────────────────────╮")
	fmt.Printf("  📅 Date:    %s\n", syncReport.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Println("  ─────────────────────────────────────────────")
	fmt.Printf("  📊 Total:   %d\n", syncReport.TotalRepos)
	fmt.Printf("  ✓ Synced:   %d\n", syncReport.SyncedRepos)
	fmt.Printf("  ⬇ Cloned:   %d\n", syncReport.ClonedRepos)
	fmt.Printf("  ⬆ Updated:  %d\n", syncReport.UpdatedRepos)

	if len(syncReport.Errors) > 0 {
		fmt.Println("  ─────────────────────────────────────────────")
		fmt.Printf("  ✗ Errors:   %d\n", len(syncReport.Errors))
		fmt.Println("╰────────────────────────────────────────────────╯")
		fmt.Println()
		fmt.Println("╭─ Failed Repositories ─────────────────────────╮")
		for _, syncErr := range syncReport.Errors {
			fmt.Printf("  ✗ %s\n", syncErr.RepoName)
			fmt.Printf("    → %s\n", syncErr.ErrorMsg)
		}
		fmt.Println("╰────────────────────────────────────────────────╯")
	} else {
		fmt.Printf("  ✓ Errors:   0\n")
		fmt.Println("╰────────────────────────────────────────────────╯")
	}
}

func init() {
	reportCmd.Flags().Bool("all", false, "Indicates whether to show all reports instead of just the last one")
	rootCmd.AddCommand(reportCmd)
}
