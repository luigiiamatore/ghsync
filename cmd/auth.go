package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/luigiiamatore/ghsync/internal/ui"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Save your GitHub authentication token",
	Long: `Authenticate with GitHub by storing your personal access token locally.

The token will be securely saved in ~/.ghsync/config with restricted permissions,
ensuring only you can read it. This token will be used by other commands for GitHub API access.

REQUIREMENTS:
  • A valid GitHub personal access token with appropriate permissions
  • For syncing repos: 'repo' and 'admin:org_hook' scopes
  • For archiving repos: 'repo' scope

EXAMPLE:
  ghsync auth
  
This will prompt you to enter your token interactively.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println()
		ui.PrintBox("GitHub Authentication",
			"Enter your GitHub personal access token.",
			"💡 Generate at: https://github.com/settings/tokens",
		)
		fmt.Println()

		fmt.Print("  Token: ")
		var token string
		fmt.Scanln(&token)

		token = strings.TrimSpace(token)
		if token == "" {
			fmt.Println()
			ui.PrintWarning("Token cannot be empty.")
			fmt.Println()
			return
		}

		err := saveToken(token)
		if err != nil {
			fmt.Println()
			fmt.Printf("Error saving token: %v\n", err)
			fmt.Println()
			return
		}

		fmt.Println()
		ui.PrintSuccess(
			"✓ Token saved successfully!",
			"📝 Location: ~/.ghsync/config",
		)
		fmt.Println()
	},
}

func saveToken(token string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(home, ".ghsync")
	err = os.MkdirAll(configDir, 0700)
	if err != nil {
		return err
	}

	configFile := filepath.Join(configDir, "config")
	err = os.WriteFile(configFile, []byte(token), 0600)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(authCmd)
}
