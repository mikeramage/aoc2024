package day15

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/mikeramage/aoc2024/position"
	"github.com/mikeramage/aoc2024/utils"
)

type Direction byte

type GridContent byte

const (
	box       GridContent = 'O'
	largeBoxL GridContent = '['
	largeBoxR GridContent = ']'
	robot     GridContent = '@'
	wall      GridContent = '#'
	empty     GridContent = '.'
	sUp       Direction   = '^'
	sDown     Direction   = 'v'
	sLeft     Direction   = '<'
	sRight    Direction   = '>'
)

func Day15() (int, int) {
	var part1, part2 int
	lines := utils.Lines("./input/day15.txt")
	var remainder []string
	var robotPosition, robotPosition2 position.Position

	var grid [][]GridContent
	var grid2 [][]GridContent
outer:
	for r, line := range lines {
		if strings.TrimSpace(line) == "" {
			remainder = lines[r+1:]
			break outer
		}
		var row []GridContent
		var row2 []GridContent
		for c, content := range line {
			if GridContent(content) == robot {
				robotPosition = position.Position{Row: r, Col: c}
				robotPosition2 = position.Position{Row: r, Col: 2 * c}
				row = append(row, empty)
				for i := 0; i < 2; i++ {
					row2 = append(row2, empty)
				}
			} else {
				row = append(row, GridContent(content))
				switch GridContent(content) {
				case wall:
					for i := 0; i < 2; i++ {
						row2 = append(row2, wall)
					}
				case empty:
					for i := 0; i < 2; i++ {
						row2 = append(row2, empty)
					}
				case box:
					row2 = append(row2, largeBoxL)
					row2 = append(row2, largeBoxR)
				}
			}
		}
		grid = append(grid, row)
		grid2 = append(grid2, row2)
	}

	var moves []Direction
	for _, line := range remainder {
		for _, move := range line {
			moves = append(moves, Direction(move))
		}
	}

	visualiseWarehouse(grid, robotPosition)
	visualiseWarehouse(grid2, robotPosition2)

	for _, move := range moves {
		// visualiseWarehouse(grid2, robotPosition2)
		switch move {
		case sUp:
			doMove(&robotPosition, &grid, -1, 0)
			doMove(&robotPosition2, &grid2, -1, 0)
		case sRight:
			doMove(&robotPosition, &grid, 0, 1)
			doMove(&robotPosition2, &grid2, 0, 1)
		case sDown:
			doMove(&robotPosition, &grid, 1, 0)
			doMove(&robotPosition2, &grid2, 1, 0)
		case sLeft:
			doMove(&robotPosition, &grid, 0, -1)
			doMove(&robotPosition2, &grid2, 0, -1)
		}
	}

	// visualiseWarehouse(grid)
	visualiseWarehouse(grid2, robotPosition2)

	for r, row := range grid {
		for c, content := range row {
			if content == box {
				part1 += 100*r + c
			}
		}
	}

	for r, row := range grid2 {
		for c, content := range row {
			if content == largeBoxL {
				part2 += 100*r + c
			}
		}
	}

	return part1, part2
}

func doMove(robotPosition *position.Position, grid *[][]GridContent, dr, dc int) {
	r, c := robotPosition.Row, robotPosition.Col
	gridContent := (*grid)[r+dr][c+dc]
	switch gridContent {
	case empty:
		robotPosition.Row += dr
		robotPosition.Col += dc
	case wall:
		//Do nothing
	case box:
		//The complicated case. Keep going in direction dr, dc until we hit an empty or a wall
		numBoxes := 0
		for gridContent == box {
			numBoxes++
			gridContent = (*grid)[r+(numBoxes+1)*dr][c+(numBoxes+1)*dc]
		}

		switch gridContent {
		case wall:
			//No-op. Can't push the boxes
		case empty:
			robotPosition.Row += dr
			robotPosition.Col += dc
			(*grid)[r+dr][c+dc] = empty
			for i := 0; i < numBoxes; i++ {
				(*grid)[r+dr*(i+2)][c+dc*(i+2)] = box
			}
		default:
			panic("Expected either wall or empty!")
		}
	case largeBoxL, largeBoxR:
		//The _really_ complicated case (at least for up/down)
		if dr == 0 {
			//Horizontal move - okay ish!
			doLargeBoxHorizontalMove(gridContent, grid, r, c, dc, robotPosition)
		} else {
			//Vertical move - urgh!
			doLargeBoxVerticalMove(gridContent, grid, r, c, dr, robotPosition)
		}
	}
}

func doLargeBoxVerticalMove(gridContent GridContent, grid *[][]GridContent, r, c, dr int, robotPosition *position.Position) {
	moveMap := make(map[int]map[int]GridContent)
	rowMap := make(map[int]GridContent)
	moveMap[r+dr] = rowMap
	moveMap[r+dr][c] = gridContent
	completeBox(gridContent, moveMap, grid, r+dr, c)
	if chainBoxes(grid, moveMap, r, dr) {
		//Can move so do so
		robotPosition.Row += dr
		moveBoxes(grid, moveMap, dr)
	}
}

func moveBoxes(grid *[][]GridContent, moveMap map[int]map[int]GridContent, dr int) {
	// Start with the furthest away row and move a row up (or down) into empty space one at a time,
	// temporarily replacing the moved box part with empty space (which may be filled in at the
	// next iteration with another box part). When we're pushing down the way (direction of increasing row)
	// we want to sort the rows from largest to smallest, otherwise smallest to largest
	sortedRows := slices.Sorted(maps.Keys(moveMap))
	if dr > 0 {
		slices.Reverse(sortedRows)
	}

	for _, r := range sortedRows {
		for c := range maps.Keys(moveMap[r]) {
			(*grid)[r+dr][c] = moveMap[r][c]
			(*grid)[r][c] = empty
		}
	}
}

// Search up (or down) the rows along the chain of boxes until we reach the end and have empty space
// to move into (return true) or hit a wall (return false)
func chainBoxes(grid *[][]GridContent, moveMap map[int]map[int]GridContent, r, dr int) bool {
	currentRow, nextRow := r+dr, r+2*dr
	for {
		canMove := true
		for c := range maps.Keys(moveMap[currentRow]) {
			gridContent := (*grid)[nextRow][c] //the thing above or below the current item
			switch gridContent {
			case wall: //Box blocked - return false.
				return false
			case empty: //Empty - we might be able to move - check next in the list.
				continue
			case largeBoxL, largeBoxR:
				//Box part - can't move right now; need to check the next level - add a new row to the map if not already done so and add the box parts
				canMove = false
				if moveMap[nextRow] == nil {
					rowMap := make(map[int]GridContent)
					moveMap[nextRow] = rowMap
				}
				moveMap[nextRow][c] = gridContent
				completeBox(gridContent, moveMap, grid, nextRow, c)
			}
		}

		if canMove {
			return true
		}

		currentRow = nextRow
		nextRow += dr
	}
}

func completeBox(gridContent GridContent, moveMap map[int]map[int]GridContent, grid *[][]GridContent, r, c int) {
	switch gridContent {
	case largeBoxL:
		moveMap[r][c+1] = (*grid)[r][c+1] //must be largeBoxR, but not going to bother checking this
	case largeBoxR:
		moveMap[r][c-1] = (*grid)[r][c-1] //must be largeBoxL
	default:
		panic("Should only be called when gridContent is part of a box!")
	}
}

func doLargeBoxHorizontalMove(gridContent GridContent, grid *[][]GridContent, r, c, dc int, robotPosition *position.Position) {
	//Not complicated - just have to jump right 2 at a time.
	numBoxes := 0
	boxPart := gridContent //Must be large box part L or R but can't be bothered to check
	for gridContent == boxPart {
		numBoxes++
		gridContent = (*grid)[r][c+(2*numBoxes+1)*dc]
	}

	switch gridContent {
	case wall:
		//No-op. Can't push the boxes
	case empty:
		robotPosition.Col += dc
		(*grid)[r][c+dc] = empty
		for i := 0; i < numBoxes; i++ {
			if dc > 0 { //Pushing right
				(*grid)[r][c+dc*(2*i+2)] = largeBoxL
				(*grid)[r][c+dc*(2*i+3)] = largeBoxR
			} else { //Pushing left
				(*grid)[r][c+dc*(2*i+2)] = largeBoxR
				(*grid)[r][c+dc*(2*i+3)] = largeBoxL
			}
		}
	default:
		fmt.Println(r, c, dc, gridContent)
		panic("Expected either wall or empty!")
	}
}

func visualiseWarehouse(grid [][]GridContent, robotPosition position.Position) {
	for r, row := range grid {
		for c, col := range row {
			if r == robotPosition.Row && c == robotPosition.Col {
				fmt.Printf("@")
			} else {
				fmt.Printf("%v", string(col))
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
