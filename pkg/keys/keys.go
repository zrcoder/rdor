package keys

import "github.com/charmbracelet/bubbles/key"

var (
	Up = key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "up"),
	)
	Left = key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "left"),
	)
	Down = key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "down"),
	)
	Right = key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "right"),
	)
)
