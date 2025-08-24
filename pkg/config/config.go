package config

import "github.io/uberate/hcli/pkg/io"

type CliConfig struct {
	Templates []TemplateConfig `yaml:"Templates" comment:"Define the template of posts."`
}

func DefaultCliConfig() CliConfig {
	return CliConfig{}
}

func ReadConfig(path string) (CliConfig, error) {
	c := DefaultCliConfig()

	return c, io.ReadYaml(path, &c)
}
