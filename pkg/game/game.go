package game

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/rdor/pkg/dialog"
	"github.com/zrcoder/rdor/pkg/style"
)

type (
	InitFunc   func() tea.Cmd
	UpdateFunc func(tea.Msg) tea.Cmd
	ViewFunc   func() string
	KeyAction  func()
)

type Game struct {
	Title             string
	ViewFunc          ViewFunc
	HelpFunc          ViewFunc
	InitFunc          InitFunc
	UpdateFunc        UpdateFunc
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

func New(title string) *Game {
	game := &Game{Title: title,
		CommonKeys: Keys,
		keysHelp:   help.New(),
	}
	game.keysHelp.ShowAll = true
	return game
}

func (g *Game) SetError(err error) {
	g.err = err
}
func (g *Game) SetSuccess(msg string) {
	g.showSuccess = true
	g.successMsg = msg
}
func (g *Game) SetStars(total, erned int) {
	g.totalStars = total
	g.ernedStars = erned
}
func (g *Game) SetFailure(msg string) {
	g.showFailure = true
	g.failureMsg = msg
}

func (g *Game) Init() tea.Cmd {
	var res tea.Cmd
	if g.InitFunc != nil {
		res = g.InitFunc()
	}
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
	return res
}

func (g *Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	return g, g.UpdateFunc(msg)
}

func (g *Game) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		style.Title.Render(g.Title),
		"\n",
		g.mainView(),
		"\n",
		g.keysHelp.View(g),
	)
}

func (g *Game) mainView() string {
	if g.err != nil {
		return dialog.Error(g.err.Error()).
			WhiteSpaceChars(g.Title).
			Width(g.width).Height(g.height).
			String()
	}

	if g.showSuccess {
		return dialog.Success(g.successMsg).
			Stars(g.totalStars, g.ernedStars).
			WhiteSpaceChars(g.Title).
			Width(g.width).Height(g.height).
			String()
	}

	if g.showFailure {
		return dialog.Error(g.failureMsg).
			WhiteSpaceChars(g.Title).
			Width(g.width).Height(g.height).
			String()
	}

	if g.showHelp {
		return g.HelpFunc()
	}

	return g.ViewFunc()
}
