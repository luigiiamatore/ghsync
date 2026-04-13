package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ghsync",
	Short: "CLI tool to sync and archive your GitHub repositories",
	Long: `ghsync is a powerful CLI tool designed to simplify the management of your GitHub repositories.

FEATURES:
  • Sync all your GitHub repositories locally in bulk
  • Archive repositories with a single command
  • Generate detailed reports of sync/archive operations
  • Support for authentication with GitHub tokens
  • Automatic conflict resolution and error handling

USE CASES:
  • Backup all your repositories locally
  • Archive inactive projects and keep your GitHub organized
  • Monitor and track repository synchronization status
  • Generate audit reports for repository management

For more information on specific commands, use:
  ghsync COMMAND --help`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
