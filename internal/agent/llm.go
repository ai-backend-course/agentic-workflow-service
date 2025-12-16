package agent

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/sashabaranov/go-openai"
)

var openAIClient *openai.Client

var ErrLLMNotInitialized = errors.New("LLM client not initialized")

func InitLLM(apiKey string) {
	openAIClient = openai.NewClient(apiKey)
}

func CallLLM(prompt string) (map[string]any, error) {
	if openAIClient == nil {
		return nil, ErrLLMNotInitialized
	}
	ctx := context.Background()

	resp, err := openAIClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt + "\n\nReturn your answer as valid JSON.",
			},
		},
	})

	if err != nil {
		return nil, err
	}

	// LLM returns JSON object -> convert to map[string]any
	content := resp.Choices[0].Message.Content
	clean := extractJSON(content)

	if !strings.HasPrefix(strings.TrimSpace(clean), "{") {
		return nil, errors.New("LLM did not return valid JSON")
	}

	var out map[string]any
	if err := json.Unmarshal([]byte(clean), &out); err != nil {
		return nil, err
	}

	return out, nil

}

func extractJSON(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")

	if start == -1 || end == -1 || end <= start {
		return s
	}

	return s[start : end+1]
}
