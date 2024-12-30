package day20

import (
	"fmt"

	"github.com/mikeramage/aoc2024/position"
	"github.com/mikeramage/aoc2024/utils"
)

type GridContent byte

var (
	empty GridContent = '.'
	start GridContent = 'S'
	end   GridContent = 'E'
)

type GridCell struct {
	position position.Position
	distance int
	onPath   bool
	content  GridContent
}

type Cheat struct {
	start, end position.Position
}

func Day20() (int, int) {
	var part1, part2 int
	// Looks like BFS but isn't. There are no choices to make - just follow the path.
	// Once path is known, 2 seconds only allows removal of one piece of wall on the path.
	// But not just any wall. The wall has to connect 2 squares on the original path. If it connects
	// to some other part of the maze because there's just the single corridor we know it can't
	// ever lead to the exit. (Part 2 may well introduce 2 cheats in which case this is no longer
	// true and my simple solution for part 1 won't work, but let's cross that bridge). So:
	//
	// - Find the path - just follow the maze in the prescribed direction till the exit. Store each
	// grid location in a map of grid location to step number along the path. Start is step 0. At the
	// same time add and adjacent walls to a list of "candidate cheats" - i.e. if we remove the wall it's
	// possible it connects to a shorter path. Only add adjacent walls that have at least two adjacent
	// empty grid positions.
	//
	// - For each candidate cheat wall check to see if . For each one we find that connects to another
	// location on the path, consider it for removal. The time saved is the difference between the
	// step counts of the two locations + 2 for the two steps to get from the current location to the wall location
	// to the other location on the path.
	//
	// Done!
	//
	// So that worked for part 1, but not part 2 - the set of candidate cheat walls to remove is much larger.
	// Instead for part 2, for a different approach we're going to walk the path .
	// - For each square on the path, only consider cheating if the distance is > 100 from the end (small optimisation)
	// - If it's worth cheating, the current square is the cheat start. If that square is at (r, c) then possible
	// cheat ends are at (r+dr, c+dc) where 2 <= dr + dc <= 20 (about 200ish possibilities).
	// - It's a good cheat if distance(r+dr, c+dc) - distance(r, c) - dr - dc > 100
	// I'll adopt the same solution to recompute part 1 (but keep my old one for reference - my guess is the new one will
	// be worse than 20ms it takes for the old one, but let's see) ...
	//
	// ... nope, the new one takes 1ms for part 1 so my first attempt was a total waste of time. It does take 160ms+ for part2
	// which is rubbish, but can I be bothered to optimize?
	//

	var grid [][]GridCell
	// pathMap := make(map[position.Position]int)
	// candidateCheats := make(map[position.Position]bool)
	lines := utils.Lines("./input/day20.txt")
	minSaving := 100
	cheatDurationPart1 := 2
	cheatDurationPart2 := 20
	var startPos, endPos position.Position
	// rows := len(lines)
	// cols := len(lines[0])
	for r, line := range lines {
		var row []GridCell
		for c, content := range line {
			cell := GridCell{position: position.NewPosition(r, c), distance: -1, onPath: false, content: GridContent(content)}
			row = append(row, cell)
			if GridContent(content) == start {
				startPos = cell.position
			} else if GridContent(content) == end {
				endPos = cell.position
			}
		}
		grid = append(grid, row)
	}

	// followPath1stAttempt(grid, candidateCheats, startPos, endPos)
	// bestCheats := findBestCheats1stAttempt(candidateCheats, grid, minSaving)
	path := followPath(grid, startPos, endPos)
	bestCheatsPart1 := findBestCheats(path, grid, cheatDurationPart1, minSaving)
	bestCheatsPart2 := findBestCheats(path, grid, cheatDurationPart2, minSaving)
	part1 = len(bestCheatsPart1)
	part2 = len(bestCheatsPart2)
	visualiseTrack(grid, bestCheatsPart1)
	visualiseTrack(grid, bestCheatsPart2)

	return part1, part2
}

func followPath(grid [][]GridCell, startPos, endPos position.Position) []*GridCell {
	distance := 0
	currentPos := startPos
	var path []*GridCell
	var cell *GridCell

	for currentPos != endPos {
		cell = &grid[currentPos.Row][currentPos.Col]
		path = append(path, cell)
		cell.distance = distance
		cell.onPath = true
		var nextPos position.Position
		for _, dir := range position.DirectionsPos {
			newPos := position.Add(currentPos, dir)
			if !position.WithinBoundsPos(newPos, len(grid), len(grid[0])) {
				continue
			}
			newCell := &grid[newPos.Row][newPos.Col]
			if !newCell.onPath && (newCell.content == empty || newCell.content == end) {
				nextPos = newPos
			}
		}
		distance++
		currentPos = nextPos
	}

	endCell := &grid[currentPos.Row][currentPos.Col]
	path = append(path, endCell)
	endCell.distance = distance
	endCell.onPath = true

	return path
}

func findBestCheats(path []*GridCell, grid [][]GridCell, cheatDuration, minSaving int) map[Cheat]int {
	bestCheats := make(map[Cheat]int)
	count := 0

	for _, cheatStart := range path {
		count++
		for duration := 2; duration <= cheatDuration; duration++ {
			for dr := -duration; dr <= duration; dr++ {
				dcAbs := duration - utils.Abs(dr)
				dcs := []int{dcAbs, -dcAbs}
				for _, dc := range dcs {
					endPos := position.NewPosition(cheatStart.position.Row+dr, cheatStart.position.Col+dc)
					if position.WithinBoundsPos(endPos, len(grid), len(grid[0])) {
						cheatEnd := &grid[endPos.Row][endPos.Col]
						timeSaved := cheatEnd.distance - cheatStart.distance - duration
						if cheatEnd.onPath && timeSaved >= minSaving {
							cheat := Cheat{start: cheatStart.position, end: cheatEnd.position}
							bestYet, exists := bestCheats[cheat]
							if !exists || timeSaved > bestYet {
								bestCheats[cheat] = timeSaved
							}
						}
					}
				}
			}
		}
	}

	return bestCheats
}

// func followPath1stAttempt(grid [][]GridCell, candidateCheats map[position.Position]bool, startPos, endPos position.Position) {
// 	distance := 0
// 	currentPos := startPos
// 	var cell *GridCell

// 	for currentPos != endPos {
// 		cell = &grid[currentPos.Row][currentPos.Col]
// 		cell.distance = distance
// 		cell.onPath = true
// 		var nextPos position.Position
// 		for _, dir := range position.DirectionsPos {
// 			newPos := position.Add(currentPos, dir)
// 			if !position.WithinBoundsPos(newPos, len(grid), len(grid[0])) {
// 				continue
// 			}
// 			newCell := &grid[newPos.Row][newPos.Col]
// 			if !newCell.onPath && (newCell.content == empty || newCell.content == end) {
// 				nextPos = newPos
// 			} else if newCell.content == wall {
// 				candidateCheats[newPos] = true
// 			}
// 		}
// 		distance++
// 		currentPos = nextPos
// 	}

// 	endCell := &grid[currentPos.Row][currentPos.Col]
// 	endCell.distance = distance
// 	endCell.onPath = true
// }

// func findBestCheats1stAttempt(candidateCheats map[position.Position]bool, grid [][]GridCell, minSaving int) map[Cheat]int {
// 	bestCheats := make(map[Cheat]int)

// 	for candidateCheat := range candidateCheats {
// 		var connections []GridCell
// 		for _, dir := range position.DirectionsPos {
// 			pos := position.Add(candidateCheat, dir)
// 			if !position.WithinBoundsPos(pos, len(grid), len(grid[0])) {
// 				continue
// 			}
// 			onPath := grid[pos.Row][pos.Col].onPath
// 			if onPath { //if it's on the path then the distance is guaranteed accurate and it's guaranteed empty or S or E
// 				connections = append(connections, grid[pos.Row][pos.Col])
// 			}
// 		}
// 		slices.SortFunc(connections, func(a, b GridCell) int {
// 			return cmp.Compare(a.distance, b.distance)
// 		})

// 		//Pairwise compare connections - each cheat is uniquely defined by start and end position. It's possible for
// 		//more than 1 cheat to start on the same grid cell. Consider the maximal arrangement
// 		// ...   w e w   ...
// 		// ... w w E w w ...
// 		// ... e E W E e ...
// 		// ... w w E w w ...
// 		// ...   w e w   ...
// 		// E/e for empty, W/w for wall. The big W cheat in the middle can allow the 4 different big Es to see each other. It's possible
// 		// that they all differ by > 100 - you could have 100, 300, 500, 900 - so that's a max of 6 possible unique start
// 		// and end combinations (of course the w's diagonally adjacent would open up a subset of those paths too, and
// 		// we could probably just evaluate horizontally and vertically opposite pairs, but let's just evaluate the general case).
// 		//
// 		// Note the +2. To get to the cheat end square from the cheat start is 2 steps. Remember that the distance decreases towards the end
// 		for i, connection := range connections {
// 			for j := i + 1; j < len(connections); j++ {
// 				timeSaved := connections[j].distance - connection.distance - 2
// 				if timeSaved >= 100 {
// 					currentBest, exists := bestCheats[Cheat{connection.position, connections[j].position}]
// 					if !exists || timeSaved > currentBest {
// 						bestCheats[Cheat{connection.position, connections[j].position}] = timeSaved
// 					}
// 				}
// 			}
// 		}

// 	}

// 	return bestCheats
// }

func visualiseTrack(grid [][]GridCell, bestCheats map[Cheat]int) {

	cheatStarts, cheatEnds := make(map[position.Position]bool), make(map[position.Position]bool)
	for cheat := range bestCheats {
		cheatStarts[cheat.start] = true
		cheatEnds[cheat.end] = true
	}

	for r, row := range grid {
		for c, cell := range row {
			pos := position.NewPosition(r, c)
			if cheatStarts[pos] {
				fmt.Printf("s")
			} else if cheatEnds[pos] {
				fmt.Printf("e")
			} else {
				fmt.Printf("%v", string(cell.content))
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
