package internal

import (
	_ "embed"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const (
	minDisks = 1
	maxDisks = 7

	height              = 11
	diskWidthUnit       = 4
	horizontalSepBlanks = 2
	starCh              = "★"
	starOutlineCh       = "☆"
	poleCh              = "|"
	diskCh              = " "
	groundCh            = "‾"
	pole1Label          = "1"
	pole2Label          = "2"
	pole3Label          = "3"

	settingHint = "How many disks do you like?"
)

var (
	//go:embed help.md
	helpInfo  string
	healpHead string

	poleWidth = 1
	pileWidth = diskWidthUnit*maxDisks + poleWidth

	errDiskNum  = fmt.Errorf("disks number must be an integer between %d to %d", minDisks, maxDisks)
	errCantMove = errors.New("can not move the disk above a smaller one")
)

var (
	// rainbow colors
	red    = lipgloss.Color("#ff0000")
	orange = lipgloss.Color("#ffa500")
	yellow = lipgloss.Color("#ffff00")
	green  = lipgloss.Color("#008000")
	blue   = lipgloss.Color("#0000ff")
	indigo = lipgloss.Color("#4b0082")
	violet = lipgloss.Color("#ee82ee")

	diskStyles = [maxDisks]lipgloss.Style{
		lipgloss.NewStyle().Background(red),
		lipgloss.NewStyle().Background(orange),
		lipgloss.NewStyle().Background(yellow),
		lipgloss.NewStyle().Background(green),
		lipgloss.NewStyle().Background(blue),
		lipgloss.NewStyle().Background(indigo),
		lipgloss.NewStyle().Background(violet),
	}

	starStyle  = lipgloss.NewStyle().Foreground(orange)
	infoStyle  = lipgloss.NewStyle().Foreground(green)
	errorStyle = lipgloss.NewStyle().Foreground(red)
)

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
		if i > 1 {
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}

func init() {
	md := mdRender()
	helpInfo = strings.TrimSpace(helpInfo)
	i := strings.Index(helpInfo, "\n")
	healpHead, _ = md.Render(helpInfo[:i+1])
	helpInfo, _ = md.Render(helpInfo)
}

func shuffleDiskStyles() {
	rand.Shuffle(len(diskStyles), func(i, j int) {
		diskStyles[i], diskStyles[j] = diskStyles[j], diskStyles[i]
	})
}

func mdRender() *glamour.TermRenderer {
	styleConfig := glamour.DarkStyleConfig
	var noMargin uint = 0
	styleConfig.Document.Margin = &noMargin
	render, _ := glamour.NewTermRenderer(glamour.WithStyles(styleConfig))
	return render
}

func printError(err error) {
	fmt.Println(errorStyle.Render(err.Error()))
}

func blanks(n int) string {
	return strings.Repeat(" ", n)
}
