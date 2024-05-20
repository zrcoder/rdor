package crossword

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/rdor/pkg/game"
	"github.com/zrcoder/rdor/pkg/grid"
	"github.com/zrcoder/rdor/pkg/keys"
	"github.com/zrcoder/rdor/pkg/style"
)

func New() game.Game {
	return &crossword{Base: game.New(name)}
}

type crossword struct {
	*game.Base

	grid          *Grid
	candidates    []*Word
	buf           *strings.Builder
	candidateBuf  *strings.Builder
	state         string
	downKey       *key.Binding
	leftKey       *key.Binding
	upKey         *key.Binding
	rightKey      *key.Binding
	directions    []grid.Direction
	pos           grid.Position
	blankWord     *Word
	candidatesPos map[byte]int
	level         *Level
	levels        int
	blanks        int
}

func (c *crossword) Init() tea.Cmd {
	c.RegisterView(c.view)
	c.loadSummary()
	c.buf = &strings.Builder{}
	c.candidateBuf = &strings.Builder{}
	c.blankWord = &Word{state: WordStateBlank}
	c.directions = []grid.Direction{grid.Down, grid.Right, grid.Up, grid.Left}
	c.upKey = &keys.Up
	c.leftKey = &keys.Left
	c.downKey = &keys.Down
	c.rightKey = &keys.Right
	c.ClearGroups()
	c.AddKeyGroup(game.KeyGroup{c.upKey, c.leftKey, c.downKey, c.rightKey})
	c.set(0)
	return c.Base.Init()
}

func (c *crossword) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	b, cmd := c.Base.Update(msg)
	if b != c.Base {
		return b, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, *c.upKey):
			c.move(grid.Up)
		case key.Matches(msg, *c.downKey):
			c.move(grid.Down)
		case key.Matches(msg, *c.leftKey):
			c.move(grid.Left)
		case key.Matches(msg, *c.rightKey):
			c.move(grid.Right)
		default:
			if msg.Type == tea.KeyEnter {
				c.pick(-1)
				break
			}
			if len(msg.Runes) > 0 {
				letter := byte(unicode.ToUpper(msg.Runes[0]))
				if i, ok := c.candidatesPos[letter]; ok {
					c.pick(i)
				}
			}
		}
	}
	return c, cmd
}

func (c *crossword) view() string {
	if c.Err != nil {
		return ""
	}
	return lg.JoinVertical(lg.Left,
		c.boardView(),
		c.candidatesView(),
		"",
		c.state)
}

func (c *crossword) loadSummary() {
	path := filepath.Join("levels", "index.toml")
	data, err := lvsFS.ReadFile(path)
	if err != nil {
		c.SetError(err)
		return
	}
	type LevelSummary struct {
		Levels int
	}
	ls := &LevelSummary{}
	err = toml.Unmarshal(data, ls)
	if err != nil {
		c.SetError(err)
		return
	}
	c.levels = ls.Levels
	c.RegisterLevels(ls.Levels, c.set)
}

func (c *crossword) set(i int) {
	data, err := lvsFS.ReadFile(fmt.Sprintf("levels/%02d.toml", i))
	if err != nil {
		c.SetError(err)
		return
	}
	c.level = &Level{}
	c.blanks = 0
	err = toml.Unmarshal(data, c.level)
	if err != nil {
		c.SetError(err)
		return
	}
	if c.newGrid(); c.Err != nil {
		return
	}
	if c.newCandidates(); c.Err != nil {
		return
	}
	c.state = style.Help.Render(fmt.Sprintf("%d/%d", i+1, c.levels))
}

func (c *crossword) newGrid() {
	cfg := c.level.Grid
	if len(cfg) > size {
		c.SetError(errors.New("配置中有太多行"))
		return
	}
	c.grid = &Grid{}
	for i, row := range cfg {
		for j, v := range []rune(row) {
			if v == emptyWord {
				continue
			}
			if v != blankWord {
				c.grid[i][j] = &Word{char: v, state: WordStateRight}
				continue
			}
			c.blanks++
			c.grid[i][j] = c.blankWord
			if c.blanks == 1 {
				c.pos.Row = i
				c.pos.Col = j
			}
		}
	}
	if c.blanks == 0 {
		c.SetError(errors.New("没有空格要填"))
	}
}

func (c *crossword) newCandidates() {
	cfg := c.level.Candidates
	size := utf8.RuneCountInString(cfg)
	if size > candidatesLimit {
		c.SetError(errors.New("too many candidaters"))
		return
	}
	c.candidates = make([]*Word, size)
	c.candidatesPos = make(map[byte]int, size)
	for i, v := range []rune(cfg) {
		c.candidates[i] = &Word{char: v, candidatePos: i, destPos: c.level.AnswerPos[i]}
		c.candidatesPos[candidatesKeys[i]] = i
	}
}

func (c *crossword) boardView() string {
	c.buf.Reset()
	for i, row := range c.grid {
		for j, word := range row {
			if word == nil {
				c.buf.WriteRune(emptyWord)
				continue
			}
			s := string(word.char)
			switch word.state {
			case WordStateBlank:
				s = string(blankWord)
			case WordStateRight:
				s = rightBg.Render(s)
			case WordStateWrong:
				if i != c.pos.Row || j != c.pos.Col {
					s = wrongBg.Render(s)
				}
			}
			if i == c.pos.Row && j == c.pos.Col {
				s = curBg.Render(s)
			}
			c.buf.WriteString(s)
		}
		c.buf.WriteString("\n")
	}
	return boardStyle.Render(c.buf.String())
}

func (c *crossword) candidatesView() string {
	c.buf.Reset()
	for i, w := range c.candidates {
		c.buf.WriteByte(candidatesKeys[i])
		c.buf.WriteRune(':')
		if w != nil {
			c.buf.WriteRune(w.char)
		} else {
			c.buf.WriteRune(emptyWord)
		}
		if (i+1)%candidatesPerLine == 0 {
			c.buf.WriteRune('\n')
		} else {
			c.buf.WriteRune(emptyWord)
		}
	}
	return boardStyle.Render(c.buf.String())
}

func (c *crossword) move(d grid.Direction) {
	cnt := 0
	i, j := (c.pos.Row+size+d.Dy)%size, (c.pos.Col+size+d.Dx)%size
	for cnt < size*size {
		if i == c.pos.Row && j == c.pos.Col {
			c.moveToNextBlankPos()
			return
		}
		if c.grid[i][j] != nil && !c.grid[i][j].Fixed() {
			c.pos.Row = i
			c.pos.Col = j
			return
		}
		i, j = (i+d.Dy+size)%size, (j+d.Dx+size)%size
		cnt++
	}
	c.SetError(errors.New("无法移动"))
}

func (c *crossword) pick(i int) {
	if i == -1 {
		cur := c.curWord()
		if cur == c.blankWord || cur.Fixed() {
			return
		}
		c.setCurWord(c.blankWord)
		c.candidates[cur.candidatePos] = cur
		return
	}

	word := c.candidates[i]
	c.candidates[i] = nil
	if word == nil {
		return
	}
	cur := c.curWord()
	c.setCurWord(word)
	if cur.state != WordStateBlank {
		c.candidates[cur.candidatePos] = cur
	}
	if !c.check() {
		return
	}
	if c.success() {
		c.SetSuccess("成功")
		return
	}
	c.moveToNextBlankPos()
}

func (c *crossword) success() bool {
	return c.blanks == 0
}

func (c *crossword) check() bool {
	return c.checkHorizental() && c.checkVertical()
}

func (c *crossword) checkHorizental() bool {
	left, right := c.pos.Col, c.pos.Col
	for ; left > 0 && c.grid[c.pos.Row][left] != nil && c.grid[c.pos.Row][left].state != WordStateBlank; left-- {
	}
	for ; right < size && c.grid[c.pos.Row][right] != nil && c.grid[c.pos.Row][right].state != WordStateBlank; right++ {
	}
	if left+1+idiomLen == right {
		ok := true
		for i := left + 1; i < right; i++ {
			word := c.grid[c.pos.Row][i]
			if word.Fixed() {
				continue
			}
			pos := c.pos.Row*size + i
			if word.destPos != pos {
				ok = false
				word.state = WordStateWrong
			} else {
				c.blanks--
				word.state = WordStateRight
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

func (c *crossword) checkVertical() bool {
	up, down := c.pos.Row, c.pos.Row
	for ; up > 0 && c.grid[up][c.pos.Col] != nil && c.grid[up][c.pos.Col].state != WordStateBlank; up-- {
	}
	for ; down < size && c.grid[down][c.pos.Col] != nil && c.grid[down][c.pos.Col].state != WordStateBlank; down++ {
	}
	if up+1+idiomLen == down {
		ok := true
		for i := up + 1; i < down; i++ {
			word := c.grid[i][c.pos.Col]
			if word.Fixed() {
				continue
			}
			pos := i*size + c.pos.Col
			if word.destPos != pos {
				word.state = WordStateWrong
			} else {
				c.blanks--
				word.state = WordStateRight
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

func (c *crossword) moveToNextBlankPos() {
	seen := make(map[grid.Position]bool, size*size)
	seen[c.pos] = true
	q := []grid.Position{c.pos}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		for _, d := range c.directions {
			next := grid.TransForm(cur, d)
			if c.outOfRange(&next) || seen[next] {
				continue
			}
			seen[next] = true
			word := c.getWord(&next)
			if word != nil && !word.Fixed() {
				c.pos = next
				return
			}
			q = append(q, next)
		}
	}
}

func (c *crossword) curWord() *Word {
	return c.grid[c.pos.Row][c.pos.Col]
}

func (c *crossword) setCurWord(w *Word) {
	c.grid[c.pos.Row][c.pos.Col] = w
}

func (c *crossword) getWord(p *grid.Position) *Word {
	return c.grid[p.Row][p.Col]
}

func (c *crossword) outOfRange(p *grid.Position) bool {
	return p.Row < 0 || p.Row >= size || p.Col < 0 || p.Col >= size
}
