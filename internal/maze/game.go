package maze

import (
	"math/rand"
	"strings"
	"time"

	"github.com/zrcoder/rdor/internal/maze/levels"
	"github.com/zrcoder/rdor/pkg/game"
	"github.com/zrcoder/rdor/pkg/grid"
	"github.com/zrcoder/rdor/pkg/keys"
	"github.com/zrcoder/rdor/pkg/style"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	Name           = "Maze"
	verticalWall   = '┃'
	horizontalWall = '━'
	corner         = '•'
	me             = '⦿'
	goal           = '❀'
	blank          = ' '
)

var (
	up    = grid.Up
	down  = grid.Down
	right = grid.Right.Scale(2)
	left  = grid.Left.Scale(2)
)

func New() game.Game {
	return &maze{Base: game.New(Name)}
}

type maze struct {
	*game.Base
	charMap   map[rune]rune
	upKey     key.Binding
	downKey   key.Binding
	leftKey   key.Binding
	rightKey  key.Binding
	pickKey   key.Binding
	myPos     grid.Position
	goals     map[grid.Position]bool
	grid      *grid.Grid
	helpGrid  *grid.Grid
	rand      *rand.Rand
	levelName string
	buf       *strings.Builder
}

func (m *maze) Init() tea.Cmd {
	m.KeyActionReset = m.reset
	m.ViewFunc = m.view
	m.HelpFunc = m.helpInfo
	m.charMap = map[rune]rune{
		'|':   verticalWall,
		'-':   horizontalWall,
		'o':   corner,
		'S':   me,
		'G':   goal,
		blank: blank,
	}
	m.upKey = keys.Up
	m.leftKey = keys.Left
	m.downKey = keys.Down
	m.rightKey = keys.Right
	m.pickKey = key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "pick one random level"),
	)
	m.Keys = []key.Binding{m.upKey, m.leftKey, m.downKey, m.rightKey, m.pickKey}
	m.rand = rand.New(rand.NewSource(int64(time.Now().UnixNano())))
	m.pickOne()
	m.buf = &strings.Builder{}
	return m.Base.Init()
}

func (m *maze) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	b, cmd := m.Base.Update(msg)
	if b != m.Base {
		return b, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.upKey):
			m.move(up)
		case key.Matches(msg, m.leftKey):
			m.move(left)
		case key.Matches(msg, m.downKey):
			m.move(down)
		case key.Matches(msg, m.rightKey):
			m.move(right)
		case key.Matches(msg, m.pickKey):
			m.pickOne()
		}
	}
	return m, cmd
}

func (m *maze) view() string {
	m.buf.Reset()

	m.grid.Range(func(pos grid.Position, char rune, isLineEnd bool) (end bool) {
		m.buf.WriteRune(char)
		if isLineEnd {
			m.buf.WriteRune('\n')
		}
		return
	})

	return lipgloss.JoinVertical(lipgloss.Center,
		strings.TrimRight(m.buf.String(), "\n"),
		style.Help.Render("level: "+m.levelName),
	)

}

func (m *maze) helpInfo() string {
	return "Our goal is to take all the flowers in the maze."
}

func (m *maze) pickOne() {
	m.levelName = levels.Names[rand.Intn(len(levels.Names))]
	m.load()
}

func (m *maze) load() {
	level, err := levels.ReadLevel(m.levelName)
	if err != nil {
		panic(err)
	}
	m.goals = map[grid.Position]bool{}
	m.grid = grid.New(level)
	m.helpGrid = grid.Copy(m.grid)
	m.reMap()
}

func (m *maze) reset() {
	m.goals = map[grid.Position]bool{}
	m.grid.Copy(m.helpGrid)
	m.reMap()
}

func (m *maze) reMap() {
	m.grid.Range(func(pos grid.Position, char rune, isLineEnd bool) (end bool) {
		v, ok := m.charMap[char]
		if !ok {
			panic("invalid level config")
		}
		if v == me {
			m.myPos = pos
		} else if v == goal {
			m.goals[pos] = true
		}
		m.grid.Set(pos, v)
		return
	})
}

func (m *maze) move(d grid.Direction) {
	pos := grid.TransForm(m.myPos, d)
	obj := m.grid.Get(pos)
	if m.grid.OutBound(pos) || obj == horizontalWall || obj == verticalWall {
		return
	}
	pos = grid.TransForm(pos, d)
	m.moveMe(pos)
	if m.success() {
		m.SetSuccess("")
	}
}

func (m *maze) moveMe(pos grid.Position) {
	if m.grid.Get(pos) == goal {
		delete(m.goals, pos)
	}
	m.grid.Set(pos, me)
	m.grid.Set(m.myPos, blank)
	m.myPos = pos
}

func (m *maze) success() bool {
	return len(m.goals) == 0
}
