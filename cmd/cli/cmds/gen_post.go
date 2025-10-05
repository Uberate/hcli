package cmds

import (
	"context"
	"errors"
	"github.com/spf13/cobra"
	"github.io/uberate/hcli/pkg/config"
	"github.io/uberate/hcli/pkg/hctx"
	"github.io/uberate/hcli/pkg/template"
	"path/filepath"
	"strings"
)

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
		Args:    cobra.MaximumNArgs(1),
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

	cmd.Flags().StringVarP(&templateName, "template", "n", "", "template name")
	cmd.Flags().StringSliceVarP(&tags, "tags", "t", nil, "tags")
	cmd.Flags().StringSliceVarP(&categories, "categories", "k", nil, "categories")
	cmd.Flags().StringVarP(&title, "title", "", "", "title")
	cmd.Flags().StringSliceVarP(&customArgs, "custom-args", "a", nil, "custom args")

	return cmd
}

func GenerateNewPost(ctx context.Context, name string, title string) error {
	// read config first
	hctx.Debug(ctx, "config path: %s", hctx.GetConfigPath(ctx))
	c, err := config.ReadConfig(hctx.GetConfigPath(ctx))
	if err != nil {
		return errors.New("read config error: " + err.Error())
	}

	tp, err := c.SearchTemplate(templateName)
	if err != nil {
		return errors.New("search template error: " + err.Error())
	}

	hctx.Debug(ctx, "%+v", tp)

	args := map[string]string{}

	for _, item := range customArgs {
		values := strings.Split(item, "=")
		if len(values) < 2 {
			return errors.New("wrong custom args format, need k=v")
		}

		args[values[0]] = strings.Join(values[1:], "=")
	}

	return template.RenderToFile(ctx, &tp, name, template.RenderOption{
		Title:            title,
		AppendTags:       tags,
		AppendCategories: categories,
		CustomArgs:       args,
	})
}

func getFileNameWithoutExtension(filePath string) string {
	filename := filepath.Base(filePath)
	fileSuffix := filepath.Ext(filename)
	fileNameWithoutSuffix := strings.TrimSuffix(filename, fileSuffix)

	return fileNameWithoutSuffix
}
