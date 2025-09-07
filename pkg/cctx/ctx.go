package cctx

import (
	"context"
	"github.io/uberate/hcli/pkg/ais"
	"github.io/uberate/hcli/pkg/config"
)

const (
	configKey = "config"
	aiClientKey = "ai_client"
)

func ConfigFromContext(ctx context.Context) (config.CliConfig, bool) {
	res := ctx.Value(configKey)
	if r, ok := res.(config.CliConfig); ok {
		return r, true
	}
	return config.CliConfig{}, false
}

func AIClientFromContext(ctx context.Context) (ais.AIs, bool) {
	res := ctx.Value(aiClientKey)
	if r, ok := res.(ais.AIs); ok {
		return r, true
	}
	return nil, false
}

func SetConfig(ctx context.Context, cfg config.CliConfig) context.Context {
	return context.WithValue(ctx, configKey, cfg)
}

func SetAIClient(ctx context.Context, client ais.AIs) context.Context {
	return context.WithValue(ctx, aiClientKey, client)
}
