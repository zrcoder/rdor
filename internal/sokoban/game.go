package sokoban

import (
	"embed"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/zrcoder/rdor/pkg/grid"
	"github.com/zrcoder/rdor/pkg/style"
	"github.com/zrcoder/rdor/pkg/style/color"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sokoban struct {
	title    string
	helpInfo string
	blocks   map[rune]string

	level    int
	keys     *keyMap
	keysHelp help.Model
	input    textinput.Model

	helpGrid *grid.Grid
	grid     *grid.Grid
	err      error
	myPos    grid.Position
	buf      *strings.Builder
}

func New() tea.Model { return &sokoban{} }

const (
	maxLevel         = 51
	inputPlaceholder = "1-51"

	wall      = '#'
	me        = '@'
	blank     = ' '
	slot      = 'X'
	box       = 'O'
	boxInSlot = '*'
	meInSlot  = '.'
)

//go:embed levels
var levelsFS embed.FS

func (s *sokoban) Init() tea.Cmd {
	s.title = style.Title.Render("Sokoban")
	s.helpInfo = style.Help.Render("Our goal is to push all the boxes into the slots without been stuck somewhere.")
	s.blocks = map[rune]string{
		wall:      lipgloss.NewStyle().Background(color.Orange).Render(" = "),
		me:        " ⦿ ", // ♾ ⚉ ⚗︎ ⚘ ☻
		blank:     "   ",
		slot:      lipgloss.NewStyle().Background(color.Violet).Render("   "),
		box:       lipgloss.NewStyle().Background(color.Red).Render(" x "),
		boxInSlot: lipgloss.NewStyle().Background(color.Green).Render("   "),
		meInSlot:  lipgloss.NewStyle().Background(color.Violet).Render(" ⦿ "),
	}
	s.keys = getKeys()
	s.keysHelp = help.New()
	s.input = textinput.New()
	s.buf = &strings.Builder{}
	s.keysHelp.ShowAll = true
	s.loadLever()
	s.input.Placeholder = inputPlaceholder
	return nil
}

func (s *sokoban) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		s.err = nil
		switch {
		case key.Matches(msg, s.keys.Quit):
			return s, tea.Quit
		case key.Matches(msg, s.keys.Up):
			s.move(grid.Up)
		case key.Matches(msg, s.keys.Left):
			s.move(grid.Left)
		case key.Matches(msg, s.keys.Down):
			s.move(grid.Down)
		case key.Matches(msg, s.keys.Right):
			s.move(grid.Right)
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
				s.setted(s.input.Value())
				s.input.SetValue("")
			}
		}
	}
	return s, cmd
}

func (s *sokoban) View() string {
	s.buf.Reset()
	s.buf.WriteString("\n" + s.title + "\n\n")
	s.grid.Range(func(pos grid.Position, char rune, isLineEnd bool) (end bool) {
		s.buf.WriteString(s.blocks[char])
		if isLineEnd {
			s.buf.WriteByte('\n')
		}
		return
	})
	s.buf.WriteString(style.Help.Render(fmt.Sprintf("- %d/%d - ", s.level+1, maxLevel)))
	if s.success() {
		s.buf.WriteString(style.Success.Render("Success!"))
	}
	s.buf.WriteByte('\n')

	if s.input.Focused() {
		s.buf.WriteString("\npick a level\n")
		s.buf.WriteString(s.input.View())
	} else {
		s.buf.WriteString("\n" + s.helpInfo + "\n")
		if s.err != nil {
			s.buf.WriteString("\n" + style.Error.Render(s.err.Error()) + "\n")
		}
		s.buf.WriteString("\n" + s.keysHelp.View(s.keys))
	}
	s.buf.WriteByte('\n')
	return s.buf.String()
}

func (s *sokoban) setted(level string) {
	n, err := strconv.Atoi(level)
	if err != nil {
		s.err = errors.New("invalid number")
		return
	}
	if n < 1 || n > maxLevel+1 {
		s.err = errors.New("level out of range")
		return
	}
	s.level = n - 1
	s.loadLever()
}

func (s *sokoban) loadLever() {
	data, err := levelsFS.ReadFile("levels/" + strconv.Itoa(s.level+1) + ".txt")
	if err != nil {
		panic(err)
	}
	s.grid = grid.New(string(data))
	s.helpGrid = grid.Copy(s.grid)
	s.grid.Range(func(pos grid.Position, char rune, isLineEnd bool) (end bool) {
		if char == me || char == meInSlot {
			s.myPos = pos
			return true
		}
		return
	})
}

func (s *sokoban) move(d grid.Direction) {
	pos := grid.TransForm(s.myPos, d)
	if s.grid.OutBound(pos) {
		return
	}
	switch s.grid.Get(pos) {
	case blank, slot:
		s.moveMe(pos)
	case box, boxInSlot:
		dest := grid.TransForm(pos, d)
		if s.grid.OutBound(dest) {
			return
		}
		char := s.grid.Get(dest)
		if char == blank || char == slot {
			s.moveBox(pos, dest)
			s.moveMe(pos)
		}
	}
}

func (s *sokoban) moveMe(p grid.Position) {
	if s.grid.Get(p) == blank {
		s.grid.Set(p, me)
	} else {
		s.grid.Set(p, meInSlot)
	}
	if s.grid.Get(s.myPos) == me {
		s.grid.Set(s.myPos, blank)
	} else {
		s.grid.Set(s.myPos, slot)
	}
	s.myPos = p
}

func (s *sokoban) moveBox(src, dest grid.Position) {
	char := s.grid.Get(dest)
	if char == blank {
		s.grid.Set(dest, box)
	} else if char == slot {
		s.grid.Set(dest, boxInSlot)
	}
	if s.grid.Get(src) == box {
		s.grid.Set(src, blank)
	} else {
		s.grid.Set(src, slot)
	}
}

func (s *sokoban) success() bool {
	res := true
	s.grid.Range(func(pos grid.Position, char rune, isLineEnd bool) (end bool) {
		if char == box {
			res = false
			return true
		}
		return
	})
	return res
}

func (s *sokoban) reset() {
	s.grid.Copy(s.helpGrid)
	s.grid.Range(func(pos grid.Position, char rune, isLineEnd bool) (end bool) {
		if char == me || char == meInSlot {
			s.myPos = pos
			return true
		}
		return
	})
}
