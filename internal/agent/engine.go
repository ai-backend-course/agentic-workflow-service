package agent

import (
	"context"
	"errors"
	"log"
)

type Engine struct {
	Prompts PromptLoader
	Tools   ToolRegistry
	Eval    Evaluator
}

func NewEngine(pl PromptLoader, tr ToolRegistry, ev Evaluator) *Engine {
	return &Engine{Prompts: pl, Tools: tr, Eval: ev}
}

func (e *Engine) Run(ctx context.Context, state map[string]any) (map[string]any, error) {
	runID, _ := state["run_id"].(string)

	for _, step := range Workflow {
		log.Printf(
			"[agent] run=%s step_type=%s step=%s",
			runID,
			step.Type,
			func() string {
				if step.Prompt != "" {
					return step.Prompt
				}
				return step.Tool
			}(),
		)

		switch step.Type {

		case StepLLM:
			prompt, err := e.Prompts.Render(step.Prompt, state)
			if err != nil {
				return nil, err
			}
			out, err := CallLLM(prompt)
			if err != nil {
				return nil, err
			}
			state[step.Prompt] = out

		case StepTool:
			toolFn, ok := e.Tools[step.Tool]
			if !ok {
				return nil, errors.New("unknown tool: " + step.Tool)
			}
			out, err := toolFn(ctx, state)
			if err != nil {
				return nil, err
			}
			state[step.Tool] = out

			if step.Tool == "search" {
				state["content"] = buildSummaryContentFromAny(out["results"])
			}

		case StepEvaluate:
			score, pass := e.Eval.Evaluate(state)
			state["evaluation_score"] = score

			log.Printf(
				"[agent] run=%s evaluation_score=%.2f pass=%v",
				runID,
				score,
				pass,
			)

			if !pass {
				continue
			}
		}
	}

	return state, nil
}
