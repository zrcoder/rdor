package last

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/zrcoder/rdor/pkg/grid"
	"github.com/zrcoder/rdor/pkg/model"
	"github.com/zrcoder/rdor/pkg/style"
	"github.com/zrcoder/rdor/pkg/style/color"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type last struct {
	parent      tea.Model
	title       string
	levels      []*level
	levelIndex  int
	commonCells int // the lefted commen cells' number, exclude the 2 players
	rd          *rand.Rand
	keys        *keyMap
	keysHepl    help.Model
	grid        *grid.Grid
	helpGrid    *grid.Grid
	charDic     map[rune]string
	buf         *strings.Builder
	players     [2]grid.Position // players[0]: me, players[1]: rival
	playerIndex int
	eatingLeft  int
	eating      bool
	eatingPath  *pathStack
	setting     bool
	err         error
	showHelp    bool
}

func New() model.Game                      { return &last{} }
func (l *last) SetParent(parent tea.Model) { l.parent = parent }

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

var playSyles = []lipgloss.Style{
	lipgloss.NewStyle().Foreground(color.Red),
	lipgloss.NewStyle().Foreground(color.Orange),
	lipgloss.NewStyle().Foreground(color.Yellow),
	lipgloss.NewStyle().Foreground(color.Green),
	lipgloss.NewStyle().Foreground(color.Blue),
	lipgloss.NewStyle().Foreground(color.Indigo),
	lipgloss.NewStyle().Foreground(color.Violet),
}

func (l *last) Init() tea.Cmd {
	l.title = style.Title.Render("Last")
	l.rd = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	l.keys = getKeys()
	l.keysHepl = help.New()
	l.levels = getDefaultLevers()
	l.buf = &strings.Builder{}
	l.setLevel()
	return l.lifeTransform()
}

func (l *last) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		return l, tea.Batch(l.lifeTransform(), l.eat())
	case tea.KeyMsg:
		l.err = nil
		switch {
		case key.Matches(msg, l.keys.Home):
			return l.parent, nil
		case key.Matches(msg, l.keys.Reset):
			l.setLevel()
		case key.Matches(msg, l.keys.Next):
			l.levelIndex = (l.levelIndex + 1) % len(l.levels)
			l.setLevel()
		case key.Matches(msg, l.keys.Previous):
			l.levelIndex = (l.levelIndex + len(l.levels) - 1) % len(l.levels)
			l.setLevel()
		case key.Matches(msg, l.keys.Numbers):
			n, _ := strconv.Atoi(msg.String())
			l.playerIndex = 0
			l.eatingLeft = n
			l.eating = true
			return l, l.eat()
		case key.Matches(msg, l.keys.Help):
			l.showHelp = !l.showHelp
		case msg.String() == "ctrl+c":
			return l, tea.Quit
		default:
			if !l.setting {
				return l, nil
			}
			switch msg.String() {
			case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
				if l.setting || l.eating {
					l.err = errors.New("wait, please")
				} else {
					l.err = fmt.Errorf("only 1-%d cells can be eaten each turn", l.currentLevel().eatingMax)
				}
			case "y", "Y", "enter":
				if l.setting {
					l.setted()
				}
			case "n", "N":
				if l.setting {
					l.setted()
					l.changeTurn()
				}
			}
		}
	}
	return l, nil
}

func (l *last) View() string {
	l.buf.Reset()
	l.buf.WriteString("\n" + l.title + "\n")
	if l.err != nil {
		l.buf.WriteString(style.Error.Render(l.err.Error()))
	}
	l.buf.WriteString("\n")

	l.grid.Range(func(pos grid.Position, char rune, isLineEnd bool) (end bool) {
		l.buf.WriteString(l.charDic[char])
		if isLineEnd {
			l.buf.WriteByte('\n')
		}
		return
	})
	l.buf.WriteString("\n")

	l.buf.WriteString(l.currentLevel().shortView() + "\n")
	if l.setting {
		l.buf.WriteString(style.Warn.Render("You go first? (y/n)"))
	} else {
		l.buf.WriteString(style.Help.Render("You:") + l.charDic[me] + style.Help.Render(" Rival:") + l.charDic[rival] + " ")
		l.buf.WriteString(style.Help.Render(fmt.Sprintf("Left: %2d  Turn:", l.commonCells+2)))
		if l.playerIndex == 0 {
			l.buf.WriteString(l.charDic[me])
		} else {
			l.buf.WriteString(l.charDic[rival])
		}
	}
	if l.success() {
		l.buf.WriteString(style.Success.Render("\nYou are the last :)"))
	} else if l.fail() {
		l.buf.WriteString(style.Error.Render("\nYour rival is the last :("))
	} else { // not end
		l.buf.WriteString("\n")
	}
	if l.showHelp {
		l.buf.WriteString("\n")
		l.buf.WriteString(l.currentLevel().View() + "\n")
	}
	l.buf.WriteString("\n")
	l.buf.WriteString(l.keysHepl.View(l.keys))
	return l.buf.String()
}

func (l *last) setLevel() {
	l.setting = true // wait for the user to decide whether to get started first
	l.keys.Numbers.SetEnabled(false)
	l.eatingPath = &pathStack{}
	l.commonCells = l.currentLevel().totalCells - 2 // minus the 2 plays
	l.rd.Shuffle(len(playSyles), func(i, j int) {
		playSyles[i], playSyles[j] = playSyles[j], playSyles[i]
	})
	l.charDic = map[rune]string{
		blank: "   ",
		cell:  " ◎ ",
		me:    playSyles[0].Render(" ◉ "),
		rival: playSyles[1].Render(" ◉ "),
	}
	l.genCells()
	l.playerIndex = 0
}

func (l *last) setted() {
	l.setting = false
	l.keys.Numbers.SetEnabled(true)
	ks := []string{"1", "2", "3", "4"}
	l.keys.Numbers.SetKeys(ks[:l.currentLevel().eatingMax]...)
	l.keys.Numbers.SetHelp(fmt.Sprintf("1-%d", l.currentLevel().eatingMax), "cells to eat")
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
	if l.eating || l.commonCells <= 0 {
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
	diff := l.commonCells - cells
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
		l.commonCells--
	}
	if l.eatingLeft == 0 {
		return l.changeTurn()
	}
	return l.doTick()
}

func (l *last) bfs() {
	queue := []grid.Position{l.currentPlayer()}
	vis := map[grid.Position]bool{l.currentPlayer(): true}
	directions := map[grid.Position]grid.Direction{}
	var cur grid.Position
	getPath := func() {
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
		if l.commonCells > 0 {
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
		return l.commonCells == 0 && isOpposite(char)
	}
	for len(queue) > 0 {
		cur = queue[0]
		queue = queue[1:]
		if foundTarget() {
			getPath()
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
}

func (l *last) currentLevel() *level {
	return l.levels[l.levelIndex]
}

func (l *last) changeTurn() tea.Cmd {
	if l.ended() {
		l.keys.Numbers.SetEnabled(false)
		return nil
	}
	l.playerIndex ^= 1      // 0->1 / 1->0
	if l.playerIndex == 1 { // the rival, auto eating
		l.eating = true
		if l.canEatRival() {
			l.eatingLeft = l.commonCells + 1
		} else if !l.currentLevel().hard {
			// just take random cells
			l.eatingLeft = 1 + l.rd.Intn(l.currentLevel().eatingMax)
		} else {
			// very clever rival
			total := l.commonCells + 1
			period := l.currentLevel().eatingMax + 1
			// the best strategy
			l.eatingLeft = total % period
			if l.eatingLeft == 0 {
				l.eatingLeft = 1 + l.rd.Intn(l.currentLevel().eatingMax)
			}
		}
	} else {
		l.eating = false
	}
	return l.doTick()
}

func (l *last) ended() bool {
	return l.commonCells == -1
}

func (l *last) success() bool {
	return l.ended() && l.playerIndex == 0
}

func (l *last) fail() bool {
	return l.ended() && l.playerIndex != 0
}

func (l *last) canEatRival() bool {
	return l.commonCells+1 <= l.currentLevel().eatingMax
}

func (l *last) currentPlayer() grid.Position {
	return l.players[l.playerIndex]
}
