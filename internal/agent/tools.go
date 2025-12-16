package agent

import (
	"context"
	"fmt"
)

type ToolFn func(ctx context.Context, state map[string]any) (map[string]any, error)

type ToolRegistry map[string]ToolFn

func SimpleSearchTool(ctx context.Context, state map[string]any) (map[string]any, error) {
	input, _ := state["input"].(string)

	result := map[string]any{
		"summary": fmt.Sprintf("Pretend search results for '%s'", input),
		"score":   0.9,
	}

	return result, nil
}
