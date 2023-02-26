package last

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Numbers  key.Binding
	Next     key.Binding
	Previous key.Binding
	Reset    key.Binding
	Home     key.Binding
	Help     key.Binding
}

func (k *keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Numbers, k.Next, k.Previous, k.Reset, k.Home, k.Help}
}

func (k *keyMap) FullHelp() [][]key.Binding {
	return nil
}

func getKeys() *keyMap {
	return &keyMap{
		Numbers: key.NewBinding(
			key.WithKeys("1", "2", "3", "4"),
		),
		Next: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "next"),
		),
		Previous: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "previous"),
		),
		Reset: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reset"),
		),
		Home: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "home"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
	}
}
