package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/luigiiamatore/ghsync/internal/report"
	"github.com/luigiiamatore/ghsync/internal/ui"

	"github.com/google/go-github/v60/github"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Sync all your GitHub repositories locally",
	Long: `Sync all your GitHub repositories to your local machine. This command will:
  • Clone any repositories that don't exist locally
  • Pull the latest changes for repositories that already exist
  • Handle all your repositories in bulk with a single command

The command requires authentication to be set up first using 'ghsync auth'.

OPTIONS:
  • --dir: Specify the directory where repositories will be stored (default: ghsync-repos)

EXAMPLE:
  ghsync pull
  ghsync pull --dir /path/to/repos

STATUS:
  ✓ Clones new repositories
  ✓ Updates existing repositories
  ✓ Preserves local changes via git pull`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			fmt.Println("Error getting dir flag: ", err)
			return
		}

		token, err := readToken()
		if err != nil {
			fmt.Println("Error reading token: ", err)
			return
		}

		client := authenticate(token)

		var allRepos []*github.Repository
		opts := &github.RepositoryListByAuthenticatedUserOptions{
			Type: "owner",
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		}

		for {
			repos, resp, err := client.Repositories.ListByAuthenticatedUser(context.Background(), opts)
			if err != nil {
				fmt.Println("Error listing repositories: ", err)
				return
			}

			allRepos = append(allRepos, repos...)

			if resp.NextPage == 0 {
				break
			}
			opts.ListOptions.Page = resp.NextPage
		}

		repos := allRepos

		syncReport := report.SyncReport{
			Timestamp:    time.Now(),
			TotalRepos:   len(repos),
			SyncedRepos:  0,
			ClonedRepos:  0,
			UpdatedRepos: 0,
			Errors:       []report.SyncError{},
		}

		updatedCount := 0
		clonedCount := 0
		var errorMessages []string

		// Header
		fmt.Println()
		ui.PrintBox("Repository Sync",
			fmt.Sprintf("Syncing %d repositories to %s", len(repos), dir),
		)
		fmt.Println()

		bar := progressbar.NewOptions(len(repos),
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowCount(),
			progressbar.OptionFullWidth(),
		)

		for _, repo := range repos {
			path := filepath.Join(dir, repo.GetName())

			var repoErr error

			if _, err := os.Stat(path); err == nil {
				repoErr = exec.Command("git", "-C", path, "pull").Run()
				if repoErr != nil {
					syncReport.Errors = append(syncReport.Errors, report.SyncError{
						RepoName: repo.GetName(),
						ErrorMsg: repoErr.Error(),
					})
					errorMessages = append(errorMessages, fmt.Sprintf("✗ %s: %s", repo.GetName(), repoErr.Error()))
				} else {
					updatedCount++
				}
			} else {
				authenticatedURL := buildAuthenticatedCloneURL(repo.GetCloneURL(), token)
				repoErr = exec.Command("git", "clone", authenticatedURL, path).Run()
				if repoErr != nil {
					syncReport.Errors = append(syncReport.Errors, report.SyncError{
						RepoName: repo.GetName(),
						ErrorMsg: repoErr.Error(),
					})
					errorMessages = append(errorMessages, fmt.Sprintf("✗ %s: %s", repo.GetName(), repoErr.Error()))
				} else {
					clonedCount++
				}
			}

			bar.Add(1)
		}

		fmt.Println()

		syncReport.SyncedRepos = updatedCount + clonedCount
		syncReport.UpdatedRepos = updatedCount
		syncReport.ClonedRepos = clonedCount

		err = report.SaveSyncReport(syncReport)
		if err != nil {
			fmt.Println("Error saving sync report: ", err)
		}

		// Summary
		fmt.Println()
		summaryLines := []string{
			fmt.Sprintf("✓ Synced:  %d", syncReport.SyncedRepos),
			fmt.Sprintf("⬇ Cloned:  %d", syncReport.ClonedRepos),
			fmt.Sprintf("⬆ Updated: %d", syncReport.UpdatedRepos),
		}
		if len(errorMessages) > 0 {
			summaryLines = append(summaryLines, fmt.Sprintf("✗ Errors:  %d", len(errorMessages)))
		}
		ui.PrintBox("Sync Summary", summaryLines...)

		// Errors if any
		if len(errorMessages) > 0 {
			fmt.Println()
			ui.PrintErrors(errorMessages...)
		}

		fmt.Println()
		fmt.Printf("📄 Report: ~/.ghsync/reports/%s.json\n", syncReport.Timestamp.Format("2006-01-02T15-04-05"))
		fmt.Println()
	},
}

func readToken() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configFile := fmt.Sprintf("%s/.ghsync/config", home)
	tokenBytes, err := os.ReadFile(configFile)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(tokenBytes)), nil
}

func authenticate(token string) *github.Client {
	client := github.NewClient(nil).WithAuthToken(token)
	return client
}

func buildAuthenticatedCloneURL(cloneURL, token string) string {
	if strings.HasPrefix(cloneURL, "https://") {
		return "https://" + token + ":x-oauth-basic@github.com/" + strings.TrimPrefix(cloneURL, "https://github.com/")
	}
	return cloneURL
}

func init() {
	pullCmd.Flags().String("dir", "ghsync-repos", "Directory to store the synced repositories")
	rootCmd.AddCommand(pullCmd)
}
