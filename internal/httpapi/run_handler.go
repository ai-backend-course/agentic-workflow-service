package httpapi

import (
	"agentic-workflow-service/internal/agent"
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h Handlers) runAgent(c *fiber.Ctx) error {
	runID := uuid.NewString()

	var req struct {
		Input string `json:"input"`
	}

	// Parse incoming JSON
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Timeout for agent execution
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)
	ctx = context.WithValue(ctx, "run_id", runID)
	defer cancel()

	// Load prompts from internal/prompts/
	pl, err := agent.NewEmbeddedPromptLoader()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to load prompts",
		})
	}

	// Initialize LLM
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return c.Status(500).JSON(fiber.Map{
			"error": "missing OPENAI_API_KEY",
		})
	}
	agent.InitLLM((apiKey))

	// Tools
	tools := agent.ToolRegistry{
		"search": agent.RAGNotesSearchTool,
	}

	evaluator := agent.GroundingEvaluator{}

	engine := agent.NewEngine(pl, tools, evaluator)

	// Run agent workflow
	result, err := engine.Run(ctx, map[string]any{
		"input":  req.Input,
		"run_id": runID,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "ok",
		"run_id": runID,
		"result": result,
	})
}
