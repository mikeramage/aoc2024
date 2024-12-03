package internal

import (
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var reMul = regexp.MustCompile(`mul\((\d+),(\d+)\)`)
var reDoDont = regexp.MustCompile(`do\(\)(.*?)(don't\(\)|\z)`)

const do = "do()"

func Day3() (int, int) {
	input := Lines("./input/day3.txt")
	var part1 int
	var part2 int
	var doArr []string
	doArr = append(doArr, do)
	input = slices.Concat(doArr, input)
	allInput := strings.Join(input, "")
	part1 += doMultiply(allInput)

	matches := reDoDont.FindAllStringSubmatch(allInput, -1)
	for _, match := range matches {
		part2 += doMultiply(match[1])
	}

	// Yeah, I really should immediately have used regexes for part2. The following does all work assuming you
	// don't prepend "do()" onto the beginning of the string as I've done above, but it's a faff.
	// Not at all sure why I went down this rabbit hole. It took about 10 times as long. Keeping for posterity
	//
	// var doIndices []int
	// doIndices = append(doIndices, 0)
	// doIndices = slices.Concat(doIndices, indexAll(allInput, do))
	// dontIndices := indexAll(allInput, dont)
	// currentDo := 0
	// currentDont := 0
	// for currentDo < len(doIndices) && currentDont < len(dontIndices) {
	// 	// Find the next don't index that's less than the current do
	// 	for currentDont < len(dontIndices) && dontIndices[currentDont] < doIndices[currentDo] {
	// 		currentDont++
	// 	}

	// 	if currentDont < len(dontIndices) {
	// 		part2 += doMultiply(allInput[doIndices[currentDo]:dontIndices[currentDont]])
	// 	} else {
	// 		part2 += doMultiply(allInput[doIndices[currentDo]:])
	// 		break //Reached end of input
	// 	}

	// 	// Find the next do index that's greater than the current don't
	// 	for currentDo < len(doIndices) && doIndices[currentDo] < dontIndices[currentDont] {
	// 		currentDo++
	// 	}
	// }

	return part1, part2
}

func doMultiply(s string) int {
	var ans int
	matches := reMul.FindAllStringSubmatch(s, -1)
	for _, match := range matches {
		x, _ := strconv.Atoi(match[1])
		y, _ := strconv.Atoi(match[2])
		ans += x * y
	}
	return ans
}

// func indexAll(s string, substr string) []int {
// 	var indices []int
// 	index := strings.Index(s, substr)
// 	var cumIndex int
// 	for index != -1 {
// 		cumIndex += index
// 		indices = append(indices, cumIndex)
// 		cumIndex += len(substr)
// 		if cumIndex >= len(s) {
// 			break
// 		}
// 		index = strings.Index(s[cumIndex:], substr)
// 	}

// 	return indices
// }
