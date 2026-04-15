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
		var token string

		fmt.Print("Enter your GitHub token: ")
		fmt.Scanln(&token)

		err := saveToken(strings.TrimSpace(token))
		if err != nil {
			fmt.Println("Error saving token: ", err)
			return
		}

		fmt.Println("Token saved successfully!")
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
