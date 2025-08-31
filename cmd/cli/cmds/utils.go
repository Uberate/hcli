package cmds

import (
	"context"
	"errors"
	"github.io/uberate/hcli/pkg/config"
	"path"
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
