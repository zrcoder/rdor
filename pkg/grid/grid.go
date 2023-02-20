package grid

import "strings"

type Grid struct {
	data [][]rune
}

func New(s string) *Grid {
	lines := strings.Split(s, "\n")
	data := make([][]rune, len(lines))
	for i, line := range lines {
		data[i] = []rune(line)
	}
	return &Grid{data: data}
}

func (g *Grid) Copy(gg *Grid) {
	for i, row := range gg.data {
		copy(g.data[i], row)
	}
}

func (g *Grid) OutBound(pos Position) bool {
	return pos.Row < 0 || pos.Row >= len(g.data) ||
		pos.Col < 0 || pos.Col >= len(g.data[pos.Row])
}

type RangeAction func(pos Position, char rune, isLineEnd bool) (end bool)

func (g *Grid) Range(action RangeAction) {
	for i, row := range g.data {
		for j, v := range row {
			if action(Position{Row: i, Col: j}, v, j == len(row)-1) {
				return
			}
		}
	}
}

func (g *Grid) Get(p Position) rune {
	return g.data[p.Row][p.Col]
}

func (g *Grid) Set(p Position, val rune) {
	g.data[p.Row][p.Col] = val
}

type Position struct {
	Row int
	Col int
}

type Direction struct {
	Dx int
	Dy int
}

func (d Direction) Scale(n int) Direction {
	return Direction{Dx: d.Dx * n, Dy: d.Dy * n}
}

var (
	Up    = Direction{Dx: 0, Dy: -1}
	Down  = Direction{Dx: 0, Dy: 1}
	Left  = Direction{Dx: -1, Dy: 0}
	Right = Direction{Dx: 1, Dy: 0}
)

func TransForm(p Position, d Direction) Position {
	return Position{Row: p.Row + d.Dy, Col: p.Col + d.Dx}
}
