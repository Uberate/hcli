package template

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Template struct {
	Name string `yaml:"Name" describe:"Template name, use hcli gen posts --template-name(|-n)=Name to generate a new 
posts by template value to init."`

	Categories []string `yaml:"Categories" describe:"Categories of posts, can be append by '--categories' args."`
	Tags       []string `yaml:"Tags" describe:"Tags of posts, can be append by '--tag|-t'"`

	Template string `yaml:"Template" describe:"The go template of posts."`

	Dir     string `yaml:"Dir" describe:"The directory of posts, absolute of command run path."`
	NeedDir bool   `yaml:"NeedDir" describe:"Whether to need directory, if need, hcli will create posts in a new dir
named args and set file to index.md" default:"false"`

	PicSummaryPrompt string `yaml:"PicSummaryPrompt" describe:"Pic summary prompt"`
}

func (t Template) GetFilePath(fileName string) string {

	if strings.HasSuffix(fileName, ".md") {
		fileName = strings.TrimSuffix(fileName, ".md")
	}

	res := ""
	if t.NeedDir {
		res = path.Join(fileName, "index.md")
	} else {
		res = fmt.Sprintf("%s.md", fileName)
	}

	return filepath.Join(t.Dir, res)
}

func (t Template) WritePicSummary(fileName string, summary string) error {
	dir := path.Dir(t.GetFilePath(fileName))
	summaryFileName := path.Join(dir, fmt.Sprintf("%s.summary.text", fileName))
	if FileExists(summaryFileName) {
		return fmt.Errorf("file %s already exists", summaryFileName)
	}

	return SafeWriteFile(summaryFileName, []byte(summary))
}

func (t Template) WritePoster(fileName string, picData []byte) error {
	dir := path.Dir(t.GetFilePath(fileName))
	summaryFileName := path.Join(dir, "feature.png")
	if FileExists(summaryFileName) {
		return fmt.Errorf("file %s already exists", summaryFileName)
	}

	return SafeWriteFile(summaryFileName, picData)
}

func (t Template) WriteIfNotExists(fileName string, data []byte) error {
	outputPath := t.GetFilePath(fileName)

	// Check if file already exists and create a new name with timestamp
	if FileExists(outputPath) {
		return fmt.Errorf("file already exists: %s", outputPath)
	}

	// Ensure output directory exists
	if err := EnsureDir(outputPath); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write content to file
	if err := SafeWriteFile(outputPath, data); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (t Template) ReadIfExists(fileName string) ([]byte, error) {
	outputPath := t.GetFilePath(fileName)
	if !FileExists(outputPath) {
		return nil, fmt.Errorf("file does not exist: %s", outputPath)
	}

	return os.ReadFile(outputPath)
}
