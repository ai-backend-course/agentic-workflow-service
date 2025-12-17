package agent

func DefaultToolRegistry() ToolRegistry {
	return ToolRegistry{
		"search":    RAGNotesSearchTool,
		"summarize": SummaryMicroserviceTool,
	}
}
