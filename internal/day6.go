package internal

import (
	"fmt"
	"maps"
	"slices"
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
	empty    GridContent = '.'
	obstacle GridContent = '#'
	guard    GridContent = '^' // Can assume facing up. This is transient, only used for initialization
	invalid  GridContent = '@' // Can use any symbol for this
)

type Position struct {
	row, col int
}

type GuardState struct {
	position    Position
	orientation Orientation
}

func Day6() (int, int) {
	mapRows := Lines("./input/day6.txt")
	var guardState GuardState
	visitedPositions := make(map[Position][]Orientation)
	var initialGuardState GuardState

	var mapContents [][]GridContent
	for r, mapRow := range mapRows {
		var row []GridContent
		for c := 0; c < len(mapRow); c++ {
			content := GridContent(mapRow[c])
			if GridContent(mapRow[c]) == guard {
				position := Position{r, c}
				guardState = GuardState{position, up}
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
	delete(visitedPositions, initialGuardState.position)
	origVisitedPositions := visitedPositions
	for position := range maps.Keys(origVisitedPositions) {
		guardState = initialGuardState
		visitedPositions = make(map[Position][]Orientation)
		setVisited(visitedPositions, guardState)
		mapContents[position.row][position.col] = obstacle //Add a new obstacle
		_, loopEncountered := walkMaze(&guardState, mapContents, visitedPositions, rows, cols)
		if loopEncountered {
			part2++
		}
		mapContents[position.row][position.col] = empty
	}

	return part1, part2
}

func walkMaze(guardState *GuardState,
	mapContents [][]GridContent,
	visitedPositions map[Position][]Orientation,
	rows, cols int) (int, bool) {
	next := ahead(guardState, mapContents, rows, cols)
	steps := 1 //includes staring point
	loopEncountered := false
	for next != invalid {
		switch next {
		case empty:
			move(guardState)
			orientations, visited := visitedPositions[guardState.position]
			if !visited {
				//Not been here before - increment part1 answer
				steps++
				setVisited(visitedPositions, *guardState)
			} else if slices.Contains(orientations, guardState.orientation) {
				//Guard already been here and with same orientation - this is a loop!
				loopEncountered = true
				return steps, loopEncountered
			}
		case obstacle:
			turn(guardState)
		default:
			panic(fmt.Sprintf("Bad grid content: %v", next))
		}
		next = ahead(guardState, mapContents, rows, cols)
	}

	return steps, loopEncountered
}

func setVisited(visited map[Position][]Orientation, state GuardState) {
	visited[state.position] = append(visited[state.position], state.orientation)
}

func ahead(state *GuardState, mapContents [][]GridContent, rows, cols int) GridContent {
	var r, c int
	switch state.orientation {
	case up:
		r = state.position.row - 1
		c = state.position.col
	case right:
		r = state.position.row
		c = state.position.col + 1
	case down:
		r = state.position.row + 1
		c = state.position.col
	case left:
		r = state.position.row
		c = state.position.col - 1
	default:
		panic(fmt.Sprintf("Bad orientation: %v", state.orientation))
	}

	if !withinBounds(r, c, rows, cols) {
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
		state.position.row--
	case right:
		state.position.col++
	case down:
		state.position.row++
	case left:
		state.position.col--
	default:
		panic(fmt.Sprintf("Bad orientation: %v", state.orientation))
	}
}
