package keyblock

import (
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/rdor/pkg/style/color"
)

var (
	keyStyle     = lg.NewStyle().Foreground(color.Faint)
	normalStyle  = lg.NewStyle().Padding(0, 1).Border(lg.RoundedBorder())
	pressedStyle = normalStyle.Copy().Faint(true).Foreground(color.Faint)
)

type Action func(key *Key)

type Key struct {
	Key     string
	Display string
	Once    bool
	Action  Action
	pressed bool
}

func NewKey(once bool, key string) *Key {
	return &Key{Key: key, Once: once}
}

func (k *Key) SetDisply(display string) {
	k.Display = display
}

func (k *Key) SetAction(action Action) {
	k.Action = action
}

func (k *Key) Init() tea.Cmd { return nil }

func (k *Key) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		val := string(msg.Runes)
		if val == k.Key {
			if !k.pressed && k.Action != nil {
				k.Action(k)
			}
			if k.Once {
				k.pressed = true
			}
		}
	}
	return k, nil
}

func (k *Key) View() string {
	display := k.Display
	if display == "" {
		display = k.Key
	}
	if k.pressed {
		display = pressedStyle.Render(display)
	} else {
		display = normalStyle.Render(display)
	}
	return lg.JoinVertical(lg.Center,
		display,
		keyStyle.Render(k.Key),
	)
}

type KeysLine []*Key

func NewKeysLine(once bool, keys ...string) KeysLine {
	res := make([]*Key, len(keys))
	for i, k := range keys {
		res[i] = NewKey(once, k)
	}
	return res
}

func (kl KeysLine) SetDisplays(displays ...string) {
	for i, display := range displays {
		kl[i].SetDisply(display)
	}
}

func (kl KeysLine) SetDisplay(i int, display string) {
	kl[i].SetDisply(display)
}

func (kl KeysLine) SetAction(action Action) {
	for i := range kl {
		kl[i].SetAction(action)
	}
}

func (kl KeysLine) SetActionAt(i int, action Action) {
	kl[i].SetAction(action)
}

func (kl KeysLine) Init() tea.Cmd {
	return nil
}

func (kl KeysLine) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	for _, k := range kl {
		k.Update(msg)
	}
	return kl, nil
}

func (kl KeysLine) View() string {
	views := make([]string, len(kl))
	for i, k := range kl {
		views[i] = k.View()
	}
	return lg.JoinHorizontal(lg.Center, views...)
}
