package aiclient

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultHTTPTimeout = 5 * time.Minute

type openAICompatibleClient struct {
	cfg    Config
	client *http.Client
}

func newOpenAICompatibleClient(cfg Config) *openAICompatibleClient {
	return &openAICompatibleClient{
		cfg: cfg,
		client: &http.Client{
			Timeout: defaultHTTPTimeout,
		},
	}
}

func (c *openAICompatibleClient) Config() Config {
	return c.cfg
}

type chatCompletionRequest struct {
	Model       string              `json:"model"`
	Messages    []chatMessageWire   `json:"messages"`
	Stream      bool                `json:"stream"`
	Temperature float64             `json:"temperature,omitempty"`
}

type chatMessageWire struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type streamChoiceDelta struct {
	Content string `json:"content"`
}

type streamChoice struct {
	Delta        streamChoiceDelta `json:"delta"`
	FinishReason *string           `json:"finish_reason"`
}

type streamChunkWire struct {
	Choices []streamChoice `json:"choices"`
	Error   *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (c *openAICompatibleClient) ChatStream(ctx context.Context, req ChatRequest) (<-chan StreamPart, error) {
	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("messages are required")
	}

	messages := make([]chatMessageWire, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = chatMessageWire{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}

	temp := c.cfg.ChatTemperature(req.Temperature)

	body, err := json.Marshal(chatCompletionRequest{
		Model:       c.cfg.Model,
		Messages:    messages,
		Stream:      true,
		Temperature: temp,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal chat request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.Endpoint(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create chat request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	if key := strings.TrimSpace(c.cfg.APIKey); key != "" {
		httpReq.Header.Set("Authorization", "Bearer "+key)
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("chat request failed: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		payload, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("chat api error (%d): %s", resp.StatusCode, strings.TrimSpace(string(payload)))
	}

	out := make(chan StreamPart, 32)
	go func() {
		defer close(out)
		defer resp.Body.Close()
		c.readSSE(ctx, resp.Body, out)
	}()

	return out, nil
}

func (c *openAICompatibleClient) readSSE(ctx context.Context, body io.Reader, out chan<- StreamPart) {
	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		if ctx.Err() != nil {
			emitPart(ctx, out, StreamPart{Type: StreamPartError, Error: ctx.Err().Error()})
			return
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" || !strings.HasPrefix(line, "data:") {
			continue
		}

		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			emitPart(ctx, out, StreamPart{Type: StreamPartFinish, FinishReason: "stop"})
			return
		}

		var chunk streamChunkWire
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}
		if chunk.Error != nil && chunk.Error.Message != "" {
			emitPart(ctx, out, StreamPart{Type: StreamPartError, Error: chunk.Error.Message})
			return
		}
		if len(chunk.Choices) == 0 {
			continue
		}

		choice := chunk.Choices[0]
		if text := choice.Delta.Content; text != "" {
			emitPart(ctx, out, StreamPart{Type: StreamPartTextDelta, Text: text})
		}
		if choice.FinishReason != nil && *choice.FinishReason != "" {
			emitPart(ctx, out, StreamPart{Type: StreamPartFinish, FinishReason: *choice.FinishReason})
			return
		}
	}

	if err := scanner.Err(); err != nil && ctx.Err() == nil {
		emitPart(ctx, out, StreamPart{Type: StreamPartError, Error: err.Error()})
	}
}

func emitPart(ctx context.Context, out chan<- StreamPart, part StreamPart) {
	select {
	case out <- part:
	case <-ctx.Done():
	}
}
