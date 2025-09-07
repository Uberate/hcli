package outputer

import (
	"context"
	"fmt"
	"strings"
)

func ForceFL(ctx context.Context, f string, args ...any) {
	s := fmt.Sprintf(f, args...)
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	fmt.Print(s)
}

func PrintFL(ctx context.Context, f string, args ...any) {
	if canOutput(ctx, OutputLevelNormal) {
		s := fmt.Sprintf(f, args...)
		if !strings.HasSuffix(s, "\n") {
			s += "\n"
		}
		fmt.Print(s)
	}
}

func DetailFL(ctx context.Context, f string, args ...any) {
	if canOutput(ctx, OutputLevelDetail) {
		s := fmt.Sprintf("> "+f, args...)
		if !strings.HasSuffix(s, "\n") {
			s += "\n"
		}
		fmt.Print(s)
	}
}

func InfoFL(ctx context.Context, f string, args ...any) {
	if canOutput(ctx, OutputLevelNormal) {
		s := fmt.Sprintf("ℹ️  "+f, args...)
		if !strings.HasSuffix(s, "\n") {
			s += "\n"
		}
		fmt.Print(s)
	}
}

func SuccessFL(ctx context.Context, f string, args ...any) {
	if canOutput(ctx, OutputLevelNormal) {
		s := fmt.Sprintf("✅ "+f, args...)
		if !strings.HasSuffix(s, "\n") {
			s += "\n"
		}
		fmt.Print(s)
	}
}

func WarningFL(ctx context.Context, f string, args ...any) {
	if canOutput(ctx, OutputLevelNormal) {
		s := fmt.Sprintf("⚠️  "+f, args...)
		if !strings.HasSuffix(s, "\n") {
			s += "\n"
		}
		fmt.Print(s)
	}
}

func ErrorFL(ctx context.Context, f string, args ...any) {
	if canOutput(ctx, OutputLevelNormal) {
		s := fmt.Sprintf("❌ "+f, args...)
		if !strings.HasSuffix(s, "\n") {
			s += "\n"
		}
		fmt.Print(s)
	}
}

func StepFL(ctx context.Context, f string, args ...any) {
	if canOutput(ctx, OutputLevelNormal) {
		s := fmt.Sprintf("🚀 "+f, args...)
		if !strings.HasSuffix(s, "\n") {
			s += "\n"
		}
		fmt.Print(s)
	}
}
