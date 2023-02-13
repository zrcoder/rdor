package internal

import (
	_ "embed"
	"errors"
	"fmt"
	"math/rand"

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
	//go:embed head.md
	head string
	//go:embed helpinfo.md
	helpInfo string

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

func init() {
	rand.Shuffle(len(diskStyles), func(i, j int) {
		diskStyles[i], diskStyles[j] = diskStyles[j], diskStyles[i]
	})
	md := mdRender()
	head, _ = md.Render(head)
	helpInfo, _ = md.Render(helpInfo)
}

func mdRender() *glamour.TermRenderer {
	styleConfig := glamour.DarkStyleConfig
	var noMargin uint = 0
	styleConfig.Document.Margin = &noMargin
	render, _ := glamour.NewTermRenderer(glamour.WithStyles(styleConfig))
	return render
}
