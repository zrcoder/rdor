package internal

import "github.com/charmbracelet/bubbles/key"

var (
	QuitKey = key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	)
	HelpKey = key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toogle help"),
	)
	ResetKey = key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reset"),
	)
	HomeKey = key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "home"),
	)
	NextKey = key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "next"),
	)
	PreviousKey = key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "previous"),
	)
)

type KeyMap struct {
	Next     key.Binding
	Previous key.Binding
	Reset    key.Binding
	Home     key.Binding
	Quit     key.Binding
	Help     key.Binding
}

var Keys = KeyMap{
	Next:     NextKey,
	Previous: PreviousKey,
	Reset:    ResetKey,
	Home:     HomeKey,
	Quit:     QuitKey,
	Help:     HelpKey,
}

func (g *Game) ShortHelp() []key.Binding {
	res := []key.Binding{}
	for _, k := range g.extroKeys {
		res = append(res, k)
	}
	res = append(res, g.Keys.Quit)
	return res
}

func (g *Game) FullHelp() [][]key.Binding {
	return [][]key.Binding{g.extroKeys,
		{
			g.Keys.Reset,
			g.Keys.Previous,
			g.Keys.Next,
			g.Keys.Help,
			g.Keys.Home,
			g.Keys.Quit},
	}
}
