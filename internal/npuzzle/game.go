package npuzzle

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/rdor/pkg/game"
	"github.com/zrcoder/rdor/pkg/grid"
	"github.com/zrcoder/rdor/pkg/keys"
	"github.com/zrcoder/rdor/pkg/model"
	"github.com/zrcoder/rdor/pkg/style"
)

type nPuzzle struct {
	*game.Game
	n          int
	state      string
	grid       *grid.Grid
	target     *grid.Grid
	shuffleKey key.Binding
	upKey      key.Binding
	leftKey    key.Binding
	downKey    key.Binding
	rightKey   key.Binding
	buf        *strings.Builder
	rows       []string
	cols       []string
	blank      grid.Position
	directions []grid.Direction
	rd         *rand.Rand
}

func New() model.Game {
	base := game.New(Name)
	res := &nPuzzle{Game: base}
	base.InitFunc = res.initialize
	base.UpdateFunc = res.update
	base.ViewFunc = res.view
	base.KeyActionNext = res.nextLevel
	return res
}

func (p *nPuzzle) SetParent(parent tea.Model) { p.Parent = parent }

const (
	Name     = "N-Puzzle"
	defaultN = 3
)

func (p *nPuzzle) initialize() tea.Cmd {
	p.n = defaultN
	p.buf = &strings.Builder{}
	p.rows = []string{"A", "B", "C", "D", "E"}
	p.cols = []string{"1", "2", "3", "4", "5"}
	p.rd = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	p.directions = []grid.Direction{grid.Up, grid.Left, grid.Down, grid.Right}
	p.shuffleKey = key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "shuffle"),
	)
	p.upKey = keys.Up
	p.leftKey = keys.Left
	p.downKey = keys.Down
	p.rightKey = keys.Right
	p.Keys = []key.Binding{p.upKey, p.leftKey, p.downKey, p.rightKey, p.shuffleKey}
	p.set()
	return nil
}
func (p *nPuzzle) update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, p.shuffleKey):
			p.shuffle()
		case key.Matches(msg, p.upKey):
			p.move(grid.Down)
		case key.Matches(msg, p.downKey):
			p.move(grid.Up)
		case key.Matches(msg, p.leftKey):
			p.move(grid.Right)
		case key.Matches(msg, p.rightKey):
			p.move(grid.Left)
		}
	}
	return nil
}

func (p *nPuzzle) view() string {
	return lipgloss.JoinVertical(lipgloss.Center,
		p.boardView(),
		p.state,
	)
}

func (p *nPuzzle) nextLevel() {
	p.n = 3 + (p.n+1)%3 // 3, 4, 5 only
	p.set()
}

func (p *nPuzzle) set() {
	p.state = style.Help.Render(fmt.Sprintf("%d✗%d", p.n, p.n))
	g := make([][]rune, p.n)
	for r := range g {
		g[r] = make([]rune, p.n)
		for c := range g[r] {
			g[r][c] = rune(r*p.n + c + 1)
		}
	}
	g[p.n-1][p.n-1] = 0
	p.grid = grid.New("")
	p.grid.SetData(g)
	p.blank = grid.Position{Row: p.n - 1, Col: p.n - 1}
	p.target = grid.Copy(p.grid)
	p.shuffle()
}

func (p *nPuzzle) boardView() string {
	p.buf.Reset()
	p.buf.WriteString("   ")
	for _, v := range p.cols[:p.n] {
		p.buf.WriteString(" " + v + "   ")
	}
	p.buf.WriteString("\n")
	p.buf.WriteString("  " + strings.Repeat("•━━━━", p.n))
	p.buf.WriteString("•\n")
	p.grid.Range(func(pos grid.Position, char rune, isLineEnd bool) (end bool) {
		s := "  "
		if char != 0 {
			r, c := int(char-1)/p.n, int(char-1)%p.n
			s = p.rows[r] + p.cols[c]
		}
		if pos.Col == 0 {
			p.buf.WriteString(p.rows[pos.Row] + " ")
		}
		p.buf.WriteString("┃ " + s + " ")
		if isLineEnd {
			p.buf.WriteString("┃\n")
			p.buf.WriteString("  " + strings.Repeat("•━━━━", p.n))
			p.buf.WriteString("•")
			if pos.Row != p.n-1 {
				p.buf.WriteString("\n")
			}
		}
		return false
	})
	return p.buf.String()
}

func (p *nPuzzle) shuffle() {
	for i := p.n * p.n * 11; i > 0; i-- {
		p.move(p.directions[p.rd.Intn(len(p.directions))])
	}
}

func (p *nPuzzle) move(d grid.Direction) {
	pos := grid.TransForm(p.blank, d)
	if p.grid.OutBound(pos) {
		return
	}
	char := p.grid.Get(pos)
	p.grid.Set(p.blank, char)
	p.grid.Set(pos, 0)
	p.blank = pos
	if p.success() {
		p.SetSuccess("")
	}
}

func (p *nPuzzle) success() bool {
	return p.grid.Equal(p.target)
}
