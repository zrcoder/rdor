package last

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/zrcoder/rdor/pkg/game"
	"github.com/zrcoder/rdor/pkg/grid"
	"github.com/zrcoder/rdor/pkg/style"
	"github.com/zrcoder/rdor/pkg/style/color"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	name         = "Last"
	width        = 10
	height       = 10
	defaultTotal = 30
	defaultLimit = 2

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

func New() game.Game {
	return &last{Base: game.New(name)}
}

type last struct {
	eatingPath *pathStack
	*game.Base
	rd          *rand.Rand
	grid        *grid.Grid[grid.Rune]
	helpGrid    *grid.Grid[grid.Rune]
	charDic     map[rune]string
	buf         *strings.Builder
	numbersKey  *key.Binding
	levels      []*level
	players     [2]grid.Position
	commonCells int
	playerIndex int
	eatingLeft  int
	levelIndex  int
	eating      bool
	setting     bool
}

type tickMsg time.Time

func (l *last) Init() tea.Cmd {
	l.levels = getDefaultLevers()
	l.RegisterLevels(len(l.levels), l.setLevel)
	l.RegisterView(l.view)
	l.rd = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	l.buf = &strings.Builder{}
	cmd := l.lifeTransform()
	return tea.Batch(l.Base.Init(), cmd)
}

func (l *last) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	b, cmd := l.Base.Update(msg)
	if b != l.Base {
		return b, cmd
	}

	switch msg := msg.(type) {
	case tickMsg:
		return l, tea.Batch(l.lifeTransform(), l.eat(), cmd)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, *l.numbersKey):
			n, _ := strconv.Atoi(msg.String())
			l.playerIndex = 0
			l.eatingLeft = n
			l.eating = true
			return l, tea.Batch(l.eat(), cmd)
		default:
			if !l.setting {
				return l, cmd
			}
			switch msg.String() {
			case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
				if l.setting || l.eating {
					err := errors.New("wait, please")
					l.SetError(err)
				} else {
					err := fmt.Errorf("only 1-%d cells can be eaten each turn", l.currentLevel().eatingMax)
					l.SetError(err)
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
	return l, cmd
}

func (l *last) view() string {
	l.buf.Reset()

	l.grid.Range(func(_ grid.Position, char grid.Rune, isLineEnd bool) (end bool) {
		l.buf.WriteString(l.charDic[rune(char)])
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
	return l.buf.String()
}

func (l *last) setLevel(i int) {
	l.setting = true // wait for the user to decide whether to get started first
	l.eatingPath = &pathStack{}
	l.levelIndex = i
	curLvl := l.currentLevel()
	l.commonCells = curLvl.totalCells - 2 // minus the 2 plays
	keys := []string{}
	for i := 1; i <= curLvl.eatingMax; i++ {
		keys = append(keys, strconv.Itoa(i))
	}
	numbersKey := key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(fmt.Sprintf("1-%d", curLvl.eatingMax), "cells to eatk"),
	)
	l.numbersKey = &numbersKey
	l.numbersKey.SetEnabled(false)
	l.ClearGroups()
	l.AddKeyGroup(game.KeyGroup{l.numbersKey})
	l.rd.Shuffle(len(playSyles), func(i, j int) {
		playSyles[i], playSyles[j] = playSyles[j], playSyles[i]
	})
	l.charDic = map[rune]string{
		blank: "     ",
		cell:  "  ◎  ",
		me:    playSyles[0].Render("  ◉  "),
		rival: playSyles[1].Render("  ◉  "),
	}
	l.genCells()
	l.playerIndex = 0
}

func (l *last) setted() {
	l.setting = false
	l.numbersKey.SetEnabled(true)
	ks := []string{"1", "2", "3", "4"}
	l.numbersKey.SetKeys(ks[:l.currentLevel().eatingMax]...)
	l.numbersKey.SetHelp(fmt.Sprintf("1-%d", l.currentLevel().eatingMax), "cells to eat")
}

func (l *last) genCells() {
	l.grid = grid.NewWithString("")
	g := make([][]grid.Rune, height)
	for i := range g {
		g[i] = make([]grid.Rune, width)
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
	l.helpGrid = l.grid.Copied()
	l.grid.Range(func(pos grid.Position, char grid.Rune, _ bool) (end bool) {
		switch char {
		case me:
			l.players[0] = pos
		case rival:
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
	l.grid.Range(func(pos grid.Position, char grid.Rune, _ bool) (end bool) {
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

func (l *last) removeCells(g *grid.Grid[grid.Rune], n int) {
	for ; n > 0; n-- {
		l.changeCell(g, cell, blank)
	}
}

func (l *last) addCells(g *grid.Grid[grid.Rune], n int) {
	for ; n > 0; n-- {
		l.changeCell(g, blank, cell)
	}
}

func (l *last) changeCell(g *grid.Grid[grid.Rune], from, to grid.Rune) {
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
	if l.success() {
		l.SetSuccess("Your are the last :)")
	} else if l.fail() {
		l.SetFailure("Your rival is the last :(")
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
	isOpposite := func(char grid.Rune) bool {
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
		l.numbersKey.SetEnabled(false)
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
