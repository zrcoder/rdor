package sokoban

import (
	"embed"
	"errors"
	"fmt"
	_ "fmt"
	"strconv"
	"strings"

	"github.com/zrcoder/tgame/pkg/style"
	"github.com/zrcoder/tgame/pkg/style/color"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Direction struct {
	x, y int
}

var (
	up    = Direction{0, -1}
	down  = Direction{0, 1}
	left  = Direction{-1, 0}
	right = Direction{1, 0}
)

const (
	maxLevel = 51

	wall      = '#'
	me        = '@'
	blank     = ' '
	slot      = 'X'
	box       = 'O'
	boxInSlot = '*'
	meInSlot  = '.'
)

var (
	blocks = map[rune]string{
		wall:      lipgloss.NewStyle().Background(color.Orange).Render(" = "),
		me:        " ◉ ",
		blank:     "   ",
		slot:      lipgloss.NewStyle().Background(color.Violet).Render("   "),
		box:       lipgloss.NewStyle().Background(color.Red).Render(" x "),
		boxInSlot: lipgloss.NewStyle().Background(color.Green).Render("   "),
		meInSlot:  lipgloss.NewStyle().Background(color.Violet).Render(" ◉ "),
	}
)

//go:embed levels
var levelsFS embed.FS

type sokoban struct {
	level    int
	keys     keyMap
	keysHelp help.Model
	input    textinput.Model

	showHelp bool
	origin   [][]rune
	main     [][]rune
	err      error
	y        int
	x        int
	debug    bool
}

func New() *sokoban {
	res := &sokoban{
		keys:     keys,
		keysHelp: help.New(),
		input:    textinput.New(),
	}
	res.keysHelp.ShowAll = true
	res.loadLever()
	res.input.Placeholder = "1-51"
	return res
}

func (s *sokoban) Init() tea.Cmd { return nil }

func (s *sokoban) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.Quit):
			return s, tea.Quit
		case key.Matches(msg, s.keys.Help):
			s.showHelp = !s.showHelp
			s.keysHelp.ShowAll = !s.showHelp
		case key.Matches(msg, s.keys.Debug):
			s.debug = !s.debug
		case key.Matches(msg, s.keys.Up):
			s.move(up)
		case key.Matches(msg, s.keys.Left):
			s.move(left)
		case key.Matches(msg, s.keys.Down):
			s.move(down)
		case key.Matches(msg, s.keys.Right):
			s.move(right)
		case key.Matches(msg, s.keys.Next):
			s.level = (s.level + 1) % maxLevel
			s.loadLever()
		case key.Matches(msg, s.keys.Previous):
			s.level = (s.level + maxLevel - 1) % maxLevel
			s.loadLever()
		case key.Matches(msg, s.keys.Set):
			return s, s.input.Focus()
		case key.Matches(msg, s.keys.Reset):
			s.reset()
		default:
			if msg.Type == tea.KeyEnter && s.input.Focused() {
				s.input.Blur()
				n, err := strconv.Atoi(s.input.Value())
				s.input.SetValue("")
				if err != nil {
					s.err = err
					return s, nil
				}
				if n < 1 || n > maxLevel+1 {
					s.err = errors.New("level out of range")
					return s, nil
				}
				s.level = n - 1
				s.loadLever()
			}
		}
	}
	return s, cmd
}

func (s *sokoban) View() string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("-- Sokoban Level %d of %d --\n", s.level+1, maxLevel))

	if s.debug {
		for _, line := range s.main {
			buf.WriteString(string(line))
		}
	} else if s.input.Focused() {
		buf.WriteString("which level do you like?\n")
		buf.WriteString(s.input.View())
		buf.WriteByte('\n')
	} else {
		for _, line := range s.main {
			for _, v := range line {
				buf.WriteString(blocks[v])
			}
			buf.WriteByte('\n')
		}
		if s.success() {
			buf.WriteString(style.Success.Render("Success!"))
		}
	}

	buf.WriteByte('\n')
	buf.WriteString(s.keysHelp.View(s.keys))
	buf.WriteByte('\n')

	return buf.String()
}

func (s *sokoban) loadLever() {
	data, err := levelsFS.ReadFile("levels/" + strconv.Itoa(s.level+1) + ".txt")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")
	maxWidth := 0
	for y, line := range lines {
		maxWidth = max(maxWidth, len(line))
		for x, v := range line {
			if v == me || v == meInSlot {
				s.y = y
				s.x = x
			}
		}
	}
	s.main = make([][]rune, len(lines))
	s.origin = make([][]rune, len(lines))
	for i, line := range lines {
		s.main[i] = make([]rune, maxWidth)
		s.origin[i] = make([]rune, maxWidth)
		copy(s.main[i], []rune(line))
		copy(s.origin[i], []rune(line))
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (s *sokoban) move(d Direction) {
	y, x := s.y+d.y, s.x+d.x
	if s.outRange(y, x) {
		return
	}
	switch s.main[y][x] {
	case blank, slot:
		s.moveMe(y, x)
	case box, boxInSlot:
		nx, ny := x+d.x, y+d.y
		if s.outRange(ny, nx) {
			return
		}
		if s.main[ny][nx] == blank || s.main[ny][nx] == slot {
			s.moveBox(y, x, ny, nx)
			s.moveMe(y, x)
		}
	}
}

func (s *sokoban) outRange(y, x int) bool {
	return y < 0 || y >= len(s.main) || x < 0 || x >= len(s.main[0])
}

func (s *sokoban) moveMe(y, x int) {
	if s.main[y][x] == blank {
		s.main[y][x] = me
	} else {
		s.main[y][x] = meInSlot
	}
	if s.main[s.y][s.x] == me {
		s.main[s.y][s.x] = blank
	} else {
		s.main[s.y][s.x] = slot
	}
	s.y, s.x = y, x
}

func (s *sokoban) moveBox(srcY, srcX, destY, destX int) {
	if s.main[destY][destX] == blank {
		s.main[destY][destX] = box
	} else if s.main[destY][destX] == slot {
		s.main[destY][destX] = boxInSlot
	}
	if s.main[srcY][srcX] == box {
		s.main[srcY][srcX] = blank
	} else {
		s.main[srcY][srcX] = slot
	}
}

func (s *sokoban) success() bool {
	for _, line := range s.main {
		for _, v := range line {
			if v == slot {
				return false
			}
		}
	}
	return true
}

func (s *sokoban) reset() {
	for i, row := range s.origin {
		for j, v := range row {
			if v == me || v == meInSlot {
				s.y = i
				s.x = j
			}
			s.main[i][j] = v
		}
	}
}
