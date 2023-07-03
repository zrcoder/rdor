package hanoi

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/zrcoder/rdor/pkg/game"
	"github.com/zrcoder/rdor/pkg/style/color"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	name                = "Hanoi"
	diskWidthUnit       = 4
	horizontalSepBlanks = 1
	poleWidth           = 1

	poleCh     = "|"
	diskCh     = " "
	groundCh   = "‾"
	pole1Label = "1"
	pole2Label = "2"
	pole3Label = "3"
)

var (
	errCantMove = errors.New("can not move the disk above a smaller one")
)

func New() game.Game {
	return &hanoi{Base: game.New(name)}
}

type hanoi struct {
	*game.Base
	levels     []int
	maxDisks   int
	pileWidth  int
	diskStyles []lipgloss.Style
	disks      int
	overDisk   *disk
	buf        *strings.Builder
	pilesKey   *key.Binding
	piles      []*pile
	steps      int
}

type disk struct {
	view string
	id   int
}

type pile struct {
	*hanoi
	disks   []*disk
	overOne bool
}

func (h *hanoi) Init() tea.Cmd {
	h.levels = []int{2, 3, 4, 5, 6, 7}
	h.maxDisks = h.levels[len(h.levels)-1]
	h.pileWidth = diskWidthUnit*(h.maxDisks) + poleWidth
	h.diskStyles = []lipgloss.Style{
		lipgloss.NewStyle().Background(color.Red),
		lipgloss.NewStyle().Background(color.Orange),
		lipgloss.NewStyle().Background(color.Yellow),
		lipgloss.NewStyle().Background(color.Green),
		lipgloss.NewStyle().Background(color.Blue),
		lipgloss.NewStyle().Background(color.Indigo),
		lipgloss.NewStyle().Background(color.Violet),
	}
	h.RegisterView(h.view)
	h.RegisterHelp(h.helpInfo)
	pilesKey := key.NewBinding(
		key.WithKeys("1", "2", "3", "j", "k", "l"),
		key.WithHelp("1-3/j,k,l", "pick a pile"),
	)
	h.pilesKey = &pilesKey
	h.ClearGroups()
	h.AddKeyGroup(game.KeyGroup{h.pilesKey})
	h.RegisterLevels(len(h.levels), h.setted)
	h.buf = &strings.Builder{}
	return h.Base.Init()
}

func (h *hanoi) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	b, cmd := h.Base.Update(msg)
	if b != h.Base { // is parent
		return b, cmd
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, *h.pilesKey):
			h.pick(msg.String())
			if h.success() {
				h.setSuccessView()
			}
		}
	}
	return h, cmd // b is base or parent
}

func (h *hanoi) view() string {
	h.buf.Reset()
	h.writePoles()
	h.writeGround()
	h.writeLabels()
	h.writeState()
	return h.buf.String()
}

func (h *hanoi) helpInfo() string {
	return "Our goal is to move all disks from pile `1` to pile `3`."
}

func (h *hanoi) setSuccessView() {
	minSteps := 1<<h.disks - 1
	totalStars := 5
	if h.steps == minSteps {
		h.SetSuccess("Fantastic! you earned all the stars!")
		h.SetStars(totalStars, totalStars)
		return
	}
	s := fmt.Sprintf("Done! Taken %d steps, can you complete it in %d step(s)? ", h.steps, minSteps)
	stars := 3
	if h.steps-minSteps > minSteps/2 {
		stars = 1
	}
	h.SetSuccess(s)
	h.SetStars(totalStars, stars)
}

func (h *hanoi) setted(level int) {
	h.disks = h.levels[level]
	h.steps = 0
	h.overDisk = nil
	h.piles = make([]*pile, 3)
	for i := range h.piles {
		h.piles[i] = &pile{hanoi: h}
	}
	h.shuffleDiskStyles()
	disks := make([]*disk, h.disks)
	for i := 1; i <= h.disks; i++ {
		disks[h.disks-i] = &disk{
			id:   i,
			view: h.diskStyles[i-1].Render(strings.Repeat(diskCh, i*diskWidthUnit)),
		}
	}
	h.piles[0].disks = disks
}

func (h *hanoi) pick(key string) {
	idx := map[string]int{
		"1": 0,
		"j": 0,
		"2": 1,
		"k": 1,
		"3": 2,
		"l": 2,
	}
	i := idx[key]
	curPile := h.piles[i]
	if h.overDisk == nil && curPile.empty() {
		return
	}
	if h.overDisk == nil {
		curPile.overOne = true
		h.overDisk = curPile.top()
		return
	}
	if !curPile.empty() && h.overDisk.id > curPile.top().id {
		h.SetError(errCantMove)
		return
	}
	if !curPile.empty() && h.overDisk == curPile.top() {
		curPile.overOne = false
		h.overDisk = nil
		return
	}
	for _, p := range h.piles {
		if p.overOne {
			h.steps++
			curPile.push(p.pop())
			p.overOne = false
			h.overDisk = nil
		}
	}
}

func (h *hanoi) writePoles() {
	views := make([]string, len(h.piles))
	for i, p := range h.piles {
		views[i] = p.view()
	}
	poles := lipgloss.JoinHorizontal(
		lipgloss.Top,
		views...,
	)
	h.buf.WriteString(poles)
	h.writeBlankLine()
}

func (h *hanoi) writeGround() {
	h.buf.WriteString(strings.Repeat(groundCh, (h.pileWidth*3 + horizontalSepBlanks*4)))
	h.writeBlankLine()
}

func (h *hanoi) writeLabels() {
	n := horizontalSepBlanks + (h.pileWidth-len(pole1Label))/2
	h.buf.WriteString(strings.Repeat(" ", n))
	h.buf.WriteString(pole1Label)
	n = (h.pileWidth-len(pole1Label))/2 + horizontalSepBlanks + (h.pileWidth-len(pole2Label))/2
	h.buf.WriteString(strings.Repeat(" ", n))
	h.buf.WriteString(pole2Label)
	n = (h.pileWidth-len(pole2Label))/2 + horizontalSepBlanks + (h.pileWidth-len(pole3Label))/2
	h.buf.WriteString(strings.Repeat(" ", n))
	h.buf.WriteString(pole3Label)
	h.writeBlankLine()
	h.writeBlankLine()
}

func (h *hanoi) writeState() {
	h.writeLine(fmt.Sprintf("steps: %d\n", h.steps))
}

func (h *hanoi) writeLine(s string) {
	h.buf.WriteString(s)
	h.writeBlankLine()
}

func (h *hanoi) writeBlankLine() {
	h.buf.WriteByte('\n')
}

func (h *hanoi) success() bool {
	last := h.piles[len(h.piles)-1]
	return len(last.disks) == h.disks
}

func (p *pile) empty() bool {
	return len(p.disks) == 0
}

func (p *pile) push(d *disk) {
	p.disks = append(p.disks, d)
}

func (p *pile) pop() *disk {
	n := len(p.disks)
	res := p.disks[n-1]
	p.disks = p.disks[:n-1]
	return res
}

func (p *pile) top() *disk {
	n := len(p.disks)
	return p.disks[n-1]
}

func (p *pile) view() string {
	buf := strings.Builder{}
	disks := p.disks
	writeDisk := func() {
		top := disks[len(disks)-1]
		buf.WriteString(strings.Repeat(" ", (p.pileWidth-poleWidth-diskWidthUnit*top.id)/2+horizontalSepBlanks))
		buf.WriteString(top.view)
		buf.WriteString(strings.Repeat(" ", (p.pileWidth-poleWidth-diskWidthUnit*top.id)/2))
		disks = disks[:len(disks)-1]
	}
	if p.overOne {
		writeDisk()
	} else {
		buf.WriteString(strings.Repeat(" ", horizontalSepBlanks+p.pileWidth))
	}
	buf.WriteByte('\n')
	for i := p.maxDisks; i > 0; i-- {
		if i == len(disks) {
			writeDisk()
		} else {
			buf.WriteString(strings.Repeat(" ", (p.pileWidth-poleWidth)/2+horizontalSepBlanks))
			buf.WriteString(poleCh)
			buf.WriteString(strings.Repeat(" ", (p.pileWidth-poleWidth)/2))
		}
		if i > 1 {
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}

func (h *hanoi) shuffleDiskStyles() {
	rand.Shuffle(len(h.diskStyles), func(i, j int) {
		h.diskStyles[i], h.diskStyles[j] = h.diskStyles[j], h.diskStyles[i]
	})
}
