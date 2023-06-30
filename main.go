package main

import (
	"fmt"
	"io"
	"os"

	"github.com/zrcoder/rdor/internal/hanoi"
	"github.com/zrcoder/rdor/internal/last"
	"github.com/zrcoder/rdor/internal/maze"
	"github.com/zrcoder/rdor/internal/npuzzle"
	"github.com/zrcoder/rdor/internal/sokoban"
	"github.com/zrcoder/rdor/pkg/model"
	"github.com/zrcoder/rdor/pkg/style"
	"github.com/zrcoder/rdor/pkg/style/color"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//go:generate go run ./internal/maze/levels/tool
func main() {
	items := []list.Item{
		item{name: hanoi.Name, game: hanoi.New()},
		item{name: sokoban.Name, game: sokoban.New()},
		item{name: maze.Name, game: maze.New()},
		item{name: npuzzle.Name, game: npuzzle.New()},
		item{name: last.Name, game: last.New()},
	}
	const listHeight = 14
	const defaultWidth = 20
	m := rdor{list: list.New(items, itemDelegate{}, defaultWidth, listHeight)}
	m.list.Title = "Welcome to `rdor`"
	m.list.Styles.Title = style.Title
	m.list.SetShowStatusBar(false)
	m.list.SetFilteringEnabled(false)
	for _, it := range items {
		it.(item).game.SetParent(m)
	}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

type item struct {
	name string
	game model.Game
}

func (i item) FilterValue() string { return i.name }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

var (
	itemRender         = lipgloss.NewStyle().PaddingLeft(4).Render
	itemSelectedRender = lipgloss.NewStyle().PaddingLeft(2).Foreground(color.Orange).Render
)

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}
	fn := itemRender
	if index == m.Index() {
		fn = func(s ...string) string {
			if len(s) == 0 {
				return "> "
			}
			return itemSelectedRender("> " + s[0])
		}
	}
	fmt.Fprint(w, fn(fmt.Sprintf("%d. %s", index+1, i.name)))
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
		switch val := msg.String(); val {
		case "enter":
			i := m.list.SelectedItem().(item)
			return i.game, i.game.Init()
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m rdor) View() string {
	return "\n" + m.list.View()
}
