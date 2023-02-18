package sokoban

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Debug key.Binding

	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding

	Next     key.Binding
	Previous key.Binding
	Set      key.Binding

	Help  key.Binding
	Quit  key.Binding
	Reset key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Left, k.Down, k.Right},
		{k.Debug, k.Reset},
		{k.Next, k.Previous, k.Set},
		{k.Quit, k.Help},
	}
}

var (
	keys = keyMap{
		Up:    upBinding,
		Left:  leftBinding,
		Down:  downBinding,
		Right: rightBinding,

		Debug: debugBinding,
		Reset: resetBinging,

		Next:     nextBinding,
		Previous: previousBinding,
		Set:      setBinding,

		Quit: quitBinding,
		Help: helpBinding,
	}
	debugBinding = key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Toggle debug"),
	)
	upBinding = key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "Move up"),
	)
	leftBinding = key.NewBinding(
		key.WithKeys("left"),
		key.WithHelp("←", "Move left"),
	)
	downBinding = key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "Move down"),
	)
	rightBinding = key.NewBinding(
		key.WithKeys("right"),
		key.WithHelp("→", "Move right"),
	)

	nextBinding = key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "Next level"),
	)
	previousBinding = key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Previous level"),
	)
	setBinding = key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Set level"),
	)

	helpBinding = key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "Toggle help"),
	)
	quitBinding = key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "Quit"),
	)
	resetBinging = key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "Reset current level"),
	)
)
