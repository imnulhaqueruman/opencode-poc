package agent

import (
	"errors"

	"github.com/imnulhaqueruman/opencode-poc/internal/app"
	"github.com/imnulhaqueruman/opencode-poc/internal/config"
	"github.com/imnulhaqueruman/opencode-poc/internal/llm/models"
	"github.com/imnulhaqueruman/opencode-poc/internal/llm/tools"
)

type taskAgent struct {
	*agent
}

func (c *taskAgent) Generate(sessionID string, content string) error {
	return c.generate(sessionID, content)
}

func NewTaskAgent(app *app.App) (Agent, error) {
	model, ok := models.SupportedModels[config.Get().Model.Coder]
	if !ok {
		return nil, errors.New("model not supported")
	}

	agentProvider, titleGenerator, err := getAgentProviders(app.Context, model)
	if err != nil {
		return nil, err
	}
	return &taskAgent{
		agent: &agent{
			App: app,
			tools: []tools.BaseTool{
				tools.NewGlobTool(),
				tools.NewGrepTool(),
				tools.NewLsTool(),
				tools.NewViewTool(),
			},
			model:          model,
			agent:          agentProvider,
			titleGenerator: titleGenerator,
		},
	}, nil
}