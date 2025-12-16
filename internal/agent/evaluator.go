package agent

import "strings"

type Evaluator interface {
	Evaluate(state map[string]any) (score float64, pass bool)
}

type GroundingEvaluator struct{}

func (e GroundingEvaluator) Evaluate(state map[string]any) (float64, bool) {
	search, ok := state["search"].(map[string]any)
	if !ok {
		return 0.2, false
	}

	results, ok := search["results"].([]map[string]any)
	if !ok || len(results) == 0 {
		return 0.3, false
	}

	final, ok := state["final_answer"].(map[string]any)
	if !ok {
		return 0.4, false
	}

	answer, _ := final["answer"].(string)

	// Simple grounding check:
	// does the answer reference any retrieved content?
	for _, r := range results {
		content, _ := r["content"].(string)
		if content != "" && strings.Contains(answer, content[:min(50, len(content))]) {
			return 0.9, true
		}
	}

	// Answer exists but not clearly grounded
	return 0.6, false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
