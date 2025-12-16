package agent

type StepType string

const (
	StepLLM      StepType = "llm"
	StepTool     StepType = "tool"
	StepEvaluate StepType = "evaluate"
)

type Step struct {
	Type   StepType
	Prompt string // used for LLM steps
	Tool   string // used for tool steps
}

// Our workflow definition (in-memory)
var Workflow = []Step{
	{Type: StepLLM, Prompt: "intent_decider"},
	{Type: StepTool, Tool: "search"},
	{Type: StepEvaluate},
	{Type: StepLLM, Prompt: "final_answer"},
}
