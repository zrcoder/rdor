package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if _, err := tea.NewProgram(New()).Run(); err != nil {

	}
}
