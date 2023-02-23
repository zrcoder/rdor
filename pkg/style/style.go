package style

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/rdor/pkg/style/color"
)

var (
	Success = lipgloss.NewStyle().Foreground(color.Green)
	Error   = lipgloss.NewStyle().Foreground(color.Red)
	Warn    = lipgloss.NewStyle().Foreground(color.Orange)
	Title   = lipgloss.NewStyle().Background(color.Blue).
		Foreground(color.White).
		PaddingLeft(1).
		PaddingRight(1)
	Help = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#909090",
		Dark:  "#626262",
	})
)
