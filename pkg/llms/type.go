package llms

import (
	"context"
	"errors"
	"os"
)

type Config struct {
	Provider         string           `yaml:"Provider" describe:"The provider of LLM, support: volc"`
	VolcEngineConfig VolcEngineConfig `yaml:"VolcEngineConfig" describe:"VolcEngineConfig"`
}

const (
	VolcEngineLLMKey = "volc"
)

type LLMTools interface {
	Text(ctx context.Context, sysPrompt, input string) (resp string, err error)
	Pic(ctx context.Context, input string) (resp []byte, err error)
}

func NewLLM(c Config) (LLMTools, error) {
	switch c.Provider {
	case VolcEngineLLMKey:
		return NewVolcEngineLLM(c.VolcEngineConfig), nil
	}

	return nil, errors.New("invalid provider, support: [volc]")
}

func NewLLMWithAutoEnv(c Config) (LLMTools, error) {
	switch c.Provider {
	case VolcEngineLLMKey:
		if c.VolcEngineConfig.ApiKey == "" {
			c.VolcEngineConfig.ApiKey = os.Getenv("VOLC_API_KEY")
		}
		if c.VolcEngineConfig.TextModel == "" {
			c.VolcEngineConfig.TextModel = os.Getenv("THINK_MODEL_ID")
		}
		if c.VolcEngineConfig.PicModel == "" {
			c.VolcEngineConfig.PicModel = os.Getenv("PIC_MODEL_ID")
		}

		return NewVolcEngineLLM(c.VolcEngineConfig), nil
	}

	return nil, errors.New("invalid provider, support: [volc]")
}
