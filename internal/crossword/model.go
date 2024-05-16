package crossword

type Level struct {
	Grid       string
	Candidates string
	Answers    map[string]Idiom
}

type Idiom struct {
	Meaning string
	Example string
}

type Grid [size][size]*Word

type Word struct {
	char         rune
	isFixed      bool
	isBlank      bool
	isWrong      bool
	candinatePos int
}
