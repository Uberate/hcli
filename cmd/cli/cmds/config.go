package cmds

import (
	"github.com/spf13/cobra"
	"github.io/uberate/hcli/pkg/config"
	"github.io/uberate/hcli/pkg/hctx"
	"github.io/uberate/hcli/pkg/yamlutil"
)

func ConfigCmd() *cobra.Command {
	configCmd := &cobra.Command{
		Use: "config",
	}

	configCmd.AddCommand(
		demoCmd(),
	)

	return configCmd
}

func demoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "demo",

		Run: func(cmd *cobra.Command, args []string) {
			exampleConfigBytes, err := yamlutil.Render(config.ExampleCliConfig())
			if err != nil {
				hctx.Err(cmd.Context(), "BUG: generate default yaml fail, err: %v", err)
			}
			hctx.Println(cmd.Context(), string(exampleConfigBytes))
		},
	}

	return cmd
}
