package keyblock

import (
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/rdor/pkg/style/color"
)

type Key struct {
	Key        string
	Display    string
	keyBinding key.Binding
	Number     int
	isNumber   bool
	Pressed    bool
}

type PressMsg = *Key

func NewKey(key string) *Key {
	return &Key{Key: key}
}

func (k *Key) SetNumber(val int) {
	k.isNumber = true
	k.Number = val
}

func (k *Key) SetDisply(display string) {
	k.Display = display
}

func (k *Key) RemoveNumber() {
	k.isNumber = false
}

func (k *Key) IsNumber() bool {
	return k.isNumber
}

func (k *Key) Init() tea.Cmd {
	k.keyBinding = key.NewBinding(key.WithKeys(k.Key))
	return nil
}

func (k *Key) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, k.keyBinding) {
			k.Pressed = !k.Pressed
		}
	}
	return k, func() tea.Msg { return PressMsg(k) }
}

func (k *Key) View() string {
	style := lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.RoundedBorder())
	if k.Pressed {
		style = style.Faint(true).Foreground(color.Faint)
	}
	display := k.Display
	if k.isNumber {
		display = strconv.Itoa(k.Number)
	}
	if display == "" {
		display = " "
	}
	return lipgloss.JoinVertical(lipgloss.Center,
		style.Render(display),
		lipgloss.NewStyle().Foreground(color.Faint).Render(k.Key),
	)
}

type KeysLine []*Key

func NewKeysLine(keys ...string) KeysLine {
	res := make([]*Key, len(keys))
	for i, k := range keys {
		res[i] = NewKey(k)
	}
	return res
}

func (kl KeysLine) SetNumbers(nums ...int) {
	for i, num := range nums {
		kl[i].SetNumber(num)
	}
}

func (kl KeysLine) SetDisplays(displays ...string) {
	for i, display := range displays {
		kl[i].SetDisply(display)
	}
}

func (kl KeysLine) Init() tea.Cmd {
	return nil
}

func (kl KeysLine) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	for _, k := range kl {
		_, c := k.Update(msg)
		cmds = append(cmds, c)
	}
	return kl, tea.Sequence(cmds...)
}

func (kl KeysLine) View() string {
	views := make([]string, len(kl))
	for i, k := range kl {
		views[i] = k.View()
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, views...)
}
