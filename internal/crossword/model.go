package crossword

import (
	"embed"

	lg "github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/rdor/pkg/style/color"
)

//go:embed levels/*.toml
var lvsFS embed.FS

var (
	curBg      = lg.NewStyle().Background(color.Yellow)
	rightBg    = lg.NewStyle().Background(color.Green)
	wrongBg    = lg.NewStyle().Background(color.Red)
	boardStyle = lg.NewStyle().Width(boardWidth)
)

const (
	name              = "成语填字"
	size              = 9
	emptyWord         = '　'
	blankWord         = '〇'
	candidatesPerLine = 5
	idiomLen          = 4
	candidatesKeys    = "ACDEFGHIJKLMOTUVWXYZ"
	candidatesLimit   = len(candidatesKeys)
	boardWidth        = 50
)

type Level struct {
	Grid       []string `toml:"grid"`
	Candidates string   `toml:"candidates"`
	AnswerPos  []int    `toml:"answerPos"`
}

type Grid [size][size]*Word

type WordState int

const (
	WordStateInit WordState = iota
	WordStateBlank
	WordStateRight
	WordStateWrong
)

type Word struct {
	state        WordState
	char         rune
	candidatePos int
	destPos      int
}

func (w Word) Fixed() bool {
	return w.state == WordStateRight
}
