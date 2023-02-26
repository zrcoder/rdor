package sokoban

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Next     key.Binding
	Previous key.Binding
	Set      key.Binding
	Home     key.Binding
	Reset    key.Binding
}

func (k *keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Set, k.Home}
}

func (k *keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Left, k.Down, k.Right},
		{k.Reset, k.Next, k.Previous},
		{k.Set, k.Home},
	}
}

func getKeys() *keyMap {
	return &keyMap{
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		Left: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "left"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
		Right: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "right"),
		),
		Reset: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reset"),
		),
		Next: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "next"),
		),
		Previous: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "previous"),
		),
		Set: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "set"),
		),
		Home: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "home"),
		),
	}
}
