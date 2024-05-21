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
	grid          *Grid
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
	l.grid = &Grid{}
	for i, row := range lines {
		if utf8.RuneCountInString(row) > size {
			l.SetError(fmt.Errorf("第%d行有太多字", i+1))
			return
		}
		for j, v := range []rune(row) {
			if v == emptyWord {
				continue
			}
			if v != blankWord {
				l.grid[i][j] = &Word{char: v, state: WordStateRight}
				continue
			}
			l.blanks++
			l.grid[i][j] = l.blankWord
			if l.blanks == 1 {
				l.pos.Row = i
				l.pos.Col = j
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
	for i, row := range l.grid {
		for j, word := range row {
			switch {
			case word == nil:
				l.buf.WriteRune(emptyWord)
			case word.state == WordStateRight:
				l.buf.WriteString(rightBg.Render(string(word.char)))
			case word.state == WordStateWrong:
				l.buf.WriteString(wrongBg.Render(string(word.char)))
			case word.state == WordStateBlank:
				if i == l.pos.Row && j == l.pos.Col {
					l.buf.WriteString(curBg.Render(string(emptyWord)))
				} else {
					l.buf.WriteString(blankBg.Render(string(emptyWord)))
				}
			default:
				l.buf.WriteString(string(word.char))
			}
		}
		l.buf.WriteString("\n")
	}
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
	for !l.outOfRange(&pos) && cnt > 0 {
		if l.getWord(&pos) != nil && !l.getWord(&pos).Fixed() {
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
	if word == nil {
		return
	}
	l.candidates[i] = nil
	cur := l.curWord()
	l.setCurWord(word)
	if cur.state != WordStateBlank {
		l.candidates.Set(cur)
	}
	if !l.check() {
		return
	}
	if l.success() {
		l.SetSuccess("成功")
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
	for ; left >= 0 && l.grid.get(l.pos.Row, left) != nil && l.grid.get(l.pos.Row, left).state != WordStateBlank; left-- {
	}
	for ; right < size && l.grid.get(l.pos.Row, right) != nil && l.grid.get(l.pos.Row, right).state != WordStateBlank; right++ {
	}
	if left+1+idiomLen != right {
		return true
	}
	ok := true
	for i := left + 1; i < right; i++ {
		word := l.grid.get(l.pos.Row, i)
		if word.Fixed() {
			continue
		}
		if word.destPos != l.pos.Row*size+i {
			ok = false
			word.state = WordStateWrong
		} else {
			l.blanks--
			word.state = WordStateRight
		}
	}
	return ok
}

func (l *Level) checkVertical() bool {
	up, down := l.pos.Row, l.pos.Row
	for ; up >= 0 && l.grid.get(up, l.pos.Col) != nil && l.grid.get(up, l.pos.Col).state != WordStateBlank; up-- {
	}
	for ; down < size && l.grid.get(down, l.pos.Col) != nil && l.grid.get(down, l.pos.Col).state != WordStateBlank; down++ {
	}
	if up+1+idiomLen != down {
		return true
	}
	ok := true
	for i := up + 1; i < down; i++ {
		word := l.grid.get(i, l.pos.Col)
		if word.Fixed() {
			continue
		}
		if word.destPos != i*size+l.pos.Col {
			word.state = WordStateWrong
		} else {
			l.blanks--
			word.state = WordStateRight
		}
	}
	return ok
}

func (l *Level) moveToNearestPos(dir ...grid.Direction) {
	var q []grid.Position
	seen := make(map[grid.Position]bool, size*size)
	invalid := func(p grid.Position) bool {
		return l.outOfRange(&p) || seen[p] || l.getWord(&p) != l.blankWord && l.getWord(&p).isEmpty()
	}
	visArea := func(startRow, startCol, endRow, endCol int) {
		for i := startRow; i <= endRow; i++ {
			for j := startCol; j <= endCol; j++ {
				pos := grid.Position{Row: i, Col: j}
				if !invalid(pos) {
					seen[pos] = true
					q = append(q, pos)
				}
			}
		}
	}
	if dir == nil {
		visArea(l.pos.Row, l.pos.Col, l.pos.Row, l.pos.Col)
	} else {
		switch dir[0].Opposite() {
		case grid.Left:
			visArea(0, 0, size-1, l.pos.Col)
		case grid.Right:
			visArea(0, l.pos.Col, size-1, size-1)
		case grid.Down:
			visArea(l.pos.Row, 0, size-1, size-1)
		case grid.Up:
			visArea(0, 0, l.pos.Row, size-1)
		}
	}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		for _, d := range l.directions {
			next := cur.TransForm(d)
			if invalid(next) {
				continue
			}
			word := l.getWord(&next)
			if word != nil && !word.Fixed() {
				l.pos = next
				return
			}
			seen[next] = true
			q = append(q, next)
		}
	}
	l.SetError(errors.New("无法移动"))
}

func (l *Level) curWord() *Word {
	return l.getWord(&l.pos)
}

func (l *Level) setCurWord(w *Word) {
	l.grid[l.pos.Row][l.pos.Col] = w
}

func (l *Level) getWord(p *grid.Position) *Word {
	return l.grid.get(p.Row, p.Col)
}

func (l *Level) outOfRange(p *grid.Position) bool {
	return p.Row < 0 || p.Row >= size || p.Col < 0 || p.Col >= size
}
