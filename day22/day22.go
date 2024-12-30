package day22

import (
	"maps"
	"slices"
	"strconv"

	"github.com/mikeramage/aoc2024/utils"
)

func Day22() (int, int) {
	var part1, part2 int

	lines := utils.Lines("./input/day22.txt")
	var allPrices [][]int
	var allDiffs [][]int
	for _, line := range lines {
		secretNumber, _ := strconv.Atoi(line)
		var prices = []int{secretNumber % 10}
		var diffs []int
		for j := 0; j < 2000; j++ {
			secretNumber = nextSecretNumber(secretNumber)
			prices = append(prices, secretNumber%10)
			diffs = append(diffs, prices[j+1]-prices[j])
		}
		allPrices = append(allPrices, prices)
		allDiffs = append(allDiffs, diffs)
		part1 += secretNumber
	}

	type MultiDiff struct {
		d1, d2, d3, d4 int
	}

	priceMap := make(map[MultiDiff]int)

	for i, diffs := range allDiffs {
		seen := make(map[MultiDiff]bool)
		for j := range diffs {
			if j > len(diffs)-4 {
				//Will reach end of diffs before end of sequence
				break
			}
			md := MultiDiff{diffs[j], diffs[j+1], diffs[j+2], diffs[j+3]}
			if !seen[md] {
				priceMap[md] += allPrices[i][j+4]
			}
			seen[md] = true
		}
	}

	part2 = slices.Max(slices.Collect(maps.Values(priceMap)))
	return part1, part2
}

func nextSecretNumber(secretNumber int) int {
	nextNumber := (secretNumber * 64) ^ secretNumber
	nextNumber %= 16777216

	nextNumber = (nextNumber / 32) ^ nextNumber
	nextNumber %= 16777216

	nextNumber = (nextNumber * 2048) ^ nextNumber

	return nextNumber % 16777216
}
