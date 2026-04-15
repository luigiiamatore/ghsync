package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/luigiiamatore/ghsync/internal/report"
)

func TestReadSyncReport(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()

	// Create a test sync report
	testReport := &report.SyncReport{
		Timestamp:    time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		TotalRepos:   10,
		SyncedRepos:  8,
		ClonedRepos:  3,
		UpdatedRepos: 5,
		Errors: []report.SyncError{
			{
				RepoName: "test-repo",
				ErrorMsg: "clone failed",
			},
		},
	}

	// Marshal to JSON and write to temporary file
	jsonData, err := json.Marshal(testReport)
	if err != nil {
		t.Fatalf("Failed to marshal report: %v", err)
	}

	reportPath := filepath.Join(tmpDir, "test-report.json")
	err = os.WriteFile(reportPath, jsonData, 0600)
	if err != nil {
		t.Fatalf("Failed to write report file: %v", err)
	}

	// Read the report using readSyncReport
	readReport, err := readSyncReport(reportPath)
	if err != nil {
		t.Fatalf("readSyncReport failed: %v", err)
	}

	// Verify report contents
	if readReport.TotalRepos != 10 {
		t.Errorf("Expected 10 total repos, got %d", readReport.TotalRepos)
	}

	if readReport.SyncedRepos != 8 {
		t.Errorf("Expected 8 synced repos, got %d", readReport.SyncedRepos)
	}

	if readReport.ClonedRepos != 3 {
		t.Errorf("Expected 3 cloned repos, got %d", readReport.ClonedRepos)
	}

	if readReport.UpdatedRepos != 5 {
		t.Errorf("Expected 5 updated repos, got %d", readReport.UpdatedRepos)
	}

	if len(readReport.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(readReport.Errors))
	}

	if readReport.Errors[0].RepoName != "test-repo" {
		t.Errorf("Expected error repo name 'test-repo', got %s", readReport.Errors[0].RepoName)
	}

	t.Logf("✓ Report test passed: successfully read and verified report")
}

func TestReadSyncReportFileNotFound(t *testing.T) {
	// Try to read non-existent file
	_, err := readSyncReport("/nonexistent/path/report.json")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}

	t.Logf("✓ Report test passed: correctly handles missing file")
}
