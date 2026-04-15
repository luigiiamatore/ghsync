package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
		fmt.Println("╭─ GitHub Authentication ───────────────────────╮")
		fmt.Println("  Enter your GitHub personal access token.")
		fmt.Println("  💡 Generate at: https://github.com/settings/tokens")
		fmt.Println("╰────────────────────────────────────────────────╯")
		fmt.Println()

		fmt.Print("  Token: ")
		var token string
		fmt.Scanln(&token)

		token = strings.TrimSpace(token)
		if token == "" {
			fmt.Println("\n  ✗ Token cannot be empty.")
			return
		}

		err := saveToken(token)
		if err != nil {
			fmt.Printf("\n  ✗ Error saving token: %v\n\n", err)
			return
		}

		fmt.Println()
		fmt.Println("╭─ Success ─────────────────────────────────────╮")
		fmt.Println("  ✓ Token saved successfully!")
		fmt.Printf("  📝 Location: ~/.ghsync/config\n")
		fmt.Println("╰────────────────────────────────────────────────╯")
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
