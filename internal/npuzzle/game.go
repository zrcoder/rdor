package npuzzle

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/zrcoder/rdor/pkg/game"
	"github.com/zrcoder/rdor/pkg/grid"
	"github.com/zrcoder/rdor/pkg/keys"
	"github.com/zrcoder/rdor/pkg/style"
)

const (
	name     = "N-Puzzle"
	defaultN = 3
)

var (
	rows = []string{"A", "B", "C", "D", "E"}
	cols = []string{"1", "2", "3", "4", "5"}
)

func New() game.Game {
	return &nPuzzle{Base: game.New(name)}
}

type nPuzzle struct {
	*game.Base

	rd         *rand.Rand
	grid       *grid.Grid[string]
	state      string
	downKey    *key.Binding
	leftKey    *key.Binding
	upKey      *key.Binding
	rightKey   *key.Binding
	rows       []string
	cols       []string
	directions []grid.Direction
	blank      grid.Position
	n          int
}

func (p *nPuzzle) Init() tea.Cmd {
	p.RegisterView(p.view)
	p.n = defaultN
	p.RegisterLevels(p.n, p.set)
	p.rd = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	p.directions = []grid.Direction{grid.Up, grid.Left, grid.Down, grid.Right}
	p.upKey = &keys.Up
	p.leftKey = &keys.Left
	p.downKey = &keys.Down
	p.rightKey = &keys.Right
	p.ClearGroups()
	p.AddKeyGroup(game.KeyGroup{p.upKey, p.leftKey, p.downKey, p.rightKey})
	p.set(0)
	return p.Base.Init()
}

func (p *nPuzzle) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	b, cmd := p.Base.Update(msg)
	if b != p.Base {
		return b, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, *p.upKey):
			p.move(grid.Down)
		case key.Matches(msg, *p.downKey):
			p.move(grid.Up)
		case key.Matches(msg, *p.leftKey):
			p.move(grid.Right)
		case key.Matches(msg, *p.rightKey):
			p.move(grid.Left)
		}
	}
	return p, cmd
}

func (p *nPuzzle) view() string {
	return lipgloss.JoinVertical(lipgloss.Center,
		p.boardView(),
		p.state,
	)
}

func (p *nPuzzle) set(i int) {
	p.n = i + 3
	p.rows = rows[:p.n]
	p.cols = cols[:p.n]
	p.state = style.Help.Render(fmt.Sprintf("%d✗%d", p.n, p.n))
	g := make([][]string, p.n)
	for r := range g {
		g[r] = make([]string, p.n)
		for c := range g[r] {
			g[r][c] = p.rows[r] + p.cols[c]
		}
	}
	g[p.n-1][p.n-1] = ""
	p.grid = grid.New[string](len(p.rows), len(p.cols))
	p.grid.SetData(g)
	p.blank = grid.Position{Row: p.n - 1, Col: p.n - 1}
	p.shuffle()
}

func (p *nPuzzle) boardView() string {
	t := table.New().Border(lg.NormalBorder()).BorderRow(true).StyleFunc(func(row, col int) lg.Style {
		return lg.NewStyle().Padding(0, 1)
	})
	p.grid.RangeRows(func(r int, row []string, isLast bool) (end bool) {
		t.Row(row...)
		return false
	})
	return lg.JoinVertical(lg.Center,
		strings.Join(p.cols, "    "),
		lg.JoinHorizontal(lg.Center, strings.Join(p.rows, "\n\n"), " ", t.String()))
}

func (p *nPuzzle) shuffle() {
	for i := p.n * p.n * 8; i > 0; i-- {
		p.move(p.directions[p.rd.Intn(len(p.directions))])
	}
}

func (p *nPuzzle) move(d grid.Direction) {
	pos := grid.TransForm(p.blank, d)
	if p.grid.OutBound(pos) {
		return
	}
	s := p.grid.Get(pos)
	p.grid.Set(p.blank, s)
	p.grid.Set(pos, "")
	p.blank = pos
	if p.success() {
		p.SetSuccess("")
	}
}

func (p *nPuzzle) success() bool {
	res := true
	p.grid.Range(func(pos grid.Position, s string, isLineEnd bool) (end bool) {
		if len(s) == 2 && (p.rows[pos.Row][0] != s[0] || p.cols[pos.Col][0] != s[1]) {
			res = false
			return true
		}
		return false
	})
	return res
}
