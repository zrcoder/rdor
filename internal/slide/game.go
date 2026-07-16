package slide

import (
	icolor "image/color"
	"math/rand"
	"time"

	"github.com/zrcoder/rdor/pkg/game"
	"github.com/zrcoder/rdor/pkg/grid"
	"github.com/zrcoder/rdor/pkg/style/color"

	tea "charm.land/bubbletea/v2"
	lg "charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

const (
	name = "Slippy Slide"

	size = 9

	player = 'p'
	target = 't'
	water  = 'w'
	ice    = 'I'
	block  = 'b'
)

type Dir int

const (
	DirUp Dir = iota
	DirLeft
	DirDown
	DirRight
)

func New() game.Game {
	return &slide{Base: game.New(name)}
}

type slide struct {
	*game.Base
	rd          *rand.Rand
	grid        *grid.Grid[rune]
	helpGrid    *grid.Grid[rune]
	table       table.Table
	playerPos   grid.Position
	targetPos   grid.Position
	charViewDic map[rune]string
}

type moveMsg = Dir

func (s *slide) Init() tea.Cmd {
	s.RegisterLevels(1, s.setLevel)
	s.RegisterView(s.view)
	s.table = *table.New().BorderRow(true).BaseStyle(lg.NewStyle().Background(color.IceBlue))
	s.rd = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	s.charViewDic = map[rune]string{
		player: fgCell(color.White, " ⛇ "),
		target: fgCell(color.Red, " ⚑ "),
		block:  bgCell(color.Gray),
		water:  bgCell(color.LightBlue),
		ice:    "   ",
	}
	s.Base.DisabledNextKey()
	s.Base.DisabledPrevKey()
	s.Base.DisabledSetKey()
	return s.Base.Init()
}

func fgCell(fg icolor.Color, text string) string {
	return lg.NewStyle().Foreground(fg).Render(text)
}

func bgCell(bg icolor.Color) string {
	return lg.NewStyle().Background(bg).Render("   ")
}

func (s *slide) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	b, cmd := s.Base.Update(msg)
	if b != s.Base {
		return b, cmd
	}

	switch msg := msg.(type) {
	case moveMsg:
		cmd = s.move(msg)
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			cmd = s.move(DirUp)
		case "left", "h":
			cmd = s.move(DirLeft)
		case "down", "j":
			cmd = s.move(DirDown)
		case "right", "l":
			cmd = s.move(DirRight)
		}
	}
	return s, cmd
}

func (s *slide) move(dir Dir) tea.Cmd {
	pos := s.playerPos
	switch dir {
	case DirUp:
		pos.Row--
	case DirLeft:
		pos.Col--
	case DirDown:
		pos.Row++
	case DirRight:
		pos.Col++
	}
	if pos.Row < 0 || pos.Row >= size || pos.Col < 0 || pos.Col >= size {
		return nil
	}
	dst := s.grid.Get(pos)
	if dst == block {
		return nil
	}
	if dst == water {
		s.Base.SetFailure("Failed")
		return nil
	}
	if dst == target {
		s.Base.SetSuccess("Done")
		return nil
	}
	s.grid.Set(s.playerPos, ice)
	s.playerPos = pos
	s.grid.Set(s.playerPos, player)
	time.Sleep(300 * time.Millisecond)
	return func() tea.Msg {
		return dir
	}
}

func (s *slide) view() string {
	s.table.ClearRows()
	s.grid.RangeRows(func(r int, row []rune, isLast bool) (end bool) {
		blocks := make([]string, len(row))
		for i, v := range row {
			blocks[i] = s.charViewDic[v]
		}
		s.table.Row(blocks...)
		return false
	})
	return s.table.String()
}

func (s *slide) setLevel(i int) {
	s.genScene()
}

func (s *slide) genScene() {
	s.grid = grid.NewWithString("")
	g := make([][]rune, size)
	for i := range g {
		g[i] = make([]rune, size)
	}
	g[0][0] = player
	g[0][1] = target
	from, lim := 2, 2+size
	for i := from; i < lim; i++ {
		g[i/size][i%size] = block
	}
	from, lim = lim, 2+size*2
	for i := from; i < lim; i++ {
		g[i/size][i%size] = water
	}
	from, lim = lim, size*size
	for i := from; i < lim; i++ {
		g[i/size][i%size] = ice
	}
	s.rd.Shuffle(lim, func(i, j int) {
		ir, ic := i/size, i%size
		jr, jc := j/size, j%size
		g[ir][ic], g[jr][jc] = g[jr][jc], g[ir][ic]
	})
	s.grid.SetData(g)
	s.helpGrid = s.grid.Copied()
	s.grid.Range(func(pos grid.Position, char rune, _ bool) (end bool) {
		switch char {
		case player:
			s.playerPos = pos
		case target:
			s.targetPos = pos
		}
		return
	})
}
