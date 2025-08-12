package page

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/imnulhaqueruman/opencode-poc/internal/app"
	"github.com/imnulhaqueruman/opencode-poc/internal/tui/components/repl"
	"github.com/imnulhaqueruman/opencode-poc/internal/tui/layout"
)

var ReplPage PageID = "repl"

func NewReplPage(app *app.App) tea.Model {
	return layout.NewBentoLayout(
		layout.BentoPanes{
			layout.BentoLeftPane:        repl.NewSessionsCmp(app),
			layout.BentoRightTopPane:    repl.NewMessagesCmp(app),
			layout.BentoRightBottomPane: repl.NewEditorCmp(app),
		},
		layout.WithBentoLayoutCurrentPane(layout.BentoRightBottomPane),
	)
}
