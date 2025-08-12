package repl

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/imnulhaqueruman/opencode-poc/internal/app"
	"github.com/imnulhaqueruman/opencode-poc/internal/llm/agent"
	"github.com/imnulhaqueruman/opencode-poc/internal/message"
	"github.com/imnulhaqueruman/opencode-poc/internal/pubsub"
	"github.com/imnulhaqueruman/opencode-poc/internal/session"
	"github.com/imnulhaqueruman/opencode-poc/internal/tui/layout"
	"github.com/imnulhaqueruman/opencode-poc/internal/tui/styles"
)

type MessagesCmp interface {
	tea.Model
	layout.Focusable
	layout.Bordered
	layout.Sizeable
	layout.Bindings
}

type messagesCmp struct {
	app            *app.App
	messages       []message.Message
	selectedMsgIdx int // Index of the selected message
	session        session.Session
	viewport       viewport.Model
	mdRenderer     *glamour.TermRenderer
	width          int
	height         int
	focused        bool
	cachedView     string
}

func (m *messagesCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case pubsub.Event[message.Message]:
		if msg.Type == pubsub.CreatedEvent {
			if msg.Payload.SessionID == m.session.ID {
				m.messages = append(m.messages, msg.Payload)
				m.renderView()
				if m.viewport.Width > 0 && m.viewport.Height > 0 {
					m.viewport.GotoBottom()
				}
			}
			for _, v := range m.messages {
				for _, c := range v.ToolCalls {
					if c.ID == msg.Payload.SessionID {
						m.renderView()
						if m.viewport.Width > 0 && m.viewport.Height > 0 {
							m.viewport.GotoBottom()
						}
					}
				}
			}
		} else if msg.Type == pubsub.UpdatedEvent && msg.Payload.SessionID == m.session.ID {
			for i, v := range m.messages {
				if v.ID == msg.Payload.ID {
					m.messages[i] = msg.Payload
					m.renderView()
					if i == len(m.messages)-1 && m.viewport.Width > 0 && m.viewport.Height > 0 {
						m.viewport.GotoBottom()
					}
					break
				}
			}
		}
	case pubsub.Event[session.Session]:
		if msg.Type == pubsub.UpdatedEvent && m.session.ID == msg.Payload.ID {
			m.session = msg.Payload
		}
	case SelectedSessionMsg:
		m.session, _ = m.app.Sessions.Get(msg.SessionID)
		m.messages, _ = m.app.Messages.List(m.session.ID)
		m.renderView()
		// Only go to bottom if viewport has valid dimensions
		if m.viewport.Width > 0 && m.viewport.Height > 0 {
			m.viewport.GotoBottom()
		}
	}
	if m.focused {
		u, cmd := m.viewport.Update(msg)
		m.viewport = u
		return m, cmd
	}
	return m, nil
}

func borderColor(role message.MessageRole) lipgloss.TerminalColor {
	switch role {
	case message.Assistant:
		return styles.Mauve
	case message.User:
		return styles.Rosewater
	}
	return styles.Blue
}

func borderText(msgRole message.MessageRole, currentMessage int) map[layout.BorderPosition]string {
	role := ""
	icon := ""
	switch msgRole {
	case message.Assistant:
		role = "Assistant"
		icon = styles.BotIcon
	case message.User:
		role = "User"
		icon = styles.UserIcon
	}
	return map[layout.BorderPosition]string{
		layout.TopLeftBorder: lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true).
			Foreground(styles.Crust).
			Background(borderColor(msgRole)).
			Render(fmt.Sprintf("%s %s ", role, icon)),
		layout.TopRightBorder: lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true).
			Foreground(styles.Crust).
			Background(borderColor(msgRole)).
			Render(fmt.Sprintf("#%d ", currentMessage)),
	}
}

func hasUnfinishedMessages(messages []message.Message) bool {
	if len(messages) == 0 {
		return false
	}
	for _, msg := range messages {
		if !msg.Finished {
			return true
		}
	}
	lastMessage := messages[len(messages)-1]
	return lastMessage.Role != message.Assistant
}

func (m *messagesCmp) renderMessageWithToolCall(content string, tools []message.ToolCall, futureMessages []message.Message) string {
	allParts := []string{content}

	leftPaddingValue := 4
	connectorStyle := lipgloss.NewStyle().
		Foreground(styles.Peach).
		Bold(true)

	toolCallStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Peach).
		Width(m.width-leftPaddingValue-5).
		Padding(0, 1)

	toolResultStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Green).
		Width(m.width-leftPaddingValue-5).
		Padding(0, 1)

	leftPadding := lipgloss.NewStyle().Padding(0, 0, 0, leftPaddingValue)

	runningStyle := lipgloss.NewStyle().
		Foreground(styles.Peach).
		Bold(true)

	renderTool := func(toolCall message.ToolCall) string {
		toolHeader := lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.Blue).
			Render(fmt.Sprintf("%s %s", styles.ToolIcon, toolCall.Name))

		var paramLines []string
		var args map[string]interface{}
		var paramOrder []string

		json.Unmarshal([]byte(toolCall.Input), &args)

		for key := range args {
			paramOrder = append(paramOrder, key)
		}
		sort.Strings(paramOrder)

		for _, name := range paramOrder {
			value := args[name]
			paramName := lipgloss.NewStyle().
				Foreground(styles.Peach).
				Bold(true).
				Render(name)

			truncate := m.width - leftPaddingValue*2 - 10
			if len(fmt.Sprintf("%v", value)) > truncate {
				value = fmt.Sprintf("%v", value)[:truncate] + lipgloss.NewStyle().Foreground(styles.Blue).Render("... (truncated)")
			}
			paramValue := fmt.Sprintf("%v", value)
			paramLines = append(paramLines, fmt.Sprintf("  %s: %s", paramName, paramValue))
		}

		paramBlock := lipgloss.JoinVertical(lipgloss.Left, paramLines...)

		toolContent := lipgloss.JoinVertical(lipgloss.Left, toolHeader, paramBlock)
		return toolCallStyle.Render(toolContent)
	}

	findToolResult := func(toolCallID string, messages []message.Message) *message.ToolResult {
		for _, msg := range messages {
			if msg.Role == message.Tool {
				for _, result := range msg.ToolResults {
					if result.ToolCallID == toolCallID {
						return &result
					}
				}
			}
		}
		return nil
	}

	renderToolResult := func(result message.ToolResult) string {
		resultHeader := lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.Green).
			Render(fmt.Sprintf("%s Result", styles.CheckIcon))
		if result.IsError {
			resultHeader = lipgloss.NewStyle().
				Bold(true).
				Foreground(styles.Red).
				Render(fmt.Sprintf("%s Error", styles.ErrorIcon))
		}

		truncate := 200
		content := result.Content
		if len(content) > truncate {
			content = content[:truncate] + lipgloss.NewStyle().Foreground(styles.Blue).Render("... (truncated)")
		}

		resultContent := lipgloss.JoinVertical(lipgloss.Left, resultHeader, content)
		return toolResultStyle.Render(resultContent)
	}

	connector := connectorStyle.Render("└─> Tool Calls:")
	allParts = append(allParts, connector)

	for _, toolCall := range tools {
		toolOutput := renderTool(toolCall)
		allParts = append(allParts, leftPadding.Render(toolOutput))

		result := findToolResult(toolCall.ID, futureMessages)
		if result != nil {

			resultOutput := renderToolResult(*result)
			allParts = append(allParts, leftPadding.Render(resultOutput))

		} else if toolCall.Name == agent.AgentToolName {

			runningIndicator := runningStyle.Render(fmt.Sprintf("%s Running...", styles.SpinnerIcon))
			allParts = append(allParts, leftPadding.Render(runningIndicator))
			taskSessionMessages, _ := m.app.Messages.List(toolCall.ID)
			for _, msg := range taskSessionMessages {
				if msg.Role == message.Assistant {
					for _, toolCall := range msg.ToolCalls {
						toolHeader := lipgloss.NewStyle().
							Bold(true).
							Foreground(styles.Blue).
							Render(fmt.Sprintf("%s %s", styles.ToolIcon, toolCall.Name))

						var paramLines []string
						var args map[string]interface{}
						var paramOrder []string

						json.Unmarshal([]byte(toolCall.Input), &args)

						for key := range args {
							paramOrder = append(paramOrder, key)
						}
						sort.Strings(paramOrder)

						for _, name := range paramOrder {
							value := args[name]
							paramName := lipgloss.NewStyle().
								Foreground(styles.Peach).
								Bold(true).
								Render(name)

							truncate := 50
							if len(fmt.Sprintf("%v", value)) > truncate {
								value = fmt.Sprintf("%v", value)[:truncate] + lipgloss.NewStyle().Foreground(styles.Blue).Render("... (truncated)")
							}
							paramValue := fmt.Sprintf("%v", value)
							paramLines = append(paramLines, fmt.Sprintf("  %s: %s", paramName, paramValue))
						}

						paramBlock := lipgloss.JoinVertical(lipgloss.Left, paramLines...)
						toolContent := lipgloss.JoinVertical(lipgloss.Left, toolHeader, paramBlock)
						toolOutput := toolCallStyle.BorderForeground(styles.Teal).MaxWidth(m.width - leftPaddingValue*2 - 2).Render(toolContent)
						allParts = append(allParts, lipgloss.NewStyle().Padding(0, 0, 0, leftPaddingValue*2).Render(toolOutput))
					}
				}
			}

		} else {
			runningIndicator := runningStyle.Render(fmt.Sprintf("%s Running...", styles.SpinnerIcon))
			allParts = append(allParts, "    "+runningIndicator)
		}
	}

	for _, msg := range futureMessages {
		if msg.Content != "" {
			break
		}

		for _, toolCall := range msg.ToolCalls {
			toolOutput := renderTool(toolCall)
			allParts = append(allParts, "    "+strings.ReplaceAll(toolOutput, "\n", "\n    "))

			result := findToolResult(toolCall.ID, futureMessages)
			if result != nil {
				resultOutput := renderToolResult(*result)
				allParts = append(allParts, "    "+strings.ReplaceAll(resultOutput, "\n", "\n    "))
			} else {
				runningIndicator := runningStyle.Render(fmt.Sprintf("%s Running...", styles.SpinnerIcon))
				allParts = append(allParts, "    "+runningIndicator)
			}
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, allParts...)
}

func (m *messagesCmp) renderView() {
	stringMessages := make([]string, 0)
	r, _ := glamour.NewTermRenderer(
		glamour.WithStyles(styles.CatppuccinMarkdownStyle()),
		glamour.WithWordWrap(m.width-20),
		glamour.WithEmoji(),
	)
	textStyle := lipgloss.NewStyle().Width(m.width - 4)
	currentMessage := 1
	displayedMsgCount := 0 // Track the actual displayed messages count

	prevMessageWasUser := false
	for inx, msg := range m.messages {
		content := msg.Content
		if content != "" || prevMessageWasUser {
			if msg.Thinking != "" && content == "" {
				content = msg.Thinking
			} else if content == "" {
				content = "..."
			}
			content, _ = r.Render(content)

			isSelected := inx == m.selectedMsgIdx

			border := lipgloss.DoubleBorder()
			activeColor := borderColor(msg.Role)

			if isSelected {
				activeColor = styles.Primary // Use primary color for selected message
			}

			content = layout.Borderize(
				textStyle.Render(content),
				layout.BorderOptions{
					InactiveBorder: border,
					ActiveBorder:   border,
					ActiveColor:    activeColor,
					InactiveColor:  borderColor(msg.Role),
					EmbeddedText:   borderText(msg.Role, currentMessage),
				},
			)
			if len(msg.ToolCalls) > 0 {
				content = m.renderMessageWithToolCall(content, msg.ToolCalls, m.messages[inx+1:])
			}
			stringMessages = append(stringMessages, content)
			currentMessage++
			displayedMsgCount++
		}
		if msg.Role == message.User && msg.Content != "" {
			prevMessageWasUser = true
		} else {
			prevMessageWasUser = false
		}
	}
	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Top, stringMessages...))
}

func (m *messagesCmp) View() string {
	return lipgloss.NewStyle().Padding(1).Render(m.viewport.View())
}

func (m *messagesCmp) BindingKeys() []key.Binding {
	keys := layout.KeyMapToSlice(m.viewport.KeyMap)

	return keys
}

func (m *messagesCmp) Blur() tea.Cmd {
	m.focused = false
	return nil
}

func (m *messagesCmp) BorderText() map[layout.BorderPosition]string {
	title := m.session.Title
	titleWidth := m.width / 2
	if len(title) > titleWidth {
		title = title[:titleWidth] + "..."
	}
	if m.focused {
		title = lipgloss.NewStyle().Foreground(styles.Primary).Render(title)
	}
	borderTest := map[layout.BorderPosition]string{
		layout.TopLeftBorder:     title,
		layout.BottomRightBorder: formatTokensAndCost(m.session.CompletionTokens+m.session.PromptTokens, m.session.Cost),
	}
	if hasUnfinishedMessages(m.messages) {
		borderTest[layout.BottomLeftBorder] = lipgloss.NewStyle().Foreground(styles.Peach).Render("Thinking...")
	} else {
		borderTest[layout.BottomLeftBorder] = lipgloss.NewStyle().Foreground(styles.Text).Render("Sleeping " + styles.SleepIcon + " ")
	}

	return borderTest
}

func (m *messagesCmp) Focus() tea.Cmd {
	m.focused = true
	return nil
}

func (m *messagesCmp) GetSize() (int, int) {
	return m.width, m.height
}

func (m *messagesCmp) IsFocused() bool {
	return m.focused
}

func (m *messagesCmp) SetSize(width int, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = width - 2   // padding
	m.viewport.Height = height - 2 // padding
	m.renderView()
}

func (m *messagesCmp) Init() tea.Cmd {
	return nil
}

func NewMessagesCmp(app *app.App) MessagesCmp {
	return &messagesCmp{
		app:      app,
		messages: []message.Message{},
		viewport: viewport.New(0, 0),
	}
}
