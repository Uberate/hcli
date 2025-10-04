package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRenderTemplate(t *testing.T) {
	tmpl := &Template{
		Name:       "test",
		Template:   "Hello, {{.Data.Name}}!",
		Categories: []string{"test"},
		Tags:       []string{"demo"},
		Dir:        "test_output",
	}

	data := map[string]interface{}{
		"Name": "World",
	}

	// Test rendering
	result, err := RenderTemplate(tmpl, data)
	if err != nil {
		t.Fatalf("RenderTemplate failed: %v", err)
	}

	expected := "Hello, World!"
	if result != expected {
		t.Fatalf("Expected: %s, Got: %s", expected, result)
	}
}

func TestRenderToFile(t *testing.T) {
	tmpl := &Template{
		Name:       "test_file",
		Template:   "Title: {{.Data.Title}}\nContent: {{.Data.Content}}",
		Categories: []string{"test"},
		Tags:       []string{"file"},
		Dir:        "test_output",
		NeedDir:    false,
	}

	data := map[string]interface{}{
		"Title":   "Test Title",
		"Content": "This is test content",
	}

	// Clean up
	defer os.RemoveAll("test_output")

	// Test render to file
	err := RenderToFile(tmpl, data)
	if err != nil {
		t.Fatalf("RenderToFile failed: %v", err)
	}

	// Check that a file was created with timestamp in filename
	files, err := filepath.Glob("test_output/test_file_*.md")
	if err != nil {
		t.Fatalf("Failed to list files: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("No files were created in the output directory")
	}

	// Verify the file content
	content, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatalf("Failed to read created file: %v", err)
	}

	expectedContent := "Title: Test Title\nContent: This is test content"
	if string(content) != expectedContent {
		t.Fatalf("File content mismatch. Expected: %s, Got: %s", expectedContent, string(content))
	}
}

func TestRenderToFileWithNeedDir(t *testing.T) {
	tmpl := &Template{
		Name:       "test_need_dir",
		Template:   "Data: {{.Data.Value}}",
		Categories: []string{"test"},
		Tags:       []string{"dir"},
		Dir:        "test_output_need_dir",
		NeedDir:    true,
	}

	data := map[string]interface{}{
		"Value": "test_value",
	}

	// Clean up
	defer os.RemoveAll("test_output_need_dir")

	// Test render to file with NeedDir=true
	err := RenderToFile(tmpl, data)
	if err != nil {
		t.Fatalf("RenderToFile with NeedDir failed: %v", err)
	}

	// Should create index.md in the directory
	expectedPath := "test_output_need_dir/index.md"
	if !FileExists(expectedPath) {
		t.Fatalf("index.md was not created: %s", expectedPath)
	}

	// Verify the file content
	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("Failed to read index.md: %v", err)
	}

	expectedContent := "Data: test_value"
	if string(content) != expectedContent {
		t.Fatalf("File content mismatch. Expected: %s, Got: %s", expectedContent, string(content))
	}
}