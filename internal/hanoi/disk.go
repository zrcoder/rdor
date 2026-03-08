package hanoi

import (
	"strings"

	"charm.land/lipgloss/v2"
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
