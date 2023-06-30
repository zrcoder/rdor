package npuzzle

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zrcoder/rdor/pkg/dialog"
	"github.com/zrcoder/rdor/pkg/grid"
	"github.com/zrcoder/rdor/pkg/model"
	"github.com/zrcoder/rdor/pkg/style"
)

type nPuzzle struct {
	parent      tea.Model
	n           int
	title       string
	help        string
	grid        *grid.Grid
	target      *grid.Grid
	keys        *keyMap
	keysHelp    help.Model
	buf         *strings.Builder
	rows        []string
	cols        []string
	blank       grid.Position
	directions  []grid.Direction
	rd          *rand.Rand
	showSuccess bool
}

func New() model.Game { return &nPuzzle{} }

func (p *nPuzzle) SetParent(parent tea.Model) { p.parent = parent }

const (
	Name     = "N-Puzzle"
	defaultN = 3
)

func (p *nPuzzle) Init() tea.Cmd {
	p.n = defaultN
	p.title = style.Title.Render("N-Puzzle")
	p.buf = &strings.Builder{}
	p.rows = []string{"A", "B", "C", "D", "E"}
	p.cols = []string{"1", "2", "3", "4", "5"}
	p.rd = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	p.directions = []grid.Direction{grid.Up, grid.Left, grid.Down, grid.Right}
	p.keys = getKeys()
	p.keysHelp = help.New()
	p.keysHelp.ShowAll = true
	p.set()
	return nil
}
func (p *nPuzzle) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		p.showSuccess = false
		switch {
		case key.Matches(msg, p.keys.Home):
			return p.parent, nil
		case key.Matches(msg, p.keys.Shuffle):
			p.shuffle()
		case key.Matches(msg, p.keys.Next):
			p.n = 3 + (p.n+1)%3 // 3, 4, 5 only
			p.set()
		case key.Matches(msg, p.keys.Up):
			p.move(grid.Down)
		case key.Matches(msg, p.keys.Down):
			p.move(grid.Up)
		case key.Matches(msg, p.keys.Left):
			p.move(grid.Right)
		case key.Matches(msg, p.keys.Right):
			p.move(grid.Left)
		case msg.String() == "ctrl+c":
			return p, tea.Quit
		}
	}
	return p, nil
}
func (p *nPuzzle) View() string {
	if p.showSuccess {
		return dialog.Success("").WhiteSpaceChars(Name).String()
	}

	curBoard := p.drawBoard()
	p.buf.Reset()
	p.buf.WriteString("\n" + p.title + "\n\n")
	p.buf.WriteString(curBoard)
	p.buf.WriteString("\n\n")
	p.buf.WriteString(p.help)
	p.buf.WriteString("\n\n")
	p.buf.WriteString(p.keysHelp.View(p.keys))
	p.buf.WriteString("\n\n")
	return p.buf.String()
}

func (p *nPuzzle) set() {
	p.help = style.Help.Render(fmt.Sprintf("%d✗%d", p.n, p.n))
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

func (p *nPuzzle) drawBoard() string {
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
	p.showSuccess = p.success()
}

func (p *nPuzzle) success() bool {
	return p.grid.Equal(p.target)
}
