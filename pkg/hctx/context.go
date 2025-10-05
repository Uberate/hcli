package hctx

import (
	"context"
	"fmt"
	"github.io/uberate/hcli/pkg/output"
	"os"
)

const (
	kTitle      = "k_title"
	kOutputter  = "k_outputter"
	kConfig     = "k_config"
	kConfigPath = "k_config_path"
)

func SetConfigPath(ctx context.Context, configPath string) context.Context {
	return context.WithValue(ctx, kConfigPath, configPath)
}

func GetConfigPath(ctx context.Context) string {
	res := ctx.Value(kConfigPath)
	if r, ok := res.(string); ok {
		return r
	}
	return ""
}

// ---------------------- command args if exists

// ---------------------- outputs

func SetOutputter(ctx context.Context, outputer *output.Outputter) context.Context {
	return context.WithValue(ctx, kOutputter, outputer)
}

func getOutputter(ctx context.Context) *output.Outputter {
	res := ctx.Value(kOutputter)
	if out, ok := res.(*output.Outputter); ok {
		return out
	}

	return output.NewOutputter(output.LevelInfo, os.Stdout)
}

func Println(ctx context.Context, format string, args ...interface{}) {
	op := getOutputter(ctx)
	op.Println(fmt.Sprintf(format, args...))
}

func Err(ctx context.Context, format string, args ...interface{}) {
	op := getOutputter(ctx)
	op.Error(format, args...)
}

func Debug(ctx context.Context, format string, args ...interface{}) {
	op := getOutputter(ctx)
	op.Debug(format, args...)
}
