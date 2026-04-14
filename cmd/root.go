package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ghsync",
	Short: "Sync and manage your GitHub repositories locally",
	Long: `ghsync is a CLI tool that helps you manage your GitHub repositories efficiently.

FEATURES:
  • Authenticate securely with your GitHub token
  • Sync all your GitHub repositories locally in bulk
  • Clone new repositories and pull latest changes automatically
  • Generate detailed reports of sync operations
  • Support for both public and private repositories

COMMANDS:
  • auth   - Save your GitHub authentication token
  • pull   - Sync all repositories locally
  • report - View sync operation results

EXAMPLES:
  ghsync auth           # Save your GitHub token
  ghsync pull           # Sync all repositories
  ghsync report         # View the last sync report
  ghsync report --all   # View all sync reports

For more information on a specific command, use:
  ghsync COMMAND --help`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
