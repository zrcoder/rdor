package point24

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/rdor/pkg/game"
	"github.com/zrcoder/rdor/pkg/keyblock"
	"github.com/zrcoder/rdor/pkg/style/color"
)

const (
	name  = "24 Points"
	dest  = 24
	plus  = "+"
	minus = "-"
	times = "×"
	divid = "÷"
)

var (
	resStyle     = lg.NewStyle().Foreground(color.Orange).Border(lg.NormalBorder(), false, false, true, false)
	successStyle = lg.NewStyle().Foreground(color.Green).Border(lg.NormalBorder(), false, false, true, false)
)

type point24 struct {
	*game.Base
	levels [][4]int
	oper   string
	nums   keyblock.KeysLine
	opers  keyblock.KeysLine
	num    int
}

func New() game.Game {
	return &point24{Base: game.New(name)}
}

func (p *point24) Init() tea.Cmd {
	p.levels = getLevers()
	p.RegisterView(p.view)
	p.RegisterLevels(len(p.levels), p.setLever)
	p.DisabledSetKey()

	return p.Base.Init()
}

func (p *point24) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	b, cmd := p.Base.Update(msg)
	if b != p.Base {
		return b, cmd
	}
	p.nums.Update(msg)
	p.opers.Update(msg)
	return p, cmd
}

func (p *point24) view() string {
	keysView := lg.JoinHorizontal(lg.Center,
		p.nums.View(),
		"        ",
		p.opers.View(),
	)
	if p.num == -1 {
		return keysView
	}

	resView := strconv.Itoa(p.num)
	if p.oper != "" {
		resView += p.oper
	}
	if p.num == 24 {
		resView += "!"
	}
	switch len(resView) {
	case 1:
		resView = "  " + resView + "  "
	case 2:
		resView = " " + resView + "  "
	default:
		resView = " " + resView + " "
	}
	if p.num == dest {
		resView = successStyle.Render(resView)
	} else {
		resView = resStyle.Render(resView)
	}
	return lg.JoinVertical(lg.Center,
		keysView,
		"",
		resView,
	)
}

func (p *point24) setLever(i int) {
	p.nums = keyblock.NewKeysLine("a", "s", "d", "f")
	p.nums.SetOnce(true)
	p.nums.SetAction(p.numAction)
	p.opers = keyblock.NewKeysLine("h", "j", "k", "l")
	p.opers.SetOnce(true)
	p.opers.SetDisplays(plus, minus, times, divid)
	p.opers.SetAction(p.operAction)
	p.num = -1
	level := p.levels[i]
	for i, v := range level {
		p.nums.SetDisplay(i, strconv.Itoa(v))
	}
}

func (p *point24) operAction(key *keyblock.Key) {
	p.oper = key.Display
}

func (p *point24) numAction(key *keyblock.Key) {
	num, err := strconv.Atoi(key.Display)
	if err != nil {
		p.SetError(err)
		return
	}
	if p.oper == "" {
		p.num = num
		return
	}
	switch p.oper {
	case plus:
		p.num += num
	case minus:
		p.num -= num
	case times:
		p.num *= num
	case divid:
		p.num /= num
	}
	p.oper = ""
}
