package crossword

import (
	"errors"
	"fmt"
	"unicode/utf8"

	"github.com/zrcoder/rdor/pkg/grid"
)

type Level struct {
	*crossword
	Grid          []string `toml:"grid"`
	Candidates    string   `toml:"candidates"`
	AnswerPos     []int    `toml:"answerPos"`
	candidates    Candidates
	candidatesPos map[byte]int
	blanks        int
	grid          *grid.Grid[*Word]
}

func (l *Level) adapt() {
	if l.adaptGrid(); l.Err != nil {
		return
	}
	l.adaptCandidates()
}

func (l *Level) adaptGrid() {
	lines := l.Grid
	if len(lines) > size {
		l.SetError(errors.New("配置中有太多行"))
		return
	}
	l.grid = grid.New[*Word](size, size)
	for i, row := range lines {
		if utf8.RuneCountInString(row) > size {
			l.SetError(fmt.Errorf("第%d行有太多字", i+1))
			return
		}
		for j, v := range []rune(row) {
			if v == emptyWord {
				continue
			}
			pos := grid.Position{Row: i, Col: j}
			if v != blankWord {
				l.grid.Set(pos, &Word{char: v, state: WordStateRight})
				continue
			}
			l.blanks++
			l.grid.Set(pos, l.blankWord)
			if l.blanks == 1 {
				l.pos = pos
			}
		}
	}
	if l.blanks == 0 {
		l.SetError(errors.New("没有空格要填"))
	}
}

func (l *Level) adaptCandidates() {
	cfg := l.Candidates
	n := utf8.RuneCountInString(cfg)
	if n > candidatesLimit {
		l.SetError(errors.New("候选字过多"))
		return
	}
	l.candidates = make([]*Word, n)
	l.candidatesPos = make(map[byte]int, n)
	for i, v := range []rune(cfg) {
		l.candidates[i] = &Word{char: v, candidatePos: i, destPos: l.AnswerPos[i]}
		l.candidatesPos[candidatesKeys[i]] = i
	}
}

func (l *Level) boardView() string {
	l.buf.Reset()
	l.grid.Range(func(pos grid.Position, word *Word, isLineEnd bool) (end bool) {
		if pos == l.pos && !word.Fixed() {
			l.buf.WriteString(curBg.Render(string(word.char)))
		} else {
			l.buf.WriteString(word.View())
		}
		if isLineEnd {
			l.buf.WriteString("\n")
		}
		return false
	})
	return boardStyle.Render(l.buf.String())
}

func (l *Level) candidatesView() string {
	l.buf.Reset()
	for i, w := range l.candidates {
		l.buf.WriteByte(candidatesKeys[i])
		l.buf.WriteRune(':')
		if w != nil {
			l.buf.WriteRune(w.char)
		} else {
			l.buf.WriteRune(emptyWord)
		}
		if (i+1)%candidatesPerLine == 0 {
			l.buf.WriteRune('\n')
		} else {
			l.buf.WriteRune(emptyWord)
		}
	}
	return boardStyle.Render(l.buf.String())
}

func (l *Level) move(d grid.Direction) {
	pos := l.pos.TransForm(d)
	cnt := size*size - 1
	for !l.grid.OutBound(pos) && cnt > 0 {
		word := l.grid.Get(pos)
		if word != nil && !word.Fixed() {
			l.pos = pos
			return
		}
		pos = pos.TransForm(d)
		cnt--
	}
	l.moveToNearestPos(d)
}

func (l *Level) pick(i int) {
	if i == -1 {
		cur := l.curWord()
		if cur == l.blankWord || cur.Fixed() {
			return
		}
		l.setCurWord(l.blankWord)
		l.candidates.Set(cur)
		return
	}
	word := l.candidates[i]
	cur := l.curWord()
	if word == nil || cur != nil && cur.Fixed() {
		return
	}
	l.candidates[i] = nil
	l.setCurWord(word)
	if cur.state != WordStateBlank {
		l.candidates.Set(cur)
	}
	if !l.check() || l.success() {
		return
	}
	l.moveToNearestPos()
}

func (l *Level) success() bool {
	return l.blanks == 0
}

func (l *Level) check() bool {
	return l.checkHorizental() && l.checkVertical()
}

func (l *Level) checkHorizental() bool {
	left, right := l.pos.Col, l.pos.Col
	for ; left >= 0 && l.grid.Getrc(l.pos.Row, left) != nil && l.grid.Getrc(l.pos.Row, left).state != WordStateBlank; left-- {
	}
	for ; right < size && l.grid.Getrc(l.pos.Row, right) != nil && l.grid.Getrc(l.pos.Row, right).state != WordStateBlank; right++ {
	}
	if left+1+idiomLen != right {
		return true
	}
	return l.checkIdiom(l.pos.Row, left+1, l.pos.Row, right-1)
}

func (l *Level) checkVertical() bool {
	up, down := l.pos.Row, l.pos.Row
	for ; up >= 0 && l.grid.Getrc(up, l.pos.Col) != nil && l.grid.Getrc(up, l.pos.Col).state != WordStateBlank; up-- {
	}
	for ; down < size && l.grid.Getrc(down, l.pos.Col) != nil && l.grid.Getrc(down, l.pos.Col).state != WordStateBlank; down++ {
	}
	if up+1+idiomLen != down {
		return true
	}
	return l.checkIdiom(up+1, l.pos.Col, down-1, l.pos.Col)
}

func (l *Level) checkIdiom(startR, startC, endR, endC int) bool {
	ok := true
	for i := startR; i <= endR; i++ {
		for j := startC; j <= endC; j++ {
			word := l.grid.Getrc(i, j)
			if !word.Fixed() && word.destPos != i*size+j {
				ok = false
				break
			}
		}
	}
	state := WordStateRight
	delta := -1
	if !ok {
		state = WordStateWrong
		delta = 0
	}
	for i := startR; i <= endR; i++ {
		for j := startC; j <= endC; j++ {
			word := l.grid.Getrc(i, j)
			if !word.Fixed() {
				word.state = state
				l.blanks += delta
			}
		}
	}
	return ok
}

func (l *Level) moveToNearestPos(dir ...grid.Direction) {
	pos := l.grid.Nearest(l.pos, l.directions, func(p grid.Position) bool {
		word := l.grid.Get(p)
		return word != nil && !word.Fixed()
	}, dir...)
	if pos == nil {
		l.SetError(errors.New("无法移动"))
	} else {
		l.pos = *pos
	}
}

func (l *Level) curWord() *Word {
	return l.grid.Get(l.pos)
}

func (l *Level) setCurWord(w *Word) {
	l.grid.Set(l.pos, w)
}
