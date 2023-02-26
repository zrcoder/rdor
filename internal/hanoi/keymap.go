package hanoi

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Piles    key.Binding
	Next     key.Binding
	Previous key.Binding
	Reset    key.Binding
	Home     key.Binding
}

func (k *keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Piles, k.Reset, k.Next, k.Previous, k.Home}
}

func (k *keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Piles},
		{k.Next, k.Previous},
		{k.Reset, k.Home},
	}
}

func getKeys() *keyMap {
	return &keyMap{
		Piles: key.NewBinding(
			key.WithKeys("1", "2", "3", "j", "k", "l"),
			key.WithHelp("1-3/j,k,l", "pick a pile"),
		),
		Next: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "next"),
		),
		Previous: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "previous"),
		),
		Reset: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reset"),
		),
		Home: key.NewBinding(
			key.WithKeys("h"),
			key.WithHelp("h", "home"),
		),
	}
}
