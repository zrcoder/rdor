package last

type level struct {
	totalCells int // 30-50, include the two players
	eatingMax  int // 2-5 every turn, and the `min` limit is 1
	hard       bool
}

func getDefaultLevers() []*level {
	return []*level{
		{totalCells: 30, eatingMax: 2},
		{totalCells: 41, eatingMax: 3},
		{totalCells: 30, eatingMax: 2, hard: true},
		{totalCells: 40, eatingMax: 3, hard: true},
		{totalCells: 50, eatingMax: 4, hard: true},
		{totalCells: 50, eatingMax: 5, hard: true},
	}
}
