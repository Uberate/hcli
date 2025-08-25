package cmds

import "github.com/spf13/cobra"

func ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "config",
	}

	cmd.AddCommand(demoConfig())

	return cmd
}
