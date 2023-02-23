package model

import tea "github.com/charmbracelet/bubbletea"

type Game interface {
	tea.Model
	SetParent(tea.Model)
}
