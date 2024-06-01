package game

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	help, back, quit *key.Binding

	reset, next, previous, setLevel *key.Binding

	groups []KeyGroup
}

type KeyGroup []*key.Binding

func (g KeyGroup) ShortHelp() []key.Binding {
	ks := make([]key.Binding, len(g))
	for i, k := range g {
		ks[i] = *k
	}
	return ks
}

func (g KeyGroup) FullHelp() [][]key.Binding {
	return [][]key.Binding{g.ShortHelp()}
}

func getCommonKeys() *KeyMap {
	reset := key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reset"),
		key.WithDisabled(),
	)
	next := key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "next"),
		key.WithDisabled(),
	)
	previous := key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "previous"),
		key.WithDisabled(),
	)
	setLevel := key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "set level"),
		key.WithDisabled(),
	)
	help := key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
		key.WithDisabled(),
	)
	back := key.NewBinding(
		key.WithKeys("ctrl+h"),
		key.WithHelp("ctrl+h", "back home"),
	)
	quit := key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	)
	res := &KeyMap{
		reset:    &reset,
		next:     &next,
		previous: &previous,
		setLevel: &setLevel,
		help:     &help,
		back:     &back,
		quit:     &quit,
	}
	res.groups = []KeyGroup{
		{res.help},
		{res.reset, res.next, res.previous, res.setLevel},
		{res.back, res.quit},
	}
	return res
}
