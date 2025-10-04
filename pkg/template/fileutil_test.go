package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureDir(t *testing.T) {
	testDir := "testdata/subdir/nested"

	// Clean up after test
	defer os.RemoveAll("testdata")

	err := EnsureDir(testDir)
	if err != nil {
		t.Fatalf("EnsureDir failed: %v", err)
	}

	// Verify directory exists
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Fatalf("Directory was not created: %s", testDir)
	}

	// Test with file path
	filePath := "testdata/files/output.txt"
	err = EnsureDir(filePath)
	if err != nil {
		t.Fatalf("EnsureDir with file path failed: %v", err)
	}

	// Should create the directory portion
	if _, err := os.Stat(filepath.Dir(filePath)); os.IsNotExist(err) {
		t.Fatalf("Parent directory was not created for file path: %s", filePath)
	}
}

func TestSafeWriteFile(t *testing.T) {
	testFile := "testdata/output/written_file.txt"
	testContent := "Hello, World!"

	// Clean up after test
	defer os.RemoveAll("testdata")

	err := SafeWriteString(testFile, testContent)
	if err != nil {
		t.Fatalf("SafeWriteString failed: %v", err)
	}

	// Verify file exists and content is correct
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}

	if string(content) != testContent {
		t.Fatalf("File content mismatch. Expected: %s, Got: %s", testContent, string(content))
	}
}

func TestFileExists(t *testing.T) {
	// Test with non-existent file
	if FileExists("nonexistent_file.txt") {
		t.Error("FileExists should return false for non-existent file")
	}

	// Create a test file
	testFile := "test_existence.txt"
	os.WriteFile(testFile, []byte("test"), 0644)
	defer os.Remove(testFile)

	if !FileExists(testFile) {
		t.Error("FileExists should return true for existing file")
	}
}

func TestIsDir(t *testing.T) {
	// Test with file
	testFile := "test_file.txt"
	os.WriteFile(testFile, []byte("test"), 0644)
	defer os.Remove(testFile)

	if IsDir(testFile) {
		t.Error("IsDir should return false for files")
	}

	// Test with directory
	testDir := "test_dir"
	os.Mkdir(testDir, 0755)
	defer os.Remove(testDir)

	if !IsDir(testDir) {
		t.Error("IsDir should return true for directories")
	}
}