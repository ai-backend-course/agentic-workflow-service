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

type RAGSemanticSearchResponse struct {
	Query   string           `json:"query"`
	Results []map[string]any `json:"results"`
}

func RAGNotesSearchTool(ctx context.Context, state map[string]any) (map[string]any, error) {
	query, ok := state["input"].(string)
	if !ok || query == "" {
		return nil, errors.New("missing input for RAG search")
	}

	baseURL := os.Getenv("RAG_API_URL")
	if baseURL == "" {
		return nil, errors.New("missing RAG_API_URL")
	}

	payload := map[string]string{
		"query": query,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+"/search",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("RAG semantic search failed")
	}

	var ragResp RAGSemanticSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&ragResp); err != nil {
		return nil, err
	}

	return map[string]any{
		"query":   ragResp.Query,
		"results": ragResp.Results,
	}, nil
}
