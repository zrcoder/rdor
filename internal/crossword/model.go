package crossword

import (
	"embed"

	lg "github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/rdor/pkg/style/color"
)

//go:embed levels/*.toml
var lvsFS embed.FS

var (
	rightBg    = lg.NewStyle().Background(color.Faint)
	wrongBg    = lg.NewStyle().Background(color.Red)
	blankBg    = lg.NewStyle().Background(color.Orange)
	curBg      = lg.NewStyle().Background(color.Violet)
	successBg  = lg.NewStyle().Background(color.Green)
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

func (w *Word) isEmpty() bool {
	return w == nil || w.char == emptyWord
}

type (
	Grid       [size][size]*Word
	Candidates []*Word
)

func (g *Grid) get(i, j int) *Word {
	return g[i][j]
}

func (cs Candidates) Set(word *Word) {
	word.state = WordStateInit
	cs[word.candidatePos] = word
}
