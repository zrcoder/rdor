package slide

import (
	icolor "image/color"
	"math/rand"
	"strings"
	"time"

	"github.com/zrcoder/rdor/pkg/game"
	"github.com/zrcoder/rdor/pkg/grid"
	"github.com/zrcoder/rdor/pkg/style/color"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	name = "Slippy Slide"

	size = 9

	player = 'p'
	target = 't'
	water  = 'w'
	land   = 'l'
	ice    = 'I'
	block  = 'b'
)

func New() game.Game {
	return &slide{Base: game.New(name)}
}

type slide struct {
	*game.Base
	rd          *rand.Rand
	bgStyle     lipgloss.Style
	grid        *grid.Grid[rune]
	helpGrid    *grid.Grid[rune]
	buf         *strings.Builder
	playerPos   grid.Position
	targetPos   grid.Position
	charViewDic map[rune]string
}

type tickMsg time.Time

func (s *slide) Init() tea.Cmd {
	s.RegisterLevels(1, s.setLevel)
	s.RegisterView(s.view)
	s.bgStyle = lipgloss.NewStyle().Background(color.IceBlue)
	s.rd = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	s.buf = &strings.Builder{}
	s.charViewDic = map[rune]string{
		player: colorStr(color.White, " ⛇ "),
		target: colorStr(color.Red, " ⚑ "),
		block:  colorStr(color.Gray, " ▲ "),
		water:  colorStr(color.LightBlue, " ⏺ "),
		ice:    "   ",
		land:   colorStr(color.LimeGreen, " ◉ "),
	}
	return s.Base.Init()
}

func colorStr(c icolor.Color, text string) string {
	return lipgloss.NewStyle().Foreground(c).Render(text)
}

func (s *slide) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	b, cmd := s.Base.Update(msg)
	if b != s.Base {
		return b, cmd
	}

	switch msg := msg.(type) {
	case tickMsg: // TODO
		return s, tea.Batch(cmd)
	case tea.KeyPressMsg:

		switch msg.String() {
		case "up":

		case "down":

		case "left":
		case "right":

		}
	}
	return s, cmd
}

func (s *slide) view() string {
	s.buf.Reset()
	s.grid.Range(func(_ grid.Position, char rune, isLineEnd bool) (end bool) {
		s.buf.WriteString(s.bgStyle.Render(s.charViewDic[char]))
		if isLineEnd {
			s.buf.WriteByte('\n')
		}
		return
	})
	return s.buf.String()
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
		g[i/size][i%size] = land
	}
	from, lim = lim, 2+size*2+size*2/3
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

func (s *slide) doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
