package last

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Numbers  key.Binding
	Next     key.Binding
	Previous key.Binding
	Reset    key.Binding
	Quit     key.Binding
	Help     key.Binding
}

func (k *keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Numbers, k.Next, k.Previous, k.Reset, k.Quit, k.Help}
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
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
	}
}
