package day11

import (
	"maps"
	"strconv"
	"strings"

	"github.com/mikeramage/aoc2024/utils"
)

func Day11() (int, int) {
	var part1, part2 int
	stones := parseStones()
	stones = blink(stones, 25)

	for stone := range maps.Keys(stones) {
		part1 += stones[stone]
	}

	stones = blink(stones, 50)
	for stone := range maps.Keys(stones) {
		part2 += stones[stone]
	}

	return part1, part2
}

func blink(stones map[int]int, blinks int) map[int]int {
	for b := 0; b < blinks; b++ {
		newStones := make(map[int]int)
		for stone := range maps.Keys(stones) {
			occurrences := stones[stone]
			stoneAsStr := strconv.Itoa(stone)
			if stone == 0 {
				newStones[1] += occurrences
			} else if len(stoneAsStr)%2 == 0 {
				stone1, _ := strconv.Atoi(stoneAsStr[:len(stoneAsStr)/2])
				stone2, _ := strconv.Atoi(stoneAsStr[len(stoneAsStr)/2:])
				newStones[stone1] += occurrences
				newStones[stone2] += occurrences
			} else {
				newStones[stone*2024] += occurrences
			}
		}
		stones = newStones
	}

	return stones
}

func parseStones() map[int]int {
	input := utils.Lines("./input/day11.txt")
	stonesStr := strings.Fields(input[0])

	stones := make(map[int]int)
	for _, stone := range stonesStr {
		stoneInt, _ := strconv.Atoi(stone)
		stones[stoneInt]++
	}
	return stones
}
