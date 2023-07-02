package color

import "github.com/charmbracelet/lipgloss"

var (
	Faint = lipgloss.AdaptiveColor{
		Light: "#D9DCCF",
		Dark:  "#383838",
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
