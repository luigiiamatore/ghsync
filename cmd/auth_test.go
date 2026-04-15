package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveToken(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()

	// Mock home directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	testToken := "ghp_test_token_12345"

	// Call saveToken function
	err := saveToken(testToken)
	if err != nil {
		t.Fatalf("saveToken failed: %v", err)
	}

	// Verify file was created
	configFile := filepath.Join(tmpDir, ".ghsync", "config")
	if _, err := os.Stat(configFile); err != nil {
		t.Errorf("Token file not found: %v", err)
	}

	// Verify file content
	content, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Failed to read token file: %v", err)
	}

	if string(content) != testToken {
		t.Errorf("Token mismatch. Expected %s, got %s", testToken, string(content))
	}

	// Verify file permissions are 0600 (read/write only for owner)
	fileInfo, err := os.Stat(configFile)
	if err != nil {
		t.Fatalf("Failed to stat token file: %v", err)
	}

	expectedMode := os.FileMode(0600)
	actualMode := fileInfo.Mode().Perm()
	if actualMode != expectedMode {
		t.Errorf("Token file permissions incorrect. Expected %o, got %o", expectedMode, actualMode)
	}

	t.Logf("✓ Auth test passed: token saved with correct permissions (0600)")
}
