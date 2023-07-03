package game

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
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
	ViewFunc       func() string
	SetLevelAction func(int)
)

type Base struct {
	keyMap         *KeyMap
	keyGroups      []KeyGroup
	levels         int
	input          textinput.Model
	currentLevel   int
	setLevelAction SetLevelAction
	name           string
	parent         tea.Model
	err            error
	viewFunc       ViewFunc
	helpFunc       ViewFunc
	keysHelp       help.Model
	successMsg     string
	failureMsg     string
	width          int
	height         int
	totalStars     int
	ernedStars     int
	showSuccess    bool
	showFailure    bool
	showHelp       bool
}

func New(name string) *Base {
	return &Base{name: name, keyMap: getCommonKeys()}
}

func (b *Base) Name() string {
	return b.name
}

func (b *Base) FilterValue() string {
	return b.name
}

func (b *Base) SetParent(parent tea.Model) {
	b.parent = parent
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

func (g *Base) RegisterView(action ViewFunc) {
	g.viewFunc = action
}

func (g *Base) ClearGroups() {
	if g.keyGroups == nil {
		return
	}
	g.keyGroups = g.keyGroups[:0]
}

func (g *Base) AddKeyGroup(group KeyGroup) {
	g.keyGroups = append(g.keyGroups, group)
}

func (b *Base) RegisterLevels(total int, action SetLevelAction) {
	b.levels = total
	b.setLevelAction = action
	if b.levels > 0 {
		b.keyMap.reset.SetEnabled(true)
		b.keyMap.next.SetEnabled(true)
		b.keyMap.previous.SetEnabled(true)
		b.keyMap.setLevel.SetEnabled(true)
	}
}

func (g *Base) RegisterHelp(action ViewFunc) {
	g.keyMap.help.SetEnabled(true)
	g.helpFunc = action
}

func (g *Base) Init() tea.Cmd {
	g.keysHelp = help.New()
	g.keysHelp.ShowAll = true
	g.input = textinput.New()
	g.setLevelAction(0)
	return nil
}

func (g *Base) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		g.err = nil
		g.showFailure = false
		g.showSuccess = false
		switch {
		case key.Matches(msg, *g.keyMap.quit):
			return g, tea.Quit
		case key.Matches(msg, *g.keyMap.back):
			return g.parent, nil
		case key.Matches(msg, *g.keyMap.help):
			g.showHelp = !g.showHelp
		case key.Matches(msg, *g.keyMap.reset):
			g.setLevelAction(g.currentLevel)
		case key.Matches(msg, *g.keyMap.next):
			g.currentLevel = (g.currentLevel + 1) % g.levels
			g.setLevelAction(g.currentLevel)
		case key.Matches(msg, *g.keyMap.previous):
			g.currentLevel = (g.currentLevel - 1 + g.levels) % g.levels
			g.setLevelAction(g.currentLevel)
		case key.Matches(msg, *g.keyMap.setLevel):
			g.input.Placeholder = fmt.Sprintf("1-%d", g.levels)
			cmd = g.input.Focus()
		default:
			if msg.Type == tea.KeyEnter && g.input.Focused() {
				g.input.Blur()
				g.pickLevel(g.input.Value())
				g.input.SetValue("")
			}
		}
	case tea.WindowSizeMsg:
		g.width = msg.Width
		g.height = msg.Height
		g.keysHelp.Width = msg.Width
	}
	return g, cmd
}

func (g *Base) View() string {
	return lipgloss.NewStyle().Padding(1, 3).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			style.Title.Render(g.name),
			"",
			lipgloss.JoinHorizontal(lipgloss.Top,
				g.mainView(),
				"    ",
				g.keysHelpView(),
			),
		),
	)
}

func (g *Base) pickLevel(s string) {
	n, err := strconv.Atoi(s)
	if err != nil {
		g.err = err
		return
	}
	if n < 1 || n > g.levels {
		g.err = fmt.Errorf("the levels must between 1 and %d", g.levels)
	}
	g.setLevelAction(n - 1)
}

func (g *Base) mainView() string {
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

	if g.err != nil {
		return lipgloss.JoinVertical(lipgloss.Left,
			g.viewFunc(),
			style.Error.Render(g.err.Error()),
		)
	}

	if g.showHelp {
		return lipgloss.JoinVertical(lipgloss.Left,
			g.viewFunc(),
			g.helpFunc(),
		)
	}

	res := g.viewFunc()
	if g.input.Focused() {
		return lipgloss.JoinVertical(lipgloss.Left,
			res,
			"pick a level",
			g.input.View(),
		)
	}
	return res
}

func (g *Base) keysHelpView() string {
	groups := append(g.keyGroups, g.keyMap.groups...)
	views := make([]string, 0, 2*len(g.keyGroups)+1)
	for i, group := range groups {
		view := g.keysHelp.View(group)
		if view == "" {
			continue
		}
		views = append(views, view)
		if i < len(groups)-1 {
			views = append(views, "")
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left, views...)
}
