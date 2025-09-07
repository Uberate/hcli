package cmds

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.io/uberate/hcli/pkg/cctx"
	"github.io/uberate/hcli/pkg/config"
	"os"
	"path/filepath"
	"strings"
)

var picTemplateName string

func genPic() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pic",
		Aliases: []string{"pictures", "image"},
		Short:   "generate pictures from template content",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("filename argument is required")
			}

			ctx := cmd.Context()
			fileName := args[0]

			if strings.HasSuffix(fileName, ".md") {
				fileName = strings.TrimSuffix(fileName, ".md")
			}

			return GeneratePictureFromTemplate(ctx, fileName, picTemplateName)
		},
	}

	cmd.Flags().StringVarP(&picTemplateName, "template-name", "n", "", "the template name for picture generation")

	return cmd
}

func GeneratePictureFromTemplate(ctx context.Context, fileName, templateName string) error {

	// Get AI client from context
	aiClient, ok := cctx.AIClientFromContext(ctx)
	if !ok {
		return fmt.Errorf("AI client not found in context")
	}

	// Search for template
	tp, err := searchTemplate(ctx, templateName)
	if err != nil {
		return err
	}

	// Read file content
	fileContent, err := readFileContent(fileName)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", fileName, err)
	}

	aiClient.CreatePICSummary(ctx)

	// Generate image using file content as prompt
	imageData, err := aiClient.GenPic(ctx, fileContent)
	if err != nil {
		return fmt.Errorf("failed to generate image: %w", err)
	}

	// Save or output the generated image
	return saveGeneratedImage(ctx, fileName, tp, imageData)
}

func readFileContent(fileName string) (string, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func saveGeneratedImage(ctx context.Context, fileName string, tp config.TemplateConfig, imageData []byte) error {
	// Create output filename based on template configuration
	outputFileName := getOutputFileName(fileName, tp)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(outputFileName), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write image data to file
	if err := os.WriteFile(outputFileName, imageData, 0644); err != nil {
		return fmt.Errorf("failed to write image file: %w", err)
	}

	fmt.Printf("Generated picture saved to: %s\n", outputFileName)
	fmt.Printf("Image size: %d bytes\n", len(imageData))

	return nil
}

func getOutputFileName(originalName string, tp config.TemplateConfig) string {
	baseName := filepath.Base(originalName)
	nameWithoutExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))

	if tp.NeedDir && tp.Dir != "" {
		return filepath.Join(tp.Dir, nameWithoutExt, "feature.png")
	}

	return nameWithoutExt + ".png"
}
