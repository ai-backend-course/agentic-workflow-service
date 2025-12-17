package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

type SummaryServiceResponse struct {
	Summary string `json:"summary"`
}

func SummaryMicroserviceTool(ctx context.Context, state map[string]any) (map[string]any, error) {
	content, ok := state["content"].(string)
	if !ok || content == "" {
		return nil, errors.New("missing content for summarization")
	}

	baseURL := os.Getenv("SUMMARY_API_URL")
	if baseURL == "" {
		return nil, errors.New("missing SUMMARY_API_URL")
	}

	payload := map[string]string{
		"text": content,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+"/summary",
		bytes.NewBuffer((body)),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("summary microservice failed")
	}

	var summaryResp SummaryServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&summaryResp); err != nil {
		return nil, err
	}

	if summaryResp.Summary == "" {
		return nil, errors.New("empty summary returned")
	}

	return map[string]any{
		"summary": summaryResp.Summary,
	}, nil
}
