package outputer

import (
	"context"
)

const (
	outputLevelKey = "outputLevel"

	OutputLevelSilence = iota
	OutputLevelNormal
	OutputLevelDetail
)

func SetLevel(ctx context.Context, level int) context.Context {
	return context.WithValue(ctx, outputLevelKey, level)
}

func canOutput(ctx context.Context, l int) bool {
	v := ctx.Value(outputLevelKey)
	if level, ok := v.(int); ok {
		return l <= level
	}
	return l <= OutputLevelNormal // default is normal
}
