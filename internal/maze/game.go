package maze

import (
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zrcoder/rdor/internal/maze/levels"
	"github.com/zrcoder/rdor/pkg/style"
)

type maze struct {
	title     string
	helpInfo  string
	charMap   map[rune]rune
	keys      keyMap
	keysHelp  help.Model
	me        position
	goals     map[position]bool
	main      [][]rune
	origin    [][]rune
	rand      *rand.Rand
	levelName string
	buf       *strings.Builder
}

func New() tea.Model { return &maze{} }

type position struct {
	x, y int
}

type direction int

const (
	up direction = iota
	down
	right
	left
)

const (
	verticalWall   = '┃'
	horizontalWall = '━'
	corner         = '•'
	me             = '◉'
	goal           = '★'
	blank          = ' '
)

func (m *maze) Init() tea.Cmd {
	m.title = style.Title.Render("Maze")
	m.helpInfo = style.Help.Render("Our goal is to take all the stars in the maze.")
	m.charMap = map[rune]rune{
		'|': verticalWall,
		'-': horizontalWall,
		'o': corner,
		'S': me,
		'G': goal,
	}
	m.rand = rand.New(rand.NewSource(int64(time.Now().UnixNano())))
	err := m.pickOne()
	if err != nil {
		panic(err)
	}
	m.keys = keys
	m.keysHelp = help.New()
	m.keysHelp.ShowAll = true
	m.buf = &strings.Builder{}
	return nil
}
func (m *maze) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Up):
			m.move(up)
		case key.Matches(msg, m.keys.Left):
			m.move(left)
		case key.Matches(msg, m.keys.Down):
			m.move(down)
		case key.Matches(msg, m.keys.Right):
			m.move(right)
		case key.Matches(msg, m.keys.Pick):
			err := m.pickOne()
			if err != nil {
				panic(err)
			}
		case key.Matches(msg, m.keys.Reset):
			m.reset()
		}
	}
	return m, nil
}

func (m *maze) View() string {
	m.buf.Reset()
	m.buf.WriteString("\n" + m.title + "\n\n")
	for _, line := range m.main {
		m.buf.WriteString(string(line))
		m.buf.WriteByte('\n')
	}
	m.buf.WriteByte('\n')
	m.buf.WriteString(style.Help.Render("level: " + m.levelName))
	if m.success() {
		m.buf.WriteString(style.Success.Render("  Success!"))
	}
	m.buf.WriteString("\n\n")
	m.buf.WriteString(style.Help.Render(m.helpInfo))
	m.buf.WriteString("\n\n")
	m.buf.WriteString(m.keysHelp.View(m.keys))
	m.buf.WriteByte('\n')
	return m.buf.String()
}

func (m *maze) pickOne() error {
	m.levelName = levels.Names[rand.Intn(len(levels.Names))]
	return m.load()
}

func (m *maze) load() error {
	level, err := levels.ReadLevel(m.levelName)
	if err != nil {
		return err
	}
	m.goals = map[position]bool{}
	lines := strings.Split(level, "\n")
	m.main = make([][]rune, len(lines))
	m.origin = make([][]rune, len(lines))
	for i, line := range lines {
		m.main[i] = make([]rune, 0, len(line))
		m.origin[i] = make([]rune, 0, len(line))
		for j, v := range line {
			v, ok := m.charMap[v]
			if ok {
				if v == me {
					m.me = position{x: j, y: i}
				} else if v == goal {
					m.goals[position{x: j, y: i}] = true
				}
				m.main[i] = append(m.main[i], v)
				m.origin[i] = append(m.origin[i], v)
			} else {
				m.main[i] = append(m.main[i], blank)
				m.origin[i] = append(m.origin[i], blank)
			}
		}
	}
	return nil
}

func (m *maze) reset() {
	m.goals = map[position]bool{}
	for i, line := range m.origin {
		for j, v := range line {
			if v == me {
				m.me = position{x: j, y: i}
			} else if v == goal {
				m.goals[position{x: j, y: i}] = true
			}
			m.main[i][j] = v
		}
	}
}

func (m *maze) move(d direction) {
	y, x := m.me.y, m.me.x
	switch d {
	case up:
		y--
	case left:
		x -= 2
	case down:
		y++
	case right:
		x += 2
	}
	if m.outBound(y, x) || m.main[y][x] == horizontalWall || m.main[y][x] == verticalWall {
		return
	}
	switch d {
	case up:
		y--
	case left:
		x -= 2
	case down:
		y++
	case right:
		x += 2
	}
	m.moveMe(y, x)
}

func (m *maze) outBound(y, x int) bool {
	return y < 0 || y >= len(m.main) || x < 0 || x >= len(m.main[y])
}

func (m *maze) moveMe(dy, dx int) {
	if m.main[dy][dx] == goal {
		delete(m.goals, position{x: dx, y: dy})
	}
	m.main[dy][dx] = me
	m.main[m.me.y][m.me.x] = blank
	m.me = position{y: dy, x: dx}
}

func (m *maze) success() bool {
	return len(m.goals) == 0
}
