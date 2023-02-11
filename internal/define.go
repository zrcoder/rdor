package internal

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/charmbracelet/lipgloss"
)

const (
	minDisks = 1
	maxDisks = 7

	height              = 11
	diskWidthUnit       = 4
	horizontalSepBlanks = 2
	successCh           = "★"
	helpInfo            = `Welcome to hanoi game!

The goal is to move all disks from pile 1 to pile 3.
Each time, we can pick a pile and move the top disk to another pile.
> Notice that we can only place a disk above a bigger one.
`

	poleCh     = "|"
	diskCh     = " "
	groundCh   = "o"
	pole1Label = "1"
	pole2Label = "2"
	pole3Label = "3"
)

var (
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
	helpStyle  = lipgloss.NewStyle().Foreground(green)
	errorStyle = lipgloss.NewStyle().Foreground(red)
)

func init() {
	rand.Shuffle(len(diskStyles), func(i, j int) {
		diskStyles[i], diskStyles[j] = diskStyles[j], diskStyles[i]
	})
}
