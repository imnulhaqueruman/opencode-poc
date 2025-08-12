package provider

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/imnulhaqueruman/opencode-poc/internal/llm/models"
	"github.com/imnulhaqueruman/opencode-poc/internal/llm/tools"
	"github.com/imnulhaqueruman/opencode-poc/internal/message"
)

type anthropicProvider struct {
	client        *anthropic.Client
	model         models.Model
	maxTokens     int64
	apiKey        string
	systemMessage string
}

type AnthropicOption func(*anthropicProvider)

func WithAnthropicSystemMessage(message string) AnthropicOption {
	return func(a *anthropicProvider) {
		a.systemMessage = message
	}
}

func WithAnthropicMaxTokens(maxTokens int64) AnthropicOption {
	return func(a *anthropicProvider) {
		a.maxTokens = maxTokens
	}
}

func WithAnthropicModel(model models.Model) AnthropicOption {
	return func(a *anthropicProvider) {
		a.model = model
	}
}

func WithAnthropicKey(apiKey string) AnthropicOption {
	return func(a *anthropicProvider) {
		a.apiKey = apiKey
	}
}

func NewAnthropicProvider(opts ...AnthropicOption) (Provider, error) {
	provider := &anthropicProvider{
		maxTokens: 1024,
	}

	for _, opt := range opts {
		opt(provider)
	}

	if provider.systemMessage == "" {
		return nil, errors.New("system message is required")
	}

	provider.client = anthropic.NewClient(option.WithAPIKey(provider.apiKey))
	return provider, nil
}

func (a *anthropicProvider) SendMessages(ctx context.Context, messages []message.Message, tools []tools.BaseTool) (*ProviderResponse, error) {
	anthropicMessages := a.convertToAnthropicMessages(messages)
	anthropicTools := a.convertToAnthropicTools(tools)

	response, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:       anthropic.F(anthropic.Model(a.model.APIModel)),
		MaxTokens:   anthropic.F(a.maxTokens),
		Temperature: anthropic.F(0.0),
		Messages:    anthropic.F(anthropicMessages),
		Tools:       anthropic.F(anthropicTools),
		System: anthropic.F([]anthropic.TextBlockParam{
			{
				Text: anthropic.F(a.systemMessage),
				CacheControl: anthropic.F(anthropic.CacheControlEphemeralParam{
					Type: anthropic.F(anthropic.CacheControlEphemeralType("ephemeral")),
				}),
			},
		}),
	})
	if err != nil {
		return nil, err
	}

	content := ""
	for _, block := range response.Content {
		if block.Type == anthropic.ContentBlockTypeText {
			content += block.Text
		}
	}

	toolCalls := a.extractToolCalls(response.Content)
	tokenUsage := a.extractTokenUsage(response.Usage)

	return &ProviderResponse{
		Content:   content,
		ToolCalls: toolCalls,
		Usage:     tokenUsage,
	}, nil
}

func (a *anthropicProvider) StreamResponse(ctx context.Context, messages []message.Message, tools []tools.BaseTool) (<-chan ProviderEvent, error) {
	anthropicMessages := a.convertToAnthropicMessages(messages)
	anthropicTools := a.convertToAnthropicTools(tools)

	temperature := anthropic.F(0.0)
	lastMessage := messages[len(messages)-1]
	if lastMessage.Role == message.User && strings.Contains(strings.ToLower(lastMessage.Content), "think") {
		temperature = anthropic.F(1.0)
	}

	stream := a.client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
		Model:       anthropic.F(anthropic.Model(a.model.APIModel)),
		MaxTokens:   anthropic.F(a.maxTokens),
		Temperature: temperature,
		Messages:    anthropic.F(anthropicMessages),
		Tools:       anthropic.F(anthropicTools),
		System: anthropic.F([]anthropic.TextBlockParam{
			{
				Text: anthropic.F(a.systemMessage),
				CacheControl: anthropic.F(anthropic.CacheControlEphemeralParam{
					Type: anthropic.F(anthropic.CacheControlEphemeralType("ephemeral")),
				}),
			},
		}),
	})

	eventChan := make(chan ProviderEvent)

	go func() {
		defer close(eventChan)

		accumulatedMessage := anthropic.Message{}

		for stream.Next() {
			event := stream.Current()
			err := accumulatedMessage.Accumulate(event)
			if err != nil {
				eventChan <- ProviderEvent{Type: EventError, Error: err}
				return
			}

			switch event.Type {
			case "content_block_start":
				eventChan <- ProviderEvent{Type: EventContentStart}

			case "content_block_delta":
				// Skip delta handling for now to get build working

			case "content_block_stop":
				eventChan <- ProviderEvent{Type: EventContentStop}

			case "message_stop":
				content := ""
				for _, block := range accumulatedMessage.Content {
					if block.Type == anthropic.ContentBlockTypeText {
						content += block.Text
					}
				}

				toolCalls := a.extractToolCalls(accumulatedMessage.Content)
				tokenUsage := a.extractTokenUsage(accumulatedMessage.Usage)

				eventChan <- ProviderEvent{
					Type: EventComplete,
					Response: &ProviderResponse{
						Content:   content,
						ToolCalls: toolCalls,
						Usage:     tokenUsage,
					},
				}
			}
		}

		if stream.Err() != nil {
			eventChan <- ProviderEvent{Type: EventError, Error: stream.Err()}
		}
	}()

	return eventChan, nil
}

func (a *anthropicProvider) extractToolCalls(content []anthropic.ContentBlock) []message.ToolCall {
	var toolCalls []message.ToolCall

	for _, block := range content {
		if block.Type == anthropic.ContentBlockTypeToolUse {
			toolCall := message.ToolCall{
				ID:    block.ID,
				Name:  block.Name,
				Input: string(block.Input),
				Type:  "tool_use",
			}
			toolCalls = append(toolCalls, toolCall)
		}
	}

	return toolCalls
}

func (a *anthropicProvider) extractTokenUsage(usage anthropic.Usage) TokenUsage {
	return TokenUsage{
		InputTokens:         usage.InputTokens,
		OutputTokens:        usage.OutputTokens,
		CacheCreationTokens: usage.CacheCreationInputTokens,
		CacheReadTokens:     usage.CacheReadInputTokens,
	}
}

func (a *anthropicProvider) convertToAnthropicTools(tools []tools.BaseTool) []anthropic.ToolParam {
	anthropicTools := make([]anthropic.ToolParam, len(tools))

	for i, tool := range tools {
		info := tool.Info()
		toolParam := anthropic.ToolParam{
			Name:        anthropic.F(info.Name),
			Description: anthropic.F(info.Description),
			InputSchema: anthropic.F(interface{}(info.Parameters)),
		}

		if i == len(tools)-1 {
			toolParam.CacheControl = anthropic.F(anthropic.CacheControlEphemeralParam{
				Type: anthropic.F(anthropic.CacheControlEphemeralType("ephemeral")),
			})
		}

		anthropicTools[i] = toolParam
	}

	return anthropicTools
}

func (a *anthropicProvider) convertToAnthropicMessages(messages []message.Message) []anthropic.MessageParam {
	anthropicMessages := make([]anthropic.MessageParam, len(messages))
	cachedBlocks := 0

	for i, msg := range messages {
		switch msg.Role {
		case message.User:
			content := anthropic.TextBlockParam{
				Type: anthropic.F(anthropic.TextBlockParamTypeText),
				Text: anthropic.F(msg.Content),
			}
			if cachedBlocks < 2 {
				content.CacheControl = anthropic.F(anthropic.CacheControlEphemeralParam{
					Type: anthropic.F(anthropic.CacheControlEphemeralType("ephemeral")),
				})
				cachedBlocks++
			}
			anthropicMessages[i] = anthropic.MessageParam{
				Role:    anthropic.F(anthropic.MessageParamRoleUser),
				Content: anthropic.F([]anthropic.ContentBlockParamUnion{content}),
			}

		case message.Assistant:
			blocks := []anthropic.ContentBlockParamUnion{}
			if msg.Content != "" {
				content := anthropic.TextBlockParam{
					Type: anthropic.F(anthropic.TextBlockParamTypeText),
					Text: anthropic.F(msg.Content),
				}
				if cachedBlocks < 2 {
					content.CacheControl = anthropic.F(anthropic.CacheControlEphemeralParam{
						Type: anthropic.F(anthropic.CacheControlEphemeralType("ephemeral")),
					})
					cachedBlocks++
				}
				blocks = append(blocks, content)
			}

			for _, toolCall := range msg.ToolCalls {
				var inputMap map[string]any
				err := json.Unmarshal([]byte(toolCall.Input), &inputMap)
				if err != nil {
					continue
				}
				toolBlock := anthropic.ToolUseBlockParam{
					Type:  anthropic.F(anthropic.ToolUseBlockParamTypeToolUse),
					ID:    anthropic.F(toolCall.ID),
					Name:  anthropic.F(toolCall.Name),
					Input: anthropic.F(interface{}(inputMap)),
				}
				blocks = append(blocks, toolBlock)
			}

			anthropicMessages[i] = anthropic.MessageParam{
				Role:    anthropic.F(anthropic.MessageParamRoleAssistant),
				Content: anthropic.F(blocks),
			}

		case message.Tool:
			results := make([]anthropic.ContentBlockParamUnion, len(msg.ToolResults))
			for j, toolResult := range msg.ToolResults {
				results[j] = anthropic.ToolResultBlockParam{
					Type:      anthropic.F(anthropic.ToolResultBlockParamTypeToolResult),
					ToolUseID: anthropic.F(toolResult.ToolCallID),
					Content:   anthropic.F([]anthropic.ToolResultBlockParamContentUnion{anthropic.TextBlockParam{Type: anthropic.F(anthropic.TextBlockParamTypeText), Text: anthropic.F(toolResult.Content)}}),
					IsError:   anthropic.F(toolResult.IsError),
				}
			}
			anthropicMessages[i] = anthropic.MessageParam{
				Role:    anthropic.F(anthropic.MessageParamRoleUser),
				Content: anthropic.F(results),
			}
		}
	}

	return anthropicMessages
}