package template

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// EnsureDir ensures that the directory exists, creating it if necessary
func EnsureDir(path string) error {
	// Clean the path to handle any relative paths
	cleanPath := filepath.Clean(path)

	// If it's a file path, get the directory
	if strings.Contains(cleanPath, ".") && !strings.HasSuffix(cleanPath, string(os.PathSeparator)) {
		cleanPath = filepath.Dir(cleanPath)
	}

	// Check if directory already exists
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		// Create all directories in the path
		if err := os.MkdirAll(cleanPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", cleanPath, err)
		}
		fmt.Printf("Created directory: %s\n", cleanPath)
	} else if err != nil {
		return fmt.Errorf("failed to check directory %s: %w", cleanPath, err)
	}

	return nil
}

// SafeWriteFile writes data to a file, ensuring the directory path exists
func SafeWriteFile(path string, data []byte) error {
	// Ensure the directory exists
	if err := EnsureDir(path); err != nil {
		return fmt.Errorf("failed to ensure directory for %s: %w", path, err)
	}

	// Write the file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	fmt.Printf("File written: %s\n", path)
	return nil
}

// SafeWriteString writes a string to a file, ensuring the directory path exists
func SafeWriteString(path, content string) error {
	return SafeWriteFile(path, []byte(content))
}

// CopyFile copies a file from source to destination, ensuring the destination directory exists
func CopyFile(src, dst string) error {
	// Ensure destination directory exists
	if err := EnsureDir(dst); err != nil {
		return fmt.Errorf("failed to ensure destination directory: %w", err)
	}

	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer destFile.Close()

	// Copy content
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	fmt.Printf("File copied: %s -> %s\n", src, dst)
	return nil
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsDir checks if the path is a directory
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// GetAbsolutePath returns the absolute path, resolving any relative paths
func GetAbsolutePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for %s: %w", path, err)
	}
	return absPath, nil
}
