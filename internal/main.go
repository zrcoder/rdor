package internal

import (
	"fmt"
	"io"

	"github.com/zrcoder/rdor/internal/ballsort"
	"github.com/zrcoder/rdor/internal/crossword"
	"github.com/zrcoder/rdor/internal/hanoi"
	"github.com/zrcoder/rdor/internal/last"
	"github.com/zrcoder/rdor/internal/maze"
	"github.com/zrcoder/rdor/internal/npuzzle"
	"github.com/zrcoder/rdor/internal/point24"
	"github.com/zrcoder/rdor/internal/sokoban"
	"github.com/zrcoder/rdor/pkg/game"
	"github.com/zrcoder/rdor/pkg/style"
	"github.com/zrcoder/rdor/pkg/style/color"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func Run() error {
	const title = "Welcome to rdor"
	items := []list.Item{
		hanoi.New(),
		sokoban.New(),
		maze.New(),
		last.New(),
		npuzzle.New(),
		point24.New(),
		crossword.New(),
		ballsort.New(),
	}
	m := &rdor{
		list: list.New(
			items,
			itemDelegate{},
			// width = screen width, see Update: tea.WindowSizeMsg
			0,
			// height = items(limit 10 every page) + title  + keys help + blank lines
			len(items)%10+7),
	}
	m.list.Title = title
	m.list.Styles.Title = style.Title
	m.list.SetShowStatusBar(false)
	m.list.SetFilteringEnabled(false)
	for _, it := range items {
		it.(game.Game).SetParent(m)
	}
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	return err
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(game.Game)
	if !ok {
		return
	}
	render := lipgloss.NewStyle().PaddingLeft(4).Render
	selectedRender := lipgloss.NewStyle().PaddingLeft(2).Foreground(color.Orange).Render
	if index == m.Index() {
		render = func(s ...string) string {
			if len(s) == 0 {
				return "> "
			}
			return selectedRender("> " + s[0])
		}
	}
	fmt.Fprint(w, render(fmt.Sprintf("%d. %s", index+1, i.Name())))
}

type rdor struct {
	list list.Model
}

func (m rdor) Init() tea.Cmd { return nil }

func (m rdor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "enter" {
			it := m.list.SelectedItem().(game.Game)
			return it, it.Init()
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m rdor) View() string {
	return "\n" + m.list.View()
}
