package cmds

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.io/uberate/hcli/pkg/config"
	"github.io/uberate/hcli/pkg/outputer"
	"github.io/uberate/hcli/pkg/template"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var overridePath string
var templateName string
var tags []string
var categories []string
var title string
var customArgs []string

func genPost() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "posts",
		Aliases: []string{"p", "post"},
		Short:   "generate post by specify template define",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fileName := args[0]
			if strings.HasSuffix(fileName, ".md") {
				fileName = strings.TrimSuffix(fileName, ".md")
			}
			if title == "" {
				title = getFileNameWithoutExtension(fileName)
			}

			return GenerateNewPost(cmd.Context(), fileName, title)
		},
	}

	cmd.Flags().StringVarP(&templateName, "template-name", "n", "", "the template names")
	cmd.Flags().StringSliceVarP(&tags, "tags", "t", nil, "the tags")
	cmd.Flags().StringSliceVarP(&categories, "categories", "", nil, "the tags")
	cmd.Flags().StringVarP(&overridePath, "override-path", "", "", "the override path")
	cmd.Flags().StringVarP(&title, "title", "", "", "change title, default will parse from name")
	cmd.Flags().StringSliceVarP(&customArgs, "custom-args", "a", []string{}, "custom args")

	return cmd
}

func GenerateNewPost(ctx context.Context, name string, title string) error {
	tp, err := searchTemplate(ctx, templateName)
	if err != nil {
		return err
	}

	rc, err := getRenderConfig(ctx, tp)
	if err != nil {
		return err
	}

	res, err := template.Render(ctx, rc)

	if err != nil {
		return err
	}

	writeToFile(ctx, tp.Dir, name, tp.NeedDir, res)

	outputer.DetailFL(ctx, "render: %s", res)

	return nil
}

func getRenderConfig(ctx context.Context, tp config.TemplateConfig) (template.RenderConfig, error) {

	rc := template.RenderConfig{
		Tags:       append(tags, tp.Tags...),
		Categories: append(categories, tp.Categories...),
		Time:       time.Now().Format(time.RFC3339),
		CustomArgs: map[string]string{},
		Title:      title,
		Temp:       tp.Template,
	}

	for _, arg := range customArgs {
		v := strings.Split(arg, "=")
		if len(v) != 2 {
			return rc, fmt.Errorf("invalid custom arg: %s, want key=value", arg)
		}
		rc.CustomArgs[v[0]] = strings.Join(v[1:], "=")
	}

	return rc, nil
}

func writeToFile(ctx context.Context, dir, fileName string, needDir bool, body string) error {
	fileDir := dir
	filePath := path.Join(fileDir, fileName+".md")
	if needDir {
		filePath = path.Join(fileDir, fileName, "index.md")
	}

	if err := os.MkdirAll(path.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	err := os.WriteFile(filePath, []byte(body), os.ModePerm)

	outputer.PrintFL(ctx, "output to: %s", filePath)
	return err
}

func searchTemplate(ctx context.Context, name string) (config.TemplateConfig, error) {
	c := config.GetConfig(ctx)
	for _, item := range c.Templates {
		if item.Name == name {
			return item, nil
		}
	}

	return config.TemplateConfig{}, errors.New("template not found")
}

func getFileNameWithoutExtension(filePath string) string {
	filename := filepath.Base(filePath)
	fileSuffix := filepath.Ext(filename)
	fileNameWithoutSuffix := strings.TrimSuffix(filename, fileSuffix)

	return fileNameWithoutSuffix
}
