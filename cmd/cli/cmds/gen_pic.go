package cmds

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.io/uberate/hcli/pkg/config"
	"github.io/uberate/hcli/pkg/hctx"
	"github.io/uberate/hcli/pkg/llms"
	"github.io/uberate/hcli/pkg/poster"
)

var picTemplateName string

func genPic() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pic",
		Aliases: []string{"pictures", "image"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			fileName := args[0]

			return GeneratePictureFromTemplate(ctx, fileName, picTemplateName)
		},
	}

	cmd.Flags().StringVarP(&picTemplateName, "template-name", "n", "", "the template name for picture generation")

	return cmd
}

func GeneratePictureFromTemplate(ctx context.Context, fileName, templateName string) error {
	c, err := config.ReadConfig(hctx.GetConfigPath(ctx))
	if err != nil {
		return err
	}

	llmTools, err := llms.NewLLMWithAutoEnv(c.LLMs)
	if err != nil {
		return err
	}

	tp, err := c.SearchTemplate(templateName)
	if err != nil {
		return err
	}

	if !tp.NeedDir {
		return errors.New(fmt.Sprintf("template %s has no dir, pic can't generate", templateName))
	}

	generateResult, err := poster.GeneratePoster(ctx, poster.GeneratePosterArgs{
		TP:       &tp,
		LLMTools: llmTools,
		FileName: fileName,
	})

	if err != nil {
		return err
	}

	if err = tp.WritePicSummary(fileName, generateResult.Summary); err != nil {
		return err
	}
	if err = tp.WritePoster(fileName, generateResult.Pic); err != nil {
		return err
	}

	return nil
}
