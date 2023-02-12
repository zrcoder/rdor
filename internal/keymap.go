package internal

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Nums  key.Binding
	Help  key.Binding
	Quit  key.Binding
	Reset key.Binding
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Nums},                  // first column
		{k.Help, k.Quit, k.Reset}, // second column
	}
}

var (
	keysSetting = keyMap{
		Nums: disksBinding,
		Help: helpBinding,
		Quit: quitBinding,
	}
	keysSetted = keyMap{
		Nums:  pilesBinding,
		Help:  helpBinding,
		Reset: resetBinging,
		Quit:  quitBinding,
	}
	keysHealping = keyMap{
		Help: helpBinding,
		Quit: quitBinding,
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
		key.WithHelp("1-3(j,k,l):", "Pick the special pile"),
	)
	resetBinging = key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r:", "Reset"),
	)
)

func contains(keys []string, key string) bool {
	for _, v := range keys {
		if v == key {
			return true
		}
	}
	return false
}
