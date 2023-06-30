package internal

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
	KeyFunc    func()
)

type Game struct {
	Title           string
	ViewFunc        ViewFunc
	InitFunc        InitFunc
	UpdateFunc      UpdateFunc
	KeyFuncReset    KeyFunc
	KeyFuncNext     KeyFunc
	KeyFuncPrevious KeyFunc
	Parent          tea.Model

	Keys        KeyMap
	extroKeys   []key.Binding
	err         error
	successMsg  string
	failureView string
	showSuccess bool
	showFailure bool
	showHelp    bool
	width       int
	height      int
	keysHelp    help.Model
	totalStars  int
	ernedStars  int
}

func New(title string) *Game {
	g := &Game{Title: title,
		Keys:     Keys,
		keysHelp: help.New(),
	}
	g.keysHelp.ShowAll = true
	return g
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
	g.failureView = msg
}

func (g *Game) Init() tea.Cmd {
	res := g.InitFunc()
	if g.KeyFuncReset == nil {
		g.Keys.Reset.SetEnabled(false)
	}
	if g.KeyFuncPrevious == nil {
		g.Keys.Previous.SetEnabled(false)
	}
	if g.KeyFuncNext == nil {
		g.Keys.Next.SetEnabled(false)
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
		case key.Matches(msg, g.Keys.Quit):
			return g, tea.Quit
		case key.Matches(msg, g.Keys.Home):
			return g.Parent, nil
		case key.Matches(msg, g.Keys.Help):
			g.showHelp = !g.showHelp
		case key.Matches(msg, g.Keys.Reset):
			g.KeyFuncReset()
		case key.Matches(msg, g.Keys.Previous):
			g.KeyFuncPrevious()
		case key.Matches(msg, g.Keys.Next):
			g.KeyFuncNext()
		}
	case tea.WindowSizeMsg:
		g.width = msg.Width
		g.height = msg.Height
		g.keysHelp.Width = msg.Width
	}
	return g, g.UpdateFunc(msg)
}
func (g *Game) View() string {
	head := style.Title.Render(g.Title)

	if g.err != nil {
		ed := dialog.Error(g.err.Error()).WhiteSpaceChars(g.Title).Width(g.width).Height(g.height)
		return lipgloss.JoinVertical(lipgloss.Left, head, ed.String(), g.keysHelp.View(g))
	}

	if g.showSuccess {
		sd := dialog.Success(g.successMsg).WhiteSpaceChars(g.Title).Width(g.width).Height(g.height).Stars(g.totalStars, g.ernedStars)
		return lipgloss.JoinVertical(lipgloss.Left, head, sd.String(), g.keysHelp.View(g))
	}

	if g.showFailure {
		fd := dialog.Error(g.failureView).WhiteSpaceChars(g.Title).Width(g.width).Height(g.height)
		return lipgloss.JoinVertical(lipgloss.Left, head, fd.String(), g.keysHelp.View(g))
	}

	return lipgloss.JoinVertical(lipgloss.Left, head, g.ViewFunc(), g.keysHelp.View(g))
}

func (g *Game) SetExtraKeys(keys []key.Binding) {
	g.extroKeys = keys
}
