package day19

import (
	"strings"

	"github.com/mikeramage/aoc2024/utils"
)

func Day19() (int, int) {
	var part1, part2 int

	lines := utils.Lines("./input/day19.txt")
	var patterns []string
	var designs []string

	parsePatterns := true
	for _, line := range lines {
		if parsePatterns {
			patterns = strings.Split(line, ", ")
			parsePatterns = false
		} else if strings.TrimSpace(line) != "" {
			designs = append(designs, line)
		}
	}

	cache := make(map[string]int)
	for _, design := range designs {
		possible, ways := canMake(design, patterns, cache)
		if possible {
			part1++
			part2 += ways
		}
	}

	return part1, part2
}

func canMake(design string, patterns []string, cache map[string]int) (bool, int) {
	var possible bool
	ways, exists := cache[design]
	possible = ways > 0

	if len(design) == 0 {
		return true, 1
	}

	if possible {
		return true, ways
	}

	if exists && !possible {
		return false, ways
	}

	remaining := design
	for _, pattern := range patterns {
		if strings.HasPrefix(remaining, pattern) {
			subPossible, subWays := canMake(remaining[len(pattern):], patterns, cache)
			if subPossible {
				cache[remaining[len(pattern):]] = subWays
				ways += subWays
				possible = true
			} else {
				cache[remaining[len(pattern):]] = 0
			}
		}
	}

	if possible {
		return true, ways
	}

	return false, 0
}
