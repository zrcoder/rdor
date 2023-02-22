package last

import "github.com/zrcoder/rdor/pkg/grid"

type pathStack struct {
	path []grid.Direction
}

func (ps *pathStack) empty() bool           { return len(ps.path) == 0 }
func (ps *pathStack) push(d grid.Direction) { ps.path = append(ps.path, d) }
func (ps *pathStack) pop() grid.Direction {
	n := len(ps.path)
	res := ps.path[n-1]
	ps.path = ps.path[:n-1]
	return res
}
