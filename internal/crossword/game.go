package crossword

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/rdor/pkg/game"
	"github.com/zrcoder/rdor/pkg/grid"
	"github.com/zrcoder/rdor/pkg/keys"
	"github.com/zrcoder/rdor/pkg/style"
)

func New() game.Game {
	return &crossword{Base: game.New(name)}
}

type crossword struct {
	*game.Base
	*Level
	buf        *strings.Builder
	state      string
	directions []grid.Direction
	pos        grid.Position
	blankWord  *Word
	levels     int
}

func (c *crossword) Init() tea.Cmd {
	c.RegisterView(c.view)
	c.loadSummary()
	c.buf = &strings.Builder{}
	c.blankWord = &Word{state: WordStateBlank, char: emptyWord}
	c.directions = []grid.Direction{grid.Down, grid.Right, grid.Up, grid.Left}
	c.ClearGroups()
	c.AddKeyGroup(game.KeyGroup{&keys.Up, &keys.Left, &keys.Down, &keys.Right})
	c.set(0)
	return c.Base.Init()
}

func (c *crossword) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	b, cmd := c.Base.Update(msg)
	if b != c.Base {
		return b, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Up):
			c.move(grid.Up)
		case key.Matches(msg, keys.Down):
			c.move(grid.Down)
		case key.Matches(msg, keys.Left):
			c.move(grid.Left)
		case key.Matches(msg, keys.Right):
			c.move(grid.Right)
		default:
			if msg.Type == tea.KeyEnter {
				c.pick(-1)
				break
			}
			if len(msg.Runes) > 0 {
				letter := byte(unicode.ToUpper(msg.Runes[0]))
				if i, ok := c.Level.candidatesPos[letter]; ok {
					c.pick(i)
				}
			}
		}
	}
	return c, cmd
}

func (c *crossword) view() string {
	if c.Err != nil {
		return ""
	}
	last := c.state
	if c.success() {
		last = lg.JoinHorizontal(lg.Top, last, "  ", successBg.Render(" 成功 "))
	}
	return lg.JoinVertical(lg.Left,
		c.boardView(),
		c.candidatesView(),
		"",
		last)
}

func (c *crossword) loadSummary() {
	path := filepath.Join("levels", "index.toml")
	data, err := lvsFS.ReadFile(path)
	if err != nil {
		c.SetError(err)
		return
	}
	ls := &struct{ Levels int }{}
	err = toml.Unmarshal(data, ls)
	if err != nil {
		c.SetError(err)
		return
	}
	c.levels = ls.Levels
	c.RegisterLevels(ls.Levels, c.set)
}

func (c *crossword) set(i int) {
	path := filepath.Join("levels", fmt.Sprintf("%02d.toml", i))
	data, err := lvsFS.ReadFile(path)
	if err != nil {
		c.SetError(err)
		return
	}
	c.Level = &Level{crossword: c}
	err = toml.Unmarshal(data, c.Level)
	if err != nil {
		c.SetError(err)
		return
	}
	if c.adapt(); c.Err != nil {
		return
	}
	c.state = style.Help.Render(fmt.Sprintf("%d/%d", i+1, c.levels))
}
