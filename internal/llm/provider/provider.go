package provider

import (
	"context"

	"github.com/imnulhaqueruman/opencode-poc/internal/llm/tools"
	"github.com/imnulhaqueruman/opencode-poc/internal/message"
)

// EventType represents the type of streaming event
type EventType string

const (
	EventContentStart  EventType = "content_start"
	EventContentDelta  EventType = "content_delta"
	EventThinkingDelta EventType = "thinking_delta"
	EventContentStop   EventType = "content_stop"
	EventComplete      EventType = "complete"
	EventError         EventType = "error"
)

type TokenUsage struct {
	InputTokens         int64
	OutputTokens        int64
	CacheCreationTokens int64
	CacheReadTokens     int64
}

type ProviderResponse struct {
	Content   string
	ToolCalls []message.ToolCall
	Usage     TokenUsage
}

type ProviderEvent struct {
	Type     EventType
	Content  string
	Thinking string
	ToolCall *message.ToolCall
	Error    error
	Response *ProviderResponse
}

type Provider interface {
	SendMessages(ctx context.Context, messages []message.Message, tools []tools.BaseTool) (*ProviderResponse, error)

	StreamResponse(ctx context.Context, messages []message.Message, tools []tools.BaseTool) (<-chan ProviderEvent, error)
}