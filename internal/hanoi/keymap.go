package hanoi

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Disks    key.Binding
	Piles    key.Binding
	Next     key.Binding
	Previous key.Binding
	Reset    key.Binding
	Set      key.Binding
	Home     key.Binding
}

func (k *keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Disks, k.Piles, k.Set, k.Home}
}

func (k *keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Disks, k.Piles, k.Reset},
		{k.Next, k.Previous},
		{k.Set, k.Home},
	}
}

func getKeys() *keyMap {
	return &keyMap{
		Disks: key.NewBinding(
			key.WithKeys("1", "2", "3", "4", "5", "6", "7"),
			key.WithHelp("1-7", "set disks"),
		),
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
