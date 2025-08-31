package cmds

import (
	"github.com/spf13/cobra"
)

func GenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "gen",
		Aliases: []string{"g"},
		Run:     func(cmd *cobra.Command, args []string) {},
	}

	cmd.AddCommand(
		genPost(),
		genPic(),
	)

	return cmd
}
