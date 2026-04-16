package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/luigiiamatore/ghsync/internal/report"
	"github.com/luigiiamatore/ghsync/internal/ui"
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
			ui.PrintBox("All Sync Reports",
				fmt.Sprintf("Found %d report(s)", len(entries)),
			)
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
			fmt.Println()
			ui.PrintWarning("No reports found. Run 'ghsync pull' to generate a report.")
			fmt.Println()
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
	reportLines := []string{
		fmt.Sprintf("📅 Date:    %s", syncReport.Timestamp.Format("2006-01-02 15:04:05")),
		"─────────────────────────────────────────────",
		fmt.Sprintf("📊 Total:   %d", syncReport.TotalRepos),
		fmt.Sprintf("✓ Synced:   %d", syncReport.SyncedRepos),
		fmt.Sprintf("⬇ Cloned:   %d", syncReport.ClonedRepos),
		fmt.Sprintf("⬆ Updated:  %d", syncReport.UpdatedRepos),
	}

	if len(syncReport.Errors) > 0 {
		reportLines = append(reportLines, "─────────────────────────────────────────────")
		reportLines = append(reportLines, fmt.Sprintf("✗ Errors:   %d", len(syncReport.Errors)))
	} else {
		reportLines = append(reportLines, fmt.Sprintf("✓ Errors:   0"))
	}

	ui.PrintBox("Sync Report", reportLines...)

	if len(syncReport.Errors) > 0 {
		fmt.Println()
		errorLines := []string{}
		for _, syncErr := range syncReport.Errors {
			errorLines = append(errorLines,
				fmt.Sprintf("✗ %s", syncErr.RepoName),
				fmt.Sprintf("  → %s", syncErr.ErrorMsg),
			)
		}
		ui.PrintBox("Failed Repositories", errorLines...)
	}
}

func init() {
	reportCmd.Flags().Bool("all", false, "Indicates whether to show all reports instead of just the last one")
	rootCmd.AddCommand(reportCmd)
}
