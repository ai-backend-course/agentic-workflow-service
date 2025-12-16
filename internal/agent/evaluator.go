package agent

type Evaluator interface {
	Evaluate(state map[string]any) (score float64, pass bool)
}

type DefaultEvaluator struct{}

// Simple evaluator: always returns pass=true
func (e DefaultEvaluator) Evaluate(state map[string]any) (float64, bool) {
	return 1.0, true
}
