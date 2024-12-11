package internal

import (
	"fmt"
	"strconv"
)

type Cell struct {
	pos    Position
	height int
}

func (c *Cell) String() string {
	return fmt.Sprintf("Cell {%v %v}", c.pos, c.height)
}

func Day10() (int, int) {
	var part1, part2 int

	lines := Lines("./input/day10.txt")
	rows := len(lines)
	cols := len(lines[0])
	grid, trailheads := parseInput(lines)
	ratings := make([]int, len(trailheads))
	for i, trailhead := range trailheads {
		seen := make(map[Position]bool)
		cellQ := newQ[*Cell]()
		cellQ.Append(trailhead)
		for cellQ.Len() > 0 {
			cell := cellQ.PopFront()
			dPos := []Position{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
			for _, dP := range dPos {
				newPos := Position{cell.pos.row + dP.row, cell.pos.col + dP.col}
				if withinBoundsPos(newPos, rows, cols) && grid[newPos].height == cell.height+1 {
					if grid[newPos].height == 9 {
						ratings[i]++
						seen[newPos] = true
					} else {
						//Stick it in the queue
						cellQ.Append(grid[newPos])
					}
				}
			}
		}
		part1 += len(seen)
	}

	for _, rating := range ratings {
		part2 += rating
	}

	return part1, part2
}

func parseInput(lines []string) (map[Position]*Cell, []*Cell) {
	grid := make(map[Position]*Cell)
	trailheads := make([]*Cell, 0)
	for r, line := range lines {
		for c, char := range line {
			height, err := strconv.Atoi(string(char))
			if err != nil {
				height = 20
			}
			pos := Position{r, c}
			cell := Cell{pos: pos, height: height}
			grid[pos] = &cell
			if height == 0 {
				//trailhead
				trailheads = append(trailheads, &cell)
			}
		}
	}

	return grid, trailheads
}
