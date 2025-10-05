package template

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"strings"
	"time"
)

type RenderOption struct {
	Title            string
	AppendTags       []string
	AppendCategories []string
	CustomArgs       map[string]string
}

// RenderTemplate renders a template with the given data
func RenderTemplate(ctx context.Context, tmpl *Template, option RenderOption) (string, error) {

	if tmpl == nil {
		return "", errors.New("template is nil")
	}

	vars := map[string]string{
		"title":      option.Title,
		"createAt":   time.Now().Format(time.RFC3339),
		"tags":       fmt.Sprintf("[%s]", strings.Join(append(tmpl.Tags, option.AppendTags...), ",")),
		"categories": fmt.Sprintf("[%s]", strings.Join(append(tmpl.Categories, option.AppendCategories...), ",")),
	}

	for k, v := range option.CustomArgs {
		vars[k] = v
	}

	// Create Go template
	t, err := template.New(tmpl.Name).Parse(tmpl.Template)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template
	var buf bytes.Buffer
	if err := t.Execute(&buf, vars); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// RenderToFile renders a template and writes it to a file using the template's configuration
// Automatically handles directory creation and file path generation
func RenderToFile(ctx context.Context, tmpl *Template, fileName string, data RenderOption) error {
	if tmpl.Dir == "" {
		return fmt.Errorf("template Dir field is empty")
	}

	// Render template content
	content, err := RenderTemplate(ctx, tmpl, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return tmpl.WriteIfNotExists(fileName, []byte(content))
}
