package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type SyncError struct {
	RepoName string
	ErrorMsg string
}

type SyncReport struct {
	Timestamp    time.Time
	TotalRepos   int
	SyncedRepos  int
	ClonedRepos  int
	UpdatedRepos int
	Errors       []SyncError
}

func SaveSyncReport(report SyncReport) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	reportDir := filepath.Join(home, ".ghsync", "reports")
	fileName := fmt.Sprintf("%s.json", report.Timestamp.Format("2006-01-02T15-04-05"))

	if _, err := os.Stat(reportDir); os.IsNotExist(err) {
		err = os.Mkdir(reportDir, 0755)
		if err != nil {
			return err
		}
	}

	reportPath := filepath.Join(reportDir, fileName)
	reportData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(reportPath, reportData, 0644)
	if err != nil {
		return err
	}

	return nil
}
