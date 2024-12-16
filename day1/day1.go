package day1

import (
	"slices"
	"strconv"
	"strings"

	"github.com/mikeramage/aoc2024/utils"
)

func Day1() (int, int) {
	lines := utils.Lines("./input/day1.txt")
	var left, right []int

	for _, line := range lines {
		input := strings.Fields(line)
		l, _ := strconv.Atoi(input[0])
		r, _ := strconv.Atoi(input[1])
		left = append(left, l)
		right = append(right, r)
	}

	slices.Sort(left)
	slices.Sort(right)

	part1 := 0
	m := make(map[int]int)
	for i := 0; i < len(left); i++ {
		diff := utils.Abs(left[i] - right[i])
		part1 += diff
		m[right[i]]++
	}

	part2 := 0
	for _, l := range left {
		part2 += m[l] * l
	}

	return part1, part2
}
