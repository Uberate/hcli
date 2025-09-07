package config

import "github.io/uberate/hcli/pkg/io"

type AIConfig struct {
	Provider     string            `yaml:"Provider" comment:"AI provider, e.g. 'volc' for VolcEngine."`
	APIKey       string            `yaml:"APIKey" comment:"API key for the AI provider."`
	ThinkModel   string            `yaml:"ThinkModel" comment:"Model ID for text generation."`
	PicModel     string            `yaml:"PicModel" comment:"Model ID for image generation."`
	CustomPrompt map[string]string `yaml:"CustomPrompt" comment:"Custom prompts for different scenarios."`
}

type CliConfig struct {
	Templates []TemplateConfig `yaml:"Templates" comment:"Define the template of posts."`
	AI        AIConfig         `yaml:"AI" comment:"AI configuration for content generation."`
}

func DefaultCliConfig() CliConfig {
	return CliConfig{
		AI: AIConfig{
			Provider: "volc",
		},
	}
}

func ReadConfig(path string) (CliConfig, error) {
	c := DefaultCliConfig()

	return c, io.ReadYaml(path, &c)
}
