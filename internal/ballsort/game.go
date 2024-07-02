package ballsort

import (
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/rdor/pkg/game"
	"github.com/zrcoder/rdor/pkg/style/color"
)

var (
	ballStyles    []lg.Style
	fullTubeNames = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}
)

const (
	name          = "Ball Sort"
	tubeCap       = 4
	defaultColors = 5
	emptyTubes    = 2
	levels        = 3
)

func New() game.Game {
	return &ballSort{Base: game.New(name)}
}

type ballSort struct {
	*game.Base
	rd        *rand.Rand
	buf       *strings.Builder
	tubes     map[string]*Tube
	overBall  *Ball
	balls     []*Ball
	tubeNames []string
	colors    int
}

func (p *ballSort) Init() tea.Cmd {
	ballStyles = []lg.Style{
		lg.NewStyle().Foreground(color.Red),
		lg.NewStyle().Foreground(color.Orange),
		lg.NewStyle().Foreground(color.Yellow),
		lg.NewStyle().Foreground(color.Green),
		lg.NewStyle().Foreground(color.Blue),
		lg.NewStyle().Foreground(color.Violet),
		lg.NewStyle(), // black/white as default
	}
	p.RegisterView(p.view)
	p.RegisterLevels(levels, p.set)
	p.DisabledPrevKey()
	p.DisabledSetKey()
	p.buf = &strings.Builder{}

	p.rd = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	p.set(0)
	return p.Base.Init()
}

func (p *ballSort) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	b, cmd := p.Base.Update(msg)
	if b != p.Base {
		return b, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := strings.ToUpper(msg.String())
		if tube, ok := p.tubes[key]; ok {
			p.pick(tube)
		}
	}
	return p, cmd
}

func (p *ballSort) view() string {
	views := make([]string, 0, len(p.tubeNames))
	for _, name := range p.tubeNames {
		views = append(views, lg.JoinVertical(lg.Center,
			p.tubes[name].view(),
			name,
		), " ")
	}
	return lg.JoinVertical(lg.Center,
		lg.JoinHorizontal(lg.Top, views...),
	)
}

func (p *ballSort) set(levle int) {
	p.colors = defaultColors + levle
	p.tubeNames = fullTubeNames[:p.colors+emptyTubes]
	p.rd.Shuffle(len(ballStyles), func(i, j int) {
		ballStyles[i], ballStyles[j] = ballStyles[j], ballStyles[i]
	})
	p.tubes = make(map[string]*Tube, len(p.tubeNames))
	for _, name := range p.tubeNames {
		p.tubes[name] = NewTube(tubeCap)
	}
	p.balls = make([]*Ball, 0, p.colors*tubeCap)
	for i := 0; i < p.colors; i++ {
		for j := 0; j < tubeCap; j++ {
			p.balls = append(p.balls, &Ball{id: i})
		}
	}
	p.rd.Shuffle(len(p.balls), func(i, j int) {
		p.balls[i], p.balls[j] = p.balls[j], p.balls[i]
	})
	i := 0
	for _, tube := range p.tubes {
		for i < len(p.balls) && !tube.full() {
			tube.Push(p.balls[i])
			i++
		}
		if i == len(p.balls) {
			break
		}
	}
}

func (p *ballSort) pick(tube *Tube) {
	if p.overBall == nil {
		if tube.empty() || tube.done() {
			return
		}
		p.overBall = tube.top()
		tube.holdTop = true
		return
	}

	if !tube.empty() && p.overBall == tube.top() {
		tube.holdTop = false
		p.overBall = nil
		return
	}
	if tube.empty() || !tube.full() && p.overBall.id == tube.top().id {
		old := p.getOverBallTube()
		tube.Push(old.pop())
		old.holdTop = false
		p.overBall = nil
		return
	}
	old := p.getOverBallTube()
	old.holdTop = false
	p.overBall = tube.top()
	tube.holdTop = true
}

func (b *ballSort) getOverBallTube() *Tube {
	for _, tube := range b.tubes {
		if !tube.empty() && tube.top() == b.overBall {
			return tube
		}
	}
	return nil
}
