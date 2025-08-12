package repl

import (
	"errors"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/imnulhaqueruman/opencode-poc/internal/app"
	"github.com/imnulhaqueruman/opencode-poc/internal/llm/agent"
	"github.com/imnulhaqueruman/opencode-poc/internal/tui/layout"
	"github.com/imnulhaqueruman/opencode-poc/internal/tui/styles"
	"github.com/imnulhaqueruman/opencode-poc/internal/tui/util"
	"github.com/kujtimiihoxha/vimtea"
)

type EditorCmp interface {
	tea.Model
	layout.Focusable
	layout.Sizeable
	layout.Bordered
	layout.Bindings
}

type editorCmp struct {
	app        *app.App
	editor     vimtea.Editor
	editorMode vimtea.EditorMode
	sessionID  string
	focused    bool
	width      int
	height     int
}

type editorKeyMap struct {
	SendMessage    key.Binding
	SendMessageI   key.Binding
	InsertMode     key.Binding
	NormaMode      key.Binding
	VisualMode     key.Binding
	VisualLineMode key.Binding
}

var editorKeyMapValue = editorKeyMap{
	SendMessage: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "send message normal mode"),
	),
	SendMessageI: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "send message insert mode"),
	),
	InsertMode: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "insert mode"),
	),
	NormaMode: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "normal mode"),
	),
	VisualMode: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "visual mode"),
	),
	VisualLineMode: key.NewBinding(
		key.WithKeys("V"),
		key.WithHelp("V", "visual line mode"),
	),
}

func (m *editorCmp) Init() tea.Cmd {
	return m.editor.Init()
}

func (m *editorCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case vimtea.EditorModeMsg:
		m.editorMode = msg.Mode
	case SelectedSessionMsg:
		if msg.SessionID != m.sessionID {
			m.sessionID = msg.SessionID
		}
	}
	if m.IsFocused() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, editorKeyMapValue.SendMessage):
				if m.editorMode == vimtea.ModeNormal {
					return m, m.Send()
				}
			case key.Matches(msg, editorKeyMapValue.SendMessageI):
				if m.editorMode == vimtea.ModeInsert {
					return m, m.Send()
				}
			}
		}
		u, cmd := m.editor.Update(msg)
		m.editor = u.(vimtea.Editor)
		return m, cmd
	}
	return m, nil
}

func (m *editorCmp) Blur() tea.Cmd {
	m.focused = false
	return nil
}

func (m *editorCmp) BorderText() map[layout.BorderPosition]string {
	title := "New Message"
	if m.focused {
		title = lipgloss.NewStyle().Foreground(styles.Primary).Render(title)
	}
	return map[layout.BorderPosition]string{
		layout.BottomLeftBorder: title,
	}
}

func (m *editorCmp) Focus() tea.Cmd {
	m.focused = true
	return nil
}

func (m *editorCmp) GetSize() (int, int) {
	return m.width, m.height
}

func (m *editorCmp) IsFocused() bool {
	return m.focused
}

func (m *editorCmp) SetSize(width int, height int) {
	m.width = width
	m.height = height
	m.editor.SetSize(width, height)
}

func (m *editorCmp) Send() tea.Cmd {
	return func() tea.Msg {
		if m.sessionID == "" {
			return util.ErrorMsg(errors.New("No session selected"))
		}
		
		messages, _ := m.app.Messages.List(m.sessionID)
		if hasUnfinishedMessages(messages) {
			return util.InfoMsg("Assistant is still working on the previous message")
		}
		a, err := agent.NewCoderAgent(m.app)
		if err != nil {
			return util.ErrorMsg(err)
		}

		if m.editor == nil {
			return util.ErrorMsg(errors.New("Editor is not initialized"))
		}
		
		buffer := m.editor.GetBuffer()
		if buffer == nil {
			return util.ErrorMsg(errors.New("Editor buffer is not available"))
		}
		
		content := strings.Join(buffer.Lines(), "\n")
		go a.Generate(m.sessionID, content)

		m.editor = vimtea.NewEditor(
			vimtea.WithFileName("message.md"),
		)
		return nil
	}
}

func (m *editorCmp) View() string {
	return m.editor.View()
}

func (m *editorCmp) BindingKeys() []key.Binding {
	return layout.KeyMapToSlice(editorKeyMapValue)
}

func NewEditorCmp(app *app.App) EditorCmp {
	editor := vimtea.NewEditor(
		vimtea.WithFileName("message.md"),
	)
	return &editorCmp{
		app:    app,
		editor: editor,
	}
}
