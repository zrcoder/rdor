package hanoi

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type disk struct {
	id    int
	width int
	view  string
}

func newDisk(id int, sty lipgloss.Style) *disk {
	view := sty.Render(strings.Repeat(diskCh, id*diskWidthUnit))
	width, _ := lipgloss.Size(view)
	return &disk{
		id:    id,
		view:  view,
		width: width,
	}
}
