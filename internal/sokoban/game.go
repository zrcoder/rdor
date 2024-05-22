package sokoban

import (
	"embed"
	"strconv"
	"strings"

	"github.com/zrcoder/rdor/pkg/game"
	"github.com/zrcoder/rdor/pkg/grid"
	"github.com/zrcoder/rdor/pkg/keys"
	"github.com/zrcoder/rdor/pkg/style/color"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	name     = "Sokoban"
	maxLevel = 51

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

func New() game.Game {
	return &sokoban{Base: game.New(name)}
}

type sokoban struct {
	blocks map[rune]string
	buf    *strings.Builder
	*game.Base
	grid     *grid.Grid[rune]
	helpGrid *grid.Grid[rune]
	upKey    *key.Binding
	rightKey *key.Binding
	downKey  *key.Binding
	leftKey  *key.Binding
	myPos    grid.Position
}

func (s *sokoban) Init() tea.Cmd {
	s.RegisterView(s.view)
	s.RegisterHelp(s.helpInfo)
	s.RegisterLevels(maxLevel, s.loadLever)
	s.blocks = map[rune]string{
		wall:      lipgloss.NewStyle().Background(color.Orange).Render(" = "),
		me:        " ⦿ ", // ♾ ⚉ ⚗︎ ⚘ ☻
		blank:     "   ",
		slot:      lipgloss.NewStyle().Background(color.Violet).Render("   "),
		box:       lipgloss.NewStyle().Background(color.Red).Render(" x "),
		boxInSlot: lipgloss.NewStyle().Background(color.Green).Render("   "),
		meInSlot:  lipgloss.NewStyle().Background(color.Violet).Render(" ⦿ "),
	}
	s.upKey = &keys.Up
	s.leftKey = &keys.Left
	s.downKey = &keys.Down
	s.rightKey = &keys.Right
	s.ClearGroups()
	s.AddKeyGroup(game.KeyGroup{s.upKey, s.leftKey, s.downKey, s.rightKey})
	s.buf = &strings.Builder{}
	return s.Base.Init()
}

func (s *sokoban) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	b, bcmd := s.Base.Update(msg)
	if b != s.Base {
		return b, bcmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, *s.upKey):
			s.move(grid.Up)
		case key.Matches(msg, *s.leftKey):
			s.move(grid.Left)
		case key.Matches(msg, *s.downKey):
			s.move(grid.Down)
		case key.Matches(msg, *s.rightKey):
			s.move(grid.Right)
		}
	}
	return s, bcmd
}

func (s *sokoban) view() string {
	s.buf.Reset()
	s.grid.Range(func(_ grid.Position, char rune, isLineEnd bool) (end bool) {
		s.buf.WriteString(s.blocks[char])
		if isLineEnd {
			s.buf.WriteByte('\n')
		}
		return
	})
	return s.buf.String()
}

func (s *sokoban) helpInfo() string {
	return "Our goal is to push all the boxes into the slots without been stuck somewhere."
}

func (s *sokoban) loadLever(i int) {
	data, err := levelsFS.ReadFile("levels/" + strconv.Itoa(i+1) + ".txt")
	if err != nil {
		panic(err)
	}
	s.grid = grid.NewWithString(string(data))
	s.helpGrid = s.grid.Copied()
	s.grid.Range(func(pos grid.Position, char rune, _ bool) (end bool) {
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
	if s.success() {
		s.SetSuccess("")
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
	switch char {
	case blank:
		s.grid.Set(dest, box)
	case slot:
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
	s.grid.Range(func(_ grid.Position, char rune, _ bool) (end bool) {
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
	s.grid.Range(func(pos grid.Position, char rune, _ bool) (end bool) {
		if char == me || char == meInSlot {
			s.myPos = pos
			return true
		}
		return
	})
}
