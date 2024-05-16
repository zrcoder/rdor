package crossword

import (
	"embed"
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/rdor/pkg/game"
	"github.com/zrcoder/rdor/pkg/grid"
	"github.com/zrcoder/rdor/pkg/keys"
	"github.com/zrcoder/rdor/pkg/style/color"
)

//go:embed levels/*.toml
var lvsFS embed.FS

var (
	currentBg = lipgloss.NewStyle().Background(color.Yellow)
	fixBg     = lipgloss.NewStyle().Background(color.Green)
	errGg     = lipgloss.NewStyle().Background(color.Red)
)

const (
	name                     = "成语填字"
	size                     = 9
	candidatesLimit          = 26
	emptyWord                = '　'
	blankWord                = '〇'
	candidatesCountInOneLine = 6
)

func New() game.Game {
	return &crossword{Base: game.New(name)}
}

type crossword struct {
	*game.Base

	grid         *Grid
	cadidates    []*Word
	boardBuf     *strings.Builder
	candidateBuf *strings.Builder
	state        string
	downKey      *key.Binding
	leftKey      *key.Binding
	upKey        *key.Binding
	rightKey     *key.Binding
	directions   []grid.Direction
	pos          grid.Position
}

func (c *crossword) Init() tea.Cmd {
	c.RegisterView(c.view)
	c.RegisterLevels(1, c.set)
	c.boardBuf = &strings.Builder{}
	c.candidateBuf = &strings.Builder{}
	c.directions = []grid.Direction{grid.Up, grid.Left, grid.Down, grid.Right}
	c.upKey = &keys.Up
	c.leftKey = &keys.Left
	c.downKey = &keys.Down
	c.rightKey = &keys.Right
	c.ClearGroups()
	c.AddKeyGroup(game.KeyGroup{c.upKey, c.leftKey, c.downKey, c.rightKey})
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
		case key.Matches(msg, *c.upKey):
			c.move(grid.Up)
		case key.Matches(msg, *c.downKey):
			c.move(grid.Down)
		case key.Matches(msg, *c.leftKey):
			c.move(grid.Left)
		case key.Matches(msg, *c.rightKey):
			c.move(grid.Right)
		default:
			if len(msg.Runes) > 0 && unicode.IsLetter(msg.Runes[0]) {
				i := int(unicode.ToLower(msg.Runes[0]) - 'a')
				c.pick(i)
			}
		}
	}
	return c, cmd
}

func (c *crossword) view() string {
	if c.Err != nil {
		return ""
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		c.boardView(),
		c.candidatesView(),
		c.state)
}

func (c *crossword) set(i int) {
	data, err := lvsFS.ReadFile(fmt.Sprintf("levels/%02d.toml", i))
	if err != nil {
		c.SetError(err)
		return
	}
	lvl := &Level{}
	err = toml.Unmarshal(data, lvl)
	if err != nil {
		c.SetError(err)
		return
	}
	if c.grid, err = c.newGrid(lvl.Grid); err != nil {
		c.SetError(err)
		return
	}
	if c.cadidates, err = c.newCandidates(lvl.Candidates); err != nil {
		c.SetError(err)
		return
	}
}

func (c *crossword) newGrid(cfg string) (*Grid, error) {
	if strings.Count(cfg, "\n") > size {
		return nil, errors.New("配置中有太多行")
	}
	g := &Grid{}
	c.pos.Row = -1
	for i, row := range strings.SplitN(cfg, "\n", size) {
		if utf8.RuneCountInString(row) > size {
			return nil, errors.New("一行中有太多字")
		}
		for j, v := range []rune(row) {
			if v == emptyWord {
				continue
			}
			g[i][j] = &Word{char: v, isFixed: v != blankWord, isBlank: v == blankWord}
			if c.pos.Row == -1 && g[i][j].isBlank {
				c.pos.Row = i
				c.pos.Col = j
			}
		}
	}
	if c.pos.Row == -1 {
		return nil, errors.New("没有空格要填")
	}
	return g, nil
}

func (c *crossword) newCandidates(cfg string) ([]*Word, error) {
	size := utf8.RuneCountInString(cfg)
	if size > candidatesLimit {
		return nil, errors.New("too many candidaters")
	}
	res := make([]*Word, size)
	for i, v := range []rune(cfg) {
		res[i] = &Word{char: v, candinatePos: i}
	}
	return res, nil
}

func (c *crossword) boardView() string {
	c.boardBuf.Reset()
	for i, row := range c.grid {
		for j, word := range row {
			if word == nil {
				c.boardBuf.WriteRune(emptyWord)
				continue
			}
			s := string(word.char)
			switch {
			case word.isBlank:
				s = string(blankWord)
			case word.isFixed:
				s = fixBg.Render(s)
			case word.isWrong:
				s = errGg.Render(s)
			}
			if i == c.pos.Row && j == c.pos.Col {
				s = currentBg.Render(s)
			}
			c.boardBuf.WriteString(s)
		}
		c.boardBuf.WriteString("\n")
	}
	return c.boardBuf.String()
}

func (c *crossword) candidatesView() string {
	c.boardBuf.Reset()
	for i, w := range c.cadidates {
		c.boardBuf.WriteRune(rune('A' + i))
		c.boardBuf.WriteRune(':')
		if w != nil {
			c.boardBuf.WriteRune(w.char)
		} else {
			c.boardBuf.WriteRune(emptyWord)
		}
		if (i+1)%candidatesCountInOneLine == 0 {
			c.boardBuf.WriteRune('\n')
		} else {
			c.boardBuf.WriteRune(emptyWord)
		}
	}
	return c.boardBuf.String()
}

func (c *crossword) move(d grid.Direction) {
	cnt := 0
	for i, j := (c.pos.Row+size+d.Dy)%size, (c.pos.Col+size+d.Dx)%size; cnt < size*size; cnt++ {
		if i == c.pos.Row && j == c.pos.Col {
			if d.Dy == 0 {
				i = (i + 1) % size
				j = 0
			} else {
				i = 0
				j = (j + 1) % size
			}
		} else {
			if c.grid[i][j] != nil && !c.grid[i][j].isFixed {
				c.pos.Row = i
				c.pos.Col = j
				break
			}
			i, j = (i+d.Dy+size)%size, (j+d.Dx+size)%size
		}
	}
	if cnt == size*size {
		c.SetError(errors.New("无法移动"))
		return
	}
}

func (c *crossword) pick(i int) {
	word := c.cadidates[i]
	if word == nil {
		return
	}
	tmp := c.grid[c.pos.Row][c.pos.Col]
	c.grid[c.pos.Row][c.pos.Col] = word
	if tmp != nil {
		c.cadidates[tmp.candinatePos] = tmp
	}
	if c.success() {
		c.SetSuccess("")
	}
}

func (c *crossword) success() bool {
	return false
}
