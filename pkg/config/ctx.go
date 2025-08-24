package config

import "context"

const (
	configKey = "config"
)

func SetConfig(ctx context.Context, c CliConfig) context.Context {
	return context.WithValue(ctx, configKey, c)
}

func GetConfig(ctx context.Context) CliConfig {
	res := ctx.Value(configKey)
	if r, ok := res.(CliConfig); ok {
		return r
	}

	panic("config was not set")
}
