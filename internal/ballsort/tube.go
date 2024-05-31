package ballsort

import (
	lg "github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/rdor/pkg/style/color"
)

var (
	tubeStyle = lg.NewStyle().Border(lg.RoundedBorder(), false, true, true)
	doneStyle = lg.NewStyle().Foreground(color.Green)
)

type Ball struct {
	id int
}

func (b Ball) view() string {
	return ballStyles[b.id].Render("◉")
}

type Tube struct {
	holdTop bool
	balls   []*Ball
}

func NewTube(cap int) *Tube {
	return &Tube{balls: make([]*Ball, 0, cap)}
}

func (t *Tube) full() bool  { return len(t.balls) == cap(t.balls) }
func (t *Tube) empty() bool { return len(t.balls) == 0 }

func (t *Tube) Push(ball *Ball) {
	t.balls = append(t.balls, ball)
}

func (t *Tube) pop() *Ball {
	n := len(t.balls)
	x := t.balls[n-1]
	t.balls = t.balls[:n-1]
	return x
}

func (t *Tube) top() *Ball {
	return t.balls[len(t.balls)-1]
}

func (t *Tube) done() bool {
	if !t.full() {
		return false
	}
	for i := 1; i < len(t.balls); i++ {
		if t.balls[i].id != t.balls[i-1].id {
			return false
		}
	}
	return true
}

func (t *Tube) view() string {
	balls := t.balls
	topView := " "
	if t.holdTop {
		topView = t.top().view()
		balls = balls[:len(balls)-1]
	} else if t.done() {
		topView = doneStyle.Render("✓")
	}
	views := make([]string, 0, cap(t.balls))
	for n := cap(t.balls) - len(balls); n > 0; n-- {
		views = append(views, " ")
	}
	for i := len(balls) - 1; i >= 0; i-- {
		views = append(views, balls[i].view())
	}
	ballsView := lg.JoinVertical(lg.Center,
		views...)
	return lg.JoinVertical(lg.Center,
		topView,
		tubeStyle.Render(ballsView),
	)
}
