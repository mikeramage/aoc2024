package day6

import (
	"fmt"
	"maps"
	"slices"

	"github.com/mikeramage/aoc2024/position"
	"github.com/mikeramage/aoc2024/utils"
)

type Orientation int

const (
	up Orientation = iota
	right
	down
	left
)

type GridContent byte

const (
	empty   GridContent = '.'
	wall    GridContent = '#'
	guard   GridContent = '^' // Can assume facing up. This is transient, only used for initialization
	invalid GridContent = 'X' // Can use any symbol for this
)

type GuardState struct {
	pos         position.Position
	orientation Orientation
}

func Day6() (int, int) {
	mapRows := utils.Lines("./input/day6.txt")
	var guardState GuardState
	visitedPositions := make(map[position.Position][]Orientation)
	var initialGuardState GuardState

	var mapContents [][]GridContent
	for r, mapRow := range mapRows {
		var row []GridContent
		for c := 0; c < len(mapRow); c++ {
			content := GridContent(mapRow[c])
			if GridContent(mapRow[c]) == guard {
				pos := position.Position{Row: r, Col: c}
				guardState = GuardState{pos, up}
				initialGuardState = guardState //Store this for part 2 - can't put an obstacle on home square
				setVisited(visitedPositions, guardState)
				content = empty // The grid square is empty at this location (we ignore the guard)
			}
			row = append(row, content)

		}
		mapContents = append(mapContents, row)
	}

	rows := len(mapContents)
	cols := len(mapContents[0])

	part1, _ :=
		walkMaze(&guardState,
			mapContents,
			visitedPositions,
			rows,
			cols)

	var part2 int
	//Don't want to consider the initial position for placing an obstacle
	delete(visitedPositions, initialGuardState.pos)
	origVisitedPositions := visitedPositions
	for pos := range maps.Keys(origVisitedPositions) {
		guardState = initialGuardState
		visitedPositions = make(map[position.Position][]Orientation)
		setVisited(visitedPositions, guardState)
		mapContents[pos.Row][pos.Col] = wall //Add a new obstacle
		_, loopEncountered := walkMaze(&guardState, mapContents, visitedPositions, rows, cols)
		if loopEncountered {
			part2++
		}
		mapContents[pos.Row][pos.Col] = empty
	}

	return part1, part2
}

func walkMaze(guardState *GuardState,
	mapContents [][]GridContent,
	visitedPositions map[position.Position][]Orientation,
	rows, cols int) (int, bool) {
	next := ahead(guardState, mapContents, rows, cols)
	steps := 1 //includes staring point
	loopEncountered := false
	for next != invalid {
		switch next {
		case empty:
			move(guardState)
			orientations, visited := visitedPositions[guardState.pos]
			if !visited {
				//Not been here before - increment part1 answer
				steps++
				setVisited(visitedPositions, *guardState)
			} else if slices.Contains(orientations, guardState.orientation) {
				//Guard already been here and with same orientation - this is a loop!
				loopEncountered = true
				return steps, loopEncountered
			}
		case wall:
			turn(guardState)
		default:
			panic(fmt.Sprintf("Bad grid content: %v", next))
		}
		next = ahead(guardState, mapContents, rows, cols)
	}

	return steps, loopEncountered
}

func setVisited(visited map[position.Position][]Orientation, state GuardState) {
	visited[state.pos] = append(visited[state.pos], state.orientation)
}

func ahead(state *GuardState, mapContents [][]GridContent, rows, cols int) GridContent {
	var r, c int
	switch state.orientation {
	case up:
		r = state.pos.Row - 1
		c = state.pos.Col
	case right:
		r = state.pos.Row
		c = state.pos.Col + 1
	case down:
		r = state.pos.Row + 1
		c = state.pos.Col
	case left:
		r = state.pos.Row
		c = state.pos.Col - 1
	default:
		panic(fmt.Sprintf("Bad orientation: %v", state.orientation))
	}

	if !position.WithinBounds(r, c, rows, cols) {
		return invalid
	} else {
		return mapContents[r][c]
	}
}

func turn(state *GuardState) {
	switch state.orientation {
	case up:
		state.orientation = right
	case right:
		state.orientation = down
	case down:
		state.orientation = left
	case left:
		state.orientation = up
	default:
		panic(fmt.Sprintf("Bad orientation: %v", state.orientation))
	}
}

func move(state *GuardState) {
	switch state.orientation {
	case up:
		state.pos.Row--
	case right:
		state.pos.Col++
	case down:
		state.pos.Row++
	case left:
		state.pos.Col--
	default:
		panic(fmt.Sprintf("Bad orientation: %v", state.orientation))
	}
}
