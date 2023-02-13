package internal

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Disks key.Binding
	Piles key.Binding
	Help  key.Binding
	Quit  key.Binding
	Reset key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Disks, k.Piles},
		{k.Help, k.Quit, k.Reset},
	}
}

var (
	keys = keyMap{
		Disks: disksBinding,
		Piles: pilesBinding,
		Help:  helpBinding,
		Quit:  quitBinding,
		Reset: resetBinging,
	}

	helpBinding = key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h:", "Toggle help"),
	)
	quitBinding = key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q:", "Quit"),
	)
	disksBinding = key.NewBinding(
		key.WithKeys("1", "2", "3", "4", "5", "6", "7"),
		key.WithHelp("1-7:", "Set the total disks"),
	)
	pilesBinding = key.NewBinding(
		key.WithKeys("1", "2", "3", "j", "k", "l"),
		key.WithHelp("1-3 / j, k, l :", "Pick the special pile"),
	)
	resetBinging = key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r:", "Reset"),
	)
)
