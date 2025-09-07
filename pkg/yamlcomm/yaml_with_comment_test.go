package yamlcomm

import (
	"github.io/uberate/hcli/pkg/config"
	"strings"
	"testing"
)

func TestMarshalWithComments(t *testing.T) {
	t.Run("CliConfig with comments", func(t *testing.T) {
		c := config.CliConfig{
			Templates: []config.TemplateConfig{
				{
					Name:       "test name",
					Categories: []string{"test category"},
					Template:   "test template",
				},
			},
		}

		v, err := MarshalWithComments(c)
		if err != nil {
			t.Fatalf("MarshalWithComments failed: %v", err)
		}

		result := string(v)
		if !strings.Contains(result, "Templates:") {
			t.Errorf("Expected output to contain 'Templates:', got: %s", result)
		}
	})

	t.Run("String slice", func(t *testing.T) {
		strSlice := []string{"item1", "item2", "item3"}
		v, err := MarshalWithComments(strSlice)
		if err != nil {
			t.Fatalf("MarshalWithComments failed: %v", err)
		}

		result := string(v)
		expected := "- item1\n- item2\n- item3\n"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("Slice struct with comments", func(t *testing.T) {
		type SliceStruct struct {
			Items []string `yaml:"items" comment:"This is a list of items.\nEach item is a string."`
		}

		sliceStruct := SliceStruct{
			Items: []string{"item1", "item2"},
		}
		v, err := MarshalWithComments(sliceStruct)
		if err != nil {
			t.Fatalf("MarshalWithComments failed: %v", err)
		}

		result := string(v)
		if !strings.Contains(result, "# This is a list of items.") {
			t.Errorf("Expected comment in output, got: %s", result)
		}
		if !strings.Contains(result, "- item1") {
			t.Errorf("Expected items in output, got: %s", result)
		}
	})

	t.Run("Long comment auto-wrapping", func(t *testing.T) {
		type LongCommentStruct struct {
			LongList []string `yaml:"longList" comment:"This is a very long comment that should be automatically wrapped to multiple lines when it exceeds the maximum line length limit of 100 characters. This is a very long comment that should be automatically wrapped to multiple lines when it exceeds the maximum line length limit of 100 characters."`
		}

		longCommentStruct := LongCommentStruct{
			LongList: []string{"item1", "item2"},
		}
		v, err := MarshalWithComments(longCommentStruct)
		if err != nil {
			t.Fatalf("MarshalWithComments failed: %v", err)
		}

		result := string(v)
		// Should have multiple comment lines due to wrapping
		lines := strings.Split(result, "\n")
		commentLines := 0
		for _, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "#") {
				commentLines++
			}
		}
		if commentLines < 2 {
			t.Errorf("Expected multiple comment lines due to wrapping, got %d lines. Output: %s", commentLines, result)
		}
	})

	t.Run("Empty struct", func(t *testing.T) {
		type EmptyStruct struct {
			Field string `yaml:"field"`
		}
		
		empty := EmptyStruct{}
		v, err := MarshalWithComments(empty)
		if err != nil {
			t.Fatalf("MarshalWithComments failed: %v", err)
		}

		result := string(v)
		expected := "field: \"\" # [string]\n"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("Nested struct", func(t *testing.T) {
		type Inner struct {
			Name string `yaml:"name" comment:"Inner name field"`
		}
		type Outer struct {
			Inner Inner `yaml:"inner" comment:"Nested inner struct"`
		}

		nested := Outer{
			Inner: Inner{Name: "test"},
		}
		v, err := MarshalWithComments(nested)
		if err != nil {
			t.Fatalf("MarshalWithComments failed: %v", err)
		}

		result := string(v)
		if !strings.Contains(result, "# Inner name field") || !strings.Contains(result, "# Nested inner struct") {
			t.Errorf("Expected nested comments, got: %s", result)
		}
	})
}
