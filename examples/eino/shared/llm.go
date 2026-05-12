package shared

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/yong/doc-generation-mcp-server/internal/formaldoc"
)

type LLMConfig struct {
	Mode    string
	BaseURL string
	APIKey  string
	Model   string
	Timeout time.Duration
}

func LoadLLMConfig() LLMConfig {
	mode := os.Getenv("LLM_MODE")
	if mode == "" {
		mode = "mock"
	}
	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-4o-mini"
	}
	return LLMConfig{
		Mode:    mode,
		BaseURL: strings.TrimRight(baseURL, "/"),
		APIKey:  os.Getenv("OPENAI_API_KEY"),
		Model:   model,
		Timeout: 60 * time.Second,
	}
}

func BuildDraftWithLLM(ctx context.Context, cfg LLMConfig, input any) (formaldoc.Draft, string, error) {
	if cfg.Mode != "openai" {
		return formaldoc.Draft{}, "", fmt.Errorf("unsupported llm mode: %s", cfg.Mode)
	}
	if strings.TrimSpace(cfg.APIKey) == "" {
		return formaldoc.Draft{}, "", fmt.Errorf("OPENAI_API_KEY is required when LLM_MODE=openai")
	}
	userJSON, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return formaldoc.Draft{}, "", err
	}

	type chatMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	type chatRequest struct {
		Model          string        `json:"model"`
		Messages       []chatMessage `json:"messages"`
		Temperature    float64       `json:"temperature"`
		ResponseFormat any           `json:"response_format,omitempty"`
	}
	type chatResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	reqBody := chatRequest{
		Model:       cfg.Model,
		Temperature: 0.2,
		Messages: []chatMessage{
			{
				Role: "system",
				Content: strings.TrimSpace(`You are a formal document drafting assistant.
Output only valid JSON matching FormalDocumentDraftV1.
Rules:
- schema_version must be "1.0"
- document_type must be one of: project_proposal, weekly_report, business_letter
- tone must be "formal"
- language must be "zh-CN"
- do not output markdown
- do not explain
- do not invent missing facts; use "待补充" or review_notes when needed`),
			},
			{
				Role:    "user",
				Content: "Build a FormalDocumentDraftV1 JSON document from the following input:\n" + string(userJSON),
			},
		},
		ResponseFormat: map[string]string{"type": "json_object"},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return formaldoc.Draft{}, "", err
	}

	httpClient := &http.Client{Timeout: cfg.Timeout}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.BaseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return formaldoc.Draft{}, "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return formaldoc.Draft{}, "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return formaldoc.Draft{}, "", err
	}
	if resp.StatusCode >= 400 {
		return formaldoc.Draft{}, "", fmt.Errorf("openai request failed: %s", strings.TrimSpace(string(body)))
	}

	var decoded chatResponse
	if err := json.Unmarshal(body, &decoded); err != nil {
		return formaldoc.Draft{}, "", err
	}
	if len(decoded.Choices) == 0 {
		return formaldoc.Draft{}, "", fmt.Errorf("openai response contains no choices")
	}
	raw := decoded.Choices[0].Message.Content
	var draft formaldoc.Draft
	if err := json.Unmarshal([]byte(raw), &draft); err != nil {
		return formaldoc.Draft{}, raw, err
	}
	return draft, raw, nil
}
