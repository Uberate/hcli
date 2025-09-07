package cmds

import (
	"context"
	"errors"
	"github.io/uberate/hcli/pkg/config"
	"os"
	"path"
	"path/filepath"
)

func searchTemplate(ctx context.Context, name string) (config.TemplateConfig, error) {
	c := config.GetConfig(ctx)
	for _, item := range c.Templates {
		if item.Name == name {
			return item, nil
		}
	}

	return config.TemplateConfig{}, errors.New("template not found")
}

func parseFileName(ctx context.Context, dir, fileName string, needDir bool) string {
	fileDir := dir
	filePath := path.Join(fileDir, fileName+".md")
	if needDir {
		filePath = path.Join(fileDir, fileName, "index.md")
	}

	return filePath
}

func parseFileDir(ctx context.Context, dir, fileName string, needDir bool) string {
	fileDir := dir
	if needDir {
		fileDir = path.Join(fileDir, fileName)
	}

	return fileDir
}

func getActualFilePath(fileName string, tp config.TemplateConfig) (string, error) {
	// Use the same logic as gen_post to determine the actual file path
	if tp.NeedDir {
		// For NeedDir=true, the file should be in subdirectory/index.md
		return filepath.Join(tp.Dir, fileName, "index.md"), nil
	} else {
		// For NeedDir=false, the file should be fileName.md in the template directory
		return filepath.Join(tp.Dir, fileName+".md"), nil
	}
}

func readFileContent(fileName string) (string, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
