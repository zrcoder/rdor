package grid

import (
	"strings"
)

type Grid[T comparable] struct {
	rows int
	cols int
	data [][]T
}

func New[T comparable](rows, cols int) *Grid[T] {
	data := make([][]T, rows)
	for i := range data {
		data[i] = make([]T, cols)
	}
	return &Grid[T]{data: data, rows: rows, cols: cols}
}

func NewWithString(s string) *Grid[rune] {
	lines := strings.Split(s, "\n")
	data := make([][]rune, len(lines))
	cols := 0
	for i, line := range lines {
		data[i] = []rune(line)
		cols = max(cols, len(data[i]))
	}
	return &Grid[rune]{data: data, rows: len(lines), cols: cols}
}

func (g *Grid[T]) SetData(data [][]T) {
	g.rows = len(data)
	if g.rows == 0 {
		return
	}
	g.cols = len(data[0])
	g.data = data
}

func (g *Grid[T]) Copied() *Grid[T] {
	data := make([][]T, len(g.data))
	for i, row := range g.data {
		data[i] = make([]T, len(row))
		copy(data[i], row)
	}
	return &Grid[T]{data: data}
}

func (g *Grid[T]) Copy(gg *Grid[T]) {
	for i, row := range gg.data {
		copy(g.data[i], row)
	}
}

func (g *Grid[T]) CopyTo(gg *Grid[T]) {
	for i, row := range g.data {
		copy(gg.data[i], row)
	}
}

func (g *Grid[T]) Equal(f *Grid[T]) bool {
	if len(g.data) != len(f.data) {
		return false
	}
	for r := range g.data {
		if len(g.data[r]) != len(f.data[r]) {
			return false
		}
		for c := range g.data[r] {
			if g.data[r][c] != f.data[r][c] {
				return false
			}
		}
	}
	return true
}

func (g *Grid[T]) OutBound(pos Position) bool {
	return pos.Row < 0 || pos.Row >= len(g.data) ||
		pos.Col < 0 || pos.Col >= len(g.data[pos.Row])
}

type (
	RangeAction[T comparable]     func(pos Position, char T, isLineEnd bool) (end bool)
	RangeRowsAction[T comparable] func(r int, row []T, isLast bool) (end bool)
)

func (g *Grid[T]) Range(action RangeAction[T]) {
	for i, row := range g.data {
		for j, v := range row {
			if action(Position{Row: i, Col: j}, v, j == len(row)-1) {
				return
			}
		}
	}
}

func (g *Grid[T]) RangeRows(action RangeRowsAction[T]) {
	for r, row := range g.data {
		if action(r, row, r == len(g.data)-1) {
			return
		}
	}
}

func (g *Grid[T]) Get(p Position) T {
	return g.data[p.Row][p.Col]
}

func (g *Grid[T]) Getrc(row, col int) T {
	return g.data[row][col]
}

func (g *Grid[T]) Set(p Position, val T) {
	g.data[p.Row][p.Col] = val
}

func (g *Grid[T]) Nearest(from Position, dirs []Direction, ok func(Position) bool, dir ...Direction) *Position {
	var q []Position
	seen := make(map[Position]bool, g.rows*g.cols)
	invalid := func(p Position) bool {
		return g.OutBound(p) || seen[p]
	}
	visPos := func(p Position) {
		seen[p] = true
		q = append(q, p)
	}
	visArea := func(startRow, startCol, endRow, endCol int) {
		for i := startRow; i <= endRow; i++ {
			for j := startCol; j <= endCol; j++ {
				visPos(Position{Row: i, Col: j})
			}
		}
	}
	if dir == nil {
		visPos(from)
	} else {
		switch dir[0].Opposite() {
		case Left:
			visArea(0, 0, g.rows-1, from.Col)
		case Right:
			visArea(0, from.Col, g.rows-1, g.cols-1)
		case Down:
			visArea(from.Row, 0, g.rows-1, g.cols-1)
		case Up:
			visArea(0, 0, from.Row, g.cols-1)
		default:
			// TODO, for other directions
		}
	}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		for _, d := range dirs {
			next := cur.TransForm(d)
			if invalid(next) {
				continue
			}
			if ok(next) {
				return &next
			}
			seen[next] = true
			q = append(q, next)
		}
	}
	return nil
}
