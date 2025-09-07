package cmds

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.io/uberate/hcli/pkg/cctx"
	"github.io/uberate/hcli/pkg/config"
	"github.io/uberate/hcli/pkg/outputer"
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

	// Read file content - determine correct file path based on template config
	actualFilePath, err := getActualFilePath(fileName, tp)
	if err != nil {
		return fmt.Errorf("failed to determine file path for template %s: %w", tp.Name, err)
	}

	fileContent, err := readFileContent(actualFilePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s for template %s: %w", actualFilePath, tp.Name, err)
	}

	desc, err := aiClient.CreatePICSummary(ctx, fileContent)
	if err != nil {
		return err
	}

	// Save the prompt to file
	if err := savePromptToFile(ctx, fileName, tp, desc); err != nil {
		return fmt.Errorf("failed to save prompt: %w", err)
	}

	// Generate image using file content as prompt
	imageData, err := aiClient.GenPic(ctx, desc)
	if err != nil {
		return fmt.Errorf("failed to generate image: %w", err)
	}

	// Save or output the generated image
	return saveGeneratedImage(ctx, fileName, tp, imageData)
}

func savePromptToFile(ctx context.Context, fileName string, tp config.TemplateConfig, prompt string) error {
	// Create prompt filename based on image filename
	imageFileName := getOutputFileName(fileName, tp)
	promptFileName := strings.TrimSuffix(imageFileName, filepath.Ext(imageFileName)) + ".prompt.txt"

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(promptFileName), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write prompt to file
	if err := os.WriteFile(promptFileName, []byte(prompt), 0644); err != nil {
		return fmt.Errorf("failed to write prompt file: %w", err)
	}

	outputer.SuccessFL(ctx, "Prompt saved to: %s", promptFileName)
	return nil
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

	outputer.SuccessFL(ctx, "Generated picture saved to: %s", outputFileName)
	outputer.InfoFL(ctx, "Image size: %d bytes", len(imageData))

	return nil
}

func getOutputFileName(originalName string, tp config.TemplateConfig) string {
	baseName := filepath.Base(originalName)
	nameWithoutExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))

	// Use the same logic as gen_post for directory structure
	if tp.NeedDir && tp.Dir != "" {
		return getUniqueFilename(filepath.Join(tp.Dir, nameWithoutExt, "feature.png"))
	}

	// For non-dir templates, save in the same directory as original file
	return getUniqueFilename(filepath.Join(filepath.Dir(originalName), "feature.png"))
}

func getUniqueFilename(filePath string) string {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return filePath
	}

	// File exists, add suffix
	ext := filepath.Ext(filePath)
	base := strings.TrimSuffix(filePath, ext)

	for i := 1; i < 100; i++ {
		newPath := fmt.Sprintf("%s_%d%s", base, i, ext)
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
	}

	return filePath
}
