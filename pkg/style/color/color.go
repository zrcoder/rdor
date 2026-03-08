package color

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/compat"
)

var (
	Faint = compat.AdaptiveColor{
		Light: lipgloss.Color("#D9DCCF"),
		Dark:  lipgloss.Color("#383838"),
	}
	White = lipgloss.Color("#ffffff")
	// rainbow colors
	Red    = lipgloss.Color("#ff0000")
	Orange = lipgloss.Color("#ffa500")
	Yellow = lipgloss.Color("#ffff00")
	Green  = lipgloss.Color("#008000")
	Blue   = lipgloss.Color("#0000ff")
	Indigo = lipgloss.Color("#4b0082")
	Violet = lipgloss.Color("#ee82ee")
)
