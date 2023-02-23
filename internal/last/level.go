package last

import (
	"fmt"

	"github.com/zrcoder/rdor/pkg/style"
)

type level struct {
	id         int
	totalCells int // 30-50, include the two players
	eatingMax  int // 2-4 every turn, and the `min` limit is 1
	hard       bool
}

func (l level) shortView() string {
	s := fmt.Sprintf("Level: %d  Total: %d  limit: %d", l.id, l.totalCells, l.eatingMax)
	s = style.Help.Render(s)
	if l.hard {
		s += style.Warn.Render("  hard")
	}
	return s
}

func (l level) View() string {
	s := fmt.Sprintf("There are totally %d cells(include you and your rival) in the world now.\nEach turn, the current player should eat at least 1 cell, and no more than %d.\n> The player can only eat his rival in the `last` turn~", l.totalCells, l.eatingMax)
	return style.Help.Render(s)
}

func getDefaultLevers() []*level {
	return []*level{
		// the first hand is advantageous
		{totalCells: 30, eatingMax: 2},
		{totalCells: 30, eatingMax: 2, hard: true},
		// the second hand is advantageous
		{totalCells: 34, eatingMax: 2},
		{totalCells: 34, eatingMax: 2, hard: true},
		// the first hand is advantageous
		{totalCells: 40, eatingMax: 3},
		{totalCells: 40, eatingMax: 3, hard: true},
		// the second hand is advantageous
		{totalCells: 56, eatingMax: 4, hard: true},
	}
}
