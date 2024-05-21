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

func (g *Grid) SetData(data [][]rune) {
	g.data = data
}

func Copy(g *Grid) *Grid {
	data := make([][]rune, len(g.data))
	for r, row := range g.data {
		data[r] = make([]rune, len(row))
		copy(data[r], row)
	}
	return &Grid{data: data}
}

func (g *Grid) Copy(gg *Grid) {
	for i, row := range gg.data {
		copy(g.data[i], row)
	}
}

func (g *Grid) Equal(f *Grid) bool {
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

func (g *Grid) String() string {
	buf := strings.Builder{}
	for _, line := range g.data {
		buf.WriteString(string(line))
		buf.WriteByte('\n')
	}
	return buf.String()
}

type Position struct {
	Row int
	Col int
}

type Direction struct {
	Dx int
	Dy int
}

func (p Position) TransForm(d Direction) Position {
	return TransForm(p, d)
}

func (d Direction) Scale(n int) Direction {
	return Direction{Dx: d.Dx * n, Dy: d.Dy * n}
}

func (d Direction) Opposite() Direction {
	return Direction{Dx: -d.Dx, Dy: -d.Dy}
}

func (d Direction) Rotate() Direction {
	return Direction{Dx: -d.Dy, Dy: d.Dx}
}

var (
	Up        = Direction{Dx: 0, Dy: -1}
	Down      = Direction{Dx: 0, Dy: 1}
	Left      = Direction{Dx: -1, Dy: 0}
	Right     = Direction{Dx: 1, Dy: 0}
	UpLeft    = Direction{Dx: -1, Dy: -1}
	UpRight   = Direction{Dx: 1, Dy: -1}
	DownLeft  = Direction{Dx: -1, Dy: 1}
	DownRight = Direction{Dx: 1, Dy: 1}

	NormalDirections = []Direction{Up, Down, Left, Right}
	AllDirections    = []Direction{Up, Down, Left, Right, UpLeft, UpRight, DownLeft, DownRight}
)

func TransForm(p Position, d Direction) Position {
	return Position{Row: p.Row + d.Dy, Col: p.Col + d.Dx}
}
