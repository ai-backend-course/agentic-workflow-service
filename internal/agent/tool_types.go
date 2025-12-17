package agent

import (
	"context"
)

type ToolFn func(ctx context.Context, state map[string]any) (map[string]any, error)

type ToolRegistry map[string]ToolFn
