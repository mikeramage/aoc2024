package day2

import (
	"slices"
	"strconv"
	"strings"

	"github.com/mikeramage/aoc2024/utils"
)

func Day2() (int, int) {
	lines := utils.Lines("./input/day2.txt")
	var levels [][]int

	for _, line := range lines {
		var row []int
		for _, level := range strings.Fields(line) {
			l, _ := strconv.Atoi(level)
			row = append(row, l)
		}
		levels = append(levels, row)
	}

	part1 := safeReports(levels, false) //Swap with safeReportsAlt for a simper (but less optimized) algorithm
	part2 := safeReports(levels, true)

	return part1, part2
}

func safeReports(levels [][]int, tolerateOneBad bool) int {
	var safe int

	for _, row := range levels {
		ok, index := checkRow(row, -1)

		if ok {
			safe++
		} else if tolerateOneBad {
			for i := index - 1; i <= index+1; i++ {
				success, _ := checkRow(row, i)
				if success {
					safe++
					break
				}
			}
		}
	}

	return safe
}

func checkRow(row []int, indexToSkip int) (bool, int) {
	var current, next int
	next = 1
	if current == indexToSkip {
		current += 1
		next += 1
	} else if next == indexToSkip {
		next += 1
	}

	if next >= len(row) { // Shouldn't hit this, but in case we do ...
		return true, current
	}

	sign := row[next] - row[current]
	for next < len(row) {
		if !checkPair(row, current, next, sign) {
			return false, current
		}

		current = next
		next++

		if next == indexToSkip {
			next++
		}
	}

	return true, current
}

func checkPair(row []int, current, next, sign int) bool {
	diff := row[next] - row[current]
	if (sign < 0 && diff > 0) || (sign > 0 && diff < 0) || utils.Abs(diff) < 1 || utils.Abs(diff) > 3 {
		return false
	}

	return true
}

//lint:ignore U1000 Don't care about unused code
func safeReportsAlt(levels [][]int, tolerateOneBad bool) int {
	var safe int

	for _, row := range levels {
		ok := checkRowAlt(row)

		if ok {
			safe++
		} else if tolerateOneBad {
			for i := 0; i < len(row); i++ {
				partial := slices.Concat(row[:i], row[i+1:])
				ok = checkRowAlt(partial)
				if ok {
					safe++
					break
				}
			}
		}
	}

	return safe
}

func checkRowAlt(row []int) bool {
	sortedRow := slices.Sorted(slices.Values(row))
	revSortedRow := make([]int, len(sortedRow))
	copy(revSortedRow, sortedRow)
	slices.Reverse(revSortedRow)

	inc := true
	for i, v := range row {
		if v != sortedRow[i] {
			inc = false
			break
		}
	}

	dec := true
	if !inc {
		for i, v := range row {
			if v != revSortedRow[i] {
				dec = false
				break
			}
		}
	}

	if !inc && !dec {
		return false
	}

	for i := 0; i < len(row)-1; i++ {
		absDiff := utils.Abs(row[i+1] - row[i])
		if absDiff < 1 || absDiff > 3 {
			return false
		}
	}

	return true

}
