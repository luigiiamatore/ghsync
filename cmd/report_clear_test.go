package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestClearReports(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()

	// Mock home directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	// Create .ghsync/reports directory with test files
	reportsDir := filepath.Join(tmpDir, ".ghsync", "reports")
	err := os.MkdirAll(reportsDir, 0700)
	if err != nil {
		t.Fatalf("Failed to create reports directory: %v", err)
	}

	// Create test report files
	testFiles := []string{"report1.json", "report2.json", "report3.json"}
	for _, filename := range testFiles {
		filepath := filepath.Join(reportsDir, filename)
		err := os.WriteFile(filepath, []byte("dummy report"), 0600)
		if err != nil {
			t.Fatalf("Failed to create test report file: %v", err)
		}
	}

	// Verify files exist before clearing
	entries, err := os.ReadDir(reportsDir)
	if err != nil {
		t.Fatalf("Failed to read reports directory: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("Expected 3 files before clear, got %d", len(entries))
	}

	// Clear reports
	err = clearReports()
	if err != nil {
		t.Fatalf("clearReports failed: %v", err)
	}

	// Verify directory was removed
	if _, err := os.Stat(reportsDir); err == nil {
		t.Error("Reports directory should have been removed")
	}

	t.Logf("✓ Clear reports test passed: successfully removed reports directory")
}

func TestClearReportsNoDirectory(t *testing.T) {
	// Create temporary directory structure without reports dir
	tmpDir := t.TempDir()

	// Mock home directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	// Try to clear when reports directory doesn't exist
	// This should not error (RemoveAll succeeds on non-existent paths)
	err := clearReports()
	if err != nil {
		t.Errorf("clearReports should handle non-existent directory, got error: %v", err)
	}

	t.Logf("✓ Clear reports test passed: handles non-existent directory")
}
