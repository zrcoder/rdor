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
	boardStyle = lg.NewStyle().Width(boardWidth).Border(lg.NormalBorder()).BorderForeground(color.Faint)
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
	boardWidth        = 35
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

func (word *Word) View() string {
	if word == nil {
		return string(emptyWord)
	}
	bg := blankBg
	s := string(word.char)
	switch word.state {
	case WordStateRight:
		bg = rightBg
	case WordStateWrong:
		bg = wrongBg
	case WordStateBlank:
		s = string(emptyWord)
	}
	return bg.Render(s)
}

func (w *Word) Fixed() bool {
	return w != nil && w.state == WordStateRight
}

type Candidates []*Word

func (cs Candidates) Set(word *Word) {
	word.state = WordStateInit
	cs[word.candidatePos] = word
}
