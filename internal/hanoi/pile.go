package hanoi

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	poleWidth     = 1
	diskWidthUnit = 4

	poleCh   = "|"
	diskCh   = " "
	groundCh = "‾"
)

type pile struct {
	*hanoi
	name    string
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
	lines := make([]string, p.maxDisks+4)
	lines[0] = strings.Repeat(" ", p.maxDisks*diskWidthUnit)
	disks := p.disks
	writeDisk := func(i int) {
		lines[i] = disks[len(disks)-1].view
		disks = disks[:len(disks)-1]
	}
	if p.overOne {
		writeDisk(1)
	}
	for i := p.maxDisks; i > 0; i-- {
		j := p.maxDisks - i + 2
		if i == len(disks) {
			writeDisk(j)
		} else {
			lines[j] = poleCh
		}
	}
	lines[len(lines)-1] = p.name
	lines[len(lines)-2] = strings.Repeat(groundCh, p.maxDisks*diskWidthUnit)
	return lipgloss.NewStyle().Width(p.maxDisks * diskWidthUnit).Render(
		lipgloss.JoinVertical(lipgloss.Center, lines...),
	)
}
