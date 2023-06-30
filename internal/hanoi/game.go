package hanoi

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/zrcoder/rdor/internal/internal"
	"github.com/zrcoder/rdor/pkg/model"
	"github.com/zrcoder/rdor/pkg/style"
	"github.com/zrcoder/rdor/pkg/style/color"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type hanoi struct {
	*internal.Game
	helpInfo string
	disks    int
	piles    []*pile
	overDisk *disk
	pilesKey key.Binding
	buf      *strings.Builder
	steps    int
}

func New() model.Game {
	g := internal.New(Name)
	res := &hanoi{Game: g}
	g.InitFunc = res.initialize
	g.UpdateFunc = res.update
	g.ViewFunc = res.view
	res.pilesKey = key.NewBinding(
		key.WithKeys("1", "2", "3", "j", "k", "l"),
		key.WithHelp("1-3/j,k,l", "pick a pile"),
	)
	g.SetExtraKeys([]key.Binding{res.pilesKey})
	g.KeyFuncReset = res.reset
	g.KeyFuncPrevious = res.previousLevel
	g.KeyFuncNext = res.nextLevel
	return res
}
func (h *hanoi) SetParent(parent tea.Model) { h.Parent = parent }

func (h *hanoi) reset() {
	h.setted(h.disks)
}
func (h *hanoi) previousLevel() {
	if h.disks-1 > 0 {
		h.setted(h.disks - 1)
	}
}
func (h *hanoi) nextLevel() {
	if h.disks+1 <= maxDisks {
		h.setted(h.disks + 1)
	}

}

type disk struct {
	id   int
	view string
}

type pile struct {
	disks   []*disk
	overOne bool
}

const (
	Name         = "Hanoi"
	minDisks     = 1
	maxDisks     = 7
	defaultDisks = 3

	diskWidthUnit       = 4
	horizontalSepBlanks = 2
	poleWidth           = 1

	poleCh     = "|"
	diskCh     = " "
	groundCh   = "‾"
	pole1Label = "1"
	pole2Label = "2"
	pole3Label = "3"
)

var (
	pileWidth   = diskWidthUnit*maxDisks + poleWidth
	errDiskNum  = fmt.Errorf("disks number must be an integer between %d to %d", minDisks, maxDisks)
	errCantMove = errors.New("can not move the disk above a smaller one")

	diskStyles = [maxDisks]lipgloss.Style{
		lipgloss.NewStyle().Background(color.Red),
		lipgloss.NewStyle().Background(color.Orange),
		lipgloss.NewStyle().Background(color.Yellow),
		lipgloss.NewStyle().Background(color.Green),
		lipgloss.NewStyle().Background(color.Blue),
		lipgloss.NewStyle().Background(color.Indigo),
		lipgloss.NewStyle().Background(color.Violet),
	}
)

func (h *hanoi) initialize() tea.Cmd {
	h.helpInfo = style.Help.Render("Our goal is to move all disks from pile `1` to pile `3`.")
	h.buf = &strings.Builder{}
	h.setted(defaultDisks)
	return nil
}

func (h *hanoi) update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, h.pilesKey):
			h.pick(msg.String())
			if h.success() {
				h.setSuccessView()
			}
		}
	}
	return nil
}

func (h *hanoi) view() string {
	h.buf.Reset()
	h.writePoles()
	h.writeGround()
	h.writeLabels()
	// TODO
	h.writeState()
	h.writeHelpInfo()
	return h.buf.String()
}

func (h *hanoi) setSuccessView() {
	minSteps := 1<<h.disks - 1
	totalStars := 5
	if h.steps == minSteps {
		h.SetSuccess("Fantastic! you earned all the stars!")
		h.SetStars(totalStars, totalStars)
		// TODO	Height(10).
		return
	}
	s := fmt.Sprintf("Done! Taken %d steps, can you complete it in %d step(s)? ", h.steps, minSteps)
	stars := 3
	if h.steps-minSteps > minSteps/2 {
		stars = 1
	}
	h.SetSuccess(s)
	h.SetStars(totalStars, stars)
	// TODO	Height(11).
}

func (h *hanoi) setted(n int) {
	h.disks = n
	h.steps = 0
	h.overDisk = nil
	h.piles = make([]*pile, 3)
	for i := range h.piles {
		h.piles[i] = &pile{}
	}
	shuffleDiskStyles()
	disks := make([]*disk, n)
	for i := 1; i <= n; i++ {
		disks[n-i] = &disk{
			id:   i,
			view: diskStyles[i-1].Render(strings.Repeat(diskCh, i*diskWidthUnit)),
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
	h.buf.WriteString(strings.Repeat(groundCh, (pileWidth*3 + horizontalSepBlanks*4)))
	h.writeBlankLine()
}

func (h *hanoi) writeLabels() {
	n := horizontalSepBlanks + (pileWidth-len(pole1Label))/2
	h.buf.WriteString(strings.Repeat(" ", n))
	h.buf.WriteString(pole1Label)
	n = (pileWidth-len(pole1Label))/2 + horizontalSepBlanks + (pileWidth-len(pole2Label))/2
	h.buf.WriteString(strings.Repeat(" ", n))
	h.buf.WriteString(pole2Label)
	n = (pileWidth-len(pole2Label))/2 + horizontalSepBlanks + (pileWidth-len(pole3Label))/2
	h.buf.WriteString(strings.Repeat(" ", n))
	h.buf.WriteString(pole3Label)
	h.writeBlankLine()
	h.writeBlankLine()
}

func (h *hanoi) writeState() {
	h.writeLine(fmt.Sprintf("steps: %d\n", h.steps))
}

func (h *hanoi) writeHelpInfo() {
	h.buf.WriteString(h.helpInfo + "\n\n")
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
		buf.WriteString(strings.Repeat(" ", (pileWidth-poleWidth-diskWidthUnit*top.id)/2+horizontalSepBlanks))
		buf.WriteString(top.view)
		buf.WriteString(strings.Repeat(" ", (pileWidth-poleWidth-diskWidthUnit*top.id)/2))
		disks = disks[:len(disks)-1]
	}
	if p.overOne {
		writeDisk()
	} else {
		buf.WriteString(strings.Repeat(" ", horizontalSepBlanks+pileWidth))
	}
	buf.WriteByte('\n')
	for i := maxDisks; i > 0; i-- {
		if i == len(disks) {
			writeDisk()
		} else {
			buf.WriteString(strings.Repeat(" ", (pileWidth-poleWidth)/2+horizontalSepBlanks))
			buf.WriteString(poleCh)
			buf.WriteString(strings.Repeat(" ", (pileWidth-poleWidth)/2))
		}
		if i > 1 {
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}

func shuffleDiskStyles() {
	rand.Shuffle(len(diskStyles), func(i, j int) {
		diskStyles[i], diskStyles[j] = diskStyles[j], diskStyles[i]
	})
}
