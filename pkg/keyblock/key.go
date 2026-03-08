package keyblock

import (
	tea "charm.land/bubbletea/v2"
	lg "charm.land/lipgloss/v2"
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
	once    bool
	action  Action
	pressed bool
}

func NewKey(key string) *Key {
	return &Key{Key: key}
}

func (k *Key) SetOnce(b bool) {
	k.once = true
}

func (k *Key) SetDisply(display string) {
	k.Display = display
}

func (k *Key) SetAction(action Action) {
	k.action = action
}

func (k *Key) Init() tea.Cmd { return nil }

func (k *Key) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		val := msg.Text
		if val == k.Key {
			if !k.pressed && k.action != nil {
				k.action(k)
			}
			if k.once {
				k.pressed = true
			}
		}
	}
	return k, nil
}

func (k *Key) View() tea.View {
	display := k.Display
	if display == "" {
		display = k.Key
	}
	if k.pressed {
		display = pressedStyle.Render(display)
	} else {
		display = normalStyle.Render(display)
	}
	return tea.NewView(lg.JoinVertical(lg.Center,
		display,
		keyStyle.Render(k.Key),
	))
}

type KeysLine []*Key

func NewKeysLine(keys ...string) KeysLine {
	res := make([]*Key, len(keys))
	for i, k := range keys {
		res[i] = NewKey(k)
	}
	return res
}

func (kl KeysLine) SetOnce(b bool) {
	for i := range kl {
		kl[i].SetOnce(b)
	}
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

func (kl KeysLine) View() tea.View {
	views := make([]string, len(kl))
	for i, k := range kl {
		views[i] = k.View().Content
	}
	return tea.NewView(lg.JoinHorizontal(lg.Center, views...))
}
