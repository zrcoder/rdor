package npuzzle

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up      key.Binding
	Left    key.Binding
	Down    key.Binding
	Right   key.Binding
	Shuffle key.Binding
	Next    key.Binding
	Home    key.Binding
}

func (k *keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Left, k.Down, k.Right, k.Shuffle, k.Next, k.Home}
}

func (k *keyMap) FullHelp() [][]key.Binding { return nil }

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
		Shuffle: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "shuffle"),
		),
		Next: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "next"),
		),
		Home: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "home"),
		),
	}
}
