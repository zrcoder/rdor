package point24

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/rdor/pkg/game"
	"github.com/zrcoder/rdor/pkg/keyblock"
)

const (
	name = "24 Points"

	plus  = "+"
	minus = "-"
	times = "×"
	divid = "÷"
)

type point24 struct {
	*game.Base
	levels [][4]int
	oper   string
	nums   keyblock.KeysLine
	opers  keyblock.KeysLine
	num    int
	point  int
}

func New() game.Game {
	return &point24{Base: game.New(name)}
}

func (p *point24) Init() tea.Cmd {
	p.levels = getLevers()
	p.RegisterView(p.view)
	p.RegisterView(p.view)
	p.RegisterLevels(len(p.levels), p.setLever)
	p.nums = keyblock.NewKeysLine("A", "S", "D", "F")
	p.opers = keyblock.NewKeysLine("H", "J", "K", "L")
	p.opers.SetDisplays(plus, minus, times, divid)
	p.num = -1
	return p.Base.Init()
}

func (p *point24) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	b, cmd := p.Base.Update(msg)
	if b != p.Base {
		return b, cmd
	}
	switch msg := msg.(type) {
	case keyblock.PressMsg:
		if msg.IsNumber() {
			if msg.Pressed {
				p.caculate(msg.Number)
			}
		} else {
			if p.num != -1 {
				p.oper = msg.Display
			}
		}
	}
	_, numCmd := p.nums.Update(msg)
	_, operCmd := p.opers.Update(msg)
	return p, tea.Batch(cmd, numCmd, operCmd)
}

func (p *point24) view() string {
	keysView := lipgloss.JoinHorizontal(lipgloss.Center,
		p.nums.View(),
		"        ",
		p.opers.View(),
	)
	if p.num == -1 {
		return keysView
	}

	resView := strconv.Itoa(p.num)
	if p.oper != "" {
		resView += " " + p.oper
	}
	return lipgloss.JoinVertical(lipgloss.Center,
		keysView,
		"",
		resView,
	)
}

func (p *point24) setLever(i int) {
	level := p.levels[i]
	p.nums.SetNumbers(level[:]...)
}

func (p *point24) caculate(num int) {
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
