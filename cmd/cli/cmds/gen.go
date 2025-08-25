package cmds

import (
	"context"
	"github.com/spf13/cobra"
	"github.io/uberate/hcli/pkg/config"
	"github.io/uberate/hcli/pkg/outputer"
)

func GenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "gen",
		Aliases: []string{"g"},
		Run: func(cmd *cobra.Command, args []string) {
			list(cmd.Context())
		},
	}

	cmd.AddCommand(genPost())

	return cmd
}

func list(ctx context.Context) {
	c := config.GetConfig(ctx)
	for _, item := range c.Templates {
		outputer.PrintFL(ctx, "---")
		outputer.PrintFL(ctx, "Name: %s", item.Name)
		outputer.PrintFL(ctx, "Categories: %v", item.Categories)
		outputer.PrintFL(ctx, "Tags: %v", item.Tags)
		outputer.PrintFL(ctx, "Dir: %s", item.Dir)
		outputer.PrintFL(ctx, "NeedDir: %v", item.NeedDir)
	}
}
