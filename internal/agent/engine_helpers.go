package agent

import (
	"fmt"
	"strings"
)

func buildSummaryContentFromRAG(results []map[string]any) string {
	var buf strings.Builder

	for _, r := range results {
		if text, ok := r["content"].(string); ok {
			buf.WriteString(text)
			buf.WriteString("\n\n")
		}
	}

	return buf.String()
}

func buildSummaryContentFromAny(v any) string {
	// Case 1: already the ideal type
	if results, ok := v.([]map[string]any); ok {
		return buildSummaryContentFromRAG(results)
	}

	// Case 2: []any (common when decoding JSON to interface{})
	if arr, ok := v.([]any); ok {
		converted := make([]map[string]any, 0, len(arr))
		for _, item := range arr {
			if m, ok := item.(map[string]any); ok {
				converted = append(converted, m)
			}
		}
		return buildSummaryContentFromRAG(converted)
	}

	// Fallback: stringify whatever it is
	return fmt.Sprintf("%v", v)
}
