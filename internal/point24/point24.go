package point24

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/rdor/pkg/game"
	"github.com/zrcoder/rdor/pkg/keyblock"
)

const (
	Name = "24 Points"

	plus  = "+"
	minus = "-"
	times = "×"
	divid = "÷"
)

type point24 struct {
	*game.Base
	levels   [][4]int
	curLevel int
	point    int
	nums     keyblock.KeysLine
	opers    keyblock.KeysLine

	oper string
	num  int

	a, s, d, f, h, j, k, l *keyblock.Key
}

func New() game.Game {
	return &point24{Base: game.New(Name)}
}

func (p *point24) Init() tea.Cmd {
	p.ViewFunc = p.view
	p.levels = getLevers()
	p.KeyActionNext = func() {
		if p.curLevel < len(p.levels)-1 {
			p.setLever(p.curLevel + 1)
		}
	}
	p.KeyActionPrevious = func() {
		if p.curLevel > 0 {
			p.setLever(p.curLevel - 1)
		}
	}
	p.KeyActionReset = func() { p.setLever(p.curLevel) }
	p.a = keyblock.New("A")
	p.s = keyblock.New("S")
	p.d = keyblock.New("D")
	p.f = keyblock.New("F")
	p.h = keyblock.New("H")
	p.j = keyblock.New("J")
	p.k = keyblock.New("K")
	p.l = keyblock.New("L")
	p.h.Display = plus
	p.j.Display = minus
	p.k.Display = times
	p.l.Display = divid
	p.nums = []*keyblock.Key{p.a, p.s, p.d, p.f}
	p.opers = []*keyblock.Key{p.h, p.j, p.k, p.l}
	p.setLever(p.curLevel)
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
	p.curLevel = i
	level := p.levels[i]
	for j, v := range level {
		p.nums[j].SetNumber(v)
	}
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
