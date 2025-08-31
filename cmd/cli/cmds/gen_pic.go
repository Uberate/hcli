package cmds

import "github.com/spf13/cobra"

func genPic() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pic",
		Short: "generate pictures",
	}

	return cmd
}
