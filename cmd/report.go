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

func printSyncReport(syncReport *report.SyncReport, indent string) {
	fmt.Printf("%sReport - %s\n", indent, syncReport.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("%s  Total: %d\n", indent, syncReport.TotalRepos)
	fmt.Printf("%s  Synced: %d\n", indent, syncReport.SyncedRepos)
	fmt.Printf("%s  Cloned: %d\n", indent, syncReport.ClonedRepos)
	fmt.Printf("%s  Updated: %d\n", indent, syncReport.UpdatedRepos)

	if len(syncReport.Errors) > 0 {
		fmt.Printf("%s\n%s  Errors (%d):\n", indent, indent, len(syncReport.Errors))
		for _, syncErr := range syncReport.Errors {
			fmt.Printf("%s    ✗ %s\t%s\n", indent, syncErr.RepoName, syncErr.ErrorMsg)
		}
	} else {
		fmt.Printf("%s  No errors reported.\n", indent)
	}
}

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
			printSyncReport(syncReport, "")
			fmt.Println()
		} else {
			fmt.Printf("\nAll sync reports:\n")
			for _, entry := range entries {
				reportPath := filepath.Join(reportDir, entry.Name())
				syncReport, err := readSyncReport(reportPath)
				if err != nil {
					fmt.Printf("Error reading report file: %v\n", err)
					continue
				}

				printSyncReport(syncReport, "  ")
				fmt.Println()
			}
		}
	},
}

func init() {
	reportCmd.Flags().Bool("all", false, "Indicates whether to show all reports instead of just the last one")
	rootCmd.AddCommand(reportCmd)
}
