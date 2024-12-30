package day25

import (
	"strings"

	"github.com/mikeramage/aoc2024/utils"
)

type Lock struct {
	heights []int
}

type Key struct {
	heights []int
}

func Day25() (int, int) {
	var part1, part2 int
	var inLock, inKey bool
	lines := utils.Lines("./input/day25.txt")
	var locks []Lock
	var keys []Key

	currentHeights := make([]int, 5)
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			if inLock {
				locks = append(locks, Lock{heights: currentHeights})
				currentHeights = make([]int, 5)
				inLock = false
				continue
			}

			if inKey {
				// We've overcounted by 1 because of the final row.
				for i := range currentHeights {
					currentHeights[i]--
				}
				keys = append(keys, Key{heights: currentHeights})
				currentHeights = make([]int, 5)
				inKey = false
				continue
			}
		}
		if !inLock && !inKey {
			if line[0] == '.' {
				inKey = true
				continue

			} else if line[0] == '#' {
				inLock = true
				continue
			}
		}

		for i, char := range line {
			if char == '#' {
				currentHeights[i]++
			}
		}
	}

	if inLock {
		locks = append(locks, Lock{heights: currentHeights})
	} else if inKey {
		for i := range currentHeights {
			currentHeights[i]--
		}
		keys = append(keys, Key{heights: currentHeights})
	}

	for _, key := range keys {
		for _, lock := range locks {
			fit := true
			for i, h := range key.heights {
				if h+lock.heights[i] > 5 {
					fit = false
				}
			}
			if fit {
				part1++
			}
		}
	}

	return part1, part2
}
