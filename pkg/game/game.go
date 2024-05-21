package game

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/rdor/pkg/dialog"
	"github.com/zrcoder/rdor/pkg/style"
	"github.com/zrcoder/rdor/pkg/style/color"
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
	input          *huh.Input
	showInput      bool
	currentLevel   int
	setLevelAction SetLevelAction
	name           string
	parent         tea.Model
	Err            error
	viewFunc       ViewFunc
	helpFunc       ViewFunc
	keysHelp       help.Model
	keysHelpStyle  lipgloss.Style
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

func (b *Base) SetError(err error) {
	b.Err = err
}

func (b *Base) SetSuccess(msg string) {
	b.showSuccess = true
	b.successMsg = msg
}

func (b *Base) SetStars(total, erned int) {
	b.totalStars = total
	b.ernedStars = erned
}

func (b *Base) SetFailure(msg string) {
	b.showFailure = true
	b.failureMsg = msg
}

func (b *Base) RegisterView(action ViewFunc) {
	b.viewFunc = action
}

func (b *Base) ClearGroups() {
	if b.keyGroups == nil {
		return
	}
	b.keyGroups = b.keyGroups[:0]
}

func (b *Base) AddKeyGroup(group KeyGroup) {
	b.keyGroups = append(b.keyGroups, group)
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

func (b *Base) RegisterHelp(action ViewFunc) {
	b.keyMap.help.SetEnabled(true)
	b.helpFunc = action
}

func (b *Base) Init() tea.Cmd {
	b.keysHelp = help.New()
	b.keysHelp.ShowAll = true
	b.setLevelAction(0)
	b.newInput()
	b.keysHelpStyle = lipgloss.NewStyle().Border(
		lipgloss.ThickBorder()).
		Padding(0, 1).
		BorderForeground(color.Faint)
	return nil
}

func (b *Base) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	orimsg := msg
	switch msg := orimsg.(type) {
	case tea.KeyMsg:
		b.Err = nil
		b.showFailure = false
		b.showSuccess = false
		switch {
		case key.Matches(msg, *b.keyMap.quit):
			return b, tea.Quit
		case key.Matches(msg, *b.keyMap.back):
			return b.parent, nil
		case key.Matches(msg, *b.keyMap.help):
			b.showHelp = !b.showHelp
		case key.Matches(msg, *b.keyMap.reset):
			b.setLevelAction(b.currentLevel)
		case key.Matches(msg, *b.keyMap.next):
			b.currentLevel = (b.currentLevel + 1) % b.levels
			b.setLevelAction(b.currentLevel)
			b.newInput()
		case key.Matches(msg, *b.keyMap.previous):
			b.currentLevel = (b.currentLevel - 1 + b.levels) % b.levels
			b.setLevelAction(b.currentLevel)
			b.newInput()
		case key.Matches(msg, *b.keyMap.setLevel):
			b.showInput = true
			b.input.Placeholder(fmt.Sprintf("1-%d", b.levels))
			b.input.Focus()
			cmd = b.input.Focus()
		case b.showInput && msg.Type == tea.KeyEnter:
			b.showInput = false
			b.input.Blur()
			b.pickLevel(b.input.GetValue().(string))
			b.newInput()
		default:
			if b.showInput {
				_, cmd = b.input.Update(orimsg)
			}
		}
	case tea.WindowSizeMsg:
		b.width = msg.Width
		b.height = msg.Height
		b.keysHelp.Width = msg.Width
	}
	return b, cmd
}

func (b *Base) View() string {
	return lipgloss.NewStyle().Padding(1, 3).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			style.Title.Render(b.name),
			"",
			lipgloss.JoinHorizontal(lipgloss.Top,
				b.mainView(),
				"    ",
				b.keysHelpView(),
			),
		),
	)
}

func (b *Base) pickLevel(s string) {
	n, err := strconv.Atoi(s)
	if err != nil {
		b.Err = err
		return
	}
	if n < 1 || n > b.levels {
		b.Err = fmt.Errorf("the levels must between 1 and %d", b.levels)
		return
	}
	b.setLevelAction(n - 1)
}

func (b *Base) mainView() string {
	if b.showSuccess {
		return dialog.Success(b.successMsg).
			Stars(b.totalStars, b.ernedStars).
			WhiteSpaceChars(b.name).
			Width(b.width).Height(b.height).
			String()
	}

	if b.showFailure {
		return dialog.Error(b.failureMsg).
			WhiteSpaceChars(b.name).
			Width(b.width).Height(b.height).
			String()
	}

	if b.Err != nil {
		return lipgloss.JoinVertical(lipgloss.Left,
			b.viewFunc(),
			style.Error.Render(b.Err.Error()),
		)
	}

	if b.showHelp {
		return lipgloss.JoinVertical(lipgloss.Left,
			b.viewFunc(),
			b.helpFunc(),
		)
	}

	res := b.viewFunc()
	if b.showInput {
		return lipgloss.JoinVertical(lipgloss.Left,
			res,
			"",
			b.input.View(),
		)
	}
	return res
}

func (b *Base) keysHelpView() string {
	groups := append(b.keyGroups, b.keyMap.groups...)
	views := make([]string, 0, 2*len(b.keyGroups)+1)
	for i, group := range groups {
		view := b.keysHelp.View(group)
		if view == "" {
			continue
		}
		views = append(views, view)
		if i < len(groups)-1 {
			views = append(views, "")
		}
	}
	return b.keysHelpStyle.Render(lipgloss.JoinVertical(lipgloss.Left, views...))
}

func (b *Base) newInput() {
	b.input = huh.NewInput().Title("pick a level").Inline(true)
}
