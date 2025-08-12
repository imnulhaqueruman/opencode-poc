package agent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/imnulhaqueruman/opencode-poc/internal/app"
	"github.com/imnulhaqueruman/opencode-poc/internal/config"
	"github.com/imnulhaqueruman/opencode-poc/internal/llm/models"
	"github.com/imnulhaqueruman/opencode-poc/internal/llm/prompt"
	"github.com/imnulhaqueruman/opencode-poc/internal/llm/provider"
	"github.com/imnulhaqueruman/opencode-poc/internal/llm/tools"
	"github.com/imnulhaqueruman/opencode-poc/internal/message"
)

type Agent interface {
	Generate(sessionID string, content string) error
}

type agent struct {
	*app.App
	model          models.Model
	tools          []tools.BaseTool
	agent          provider.Provider
	titleGenerator provider.Provider
}

func (c *agent) handleTitleGeneration(sessionID, content string) {
	response, err := c.titleGenerator.SendMessages(
		c.Context,
		[]message.Message{
			{
				Role:    message.User,
				Content: content,
			},
		},
		nil,
	)
	if err != nil {
		return
	}

	session, err := c.Sessions.Get(sessionID)
	if err != nil {
		return
	}
	if response.Content != "" {
		session.Title = response.Content
		c.Sessions.Save(session)
	}
}

func (c *agent) TrackUsage(sessionID string, model models.Model, usage provider.TokenUsage) error {
	session, err := c.Sessions.Get(sessionID)
	if err != nil {
		return err
	}

	cost := model.CostPer1MInCached/1e6*float64(usage.CacheCreationTokens) +
		model.CostPer1MOutCached/1e6*float64(usage.CacheReadTokens) +
		model.CostPer1MIn/1e6*float64(usage.InputTokens) +
		model.CostPer1MOut/1e6*float64(usage.OutputTokens)

	session.Cost += cost
	session.CompletionTokens += usage.OutputTokens
	session.PromptTokens += usage.InputTokens

	_, err = c.Sessions.Save(session)
	return err
}

func (c *agent) processEvent(
	sessionID string,
	assistantMsg *message.Message,
	event provider.ProviderEvent,
) error {
	switch event.Type {
	case provider.EventThinkingDelta:
		assistantMsg.Thinking += event.Thinking
		return c.Messages.Update(*assistantMsg)
	case provider.EventContentDelta:
		assistantMsg.Content += event.Content
		return c.Messages.Update(*assistantMsg)
	case provider.EventError:
		log.Println("error", event.Error)
		return event.Error

	case provider.EventComplete:
		assistantMsg.ToolCalls = event.Response.ToolCalls
		err := c.Messages.Update(*assistantMsg)
		if err != nil {
			return err
		}
		return c.TrackUsage(sessionID, c.model, event.Response.Usage)
	}

	return nil
}

func (c *agent) ExecuteTools(ctx context.Context, toolCalls []message.ToolCall, tls []tools.BaseTool) ([]message.ToolResult, error) {
	var wg sync.WaitGroup
	toolResults := make([]message.ToolResult, len(toolCalls))
	mutex := &sync.Mutex{}

	for i, tc := range toolCalls {
		wg.Add(1)
		go func(index int, toolCall message.ToolCall) {
			defer wg.Done()

			response := ""
			isError := false
			found := false

			for _, tool := range tls {
				if tool.Info().Name == toolCall.Name {
					found = true
					toolResult, toolErr := tool.Run(ctx, tools.ToolCall{
						ID:    toolCall.ID,
						Name:  toolCall.Name,
						Input: toolCall.Input,
					})
					if toolErr != nil {
						response = fmt.Sprintf("error running tool: %s", toolErr)
						isError = true
					} else {
						response = toolResult.Content
						isError = toolResult.IsError
					}
					break
				}
			}

			if !found {
				response = fmt.Sprintf("tool not found: %s", toolCall.Name)
				isError = true
			}

			mutex.Lock()
			defer mutex.Unlock()

			toolResults[index] = message.ToolResult{
				ToolCallID: toolCall.ID,
				Content:    response,
				IsError:    isError,
			}
		}(i, tc)
	}

	wg.Wait()
	return toolResults, nil
}

func (c *agent) handleToolExecution(
	ctx context.Context,
	assistantMsg message.Message,
) (*message.Message, error) {
	if len(assistantMsg.ToolCalls) == 0 {
		return nil, nil
	}

	toolResults, err := c.ExecuteTools(ctx, assistantMsg.ToolCalls, c.tools)
	if err != nil {
		return nil, err
	}

	msg, err := c.Messages.Create(assistantMsg.SessionID, message.CreateMessageParams{
		Role:        message.Tool,
		ToolResults: toolResults,
	})

	return &msg, err
}

func (c *agent) generate(sessionID string, content string) error {
	messages, err := c.Messages.List(sessionID)
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		go c.handleTitleGeneration(sessionID, content)
	}

	userMsg, err := c.Messages.Create(sessionID, message.CreateMessageParams{
		Role:    message.User,
		Content: content,
	})
	if err != nil {
		return err
	}

	messages = append(messages, userMsg)
	for {

		eventChan, err := c.agent.StreamResponse(c.Context, messages, c.tools)
		if err != nil {
			return err
		}

		assistantMsg, err := c.Messages.Create(sessionID, message.CreateMessageParams{
			Role:    message.Assistant,
			Content: "",
		})
		if err != nil {
			return err
		}
		for event := range eventChan {
			err = c.processEvent(sessionID, &assistantMsg, event)
			if err != nil {
				assistantMsg.Finished = true
				c.Messages.Update(assistantMsg)
				return err
			}
		}

		msg, err := c.handleToolExecution(c.Context, assistantMsg)
		assistantMsg.Finished = true
		c.Messages.Update(assistantMsg)
		if err != nil {
			return err
		}

		if len(assistantMsg.ToolCalls) == 0 {
			break
		}

		messages = append(messages, assistantMsg)
		if msg != nil {
			messages = append(messages, *msg)
		}
	}
	return nil
}

func getAgentProviders(ctx context.Context, model models.Model) (provider.Provider, provider.Provider, error) {
	maxTokens := config.Get().Model.CoderMaxTokens

	providerConfig, ok := config.Get().Providers[model.Provider]
	if !ok || !providerConfig.Enabled {
		return nil, nil, errors.New("provider is not enabled")
	}
	var agentProvider provider.Provider
	var titleGenerator provider.Provider

	switch model.Provider {
	case models.ProviderOpenAI:
		var err error
		agentProvider, err = provider.NewOpenAIProvider(
			provider.WithOpenAISystemMessage(
				prompt.CoderOpenAISystemPrompt(),
			),
			provider.WithOpenAIMaxTokens(maxTokens),
			provider.WithOpenAIModel(model),
			provider.WithOpenAIKey(providerConfig.APIKey),
		)
		if err != nil {
			return nil, nil, err
		}
		titleGenerator, err = provider.NewOpenAIProvider(
			provider.WithOpenAISystemMessage(
				prompt.TitlePrompt(),
			),
			provider.WithOpenAIMaxTokens(80),
			provider.WithOpenAIModel(model),
			provider.WithOpenAIKey(providerConfig.APIKey),
		)
		if err != nil {
			return nil, nil, err
		}
	case models.ProviderAnthropic:
		var err error
		agentProvider, err = provider.NewAnthropicProvider(
			provider.WithAnthropicSystemMessage(
				prompt.CoderAnthropicSystemPrompt(),
			),
			provider.WithAnthropicMaxTokens(maxTokens),
			provider.WithAnthropicKey(providerConfig.APIKey),
			provider.WithAnthropicModel(model),
		)
		if err != nil {
			return nil, nil, err
		}
		titleGenerator, err = provider.NewAnthropicProvider(
			provider.WithAnthropicSystemMessage(
				prompt.TitlePrompt(),
			),
			provider.WithAnthropicMaxTokens(80),
			provider.WithAnthropicKey(providerConfig.APIKey),
			provider.WithAnthropicModel(model),
		)
		if err != nil {
			return nil, nil, err
		}

	case models.ProviderGemini:
		var err error
		agentProvider, err = provider.NewGeminiProvider(
			ctx,
			provider.WithGeminiSystemMessage(
				prompt.CoderOpenAISystemPrompt(),
			),
			provider.WithGeminiMaxTokens(int32(maxTokens)),
			provider.WithGeminiKey(providerConfig.APIKey),
			provider.WithGeminiModel(model),
		)
		if err != nil {
			return nil, nil, err
		}
		titleGenerator, err = provider.NewGeminiProvider(
			ctx,
			provider.WithGeminiSystemMessage(
				prompt.TitlePrompt(),
			),
			provider.WithGeminiMaxTokens(80),
			provider.WithGeminiKey(providerConfig.APIKey),
			provider.WithGeminiModel(model),
		)
		if err != nil {
			return nil, nil, err
		}
	case models.ProviderGROQ:
		var err error
		agentProvider, err = provider.NewOpenAIProvider(
			provider.WithOpenAISystemMessage(
				prompt.CoderAnthropicSystemPrompt(),
			),
			provider.WithOpenAIMaxTokens(maxTokens),
			provider.WithOpenAIModel(model),
			provider.WithOpenAIKey(providerConfig.APIKey),
			provider.WithOpenAIBaseURL("https://api.groq.com/openai/v1"),
		)
		if err != nil {
			return nil, nil, err
		}
		titleGenerator, err = provider.NewOpenAIProvider(
			provider.WithOpenAISystemMessage(
				prompt.TitlePrompt(),
			),
			provider.WithOpenAIMaxTokens(80),
			provider.WithOpenAIModel(model),
			provider.WithOpenAIKey(providerConfig.APIKey),
			provider.WithOpenAIBaseURL("https://api.groq.com/openai/v1"),
		)
		if err != nil {
			return nil, nil, err
		}

	}

	return agentProvider, titleGenerator, nil
}