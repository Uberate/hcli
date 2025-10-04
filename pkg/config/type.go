package config

import (
	"errors"
	"github.io/uberate/hcli/pkg/fileio"
	"github.io/uberate/hcli/pkg/llms"
	"github.io/uberate/hcli/pkg/template"
)

type CliConfig struct {
	Templates []template.Template `yaml:"Templates" describe:"Define the template of posts."`
	LLMs      llms.Config         `yaml:"LLMs" describe:"Define the LLMs configuration"`
}

func (cc CliConfig) SearchTemplate(name string) (template.Template, error) {
	for _, t := range cc.Templates {
		if t.Name == name {
			return t, nil
		}
	}
	return template.Template{}, errors.New("template not found")
}

func DefaultCliConfig() CliConfig {
	return CliConfig{}
}

func ExampleCliConfig() CliConfig {
	return CliConfig{
		Templates: []template.Template{
			{
				Name: "template-name",
				Categories: []string{
					"categories",
				},
				Tags: []string{
					"tag1",
					"tag2",
				},
				Template:         "",
				Dir:              "",
				NeedDir:          false,
				PicSummaryPrompt: "",
			},
		},
		LLMs: llms.Config{
			Provider: "volc",
			VolcEngineConfig: llms.VolcEngineConfig{
				ApiKey:    "123",
				TextModel: "123",
				PicModel:  "123",
			},
		},
	}
}

func ReadConfig(path string) (CliConfig, error) {
	c := DefaultCliConfig()
	return c, fileio.ReadYaml(path, &c)
}
