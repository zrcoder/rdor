package hanoi

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/tgame/pkg/style"
	"github.com/zrcoder/tgame/pkg/util"
)

type errMsg error

type hanoi struct {
	disks    int
	piles    []*pile
	keys     keyMap
	keysHelp help.Model

	setting  bool
	showHelp bool
	buf      *strings.Builder
	steps    int
	err      error
	overDisk *disk
}

func New() *hanoi {
	return &hanoi{
		setting:  true,
		keys:     keys,
		keysHelp: help.New(),
		buf:      &strings.Builder{},
	}
}

func (h *hanoi) Init() tea.Cmd {
	h.keys.Piles.SetEnabled(false)
	h.keys.Disks.SetEnabled(true)
	h.keys.Reset.SetEnabled(false)
	h.keysHelp.ShowAll = true
	return nil
}

func (h *hanoi) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		h.err = msg
	case tea.KeyMsg:
		h.err = nil
		switch {
		case key.Matches(msg, h.keys.Quit):
			return h, tea.Quit
		case key.Matches(msg, h.keys.Help):
			h.showHelp = !h.showHelp
			h.keysHelp.ShowAll = !h.showHelp
			h.keys.Reset.SetEnabled(!h.setting && !h.showHelp)
			h.keys.Disks.SetEnabled(h.setting && !h.showHelp)
			h.keys.Piles.SetEnabled(!h.setting && !h.showHelp)
		case key.Matches(msg, h.keys.Reset):
			h.setting = true
			return h, h.Init()
		case key.Matches(msg, h.keys.Disks):
			n, _ := strconv.Atoi(msg.String())
			h.setted(n)
		case key.Matches(msg, h.keys.Piles):
			return h, h.pick(msg.String())
		default:
			if h.setting {
				h.err = errDiskNum
			}
		}
	}
	return h, nil
}

func (h *hanoi) View() string {
	h.buf.Reset()
	if h.showHelp {
		h.writeHelpInfo()
	} else {
		h.writeHead()
		if h.setting {
			h.writeSettingView()
		} else {
			h.writePoles()
			h.writeGround()
			h.writeLabels()
			h.writeState()
		}
	}
	h.writeKeysHelp()
	h.writeBlankLine()
	return h.buf.String()
}

func (h *hanoi) setted(n int) {
	h.setting = false
	h.disks = n
	h.steps = 0
	h.overDisk = nil
	h.err = nil
	h.showHelp = false
	h.keys.Disks.SetEnabled(false)
	h.keys.Piles.SetEnabled(true)
	h.keys.Reset.SetEnabled(true)
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

func (h *hanoi) pick(key string) tea.Cmd {
	if h.success() {
		return nil
	}
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
	return func() tea.Msg {
		if h.overDisk == nil && curPile.empty() {
			return nil
		}
		if h.overDisk == nil {
			curPile.overOne = true
			h.overDisk = curPile.top()
			return nil
		}
		if !curPile.empty() && h.overDisk.id > curPile.top().id {
			return errMsg(errCantMove)
		}
		if !curPile.empty() && h.overDisk == curPile.top() {
			curPile.overOne = false
			h.overDisk = nil
			return nil
		}
		for _, p := range h.piles {
			if p.overOne {
				h.steps++
				curPile.push(p.pop())
				p.overOne = false
				h.overDisk = nil
			}
		}
		return nil
	}
}

func (h *hanoi) writeHead() {
	h.buf.WriteString(healpHead)
}
func (h *hanoi) success() bool {
	last := h.piles[len(h.piles)-1]
	return len(last.disks) == h.disks
}

func (h *hanoi) writeSettingView() {
	h.writeLine(settingHint)
	if h.err != nil {
		h.writeError(h.err)
	}
	h.writeBlankLine()
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
	h.buf.WriteString(util.Blanks(n))
	h.buf.WriteString(pole1Label)
	n = (pileWidth-len(pole1Label))/2 + horizontalSepBlanks + (pileWidth-len(pole2Label))/2
	h.buf.WriteString(util.Blanks(n))
	h.buf.WriteString(pole2Label)
	n = (pileWidth-len(pole2Label))/2 + horizontalSepBlanks + (pileWidth-len(pole3Label))/2
	h.buf.WriteString(util.Blanks(n))
	h.buf.WriteString(pole3Label)
	h.writeBlankLine()
	h.writeBlankLine()
}

func (h *hanoi) writeState() {
	if h.success() {
		minSteps := 1<<h.disks - 1
		totalStart := 5
		if h.steps == minSteps {
			h.buf.WriteString(style.Success.Render("Fantastic! you earned all the stars! "))
			h.buf.WriteString(starStyle.Render(strings.Repeat(starCh, totalStart)))
		} else {
			s := fmt.Sprintf("Done! Taken %d steps, can you complete it in %d step(s)? ", h.steps, minSteps)
			h.buf.WriteString(style.Success.Render(s))
			stars := 3
			if h.steps-minSteps > minSteps/2 {
				stars = 1
			}
			s = strings.Repeat(starCh, stars) + strings.Repeat(starOutlineCh, totalStart-stars)
			h.buf.WriteString(starStyle.Render(s))
		}
		h.writeBlankLine()
	} else if h.err != nil {
		h.writeError(h.err)
		h.writeBlankLine()
	} else {
		h.writeLine(fmt.Sprintf("steps: %d", h.steps))
	}
	h.writeBlankLine()
}

func (h *hanoi) writeKeysHelp() {
	h.buf.WriteString(h.keysHelp.View(h.keys))
	h.writeBlankLine()
}

func (h *hanoi) writeHelpInfo() {
	h.buf.WriteString(helpInfo)
}

func (h *hanoi) writeBlankLine() {
	h.buf.WriteByte('\n')
}

func (h *hanoi) writeError(err error) {
	h.buf.WriteString(style.Error.Render(err.Error()))
}

func (h *hanoi) writeLine(s string) {
	h.buf.WriteString(s)
	h.writeBlankLine()
}

type disk struct {
	id   int
	view string
}

type pile struct {
	disks   []*disk
	overOne bool
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
		buf.WriteString(util.Blanks((pileWidth-poleWidth-diskWidthUnit*top.id)/2 + horizontalSepBlanks))
		buf.WriteString(top.view)
		buf.WriteString(util.Blanks((pileWidth - poleWidth - diskWidthUnit*top.id) / 2))
		disks = disks[:len(disks)-1]
	}
	if p.overOne {
		writeDisk()
	} else {
		buf.WriteString(util.Blanks(horizontalSepBlanks + pileWidth))
	}
	buf.WriteByte('\n')
	for i := maxDisks; i > 0; i-- {
		if i == len(disks) {
			writeDisk()
		} else {
			buf.WriteString(util.Blanks((pileWidth-poleWidth)/2 + horizontalSepBlanks))
			buf.WriteString(poleCh)
			buf.WriteString(util.Blanks((pileWidth - poleWidth) / 2))
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
