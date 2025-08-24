package cmds

import (
	"github.com/spf13/cobra"
	"github.io/uberate/hcli/pkg/config"
	"github.io/uberate/hcli/pkg/outputer"
	"github.io/uberate/hcli/pkg/yamlcomm"
)

func demoConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "demo-config",
		Aliases: []string{"demo"},
		Short:   "show demo config",

		RunE: func(cmd *cobra.Command, args []string) error {
			outputer.DetailFL(cmd.Context(), "show demo config")
			v, err := yamlcomm.MarshalWithComments(demoConfigValue)
			if err != nil {
				return err
			}
			outputer.ForceFL(cmd.Context(), string(v))
			return nil
		},
	}

	return cmd
}

var demoConfigValue = config.CliConfig{
	Templates: []config.TemplateConfig{
		{
			Name:       "test-template",
			Categories: []string{"tech"},
			Template:   "+++",
		},
	},
}
