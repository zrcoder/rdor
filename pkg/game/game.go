package game

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/rdor/pkg/dialog"
	"github.com/zrcoder/rdor/pkg/style"
)

type Game interface {
	Name() string
	SetParent(tea.Model)
	tea.Model
	list.Item
}

type (
	ViewFunc  func() string
	KeyAction func()
)

type Base struct {
	name              string
	ViewFunc          ViewFunc
	HelpFunc          ViewFunc
	KeyActionReset    KeyAction
	KeyActionNext     KeyAction
	KeyActionPrevious KeyAction
	Parent            tea.Model
	CommonKeys        KeyMap
	Keys              []key.Binding

	err         error
	successMsg  string
	failureMsg  string
	showSuccess bool
	showFailure bool
	showHelp    bool
	width       int
	height      int
	totalStars  int
	ernedStars  int
	keysHelp    help.Model
}

func New(name string) *Base {
	return &Base{name: name}
}

func (b *Base) Name() string {
	return b.name
}
func (b *Base) FilterValue() string {
	return b.name
}

func (b *Base) SetParent(parent tea.Model) {
	b.Parent = parent
}

func (g *Base) SetError(err error) {
	g.err = err
}
func (g *Base) SetSuccess(msg string) {
	g.showSuccess = true
	g.successMsg = msg
}
func (g *Base) SetStars(total, erned int) {
	g.totalStars = total
	g.ernedStars = erned
}
func (g *Base) SetFailure(msg string) {
	g.showFailure = true
	g.failureMsg = msg
}

func (g *Base) Init() tea.Cmd {
	g.CommonKeys = getKeys()
	g.keysHelp = help.New()
	g.keysHelp.ShowAll = true
	if g.KeyActionReset == nil {
		g.CommonKeys.Reset.SetEnabled(false)
	}
	if g.KeyActionPrevious == nil {
		g.CommonKeys.Previous.SetEnabled(false)
	}
	if g.KeyActionNext == nil {
		g.CommonKeys.Next.SetEnabled(false)
	}
	if g.HelpFunc == nil {
		g.CommonKeys.Help.SetEnabled(false)
	}
	return nil
}

func (g *Base) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		g.err = nil
		g.showFailure = false
		g.showSuccess = false
		switch {
		case key.Matches(msg, g.CommonKeys.Quit):
			return g, tea.Quit
		case key.Matches(msg, g.CommonKeys.Home):
			return g.Parent, nil
		case key.Matches(msg, g.CommonKeys.Help):
			g.showHelp = !g.showHelp
		case key.Matches(msg, g.CommonKeys.Reset):
			g.KeyActionReset()
		case key.Matches(msg, g.CommonKeys.Previous):
			g.KeyActionPrevious()
		case key.Matches(msg, g.CommonKeys.Next):
			g.KeyActionNext()
		}
	case tea.WindowSizeMsg:
		g.width = msg.Width
		g.height = msg.Height
		g.keysHelp.Width = msg.Width
	}
	return g, nil
}

func (g *Base) View() string {
	return lipgloss.NewStyle().Padding(1, 3).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			style.Title.Render(g.name),
			"\n",
			g.mainView(),
			"\n",
			g.keysHelp.View(g),
		),
	)
}

func (g *Base) mainView() string {
	if g.err != nil {
		return dialog.Error(g.err.Error()).
			WhiteSpaceChars(g.name).
			Width(g.width).Height(g.height).
			String()
	}

	if g.showSuccess {
		return dialog.Success(g.successMsg).
			Stars(g.totalStars, g.ernedStars).
			WhiteSpaceChars(g.name).
			Width(g.width).Height(g.height).
			String()
	}

	if g.showFailure {
		return dialog.Error(g.failureMsg).
			WhiteSpaceChars(g.name).
			Width(g.width).Height(g.height).
			String()
	}

	if g.showHelp {
		return g.HelpFunc()
	}

	return g.ViewFunc()
}
