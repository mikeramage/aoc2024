package day18

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mikeramage/aoc2024/position"
	"github.com/mikeramage/aoc2024/queue"
	"github.com/mikeramage/aoc2024/utils"
)

type Space struct {
	Position position.Position
	Parent   *Space
}

func Day18() (int, int) {
	var part1, part2 int

	lines := utils.Lines("./input/day18.txt")
	rows, cols := 71, 71
	kB := 1024
	positionToTimeCorrupted := make(map[position.Position]int)
	timeCorruptedToPosition := make(map[int]position.Position) //inverse map

	for time, line := range lines {
		coords := strings.Split(line, ",")
		row, _ := strconv.Atoi(coords[1]) // x, y -> c, r
		col, _ := strconv.Atoi(coords[0])
		pos := position.NewPosition(row, col)
		positionToTimeCorrupted[pos] = time
		timeCorruptedToPosition[time] = pos
	}

	//Part 1 simple BFS, checking frontier for existing nodes
	solution, reached := doBFS(positionToTimeCorrupted, rows, cols, kB)

	visualiseSpace(positionToTimeCorrupted, kB, rows, cols, solution)

	if reached {
		parent := solution.Parent
		for parent != nil {
			part1++
			parent = parent.Parent
		}
	}

	//Part 2 - binary chop time until we find the earliest time that fails, and look up the
	//reverse corrupted map by that time for the position of the fatal corruption (and remember to
	// convert row <-> y, col <-> x)
	var finalNode *Space
	var lastGoodSolution, firstBlockedSolution *Space
	low, high := 0, len(positionToTimeCorrupted)
	for low < high {
		mid := (low + high) / 2
		finalNode, reached = doBFS(positionToTimeCorrupted, rows, cols, mid)
		if reached {
			low = mid + 1
			lastGoodSolution = finalNode
		} else {
			high = mid
			firstBlockedSolution = finalNode
		}
	}

	//Low is 1 more than the time input to BFS for the last good solution. When doing BFS, the time input
	//is 1 greater than the time of the last corrupted byte. So the corrupted byte fell at time = the time
	//input for the last successful solution, i.e. low - 1
	blockingTime := low - 1

	if lastGoodSolution == nil {
		panic("No solution found")
	}

	blockingPosition := timeCorruptedToPosition[blockingTime]
	fmt.Println("Last good solution\n------------------------------------------------------")
	visualiseSpace(positionToTimeCorrupted, blockingTime, rows, cols, lastGoodSolution)
	fmt.Println("Blocked path\n------------------------------------------------------------")
	visualiseSpace(positionToTimeCorrupted, blockingTime+1, rows, cols, firstBlockedSolution)

	fmt.Println("x:", blockingPosition.Col, "y:", blockingPosition.Row)

	return part1, part2
}

func doBFS(positionToTimeCorrupted map[position.Position]int, rows, cols, time int) (*Space, bool) {
	seen := make(map[position.Position]bool)
	frontierMap := make(map[position.Position]bool)
	frontier := queue.NewQueue[*Space]()
	initialSpace := &Space{Position: position.NewPosition(0, 0), Parent: nil}
	frontier.Append(initialSpace)
	var space *Space
	goalPosition := position.NewPosition(rows-1, cols-1)
	reached := false
	for frontier.Len() > 0 {
		space = frontier.PopFront()
		delete(frontierMap, space.Position)
		if space.Position.Equal(goalPosition) {
			reached = true
			break
		}
		seen[space.Position] = true
		for _, direction := range position.DirectionsPos {
			pos := position.Add(space.Position, direction)
			corruptedTime, blocked := positionToTimeCorrupted[pos]
			if !seen[pos] &&
				(!blocked || corruptedTime >= time) &&
				position.WithinBoundsPos(pos, rows, cols) &&
				!frontierMap[pos] {
				frontier.Append(&Space{Position: pos, Parent: space})
				frontierMap[pos] = true
			}
		}
	}

	return space, reached
}

func visualiseSpace(positionToTimeCorrupted map[position.Position]int, after, rows, cols int, finalNode *Space) {

	path := make(map[position.Position]bool)
	if finalNode != nil {
		path[finalNode.Position] = true
		parent := finalNode.Parent
		for parent != nil {
			path[parent.Position] = true
			parent = parent.Parent
		}
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			pos := position.NewPosition(r, c)
			time, exists := positionToTimeCorrupted[pos]
			if exists && time < after {
				fmt.Printf("#")
			} else if path[pos] {
				fmt.Printf("O")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
