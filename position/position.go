package position

type Position struct {
	Row, Col int
}

func NewPosition(row, col int) Position {
	return Position{Row: row, Col: col}
}

func ComparePositions(a, b Position) int {
	if a.Row > b.Row {
		return 1
	}

	if a.Row < b.Row {
		return -1
	}

	if a.Col > b.Col {
		return 1
	}

	if a.Col < b.Col {
		return -1
	}

	return 0
}

func Add(p Position, q Position) Position {
	r := p.Row + q.Row
	c := p.Col + q.Col
	return NewPosition(r, c)
}

func (p Position) Equal(q Position) bool {
	return p.Row == q.Row && p.Col == q.Col
}

var DirectionsPos = []Position{{Row: 1, Col: 0}, {Row: -1, Col: 0}, {Row: 0, Col: 1}, {Row: 0, Col: -1}}
var Directions = [][]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
var DirectionsDiag = [][]int{{-1, 0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}, {-1, -1}}

func WithinBounds(r, c, rows, cols int) bool {
	return r < rows && r >= 0 && c < cols && c >= 0
}

func WithinBoundsPos(p Position, rows, cols int) bool {
	return p.Row < rows && p.Row >= 0 && p.Col < cols && p.Col >= 0
}
