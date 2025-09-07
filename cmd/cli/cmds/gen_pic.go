package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.io/uberate/hcli/pkg/cctx"
)

func genPic() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pic",
		Short: "generate pictures",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			prompt := args[0]

			// Get config from context
			cfg, ok := cctx.ConfigFromContext(ctx)
			if !ok {
				return fmt.Errorf("config not found in context")
			}

			// Get AI client from context
			aiClient, ok := cctx.AIClientFromContext(ctx)
			if !ok {
				return fmt.Errorf("AI client not found in context")
			}

			// Generate image using GenPic
			imageData, err := aiClient.GenPic(ctx, prompt)
			if err != nil {
				return fmt.Errorf("failed to generate image: %w", err)
			}

			// Output the generated image data
			fmt.Printf("Generated image data (length: %d bytes)\n", len(imageData))
			fmt.Printf("Image format: PNG (base64 decoded)\n")

			return nil
		},
	}

	return cmd
}
