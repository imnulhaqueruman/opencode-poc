package core

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/imnulhaqueruman/opencode-poc/internal/config"
	"github.com/imnulhaqueruman/opencode-poc/internal/llm/models"
	"github.com/imnulhaqueruman/opencode-poc/internal/tui/styles"
	"github.com/imnulhaqueruman/opencode-poc/internal/tui/util"
	"github.com/imnulhaqueruman/opencode-poc/internal/version"
)

type statusCmp struct {
	err   error
	info  string
	width int
}

func (m statusCmp) Init() tea.Cmd {
	return nil
}

func (m statusCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case util.ErrorMsg:
		m.err = msg
	case util.InfoMsg:
		m.info = string(msg)
	}
	return m, nil
}

var (
	versionWidget = styles.Padded.Background(styles.DarkGrey).Foreground(styles.Text).Render(version.Version)
	helpWidget    = styles.Padded.Background(styles.Grey).Foreground(styles.Text).Render("? help")
)

func (m statusCmp) View() string {
	status := styles.Padded.Background(styles.Grey).Foreground(styles.Text).Render("? help")

	if m.err != nil {
		status += styles.Regular.Padding(0, 1).
			Background(styles.Red).
			Foreground(styles.Text).
			Width(m.availableFooterMsgWidth()).
			Render(m.err.Error())
	} else if m.info != "" {
		status += styles.Padded.
			Foreground(styles.Base).
			Background(styles.Green).
			Width(m.availableFooterMsgWidth()).
			Render(m.info)
	} else {
		status += styles.Padded.
			Foreground(styles.Base).
			Background(styles.LightGrey).
			Width(m.availableFooterMsgWidth()).
			Render(m.info)
	}
	status += m.model()
	status += versionWidget
	return status
}

func (m statusCmp) availableFooterMsgWidth() int {
	// -2 to accommodate padding
	return max(0, m.width-lipgloss.Width(helpWidget)-lipgloss.Width(versionWidget)-lipgloss.Width(m.model()))
}

func (m statusCmp) model() string {
	cfg := config.Get()
	if cfg.Model == nil {
		return styles.Padded.Background(styles.Grey).Foreground(styles.Text).Render("No Model")
	}
	model, exists := models.SupportedModels[cfg.Model.Coder]
	if !exists {
		return styles.Padded.Background(styles.Grey).Foreground(styles.Text).Render("Unknown Model")
	}
	return styles.Padded.Background(styles.Grey).Foreground(styles.Text).Render(model.Name)
}

func NewStatusCmp() tea.Model {
	return &statusCmp{}
}