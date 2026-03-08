package style

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
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
	Help = lipgloss.NewStyle().Foreground(compat.AdaptiveColor{
		Light: lipgloss.Color("#909090"),
		Dark:  lipgloss.Color("#626262"),
	})
)
