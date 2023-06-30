package maze

import (
	"math/rand"
	"strings"
	"time"

	"github.com/zrcoder/rdor/internal/maze/levels"
	"github.com/zrcoder/rdor/pkg/dialog"
	"github.com/zrcoder/rdor/pkg/grid"
	"github.com/zrcoder/rdor/pkg/model"
	"github.com/zrcoder/rdor/pkg/style"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type maze struct {
	parent      tea.Model
	title       string
	helpInfo    string
	charMap     map[rune]rune
	keys        *keyMap
	keysHelp    help.Model
	myPos       grid.Position
	goals       map[grid.Position]bool
	grid        *grid.Grid
	helpGrid    *grid.Grid
	rand        *rand.Rand
	levelName   string
	buf         *strings.Builder
	showSuccess bool
}

func New() model.Game                      { return &maze{} }
func (m *maze) SetParent(parent tea.Model) { m.parent = parent }

var (
	up    = grid.Up
	down  = grid.Down
	right = grid.Right.Scale(2)
	left  = grid.Left.Scale(2)
)

const (
	Name           = "Maze"
	verticalWall   = '┃'
	horizontalWall = '━'
	corner         = '•'
	me             = '⦿'
	goal           = '❀' // ★
	blank          = ' '
)

func (m *maze) Init() tea.Cmd {
	m.title = style.Title.Render("Maze")
	m.helpInfo = style.Help.Render("Our goal is to take all the flowers in the maze.")
	m.charMap = map[rune]rune{
		'|':   verticalWall,
		'-':   horizontalWall,
		'o':   corner,
		'S':   me,
		'G':   goal,
		blank: blank,
	}
	m.rand = rand.New(rand.NewSource(int64(time.Now().UnixNano())))
	m.pickOne()
	m.keys = getKeys()
	m.keysHelp = help.New()
	m.keysHelp.ShowAll = true
	m.buf = &strings.Builder{}
	return nil
}

func (m *maze) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.showSuccess = false
		switch {
		case key.Matches(msg, m.keys.Home):
			return m.parent, nil
		case key.Matches(msg, m.keys.Up):
			m.move(up)
		case key.Matches(msg, m.keys.Left):
			m.move(left)
		case key.Matches(msg, m.keys.Down):
			m.move(down)
		case key.Matches(msg, m.keys.Right):
			m.move(right)
		case key.Matches(msg, m.keys.Pick):
			m.pickOne()
		case key.Matches(msg, m.keys.Reset):
			m.reset()
		case msg.String() == "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *maze) View() string {
	m.buf.Reset()
	m.buf.WriteString("\n" + m.title + "\n\n")

	if m.showSuccess {
		m.buf.WriteString(dialog.Success("").WhiteSpaceChars(Name).String())
		return m.buf.String()
	}

	m.grid.Range(func(pos grid.Position, char rune, isLineEnd bool) (end bool) {
		m.buf.WriteRune(char)
		if isLineEnd {
			m.buf.WriteRune('\n')
		}
		return
	})

	m.buf.WriteString(style.Help.Render("level: " + m.levelName))
	m.buf.WriteString("\n\n")
	m.buf.WriteString(style.Help.Render(m.helpInfo))
	m.buf.WriteString("\n\n")
	m.buf.WriteString(m.keysHelp.View(m.keys))
	m.buf.WriteByte('\n')
	return m.buf.String()
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
	m.showSuccess = m.success()
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
