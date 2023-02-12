package internal

import "strings"

type disk struct {
	id   int
	view string
}

type pile struct {
	disks   []*disk
	overOne bool
}

func (s *pile) empty() bool {
	return len(s.disks) == 0
}
func (s *pile) push(d *disk) {
	s.disks = append(s.disks, d)
}
func (s *pile) pop() *disk {
	n := len(s.disks)
	res := s.disks[n-1]
	s.disks = s.disks[:n-1]
	return res
}
func (s *pile) top() *disk {
	n := len(s.disks)
	return s.disks[n-1]
}

func (p *pile) view() string {
	buf := strings.Builder{}
	disks := p.disks
	writeDisk := func() {
		top := disks[len(disks)-1]
		buf.WriteString(blanks((pileWidth-poleWidth-diskWidthUnit*top.id)/2 + horizontalSepBlanks))
		buf.WriteString(top.view)
		buf.WriteString(blanks((pileWidth - poleWidth - diskWidthUnit*top.id) / 2))
		disks = disks[:len(disks)-1]
	}
	if p.overOne {
		writeDisk()
	} else {
		buf.WriteString(blanks(horizontalSepBlanks + pileWidth))
	}
	buf.WriteByte('\n')
	for i := maxDisks; i > 0; i-- {
		if i == len(disks) {
			writeDisk()
		} else {
			buf.WriteString(blanks((pileWidth-poleWidth)/2 + horizontalSepBlanks))
			buf.WriteString(poleCh)
			buf.WriteString(blanks((pileWidth - poleWidth) / 2))
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}
