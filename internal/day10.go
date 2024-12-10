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
	scores := make([]int, len(trailheads))
	ratings := make([]int, len(trailheads))
	for i, trailhead := range trailheads {
		cellQ := newQ[*Cell]()
		cellQ.Append(trailhead)
		seen := make(map[Position]bool)
		ratings[i]++ //There's always at least one path
		for cellQ.Len() > 0 {
			cell := cellQ.PopFront()
			seen[cell.pos] = true
			dPos := []Position{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
			numBranches := 0
			for _, dP := range dPos {
				newPos := Position{cell.pos.row + dP.row, cell.pos.col + dP.col}
				if withinBoundsPos(newPos, rows, cols) && grid[newPos].height == cell.height+1 {
					//Got a cell that's the right height, and it's inside the grid.
					//For considering unique paths that lead to the 9, we count branchings even if we've seen it before.
					//Add to the queue unless it's a 9, in which case we mark it as seen and increment the score for this trailhead
					numBranches++
					if grid[newPos].height == 9 && !seen[newPos] {
						scores[i]++
						seen[newPos] = true
					} else {
						//Stick it in the queue
						cellQ.Append(grid[newPos])
					}
				}
			}
			if numBranches != 0 { //If the paths fork there are numBranches - 1 new paths (one of them is the existing path)
				ratings[i] += numBranches - 1
			} else if cell.height != 9 {
				//If we hit a dead end before the 9 it's not a trail. Subtract 1 as we've overcounted.
				ratings[i]--
			}
		}
	}

	for i, score := range scores {
		part1 += score
		part2 += ratings[i]
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
