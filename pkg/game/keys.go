package game

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Next     key.Binding
	Previous key.Binding
	Reset    key.Binding
	Home     key.Binding
	Quit     key.Binding
	Help     key.Binding
}

func (g *Base) ShortHelp() []key.Binding {
	res := []key.Binding{}
	for _, k := range g.Keys {
		res = append(res, k)
	}
	res = append(res, g.CommonKeys.Quit)
	return res
}

func (g *Base) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		g.Keys,
		{
			g.CommonKeys.Reset,
			g.CommonKeys.Next,
			g.CommonKeys.Previous,
		},
		{
			g.CommonKeys.Help,
			g.CommonKeys.Home,
			g.CommonKeys.Quit,
		},
	}
}

func getKeys() KeyMap {
	return KeyMap{
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
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toogle help"),
		),
	}
}
