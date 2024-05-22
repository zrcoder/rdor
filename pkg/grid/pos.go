package grid

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
