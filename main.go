package main

import (
	"fmt"
	"io"
	"os"

	"github.com/zrcoder/tgame/internal/hanoi"
	"github.com/zrcoder/tgame/internal/sokoban"
	"github.com/zrcoder/tgame/pkg/style"
	"github.com/zrcoder/tgame/pkg/style/color"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	items := []list.Item{
		item{name: "hanoi", instance: hanoi.New()},
		item{name: "sokoban", instance: sokoban.New()},
	}
	const listHeight = 14
	const defaultWidth = 20
	m := tgame{list: list.New(items, itemDelegate{}, defaultWidth, listHeight)}
	m.list.Title = "Welcome to `tgame`"
	m.list.Styles.Title = style.Title
	m.list.SetShowStatusBar(false)
	m.list.SetFilteringEnabled(false)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

type item struct {
	name     string
	instance tea.Model
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
		fn = func(s string) string {
			return itemSelectedRender("> " + s)
		}
	}
	fmt.Fprint(w, fn(fmt.Sprintf("%d. %s", index+1, i.name)))
}

type tgame struct {
	list list.Model
}

func (m tgame) Init() tea.Cmd { return nil }

func (m tgame) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch val := msg.String(); val {
		case "enter":
			i := m.list.SelectedItem().(item)
			return i.instance, i.instance.Init()
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m tgame) View() string {
	return "\n" + m.list.View()
}
