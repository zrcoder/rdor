package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/zrcoder/tgame/pkg/print"
	"github.com/zrcoder/tgame/pkg/style/color"
	"github.com/zrcoder/tgame/pkg/util"

	tea "github.com/charmbracelet/bubbletea"
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
	diskStyles = [maxDisks]lipgloss.Style{
		lipgloss.NewStyle().Background(color.Red),
		lipgloss.NewStyle().Background(color.Orange),
		lipgloss.NewStyle().Background(color.Yellow),
		lipgloss.NewStyle().Background(color.Green),
		lipgloss.NewStyle().Background(color.Blue),
		lipgloss.NewStyle().Background(color.Indigo),
		lipgloss.NewStyle().Background(color.Violet),
	}

	starStyle = lipgloss.NewStyle().Foreground(color.Orange)
)

func init() {
	md := util.GetMarkdowdRender()
	helpInfo = strings.TrimSpace(helpInfo)
	i := strings.Index(helpInfo, "\n")
	healpHead, _ = md.Render(helpInfo[:i+1])
	helpInfo, _ = md.Render(helpInfo)
}

func main() {
	if _, err := tea.NewProgram(New()).Run(); err != nil {
		print.Errorln(err)
		os.Exit(1)
	}
}
