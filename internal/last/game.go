package last

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/zrcoder/rdor/pkg/grid"
	"github.com/zrcoder/rdor/pkg/style"
	"github.com/zrcoder/rdor/pkg/style/color"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type last struct {
	title       string
	levels      []*level
	levelIndex  int
	left        int
	rd          *rand.Rand
	grid        *grid.Grid
	helpGrid    *grid.Grid
	charDic     map[rune]string
	buf         *strings.Builder
	players     [2]grid.Position // players[0]: me, players[1]: rival
	playerIndex int
	eatingLeft  int
	eating      bool
	eatingPath  *pathStack
	err         error
	logFile     *os.File
}

func New() tea.Model { return &last{} }

type tickMsg time.Time

const (
	width        = 10
	height       = 10
	defaultTotal = 30
	defaultLimit = 2
)

const (
	blank = '.'
	cell  = 'c'
	me    = 'm'
	rival = 'r'
)

var (
	playSyles = []lipgloss.Style{
		lipgloss.NewStyle().Foreground(color.Orange),
		lipgloss.NewStyle().Foreground(color.Green),
	}
)

func (l *last) Init() tea.Cmd {
	// below 4 lines are just for debug, should delete
	// var err error
	// l.logFile, err = tea.LogToFile("last.log", "")
	// if err != nil {
	// 	panic(err)
	// }

	l.title = style.Title.Render("Last")
	l.rd = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	l.levels = getDefaultLevers()
	l.buf = &strings.Builder{}
	l.set()
	return l.lifeTransform()
}

func (l *last) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		return l, tea.Batch(l.lifeTransform(), l.eat())
	case tea.KeyMsg:
		val := msg.String()
		switch val {
		case "q":
			return l, tea.Quit
		case "1", "2":
			if l.eating {
				return l, nil
			}
			n, _ := strconv.Atoi(val)
			l.playerIndex = 0
			l.eatingLeft = n
			l.eating = true
			return l, l.eat()
		}
	}
	return l, nil
}

func (l *last) View() string {
	l.buf.Reset()
	l.buf.WriteString("\n" + l.title + "\n")
	l.grid.Range(func(pos grid.Position, char rune, isLineEnd bool) (end bool) {
		l.buf.WriteString(l.charDic[char])
		if isLineEnd {
			l.buf.WriteByte('\n')
		}
		return
	})
	l.buf.WriteString("\n\n")
	l.buf.WriteString(fmt.Sprintf("Cells: %d\n", l.left))
	l.buf.WriteString("You: " + l.charDic[me] + ", your rival: " + l.charDic[rival] + ".\n")
	if l.success() {
		l.buf.WriteString(style.Success.Render("success!\n"))
	} else if l.fail() {
		l.buf.WriteString(style.Error.Render("failed~\n"))
	} else { // not end
		if l.playerIndex == 0 {
			l.buf.WriteString("Your turn now\n")
		} else {
			l.buf.WriteString("The rival's turn\n")
		}
	}
	l.buf.WriteString("\n")

	return l.buf.String() // + "\n\nDebug:\n" + l.grid.String()
}

func (l *last) set() {
	l.eatingPath = &pathStack{}
	l.left = l.currentLevel().totalCells - 2 // minus the 2 plays
	l.rd.Shuffle(len(playSyles), func(i, j int) {
		playSyles[i], playSyles[j] = playSyles[j], playSyles[i]
	})
	l.charDic = map[rune]string{
		blank: "   ",
		cell:  " ◎ ",
		me:    playSyles[0].Render(" ● "),
		rival: playSyles[1].Render(" ● "),
	}
	l.genCells()
}

func (l *last) genCells() {
	l.grid = grid.New("")
	g := make([][]rune, height)
	for i := range g {
		g[i] = make([]rune, width)
	}
	g[0][0] = me
	g[0][1] = rival
	for i := 2; i < width*height; i++ {
		r, c := i/width, i%width
		if i < l.currentLevel().totalCells {
			g[r][c] = cell
		} else {
			g[r][c] = blank
		}
	}
	l.rd.Shuffle(width*height, func(i, j int) {
		ir, ic := i/width, i%width
		jr, jc := j/width, j%width
		g[ir][ic], g[jr][jc] = g[jr][jc], g[ir][ic]
	})
	l.grid.SetData(g)
	l.helpGrid = grid.Copy(l.grid)
	l.grid.Range(func(pos grid.Position, char rune, isLineEnd bool) (end bool) {
		if char == me {
			l.players[0] = pos
		} else if char == rival {
			l.players[1] = pos
		}
		return
	})
}

func (l *last) doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (l *last) lifeTransform() tea.Cmd {
	if l.eating || l.left <= 0 {
		return nil
	}
	cells := 0
	l.grid.Range(func(pos grid.Position, char rune, isLineEnd bool) (end bool) {
		l.helpGrid.Set(pos, l.grid.Get(pos))
		if char == me || char == rival {
			return
		}
		// Conway's Game of Life
		// See https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life
		switch l.countAliveNeighbours(pos) {
		case 2: // do nothing
		case 3:
			l.helpGrid.Set(pos, cell)
		default:
			l.helpGrid.Set(pos, blank)
		}
		if l.helpGrid.Get(pos) == cell {
			cells++
		}
		return
	})
	diff := l.left - cells
	if diff < 0 {
		l.removeCells(l.helpGrid, -diff)
	} else if diff > 0 {
		l.addCells(l.helpGrid, diff)
	}
	l.grid.Copy(l.helpGrid)
	return l.doTick()
}

func (l *last) countAliveNeighbours(pos grid.Position) int {
	cnt := 0
	for _, d := range grid.AllDirections {
		p := grid.TransForm(pos, d)
		if !l.grid.OutBound(p) && l.grid.Get(p) == cell {
			cnt++
		}
	}
	return cnt
}

func (l *last) removeCells(g *grid.Grid, n int) {
	for ; n > 0; n-- {
		l.changeCell(g, cell, blank)
	}
}
func (l *last) addCells(g *grid.Grid, n int) {
	for ; n > 0; n-- {
		l.changeCell(g, blank, cell)
	}
}

func (l *last) changeCell(g *grid.Grid, from, to rune) {
	for {
		i := l.rd.Intn(width * height)
		pos := grid.Position{Row: i / width, Col: i % width}
		if g.Get(pos) == from {
			g.Set(pos, to)
			return
		}
	}
}

func (l *last) eat() tea.Cmd {
	if !l.eating || l.ended() {
		return nil
	}
	if l.eatingPath.empty() {
		l.bfs()
	}
	// move one step
	pos := l.currentPlayer()
	char := l.grid.Get(pos)
	l.grid.Set(pos, blank)
	d := l.eatingPath.pop()
	pos = grid.TransForm(pos, d)
	l.grid.Set(pos, char)
	l.players[l.playerIndex] = pos
	// end moving
	if l.eatingPath.empty() {
		l.eatingLeft--
		l.left--
	}
	if l.eatingLeft == 0 {
		l.eating = false
		return l.changeTurn()
	}
	return l.doTick()
}

func (l *last) bfs() {
	l.logf("begin bfs, current player: %#v, cells left: %d,  current board:\n%s",
		l.currentPlayer(), l.left, l.grid.String())
	queue := []grid.Position{l.currentPlayer()}
	vis := map[grid.Position]bool{l.currentPlayer(): true}
	directions := map[grid.Position]grid.Direction{}
	var cur grid.Position
	getPath := func() {
		l.logf("begin to get path, cur: %#v, current player: %#v\ndirections: %d, vis: %d", cur, l.currentPlayer(), len(directions), len(vis))
		for cur != l.currentPlayer() {
			d := directions[cur]
			l.eatingPath.push(d)
			cur = grid.TransForm(cur, d.Opposite())
		}
	}
	isOpposite := func(char rune) bool {
		if l.playerIndex == 0 {
			return char == rival
		}
		return char == me
	}
	foundTarget := func() bool {
		char := l.grid.Get(cur)
		if char == cell {
			return true
		}
		if l.left > 0 {
			return false
		}
		return isOpposite(char)
	}
	canVisit := func(pos grid.Position) bool {
		if l.grid.OutBound(pos) || vis[pos] {
			return false
		}
		char := l.grid.Get(pos)
		if char == blank || char == cell {
			return true
		}
		return isOpposite(char)
	}
	for len(queue) > 0 {
		cur = queue[0]
		queue = queue[1:]
		if foundTarget() {
			getPath()
			l.logf("bfs found target, path is:\n%#v", l.eatingPath)
			return
		}
		for _, d := range grid.NormalDirections {
			pos := grid.TransForm(cur, d)
			if !canVisit(pos) {
				continue
			}
			vis[pos] = true
			directions[pos] = d
			queue = append(queue, pos)
		}
	}
	l.logf("bfs not found target!!!")
}

func (l *last) currentLevel() *level {
	return l.levels[l.levelIndex]
}

func (l *last) changeTurn() tea.Cmd {
	if l.ended() {
		return nil
	}
	l.playerIndex ^= 1      // 0->1 / 1->0
	if l.playerIndex == 1 { // the rival, auto eating
		l.eating = true
		if l.canEatRival() {
			l.eatingLeft = l.left + 1
		} else if !l.currentLevel().hard {
			// just take random cells
			l.eatingLeft = 1 + l.rd.Intn(l.currentLevel().eatingMax)
		} else {
			// very clever rival
			total := l.left + 1
			period := l.currentLevel().eatingMax + 1
			// the best strategy
			l.eatingLeft = total % period
			if l.eatingLeft == 0 {
				l.eatingLeft = 1 + l.rd.Intn(l.currentLevel().eatingMax)
			}
		}
	}
	return l.doTick()
}

func (l *last) ended() bool {
	return l.left == -1
}

func (l *last) success() bool {
	return l.ended() && l.playerIndex == 0
}

func (l *last) fail() bool {
	return l.ended() && l.playerIndex != 0
}

func (l *last) canEatRival() bool {
	return l.left+1 <= l.currentLevel().eatingMax
}

func (l *last) currentPlayer() grid.Position {
	return l.players[l.playerIndex]
}

func (l *last) logf(s string, v ...any) {
	if l.logFile != nil {
		l.logFile.WriteString(time.Now().String())
		l.logFile.WriteString(fmt.Sprintf(s, v...) + "\n")
	}
}
